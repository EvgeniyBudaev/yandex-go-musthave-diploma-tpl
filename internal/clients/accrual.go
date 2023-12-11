package clients

import (
	wrapError "github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/errors"
	"github.com/EvgeniyBudaev/yandex-go-musthave-diploma-tpl/internal/storage"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type ClientAccrual struct {
	client    *resty.Client
	serverURL string
}

func NewClientAccrual(client *resty.Client, serverURL string) *ClientAccrual {
	return &ClientAccrual{client: client, serverURL: serverURL}
}

func (ca *ClientAccrual) CheckAccrual(number string) (*storage.AccrualDto, error) {
	accrual := storage.AccrualDto{}
	response, err := ca.client.R().
		SetResult(&accrual).
		SetRawPathParam("number", number).
		Get(ca.serverURL + "/api/orders/{number}")
	if response.StatusCode() == http.StatusTooManyRequests {
		return nil, wrapError.ErrTooManyRequests
	}
	if strings.Contains(response.Status(), http.StatusText(http.StatusNoContent)) {
		return nil, wrapError.ErrNoContent
	}
	if err != nil {
		return nil, err
	}
	return &accrual, nil
}
