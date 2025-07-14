package driver

import (
	"strconv"
	"strings"
)

type Query func(query *map[string]string)

// QueryLimit set query limit
func QueryLimit(limit int) Query {
	return func(query *map[string]string) {
		(*query)["limit"] = strconv.FormatInt(int64(limit), 10)
	}
}

// QueryOffset set query offset
func QueryOffset(offset int) Query {
	return func(query *map[string]string) {
		(*query)["offset"] = strconv.FormatInt(int64(offset), 10)
	}
}

// GetShareSnap get share snap info
func (c *Pan115Client) GetShareSnap(shareCode, receiveCode, dirID string, Queries ...Query) (*ShareSnapResp, error) {
	if isCalledByAlistV3() {
		return nil, ErrorNotSupportAlist
	}
	result := ShareSnapResp{}
	query := map[string]string{
		"share_code":   shareCode,
		"receive_code": receiveCode,
		"cid":          dirID,
		"limit":        "20",
		"offset":       "0",
	}

	for _, q := range Queries {
		q(&query)
	}

	req := c.NewRequest().
		SetQueryParams(query).
		ForceContentType("application/json;charset=UTF-8").
		SetResult(&result)
	resp, err := req.Get(ApiShareSnap)
	if err := CheckErr(err, &result, resp); err != nil {
		return nil, err
	}

	return &result, nil
}

// ReceiveShare received files or folders from the shared link
func (c *Pan115Client) ReceiveShare(shareCode, receiveCode, dirID string, fileIDs ...string) error {
	if len(fileIDs) == 0 {
		return nil
	}

	form := map[string]string{
		"share_code":   shareCode,
		"receive_code": receiveCode,
		"file_id":      strings.Join(fileIDs, ","),
		"cid":          dirID,
		"user_id":      strconv.FormatInt(c.UserID, 10),
	}

	result := BasicResp{}
	req := c.NewRequest().
		SetFormData(form).
		ForceContentType("application/json;charset=UTF-8").
		SetResult(&result)
	resp, err := req.Post(ApiShareReceive)
	return CheckErr(err, &result, resp)
}
