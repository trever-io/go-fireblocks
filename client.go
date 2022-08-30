package fireblocks

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const BASE_URL = "https://api.fireblocks.io"

var (
	ErrInvalidPrivateKey = errors.New("private key is no valid private key (PEM)")
)

type Error struct {
	status int
	body   string
}

func (e *Error) Error() string {
	return fmt.Sprintf("fireblocks error: %d: %v", e.status, e.body)
}

type Client interface {
	ListInternalWallets(ctx context.Context) ([]*InternalWallet, error)
	CreateInternalWallet(ctx context.Context, req *CreateInternalWalletRequest) (*InternalWallet, error)
}

type client struct {
	key       string
	secretKey *rsa.PrivateKey

	nonceMtx sync.Mutex
}

func NewClient(key string, privateKey string) (Client, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, ErrInvalidPrivateKey
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing rsa key: %w", err)
	}

	secretKey, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidPrivateKey
	}

	return &client{
		key:       key,
		secretKey: secretKey,
	}, nil
}

func (c *client) signRequest(uri string, req *http.Request) error {
	var bodyBytes []byte

	if req.Body != nil {
		tmp, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("error reading request body: %w", err)
		}
		err = req.Body.Close()

		bodyBytes = make([]byte, len(tmp))
		copy(bodyBytes, tmp)

		if err != nil {
			return fmt.Errorf("error closing body reader: %w", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(tmp))
		req.Header.Add("content-type", "application/json")
	} else {
		bodyBytes = []byte("")
	}

	shaHash := sha256.New()
	_, err := shaHash.Write(bodyBytes)
	if err != nil {
		return fmt.Errorf("error writing sha256 hash: %w", err)
	}
	shaSum := shaHash.Sum(nil)
	shaSumStr := hex.EncodeToString(shaSum)

	c.nonceMtx.Lock()
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"uri":      uri,
		"nonce":    now.UnixNano(),
		"iat":      now.Unix(),
		"exp":      now.Add(55 * time.Second).Unix(),
		"sub":      c.key,
		"bodyHash": shaSumStr,
	})
	c.nonceMtx.Unlock()

	jwtString, err := token.SignedString(c.secretKey)
	if err != nil {
		return fmt.Errorf("error signing jwt token: %w", err)
	}

	req.Header.Add("X-API-key", c.key)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", jwtString))
	return nil
}

func (c *client) getRequest(ctx context.Context, uri string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%v%v", BASE_URL, uri), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating get request: %w", err)
	}

	err = c.signRequest(uri, req)
	if err != nil {
		return nil, fmt.Errorf("error signing request: %w", err)
	}

	return c.doRequest(req)
}

func (c *client) postRequest(ctx context.Context, uri string, body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%v%v", BASE_URL, uri), bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("error creating post request: %w", err)
	}

	err = c.signRequest(uri, req)
	if err != nil {
		return nil, fmt.Errorf("error signing request: %w", err)
	}

	return c.doRequest(req)
}

func (c *client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during request: %w", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		apiError := &Error{
			status: resp.StatusCode,
			body:   string(b),
		}

		return nil, apiError
	}

	return b, nil
}
