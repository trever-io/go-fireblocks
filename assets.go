package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	SUPPORTED_ASSETS = "/v1/supported_assets"
	FIAT_ASSET_TYPE  = "FIAT"
)

type Asset struct {
	Id              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	ContractAddress string  `json:"contractAddress"`
	NativeAsset     string  `json:"nativeAsset"`
	Decimals        float64 `json:"decimals"`
}

func (c *client) GetSupportedAssets(ctx context.Context) ([]*Asset, error) {
	data, _, err := c.getRequest(ctx, SUPPORTED_ASSETS)
	if err != nil {
		return nil, fmt.Errorf("error during list supported assets: %w", err)
	}

	tmp := make([]*Asset, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) GetFiatAssets(ctx context.Context) ([]*Asset, error) {
	data, _, err := c.getRequest(ctx, SUPPORTED_ASSETS)
	if err != nil {
		return nil, fmt.Errorf("error during list supported assets: %w", err)
	}

	tmp := make([]*Asset, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	fiats := make([]*Asset, 0)
	for _, asset := range tmp {
		if asset.Type == FIAT_ASSET_TYPE {
			fiats = append(fiats, asset)
		}
	}
	return fiats, nil
}

func (c *client) GetFiatAssetIds(ctx context.Context) ([]string, error) {
	data, _, err := c.getRequest(ctx, SUPPORTED_ASSETS)
	if err != nil {
		return nil, fmt.Errorf("error during list supported assets: %w", err)
	}

	tmp := make([]*Asset, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	fiats := make([]string, 0)
	for _, asset := range tmp {
		if asset.Type == FIAT_ASSET_TYPE {
			fiats = append(fiats, asset.Id)
		}
	}
	return fiats, nil
}
