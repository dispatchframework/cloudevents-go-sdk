package event_test

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/cloudevents/sdk-go/v2/types"
)

var Event event.Event
var Bytes []byte
var Error error

func BenchmarkMarshal(b *testing.B) {
	now := types.Timestamp{Time: time.Now().UTC()}
	sourceUrl, _ := url.Parse("http://example.com/source")
	sourceV03 := &types.URIRef{URL: *sourceUrl}
	sourceV1 := &types.URIRef{URL: *sourceUrl}

	schemaUrl, _ := url.Parse("http://example.com/schema")
	schemaV03 := &types.URIRef{URL: *schemaUrl}
	schemaV1 := &types.URI{URL: *schemaUrl}

	testCases := map[string]struct {
		event           event.Event
		eventExtensions map[string]interface{}
	}{
		"struct data v0.3": {
			event: func() event.Event {
				e := event.Event{
					Context: event.EventContextV03{
						Type:      "com.example.test",
						Source:    *sourceV03,
						SchemaURL: schemaV03,
						ID:        "ABC-123",
						Time:      &now,
					}.AsV03(),
				}
				_ = e.SetData(event.ApplicationJSON, DataExample{
					AnInt:   42,
					AString: "testing",
				})
				return e
			}(),
			eventExtensions: map[string]interface{}{
				"exbool":   true,
				"exint":    int32(42),
				"exstring": "exstring",
				"exbinary": []byte{0, 1, 2, 3},
				"exurl":    sourceV03,
				"extime":   &now,
			},
		},
		"nil data v0.3": {
			event: event.Event{
				Context: event.EventContextV03{
					Type:            "com.example.test",
					Source:          *sourceV03,
					SchemaURL:       schemaV03,
					ID:              "ABC-123",
					Time:            &now,
					DataContentType: event.StringOfApplicationJSON(),
				}.AsV03(),
			},
			eventExtensions: map[string]interface{}{
				"exbool":   true,
				"exint":    int32(42),
				"exstring": "exstring",
				"exbinary": []byte{0, 1, 2, 3},
				"exurl":    sourceV03,
				"extime":   &now,
			},
		},
		"string data v0.3": {
			event: func() event.Event {
				e := event.Event{
					Context: event.EventContextV03{
						Type:      "com.example.test",
						Source:    *sourceV03,
						SchemaURL: schemaV03,
						ID:        "ABC-123",
						Time:      &now,
					}.AsV03(),
				}
				_ = e.SetData(event.ApplicationJSON, "This is a string.")
				return e
			}(),
			eventExtensions: map[string]interface{}{
				"exbool":   true,
				"exint":    int32(42),
				"exstring": "exstring",
				"exbinary": []byte{0, 1, 2, 3},
				"exurl":    sourceV03,
				"extime":   &now,
			},
		},
		"struct data v1.0": {
			event: func() event.Event {
				e := event.Event{
					Context: event.EventContextV1{
						Type:       "com.example.test",
						Source:     *sourceV1,
						DataSchema: schemaV1,
						ID:         "ABC-123",
						Time:       &now,
					}.AsV1(),
				}
				_ = e.SetData(event.ApplicationJSON, DataExample{
					AnInt:   42,
					AString: "testing",
				})
				return e
			}(),
			eventExtensions: map[string]interface{}{
				"exbool":   true,
				"exint":    int32(42),
				"exstring": "exstring",
				"exbinary": []byte{0, 1, 2, 3},
				"exurl":    sourceV1,
				"extime":   &now,
			},
		},
		"nil data v1.0": {
			event: event.Event{
				Context: event.EventContextV1{
					Type:            "com.example.test",
					Source:          *sourceV1,
					DataSchema:      schemaV1,
					ID:              "ABC-123",
					Time:            &now,
					DataContentType: event.StringOfApplicationJSON(),
				}.AsV1(),
			},
			eventExtensions: map[string]interface{}{
				"exbool":   true,
				"exint":    int32(42),
				"exstring": "exstring",
				"exbinary": []byte{0, 1, 2, 3},
				"exurl":    sourceV1,
				"extime":   &now,
			},
		},
		"string data v1.0": {
			event: func() event.Event {
				e := event.Event{
					Context: event.EventContextV1{
						Type:       "com.example.test",
						Source:     *sourceV1,
						DataSchema: schemaV1,
						ID:         "ABC-123",
						Time:       &now,
					}.AsV1(),
				}
				_ = e.SetData(event.ApplicationJSON, "This is a string.")
				return e
			}(),
			eventExtensions: map[string]interface{}{
				"exbool":   true,
				"exint":    int32(42),
				"exstring": "exstring",
				"exbinary": []byte{0, 1, 2, 3},
				"exurl":    sourceV1,
				"extime":   &now,
			},
		},
		"base64 json encoded data v1.0": {
			event: func() event.Event {
				e := event.Event{
					Context: event.EventContextV1{
						Type:       "com.example.test",
						Source:     *sourceV1,
						DataSchema: schemaV1,
						ID:         "ABC-123",
						Time:       &now,
					}.AsV1(),
				}
				_ = e.SetData(event.ApplicationJSON, []byte(`{"hello": "world"}`))
				return e
			}(),
		},
		"number data v1.0": {
			event: func() event.Event {
				e := event.Event{
					Context: event.EventContextV1{
						Type:       "com.example.test",
						Source:     *sourceV1,
						DataSchema: schemaV1,
						ID:         "ABC-123",
						Time:       &now,
					}.AsV1(),
				}
				_ = e.SetData(event.ApplicationJSON, 101)
				return e
			}(),
		},
	}
	for n, tc := range testCases {
		ev := tc.event.Clone()
		for k, v := range tc.eventExtensions {
			ev.SetExtension(k, v)
		}
		b.Run(n, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Bytes, Error = json.Marshal(ev)
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	now := types.Timestamp{Time: time.Now().UTC()}

	testCases := map[string][]byte{
		"struct data v0.3": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "0.3",
			"datacontenttype": "application/json",
			"data": map[string]interface{}{
				"a": 42,
				"b": "testing",
			},
			"id":        "ABC-123",
			"time":      now.Format(time.RFC3339Nano),
			"type":      "com.example.test",
			"exbool":    true,
			"exint":     42,
			"exstring":  "exstring",
			"exbinary":  "AAECAw==",
			"exurl":     "http://example.com/source",
			"extime":    now.Format(time.RFC3339Nano),
			"schemaurl": "http://example.com/schema",
			"source":    "http://example.com/source",
		}),
		"string data v0.3": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "0.3",
			"datacontenttype": "application/json",
			"data":            "This is a string.",
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"exbool":          true,
			"exint":           42,
			"exstring":        "exstring",
			"exbinary":        "AAECAw==",
			"exurl":           "http://example.com/source",
			"extime":          now.Format(time.RFC3339Nano),
			"schemaurl":       "http://example.com/schema",
			"source":          "http://example.com/source",
		}),
		"nil data v0.3": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "0.3",
			"datacontenttype": "application/json",
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"exbool":          true,
			"exint":           42,
			"exstring":        "exstring",
			"exbinary":        "AAECAw==",
			"exurl":           "http://example.com/source",
			"extime":          now.Format(time.RFC3339Nano),
			"schemaurl":       "http://example.com/schema",
			"source":          "http://example.com/source",
		}),
		"data, attributes and extensions and specversion with struct data v0.3": mustJsonMarshal(b, map[string]interface{}{
			"data": map[string]interface{}{
				"a": 42,
				"b": "testing",
			},
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"exbool":          true,
			"exint":           42,
			"exstring":        "exstring",
			"exbinary":        "AAECAw==",
			"exurl":           "http://example.com/source",
			"extime":          now.Format(time.RFC3339Nano),
			"schemaurl":       "http://example.com/schema",
			"source":          "http://example.com/source",
			"datacontenttype": "application/json",
			"specversion":     "0.3",
		}),
		"struct data v1.0": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "1.0",
			"datacontenttype": "application/json",
			"data": map[string]interface{}{
				"a": 42,
				"b": "testing",
			},
			"id":         "ABC-123",
			"time":       now.Format(time.RFC3339Nano),
			"type":       "com.example.test",
			"exbool":     true,
			"exint":      42,
			"exstring":   "exstring",
			"exbinary":   "AAECAw==",
			"exurl":      "http://example.com/source",
			"extime":     now.Format(time.RFC3339Nano),
			"dataschema": "http://example.com/schema",
			"source":     "http://example.com/source",
		}),
		"data, attributes and extensions and specversion with struct data v1.0": mustJsonMarshal(b, map[string]interface{}{
			"data": map[string]interface{}{
				"a": 42,
				"b": "testing",
			},
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"exbool":          true,
			"exint":           42,
			"exstring":        "exstring",
			"exbinary":        "AAECAw==",
			"exurl":           "http://example.com/source",
			"extime":          now.Format(time.RFC3339Nano),
			"dataschema":      "http://example.com/schema",
			"source":          "http://example.com/source",
			"datacontenttype": "application/json",
			"specversion":     "1.0",
		}),
		"string data v1.0": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "1.0",
			"datacontenttype": "application/json",
			"data":            "This is a string.",
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"exbool":          true,
			"exint":           42,
			"exstring":        "exstring",
			"exbinary":        "AAECAw==",
			"exurl":           "http://example.com/source",
			"extime":          now.Format(time.RFC3339Nano),
			"dataschema":      "http://example.com/schema",
			"source":          "http://example.com/source",
		}),
		"base64 json encoded data v1.0": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "1.0",
			"datacontenttype": "application/json",
			"data_base64":     base64.StdEncoding.EncodeToString([]byte(`{"hello":"world"}`)),
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"dataschema":      "http://example.com/schema",
			"source":          "http://example.com/source",
		}),
		"base64 xml encoded data v1.0": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "1.0",
			"datacontenttype": "application/json",
			"data_base64":     base64.StdEncoding.EncodeToString(mustEncodeWithDataCodec(b, event.ApplicationXML, &XMLDataExample{AnInt: 10})),
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"dataschema":      "http://example.com/schema",
			"source":          "http://example.com/source",
		}),
		"nil data v1.0": mustJsonMarshal(b, map[string]interface{}{
			"specversion":     "1.0",
			"datacontenttype": "application/json",
			"id":              "ABC-123",
			"time":            now.Format(time.RFC3339Nano),
			"type":            "com.example.test",
			"exbool":          true,
			"exint":           42,
			"exstring":        "exstring",
			"exbinary":        "AAECAw==",
			"exurl":           "http://example.com/source",
			"extime":          now.Format(time.RFC3339Nano),
			"dataschema":      "http://example.com/schema",
			"source":          "http://example.com/source",
		}),
	}
	for n, tc := range testCases {
		bytes := tc
		b.Run(n, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Event = event.Event{}
				Error = json.Unmarshal(bytes, &Event)
			}
		})
	}
}
