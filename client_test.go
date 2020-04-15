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
	file       *os.File
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
	cheyenneWhyoming = &api.Point{
		Lat: 41.1353874206543,
		Lon: -104.8226089477539,
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
	//create a driver object
	driverDetail, err := client.Set(context.Background(), &api.SetRequest{
		Object: &api.Object{
			Key: "driver_1", //its recommended to prefix keys or create a regex pattern for ease of querying data
			Point: &api.Point{
				Lat: 39.756378173828125,
				Lon: -104.99414825439453,
			},
			Radius: 100, //object radius for determining when objects intersect
			Metadata: map[string]string{ //optional object metadata
				"type": "driver",
			},
			GetAddress:  true, //true to get human readable address of current location on object detail
			GetTimezone: true, //true to get timezone of current location on object detail
			ExpiresUnix: 0,    //unix expiration timestamp, 0 for never expire
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("creating driver", driverDetail)

	//create a rider object
	riderDetail, err := client.Set(context.Background(), &api.SetRequest{
		Object: &api.Object{
			Key: "rider_1", //its recommended to prefix keys or create a regex pattern for ease of querying data
			Point: &api.Point{
				Lat: 39.74863815307617,
				Lon: -105.00762176513672,
			},
			Radius: 100, //object radius for determining when objects intersect
			Tracking: &api.ObjectTracking{ //optional, defaults to driving
				TravelMode: api.TravelMode_Driving,
				Trackers: []*api.ObjectTracker{
					{
						TargetObjectKey: "driver_1",
						TrackDirections: false, //
						TrackDistance:   true,
						TrackEta:        true,
					},
				},
			},
			Metadata: map[string]string{ //optional object metadata
				"type": "rider",
			},
			GetAddress:  true, //true to get human readable address of current location on object detail
			GetTimezone: true, //true to get timezone of current location on object detail
			ExpiresUnix: 0,    //unix expiration timestamp, 0 for never expire
		},
	})

	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("creating rider", riderDetail)

	//create a rider destination object
	riderDestinationDetail, err := client.Set(context.Background(), &api.SetRequest{
		Object: &api.Object{
			Key: "destination_1", //its recommended to prefix keys or create a regex pattern for ease of querying data
			Point: &api.Point{
				Lat: 39.71670913696289,
				Lon: -104.95344543457031,
			},
			Radius: 100, //object radius for determining when objects intersect
			Metadata: map[string]string{ //optional object metadata
				"type": "destination",
			},
			GetAddress:  true,                                      //true to get human readable address of current location on object detail
			GetTimezone: true,                                      //true to get timezone of current location on object detail
			ExpiresUnix: time.Now().Add(24 * time.Hour).UnixNano(), //automatically cleanup destination
		},
	})

	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("creating rider destination", riderDestinationDetail)

	//update driver to pickup rider with a tracker to get google maps directions, eta, and travel distance
	driverDetail, err = client.Set(context.Background(), &api.SetRequest{
		Object: &api.Object{
			Key: "driver_1", //its recommended to prefix keys or create a regex pattern for ease of querying data
			Point: &api.Point{
				Lat: 39.756378173828125,
				Lon: -104.99414825439453,
			},
			Radius: 100, //object radius for determining when objects intersect
			Metadata: map[string]string{ //optional object metadata
				"type": "driver",
			},
			Tracking: &api.ObjectTracking{ //optional, defaults to driving
				TravelMode: api.TravelMode_Driving,
				Trackers: []*api.ObjectTracker{
					{
						//track relationship to rider
						TargetObjectKey: "rider_1",
						TrackDirections: false, //get directions to pickup rider
						TrackDistance:   true,  //track distance to rider
						TrackEta:        true,  //track eta to rider
					},
					{
						//track relationship to destination
						TargetObjectKey: "destination_1",
						TrackDirections: true, //get directions to dropoff rider
						TrackDistance:   true, //track distance to riders destination
						TrackEta:        true, //track eta to rider destination
					},
				},
			},
			GetAddress:  true,                                       //true to get human readable address of current location on object detail
			GetTimezone: true,                                       //true to get timezone of current location on object detail
			ExpiresUnix: time.Now().Add(5 * time.Minute).UnixNano(), //unix expiration timestamp, 0 for never expire
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("updating driver", driverDetail)

	//update driver location to near the driver to simulate a pickup
	driverArrivalDetail, err := client.Set(context.Background(), &api.SetRequest{
		Object: &api.Object{
			Key: "driver_1", //its recommended to prefix keys or create a regex pattern for ease of querying data
			Point: &api.Point{
				Lat: 39.74863815307619,
				Lon: -105.0076217651367,
			},
			Radius: 100, //object radius for determining when objects intersect
			Metadata: map[string]string{ //optional object metadata
				"type": "driver",
			},
			Tracking: &api.ObjectTracking{ //optional, defaults to driving
				TravelMode: api.TravelMode_Driving,
				Trackers: []*api.ObjectTracker{
					{
						//track relationship to rider
						TargetObjectKey: "rider_1",
						TrackDirections: false, //get directions to pickup rider
						TrackDistance:   true,  //track distance to rider
						TrackEta:        true,  //track eta to rider
					},
					{
						//track relationship to destination
						TargetObjectKey: "destination_1",
						TrackDirections: true, //get directions to dropoff rider
						TrackDistance:   true, //track distance to riders destination
						TrackEta:        true, //track eta to rider destination
					},
				},
			},
			GetAddress:  true,                                       //true to get human readable address of current location on object detail
			GetTimezone: true,                                       //true to get timezone of current location on object detail
			ExpiresUnix: time.Now().Add(5 * time.Minute).UnixNano(), //unix expiration timestamp, 0 for never expire
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("updating driver for arrival", driverArrivalDetail)
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
		Prefix: "driver_",
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
		Regex: "driver_*",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(resp.Objects) != 1 {
		t.Fatal("expected 1 result")
	}
	prettyJson("TestClient_GetRegex", resp)
}

func TestClient_GetKeys(t *testing.T) {
	resp, err := client.GetKeys(context.Background(), &api.GetKeysRequest{})
	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("TestClient_GetKeys", resp)
}

func TestClient_GetPrefixKeys(t *testing.T) {
	resp, err := client.GetPrefixKeys(context.Background(), &api.GetPrefixKeysRequest{
		Prefix: "driver_",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(resp.Keys) != 1 {
		t.Fatal("expected 1 result")
	}
	prettyJson("TestClient_GetPrefixKeys", resp)
}

func TestClient_GetRegexKeys(t *testing.T) {
	resp, err := client.GetRegexKeys(context.Background(), &api.GetRegexKeysRequest{
		Regex: "driver_*",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(resp.Keys) != 1 {
		t.Fatal("expected 1 result")
	}
	prettyJson("TestClient_GetRegexKeys", resp)
}

func TestClient_ScanBound(t *testing.T) {
	//ScanBound scans a give geolocation boundary for objects, use regex/prefix methods to filter objects
	resp, err := client.ScanBound(context.Background(), &api.ScanBoundRequest{
		//a bound is like a circle on a map
		Bound: &api.Bound{
			Center: pepsiCenter, //center
			Radius: 5000,        //radius
		},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	prettyJson("TestClient_ScanBound", resp)
}

func TestClient_Delete(t *testing.T) {
	//Delete deletes an array of objects. if the first string is *, all objects will be dropped from the database
	_, err := client.Delete(context.Background(), &api.DeleteRequest{
		Keys: []string{"*"},
	})
	if err != nil {
		t.Fatal(err.Error())
	}
}
