package v01_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dispatchframework/cloudevents-go-sdk/v01"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPMarshallerFromRequestBinaryBase64Success(t *testing.T) {
	factory := v01.NewDefaultHTTPMarshaller()

	header := http.Header{}
	header.Set("Content-Type", "application/octet-stream")
	header.Set("CE-EventType", "dispatch")
	header.Set("CE-Source", "dispatch")
	header.Set("CE-EventID", "00001")
	header.Set("CE-MyExtension", "myvalue")
	header.Set("CE-AnotherExtension", "anothervalue")
	header.Set("CE-EventTime", "2018-08-08T15:00:00-07:00")

	body := bytes.NewBufferString("This is a byte array of data.")
	req := httptest.NewRequest("GET", "localhost:8080", ioutil.NopCloser(body))
	req.Header = header

	actual, err := factory.FromRequest(req)
	require.NoError(t, err)

	timestamp, err := time.Parse(time.RFC3339, "2018-08-08T15:00:00-07:00")
	expected := &v01.Event{
		ContentType: "application/octet-stream",
		EventType:   "dispatch",
		Source:      "dispatch",
		EventID:     "00001",
		EventTime:   &timestamp,
		Data:        []byte("This is a byte array of data."),
	}

	expected.Set("myextension", "myvalue")
	expected.Set("anotherextension", "anothervalue")

	assert.EqualValues(t, expected, actual)
}

func TestHTTPMarshallerToRequestBinaryBase64Success(t *testing.T) {
	factory := v01.NewDefaultHTTPMarshaller()

	event := v01.Event{
		CloudEventsVersion: "0.1",
		EventType:          "dispatch",
		EventTypeVersion:   "0.1",
		EventID:            "00001",
		Source:             "dispatch",
		ContentType:        "application/octet-stream",
		Data:               []byte("This is a byte array of data"),
	}

	event.Set("myfloat", 100e+3)
	event.Set("myint", 100)
	event.Set("mybool", true)
	event.Set("mystring", "string")

	actual, _ := http.NewRequest("GET", "localhost:8080", nil)
	err := factory.ToRequest(actual, &event)
	require.NoError(t, err)

	buffer := bytes.NewBufferString("This is a byte array of data")

	expected, _ := http.NewRequest("GET", "localhost:8080", buffer)
	expected.Header.Set("CE-CloudEventsVersion", "0.1")
	expected.Header.Set("CE-EventID", "00001")
	expected.Header.Set("CE-EventType", "dispatch")
	expected.Header.Set("CE-EventTypeVersion", "0.1")
	expected.Header.Set("CE-Source", "dispatch")
	expected.Header.Set("CE-Myfloat", "100000")
	expected.Header.Set("CE-Myint", "100")
	expected.Header.Set("CE-Mybool", "true")
	expected.Header.Set("CE-Mystring", "string")
	expected.Header.Set("Content-Type", "application/octet-stream")

	// Can't test function equality
	expected.GetBody = nil
	actual.GetBody = nil

	assert.EqualValues(t, expected, actual)
}
