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

	// 4.               50 =  50 -50 \
	// 1.   66     + 4     =  70       -40 -30 \
	// 5.              -50 = -50 +50 \
	// 3. - 33 - 5 - 2     = -40       +40 \
	// 2. - 33 + 5 - 2     = -30           +30 \

	out := []bill.Transcation{
		{
			From:   5,
			To:     4,
			Amount: 50,
		},
		{
			From:   3,
			To:     1,
			Amount: 40,
		},
		{
			From:   2,
			To:     1,
			Amount: 30,
		},
	}
	transcations, creditors, debtors := bill.Split(records)
	t.Log("Transcations:")
	for _, tr := range transcations {
		t.Log(tr)
	}
	t.Log("Creditors:")
	for _, c := range creditors {
		t.Log(c)
	}
	t.Log("Debtors:")
	for _, d := range debtors {
		t.Log(d)
	}

	if transcations == nil {
		t.Fatal("Expected transcations to not be nil")
	}
	for i, tr := range transcations {
		if tr.From != out[i].From || tr.To != out[i].To || tr.Amount != out[i].Amount {
			t.Errorf("Expected transcation %d to be %v, got %v", i, out[i], tr)
		}
	}
}
