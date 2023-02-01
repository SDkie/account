package account

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	CREATE_PATH = "%s/v1/organisation/accounts"
	FETCH_PATH  = "%s/v1/organisation/accounts/%s"
	DELETE_PATH = "%s/v1/organisation/accounts/%s?version=%d"
	LIST_PATH   = "%s/v1/organisation/accounts"
)

type Client struct {
	serverURL string
}

// New returns *Client using ACCOUNTS_API_URL env variable
func New() (*Client, error) {
	rURL := os.Getenv("ACCOUNTS_API_URL")
	if rURL == "" {
		err := fmt.Errorf("env ACCOUNTS_API_URL is required")
		log.Println(err)
		return nil, err
	}

	_, err := url.ParseRequestURI(rURL)
	if err != nil {
		log.Printf("Parsing ACCOUNTS_API_URL failed with err: %s", err)
		return nil, err
	}

	return &Client{serverURL: rURL}, nil
}

// AccountData is the structure for an account
type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

// AccountAttributes is the structure for an account's attributes
type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

// AccountResponse is the return type for Create and Fetch API
type AccountResponse struct {
	AccountData AccountData `json:"data"`
	Links       struct {
		Self string `json:"self"`
	} `json:"links"`
}

// Create sends request for Creating Account
func (c *Client) Create(data AccountData) (*AccountResponse, error) {
	url := fmt.Sprintf(CREATE_PATH, c.serverURL)

	reader, err := encodeRequest(data)
	if err != nil {
		return nil, err
	}

	httpResp, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Printf("http.Post request failed with err: %s\n", err)
		return nil, err
	}

	var response AccountResponse
	err = decodeResponse(httpResp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Fetch sends request to get Account by id
func (c *Client) Fetch(id string) (*AccountResponse, error) {
	url := fmt.Sprintf(FETCH_PATH, c.serverURL, id)

	httpResp, err := http.Get(url)
	if err != nil {
		log.Printf("http.Get request failed with err: %s\n", err)
		return nil, err
	}

	var response AccountResponse
	err = decodeResponse(httpResp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete sends request to delete Account by id
func (c *Client) Delete(id string, version int) error {
	url := fmt.Sprintf(DELETE_PATH, c.serverURL, id, version)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Printf("http.NewRequest failed with err: %s\n", err)
		return err
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Do(req) failed with err: %s\n", err)
		return err
	}

	if httpResp.StatusCode != http.StatusNoContent {
		err := fmt.Errorf("request failed with httpResp Status: %s", httpResp.Status)
		log.Println(err)
		return err
	}

	return nil
}

// ListResponse is responded by List API
type ListResponse struct {
	Accounts []AccountData `json:"data"`
	Links    struct {
		First string `json:"first"`
		Last  string `json:"last"`
		Prev  string `json:"prev"`
		Self  string `json:"self"`
	} `json:"links"`
}

// List sends request to List accounts based on pageNo
// pageNo can be 0, 1, 2, 3... or first, last, prev, self
// pageNo can also be empty string
func (c *Client) List(pageNo string) (*ListResponse, error) {
	rURL := fmt.Sprintf(LIST_PATH, c.serverURL)
	if pageNo != "" {
		params := url.Values{}
		params.Add("page[number]", pageNo)
		params.Encode()

		rURL = fmt.Sprintf("%s?%s", rURL, params.Encode())
	}

	httpResp, err := http.Get(rURL)
	if err != nil {
		log.Printf("http.Get request failed with err: %s\n", err)
		return nil, err
	}

	var response ListResponse
	err = decodeResponse(httpResp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
