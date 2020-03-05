package binding

import "context"

// Message is the interface to a binding-specific message containing an event.
//
// Reliable Delivery
//
// There are 3 reliable qualities of service for messages:
//
// 0/at-most-once/unreliable: messages can be dropped silently.
//
// 1/at-least-once: messages are not dropped without signaling an error
// to the sender, but they may be duplicated in the event of a re-send.
//
// 2/exactly-once: messages are never dropped (without error) or
// duplicated, as long as both sending and receiving ends maintain
// some binding-specific delivery state. Whether this is persisted
// depends on the configuration of the binding implementations.
//
// The Message interface supports QoS 0 and 1, the ExactlyOnceMessage interface
// supports QoS 2
//
// The Structured and Binary methods allows to perform an optimized encoding of a Message
// to a specific data structure. A Sender should try each method of interest and fall back to ToEvent() if none are supported.
// For encoding a message, look at binding.Encode function.
//
// Most binding.Message implementations can be encoded only one time, because the encoding process drain the message itself.
// In order to consume a message several times,
// the binding/buffering module provides several APIs to buffer the Message in order to visit it more times.
//
// When a Message can be forgotten by the entity who produced the message, Message.Finish() *must* be invoked.
type Message interface {
	// Return the type of the message Encoding.
	// The encoding should be preferably computed when the message is constructed.
	Encoding() Encoding

	// Structured transfers a structured-mode event to a StructuredEncoder.
	// It must return ErrNotStructured if message is not in structured mode.
	//
	// Returns a different err if something wrong happened while trying to read the structured event.
	// In this case, the caller must Finish the message with appropriate error.
	//
	// This allows Senders to avoid re-encoding messages that are
	// already in suitable structured form.
	Structured(context.Context, StructuredEncoder) error

	// Binary transfers a binary-mode event to an BinaryEncoder.
	// It must return ErrNotBinary if message is not in binary mode.
	//
	// Returns a different err if something wrong happened while trying to read the binary event
	// In this case, the caller must Finish the message with appropriate error
	//
	// This allows Senders to avoid re-encoding messages that are
	// already in suitable binary form.
	Binary(context.Context, BinaryEncoder) error

	// Finish *must* be called when message from a Receiver can be forgotten by
	// the receiver. A QoS 1 sender should not call Finish() until it gets an acknowledgment of
	// receipt on the underlying transport.  For QoS 2 see ExactlyOnceMessage.
	//
	// Note that, depending on the Message implementation, forgetting to Finish the message
	// could produce memory/resources leaks!
	//
	// Passing a non-nil err indicates sending or processing failed.
	// A non-nil return indicates that the message was not accepted
	// by the receivers peer.
	Finish(error) error
}

// ExactlyOnceMessage is implemented by received Messages
// that support QoS 2.  Only transports that support QoS 2 need to
// implement or use this interface.
type ExactlyOnceMessage interface {
	Message

	// Received is called by a forwarding QoS2 Sender when it gets
	// acknowledgment of receipt (e.g. AMQP 'accept' or MQTT PUBREC)
	//
	// The receiver must call settle(nil) when it get's the ack-of-ack
	// (e.g. AMQP 'settle' or MQTT PUBCOMP) or settle(err) if the
	// transfer fails.
	//
	// Finally the Sender calls Finish() to indicate the message can be
	// discarded.
	//
	// If sending fails, or if the sender does not support QoS 2, then
	// Finish() may be called without any call to Received()
	Received(settle func(error))
}

// Message Wrapper interface is used to walk through a decorated Message and unwrap it.
type MessageWrapper interface {
	Message

	// Method to get the wrapped message
	GetWrappedMessage() Message
}
