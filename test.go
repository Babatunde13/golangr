// Marshal a struct into JSON. Our objective is to print a JSON string from a slice of user.

//i.	Ignore an empty field you don't want to initialize e.g. Balance field
//ii. 	Don't display the user's bvn in the result

package main

import (
	"encoding/json"
	"fmt"
)

type user struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Bvn     string `json:"-"`
	Account account `json:"account"`
}

type account struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Balance float64 `json:"balance,omitempty"`
}

func (a *account) String() string {
	return fmt.Sprintf("Account: { id: %d, name: %s }", a.Id, a.Name)
}

func (u *user) String () string {
	return fmt.Sprintf("User: { id: %d, name: %s, bvn: %s, account: %v }", u.Id, u.Name, u.Bvn, u.Account )
}

func Test() {
	account1 := account{ Id: 1, Name: "Savings", Balance: 1000.00 }
	account2 := account{ Id: 2, Name: "Current", Balance: 2000.00 }
	account3 := account{ Id: 3, Name: "Savings", Balance: 3000.00 }
	users := []user{
		{ Id: 1, Name: "John Doe", Bvn: "1234567890", Account: account1 },
		{ Id: 2, Name: "Jane Doe", Bvn: "0987654321", Account: account2 },
		{ Id: 3, Name: "John Smith", Bvn: "1234567890", Account:account3 },
	}

	usersJSON, err := json.Marshal(users); if err != nil {
		panic(err)
	}

	println(string(usersJSON))
}