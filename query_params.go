package chainpulse

import (
	"net/url"
	"strconv"
	"strings"
)

// QueryParams 为 Query API 公共查询参数。
type QueryParams struct {
	ChainID        uint64
	Consistency    string // stable | pending | finalized
	IncludeRemoved bool
	Limit          uint64
}

func applyQueryParams(query url.Values, params ...QueryParams) url.Values {
	if query == nil {
		query = url.Values{}
	}
	if len(params) == 0 {
		return query
	}
	p := params[0]
	if p.ChainID > 0 {
		query.Set("chainId", strconv.FormatUint(p.ChainID, 10))
	}
	if c := strings.TrimSpace(p.Consistency); c != "" {
		query.Set("consistency", c)
	}
	if p.IncludeRemoved {
		query.Set("includeRemoved", "true")
	}
	if p.Limit > 0 {
		query.Set("limit", strconv.FormatUint(p.Limit, 10))
	}
	return query
}
