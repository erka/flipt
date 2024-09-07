package auth

import (
	"context"

	"go.flipt.io/flipt/rpc/flipt"
)

func (req *CreateTokenRequest) Request(context.Context) flipt.Request {
	return flipt.NewRequest(flipt.ResourceAuthentication, flipt.ActionCreate, flipt.WithSubject(flipt.SubjectToken))
}

func (req *ListAuthenticationsRequest) Request(context.Context) flipt.Request {
	return flipt.NewRequest(flipt.ResourceAuthentication, flipt.ActionRead)
}

func (req *GetAuthenticationRequest) Request(context.Context) flipt.Request {
	return flipt.NewRequest(flipt.ResourceAuthentication, flipt.ActionRead)
}

func (req *DeleteAuthenticationRequest) Request(context.Context) flipt.Request {
	return flipt.NewRequest(flipt.ResourceAuthentication, flipt.ActionDelete)
}
