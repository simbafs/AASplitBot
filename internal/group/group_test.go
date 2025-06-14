package group_test

import (
	"strconv"
	"testing"

	"splitbill/internal/bill"
	"splitbill/internal/group"
)

func username(id int64) (string, error) {
	return strconv.Itoa(int(id)), nil
}

func TestResultMsg(t *testing.T) {
	records := []bill.Record{
		{
			User:   1,
			Amount: 100,
			Shared: []int64{1, 2, 3},
		},
		{
			User:   2,
			Amount: 10,
			Shared: []int64{2, 3},
		},
		{
			User:   1,
			Amount: 7,
			Shared: []int64{1, 2, 3},
		},
		{
			User:   4,
			Amount: 100,
			Shared: []int64{4, 5},
		},
	}

	g := group.New(0, username)

	for _, r := range records {
		g.AddRecord(r.User, r.Shared, float64(r.Amount))
	}

	out, err := g.RecordsMsg()
	if err != nil {
		t.Fatalf("Failed to get records message: %v", err)
	}

	expected := "Records:\n$100(1)\n1, 2, 3\n$10(2)\n2, 3\n$7(1)\n1, 2, 3\n$100(4)\n4, 5\n"
	if out != expected {
		t.Errorf("Expected message:\n%s\nGot:\n%s", expected, out)
	}
}
