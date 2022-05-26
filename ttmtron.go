package ttmtron

import (
	"context"

	"github.com/pkg/errors"
	"github.com/trustwallet/go-libs/client"
)

// Init TronRequest lib
func Init(baseURL string) *TronRequest {
	return &TronRequest{client.Request{
		BaseURL:          baseURL,
		Headers:          make(map[string]string),
		HttpClient:       DefaultClient,
		HttpErrorHandler: DefaultErrorHandler,
	}}
}

// CurrentBlockNumber return current block number
func (c *TronRequest) CurrentBlockNumber(ctx context.Context) (int64, error) {
	//log := logger.FromContext(ctx).WithField("m", "Client::CurrentBlockNumber")
	//log.Debugf("Client::CurrentBlockNumber:: ")

	var block Block
	err := c.Post(&block, "wallet/getnowblock", nil)

	return block.BlockHeader.Data.Number, errors.Wrap(err,"unable to post wallet/getnowblock")
}
