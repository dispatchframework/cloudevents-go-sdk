package transcoder

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/spec"
	"github.com/cloudevents/sdk-go/pkg/binding/test"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

func TestVersionTranscoder(t *testing.T) {
	var testEventV02 = cloudevents.Event{
		Context: cloudevents.EventContextV02{
			Source:      types.URLRef{URL: url.URL{Path: "source"}},
			ContentType: cloudevents.StringOfApplicationJSON(),
			ID:          "id",
			Type:        "type",
		}.AsV02(),
	}

	var testEventV1 = testEventV02
	testEventV1.Context = testEventV02.Context.AsV1()

	data := []byte("\"data\"")
	err := testEventV02.SetData(data)
	require.NoError(t, err)
	err = testEventV1.SetData(data)
	require.NoError(t, err)

	test.RunTranscoderTests(t, context.Background(), []test.TranscoderTestArgs{
		{
			name:         "V02 -> V1 with Mock Structured message",
			inputMessage: test.NewMockStructuredMessage(test.CopyEventContext(testEventV02)),
			wantEvent:    test.CopyEventContext(testEventV1),
			transformer:  binding.TransformerFactories{Version(spec.V1)},
		},
		{
			name:         "V02 -> V1 with Mock Binary message",
			inputMessage: test.NewMockBinaryMessage(test.CopyEventContext(testEventV02)),
			wantEvent:    test.CopyEventContext(testEventV1),
			transformer:  binding.TransformerFactories{Version(spec.V1)},
		},
		{
			name:         "V02 -> V1 with Event message",
			inputMessage: binding.EventMessage(test.CopyEventContext(testEventV02)),
			wantEvent:    test.CopyEventContext(testEventV1),
			transformer:  binding.TransformerFactories{Version(spec.V1)},
		},
	})
}
