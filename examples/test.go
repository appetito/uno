package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"math/rand"

	"github.com/rs/zerolog/log"

	"github.com/appetito/uno"
	analyticsapi "github.com/appetito/uno/examples/greetanalytics/api"
	greetapi "github.com/appetito/uno/examples/greeter/api"
	"github.com/nats-io/nats.go"
)

var names = []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Helen", "Ivan", "Jack"}


func main(){

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	ac := analyticsapi.NewGreetAnalyticsClient(nc, &uno.UnoClientConfig{})
	gc := greetapi.NewGreeterClient(nc, &uno.UnoClientConfig{})

	workersFlag := flag.Int("workers", 10, "Number of workers")
	greetsFlag := flag.Int("greets", 1000, "Number of greets")

	flag.Parse()
	numWorkers := *workersFlag
	numJobs := *greetsFlag

	startTime := time.Now()

	fmt.Printf("Starting %d workers with %d greetings\n", numWorkers, numJobs)

	jobs := make(chan string, numJobs)
	
	wg := &sync.WaitGroup{}

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go GreetWorker(gc, jobs, wg)
	}

	for i := 0; i < numJobs; i++ {
		jobs <- names[rand.Intn(len(names))]
	}
	close(jobs)

	wg.Wait()

	resp, err := ac.TopGreetedUsers(context.Background(), analyticsapi.TopGreetedUsersRequest{Count: 10})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top greeted users")
	}
	fmt.Printf("Elapsed: %s\n", time.Now().Sub(startTime))
	fmt.Printf("Average greets per second: %f\n", float64(numJobs)/time.Now().Sub(startTime).Seconds())
	fmt.Printf("Top greeted users: %v\n", resp)
}	

func GreetWorker(c *greetapi.GreeterClient, in <-chan string, wg *sync.WaitGroup){
	defer wg.Done()
	for name := range in {
		resp, err := c.Greet(context.Background(), greetapi.GreetRequest{Name: name})
		if err != nil {
			log.Error().Err(err).Msg("Failed to greet")
		}
		if!strings.Contains(resp.Message, name) {
			log.Error().Str("name", name).Str("msg", resp.Message).Msg("Invalid response")
		}

	}
}