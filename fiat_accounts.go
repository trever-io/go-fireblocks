package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	LIST_FIAT_ACCOUNTS = "/v1/fiat_accounts"
)

type FiatAsset struct {
	Id      string `json:"id"`
	Balance string `json:"balance"`
}

type FiatAccount struct {
	Id      string       `json:"id"`
	Type    string       `json:"type"`
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Assets  []*FiatAsset `json:"assets"`
}

func (c *client) ListFiatAccounts(ctx context.Context) ([]*FiatAccount, error) {
	data, _, err := c.getRequest(ctx, LIST_FIAT_ACCOUNTS)
	if err != nil {
		return nil, fmt.Errorf("error during list fiat accounts: %w", err)
	}

	tmp := make([]*FiatAccount, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) RetrieveFiatAccount(ctx context.Context, id string) (*FiatAccount, error) {
	uri := fmt.Sprintf("%v/%v", LIST_FIAT_ACCOUNTS, id)
	data, _, err := c.getRequest(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("error during retrieve fiat account: %w", err)
	}

	tmp := new(FiatAccount)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}
