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
		CurrentBlockNumber(context.Context) (int64, error)
		GetAccountBalance(context.Context, string, string) (uint64, error)
		GetTRC20TokenSymbol(context.Context, string, string) (string, error)
		GetTRC20TokenDecimals(context.Context, string, string) (uint64, error)
		GetTRC20TokenBalance(context.Context, string, string) (uint64, error)
		GetTRC20SmartContract(context.Context, string) (*GetContractReply, error)
	}
)

var (
	DefaultClient = &http.Client{
		Timeout: time.Second * 15,
	}

	DefaultErrorHandler = func(res *http.Response, uri string) error {
		return nil
	}
)
