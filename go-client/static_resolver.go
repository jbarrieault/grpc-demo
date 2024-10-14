package main

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

type staticResolver struct {
	cc        resolver.ClientConn
	addresses []resolver.Address
}

func (r *staticResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *staticResolver) Close() {}

type staticResolverBuilder struct {
	addresses []resolver.Address
}

func (b *staticResolverBuilder) Scheme() string {
	return "static"
}

func (b *staticResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &staticResolver{cc, b.addresses}
	r.cc.UpdateState(resolver.State{Addresses: b.addresses})

	return r, nil
}

func registerStaticResolver() {
	addresses := strings.Split(*addr, ",")
	var addrs []resolver.Address
	for _, address := range addresses {
		addrs = append(addrs, resolver.Address{Addr: address})
	}

	srb := &staticResolverBuilder{addresses: addrs}
	resolver.Register(srb)
}
