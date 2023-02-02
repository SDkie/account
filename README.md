# account
This is a client library in Go to access Form3 fake account API. It provides `Create`, `Update`, `Delete` methods.

---

* [Install](#install)
* [Setup](#setup)
* [Examples](#examples)
* [Testing](#testing)
* [Author](#author)

---

## Install
    go get -u github.com/SDkie/account

## Setup
Setup `ACCOUNTS_API_URL` environement variable to Form3 accounts API server

## Examples
#### Create client of account lib

    client, err := account.New()
	
#### Create Account

	country := "GB"
	data := account.AccountData{
		ID:             uuid.New().String(),
		OrganisationID: uuid.New().String(),
		Type:           "accounts",
	}
	data.Attributes = &account.AccountAttributes{
		BankID:       "400300",
		BankIDCode:   "GBDSC",
		Bic:          "NWBKGB22",
		Country:      &country,
		BaseCurrency: "GBP",
		Name:         []string{"John Doe"},
	}
	
	resp, err := client.Create(account)


#### Fetch Account
	resp, err := client.Fetch(<id>)


#### Delete Account
	resp, err := client.Delete(<id>)


## Testing
    docker-compose up
   
   
## Author
Kumar Sukhani

kumarsukhani@gmail.com
