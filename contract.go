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
