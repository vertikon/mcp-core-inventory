package fulfillment

import (
	"fmt"
	"time"
)

// Item representa uma unidade manipulada em qualquer processo físico.
type Item struct {
	SKU      string
	Quantity int
	Batch    string
}

func (i Item) validate() error {
	if i.SKU == "" {
		return fmt.Errorf("item sku is required")
	}
	if i.Quantity <= 0 {
		return fmt.Errorf("item quantity must be positive")
	}
	return nil
}

// InboundShipment descreve o recebimento físico de mercadorias.
type InboundShipment struct {
	ID          string
	ReferenceID string
	Origin      string
	Destination string
	Status      Status
	Items       []Item
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewInboundShipment cria um novo workflow de recebimento.
func NewInboundShipment(id, referenceID, origin, destination string, items []Item) (*InboundShipment, error) {
	if id == "" || destination == "" {
		return nil, ErrInvalidPayload
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("inbound shipment requires items")
	}
	for _, item := range items {
		if err := item.validate(); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	return &InboundShipment{
		ID:          id,
		ReferenceID: referenceID,
		Origin:      origin,
		Destination: destination,
		Status:      StatusPending,
		Items:       items,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Start marca o recebimento como em progresso.
func (s *InboundShipment) Start() error {
	if s.Status != StatusPending {
		return ErrInvalidStateTransition
	}
	s.Status = StatusInProgress
	s.UpdatedAt = time.Now()
	return nil
}

// Complete marca o recebimento como concluído.
func (s *InboundShipment) Complete() error {
	if s.Status != StatusInProgress {
		return ErrInvalidStateTransition
	}
	s.Status = StatusCompleted
	s.UpdatedAt = time.Now()
	return nil
}

// FulfillmentOrder descreve a vida de um pedido em outbound.
type FulfillmentOrder struct {
	ID             string
	OrderID        string
	Destination    string
	Priority       int
	Status         Status
	Items          []Item
	IdempotencyKey string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewFulfillmentOrder cria um pedido de fulfillment.
func NewFulfillmentOrder(id, orderID, destination string, priority int, items []Item) (*FulfillmentOrder, error) {
	if id == "" || orderID == "" || destination == "" {
		return nil, ErrInvalidPayload
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("fulfillment order requires items")
	}
	for _, item := range items {
		if err := item.validate(); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	return &FulfillmentOrder{
		ID:          id,
		OrderID:     orderID,
		Destination: destination,
		Priority:    priority,
		Status:      StatusPending,
		Items:       items,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// StartPicking inicia o picking de um pedido.
func (o *FulfillmentOrder) StartPicking() error {
	if o.Status != StatusPending {
		return ErrInvalidStateTransition
	}
	o.Status = StatusInProgress
	o.UpdatedAt = time.Now()
	return nil
}

// MarkShipped finaliza o pedido como expedido.
func (o *FulfillmentOrder) MarkShipped() error {
	if o.Status != StatusInProgress {
		return ErrInvalidStateTransition
	}
	o.Status = StatusCompleted
	o.UpdatedAt = time.Now()
	return nil
}

// TransferOrder controla transferências entre locais.
type TransferOrder struct {
	ID          string
	Origin      string
	Destination string
	Status      Status
	Items       []Item
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTransferOrder cria uma transferência.
func NewTransferOrder(id, origin, destination string, items []Item) (*TransferOrder, error) {
	if id == "" || origin == "" || destination == "" {
		return nil, ErrInvalidPayload
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("transfer order requires items")
	}
	for _, item := range items {
		if err := item.validate(); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	return &TransferOrder{
		ID:          id,
		Origin:      origin,
		Destination: destination,
		Status:      StatusPending,
		Items:       items,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Start marca a transferência como em movimento.
func (t *TransferOrder) Start() error {
	if t.Status != StatusPending {
		return ErrInvalidStateTransition
	}
	t.Status = StatusInProgress
	t.UpdatedAt = time.Now()
	return nil
}

// Complete finaliza a transferência.
func (t *TransferOrder) Complete() error {
	if t.Status != StatusInProgress {
		return ErrInvalidStateTransition
	}
	t.Status = StatusCompleted
	t.UpdatedAt = time.Now()
	return nil
}

// ReturnOrder descreve uma devolução.
type ReturnOrder struct {
	ID          string
	Source      string
	Reason      string
	Status      Status
	Disposition string
	Items       []Item
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewReturnOrder cria uma devolução.
func NewReturnOrder(id, source, reason string, items []Item) (*ReturnOrder, error) {
	if id == "" || source == "" {
		return nil, ErrInvalidPayload
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("return order requires items")
	}
	for _, item := range items {
		if err := item.validate(); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	return &ReturnOrder{
		ID:        id,
		Source:    source,
		Reason:    reason,
		Status:    StatusPending,
		Items:     items,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Complete finaliza a devolução com a destinação informada.
func (r *ReturnOrder) Complete(disposition string) error {
	if r.Status != StatusPending && r.Status != StatusInProgress {
		return ErrInvalidStateTransition
	}
	if disposition == "" {
		return fmt.Errorf("disposition is required")
	}
	r.Status = StatusCompleted
	r.Disposition = disposition
	r.UpdatedAt = time.Now()
	return nil
}

// CycleCountTask descreve uma contagem cíclica.
type CycleCountTask struct {
	ID          string
	Location    string
	Status      Status
	Expected    map[string]int
	Results     map[string]int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CountedBy   string
	CompletedAt *time.Time
}

// NewCycleCountTask cria uma nova contagem cíclica.
func NewCycleCountTask(id, location string, expected map[string]int) (*CycleCountTask, error) {
	if id == "" || location == "" {
		return nil, ErrInvalidPayload
	}
	if len(expected) == 0 {
		return nil, fmt.Errorf("cycle count requires expected items")
	}
	now := time.Now()
	return &CycleCountTask{
		ID:        id,
		Location:  location,
		Status:    StatusPending,
		Expected:  expected,
		Results:   make(map[string]int),
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Start inicia a contagem.
func (c *CycleCountTask) Start() error {
	if c.Status != StatusPending {
		return ErrInvalidStateTransition
	}
	c.Status = StatusInProgress
	c.UpdatedAt = time.Now()
	return nil
}

// SubmitResults registra o resultado contado.
func (c *CycleCountTask) SubmitResults(countedBy string, results map[string]int) error {
	if c.Status != StatusInProgress {
		return ErrInvalidStateTransition
	}
	if countedBy == "" {
		return fmt.Errorf("counted by is required")
	}
	c.CountedBy = countedBy
	c.Results = results
	now := time.Now()
	c.CompletedAt = &now
	c.Status = StatusCompleted
	c.UpdatedAt = now
	return nil
}
