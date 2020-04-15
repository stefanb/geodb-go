# GeoDB-Go- The Golang Client to GeoDb - A Persistent Geospatial Database
    import "github.com/autom8ter/geodb-go"
    

GeoDB-Go is the official Golang gRPC client for [GeoDb - A Persistent Geospatial Database](https://github.com/autom8ter/geodb)
## Methodology

- Clients may query the database in three ways keys(unique ids), prefix-scanning, or regex 
- Clients can open and execute logic on object geolocation streams that can be filtered by keys(unique ids), prefix-scanning, or regex
- Clients can manage object-centric, dynamic geofences(trackers) that can be used to track an objects location in relation to other registered objects
- Haversine formula is used to calculate whether objects are overlapping using object coordinates and their radius.
- If the server has a google maps api key present in its environmental variables, all geofencing(trackers) will be enhanced with html directions, estimated time of arrival, and more.

## Use Cases
- Ride Sharing
- Food Delivery
- Asset Tracking


## Getting Started

### Docker Compose - Server

```yaml
version: '3.7'
services:
  db:
    image: colemanword/geodb:latest
    env_file:
      - geodb.env
    ports:
      - "8080:8080"
    volumes:
      - default:/tmp/geodb
    networks:
      default:
        aliases:
          - geodb
networks:
  default:

volumes:
  default:

```

geodb.env:

```.env
GEODB_PORT (optional) default: :8080
GEODB_PATH (optional) default: /tmp/geodb
GEODB_GC_INTERVAL (optional) default: 5m
GEODB_PASSWORD (optional) 
GEODB_GMAPS_KEY (optional)

```

Up:

    docker-compose -f docker-compose.yml pull
    docker-compose -f docker-compose.yml up -d

Down:

    docker-compose -f docker-compose.yml down --remove-orphans

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

### Short Example - Uber-like functionality for ride-sharing

#### Create a Driver

```go
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
```

Pretty JSON Response: 

```json
{
    "object": {
        "object": {
            "key": "driver_1",
            "point": {
                "lat": 39.756378173828125,
                "lon": -104.99414825439453
            },
            "radius": 100,
            "metadata": {
                "type": "driver"
            },
            "get_address": true,
            "get_timezone": true,
            "updated_unix": 1586917896
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
    }
}
```


#### Create a Rider

```go
//create a rider object
	riderDetail, err := client.Set(context.Background(), &api.SetRequest{
		Object: &api.Object{
			Key: "rider_1", //its recommended to prefix keys or create a regex pattern for ease of querying data
			Point: &api.Point{
				Lat: 39.74863815307617,
				Lon: -105.00762176513672,
			},
			Radius: 100, //object radius for determining when objects intersect
			Tracking: &api.ObjectTracking{ 
				//optional, defaults to driving
				TravelMode: api.TravelMode_Driving,
				//track riders location relationship to the driver
				Trackers: []*api.ObjectTracker{
					{
						TargetObjectKey: "driver_1",
						TrackDirections: false, //no need to get directions to driver
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
```

Pretty JSON Response: 

```json
{
    "object": {
        "object": {
            "key": "rider_1",
            "point": {
                "lat": 39.74863815307617,
                "lon": -105.00762176513672
            },
            "radius": 100,
            "tracking": {
                "trackers": [
                    {
                        "target_object_key": "driver_1",
                        "track_distance": true,
                        "track_eta": true
                    }
                ]
            },
            "metadata": {
                "type": "rider"
            },
            "get_address": true,
            "get_timezone": true,
            "updated_unix": 1586917896
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
                    "key": "driver_1",
                    "point": {
                        "lat": 39.756378173828125,
                        "lon": -104.99414825439453
                    },
                    "radius": 100,
                    "metadata": {
                        "type": "driver"
                    },
                    "get_address": true,
                    "get_timezone": true,
                    "updated_unix": 1586917896
                },
                "distance": 1439.4645850870015,
                "direction": {
                    "eta": 8,
                    "travel_dist": 3380
                },
                "timestamp_unix": 1586917896
            }
        ]
    }
}
```

#### Create Rider Destination 

```go
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
```

Pretty JSON Response: 

```json
{
    "object": {
        "object": {
            "key": "destination_1",
            "point": {
                "lat": 39.71670913696289,
                "lon": -104.95344543457031
            },
            "radius": 100,
            "metadata": {
                "type": "destination"
            },
            "get_address": true,
            "get_timezone": true,
            "expires_unix": 1587004296815040000,
            "updated_unix": 1586917896
        },
        "address": {
            "state": "Colorado",
            "address": "Unnamed Road, Denver, CO 80209, USA",
            "country": "United States",
            "zip": "80209",
            "county": "Denver County",
            "city": "Denver"
        },
        "timezone": "America/Denver"
    }
}
```

#### Update Driver to Track Rider & Riders Destination

```go
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
```

Pretty JSON Response: 

```json
{
    "object": {
        "object": {
            "key": "driver_1",
            "point": {
                "lat": 39.756378173828125,
                "lon": -104.99414825439453
            },
            "radius": 100,
            "tracking": {
                "trackers": [
                    {
                        "target_object_key": "rider_1",
                        "track_distance": true,
                        "track_eta": true
                    },
                    {
                        "target_object_key": "destination_1",
                        "track_directions": true,
                        "track_distance": true,
                        "track_eta": true
                    }
                ]
            },
            "metadata": {
                "type": "driver"
            },
            "get_address": true,
            "get_timezone": true,
            "expires_unix": 1586918197040534000,
            "updated_unix": 1586917897
        },
        "address": {
            "state": "Colorado",
            "address": "2001 Blake St, Denver, CO 80205, USA",
            "country": "United States",
            "zip": "80205",
            "county": "Denver County",
            "city": "Denver"
        },
        "timezone": "America/Denver",
        "events": [
            {
                "object": {
                    "key": "rider_1",
                    "point": {
                        "lat": 39.74863815307617,
                        "lon": -105.00762176513672
                    },
                    "radius": 100,
                    "tracking": {
                        "trackers": [
                            {
                                "target_object_key": "driver_1",
                                "track_distance": true,
                                "track_eta": true
                            }
                        ]
                    },
                    "metadata": {
                        "type": "rider"
                    },
                    "get_address": true,
                    "get_timezone": true,
                    "updated_unix": 1586917896
                },
                "distance": 1439.4645850870015,
                "direction": {
                    "eta": 8,
                    "travel_dist": 2770
                },
                "timestamp_unix": 1586917897
            },
            {
                "object": {
                    "key": "destination_1",
                    "point": {
                        "lat": 39.71670913696289,
                        "lon": -104.95344543457031
                    },
                    "radius": 100,
                    "metadata": {
                        "type": "destination"
                    },
                    "get_address": true,
                    "get_timezone": true,
                    "expires_unix": 1587004296815040000,
                    "updated_unix": 1586917896
                },
                "distance": 5625.029340497144,
                "direction": {
                    "html_directions": "CjxoNT5EZXN0aW5hdGlvbjogVW5uYW1lZCBSb2FkLCBEZW52ZXIsIENPIDgwMjA5LCBVU0E8L2g1PkhlYWQgPGI+bm9ydGh3ZXN0PC9iPiAtIDExNSBmdDxicj5UdXJuIDxiPnJpZ2h0PC9iPiAtIDAuMyBtaTxicj5UdXJuIDxiPnJpZ2h0PC9iPiAtIDIyMCBmdDxicj5UdXJuIDxiPnJpZ2h0PC9iPiB0b3dhcmQgPGI+UGFyayBBdmUgVzwvYj4gLSAwLjEgbWk8YnI+VHVybiA8Yj5yaWdodDwvYj4gYXQgdGhlIDFzdCBjcm9zcyBzdHJlZXQgb250byA8Yj5QYXJrIEF2ZSBXPC9iPiAtIDAuMSBtaTxicj5UdXJuIDxiPmxlZnQ8L2I+IG9udG8gPGI+V2V3YXR0YSBTdDwvYj4gLSAwLjkgbWk8YnI+VHVybiA8Yj5sZWZ0PC9iPiBvbnRvIDxiPlNwZWVyIEJsdmQ8L2I+IC0gMS43IG1pPGJyPktlZXAgPGI+bGVmdDwvYj4gdG8gc3RheSBvbiA8Yj5TcGVlciBCbHZkPC9iPiAtIDEuMiBtaTxicj5Db250aW51ZSBvbnRvIDxiPkUgMXN0IEF2ZTwvYj4gLSAxLjIgbWk8YnI+VHVybiA8Yj5yaWdodDwvYj4gb250byA8Yj5TdGVlbGUgU3Q8L2I+IC0gMC4xIG1pPGJyPlR1cm4gPGI+cmlnaHQ8L2I+IGF0IHRoZSAxc3QgY3Jvc3Mgc3RyZWV0IGF0IDxiPkUgRWxsc3dvcnRoIEF2ZTwvYj4gLSAxNjcgZnQ8YnI+Q29udGludWUgc3RyYWlnaHQgLSAwLjIgbWk8YnI+",
                    "eta": 18,
                    "travel_dist": 9584
                },
                "timestamp_unix": 1586917897
            }
        ]
    }
}
```

HTML Directions are base-64 encoded

#### Pickup Rider
Update Lat/Lon to close to the rider. If the Driver and Riders Boundaries intersect, the trackers "inside" will be true to indicate they are at the same location

```go
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
```

Pretty JSON Response: 

```json
{
    "object": {
        "object": {
            "key": "driver_1",
            "point": {
                "lat": 39.74863815307619,
                "lon": -105.0076217651367
            },
            "radius": 100,
            "tracking": {
                "trackers": [
                    {
                        "target_object_key": "rider_1",
                        "track_distance": true,
                        "track_eta": true
                    },
                    {
                        "target_object_key": "destination_1",
                        "track_directions": true,
                        "track_distance": true,
                        "track_eta": true
                    }
                ]
            },
            "metadata": {
                "type": "driver"
            },
            "get_address": true,
            "get_timezone": true,
            "expires_unix": 1586919701045993000,
            "updated_unix": 1586919401
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
                    "key": "rider_1",
                    "point": {
                        "lat": 39.74863815307617,
                        "lon": -105.00762176513672
                    },
                    "radius": 100,
                    "tracking": {
                        "trackers": [
                            {
                                "target_object_key": "driver_1",
                                "track_distance": true,
                                "track_eta": true
                            }
                        ]
                    },
                    "metadata": {
                        "type": "rider"
                    },
                    "get_address": true,
                    "get_timezone": true,
                    "updated_unix": 1586919400
                },
                "distance": 2.666476830861055e-9,
                "inside": true,
                "direction": {},
                "timestamp_unix": 1586919401
            },
            {
                "object": {
                    "key": "destination_1",
                    "point": {
                        "lat": 39.71670913696289,
                        "lon": -104.95344543457031
                    },
                    "radius": 100,
                    "metadata": {
                        "type": "destination"
                    },
                    "get_address": true,
                    "get_timezone": true,
                    "expires_unix": 1587005800393783000,
                    "updated_unix": 1586919400
                },
                "distance": 5843.27595555179,
                "direction": {
                    "html_directions": "CjxoNT5EZXN0aW5hdGlvbjogVW5uYW1lZCBSb2FkLCBEZW52ZXIsIENPIDgwMjA5LCBVU0E8L2g1PkhlYWQgPGI+bm9ydGhlYXN0PC9iPiB0b3dhcmQgPGI+MTF0aCBTdDwvYj4vPHdici8+PGI+MTJ0aCBTdDwvYj4gLSAwLjEgbWk8YnI+Q29udGludWUgb250byA8Yj4xMXRoIFN0PC9iPi88d2JyLz48Yj4xMnRoIFN0PC9iPiAtIDM1OCBmdDxicj5UdXJuIDxiPmxlZnQ8L2I+IG9udG8gPGI+Q2hvcHBlciBDaXI8L2I+IC0gMzI4IGZ0PGJyPlR1cm4gPGI+cmlnaHQ8L2I+IG9udG8gPGI+U3BlZXIgQmx2ZDwvYj4gLSAxLjcgbWk8YnI+S2VlcCA8Yj5sZWZ0PC9iPiB0byBzdGF5IG9uIDxiPlNwZWVyIEJsdmQ8L2I+IC0gMS4yIG1pPGJyPkNvbnRpbnVlIG9udG8gPGI+RSAxc3QgQXZlPC9iPiAtIDEuMiBtaTxicj5UdXJuIDxiPnJpZ2h0PC9iPiBvbnRvIDxiPlN0ZWVsZSBTdDwvYj4gLSAwLjEgbWk8YnI+VHVybiA8Yj5yaWdodDwvYj4gYXQgdGhlIDFzdCBjcm9zcyBzdHJlZXQgYXQgPGI+RSBFbGxzd29ydGggQXZlPC9iPiAtIDE2NyBmdDxicj5Db250aW51ZSBzdHJhaWdodCAtIDAuMiBtaTxicj4=",
                    "eta": 13,
                    "travel_dist": 7548
                },
                "timestamp_unix": 1586919401
            }
        ]
    }
}
```

The driver has now entered the pickup location

You may use a similar workflow to dropoff the driver at the destination

It is assumed that the client will open streams(see below) and push updates to client devices(perhaps via websockets) 

### Other Functionality

#### Streaming Object Updates

```go
if err := client.Stream(context.Background(), &api.StreamRequest{
		ClientId:             "", //if youre client is scaling horizontally, add a clientID so that messages are streamed to only one instance, otherwise leave empty
		Keys:                 nil, //add object keys for targeted object streams, leave empty to stream all object updates
	}, func(ctx context.Context, stream geodb_go.ObjectStreamGetter, err error) { //this function is executed on each object in the stream as they are received
		if err != nil {
			log.Println(err.Error())
			return
		}
		obj := stream.GetObject() //an object detail has entered the stream
		log.Println(obj.String()) //do stuff with object
	}); err != nil {
		log.Fatal(err.Error())
		return
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
		Prefix: "driver_",
	})
	for _, object := range resp.Objects {
		fmt.Println(object.String())
	}
```

#### Get Objects by Regex

```go
resp, err := client.GetRegex(context.Background(), &api.GetRegexRequest{
		Regex: "driver_*",
	})
for _, object := range resp.Objects {
		fmt.Println(object.String())
	}
```

#### Scan a Geolocation Boundary
```go
//ScanBound scans a give geolocation boundary for objects, use regex/prefix methods to filter objects
	resp, err := client.ScanBound(context.Background(), &api.ScanBoundRequest{
		//a bound is like a circle on a map
		Bound: &api.Bound{
			Center: &api.Point{
				Lat: 39.74863815307617,
				Lon: -105.00762176513672,
			}, //center
			Radius: 5000,        //radius
		},
	})
```

#### Deleting Objects

```go
//Delete deletes an array of objects. if the first string is *, all objects will be dropped from the database
	resp, err := client.Delete(context.Background(), &api.DeleteRequest{
		Keys: []string{"rider_1"},
	})
```

## API Docs

- Client API Docs can be found [here](https://github.com/autom8ter/geodb-go/blob/master/DOCS.md)
- Profobuf Contract can be found [here](https://github.com/autom8ter/geodb/blob/master/api.proto)
