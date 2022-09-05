package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	SUPPORTED_ASSETS = "/v1/supported_assets"
)

type Asset struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	ContractAddress string `json:"contractAddress"`
	NativeAsset     string `json:"nativeAsset"`
	Decimals        string `json:"decimals"`
}

func (c *client) GetSupportedAssets(ctx context.Context) ([]*Asset, error) {
	data, err := c.getRequest(ctx, SUPPORTED_ASSETS)
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
