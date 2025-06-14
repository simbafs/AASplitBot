package group

import (
	"fmt"
	"strings"

	"splitbill/internal/bill"
)

type Username = func(int64) (string, error)

type Group struct {
	chatID   int64
	bills    []bill.Record
	username Username
}

func New(chat int64, username Username) *Group {
	return &Group{
		chatID:   chat,
		bills:    []bill.Record{},
		username: username,
	}
}

func (g *Group) AddRecord(from int64, shared []int64, amount int) {
	r := bill.Record{
		User:   from,
		Shared: shared,
		Amount: amount,
	}

	g.bills = append(g.bills, r)
}

func (g *Group) Result() (result []bill.Transcation, creditors, debtors []bill.Person) {
	result, creditors, debtors = bill.Split(g.bills)

	return result, creditors, debtors
}

func (g *Group) RecordsMsg() (string, error) {
	if len(g.bills) == 0 {
		return "No records found.", nil
	}

	msg := strings.Builder{}
	msg.WriteString("Records:\n")
	for _, r := range g.bills {
		name, err := g.username(r.User)
		if err != nil {
			return "", fmt.Errorf("finding user %d: %w", r.User, err)
		}
		fmt.Fprintf(&msg, "$%d(%s)\n", r.Amount, name)
		first := true
		for _, s := range r.Shared {
			name, err = g.username(s)
			if err != nil {
				return "", fmt.Errorf("finding user %d: %w", s, err)
			}
			if first {
				fmt.Fprintf(&msg, "  %s", name)
			} else {
				fmt.Fprintf(&msg, ", %s", name)
			}
			first = false
		}
		msg.WriteString("\n")
	}

	return msg.String(), nil
}
