package memory_registry

import "testing"

func TestMemoryRegistry(t *testing.T) {
	r := NewRegistery()
	r.Register("test_service", "localhost:9999", "localhost:3333")
	r.Register("another_service", "localhost:3232")

	if len(r.services) != 2 {
		t.Errorf("expected registry to have 2 services, has '%d'", len(r.services))
	}

	s, err := r.GetService("test_service")
	if err != nil {
		t.Errorf("expected GetService not to error, got '%s'", err)
	}

	if s.Name != "test_service" {
		t.Errorf("expected name 'test_service', got '%s'", s.Name)
	}

	if len(s.addresses) != 2 {
		t.Errorf("expected service to have 2 addresses, got %d", len(s.addresses))
	}

	addr := s.addresses[0]
	if string(addr) != "localhost:9999" {
		t.Errorf("expected first address 'localhost:9999', got '%s'", addr)
	}

	addr = s.addresses[1]
	if string(addr) != "localhost:3333" {
		t.Errorf("expected second address 'localhost:3333', got '%s'", addr)
	}

	err = s.AddAddress("localhost:3005")
	if err != nil {
		t.Errorf("expected AddAddress not to error, got '%s'", err)
	}

	if len(s.addresses) != 3 {
		t.Errorf("expected service to have 3 addresses, got %d", len(s.addresses))
	}

	addr = s.addresses[2]
	if string(addr) != "localhost:3005" {
		t.Errorf("expected third address 'localhost:3005', got '%s'", addr)
	}

	err = s.RemoveAddress("localhost:3333")
	if err != nil {
		t.Errorf("expected RemoveAddress not to error, got '%s'", err)
	}

	if len(s.addresses) != 2 {
		t.Errorf("expected service to have 2 addresses, got %d", len(s.addresses))
	}

	addr = s.addresses[0]
	if string(addr) != "localhost:9999" {
		t.Errorf("expected first address 'localhost:9999', got '%s'", addr)
	}

	addr = s.addresses[1]
	if string(addr) != "localhost:3005" {
		t.Errorf("expected second address 'localhost:3005', got '%s'", addr)
	}
}
