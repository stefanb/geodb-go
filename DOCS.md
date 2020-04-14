# geodb_go
--
    import "github.com/autom8ter/geodb-go"


## Usage

#### type Client

```go
type Client struct {
}
```

Client is a geodb client

#### func  NewClient

```go
func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error)
```
NewClient creates a new GeoDb client with the given Config

#### func (Client) Delete

```go
func (c Client) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error)
```
Delete - input: an array of object key strings to delete, output: none

#### func (Client) Get

```go
func (c Client) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error)
```
Get - input: an array of object keys, output: returns an indexed map of current
object details

#### func (Client) GetKeys

```go
func (c Client) GetKeys(ctx context.Context, r *api.GetKeysRequest) (*api.GetKeysResponse, error)
```
GetKeys - input: none, output: returns all keys in database

#### func (Client) GetPoint

```go
func (c Client) GetPoint(ctx context.Context, r *api.GetPointRequest) (*api.GetPointResponse, error)
```
GetPoint can be used to get an addresses latitude/longitude

#### func (Client) GetPrefix

```go
func (c Client) GetPrefix(ctx context.Context, r *api.GetPrefixRequest) (*api.GetPrefixResponse, error)
```
GetPrefix - input: a prefix string, output: returns an indexed map of current
object details with keys that have the given prefix

#### func (Client) GetPrefixKeys

```go
func (c Client) GetPrefixKeys(ctx context.Context, r *api.GetPrefixKeysRequest) (*api.GetPrefixKeysResponse, error)
```
GetPrefixKeys - input: a prefix string, output: returns an array of of keys that
have the given prefix

#### func (Client) GetRegex

```go
func (c Client) GetRegex(ctx context.Context, r *api.GetRegexRequest) (*api.GetRegexResponse, error)
```
GetRegex - input: a regex string, output: returns an indexed map of current
object details with keys that match the regex pattern

#### func (Client) GetRegexKeys

```go
func (c Client) GetRegexKeys(ctx context.Context, r *api.GetRegexKeysRequest) (*api.GetRegexKeysResponse, error)
```
GetRegexKeys - input: a regex string, output: returns all keys in database that
match the regex pattern

#### func (Client) Ping

```go
func (c Client) Ping(ctx context.Context) (*api.PingResponse, error)
```
Ping - output: returns ok if server is healthy.

#### func (Client) ScanBound

```go
func (c Client) ScanBound(ctx context.Context, r *api.ScanBoundRequest) (*api.ScanBoundResponse, error)
```
ScanBound - input: a geolocation boundary, output: returns an indexed map of
current object details that are within the boundary

#### func (Client) ScanPrefixBound

```go
func (c Client) ScanPrefixBound(ctx context.Context, r *api.ScanPrefixBoundRequest) (*api.ScanPrefixBoundResponse, error)
```
ScanPrefexBound - input: a geolocation boundary, output: returns an indexed map
of current object details that have keys that match the prefix and are within
the boundary and

#### func (Client) ScanRegexBound

```go
func (c Client) ScanRegexBound(ctx context.Context, r *api.ScanRegexBoundRequest) (*api.ScanRegexBoundResponse, error)
```
ScanRegexBound - input: a geolocation boundary, string-array of unique object
ids(optional), output: returns an indexed map of current object details that
have keys that match the regex and are within the boundary and

#### func (Client) Set

```go
func (c Client) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error)
```
Set - input: an array of objects output: returns an indexed map of updated
object details. objects are upserted in the order they are sent

#### func (Client) Stream

```go
func (c Client) Stream(ctx context.Context, r *api.StreamRequest, handler StreamHandler) error
```
Stream - input: a clientID(optional) and an array of object keys(optional),
output: a stream of object details for realtime, targeted object geolocation
updates

#### func (Client) StreamPrefix

```go
func (c Client) StreamPrefix(ctx context.Context, r *api.StreamPrefixRequest, handler StreamHandler) error
```
StreamPrefix - input: a clientID(optional) a prefix string, output: a stream of
object details for realtime, targeted object geolocation updates that match the
prefix pattern

#### func (Client) StreamRegex

```go
func (c Client) StreamRegex(ctx context.Context, r *api.StreamRegexRequest, handler StreamHandler) error
```
StreamRegex - input: a clientID(optional) a regex string, output: a stream of
object details for realtime, targeted object geolocation updates that match the
regex pattern

#### type ClientConfig

```go
type ClientConfig struct {
	Host     string //geodb server url ex: localhost:8080
	Password string //optional - only necessary if basic auth is enabled on geodb server
	Metrics  bool   //register client side prometheus metrics interceptor
	Retry    bool   //register client side retry interceptor
}
```

ClientConfig holds configuration options for creating a Client

#### type ObjectStreamGetter

```go
type ObjectStreamGetter interface {
	GetObject() *api.ObjectDetail
}
```

ObjectStreamGetter pulls object details off of a stream

#### type StreamHandler

```go
type StreamHandler func(ctx context.Context, stream ObjectStreamGetter, err error)
```

StreamHandler executes logic on an object stream
