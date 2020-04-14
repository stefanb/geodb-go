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

### Storing Object Geolocation Data

```go
resp, err := client.Set(context.Background(), &api.SetRequest{
		Objects:              []*api.Object{
			{
				Key:    "testing_coors", //its recommended to prefix keys or create a regex pattern for ease of querying data
				Point:  &api.Point{
					Lat: 39.756378173828125,
					Lon: -104.99414825439453,
				},
				Radius: 100, //radius of the object- this is used to determine when objects are overlapping(Geofencing) to determine if theyre in the same place
				Tracking: &api.ObjectTracking{
					TravelMode: api.TravelMode_Driving,
				},
				Metadata: map[string]string{ //optional metadata about the object
						"type": "sports",
				},
				GetAddress:  true, //true if you want a human readable address in the object detail response
                GetTimezone: true, //true if you want the timezone in the object detail response
				ExpiresUnix: 0, //0 means never expire
			},
			{

				Key:    "testing_pepsi_center",
				Point:  &api.Point{
					Lat: 39.74863815307617,
					Lon: -105.00762176513672,
				}, 
				Radius: 100, //radius of the object- this is used to determine when objects are overlapping(Geofencing) to determine if theyre in the same place
				Tracking: &api.ObjectTracking{
					TravelMode: api.TravelMode_Driving, //different travel modes impact directions and eta for trackers below
					Trackers: []*api.ObjectTracker{
						{
							TargetObjectKey: "testing_coors",
							TrackDirections: true, //adds google maps directions to this object in object detail response
                            TrackDistance:   true, //adds real-distance to this object in object detail response
                            TrackEta:        true, //adds eta(depending on travel mode) to this object in object detail response
                         },
					},
				},
				Metadata: map[string]string{ //optional metadata about the object
					"type": "sports",
				},
				GetAddress:  true, //true if you want a human readable address in the object detail response
				GetTimezone: true, //true if you want the timezone in the object detail response
				ExpiresUnix: time.Now().Add(5 * time.Minute).Unix(), //automatically expire this object on the unix timestamp, leave empty if no expiration
			},
			{
				Key:    "malls_cherry_creek_mall",
				Point:  &api.Point{
					Lat: 39.71670913696289,
					Lon: -104.95344543457031,
				},
				Radius: 100, //radius of the object- this is used to determine when objects are overlapping(Geofencing) to determine if theyre in the same place
				Tracking: &api.ObjectTracking{
					TravelMode: api.TravelMode_Driving, //different travel modes impact directions and eta for trackers below
					Trackers: []*api.ObjectTracker{ //add trackers to track the objects geolocation in relation to another
						{
							TargetObjectKey: "testing_pepsi_center", //the object to track, must exist in database already
							TrackDirections: true, //adds google maps directions to this object in object detail response
							TrackDistance:   true, //adds real-distance to this object in object detail response
							TrackEta:        true, //adds eta(depending on travel mode) to this object in object detail response
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
```

Pretty JSON Response: 

