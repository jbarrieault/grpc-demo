package memory_registry

import (
	"fmt"
	"sync"
)

type MemoryRegistry struct {
	services map[string]*Service
	mu       sync.RWMutex
}

func (r *MemoryRegistry) Services() map[string]*Service {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.services
}

func NewRegistery() *MemoryRegistry {
	r := &MemoryRegistry{}
	r.services = make(map[string]*Service)
	return r
}

type Service struct {
	Name      string
	addresses []string
	r         *MemoryRegistry
	mu        sync.RWMutex
}

func (r *MemoryRegistry) newService(name string, addresses ...string) *Service {
	return &Service{r: r, Name: name, addresses: addresses}
}

func (s *Service) Addresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := make([]string, len(s.addresses))
	copy(addresses, s.addresses)
	return addresses
}

func (r *MemoryRegistry) Register(name string, addrs ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s := r.newService(name, addrs...)

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

	for _, a := range s.addresses {
		if a == addr {
			return fmt.Errorf("%s is already registered", addr)
		}
	}

	s.addresses = append(s.addresses, addr)

	return nil
}

func (s *Service) RemoveAddress(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, a := range s.addresses {
		if a == address {
			s.addresses = append(s.addresses[:i], s.addresses[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("%s is not registered", address)
}
