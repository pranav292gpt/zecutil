package rpcclient

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/wire"
)

func (c *Client) SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (string, error) {
	txHex := ""
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return "", err
		}
		txHex = hex.EncodeToString(buf.Bytes())
	}
	return c.SendRawTransactionCmd(txHex, allowHighFees)
}
