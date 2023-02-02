package account_test

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SDkie/account"
	"github.com/google/uuid"
)

const failed = "\u2717"

func TestCreate(t *testing.T) {
	t.Log("Given the need to test the account's Create() API")

	client := getAccountsClient(t)
	accountData := createTestAccount(t, client)
	cleanupTestAccount(t, client, accountData.ID)
}

func TestCreateWithCustomHttpClient(t *testing.T) {
	t.Log("Given the need to test the account's Create() API with custom httpClient")

	client := getAccountsClient(t)
	var httpClient http.Client
	httpClient.Timeout = 1 * time.Minute
	client.SetHTTPClient(httpClient)

	accountData := createTestAccount(t, client)
	cleanupTestAccount(t, client, accountData.ID)
}

func TestCreateWithoutID(t *testing.T) {
	t.Log("Given the need to test the account's Create() API Without ID")

	client := getAccountsClient(t)
	accountData := generateRandomAccountData()
	accountData.ID = ""

	t.Logf("\tWhen checking the response of Create() API for AccountID: %s", accountData.ID)
	_, err := client.Create(accountData)
	if err == nil {
		t.Fatalf("\t%s\tShould respond with an error", failed)
	}

	if !strings.Contains(err.Error(), "validation failure list:\nvalidation failure list:\nid in body is required") {
		t.Fatalf("\t%s\tShould fail with /'id in body is required/' error msg", failed)
	}
}

func TestCreateWithDuplicateConstraint(t *testing.T) {
	t.Log("Given the need to test the account's Create() API With Duplicate Constraint")

	client := getAccountsClient(t)
	accountData := createTestAccount(t, client)
	defer cleanupTestAccount(t, client, accountData.ID)

	t.Logf("\tWhen checking the response of Create() API for AccountID: %s", accountData.ID)
	_, err := client.Create(accountData)
	if err == nil {
		t.Fatalf("\t%s\tShould respond with an error", failed)
	}

	if err.Error() != account.ErrDuplicateConstraint {
		t.Fatalf("\t%s\tShould fail with '%s' error msg", failed, account.ErrDuplicateConstraint)
	}
}

func TestFetchAfterCreate(t *testing.T) {
	t.Log("Given the need to test the account's Fetch() API")

	client := getAccountsClient(t)
	accountData := createTestAccount(t, client)
	defer cleanupTestAccount(t, client, accountData.ID)

	t.Logf("\tWhen checking the response of Fetch() API for AccountID: %s", accountData.ID)
	resp, err := client.Fetch(accountData.ID)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
	if resp == nil {
		t.Fatalf("\t%s\tShould not respond with nil ResponseAccount", failed)
	}

	if !reflect.DeepEqual(resp.AccountData, accountData) {
		t.Fatalf("\t%s\tShould match AccountResponse.AccountData with given account", failed)
	}
}

func TestFetchWithInvalidID(t *testing.T) {
	t.Log("Given the need to test the account's Fetch() API with invalid ID")

	id := "0"

	t.Logf("\tWhen checking the response of Fetch() API for AccountID: %s", id)
	client := getAccountsClient(t)
	_, err := client.Fetch(id)
	if err == nil {
		t.Fatalf("\t%s\tShould fail with err", failed)
	}
	if err.Error() != account.ErrInvalidUUID {
		t.Fatalf("\t%s\tShould fail with '%s' error msg", failed, account.ErrInvalidUUID)
	}
}

func TestDeleteAfterCreate(t *testing.T) {
	t.Log("Given the need to test the account's Delete() API")

	client := getAccountsClient(t)
	accountData := createTestAccount(t, client)

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", accountData.ID)
	err := client.Delete(accountData.ID, 0)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
}

func TestDeleteWithInvalidID(t *testing.T) {
	t.Log("Given the need to test the account's Delete() API with invalid ID")

	id := "0"

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", id)
	client := getAccountsClient(t)
	err := client.Delete(id, 0)
	if err == nil {
		t.Fatalf("\t%s\tShould fail with err", failed)
	}
	if err.Error() != account.ErrInvalidUUID {
		t.Fatalf("\t%s\tShould fail with '%s' error msg", failed, account.ErrInvalidUUID)
	}
}

func TestFetchAfterDelete(t *testing.T) {
	t.Log("Given the need to test the account's Delete() API")

	client := getAccountsClient(t)
	accountData := createTestAccount(t, client)

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", accountData.ID)
	err := client.Delete(accountData.ID, 0)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}

	t.Logf("\tWhen checking the response of Fetch() API for AccountID: %s", accountData.ID)
	_, err = client.Fetch(accountData.ID)
	if err == nil {
		t.Fatalf("\t%s\tShould fail with err", failed)
	}
}

func getAccountsClient(t *testing.T) *account.Client {
	t.Logf("\tWhen creating client of account lib")
	client, err := account.New()
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error: %s", failed, err)
	}

	return client
}

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
	accountData := generateRandomAccountData()

	t.Logf("\tWhen checking the response of Create() API for AccountID: %s", accountData.ID)
	resp, err := client.Create(accountData)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
	if resp == nil {
		t.Fatalf("\t%s\tShould not respond with nil ResponseAccount", failed)
	}

	if !reflect.DeepEqual(resp.AccountData, accountData) {
		t.Fatalf("\t%s\tShould match AccountResponse.AccountData with given account", failed)
	}

	return accountData
}

func cleanupTestAccount(t *testing.T, client *account.Client, id string) {
	err := client.Delete(id, 0)
	if err != nil {
		t.Logf("error in deleting test account: %s", err)
	}
}
