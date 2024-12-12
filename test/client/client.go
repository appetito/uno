package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/appetito/uno"
	"github.com/nats-io/nats.go"

	// "github.com/nats-io/nats.go/uno"

	"github.com/rs/zerolog/log"
)




type Foo struct {
	A string `json:"a"`
}

type Bar struct {
	B string `json:"b"`
	T uint   `json:"t"`
} 





func main() {

	log.Info().Str("URL", nats.DefaultURL).Msg("Connecting to NATS")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	c := NewTestClient(nc, &uno.UnoClientConfig{TimeOut: 1 * time.Second})
	start := time.Now()
	// for i:=0; i < 1000; i++{

	// 	c.Foo(Foo{A: "foo"})
	// 	// fmt.Printf("Foo: %v, %s\n", r, err)
	
	// 	c.Bar(Bar{B: "bar", T: 0})
	// 	// fmt.Printf("Bar: %v, %s\n", r2, err)

	// }
	ctx, cancel := context.WithTimeout(context.Background(), 7 * time.Second)
	defer cancel()

	// ctx := context.Background()
	r2, err := c.Bar(ctx ,Bar{B: "bar", T: 5})
	fmt.Printf("Bar: %v, %s\n", r2, err)
	fmt.Printf("Took %s\n", time.Since(start))

}

type TestClient struct {
	uc *uno.UnoClient
}

func NewTestClient(nc *nats.Conn, cfg *uno.UnoClientConfig) *TestClient {
	c := &TestClient{
		uc: uno.NewUnoClient(nc, cfg),
	}
	return c
}

func (c *TestClient) Foo(ctx context.Context, request Foo) (resp Foo, err error) {
	reply, err := c.uc.RequestJSON(ctx, "test.struct_foo", request)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(reply.Data, &resp)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (c *TestClient)Bar(ctx context.Context, request Bar) (resp Bar, err error) {
	reply, err := c.uc.RequestJSON(ctx, "test.struct_bar", request)
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(reply.Data, &resp)
	if err != nil {
		return resp, err
	}
	return resp, err
}

