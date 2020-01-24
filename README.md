# Go SDK for [CloudEvents](https://github.com/cloudevents/spec)

[![go-doc](https://godoc.org/github.com/cloudevents/sdk-go?status.svg)](https://godoc.org/github.com/cloudevents/sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudevents/sdk-go)](https://goreportcard.com/report/github.com/cloudevents/sdk-go)
[![CircleCI](https://circleci.com/gh/cloudevents/sdk-go.svg?style=svg)](https://circleci.com/gh/cloudevents/sdk-go)
[![Releases](https://img.shields.io/github/release-pre/cloudevents/sdk-go.svg)](https://github.com/cloudevents/sdk-go/releases)
[![LICENSE](https://img.shields.io/github/license/cloudevents/sdk-go.svg)](https://github.com/cloudevents/sdk-go/blob/master/LICENSE)

**NOTE: This SDK is still considered work in progress.**

**With v1.0.0:**

The API that exists under [`pkg/cloudevents`](./pkg/cloudevents) will follow
semver rules. This applies to the root [`./alias.go`](./alias.go) file as well.

The API that exists under [`pkg/bindings`](./pkg/bindings) is a new API that
will become SDK v2.x, and will replace `pkg/cloudevents`. This area is still
under heavy development and will not be following the same semver rules as
`pkg/cloudevents`. If a release is required to ship changes to `pkg/bindings`, a
bug fix release will be issued (x.y.z+1).

We will target ~2 months of development to release v2 of this SDK with an end
date of March 27. You can read more about the plan for SDK v2 in the
[SDK v2 planning doc](./docs/SDK_v2.md).

## Working with CloudEvents

Package [cloudevents](./pkg/cloudevents) provides primitives to work with
CloudEvents specification: https://github.com/cloudevents/spec.

Import this repo to get the `cloudevents` package:

```go
import "github.com/cloudevents/sdk-go"
```

Receiving a cloudevents.Event via the HTTP Transport:

```go
func Receive(event cloudevents.Event) {
	// do something with event.Context and event.Data (via event.DataAs(foo)
}

func main() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), Receive));
}
```

Creating a minimal CloudEvent in version 0.2:

```go
event := cloudevents.NewEvent()
event.SetID("ABC-123")
event.SetType("com.cloudevents.readme.sent")
event.SetSource("http://localhost:8080/")
event.SetData(data)
```

Sending a cloudevents.Event via the HTTP Transport with Binary v0.2 encoding:

```go
t, err := cloudevents.NewHTTPTransport(
	cloudevents.WithTarget("http://localhost:8080/"),
	cloudevents.WithEncoding(cloudevents.HTTPBinaryV02),
)
if err != nil {
	panic("failed to create transport, " + err.Error())
}

c, err := cloudevents.NewClient(t)
if err != nil {
	panic("unable to create cloudevent client: " + err.Error())
}
if err := c.Send(ctx, event); err != nil {
	panic("failed to send cloudevent: " + err.Error())
}
```

Or, the transport can be set to produce CloudEvents using the selected encoding
but not change the provided event version, here the client is set to output
structured encoding:

```go
t, err := cloudevents.NewHTTPTransport(
	cloudevents.WithTarget("http://localhost:8080/"),
	cloudevents.WithStructuredEncoding(),
)
```

If you are using advanced transport features or have implemented your own
transport integration, provide it to a client so your integration does not
change:

```go
t, err := cloudevents.NewHTTPTransport(
	cloudevents.WithPort(8181),
	cloudevents.WithPath("/events/")
)
// or a custom transport: t := &custom.MyTransport{Cool:opts}

c, err := cloudevents.NewClient(t, opts...)
```

Checkout the sample [sender](./cmd/samples/http/sender) and
[receiver](./cmd/samples/http/receiver) applications for working demo.
