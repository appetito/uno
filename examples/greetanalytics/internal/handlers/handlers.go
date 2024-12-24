package handlers

import (
	"github.com/appetito/uno"

	"github.com/appetito/examples/greetanalytics/api"
)

var stats map[string]int64


//Get user's greet stats
func GetUsersStatsHandler(r uno.Request, request api.GetUsersStatsRequest){
	count := stats[request.Name]
	
	response := api.UserStats{
		Name: request.Name,
		GreetCount: count,
	}
	r.RespondJSON(response)
}

//Get top greeted users
func TopGreetedUsersHandler(r uno.Request, request api.TopGreetedUsersRequest){
	response := []api.UserStats{}
	r.RespondJSON(response)
}

