package geodb_go_test

import (
	"context"
	"fmt"
	geodb_go "github.com/autom8ter/geodb-go"
	"log"
	"os"
	"testing"
)

var client *geodb_go.Client
var err error

func TestMain(t *testing.M) {
	client, err = geodb_go.NewClient(context.Background(), &geodb_go.ClientConfig{
		Host:     "localhost:8080",
		Password: "",   //optional - only necessary if GEODB_PASSWORD is set server side
		Metrics:  true, //register client side prometheus metrics
		Retry:    true, //retry on failure
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := client.Ping(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Geodb server healthy= %v\n", resp.Ok)
	os.Exit(t.Run())
}

func TestPing(t *testing.T) {

}
