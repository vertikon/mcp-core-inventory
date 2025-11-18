// Package transformer provides inference engine for GLM-4.6
package transformer

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
)

// InferenceConfig represents inference configuration
type InferenceConfig struct {
	MaxTokens         int     `json:"max_tokens"`
	Temperature       float64 `json:"temperature"`
	TopK              int     `json:"top_k"`
	TopP              float64 `json:"top_p"`
	BeamWidth         int     `json:"beam_width"`
	UseBeamSearch     bool    `json:"use_beam_search"`
	RepetitionPenalty float64 `json:"repetition_penalty"`
}

// InferenceResult represents inference result
type InferenceResult struct {
	Tokens       []int     `json:"tokens"`
	LogProbs     []float64 `json:"log_probs"`
	FinishReason string    `json:"finish_reason"`
}

// BeamSearchCandidate represents a candidate in beam search
type BeamSearchCandidate struct {
	tokens   []int
	logProb  float64
	finished bool
}

// InferenceEngine provides inference capabilities
type InferenceEngine struct {
	transformer *GLMTransformer
	config      InferenceConfig
	mu          sync.RWMutex
}

// NewInferenceEngine creates a new inference engine
func NewInferenceEngine(transformer *GLMTransformer, config InferenceConfig) *InferenceEngine {
	return &InferenceEngine{
		transformer: transformer,
		config:      config,
	}
}

// Generate generates tokens using the transformer
func (ie *InferenceEngine) Generate(ctx context.Context, input *Tensor, config InferenceConfig) (*InferenceResult, error) {
	if config.UseBeamSearch {
		return ie.beamSearch(ctx, input, config)
	}
	return ie.sample(ctx, input, config)
}

// beamSearch performs beam search generation
func (ie *InferenceEngine) beamSearch(ctx context.Context, input *Tensor, config InferenceConfig) (*InferenceResult, error) {
	beamWidth := config.BeamWidth
	if beamWidth == 0 {
		beamWidth = 4 // Default
	}

	// Initialize beam with start token
	beams := []*BeamSearchCandidate{
		{
			tokens:   []int{0}, // Start token
			logProb:  0.0,
			finished: false,
		},
	}

	for step := 0; step < config.MaxTokens; step++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Expand all beams
		candidates := make([]*BeamSearchCandidate, 0)
		for _, beam := range beams {
			if beam.finished {
				candidates = append(candidates, beam)
				continue
			}

			// Get next token probabilities
			probs, err := ie.getNextTokenProbs(ctx, input, beam.tokens)
			if err != nil {
				return nil, fmt.Errorf("beam search error: %w", err)
			}

			// Get top-k candidates
			topK := config.TopK
			if topK == 0 {
				topK = 50 // Default
			}
			topTokens := ie.topK(probs, topK)

			// Create candidates
			for _, token := range topTokens {
				newTokens := append([]int{}, beam.tokens...)
				newTokens = append(newTokens, token.idx)
				candidates = append(candidates, &BeamSearchCandidate{
					tokens:   newTokens,
					logProb:  beam.logProb + token.logProb,
					finished: token.idx == 1, // End token
				})
			}
		}

		// Select top beamWidth beams
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].logProb > candidates[j].logProb
		})
		beams = candidates[:min(beamWidth, len(candidates))]

		// Check if all beams are finished
		allFinished := true
		for _, beam := range beams {
			if !beam.finished {
				allFinished = false
				break
			}
		}
		if allFinished {
			break
		}
	}

	// Return best beam
	if len(beams) == 0 {
		return nil, fmt.Errorf("no beams generated")
	}

	bestBeam := beams[0]
	return &InferenceResult{
		Tokens:   bestBeam.tokens,
		LogProbs: []float64{bestBeam.logProb},
		FinishReason: func() string {
			if len(bestBeam.tokens) >= config.MaxTokens {
				return "length"
			}
			return "stop"
		}(),
	}, nil
}

