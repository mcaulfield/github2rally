// Rally HTTP Web Service Client

package rally

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	baseUrl string = "https://rally1.rallydev.com/slm/webservice/v2.0"
)

type queryResponse struct {
	Result QueryResult `json:"QueryResult,omitempty"`
}

type QueryResult struct {
	Count   float64 `json:"TotalResultCount,omitempty"`
	Results []Ref   `json:"Results,omitempty"`
}

type defectEnvelope struct {
	Defect Defect `json:"Defect,omitempty"`
}

type Defect struct {
	Name          string `json:"Name,omitempty"`
	Owner         *Ref   `json:"Owner,omitempty"`
	ScheduleState string `json:"ScheduleState,omitempty"`
	State         string `json:"State,omitempty"`
	SubmittedBy   *Ref   `json:"SubmittedBy,omitempty"`
	Description   string `json:"Description,omitempty"`
}

type Ref struct {
	Ref  string `json:"_ref,omitempty"`
	Type string `json:"_type,omitempty"`
}

type Client struct {
	httpClient *http.Client
	apiKey     string
}

// NewClient creates a new *rally.Client.
func NewClient(apiKey string) *Client {
	c := new(Client)
	c.httpClient = &http.Client{}
	c.apiKey = apiKey
	return c
}

// QueryDefect sends a query request for Defect objects to Rally.
func (c *Client) QueryDefect(q string) (*QueryResult, error) {
	return c.query("defect", q)
}

// QueryUser sends a query request for User objects to Rally.
func (c *Client) QueryUser(q string) (*QueryResult, error) {
	return c.query("user", q)
}

func (c *Client) query(typ string, q string) (*QueryResult, error) {
	queryUrl := fmt.Sprintf("%v/%v?query=%v", baseUrl, typ, url.QueryEscape(q))
	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("zsessionid", c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	qr := queryResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qr)
	if err != nil {
		return nil, err
	}
	return &qr.Result, nil
}

// CreateDefect sends a POST request to create a new Defect in Rally.
func (c *Client) CreateDefect(d *Defect) error {
	de := defectEnvelope{Defect: *d}
	buf, err := json.Marshal(de)
	if err != nil {
		return err
	}
	postUrl := fmt.Sprintf("%v/defect/create", baseUrl)
	req, err := http.NewRequest("POST", postUrl, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Add("zsessionid", c.apiKey)
	req.Header.Add("Content-Type", "application/json")
	_, err = c.httpClient.Do(req)
	return err
}
