package binding_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/format"
	ce "github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/transport"
)

// ExMessage is a json.RawMessage, a byte slice containing a JSON encoded event.
// It implements binding.MockStructuredMessage
//
// Note: a good binding implementation should provide an easy way to convert
// between the Message implementation and the "native" message format.
// In this case it's as simple as:
//
//    native = ExMessage(impl)
//    impl = json.RawMessage(native)
//
// For example in a HTTP binding it should be easy to convert between
// the HTTP binding.Message implementation and net/http.Request and
// Response types.  There are no interfaces for this conversion as it
// requires the use of unknown types.
type ExMessage json.RawMessage

func (m ExMessage) Structured(b binding.StructuredMessageBuilder) error {
	return b.Event(format.JSON, bytes.NewReader([]byte(m)))
}

func (m ExMessage) Binary(binding.BinaryMessageBuilder) error {
	return binding.ErrNotBinary
}

func (m ExMessage) Event(b binding.EventMessageBuilder) error {
	e := ce.Event{}
	err := json.Unmarshal(m, &e)
	if err != nil {
		return err
	}
	return b.Encode(e)
}

func (m ExMessage) Finish(error) error { return nil }

var _ binding.Message = (*ExMessage)(nil)

// ExSender sends by writing JSON encoded events to an io.Writer
// ExSender supports transcoding
// ExSender implements directly StructuredMessageBuilder & EventMessageBuilder
type ExSender struct {
	encoder      *json.Encoder
	transcodings binding.TranscoderFactories
}

func NewExSender(w io.Writer, factories ...binding.TranscoderFactory) binding.Sender {
	return &ExSender{encoder: json.NewEncoder(w), transcodings: factories}
}

func (s *ExSender) Send(ctx context.Context, m binding.Message) error {
	// Wrap the transcoders in the structured builder
	structuredBuilder := s.transcodings.StructuredMessageTranscoder(s)

	// StructuredMessageTranscoder could return nil if one of transcoders doesn't support
	// direct structured transcoding
	if structuredBuilder != nil {
		// Fast case: Let's try to build in structured mode
		if err := m.Structured(structuredBuilder); err == nil {
			return nil
		} else if err != binding.ErrNotStructured {
			return err
		}
	}

	// Some other message encoding. Decode as generic Event and re-encode.
	eventBuilder := s.transcodings.EventMessageTranscoder(s)
	return m.Event(eventBuilder)
}

func (s *ExSender) Event(f format.Format, event io.Reader) error {
	if f == format.JSON {
		b, err := ioutil.ReadAll(event)
		if err != nil {
			return err
		}
		return s.encoder.Encode(json.RawMessage(b))
	} else {
		return binding.ErrNotStructured
	}
}

func (s *ExSender) Encode(event ce.Event) error {
	return s.encoder.Encode(&event)
}

func (s *ExSender) Close(context.Context) error { return nil }

var _ binding.Sender = (*ExSender)(nil)
var _ binding.StructuredMessageBuilder = (*ExSender)(nil)
var _ binding.EventMessageBuilder = (*ExSender)(nil)

// ExReceiver receives by reading JSON encoded events from an io.Reader
type ExReceiver struct{ decoder *json.Decoder }

func NewExReceiver(r io.Reader) binding.Receiver { return &ExReceiver{json.NewDecoder(r)} }

func (r *ExReceiver) Receive(context.Context) (binding.Message, error) {
	var rm json.RawMessage
	err := r.decoder.Decode(&rm) // This is just a byte copy.
	return ExMessage(rm), err
}
func (r *ExReceiver) Close(context.Context) error { return nil }

// NewExTransport returns a transport.Transport which is implemented by
// an ExSender and an ExReceiver
func NewExTransport(r io.Reader, w io.Writer) transport.Transport {
	return binding.NewTransport(NewExSender(w), NewExReceiver(r))
}

// Example of implementing a transport including a simple message type,
// and a transport sender and receiver.
func Example_implementing() {}
