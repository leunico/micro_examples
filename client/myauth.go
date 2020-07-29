package client

import (
	"context"

	myauth "myauth/proto/myauth"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
)

type myauthKey struct{}

// FromContext retrieves the client from the Context
func MyauthFromContext(ctx context.Context) (myauth.MyauthService, bool) {
	c, ok := ctx.Value(myauthKey{}).(myauth.MyauthService)
	return c, ok
}

// Client returns a wrapper for the MyauthClient
func MyauthWrapper(service micro.Service) server.HandlerWrapper {
	client := myauth.NewMyauthService("go.micro.service.template", service.Client())

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, myauthKey{}, client)
			return fn(ctx, req, rsp)
		}
	}
}