// sample performs sampling generation
func (ie *InferenceEngine) sample(ctx context.Context, input *Tensor, config InferenceConfig) (*InferenceResult, error) {
	tokens := []int{0} // Start token
	logProbs := []float64{0.0}

	for step := 0; step < config.MaxTokens; step++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Get next token probabilities
		probs, err := ie.getNextTokenProbs(ctx, input, tokens)
		if err != nil {
			return nil, fmt.Errorf("sampling error: %w", err)
		}

		// Apply temperature
		if config.Temperature > 0 {
			probs = ie.applyTemperature(probs, config.Temperature)
		}

		// Apply top-k or top-p
		if config.TopK > 0 {
			probs = ie.applyTopK(probs, config.TopK)
		}
		if config.TopP > 0 {
			probs = ie.applyTopP(probs, config.TopP)
		}

		// Sample token
		tokenIdx := ie.sampleToken(probs)
		tokens = append(tokens, tokenIdx)
		logProbs = append(logProbs, math.Log(probs[tokenIdx]))

		// Check for end token
		if tokenIdx == 1 { // End token
			break
		}
	}

	return &InferenceResult{
		Tokens:   tokens,
		LogProbs: logProbs,
		FinishReason: func() string {
			if len(tokens) >= config.MaxTokens {
				return "length"
			}
			return "stop"
		}(),
	}, nil
}

// getNextTokenProbs gets probabilities for next token
func (ie *InferenceEngine) getNextTokenProbs(ctx context.Context, input *Tensor, tokens []int) ([]float64, error) {
	// Simplified: return uniform distribution
	// In real implementation, this would use the transformer
	vocabSize := 10000 // Default vocab size
	probs := make([]float64, vocabSize)
	for i := range probs {
		probs[i] = 1.0 / float64(vocabSize)
	}
	return probs, nil
}

// topK returns top-k tokens with their probabilities
type tokenProb struct {
	idx     int
	logProb float64
}

func (ie *InferenceEngine) topK(probs []float64, k int) []tokenProb {
	tokens := make([]tokenProb, len(probs))
	for i, prob := range probs {
		tokens[i] = tokenProb{
			idx:     i,
			logProb: math.Log(prob + 1e-10),
		}
	}

	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].logProb > tokens[j].logProb
	})

	return tokens[:min(k, len(tokens))]
}

// applyTemperature applies temperature scaling
func (ie *InferenceEngine) applyTemperature(probs []float64, temp float64) []float64 {
	result := make([]float64, len(probs))
	sum := 0.0
	for i, prob := range probs {
		result[i] = math.Pow(prob, 1.0/temp)
		sum += result[i]
	}
	for i := range result {
		result[i] /= sum
	}
	return result
}

// applyTopK applies top-k filtering
func (ie *InferenceEngine) applyTopK(probs []float64, k int) []float64 {
	result := make([]float64, len(probs))
	copy(result, probs)

	// Get top-k indices
	indices := make([]int, len(probs))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return probs[indices[i]] > probs[indices[j]]
	})

	// Zero out non-top-k
	for i := k; i < len(indices); i++ {
		result[indices[i]] = 0.0
	}

	// Renormalize
	sum := 0.0
	for _, prob := range result {
		sum += prob
	}
	if sum > 0 {
		for i := range result {
			result[i] /= sum
		}
	}

	return result
}

// applyTopP applies nucleus (top-p) filtering
func (ie *InferenceEngine) applyTopP(probs []float64, p float64) []float64 {
	result := make([]float64, len(probs))
	copy(result, probs)

	// Sort by probability
	indices := make([]int, len(probs))
	for i := range indices {
		indices[i] = i
	}
	sort.Slice(indices, func(i, j int) bool {
		return probs[indices[i]] > probs[indices[j]]
	})

	// Find cumulative sum threshold
	cumSum := 0.0
	cutoff := 0
	for i, idx := range indices {
		cumSum += probs[idx]
		if cumSum >= p {
			cutoff = i + 1
			break
		}
	}

	// Zero out below threshold
	for i := cutoff; i < len(indices); i++ {
		result[indices[i]] = 0.0
	}

	// Renormalize
	sum := 0.0
	for _, prob := range result {
		sum += prob
	}
	if sum > 0 {
		for i := range result {
			result[i] /= sum
		}
	}

	return result
}

// sampleToken samples a token from probability distribution
func (ie *InferenceEngine) sampleToken(probs []float64) int {
	r := rand.Float64()
	cumSum := 0.0
	for i, prob := range probs {
		cumSum += prob
		if r <= cumSum {
			return i
		}
	}
	return len(probs) - 1
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
