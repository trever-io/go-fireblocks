package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	LIST_TRANSACTION   = "/v1/transactions"
	CREATE_TRANSACTION = "/v1/transactions"
)

type TransactionRequest struct {
	AssetId     string `json:"asset_id"`
	Amount      string `json:"amount"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type TransactionResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type Transaction struct {
	Id                 string  `json:"id"`
	AssetId            string  `json:"assetId"`
	Amount             float64 `json:"amount"`
	NetworkFee         float64 `json:"networkFee"`
	CreatedAt          int64   `json:"createdAt"`
	Status             string  `json:"status"`
	DestinationAddress string  `json:"destinationAddress"`
}

type GetWithdrawalHistoryOptions struct {
	StartTime int64 `url:"start_time"`
	EndTime   int64 `url:"end_time"`
}

func (c *client) ListTransactions(ctx context.Context, opts *GetWithdrawalHistoryOptions) ([]*Transaction, error) {
	queryParameters := map[string]string{
		"after":  strconv.FormatInt(opts.StartTime, 10),
		"before": strconv.FormatInt(opts.EndTime, 10),
	}
	data, err := c.getRequestWithQuery(ctx, LIST_TRANSACTION, queryParameters)
	if err != nil {
		return nil, fmt.Errorf("error during list transactions: %w", err)
	}

	tmp := make([]*Transaction, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) CreateTransactions(ctx context.Context, req *TransactionRequest) (*TransactionResponse, error) {
	data, err := c.postRequest(ctx, CREATE_TRANSACTION, req)
	if err != nil {
		return nil, fmt.Errorf("error during create transaction: %w", err)
	}

	tmp := new(TransactionResponse)
	err = json.Unmarshal(data, tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}
