package storage

import (
	"fmt"
)

type Account struct {
	Id       string
	Username string
}

func FindAccount(username string) *Account {
	db := Connect()
	defer db.Close()

	row := db.QueryRow("SELECT id, username FROM accounts WHERE username = ($1)", username)
	account := new(Account)
	err := row.Scan(&account.Id, &account.Username)
	if err != nil {
		fmt.Printf("Account Err: %s \n", err)
	}

	return account
}
