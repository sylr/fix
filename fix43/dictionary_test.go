package fix43

import "testing"

func TestDictionary(t *testing.T) {
	d := Dictionary()
	if d.Messages["D"] == nil || d.Messages["V"] == nil {
		t.Fatalf("FIX 4.3 dictionary lacks NewOrderSingle or MarketDataRequest")
	}
	if d.Major != "4" || d.Minor != "3" {
		t.Errorf("version = %s.%s, want 4.3", d.Major, d.Minor)
	}
}
