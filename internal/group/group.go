package group

import (
	"fmt"
	"strings"
	"sync"

	"splitbill/internal/bill"
)

type Group struct {
	mutex sync.Mutex
	bills []bill.Record
	users map[int64]string
}

func New() *Group {
	return &Group{
		bills: []bill.Record{},
		users: make(map[int64]string),
	}
}

func (g *Group) AddRecord(from int64, shared []int64, amount int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

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
	for _, r := range g.bills {
		name, ok := g.users[r.User]
		if !ok {
			return "", fmt.Errorf("finding user %d", r.User)
		}
		fmt.Fprintf(&msg, "$%d(%s)\n", r.Amount, name)
		first := true
		for _, s := range r.Shared {
			name, ok = g.users[s]
			if !ok {
				return "", fmt.Errorf("finding user %d", s)
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

func (g *Group) ResultMsg() (string, error) {
	transcations, _, _ := g.Result()

	msg := strings.Builder{}

	for _, tr := range transcations {
		fromName, ok := g.users[tr.From]
		if !ok {
			return "", fmt.Errorf("finding user %d", tr.From)
		}
		toName, ok := g.users[tr.To]
		if !ok {
			return "", fmt.Errorf("finding user %d", tr.To)
		}
		fmt.Fprintf(&msg, "%s -> %s $%d\n", fromName, toName, tr.Amount)
	}

	return msg.String(), nil
}
