
package handlers

import (

    "github.com/appetito/uno"

	"github.com/appetito/uno/examples/greeter/api"
	"github.com/appetito/uno/examples/greeter/api"

)


//Greet a user, with some additional information
func GreetHandler(r uno.Request, request api.GreetRequest){

	var response api.Greeting{
		Message: fmt.Sprintf("Hello, %s! This is greeting number %d", request.Name, greetNumber),
	}
	r.RespondJSON(response)
}

