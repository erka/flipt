package ofrep

import (
	"context"

	"go.flipt.io/flipt/rpc/flipt"
	"google.golang.org/grpc/metadata"
)

func (x *GetProviderConfigurationRequest) Request(ctx context.Context) flipt.Request {
	ns := GetNamespace(ctx)
	return flipt.NewRequest(flipt.ResourceOFREP, flipt.ActionRead, flipt.WithNamespace(ns))
}

func (x *EvaluateFlagRequest) Request(ctx context.Context) flipt.Request {
	ns := GetNamespace(ctx)
	return flipt.NewRequest(flipt.ResourceOFREP, flipt.ActionEvaluate, flipt.WithNamespace(ns))
}

func (x *EvaluateBulkRequest) Request(ctx context.Context) flipt.Request {
	ns := GetNamespace(ctx)
	return flipt.NewRequest(flipt.ResourceOFREP, flipt.ActionEvaluate, flipt.WithNamespace(ns))
}

func GetNamespace(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "default"
	}

	namespace := md.Get("x-flipt-namespace")
	if len(namespace) == 0 {
		return "default"
	}

	return namespace[0]
}
