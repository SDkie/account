package account_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/SDkie/account"
	"github.com/google/uuid"
)

const failed = "\u2717"

func TestMain(m *testing.M) {
	client, err := account.New()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// Waiting for the API server to be healthy
	sleepInterval := 10 * time.Millisecond
	for i := 0; i < 10; i++ {
		health, err := client.Health()
		if err != nil {
			log.Printf("Health checked failed with error: %s", err)
		} else if health.Status != "up" {
			err := fmt.Errorf("Health status is not up, Status:%s", health.Status)
			log.Println(err)
		}

		if err != nil {
			log.Printf("Sleeping:%v for Account API to be healthy", sleepInterval)
			time.Sleep(sleepInterval)
			sleepInterval *= 2
			continue
		}

		// API server is healthy, we are ready to run tests
		os.Exit(m.Run())
	}

	log.Printf("exiting as Health API is still not up after 10 retries")
	os.Exit(1)
}

func TestCreate(t *testing.T) {
	t.Log("Given the need to test the account's Create() API")

	client := getAccountsClient(t)
	account := createTestAccount(t, client)
	cleanupTestAccount(t, client, account.ID)
}

func TestCreateWithoutID(t *testing.T) {
	t.Log("Given the need to test the account's Create() API Without ID")

	client := getAccountsClient(t)
	account := generateRandomAccountData()
	account.ID = ""

	t.Logf("\tWhen checking the response of Create() API for AccountID: %s", account.ID)
	_, err := client.Create(account)
	if err == nil {
		t.Fatalf("\t%s\tShould respond with an error", failed)
	}

	if !strings.Contains(err.Error(), "validation failure list:\nvalidation failure list:\nid in body is required") {
		t.Fatalf("\t%s\tShould fail with /'id in body is required/' error msg", failed)
	}
}

func TestFetchAfterCreate(t *testing.T) {
	t.Log("Given the need to test the account's Fetch() API")

	client := getAccountsClient(t)
	account := createTestAccount(t, client)
	defer cleanupTestAccount(t, client, account.ID)

	t.Logf("\tWhen checking the response of Fetch() API for AccountID: %s", account.ID)
	resp, err := client.Fetch(account.ID)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
	if resp == nil {
		t.Fatalf("\t%s\tShould not respond with nil ResponseAccount", failed)
	}

	if !reflect.DeepEqual(resp.AccountData, account) {
		t.Fatalf("\t%s\tShould match AccountResponse.AccountData with given account", failed)
	}
}

func TestDeleteAfterCreate(t *testing.T) {
	t.Log("Given the need to test the account's Delete() API")

	client := getAccountsClient(t)
	account := createTestAccount(t, client)

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", account.ID)
	err := client.Delete(account.ID, 0)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
}

func TestFetchAfterDelete(t *testing.T) {
	t.Log("Given the need to test the account's Delete() API")

	client := getAccountsClient(t)
	account := createTestAccount(t, client)

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", account.ID)
	err := client.Delete(account.ID, 0)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}

	t.Logf("\tWhen checking the response of Fetch() API for AccountID: %s", account.ID)
	_, err = client.Fetch(account.ID)
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
		t.Logf("error in deleting test account: %s", err)
	}
}
