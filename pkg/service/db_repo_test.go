// +build !integration

package service

import "testing"

func TestAmountScan(t *testing.T) {
	tests := map[string]struct {
		input      string
		shouldFail bool
		want       amount
	}{
		"basic": {
			input:      "(5.00,USD)",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"no_parens": {
			input:      "5.00,USD",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"no_Lparen": {
			input:      "5.00,USD)",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"no_Rparen": {
			input:      "(5.00,USD",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"quotes": {
			input:      "(\"5.00\",\"USD\")",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"quotes_full": {
			input:      "(\"5.00,USD\")",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"single_quotes": {
			input:      "('5.00','USD')",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"single_quotes_full": {
			input:      "('5.00,USD')",
			shouldFail: false,
			want:       amount{"5.00", "USD"},
		},
		"not_enough_elements": {
			input:      "(5.00)",
			shouldFail: true,
			want:       amount{},
		},
		"too_much_elements": {
			input:      "(5.00,USD,GBP)",
			shouldFail: true,
			want:       amount{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			a := &amount{}
			err := a.Scan(tc.input)
			if !tc.shouldFail && err != nil {
				t.Fatal(err)
			}
			if *a != tc.want {
				t.Fatalf("got: %v, want: %v", a, tc.want)
			}
		})
	}
}
