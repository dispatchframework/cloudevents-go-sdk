package context

import (
	"context"
)

// Opaque key type used to store port
type portKeyType struct{}

var portKey = portKeyType{}

func ContextWithPort(ctx context.Context, port int) context.Context {
	return context.WithValue(ctx, portKey, port)
}

func PortFromContext(ctx context.Context) int {
	c := ctx.Value(portKey)
	if c != nil {
		if port, ok := c.(int); ok {
			return port
		}
	}
	return 0
}
