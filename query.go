package chainpulse

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

type EventQuery struct {
	Address   string
	TxHash    string
	Topic0    string
	Topic1    string
	Topic2    string
	Topic3    string
	FromBlock uint64
	ToBlock   uint64
	Limit     uint64
}

func (c *Client) Health(ctx context.Context) (HealthResponse, error) {
	var out HealthResponse
	err := c.do(ctx, c.queryBaseURL, http.MethodGet, "/healthz", nil, nil, &out, false)
	return out, err
}

func (c *Client) GetTransaction(ctx context.Context, txHash string, params ...QueryParams) (ChainTransaction, error) {
	query := applyQueryParams(nil, params...)
	var out ChainTransaction
	err := c.doQuery(ctx, http.MethodGet, "/v1/transactions/"+pathEscape(txHash), query, nil, &out)
	return out, err
}

func (c *Client) ListTransactionInternalTransactions(ctx context.Context, txHash string, params ...QueryParams) ([]InternalTransaction, error) {
	query := applyQueryParams(nil, params...)
	var out listResponse[InternalTransaction]
	err := c.doQuery(ctx, http.MethodGet, "/v1/transactions/"+pathEscape(txHash)+"/internal-transactions", query, nil, &out)
	return out.Items, err
}

func (c *Client) ListAddressTransactions(ctx context.Context, address string, params ...QueryParams) ([]ChainTransaction, error) {
	query := applyQueryParams(nil, params...)
	var out listResponse[ChainTransaction]
	err := c.doQuery(ctx, http.MethodGet, "/v1/addresses/"+pathEscape(address)+"/transactions", query, nil, &out)
	return out.Items, err
}

func (c *Client) ListAddressInternalTransactions(ctx context.Context, address string, params ...QueryParams) ([]InternalTransaction, error) {
	query := applyQueryParams(nil, params...)
	var out listResponse[InternalTransaction]
	err := c.doQuery(ctx, http.MethodGet, "/v1/addresses/"+pathEscape(address)+"/internal-transactions", query, nil, &out)
	return out.Items, err
}

func (c *Client) ListAddressBalances(ctx context.Context, address string, params ...QueryParams) ([]AddressBalance, error) {
	query := applyQueryParams(nil, params...)
	var out listResponse[AddressBalance]
	err := c.doQuery(ctx, http.MethodGet, "/v1/addresses/"+pathEscape(address)+"/balances", query, nil, &out)
	return out.Items, err
}

func (c *Client) ListNFTBalances(ctx context.Context, address string, standard string, params ...QueryParams) ([]NFTTokenBalance, error) {
	query := applyQueryParams(nil, params...)
	if standard != "" {
		query.Set("standard", standard)
	}
	var out listResponse[NFTTokenBalance]
	err := c.doQuery(ctx, http.MethodGet, "/v1/addresses/"+pathEscape(address)+"/nft-balances", query, nil, &out)
	return out.Items, err
}

func (c *Client) ListEvents(ctx context.Context, options EventQuery, params ...QueryParams) ([]ChainLog, error) {
	query := eventValues(options)
	query = applyQueryParams(query, params...)
	var out listResponse[ChainLog]
	err := c.doQuery(ctx, http.MethodGet, "/v1/events", query, nil, &out)
	return out.Items, err
}

func (c *Client) ListContractEvents(ctx context.Context, contractAddress string, options EventQuery, params ...QueryParams) ([]ChainLog, error) {
	query := eventValues(options)
	query = applyQueryParams(query, params...)
	var out listResponse[ChainLog]
	err := c.doQuery(ctx, http.MethodGet, "/v1/contracts/"+pathEscape(contractAddress)+"/events", query, nil, &out)
	return out.Items, err
}

func eventValues(options EventQuery) url.Values {
	query := url.Values{}
	if options.Address != "" {
		query.Set("address", options.Address)
	}
	if options.TxHash != "" {
		query.Set("txHash", options.TxHash)
	}
	if options.Topic0 != "" {
		query.Set("topic0", options.Topic0)
	}
	if options.Topic1 != "" {
		query.Set("topic1", options.Topic1)
	}
	if options.Topic2 != "" {
		query.Set("topic2", options.Topic2)
	}
	if options.Topic3 != "" {
		query.Set("topic3", options.Topic3)
	}
	if options.FromBlock > 0 {
		query.Set("fromBlock", strconv.FormatUint(options.FromBlock, 10))
	}
	if options.ToBlock > 0 {
		query.Set("toBlock", strconv.FormatUint(options.ToBlock, 10))
	}
	if options.Limit > 0 {
		query.Set("limit", strconv.FormatUint(options.Limit, 10))
	}
	return query
}
