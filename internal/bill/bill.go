package bill

import (
	"fmt"
	"slices"
)

type (
	Money float64
	Milli int64
)

const ratio = 100

func (m Money) Milli() Milli {
	return Milli(m * ratio)
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f", m)
}

func (m Milli) Money() Money {
	return Money(m) / ratio
}

func (m Milli) String() string {
	return fmt.Sprintf("%d", m)
}

type Record struct {
	User   int64
	Amount Money
	Shared []int64
}

type Transcation struct {
	From   int64
	To     int64
	Amount Milli
}

type Person struct {
	ID     int64
	Amount Milli
}

func Split(records []Record) (result []Transcation, creditors, debtors []Person) {
	balances := make(map[int64]float64)

	for _, r := range records {
		balances[r.User] += float64(r.Amount.Milli())

		n := Milli(len(r.Shared))
		amount := r.Amount.Milli()

		perPerson := amount / n
		remain := amount - perPerson*n

		for _, user := range r.Shared {
			balances[user] -= float64(perPerson + remain)
			// 第一個人負責負擔取整造成的剩餘
			remain = 0
		}
	}

	for id, balance := range balances {
		if balance > 0 {
			creditors = append(creditors, Person{
				ID:     id,
				Amount: Milli(balance),
			})
		} else if balance < 0 {
			debtors = append(debtors, Person{
				ID:     id,
				Amount: Milli(-balance),
			})
		}
	}

	slices.SortFunc(creditors, func(a, b Person) int {
		return int(a.Amount - b.Amount)
	})
	slices.SortFunc(debtors, func(a, b Person) int {
		return int(b.Amount - a.Amount)
	})

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
