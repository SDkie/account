package account_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
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

		// Run the tests, API is healthy
		os.Exit(m.Run())
	}

	log.Printf("Exiting as Health API is still not up after 10 retries")
	os.Exit(1)
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

func TestFetchAfterCreate(t *testing.T) {
	t.Log("Given the need to test the account's Fetch() API")
	t.Logf("\tWhen creating client of accounts lib")
	client, err := account.New()
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error: %s", failed, err)
	}

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
	t.Logf("\tWhen creating client of accounts lib")
	client, err := account.New()
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error: %s", failed, err)
	}

	account := createTestAccount(t, client)

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", account.ID)
	err = client.Delete(account.ID, 0)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}
}

func TestListAfterDelete(t *testing.T) {
	t.Log("Given the need to test the account's List() API")
	t.Logf("\tWhen creating client of accounts lib")
	client, err := account.New()
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error: %s", failed, err)
	}

	account := createTestAccount(t, client)

	t.Logf("\tWhen checking the response of Delete() API for AccountID: %s", account.ID)
	err = client.Delete(account.ID, 0)
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}

	t.Logf("\tWhen checking the response of List() API")
	accounts, err := client.List("last")
	if err != nil {
		t.Fatalf("\t%s\tShould not respond with error:%s", failed, err)
	}

	for _, acc := range accounts.Accounts {
		if acc.ID == account.ID {
			t.Fatalf("\t%s\tShould not match test account.ID with one of the element in List()", failed)
		}
	}
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
		t.Logf("error deleting test account: %s", err)
	}
}
