# GeoDB-Go- The Golang Client to GeoDb - A Persistent Geospatial Database
    import "github.com/autom8ter/geodb-go"
    

GeoDB-Go is the official Golang gRPC client for [GeoDb - A Persistent Geospatial Database](https://github.com/autom8ter/geodb)

## Getting Started

### Creating a Client

```go
client, err  = geodb_go.NewClient(context.Background(), &geodb_go.ClientConfig{
		Host:     "localhost:8080",
		Password: "", //optional - only necessary if GEODB_PASSWORD is set server side
		Metrics:  true, //register client side prometheus metrics
		Retry:    true, //retry on failure
	})
	if err != nil {
		log.Fatal(err.Error())
	}
    //check if server is responsive
	resp, err := client.Ping(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Geodb server healthy= %v\n", resp.Ok)
```

## API Docs

