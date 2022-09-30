package fireblocks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	LIST_VAULT_ACCOUNTS = "/v1/vault/accounts_paged"
	VAULT_ACCOUNTS      = "/v1/vault/accounts"
	VAULT_ASSETS        = "/v1/vault/assets"
	DEFAULT_VAULT_LIMIT = 200
)

type VaultAccountsPagedResponse struct {
	Accounts    []VaultAccount `json:"accounts"`
	Paging      Paging         `json:"paging"`
	PreviousUrl string         `json:"previousUrl"`
	NextUrl     string         `json:"nextUrl"`
}

type VaultAccount struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	HiddenOnUI bool         `json:"hiddenOnUI"`
	Assets     []VaultAsset `json:"assets"`
}

type VaultAsset struct {
	Id        string `json:"id"`
	Total     string `json:"total"`
	Available string `json:"available"`
	Pending   string `json:"pending"`
}

type AssetInformation struct {
	Id        string  `json:"id"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Pending   float64 `json:"pending"`
}

type CreateWalletInVault struct {
	VaultAccountId string `json:"vaultAccountId"`
	AssetId        string `json:"assetId"`
}

type CreateVaultAssetResponse struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

type Paging struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

func (c *client) ListVaultAccounts(ctx context.Context) (*VaultAccountsPagedResponse, error) {
	resp := &VaultAccountsPagedResponse{
		Accounts:    []VaultAccount{},
		PreviousUrl: "",
		NextUrl:     fmt.Sprintf("%v%v", BASE_URL, LIST_VAULT_ACCOUNTS),
	}

	for {
		nextUrl := strings.ReplaceAll(resp.NextUrl, BASE_URL, "")
		data, _, err := c.getRequest(ctx, nextUrl)
		if err != nil {
			return nil, fmt.Errorf("error during list vault accounts: %w", err)
		}

		tmp := &VaultAccountsPagedResponse{}
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		if len(tmp.Accounts) < DEFAULT_VAULT_LIMIT {
			resp.Accounts = append(resp.Accounts, tmp.Accounts...)
			break
		}

		resp.Accounts = append(resp.Accounts, tmp.Accounts...)
		resp.NextUrl = tmp.NextUrl

	}

	return resp, nil
}

func (c *client) CreateWalletInVault(ctx context.Context, req *CreateWalletInVault) (*CreateVaultAssetResponse, error) {
	uri := fmt.Sprintf("%v/%v/%v", VAULT_ACCOUNTS, req.VaultAccountId, req.AssetId)
	data, _, err := c.postRequest(ctx, uri, req)
	if err != nil {
		return nil, fmt.Errorf("error during create wallet in vault: %w", err)
	}

	tmp := new(CreateVaultAssetResponse)
	err = json.Unmarshal(data, tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) GetBalanceByAsset(ctx context.Context) ([]*AssetInformation, error) {
	data, _, err := c.getRequest(ctx, VAULT_ASSETS)
	if err != nil {
		return nil, fmt.Errorf("error during list balance by asset: %w", err)
	}

	tmp := make([]*AssetInformation, 0)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}

func (c *client) RetrieveVaultAccount(ctx context.Context, id string) (*VaultAccount, error) {
	uri := fmt.Sprintf("%v/%v", VAULT_ACCOUNTS, id)
	data, _, err := c.getRequest(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("error during retrieve vault account: %w", err)
	}

	tmp := new(VaultAccount)
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return tmp, nil
}
