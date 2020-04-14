package geodb_go_test

import (
	"context"
	"encoding/json"
	"fmt"
	geodb_go "github.com/autom8ter/geodb-go"
	api "github.com/autom8ter/geodb/gen/go/geodb"
	"log"
	"os"
	"testing"
	"time"
)

var (
	file *os.File
	client     *geodb_go.Client
	err        error
	coorsField = &api.Point{
		Lat: 39.756378173828125,
		Lon: -104.99414825439453,
	}
	pepsiCenter = &api.Point{
		Lat: 39.74863815307617,
		Lon: -105.00762176513672,
	}
	cherryCreekMall = &api.Point{
		Lat: 39.71670913696289,
		Lon: -104.95344543457031,
	}
	saintJosephHospital = &api.Point{
		Lat: 39.74626922607422,
		Lon: -104.97151184082031,
	}
	objects = []*api.Object{{
		Key: "testing_coors",
		Point: coorsField,
		Radius: 100,
		Tracking: &api.ObjectTracking{
			TravelMode: api.TravelMode_Driving,
		},
		Metadata: map[string]string{
			"type": "sports",
		},
		GetAddress:  true,
		GetTimezone: true,
		ExpiresUnix: 0,
	},
		{

			Key: "testing_pepsi_center",
			Point: pepsiCenter,
			Radius: 100,
			Tracking: &api.ObjectTracking{
				TravelMode: api.TravelMode_Driving,
				Trackers: []*api.ObjectTracker{
					{
						TargetObjectKey: "testing_coors",
						TrackDirections: true,
						TrackDistance:   true,
						TrackEta:        true,
					},
				},
			},
			Metadata: map[string]string{
				"type": "sports",
			},
			GetAddress:  true,
			GetTimezone: true,
			ExpiresUnix: time.Now().Add(5 * time.Minute).Unix(),
		},
		{
			Key: "malls_cherry_creek_mall",
			Point: cherryCreekMall,
			Radius: 100,
			Tracking: &api.ObjectTracking{
				TravelMode: api.TravelMode_Driving,
				Trackers: []*api.ObjectTracker{
					{
						TargetObjectKey: "testing_pepsi_center",
						TrackDirections: true,
						TrackDistance:   true,
						TrackEta:        true,
					},
				},
			},
			Metadata: map[string]string{
				"type": "mall",
			},
			GetAddress:  true,
			GetTimezone: true,
			ExpiresUnix: time.Now().Add(5 * time.Minute).Unix(),
		},
	}
)

func prettyJson(test string, obj interface{}) {
	bits, _ := json.MarshalIndent(obj, "", "    ")
	file.WriteString(fmt.Sprintln("--------------------"))
	file.WriteString(fmt.Sprintln(test))
	file.WriteString(fmt.Sprintln(string(bits)))
}

func TestMain(t *testing.M) {
	os.Remove("test_results.txt")
	file, err = os.Create("test_results.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
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

func TestClient_Set(t *testing.T) {
	resp, err := client.Set(context.Background(), &api.SetRequest{
		Objects: []*api.Object{
			{
				Key: "testing_coors",
				Point: &api.Point{
					Lat: 39.756378173828125,
					Lon: -104.99414825439453,
				},
				Radius: 100,
				Tracking: &api.ObjectTracking{
					TravelMode: api.TravelMode_Driving,
				},
				Metadata: map[string]string{
					"type": "sports",
				},
				GetAddress:  true,
				GetTimezone: true,
				ExpiresUnix: 0,
			},
			{

				Key: "testing_pepsi_center",
				Point: &api.Point{
					Lat: 39.74863815307617,
					Lon: -105.00762176513672,
				},
				Radius: 100,
				Tracking: &api.ObjectTracking{
					TravelMode: api.TravelMode_Driving,
					Trackers: []*api.ObjectTracker{
						{
							TargetObjectKey: "testing_coors",
							TrackDirections: true,
							TrackDistance:   true,
							TrackEta:        true,
						},
					},
				},
				Metadata: map[string]string{
					"type": "sports",
				},
				GetAddress:  true,
				GetTimezone: true,
				ExpiresUnix: time.Now().Add(5 * time.Minute).Unix(),
			},
			{
				Key: "malls_cherry_creek_mall",
				Point: &api.Point{
					Lat: 39.71670913696289,
					Lon: -104.95344543457031,
				},
				Radius: 100,
				Tracking: &api.ObjectTracking{
					TravelMode: api.TravelMode_Driving,
					Trackers: []*api.ObjectTracker{
						{
							TargetObjectKey: "testing_pepsi_center",
							TrackDirections: true,
							TrackDistance:   true,
							TrackEta:        true,
						},
					},
				},
				Metadata: map[string]string{
					"type": "mall",
				},
				GetAddress:  true,
				GetTimezone: true,
				ExpiresUnix: time.Now().Add(5 * time.Minute).Unix(),
			},
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("TestClient_Set", resp)
}

func TestClient_Get(t *testing.T) {
	resp, err := client.Get(context.Background(), &api.GetRequest{
		Keys: nil, //nil for getting all objects in db
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(resp.Objects) != 3 {
		t.Fatal("expected 3 results")
	}
	prettyJson("TestClient_Get", resp)
}

func TestClient_GetPrefix(t *testing.T) {
	resp, err := client.GetPrefix(context.Background(), &api.GetPrefixRequest{
		Prefix: "malls_",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(resp.Objects) != 1 {
		t.Fatal("expected 1 result")
	}
	prettyJson("TestClient_GetPrefix", resp)
}

func TestClient_GetRegex(t *testing.T) {
	resp, err := client.GetRegex(context.Background(), &api.GetRegexRequest{
		Regex: "malls_*",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(resp.Objects) != 1 {
		t.Fatal("expected 1 result")
	}
	prettyJson("TestClient_GetRegex", resp)
}
