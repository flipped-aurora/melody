package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/coreos/etcd/clientv3"
)

// Code taken from https://github.com/go-kit/kit/blob/master/sd/etcd/client.go

var (
	// ErrNoConfig is the error to be returned when there is no config with the etcd namespace
	ErrNoHost = fmt.Errorf("get 0 host from etcd client: no host")
)

const defaultTTL = 3 * time.Second

// Client is a wrapper around the etcd client.
type Client interface {
	// GetEntries queries the given prefix in etcd and returns a slice
	// containing the values of all keys found, recursively, underneath that
	// prefix.
	GetEntries(prefix string) ([]string, error)

	// WatchPrefix watches the given prefix in etcd for changes. When a change
	// is detected, it will signal on the passed channel. Clients are expected
	// to call GetEntries to update themselves with the latest set of complete
	// values. WatchPrefix will always send an initial sentinel value on the
	// channel after establishing the watch, to ensure that clients always
	// receive the latest set of values. WatchPrefix will block until the
	// context passed to the NewClient constructor is terminated.
	WatchPrefix(prefix string, ch chan struct{})
}

type client struct {
	v3  *clientv3.Client
	ctx context.Context
}

// ClientOptions defines options for the etcd client. All values are optional.
// If any duration is not specified, a default of 3 seconds will be used.
type ClientOptions struct {
	Cert          string
	Key           string
	CACert        string
	DialTimeout   time.Duration
	DialKeepAlive time.Duration
}

// NewClient returns Client with a connection to the named machines. It will
// return an error if a connection to the cluster cannot be made. The parameter
// machines needs to be a full URL with schemas. e.g. "http://localhost:2379"
// will work, but "localhost:2379" will not.
func NewClient(ctx context.Context, machines []string, options ClientOptions) (Client, error) {
	if options.DialTimeout == 0 {
		options.DialTimeout = defaultTTL
	}
	if options.DialKeepAlive == 0 {
		options.DialKeepAlive = defaultTTL
	}

	var tlsCfg *tls.Config
	if options.Cert != "" && options.Key != "" {
		tlsCert, err := tls.LoadX509KeyPair(options.Cert, options.Key)
		if err != nil {
			return nil, err
		}
		tlsCfg := &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
		}
		if caCertCt, err := ioutil.ReadFile(options.CACert); err == nil {
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCertCt)
			tlsCfg.RootCAs = caCertPool
		}
	}

	ce, err := clientv3.New(clientv3.Config{
		Endpoints:         machines,
		DialTimeout:       options.DialTimeout,
		DialKeepAliveTime: options.DialKeepAlive,
		Context:           ctx,
		TLS:               tlsCfg,
	})

	if err != nil {
		return nil, err
	}

	return &client{
		v3:  ce,
		ctx: ctx,
	}, nil
}

// GetEntries implements the etcd Client interface.
func (c *client) GetEntries(key string) ([]string, error) {
	resp, err := c.v3.Get(c.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	// Special case. Note that it's possible that len(resp.Node.Nodes) == 0 and
	// resp.Node.Value is also empty, in which case the key is empty and we
	// should not return any entries.
	if len(resp.Kvs) == 0 && resp.Count == 0 {
		return nil, ErrNoHost
	}

	entries := make([]string, resp.Count)
	for i, node := range resp.Kvs {
		entries[i] = string(node.Value)
	}
	return entries, nil
}

// WatchPrefix implements the etcd Client interface.
func (c *client) WatchPrefix(prefix string, ch chan struct{}) {
	wch := c.v3.Watcher.Watch(c.ctx, prefix, clientv3.WithPrefix())
	ch <- struct{}{}
	for {
		select {
		case resp := <-wch:
			if resp.Canceled {
				return
			} else {
				ch <- struct{}{}
			}
		}
	}
}
