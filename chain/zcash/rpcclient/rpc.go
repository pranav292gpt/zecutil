package rpcclient

type zcashConf struct {
	testNet     bool
	rpcUser     string
	rpcPassword string
	rpcPort     string
}

// GetBlockchainInfo return the zcashd rpc `getblockchaininfo` status
// https://zcash-rpc.github.io/getblockchaininfo.html
type GetBlockchainInfo struct {
	Chain                string  `json:"chain"`
	Blocks               int     `json:"blocks"`
	Difficulty           float64 `json:"difficulty"`
	VerificationProgress float64 `json:"verificationprogress"`
	SizeOnDisk           float64 `json:"size_on_disk"`
}

// GetInfo Returns an object containing various state info.
// https://zcash-rpc.github.io/getinfo.html
type GetInfo struct {
	Version int `json:"version"`
}

// GetMemPoolInfo return the zcashd rpc `getmempoolinfo`
// https://zcash-rpc.github.io/getmempoolinfo.html
type GetMemPoolInfo struct {
	Size  float64 `json:"size"`
	Bytes float64 `json:"bytes"`
	Usage float64 `json:"usage"`
}

// GetNetworkInfoResult return the zcashd rpc `getnetworkinfo`
// https://zcash.github.io/rpc/getnetworkinfo.html
type GetNetworkInfo struct {
	Version         uint32         `json:"version"`
	Subversion      string         `json:"subversion"`
	Protocolversion uint32         `json:"protocolversion"`
	Localservices   string         `json:"localservices"`
	Timeoffset      int64          `json:"timeoffset"`
	Connections     uint32         `json:"connections"`
	Networks        []Network      `json:"networks"`
	RelayFee        float64        `json:"relayfee"`
	Localaddresses  []LocalAddress `json:"localaddresses"`
	Warnings        string         `json:"warnings"`
}

// Network network info
type Network struct {
	Name                      string `json:"name"`
	Limited                   bool   `json:"limited"`
	Reachable                 bool   `json:"reachable"`
	Proxy                     string `json:"proxy"`
	ProxyRandomizeCredentials bool   `json:"proxy_randomize_credentials"`
}

// LocalAddress local address
type LocalAddress struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	Score   int    `json:"score"`
}

type Unspent struct {
	TxID          string  `json:"txid"`
	VOut          uint32  `json:"vout"`
	Generated     bool    `json:"generated"`
	Address       string  `json:"address"`
	ScriptPubKey  string  `json:"scriptPubKey"`
	Amount        float64 `json:"amount"`
	AmountZat     float64 `json:"amountZat"`
	Confirmations int     `json:"confirmations"`
	RedeemScript  string  `json:"redeemScript"`
	Spendable     bool    `json:"spendable"`
}

// ZGetTotalBalance return the node's wallet balances
// https://zcash-rpc.github.io/z_gettotalbalance.html
type ZGetTotalBalance struct {
	Transparent string `json:"transparent"`
	Private     string `json:"private"`
	Total       string `json:"total"`
}

// GetPeerInfo Returns data about each connected network node
// https://zcash-rpc.github.io/getpeerinfo.html
type GetPeerInfo []PeerInfo

type PeerInfo struct {
	ID             int     `json:"id"`
	Addr           string  `json:"addr"`
	AddrLocal      string  `json:"addrlocal"`
	Services       string  `json:"services"`
	LastSend       int     `json:"lastsend"`
	LastRecv       int     `json:"lastrecv"`
	BytesSent      int     `json:"bytessent"`
	BytesRecv      int     `json:"bytesrecv"`
	Conntime       int     `json:"conntime"`
	Timeoffset     int     `json:"timeoffset"`
	PingTime       float64 `json:"pingtime"`
	PingWait       float64 `json:"pingwait"`
	Version        int     `json:"version"`
	Subver         string  `json:"subver"`
	Inbound        bool    `json:"inbound"`
	Startingheight int     `json:"startingheight"`
	Banscore       int     `json:"banscore"`
	SyncedHeaders  int     `json:"synced_headers"`
	SyncedBlocks   int     `json:"synced_blocks"`
}

// GetChainTips Return information about all known tips in the block tree
// https://zcash-rpc.github.io/getchaintips.html
type GetChainTips []ChainTip

type ChainTip struct {
	Hash      string `json:"hash"`
	Height    int    `json:"height"`
	Branchlen int    `json:"branchlen"`
	Status    string `json:"status"`
}

// GetDeprecationInfo Returns an object containing current version and deprecation block height. Applicable only on mainnet.
// https://zcash-rpc.github.io/getdeprecationinfo.html
type GetDeprecationInfo struct {
	Version           int    `json:"version"`
	Subversion        string `json:"subversion"`
	DeprecationHeight int    `json:"deprecationheight"`
}

type Block struct {
	Hash              string        `json:"hash"`
	Confirmations     int           `json:"confirmations"`
	Size              int           `json:"size"`
	Height            int           `json:"height"`
	Version           int           `json:"version"`
	MerkleRoot        string        `json:"merkleroot"`
	FinalSaplingRoot  string        `json:"finalsaplingroot"`
	TX                []Transaction `json:"tx"`
	Time              int64         `json:"time"`
	Nonce             string        `json:"nonce"`
	Difficulty        float64       `json:"difficulty"`
	PreviousBlockHash string        `json:"previousblockhash"`
	NextBlockHash     string        `json:"nextblockhash"`
	ValuePools        []ValuePool   `json:"valuePools"`
}

func (b Block) NumberofTransactions() int {
	return len(b.TX)
}

