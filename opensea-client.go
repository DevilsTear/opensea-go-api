package opensea

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	mainnetAPI = "https://api.opensea.io"
	rinkebyAPI = "https://rinkeby-api.opensea.io"
)

type Opensea struct {
	API        string
	APIKey     string
	httpClient *http.Client
}

// GetAssetsParams model info
// @Description Asset listing params
type GetAssetsParams struct {
	Owner                  Address   `json:"owner" bson:"owner"`
	TokenIDs               []int32   `json:"token_ids" bson:"token_ids"`
	Collection             string    `json:"collection" bson:"collection"`
	CollectionSlug         string    `json:"collection_slug" bson:"collection_slug"`
	CollectionEditor       string    `json:"collection_editor" bson:"collection_editor"`
	OrderDirection         string    `json:"order_direction" bson:"order_direction"`
	AssetContractAddress   Address   `json:"asset_contract_address" bson:"asset_contract_address"`
	AssetContractAddresses []Address `json:"asset_contract_addresses" bson:"asset_contract_addresses"`
	Cursor                 *int      `json:"cursor" bson:"cursor"`
	Limit                  *int      `json:"limit" bson:"limit"`
	IncludeOrders          *bool     `json:"include_orders" bson:"include_orders"`
}

type errorResponse struct {
	Success bool `json:"success" bson:"success"`
}

func (e errorResponse) Error() string {
	return "Not success"
}

func NewOpensea(apiKey string) (*Opensea, error) {
	o := &Opensea{
		API:        mainnetAPI,
		APIKey:     apiKey,
		httpClient: defaultHttpClient(),
	}
	return o, nil
}

func NewOpenseaRinkeby(apiKey string) (*Opensea, error) {
	o := &Opensea{
		API:        rinkebyAPI,
		APIKey:     apiKey,
		httpClient: defaultHttpClient(),
	}
	return o, nil
}

func (p GetAssetsParams) Encode() string {
	q := url.Values{}

	if p.Owner.String() != "" && p.Owner != NullAddress {
		q.Set("owner", p.Owner.String())
	}

	if p.TokenIDs != nil && len(p.TokenIDs) > 0 {
		for i := 0; i < len(p.TokenIDs); i++ {
			q.Add("token_ids", fmt.Sprintf("%d", p.TokenIDs[i]))
		}
	}

	if p.Collection != "" {
		q.Set("collection", p.Collection)
	}

	if p.CollectionSlug != "" {
		q.Set("collection_slug", p.Collection)
	}

	if p.CollectionEditor != "" {
		q.Set("collection_editor", p.Collection)
	}

	if p.OrderDirection != "" {
		q.Set("order_direction", p.OrderDirection)
	}

	if p.AssetContractAddress.String() != "" && p.AssetContractAddress != NullAddress {
		q.Set("asset_contract_address", p.AssetContractAddress.String())
	}

	q.Del("asset_contract_addresses")
	if p.AssetContractAddresses != nil && len(p.AssetContractAddresses) > 0 {
		for i := 0; i < len(p.AssetContractAddresses); i++ {
			if p.AssetContractAddresses[i].String() != "" && p.AssetContractAddresses[i] != NullAddress {
				q.Add("asset_contract_addresses", p.AssetContractAddresses[i].String())
			}
		}
	}

	if p.Cursor != nil && *p.Cursor > 0 {
		q.Set("cursor", fmt.Sprintf("%d", *p.Cursor))
	}

	if p.Limit != nil && *p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", *p.Limit))
	}

	if p.IncludeOrders != nil {
		if includeOrders, err := strconv.ParseBool(fmt.Sprintf("%v", *p.IncludeOrders)); err != nil {
			q.Set("include_orders", strconv.FormatBool(includeOrders))
		}
	}

	return q.Encode()
}

func (o Opensea) GetAssetsTest(params GetAssetsParams) (*AssetResponse, error) {
	ctx := context.TODO()
	// o.SetHttpClient(http.DefaultClient)
	return o.GetAssetsWithContextTest(ctx, params)
}

func (o Opensea) GetAssets(params GetAssetsParams) (*AssetResponse, error) {
	ctx := context.TODO()
	// o.SetHttpClient(http.DefaultClient)
	return o.GetAssetsWithContext(ctx, params)
}

func (o Opensea) GetAssetsWithContext(ctx context.Context, params GetAssetsParams) (*AssetResponse, error) {
	path := "/api/v1/assets/?" + params.Encode()
	body, err := o.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	ret := new(AssetResponse)
	fmt.Println(string(body))
	return ret, json.Unmarshal(body, ret)
}

func (o Opensea) GetAssetsWithContextTest(ctx context.Context, params GetAssetsParams) (*AssetResponse, error) {
	path := "/api/v1/assets/?" + params.Encode()
	body, err := o.GetPathTest(ctx, path)
	if err != nil {
		return nil, err
	}
	ret := new(AssetResponse)
	fmt.Println(string(body))
	return ret, json.Unmarshal(body, ret)
}

func (o Opensea) GetSingleAsset(assetContractAddress string, tokenID *big.Int) (*Asset, error) {
	ctx := context.TODO()
	return o.GetSingleAssetWithContext(ctx, assetContractAddress, tokenID)
}

func (o Opensea) GetSingleAssetWithContext(ctx context.Context, assetContractAddress string, tokenID *big.Int) (*Asset, error) {
	path := fmt.Sprintf("/api/v1/asset/%s/%s", assetContractAddress, tokenID.String())
	body, err := o.GetPath(ctx, path)
	if err != nil {
		return nil, err
	}
	ret := new(Asset)
	return ret, json.Unmarshal(body, ret)
}

func (o Opensea) GetPathTest(ctx context.Context, path string) ([]byte, error) {
	return o.getURLTest(ctx, o.API+path)
}

func (o Opensea) GetPath(ctx context.Context, path string) ([]byte, error) {
	return o.getURL(ctx, o.API+path)
}

func (o Opensea) getURLTest(ctx context.Context, url string) ([]byte, error) {
	// url = "https://api.opensea.io/api/v1/assets?order_direction=desc&limit=20"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-API-KEY", "ad9d3916aa3a409f92a3bbd6aff78e8d")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println("\"n")
	fmt.Println("**********************TEST BODY START***************")
	fmt.Println("\"n")

	// fmt.Println(res)
	fmt.Println(string(body))

	fmt.Println("\"n")
	fmt.Println("**********************TEST BODY END***************")
	fmt.Println("\"n")

	return body, nil
}

func (o Opensea) getURL(ctx context.Context, url string) ([]byte, error) {
	// url = "https://api.opensea.io/api/v1/assets?order_direction=desc&limit=20"
	client := o.httpClient
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-API-KEY", o.APIKey)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		e := new(errorResponse)
		err = json.Unmarshal(body, e)
		if err != nil {
			return nil, err
		}
		if !e.Success {
			return nil, e
		}

		return nil, fmt.Errorf("backend returns status %d msg: %s", resp.StatusCode, string(body))
	}

	fmt.Println("\"n")
	fmt.Println("**********************PROD BODY START***************")
	fmt.Println("\"n")

	fmt.Println(string(body))

	fmt.Println("\"n")
	fmt.Println("**********************PROD BODY END***************")
	fmt.Println("\"n")
	return body, nil
}

// func (o Opensea) SetHttpClient(httpClient *http.Client) {
// 	o.httpClient = httpClient
// }

func defaultHttpClient() *http.Client {
	client := new(http.Client)
	var transport http.RoundTripper = &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		DisableKeepAlives:  false,
		DisableCompression: false,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 300 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	client.Transport = transport
	return client
}
