
package api

import (
    "encoding/json"
    "context"

    "github.com/appetito/uno"
    "github.com/nats-io/nats.go"
)

const (
    NS = "example"
    NAME = "GreetAnalytics"
    SERVICE_NAME = NS + "." + NAME


    //Get user's greet stats
    GET_USERS_STATS = "GetUsersStats"

    //Get top greeted users
    TOP_GREETED_USERS = "TopGreetedUsers"


    //Get user's greet stats
    GET_USERS_STATS_ENDPOINT = SERVICE_NAME + "." + GET_USERS_STATS

    //Get top greeted users
    TOP_GREETED_USERS_ENDPOINT = SERVICE_NAME + "." + TOP_GREETED_USERS


)
type (

 //
 UserStats struct {
    Name string `json:"name"`
    GreetCount int64 `json:"greet_count"`
 }

 //
 GetUsersStatsRequest struct {
    Name string `json:"name"`
 }

 //
 TopGreetedUsersRequest struct {
    Count int64 `json:"count"`
 }

)


type GreetAnalyticsClient struct {
    uc *uno.UnoClient
}


func NewGreetAnalyticsClient(nc *nats.Conn, cfg *uno.UnoClientConfig) *GreetAnalyticsClient {
    return &GreetAnalyticsClient{
        uc: uno.NewUnoClient(nc, cfg),
    }
}



//Get user's greet stats
func (c *GreetAnalyticsClient) GetUsersStats(ctx context.Context, request GetUsersStatsRequest) (response UserStats, err error) {
    reply, err := c.uc.RequestJSON(ctx, GET_USERS_STATS_ENDPOINT, request)
    if err != nil {
        return response, err
    }
    err = json.Unmarshal(reply.Data, &response)
    return response, err
}

//Get top greeted users
func (c *GreetAnalyticsClient) TopGreetedUsers(ctx context.Context, request TopGreetedUsersRequest) (response []UserStats, err error) {
    reply, err := c.uc.RequestJSON(ctx, TOP_GREETED_USERS_ENDPOINT, request)
    if err != nil {
        return response, err
    }
    err = json.Unmarshal(reply.Data, &response)
    return response, err
}

