package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	TRANSACTION               = "/v1/transactions"
	DEFAULT_TRANSACTION_LIMIT = 200
)

type TransactionRequest struct {
	AssetId     string                      `json:"assetId"`
	Amount      string                      `json:"amount"`
	Source      TransferPeerPath            `json:"source"`
	Destination DestinationTransferPeerPath `json:"destination"`
}

type TransactionResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type ResponseHeader struct {
	NextPage string `json:"next_page"`
}

type Transaction struct {
	Id                 string                   `json:"id"`
	AssetId            string                   `json:"assetId"`
	Source             TransferPeerPathResponse `json:"source"`
	Destination        TransferPeerPathResponse `json:"destination"`
	Amount             float64                  `json:"amount"`
	NetworkFee         float64                  `json:"networkFee"`
	ServiceFee         float64                  `json:"serviceFee"`
	CreatedAt          int64                    `json:"createdAt"`
	Status             string                   `json:"status"`
	DestinationAddress string                   `json:"destinationAddress"`
}

type GetHistoryOptions struct {
	StartTime int64 `url:"start_time"`
	EndTime   int64 `url:"end_time"`
}

type TransferPeerPath struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type DestinationTransferPeerPath struct {
	Type           string `json:"type"`
	Id             string `json:"id,omitempty"`
	OneTimeAddress string `json:"oneTimeAddress,omitempty"`
}

type TransferPeerPathResponse struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	SubType string `json:"subType"`
}

func (c *client) ListTransactions(ctx context.Context, opts *GetHistoryOptions) ([]*Transaction, error) {
	queryParameters := map[string]string{
		"after":  strconv.FormatInt(opts.StartTime, 10),
		"before": strconv.FormatInt(opts.EndTime, 10),
	}

	header := &ResponseHeader{
		NextPage: fmt.Sprintf("%v%v", BASE_URL, TRANSACTION),
	}
	transactions := make([]*Transaction, 0)

	for {
		nextUrl := strings.ReplaceAll(header.NextPage, BASE_URL, "")
		data, nextPage, err := c.getRequestWithQuery(ctx, nextUrl, queryParameters)
		if err != nil {
			return nil, fmt.Errorf("error during list transactions: %w", err)
		}

		tmp := make([]*Transaction, 0)
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		if len(tmp) < DEFAULT_TRANSACTION_LIMIT {
			transactions = append(transactions, tmp...)
			break
		}

		h := new(ResponseHeader)
		err = json.Unmarshal(nextPage, h)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling header response: %w", err)
		}

		header.NextPage = h.NextPage
		transactions = append(transactions, tmp...)

	}

	return transactions, nil
}

func (c *client) GetTransactionById(ctx context.Context, id string) (*Transaction, error) {
	uri := fmt.Sprintf("%v/%v", TRANSACTION, id)
	data, _, err := c.getRequest(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("error during list transaction: %w", err)
	}

	tmp := new(Transaction)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) CreateTransaction(ctx context.Context, req *TransactionRequest) (*TransactionResponse, error) {
	data, _, err := c.postRequest(ctx, TRANSACTION, req)
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
