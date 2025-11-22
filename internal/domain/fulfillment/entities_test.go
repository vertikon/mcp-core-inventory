package fulfillment

import "testing"

func TestInboundShipmentLifecycle(t *testing.T) {
	items := []Item{{SKU: "ABC", Quantity: 10}}
	shipment, err := NewInboundShipment("inb-1", "ref-1", "factory", "cd-sp", items)
	if err != nil {
		t.Fatalf("expected shipment creation, got error %v", err)
	}

	if err := shipment.Start(); err != nil {
		t.Fatalf("expected start success, got %v", err)
	}
	if shipment.Status != StatusInProgress {
		t.Fatalf("expected status IN_PROGRESS, got %s", shipment.Status)
	}

	if err := shipment.Complete(); err != nil {
		t.Fatalf("expected complete success, got %v", err)
	}
	if shipment.Status != StatusCompleted {
		t.Fatalf("expected status COMPLETED, got %s", shipment.Status)
	}
}

func TestFulfillmentOrderLifecycle(t *testing.T) {
	order, err := NewFulfillmentOrder("ful-1", "order-123", "customer", 1, []Item{{SKU: "ABC", Quantity: 1}})
	if err != nil {
		t.Fatalf("expected order creation, got %v", err)
	}

	if err := order.StartPicking(); err != nil {
		t.Fatalf("expected start picking success, got %v", err)
	}

	if err := order.MarkShipped(); err != nil {
		t.Fatalf("expected mark shipped success, got %v", err)
	}
	if order.Status != StatusCompleted {
		t.Fatalf("expected status COMPLETED, got %s", order.Status)
	}
}

func TestCycleCountTask(t *testing.T) {
	task, err := NewCycleCountTask("task-1", "cd-sp", map[string]int{"SKU1": 10})
	if err != nil {
		t.Fatalf("expected task creation, got %v", err)
	}

	if err := task.Start(); err != nil {
		t.Fatalf("expected start success, got %v", err)
	}

	results := map[string]int{"SKU1": 10}
	if err := task.SubmitResults("auditor", results); err != nil {
		t.Fatalf("expected submit success, got %v", err)
	}

	if task.Status != StatusCompleted {
		t.Fatalf("expected status COMPLETED, got %s", task.Status)
	}
	if task.CountedBy != "auditor" {
		t.Fatalf("expected CountedBy auditor, got %s", task.CountedBy)
	}
}
