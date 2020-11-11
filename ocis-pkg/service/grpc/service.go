package grpc

import (
	"fmt"
	"strings"
	"time"

	"github.com/micro/go-micro/v2"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"

	etcdr "github.com/micro/go-micro/v2/registry/etcd"

	"github.com/micro/go-plugins/wrapper/trace/opencensus/v2"
	"github.com/owncloud/ocis/ocis-pkg/wrapper/prometheus"
)

var DefaultClient = newGrpcClient()

func newGrpcClient() mclient.Client {
	c := grpc.NewClient(
		mclient.RequestTimeout(10*time.Second),
		mclient.Registry(etcdr.NewRegistry()), // this is a workaround and will force clients to ONLY use etcd as the registry. This needs to be configurable
	)
	return c
}

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewService initializes a new grpc service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	fmt.Printf("\n\n%v\n\n", sopts.Name)

	sopts.Logger.Info().
		Str("transport", "grpc").
		Str("addr", sopts.Address).
		Msg("starting server")

	mopts := []micro.Option{
		micro.Name(
			strings.Join(
				[]string{
					sopts.Namespace,
					sopts.Name,
				},
				".",
			),
		),
		micro.Client(DefaultClient),
		micro.Version(sopts.Version),
		micro.Address(sopts.Address),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
	}

	return Service{
		micro.NewService(
			mopts...,
		),
	}
}
