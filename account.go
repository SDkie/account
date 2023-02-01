package account

import (
	"fmt"
	"log"
	"net/http"
)

const (
	CREATE_PATH = "%s/v1/organisation/accounts"
)

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
func Create(data AccountData) (*AccountResponse, error) {
	url := fmt.Sprintf(CREATE_PATH, "http://localhost:8080")

	reader, err := encodeRequest(data)
	if err != nil {
		return nil, err
	}

	httpResp, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Printf("error from http.Post: %s\n", err)
		return nil, err
	}

	var response AccountResponse
	err = decodeResponse(httpResp, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
