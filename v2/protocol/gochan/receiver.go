package gochan

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudevents/sdk-go/v2/binding"
)

// ChanReceiver implements Receiver by receiving from a channel.
type ChanReceiver <-chan binding.Message

func (r ChanReceiver) Receive(ctx context.Context) (binding.Message, error) {
	if ctx == nil {
		return nil, fmt.Errorf("nil Context")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case m, ok := <-r:
		if !ok {
			return nil, io.EOF
		}
		return m, nil
	}
}
