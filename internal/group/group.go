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

func (g *Group) RecordsMsg() (string, error) {
	if len(g.Bills) == 0 {
		return "沒有紀錄", nil
	}

	if len(g.Bills) == 0 {
		return "目前沒有任何分帳紀錄", nil
	}

	var records []string
	for _, r := range g.Bills {
		usernames := make([]string, 0, len(r.Shared))
		for _, id := range r.Shared {
			usernames = append(usernames, "@"+g.Username(id))
		}

		records = append(records, fmt.Sprintf("@%s 代墊了 %d 元，%s 要付錢", g.Username(r.User), r.Amount, strings.Join(usernames, "、")))
	}

	return fmt.Sprintf("目前的分帳紀錄有：\n%s", strings.Join(records, "\n")), nil
}

func (g *Group) ResultMsg() (string, error) {
	transcations, _, _ := Split(g.Bills)

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
		fmt.Fprintf(&msg, "@%s 要給 @%s $%d 元\n", fromName, toName, tr.Amount)
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
