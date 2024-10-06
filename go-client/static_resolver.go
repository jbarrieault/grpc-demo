package main

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

type staticResolver struct {
	addresses []resolver.Address
	conn      resolver.ClientConn
	target    resolver.Target
}

func (r *staticResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *staticResolver) start() {
	r.conn.UpdateState(resolver.State{Addresses: r.addresses})
}

func (r *staticResolver) Close() {}

type staticResolverBuilder struct {
	addresses []resolver.Address
}

func (b *staticResolverBuilder) Scheme() string {
	return "static"
}

func (b *staticResolverBuilder) Build(target resolver.Target, conn resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &staticResolver{
		target:    target,
		conn:      conn,
		addresses: b.addresses,
	}

	r.start()
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
