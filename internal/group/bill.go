package group

import (
	"slices"
)

type Record struct {
	User   int64
	Amount int
	Shared []int64
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
