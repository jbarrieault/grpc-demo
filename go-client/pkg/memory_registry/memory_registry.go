package memory_registry

import (
	"fmt"
	"sync"
)

type MemoryRegistry struct {
	services map[string]*Service
	mu       sync.Mutex
}

func NewRegistery() *MemoryRegistry {
	r := &MemoryRegistry{}
	r.services = make(map[string]*Service)
	return r
}

type Service struct {
	Name      string
	Addresses []Address
	r         *MemoryRegistry
	mu        sync.Mutex
}

type Address string

func (r *MemoryRegistry) newService(name string, addrs ...string) *Service {
	addresses := make([]Address, len(addrs))
	for i, addr := range addrs {
		addresses[i] = Address(addr)
	}

	s := Service{r: r, Name: name, Addresses: addresses}

	return &s
}

func (r *MemoryRegistry) Register(name string, addrs ...string) error {
	s := r.newService(name, addrs...)
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.services[s.Name]
	if !ok {
		r.services[s.Name] = s
	} else {
		return fmt.Errorf("%s is already registered", s.Name)
	}

	return nil
}

func (r *MemoryRegistry) Deregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.services, name)
}

func (r *MemoryRegistry) GetService(name string) (*Service, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var s *Service
	s, ok := r.services[name]
	if !ok {
		return s, fmt.Errorf("%s is not registered", name)
	}

	return s, nil
}

func (s *Service) AddAddress(addr string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.Addresses {
		if string(a) == addr {
			return fmt.Errorf("%s is already registered", addr)
		}
	}

	s.Addresses = append(s.Addresses, Address(addr))

	return nil
}

func (s *Service) RemoveAddress(addr string) error {
	address := Address(addr)
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, a := range s.Addresses {
		if a == address {
			s.Addresses = append(s.Addresses[:i], s.Addresses[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("%s is not registered", addr)
}

// TODO: add a way to interact with the service using a unix socket
