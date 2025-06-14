package main

import (
	"fmt"

	"splitbill/bill"
)

func main() {
	records := []bill.Record{
		{
			User:   1,
			Amount: 100,
			Shared: []bill.UserID{1, 2, 3},
		},
		{
			User:   2,
			Amount: 10,
			Shared: []bill.UserID{2, 3},
		},
		{
			User:   1,
			Amount: 7,
			Shared: []bill.UserID{1, 2, 3},
		},
	}

	transcations, creditors, debtors := bill.Split(records)
	for _, u := range creditors {
		fmt.Printf("債主: %d, 金額: %s(%d)\n", u.ID, u.Amount.Money(), u.Amount)
	}
	for _, u := range debtors {
		fmt.Printf("債務人: %d, 金額: %s(%d)\n", u.ID, u.Amount.Money(), u.Amount)
	}
	for _, t := range transcations {
		fmt.Printf("%d 要給 %d %s(%d)\n", t.From, t.To, t.Amount.Money(), t.Amount)
	}
}
