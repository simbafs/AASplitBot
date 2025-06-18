package group

import (
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Record struct {
	User        int64
	Amount      int
	Shared      []int64
	Description string
}

// ParseRecord parses a string representation of a record and returns a Record struct.
// The string format is `$<Amount> (<User>) <SharedUser1> <SharedUser2> ... <Description>`, all field except Amount are optional.
// For example:
// - $100 (@user1) @user2 @user3
// - $50 @user1 @user2 @user3 some description
// - $200
func (g *Group) ParseRecord(str string) Record {
	var r Record

	// 移除前後空白
	str = strings.TrimSpace(str)

	// 用 regex 解析金額與可能存在的 (user)
	re := regexp.MustCompile(`^\$(\d*)\s*(?:\(\s*@(\w*)\s*\))?`)
	matches := re.FindStringSubmatch(str)
	if len(matches) == 0 {
		return r // 無法解析金額，回傳空的 Record
	}
	slog.Debug("Regex matches", "matches", fmt.Sprintf("%#v", matches))

	// 金額
	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		return r
	}
	r.Amount = amount

	// user (可選)
	if matches[2] != "" {
		id, ok := g.ID(matches[2])
		if ok {
			r.User = id
		} else {
			slog.Error("無法解析使用者 ID", "username", matches[2])
		}
	}

	// 剩下的字串
	str = strings.TrimSpace(str[len(matches[0]):])

	// 拆解 @user 或描述
	parts := strings.Fields(str)
	for i, part := range parts {
		if strings.HasPrefix(part, "@") {
			// 是共享使用者
			id, ok := g.ID(part[1:])
			if !ok {
				slog.Error("無法解析使用者 ID", "username", part)
				continue
			}
			r.Shared = append(r.Shared, id)
		} else {
			// 遇到第一個不是 @ 的，就當成是描述開頭
			r.Description = strings.Join(parts[i:], " ")
			break
		}
	}

	slog.Debug("Parsed record", "record", fmt.Sprintf("%#v", r))

	return r
}

type Transcation struct {
	From   int64
	To     int64
	Amount int
}

type Person struct {
	ID     int64
	Amount int
}

func Split(records []Record) (result []Transcation, creditors, debtors []Person) {
	balances := make(map[int64]int)

	for _, r := range records {
		n := len(r.Shared)
		amount := r.Amount

		balances[r.User] += amount

		perPerson := amount / n
		remain := amount - perPerson*n

		for _, user := range r.Shared {
			balances[user] -= perPerson + remain
			// 第一個人負責負擔取整造成的剩餘
			remain = 0
		}
	}

	for id, balance := range balances {
		if balance > 0 {
			creditors = append(creditors, Person{
				ID:     id,
				Amount: balance,
			})
		} else if balance < 0 {
			debtors = append(debtors, Person{
				ID:     id,
				Amount: -balance,
			})
		}
	}

	slices.SortFunc(creditors, func(a, b Person) int {
		return int(a.Amount - b.Amount)
	})
	slices.SortFunc(debtors, func(a, b Person) int {
		return int(b.Amount - a.Amount)
	})

	// TODO: use priority queue to optimize the matching process
	c := make([]Person, len(creditors))
	copy(c, creditors)
	d := make([]Person, len(debtors))
	copy(d, debtors)

	i, j := 0, 0
	for i < len(creditors) && j < len(debtors) {
		creditor := creditors[i]
		debtor := debtors[j]
		amount := min(creditor.Amount, debtor.Amount)

		result = append(result, Transcation{
			From:   debtor.ID,
			To:     creditor.ID,
			Amount: amount,
		})

		creditors[i].Amount -= amount
		debtors[j].Amount -= amount

		if creditors[i].Amount == 0 {
			i++
		}
		if debtors[j].Amount == 0 {
			j++
		}
	}

	return result, c, d
}
