package fix44

import "testing"

func TestDictionary(t *testing.T) {
	d := Dictionary()
	if d.Messages["D"] == nil || d.Messages["V"] == nil {
		t.Fatalf("FIX 4.4 dictionary lacks NewOrderSingle or MarketDataRequest")
	}
	if d.Major != "4" || d.Minor != "4" {
		t.Errorf("version = %s.%s, want 4.4", d.Major, d.Minor)
	}
}
