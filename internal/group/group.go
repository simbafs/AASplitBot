package group

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

type Group struct {
	mutex sync.Mutex
	Bills []Record
	Users map[int64]string
}

func New() *Group {
	return &Group{
		Bills: []Record{},
		Users: make(map[int64]string),
	}
}

func (g *Group) AddRecord(from int64, shared []int64, amount int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	r := Record{
		User:   from,
		Shared: shared,
		Amount: amount,
	}

	g.Bills = append(g.Bills, r)
}

func (g *Group) Result() (result []Transcation, creditors, debtors []Person) {
	result, creditors, debtors = Split(g.Bills)

	return result, creditors, debtors
}

func (g *Group) RecordsMsg() (string, error) {
	if len(g.Bills) == 0 {
		return "No records found.", nil
	}

	msg := strings.Builder{}
	for _, r := range g.Bills {
		name, ok := g.Users[r.User]
		if !ok {
			return "", fmt.Errorf("finding user %d", r.User)
		}
		fmt.Fprintf(&msg, "$%d(%s)\n", r.Amount, name)
		first := true
		for _, s := range r.Shared {
			name, ok = g.Users[s]
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
		fromName, ok := g.Users[tr.From]
		if !ok {
			return "", fmt.Errorf("finding user %d", tr.From)
		}
		toName, ok := g.Users[tr.To]
		if !ok {
			return "", fmt.Errorf("finding user %d", tr.To)
		}
		fmt.Fprintf(&msg, "%s -> %s $%d\n", fromName, toName, tr.Amount)
	}

	return msg.String(), nil
}

func (g *Group) Username(id int64) string {
	return g.Users[id]
}

// AddUser add new user, return false if the user already joint.
func (g *Group) AddUser(id int64, name string) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.Users == nil {
		g.Users = make(map[int64]string)
	}

	g.Users[id] = name
}

func (g *Group) Usernames() []string {
	usernames := make([]string, 0, len(g.Users))

	for _, name := range g.Users {
		usernames = append(usernames, name)
	}

	slices.Sort(usernames)

	return usernames
}

func (g *Group) IDs() []int64 {
	ids := make([]int64, 0, len(g.Users))

	for id := range g.Users {
		ids = append(ids, id)
	}

	slices.Sort(ids)

	return ids
}
