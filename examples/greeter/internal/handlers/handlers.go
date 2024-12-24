
package handlers

import (

    "github.com/appetito/uno"

	"github.com/appetito/examples/greeter/api"

)


//Greet a user, with some additional information
func GreetHandler(r uno.Request, request api.GreetRequest){
	var response api.Greeting
	r.RespondJSON(response)
}

