package account_test

import (
	"reflect"
	"testing"

	"github.com/SDkie/account"
	"github.com/google/uuid"
)

const failed = "\u2717"

func generateRandomAccountData() account.AccountData {
	country := "GB"
	var version int64
	data := account.AccountData{
		ID:             uuid.New().String(),
		OrganisationID: uuid.New().String(),
		Type:           "accounts",
		Version:        &version,
	}
	data.Attributes = &account.AccountAttributes{
		BankID:       "400300",
		BankIDCode:   "GBDSC",
		Bic:          "NWBKGB22",
		Country:      &country,
		BaseCurrency: "GBP",
		Name:         []string{"John Doe"},
	}

	return data
}

func createTestAccount(t *testing.T, client *account.Client) account.AccountData {
	account := generateRandomAccountData()

	t.Logf("\tWhen checking the response of Create() API for AccountID: %s", account.ID)
	resp, err := client.Create(account)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
	if resp == nil {
		t.Fatalf("\t%s\tShould not respond with nil ResponseAccount", failed)
	}

	if !reflect.DeepEqual(resp.AccountData, account) {
		t.Fatalf("\t%s\tShould match AccountResponse.AccountData with given account", failed)
	}

	return account
}

func cleanupTestAccount(t *testing.T, client *account.Client, id string) {
	err := client.Delete(id, 0)
	if err != nil {
		t.Logf("error deleting test account: %s", err)
	}
}

func TestCreate(t *testing.T) {
	t.Log("Given the need to test the account's Create() API")
	t.Logf("\tWhen creating client of accounts lib")
	client, err := account.New()
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error: %s", failed, err)
	}

	account := createTestAccount(t, client)
	cleanupTestAccount(t, client, account.ID)
}