```json
{
    "objects": {
        "malls_cherry_creek_mall": {
            "object": {
                "key": "malls_cherry_creek_mall",
                "point": {
                    "lat": 39.71670913696289,
                    "lon": -104.95344543457031
                },
                "radius": 100,
                "tracking": {
                    "trackers": [
                        {
                            "target_object_key": "testing_pepsi_center",
                            "track_directions": true,
                            "track_distance": true,
                            "track_eta": true
                        }
                    ]
                },
                "metadata": {
                    "type": "mall"
                },
                "get_address": true,
                "get_timezone": true,
                "expires_unix": 1586900513,
                "updated_unix": 1586900213
            },
            "address": {
                "state": "Colorado",
                "address": "Unnamed Road, Denver, CO 80209, USA",
                "country": "United States",
                "zip": "80209",
                "county": "Denver County",
                "city": "Denver"
            },
            "timezone": "America/Denver",
            "events": [
                {
                    "object": {
                        "key": "testing_pepsi_center",
                        "point": {
                            "lat": 39.74863815307617,
                            "lon": -105.00762176513672
                        },
                        "radius": 100,
                        "tracking": {
                            "trackers": [
                                {
                                    "target_object_key": "testing_coors",
                                    "track_directions": true,
                                    "track_distance": true,
                                    "track_eta": true
                                }
                            ]
                        },
                        "metadata": {
                            "type": "sports"
                        },
                        "get_address": true,
                        "get_timezone": true,
                        "expires_unix": 1586900513,
                        "updated_unix": 1586900213
                    },
                    "distance": 5843.275955551314,
                    "direction": {
                        "html_directions": "CjxoNT5EZXN0aW5hdGlvbjogMTAwMCBDaG9wcGVyIENpciwgRGVudmVyLCBDTyA4MDIwNCwgVVNBPC9oNT5IZWFkIDxiPnNvdXRoPC9iPiAtIDAuMyBtaTxicj5UdXJuIDxiPmxlZnQ8L2I+IG9udG8gPGI+U3RlZWxlIFN0PC9iPiAtIDQzMyBmdDxicj5UdXJuIDxiPmxlZnQ8L2I+IG9udG8gPGI+RSAxc3QgQXZlPC9iPiAtIDAuOSBtaTxicj5Db250aW51ZSBzdHJhaWdodCB0byBzdGF5IG9uIDxiPkUgMXN0IEF2ZTwvYj4gLSAwLjMgbWk8YnI+Q29udGludWUgb250byA8Yj5TcGVlciBCbHZkPC9iPiAtIDIuOSBtaTxicj5UdXJuIDxiPmxlZnQ8L2I+IG9udG8gPGI+Q2hvcHBlciBDaXI8L2I+IC0gMzk3IGZ0PGJyPlR1cm4gPGI+cmlnaHQ8L2I+IG9udG8gPGI+MTF0aCBTdDwvYj4vPHdici8+PGI+MTJ0aCBTdDwvYj4gLSAzNTggZnQ8YnI+U2xpZ2h0IDxiPmxlZnQ8L2I+IC0gMC4xIG1pPGJyPg==",
                        "eta": 15,
                        "travel_dist": 7581
                    },
                    "timestamp_unix": 1586900213
                }
            ]
        },
        "testing_coors": {
            "object": {
                "key": "testing_coors",
                "point": {
                    "lat": 39.756378173828125,
                    "lon": -104.99414825439453
                },
                "radius": 100,
                "tracking": {},
                "metadata": {
                    "type": "sports"
                },
                "get_address": true,
                "get_timezone": true,
                "updated_unix": 1586900213
            },
            "address": {
                "state": "Colorado",
                "address": "2001 Blake St, Denver, CO 80205, USA",
                "country": "United States",
                "zip": "80205",
                "county": "Denver County",
                "city": "Denver"
            },
            "timezone": "America/Denver"
        },
        "testing_pepsi_center": {
            "object": {
                "key": "testing_pepsi_center",
                "point": {
                    "lat": 39.74863815307617,
                    "lon": -105.00762176513672
                },
                "radius": 100,
                "tracking": {
                    "trackers": [
                        {
                            "target_object_key": "testing_coors",
                            "track_directions": true,
                            "track_distance": true,
                            "track_eta": true
                        }
                    ]
                },
                "metadata": {
                    "type": "sports"
                },
                "get_address": true,
                "get_timezone": true,
                "expires_unix": 1586900513,
                "updated_unix": 1586900213
            },
            "address": {
                "state": "Colorado",
                "address": "1000 Chopper Cir, Denver, CO 80204, USA",
                "country": "United States",
                "zip": "80204",
                "county": "Denver County",
                "city": "Denver"
            },
            "timezone": "America/Denver",
            "events": [
                {
                    "object": {
                        "key": "testing_coors",
                        "point": {
                            "lat": 39.756378173828125,
                            "lon": -104.99414825439453
                        },
                        "radius": 100,
                        "tracking": {},
                        "metadata": {
                            "type": "sports"
                        },
                        "get_address": true,
                        "get_timezone": true,
                        "updated_unix": 1586900213
                    },
                    "distance": 1439.4645850870015,
                    "direction": {
                        "html_directions": "CjxoNT5EZXN0aW5hdGlvbjogMjAwMSBCbGFrZSBTdCwgRGVudmVyLCBDTyA4MDIwNSwgVVNBPC9oNT5IZWFkIDxiPm5vcnRoZWFzdDwvYj4gdG93YXJkIDxiPjExdGggU3Q8L2I+Lzx3YnIvPjxiPjEydGggU3Q8L2I+IC0gMC4xIG1pPGJyPkNvbnRpbnVlIG9udG8gPGI+MTF0aCBTdDwvYj4vPHdici8+PGI+MTJ0aCBTdDwvYj4gLSAzNTggZnQ8YnI+VHVybiA8Yj5sZWZ0PC9iPiBvbnRvIDxiPkNob3BwZXIgQ2lyPC9iPiAtIDMyOCBmdDxicj5Db250aW51ZSBvbnRvIDxiPldld2F0dGEgU3Q8L2I+IC0gMC45IG1pPGJyPlR1cm4gPGI+cmlnaHQ8L2I+IG9udG8gPGI+MjJuZCBTdDwvYj4gLSAwLjMgbWk8YnI+VHVybiA8Yj5sZWZ0PC9iPiBvbnRvIDxiPk1hcmtldCBTdDwvYj4gLSA0ODIgZnQ8YnI+VHVybiA8Yj5sZWZ0PC9iPiBvbnRvIDxiPlBhcmsgQXZlIFc8L2I+IC0gMC4xIG1pPGJyPlR1cm4gPGI+cmlnaHQ8L2I+IGF0IDxiPldhemVlIFN0PC9iPiAtIDAuMSBtaTxicj5UdXJuIDxiPmxlZnQ8L2I+IC0gMjIwIGZ0PGJyPlR1cm4gPGI+bGVmdDwvYj4gLSAwLjMgbWk8YnI+VHVybiA8Yj5sZWZ0PC9iPiAtIDExNSBmdDxicj4=",
                        "eta": 9,
                        "travel_dist": 3380
                    },
                    "timestamp_unix": 1586900213
                }
            ]
        }
    }
}
```

### Getting Geolocation Data(Objects)

#### Get All Objects

```go
resp, err := client.Get(context.Background(), &api.GetRequest{
		Keys:                 nil, //nil for getting all objects in db
	})
for _, object := range resp.Objects {
		fmt.Println(object.String())
}
```

#### Get Objects by Prefix

```go
resp, err := client.GetPrefix(context.Background(), &api.GetPrefixRequest{
		Prefix: "malls_",
	})
	for _, object := range resp.Objects {
		fmt.Println(object.String())
	}
```

#### Get Objects by Regex

```go
resp, err := client.GetRegex(context.Background(), &api.GetRegexRequest{
		Regex: "malls_*",
	})
for _, object := range resp.Objects {
		fmt.Println(object.String())
	}
```
## API Docs

- Client API Docs can be found [here](https://github.com/autom8ter/geodb-go/blob/master/DOCS.md)
- Profobuf Contract can be found [here](https://github.com/autom8ter/geodb/blob/master/api.proto)