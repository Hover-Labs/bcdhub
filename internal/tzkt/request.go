package tzkt

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// URLs
const (
	TzKTURLV1    = "https://api.tzkt.io/v1/"
	TzKTServices = "https://services.tzkt.io"
)

// TzKT -
type TzKT struct {
	Host   string
	client http.Client

	retryCount int
}

// NewTzKT -
func NewTzKT(host string, timeout time.Duration) *TzKT {
	return &TzKT{
		Host: host,
		client: http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},

		retryCount: 3,
	}
}

func (t *TzKT) request(method, endpoint string, params map[string]string, response interface{}) error {
	uri := fmt.Sprintf("%s%s", t.Host, endpoint)

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest")
	}
	q := req.URL.Query()
	for key, value := range params {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	var resp *http.Response
	count := 0
	for ; count < t.retryCount; count++ {
		if resp, err = t.client.Do(req); err != nil {
			log.Printf("Attempt #%d: %s", count+1, err.Error())
			continue
		}
		break
	}

	if count == t.retryCount {
		return errors.New("Max HTTP request retry exceeded")
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(response); err != nil {
		return errors.Wrap(err, "json.Decode")
	}
	return nil
}

// GetHead - return head
func (t *TzKT) GetHead() (resp Head, err error) {
	err = t.request("GET", "head", nil, &resp)
	return
}

// GetOriginations - return originations
func (t *TzKT) GetOriginations(page, limit int64) (resp []Origination, err error) {
	if limit == 0 {
		limit = 1000
	}
	params := map[string]string{
		"p": fmt.Sprintf("%d", page),
		"n": fmt.Sprintf("%d", limit),
	}
	err = t.request("GET", "operations/originations", params, &resp)
	return
}

// GetOriginationsCount - return originations count
func (t *TzKT) GetOriginationsCount() (resp int64, err error) {
	err = t.request("GET", "operations/originations/count", nil, &resp)
	return
}

// GetSystemOperations - return system operations
func (t *TzKT) GetSystemOperations(page, limit int64) (resp []SystemOperation, err error) {
	if limit == 0 {
		limit = 1000
	}
	params := map[string]string{
		"p": fmt.Sprintf("%d", page),
		"n": fmt.Sprintf("%d", limit),
	}
	err = t.request("GET", "operations/system", params, &resp)
	return
}

// GetSystemOperationsCount - return system operations count
func (t *TzKT) GetSystemOperationsCount() (resp int64, err error) {
	err = t.request("GET", "operations/system/count", nil, &resp)
	return
}
