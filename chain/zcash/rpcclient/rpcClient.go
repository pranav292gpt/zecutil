package rpcclient

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"github.com/ybbus/jsonrpc"
)

type ConnConfig struct {
	// Host is the IP address and port of the RPC server you want to connect
	// to.
	Host string

	// Endpoint is the websocket endpoint on the RPC server.  This is
	// typically "ws".
	//Endpoint string

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
	//HTTPPostMode bool
}

type Client struct {
	// config holds the connection configuration assoiated with this client.
	config *ConnConfig

	rpcClient jsonrpc.RPCClient
}

func New(config *ConnConfig) (*Client, error) {
	basicAuth := base64.StdEncoding.EncodeToString([]byte(config.User + ":" + config.Pass))
	rpcClient := jsonrpc.NewClientWithOpts("http://"+config.Host,
		&jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"Authorization": "Basic " + basicAuth,
			}})

	return &Client{config: config, rpcClient: rpcClient}, nil
}

func (c *Client) GetInfo() *GetInfo {
	var info *GetInfo
	if err := c.rpcClient.CallFor(&info, "getinfo"); err != nil {
		log.Print("Error calling getinfo", err)
		return nil
	}
	return info
}

func (c *Client) GetBlockchainInfo() *GetBlockchainInfo {
	var blockInfo *GetBlockchainInfo
	if err := c.rpcClient.CallFor(&blockInfo, "getblockchaininfo"); err != nil {
		log.Print("Error calling getblockchaininfo", err)
		return nil
	}
	return blockInfo
}

func (c *Client) GetBlockCount() (int64, error) {
	var height int64
	err := c.rpcClient.CallFor(&height, "getblockcount")
	return height, err
}

func (c *Client) GetBlockHash(height int64) (string, error) {
	var hash string
	err := c.rpcClient.CallFor(&hash, "getblockhash", height)
	return hash, err
}

func (c *Client) GetNetworkInfo() (*GetNetworkInfo, error) {
	var networkInfo *GetNetworkInfo
	if err := c.rpcClient.CallFor(&networkInfo, "getnetworkinfo"); err != nil {
		return nil, err
	}
	return networkInfo, nil
}

// todo , asl Pranav if it's Ok to return upspent as an array
func (c *Client) ListUnspent() []Unspent {
	unspent := []Unspent{}
	if err := c.rpcClient.CallFor(&unspent, "listunspent"); err != nil {
		log.Println("Error calling listunspent ", err)
		return nil
	}
	return unspent
}

func (c *Client) ListUnspentMinMaxAddresses(minconf int, maxconf int, addresses []string) ([]Unspent, error) {
	unspent := []Unspent{}
	if err := c.rpcClient.CallFor(&unspent, "listunspent", minconf, maxconf, addresses); err != nil {
		return nil, err
	}
	return unspent, nil
}

func (c *Client) GetRawTransaction(txid string) (string, error) {
	var rawtx string
	err := c.rpcClient.CallFor(&rawtx, "getrawtransaction", txid)
	return rawtx, err
}

func (c *Client) GetRawTransactionVerbose(txid string) (*Transaction, error) {
	var rawtx *Transaction
	err := c.rpcClient.CallFor(&rawtx, "getrawtransaction", txid, 1)
	return rawtx, err
}

func (c *Client) GetRawMempool() ([]*chainhash.Hash, error) {
	var txHashStrs []string
	err := c.rpcClient.CallFor(&txHashStrs, "getrawmempool", "true")
	if err != nil {
		return nil, err
	}

	// Create a slice of ShaHash arrays from the string slice.
	txHashes := make([]*chainhash.Hash, 0, len(txHashStrs))
	for _, hashStr := range txHashStrs {
		txHash, err := chainhash.NewHashFromStr(hashStr)
		if err != nil {
			return nil, err
		}
		txHashes = append(txHashes, txHash)
	}

	return txHashes, nil
}

func (c *Client) GetBlockVerboseTx(hash string) (*GetBlockVerboseResult, error) {
	var result *GetBlockVerboseResult
	err := c.rpcClient.CallFor(&result, "getblock", hash, 2)
	return result, err
}

func (c *Client) SendRawTransactionCmd(hexstring string, allowhighfees bool) (string, error) {
	var result string
	err := c.rpcClient.CallFor(&result, "getblock", hexstring, allowhighfees)
	return result, err
}

func (c *Client) GetMempoolEntry(txID string) (*string, error) {
	var result []string
	err := c.rpcClient.CallFor(&result, "getrawmempool", txID)
	if err != nil {
		return nil, err
	}
	for _, tx := range result {
		if tx == txID {
			return &txID, nil
		}
	}
	return nil, fmt.Errorf("unable to finc txhash in mempool")
}

func (c *Client) GetBestBlockHash() (string, error) {
	var result string
	err := c.rpcClient.CallFor(&result, "getbestblockhash")
	return result, err
}
