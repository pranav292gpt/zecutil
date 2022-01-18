package rpcclient

import (
	"encoding/base64"

	"github.com/ybbus/jsonrpc"
)

type ConnConfig struct {
	// Host is the IP address and port of the RPC server you want to connect
	// to.
	Host string

	// Endpoint is the websocket endpoint on the RPC server.  This is
	// typically "ws".
	Endpoint string

	// User is the username to use to authenticate to the RPC server.
	User string

	// Pass is the passphrase to use to authenticate to the RPC server.
	Pass string

	// HTTPPostMode instructs the client to run using multiple independent
	// connections issuing HTTP POST requests instead of using the default
	// of websockets.  Websockets are generally preferred as some of the
	// features of the client such notifications only work with websockets,
	// however, not all servers support the websocket extensions, so this
	// flag can be set to true to use basic HTTP POST requests instead.
	HTTPPostMode bool
}

type Client struct {
	// config holds the connection configuration assoiated with this client.
	config *ConnConfig

	rpcClient *jsonrpc.RPCClient
}

func New(config *ConnConfig) (*Client, error) {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(config.User + ":" + config.Pass))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+config.Host,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})

	return &Client{config: config, rpcClient: &rpcClient}, nil
}
