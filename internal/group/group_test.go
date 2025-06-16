package group

import (
	"testing"
)

func TestRecordMsg(t *testing.T) {
	records := []Record{
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

	g := New()

	for _, r := range records {
		g.AddRecord(r.User, r.Shared, r.Amount)
	}

	out, err := g.RecordsMsg()
	if err != nil {
		t.Fatalf("Failed to get records message: %v", err)
	}

	expected := `$100(1)
  1, 2, 3
$10(2)
  2, 3
$7(1)
  1, 2, 3
$100(4)
  4, 5
`
	if out != expected {
		t.Errorf("Expected message:\n%s\nGot:\n%s", expected, out)
	}
}

func TestResultMsg(t *testing.T) {
	records := []Record{
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

	g := New()

	for _, r := range records {
		g.AddRecord(r.User, r.Shared, r.Amount)
	}

	out, err := g.ResultMsg()
	if err != nil {
		t.Fatalf("Failed to get records message: %v", err)
	}

	expected := `5 -> 4 $50
3 -> 1 $40
2 -> 1 $30
`
	if out != expected {
		t.Errorf("Expected message:\n%s\nGot:\n%s", expected, out)
	}
}
