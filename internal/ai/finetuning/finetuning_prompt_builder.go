package finetuning

import (
	"fmt"
	"strings"
)

// FinetuningPromptBuilder builds prompts for fine-tuning training
type FinetuningPromptBuilder struct {
	systemPrompt string
}

// NewFinetuningPromptBuilder creates a new fine-tuning prompt builder
func NewFinetuningPromptBuilder(systemPrompt string) *FinetuningPromptBuilder {
	return &FinetuningPromptBuilder{
		systemPrompt: systemPrompt,
	}
}

// BuildTrainingPrompt builds a prompt for training data generation
func (fpb *FinetuningPromptBuilder) BuildTrainingPrompt(input string, context map[string]interface{}) string {
	var parts []string

	// Add system prompt if available
	if fpb.systemPrompt != "" {
		parts = append(parts, fmt.Sprintf("System: %s", fpb.systemPrompt))
	}

	// Add context if provided
	if len(context) > 0 {
		contextStr := fpb.formatContext(context)
		parts = append(parts, fmt.Sprintf("Context: %s", contextStr))
	}

	// Add input
	parts = append(parts, fmt.Sprintf("Input: %s", input))

	return strings.Join(parts, "\n\n")
}

// BuildCompletionPrompt builds a completion prompt for training
func (fpb *FinetuningPromptBuilder) BuildCompletionPrompt(input string, expectedOutput string) string {
	var parts []string

	if fpb.systemPrompt != "" {
		parts = append(parts, fmt.Sprintf("System: %s", fpb.systemPrompt))
	}

	parts = append(parts, fmt.Sprintf("Input: %s", input))
	parts = append(parts, fmt.Sprintf("Output: %s", expectedOutput))

	return strings.Join(parts, "\n\n")
}

// BuildInstructionPrompt builds an instruction-following prompt
func (fpb *FinetuningPromptBuilder) BuildInstructionPrompt(instruction string, input string) string {
	var parts []string

	if fpb.systemPrompt != "" {
		parts = append(parts, fmt.Sprintf("System: %s", fpb.systemPrompt))
	}

	parts = append(parts, fmt.Sprintf("Instruction: %s", instruction))
	parts = append(parts, fmt.Sprintf("Input: %s", input))

	return strings.Join(parts, "\n\n")
}

// formatContext formats context map as string
func (fpb *FinetuningPromptBuilder) formatContext(context map[string]interface{}) string {
	var parts []string
	for k, v := range context {
		parts = append(parts, fmt.Sprintf("%s: %v", k, v))
	}
	return strings.Join(parts, ", ")
}

// BuildDatasetEntry builds a dataset entry from prompt and response
func (fpb *FinetuningPromptBuilder) BuildDatasetEntry(prompt string, response string) string {
	return fmt.Sprintf(`{"prompt": %q, "completion": %q}`, prompt, response)
}
