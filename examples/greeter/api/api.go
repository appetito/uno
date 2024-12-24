
package api

import (
    "encoding/json"
    "context"

    "github.com/appetito/uno"
    "github.com/nats-io/nats.go"
)

const (
    NS = "example"
    NAME = "Greeter"
    SERVICE_NAME = NS + "." + NAME


    //Greet a user, with some additional information
    GREET = "Greet"


    //Greet a user, with some additional information
    GREET_ENDPOINT = SERVICE_NAME + "." + GREET


)
type (

 //
 Greeting struct {
    Message string `json:"message"`
 }

 //
 GreetRequest struct {
    Name string `json:"name"`
 }

)


type GreeterClient struct {
    uc *uno.UnoClient
}


func NewGreeterClient(nc *nats.Conn, cfg *uno.UnoClientConfig) *GreeterClient {
    return &GreeterClient{
        uc: uno.NewUnoClient(nc, cfg),
    }
}



//Greet a user, with some additional information
func (c *GreeterClient) Greet(ctx context.Context, request GreetRequest) (response Greeting, err error) {
    reply, err := c.uc.RequestJSON(ctx, GREET_ENDPOINT, request)
    if err != nil {
        return response, err
    }
    err = json.Unmarshal(reply.Data, &response)
    return response, err
}

