package transcoder

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/spec"
	"github.com/cloudevents/sdk-go/pkg/binding/test"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
)

func TestAddAttribute(t *testing.T) {
	e := test.MinEvent()
	e.Context = e.Context.AsV1()

	subject := "aaa"
	expectedEventWithSubject := copyEventContext(e)
	require.NoError(t, expectedEventWithSubject.Context.SetSubject(subject))

	timestamp, err := types.ToTime(time.Now())
	require.NoError(t, err)
	expectedEventWithTime := copyEventContext(e)
	require.NoError(t, expectedEventWithTime.Context.SetTime(timestamp))

	RunTranscoderTests(t, []TranscoderTestArgs{
		{
			name:         "No change to id to Mock Structured message",
			inputMessage: test.NewMockStructuredMessage(copyEventContext(e)),
			wantEvent:    copyEventContext(e),
			transformer:  AddAttribute(spec.ID, "new-id"),
		},
		{
			name:         "No change to id to Mock Binary message",
			inputMessage: test.NewMockBinaryMessage(copyEventContext(e)),
			wantEvent:    copyEventContext(e),
			transformer:  AddAttribute(spec.ID, "new-id"),
		},
		{
			name:         "No change to id to Event message",
			inputMessage: binding.EventMessage(copyEventContext(e)),
			wantEvent:    copyEventContext(e),
			transformer:  AddAttribute(spec.ID, "new-id"),
		},
		{
			name:         "Add subject to Mock Structured message",
			inputMessage: test.NewMockStructuredMessage(copyEventContext(e)),
			wantEvent:    expectedEventWithSubject,
			transformer:  AddAttribute(spec.Subject, subject),
		},
		{
			name:         "Add subject to Mock Binary message",
			inputMessage: test.NewMockBinaryMessage(copyEventContext(e)),
			wantEvent:    expectedEventWithSubject,
			transformer:  AddAttribute(spec.Subject, subject),
		},
		{
			name:         "Add subject to Event message",
			inputMessage: binding.EventMessage(copyEventContext(e)),
			wantEvent:    expectedEventWithSubject,
			transformer:  AddAttribute(spec.Subject, subject),
		},
		{
			name:         "Add time to Mock Structured message",
			inputMessage: test.NewMockStructuredMessage(copyEventContext(e)),
			wantEvent:    expectedEventWithTime,
			transformer:  AddAttribute(spec.Time, timestamp),
		},
		{
			name:         "Add time to Mock Binary message",
			inputMessage: test.NewMockBinaryMessage(copyEventContext(e)),
			wantEvent:    expectedEventWithTime,
			transformer:  AddAttribute(spec.Time, timestamp),
		},
		{
			name:         "Add time to Event message",
			inputMessage: binding.EventMessage(copyEventContext(e)),
			wantEvent:    expectedEventWithTime,
			transformer:  AddAttribute(spec.Time, timestamp),
		},
	})
}
