package transcoder

import (
	"testing"

	"github.com/stretchr/testify/require"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/buffering"
	"github.com/cloudevents/sdk-go/pkg/binding/test"
)

type TranscoderTestArgs struct {
	name         string
	inputMessage binding.Message
	wantEvent    cloudevents.Event
	transformer  binding.TransformerFactory
}

func RunTranscoderTests(t *testing.T, tests []TranscoderTestArgs) {
	for _, tt := range tests {
		tt := tt // Don't use range variable inside scope
		t.Run(tt.name, func(t *testing.T) {
			copied, err := buffering.CopyMessage(tt.inputMessage, tt.transformer)
			require.NoError(t, err)
			e, _, err := binding.ToEvent(copied)
			require.NoError(t, err)
			test.AssertEventEquals(t, tt.wantEvent, e)
		})
	}
}
