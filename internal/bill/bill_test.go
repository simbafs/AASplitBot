package bill_test

import (
	"testing"

	"splitbill/internal/bill"
)

func TestBill(t *testing.T) {
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

	out := []bill.Transcation{
		{
			From:   5,
			To:     4,
			Amount: 5000,
		},
		{
			From:   3,
			To:     1,
			Amount: 4066,
		},
		{
			From:   2,
			To:     1,
			Amount: 3066,
		},
	}
	transcations, _, _ := bill.Split(records)
	if transcations == nil {
		t.Fatal("Expected transcations to not be nil")
	}
	for i, tr := range transcations {
		if tr.From != out[i].From || tr.To != out[i].To || tr.Amount != out[i].Amount {
			t.Errorf("Expected transcation %d to be %v, got %v", i, out[i], tr)
		}
	}
}
