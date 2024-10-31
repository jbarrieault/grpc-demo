package main

import (
	mr "github.com/jbarrieault/grpc-demo/memory-registry"
	"google.golang.org/grpc/resolver"
)

type registryResolver struct {
	service  string
	cc       resolver.ClientConn
	registry *mr.MemoryRegistry
}

func (r *registryResolver) Close() {}

func (r *registryResolver) ResolveNow(resolver.ResolveNowOptions) {
	r.cc.UpdateState(resolver.State{Addresses: r.addressesFromRegistry()})
}

func (r *registryResolver) addressesFromRegistry() []resolver.Address {
	s, err := r.registry.GetService(r.service)
	if err != nil {
		return []resolver.Address{}
	}

	addrs := s.Addresses()
	addresses := make([]resolver.Address, len(addrs))
	for i, addr := range addrs {
		addresses[i] = resolver.Address{Addr: addr}
	}

	return addresses
}

type registryResolverBuilder struct {
	registry *mr.MemoryRegistry
}

func (b *registryResolverBuilder) Scheme() string {
	return "registry"
}

func (b *registryResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &registryResolver{
		service:  target.Endpoint(),
		cc:       cc,
		registry: b.registry,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})

	return r, nil
}

func registerRegistryResolver() {
	rrb := &registryResolverBuilder{mem_reg}
	resolver.Register(rrb)
}
