//go:generate godocdown -o README.md

package geodb_go

import (
	"context"
	"fmt"
	"github.com/autom8ter/geodb/gen/go/geodb"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/piotrkowalczuk/promgrpc/v3"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//Client is a geodb client
type Client struct {
	client api.GeoDBClient
	authFn func(ctx context.Context) context.Context
}

//ClientConfig holds configuration options for creating a Client
type ClientConfig struct {
	Host     string //geodb server url ex: localhost:8080
	Password string //optional - only necessary if basic auth is enabled on geodb server
	Metrics  bool   //register client side prometheus metrics interceptor
	Retry    bool   //register client side retry interceptor
}

//NewClient creates a new GeoDb client with the given Config
func NewClient(ctx context.Context, cfg *ClientConfig) (*Client, error) {
	unary := []grpc.UnaryClientInterceptor{}
	stream := []grpc.StreamClientInterceptor{}
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	if cfg.Metrics {
		var promInterceptor = promgrpc.NewInterceptor(promgrpc.InterceptorOpts{})
		if err := prometheus.DefaultRegisterer.Register(promInterceptor); err != nil {
			return nil, err
		}
		unary = append(unary, promInterceptor.UnaryClient())
		stream = append(stream, promInterceptor.StreamClient())
		opts = append(opts, grpc.WithStatsHandler(promInterceptor))
	}
	if cfg.Retry {
		unary = append(unary, grpc_retry.UnaryClientInterceptor())
		stream = append(stream, grpc_retry.StreamClientInterceptor())
	}
	if len(unary) > 0 {
		opts = append(opts, grpc.WithChainUnaryInterceptor(unary...))
	}
	if len(stream) > 0 {
		opts = append(opts, grpc.WithChainStreamInterceptor(stream...))
	}

	conn, err := grpc.DialContext(ctx, cfg.Host, opts...)
	if err != nil {
		return nil, err
	}
	client := &Client{client: api.NewGeoDBClient(conn)}
	if cfg.Password != "" {
		client.authFn = func(ctx context.Context) context.Context {
			md := metadata.Pairs("authorization", fmt.Sprintf("basic %v", cfg.Password))
			return metautils.NiceMD(md).ToOutgoing(ctx)
		}
	} else {
		client.authFn = func(ctx context.Context) context.Context {
			return ctx
		}
	}
	return client, nil
}

//ObjectStreamGetter pulls object details off of a stream
type ObjectStreamGetter interface {
	GetObject() *api.ObjectDetail
}

//StreamHandler executes logic on an object stream
type StreamHandler func(ctx context.Context, stream ObjectStreamGetter, err error)

//Ping - output: returns ok if server is healthy.
func (c Client) Ping(ctx context.Context) (*api.PingResponse, error) {
	return c.client.Ping(c.authFn(ctx), &api.PingRequest{})
}

//Set - input: an array of objects output: returns an indexed map of updated object details.
//objects are upserted in the order they are sent
func (c Client) Set(ctx context.Context, r *api.SetRequest) (*api.SetResponse, error) {
	return c.client.Set(c.authFn(ctx), r)
}

//Get - input: an array of object keys, output: returns an indexed map of current object details
func (c Client) Get(ctx context.Context, r *api.GetRequest) (*api.GetResponse, error) {
	return c.client.Get(c.authFn(ctx), r)
}

//GetRegex - input: a regex string, output: returns an indexed map of current object details with keys that match the regex pattern
func (c Client) GetRegex(ctx context.Context, r *api.GetRegexRequest) (*api.GetRegexResponse, error) {
	return c.client.GetRegex(c.authFn(ctx), r)
}

//GetPrefix - input: a prefix string, output: returns an indexed map of current object details with keys that have the given prefix
func (c Client) GetPrefix(ctx context.Context, r *api.GetPrefixRequest) (*api.GetPrefixResponse, error) {
	return c.client.GetPrefix(c.authFn(ctx), r)
}

//GetKeys -  input: none, output: returns all keys in database
func (c Client) GetKeys(ctx context.Context, r *api.GetKeysRequest) (*api.GetKeysResponse, error) {
	return c.client.GetKeys(c.authFn(ctx), r)
}

//GetRegexKeys -  input: a regex string, output: returns all keys in database that match the regex pattern
func (c Client) GetRegexKeys(ctx context.Context, r *api.GetRegexKeysRequest) (*api.GetRegexKeysResponse, error) {
	return c.client.GetRegexKeys(c.authFn(ctx), r)
}

//GetPrefixKeys - input: a prefix string, output: returns an array of of keys that have the given prefix
func (c Client) GetPrefixKeys(ctx context.Context, r *api.GetPrefixKeysRequest) (*api.GetPrefixKeysResponse, error) {
	return c.client.GetPrefixKeys(c.authFn(ctx), r)
}

//Delete -  input: an array of object key strings to delete, output: none
func (c Client) Delete(ctx context.Context, r *api.DeleteRequest) (*api.DeleteResponse, error) {
	return c.client.Delete(c.authFn(ctx), r)
}

//Stream -  input: a clientID(optional) and an array of object keys(optional),
//output: a stream of object details for realtime, targeted object geolocation updates
func (c Client) Stream(ctx context.Context, r *api.StreamRequest, handler StreamHandler) error {
	client, err := c.client.Stream(c.authFn(ctx), r)
	if err != nil {
		return err
	}
	for {
		if ctx.Err() != nil {
			break
		}
		resp, err := client.Recv()
		handler(ctx, resp, err)
	}
	return nil
}

//StreamRegex -  input: a clientID(optional) a regex string,
//output: a stream of object details for realtime, targeted object geolocation updates that match the regex pattern
func (c Client) StreamRegex(ctx context.Context, r *api.StreamRegexRequest, handler StreamHandler) error {
	client, err := c.client.StreamRegex(c.authFn(ctx), r)
	if err != nil {
		return err
	}
	for {
		if ctx.Err() != nil {
			break
		}
		resp, err := client.Recv()
		handler(ctx, resp, err)
	}
	return nil
}

//StreamPrefix -  input: a clientID(optional) a prefix string,
//output: a stream of object details for realtime, targeted object geolocation updates that match the prefix pattern
func (c Client) StreamPrefix(ctx context.Context, r *api.StreamPrefixRequest, handler StreamHandler) error {
	client, err := c.client.StreamPrefix(c.authFn(ctx), r)
	if err != nil {
		return err
	}
	for {
		if ctx.Err() != nil {
			break
		}
		resp, err := client.Recv()
		handler(ctx, resp, err)
	}
	return nil
}

//ScanBound -  input: a geolocation boundary, output: returns an indexed map of current object details that are within the boundary
func (c Client) ScanBound(ctx context.Context, r *api.ScanBoundRequest) (*api.ScanBoundResponse, error) {
	return c.client.ScanBound(c.authFn(ctx), r)
}

//ScanRegexBound -  input: a geolocation boundary, string-array of unique object ids(optional), output: returns an indexed map of current object details that have keys that match the regex and are within the boundary and
func (c Client) ScanRegexBound(ctx context.Context, r *api.ScanRegexBoundRequest) (*api.ScanRegexBoundResponse, error) {
	return c.client.ScanRegexBound(c.authFn(ctx), r)
}

//ScanPrefexBound -  input: a geolocation boundary, output: returns an indexed map of current object details that have keys that match the prefix and are within the boundary and
func (c Client) ScanPrefixBound(ctx context.Context, r *api.ScanPrefixBoundRequest) (*api.ScanPrefixBoundResponse, error) {
	return c.client.ScanPrefixBound(c.authFn(ctx), r)
}

//GetPoint can be used to get an addresses latitude/longitude
func (c Client) GetPoint(ctx context.Context, r *api.GetPointRequest) (*api.GetPointResponse, error) {
	return c.client.GetPoint(c.authFn(ctx), r)
}
