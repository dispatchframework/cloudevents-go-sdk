package transcoder

import (
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/spec"
)

// TODO(slinkydeveloper) docs
func UpdateAttribute(attributeKind spec.Kind, updater func(interface{}) interface{}) binding.TransformerFactory {
	return updateAttributeTranscoderFactory{attributeKind: attributeKind, updater: updater}
}

type updateAttributeTranscoderFactory struct {
	attributeKind spec.Kind
	updater       func(interface{}) interface{}
}

func (a updateAttributeTranscoderFactory) StructuredTransformer(binding.StructuredEncoder) binding.StructuredEncoder {
	return nil
}

func (a updateAttributeTranscoderFactory) BinaryTransformer(encoder binding.BinaryEncoder) binding.BinaryEncoder {
	return &updateAttributeTransformer{
		BinaryEncoder: encoder,
		attributeKind: a.attributeKind,
		updater:       a.updater,
	}
}

func (a updateAttributeTranscoderFactory) EventTransformer() binding.EventTransformer {
	return func(event *cloudevents.Event) error {
		v, err := spec.VS.Version(event.SpecVersion())
		if err != nil {
			return err
		}
		if val := v.AttributeFromKind(a.attributeKind).Get(event.Context); val != nil {
			newVal := a.updater(val)
			if newVal == nil {
				return v.AttributeFromKind(a.attributeKind).Delete(event.Context)
			} else {
				return v.AttributeFromKind(a.attributeKind).Set(event.Context, newVal)
			}
		}
		return nil
	}
}

type updateAttributeTransformer struct {
	binding.BinaryEncoder
	attributeKind spec.Kind
	updater       func(interface{}) interface{}
}

func (b *updateAttributeTransformer) SetAttribute(attribute spec.Attribute, value interface{}) error {
	if attribute.Kind() == b.attributeKind {
		newVal := b.updater(value)
		if newVal != nil {
			return b.BinaryEncoder.SetAttribute(attribute, newVal)
		}
		return nil
	}
	return b.BinaryEncoder.SetAttribute(attribute, value)
}
