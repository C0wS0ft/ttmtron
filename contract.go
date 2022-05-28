package ttmtron

import (
	"context"
	"net/http"
	"time"

	"github.com/trustwallet/go-libs/client"
)

type (
	TronRequest struct {
		client.Request
	}

	TTMTron interface {
		CurrentBlockNumber(context.Context) (uint64, error)
		GetBlockByNumber(context.Context, uint64) (*Block, error)
		GetBlockByLimitNext(context.Context, uint64, uint64) ([]*Block, error)
		GetAccountBalance(context.Context, string, string) (uint64, error)
		GetTRC20TokenSymbol(context.Context, string, string) (string, error)
		GetTRC20TokenDecimals(context.Context, string, string) (uint64, error)
		GetTRC20TokenBalance(context.Context, string, string) (uint64, error)
		GetTRC20SmartContract(context.Context, string) (*GetContractReply, error)

		GetTRC10TokenInfoByID(context.Context, uint64) (*GetAssetIssueByIDReply, error)
	}
)

var (
	DefaultClient = &http.Client{ //nolint:gochecknoglobals
		Timeout: time.Second * 15,
	}

	DefaultErrorHandler = func(res *http.Response, uri string) error { //nolint:gochecknoglobals
		return nil
	}
)
