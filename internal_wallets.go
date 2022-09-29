package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	LIST_INTERNAL_WALLETS  = "/v1/internal_wallets"
	CREATE_INTERNAL_WALLET = "/v1/internal_wallets"
)

type InternalWalletAsset struct {
	Id             string `json:"id"`
	Balance        string `json:"balance"`
	LockedAmount   string `json:"lockedAmount"`
	Status         string `json:"status"`
	ActivationTime string `json:"activationTime"`
	Address        string `json:"address"`
	Tag            string `json:"tag"`
}

type InternalWallet struct {
	Id            string                 `json:"id"`
	Name          string                 `json:"name"`
	CustomerRefId string                 `json:"customerRefId"`
	Assets        []*InternalWalletAsset `json:"assets"`
}

type CreateInternalWalletRequest struct {
	Name          string  `json:"name"`
	CustomerRefId *string `json:"customerRefId,omitempty"`
}

func (c *client) ListInternalWallets(ctx context.Context) ([]*InternalWallet, error) {
	data, _, err := c.getRequest(ctx, LIST_INTERNAL_WALLETS)
	if err != nil {
		return nil, fmt.Errorf("error during list internal wallets: %w", err)
	}

	tmp := make([]*InternalWallet, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) CreateInternalWallet(ctx context.Context, req *CreateInternalWalletRequest) (*InternalWallet, error) {
	data, _, err := c.postRequest(ctx, CREATE_INTERNAL_WALLET, req)
	if err != nil {
		return nil, fmt.Errorf("error during create internal wallet: %w", err)
	}

	tmp := new(InternalWallet)
	err = json.Unmarshal(data, tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}
