package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	LIST_EXTERNAL_WALLETS = "/v1/external_wallets"
)

type ExternalWalletAsset struct {
	Id             string `json:"id"`
	Status         string `json:"status"`
	ActivationTime string `json:"activationTime"`
	Address        string `json:"address"`
	Tag            string `json:"tag"`
}

type ExternalWallet struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	CustomerRefId string                 `json:"customerRefId"`
	Assets        []*ExternalWalletAsset `json:"assets"`
}

func (c *client) ListExternalWallets(ctx context.Context) ([]*ExternalWallet, error) {
	data, _, err := c.getRequest(ctx, LIST_EXTERNAL_WALLETS)
	if err != nil {
		return nil, fmt.Errorf("error during list external wallets: %w", err)
	}

	tmp := make([]*ExternalWallet, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}