func (b Block) TransactionTypes() (tTXs, sTXs int) {
	for _, tx := range b.TX {
		// If all 3 fields are empty, the transaction is transparent
		if len(tx.VJoinSplit) > 0 ||
			len(tx.VShieldedOutput) > 0 ||
			len(tx.VShieldedSpend) > 0 {
			tTXs++
		} else {
			// Otherwise, it's a shielded transaction
			sTXs++
		}
	}
	return tTXs, sTXs
}

// Transaction describes a zcash tranaction
type Transaction struct {
	Hex             string                   `json:"hex"`
	Txid            string                   `json:"txid"`
	Version         int                      `json:"version"`
	LockTime        int                      `json:"locktime"`
	ExpiryHeight    int                      `json:"expirtheight"`
	VIn             []VIn                    `json:"vin"`
	VOut            []VOut                   `json:"vout"`
	VJoinSplit      []VJoinSplitTX           `json:"vjoinsplit"`
	ValueBalance    float64                  `json:"valueBalance"`
	VShieldedSpend  []map[string]interface{} `json:"vShieldedSpend"`
	VShieldedOutput []map[string]interface{} `json:"vShieldedOutput"`
}

// TransparentInAndOut return if there are transparent
// inputs and outputs
func (t Transaction) TransparentInAndOut() bool {
	return len(t.VIn) > 0 && len(t.VOut) > 0
}

// IsTransparent returns if the transaction contains
// no shielded addresses
func (t Transaction) IsTransparent() bool {
	return t.TransparentInAndOut() &&
		len(t.VJoinSplit) == 0 &&
		t.ValueBalance == 0 &&
		len(t.VShieldedSpend) == 0 &&
		len(t.VShieldedSpend) == 0
}

// ContainsSprout returns if a transaction contains
// sprout transaction data
func (t Transaction) ContainsSprout() bool {
	return len(t.VJoinSplit) > 0
}

// ContainsSapling returns if a transaction contains
// sapling transaction data
// Check that there is a valueBalance value (positive or negative)
// Check that there is data for either VShieldedSpend or VShieldedOutput
func (t Transaction) ContainsSapling() bool {
	return t.ValueBalance != 0 && (len(t.VShieldedSpend) > 0 ||
		len(t.VShieldedOutput) > 0)
}

// IsShielded returns if the transaction contains
// no transparent addresses
func (t Transaction) IsShielded() bool {
	return !t.TransparentInAndOut() &&
		(t.ContainsSprout() || t.ContainsSapling())
}

// IsMixed returns if the transaction contains
// transparent addresses and shielded transaction data
func (t Transaction) IsMixed() bool {
	tInOrOut := len(t.VIn) > 0 || len(t.VOut) > 0
	return tInOrOut &&
		(t.ContainsSprout() || t.ContainsSapling())
}

type VIn struct {
	Coinbase  string `json:"coinbase"`
	TxID      string `json:"txid"`
	VOut      int    `json:"vout"`
	ScriptSig ScriptSig
	Sequence  int `json:"sequemce"`
}

// IsCoinBase returns a bool to show if a Vin is a Coinbase one or not.
func (v *VIn) IsCoinBase() bool {
	return len(v.Coinbase) > 0
}

type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}
type VOut struct {
	Value        float64
	N            int
	ScriptPubKey ScriptPubKey
}
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}
type VJoinSplitTX struct {
	VPubOld float64 `json:"vpub_old"`
	VPubNew float64 `json:"vpub_new"`
}
type ValuePool struct {
	ID            string  `json:"id"`
	Monitored     bool    `json:"monitored"`
	ChainValue    float64 `json:"chainValue"`
	ChainValueZat float64 `json:"chainValueZat"`
	ValueDelta    float64 `json:"valueDelta"`
	ValueDeltaZat float64 `json:"valueDeltaZat"`
}

type TXOutSetInfo struct {
	Height       int     `json:"height"`
	BestBlock    string  `json:"bestblock"`
	Transactions int     `json:"transactions"`
	TXOuts       int     `json:"txouts"`
	TotalAmount  float64 `json:"total_amount"`
}

// ScriptPubKeyResult models the scriptPubKey data of a tx script.  It is
// defined separately since it is used by multiple commands.
type ScriptPubKeyResult struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex,omitempty"`
	ReqSigs   int32    `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

type GetBlockVerboseResult struct {
	Hash             string        `json:"hash"`
	Confirmations    int           `json:"confirmations"`
	Size             int           `json:"size"`
	Height           int64         `json:"height"`
	Version          int           `json:"version"`
	MerkleRoot       string        `json:"merkleroot"`
	FinalSaplingRoot string        `json:"finalsaplingroot"`
	FinalOrchardRoot string        `json:"finalorchardroot"`
	Tx               []Transaction `json:"tx"`
	Time             int64         `json:"time"`
	Nonce            uint32        `json:"nonce"`
	Bits             string        `json:"bits"`
	Difficulty       float64       `json:"difficulty"`
	PreviousHash     string        `json:"previousblockhash"`
	NextHash         string        `json:"nextblockhash,omitempty"`
}

//type Unspent struct {
//	TxID          string  `json:"txid"`
//	VOut          int     `json:"vout"`
//	Generated     bool    `json:"generated"`
//	Address       string  `json:"address"`
//	ScriptPubKey  string  `json:"scriptPubKey"`
//	Amount        float64 `json:"amount"`
//	AmountZat     uint64  `json:"amountZat"`
//	Confirmations int     `json:"confirmations"`
//	RedeemScript  string  `json:"redeemScript"`
//	Spendable     bool    `json:"spendable"`
//}
//
//type ListUnspent struct {
//
//}

type RawMemPool struct {
}
