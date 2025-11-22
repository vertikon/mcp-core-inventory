package fulfillment

import "errors"

// OperationType representa o tipo de operação física processada pelo MCP.
type OperationType string

const (
	OpInbound    OperationType = "INBOUND"
	OpOutbound   OperationType = "OUTBOUND"
	OpTransfer   OperationType = "TRANSFER"
	OpReturn     OperationType = "RETURN"
	OpCycleCount OperationType = "CYCLE_COUNT"
)

// Status representa o estado atual de um workflow operacional.
type Status string

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
	StatusFailed     Status = "FAILED"
	StatusCancelled  Status = "CANCELLED"
)

var (
	// ErrInvalidStateTransition indica que uma transição de estado não é permitida.
	ErrInvalidStateTransition = errors.New("invalid state transition")
	// ErrInvalidPayload indica que a entidade está inconsistente.
	ErrInvalidPayload = errors.New("invalid payload")
)
