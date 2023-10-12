package validator

import (
	"strconv"
	"testing"
	"time"
)

func TestIsValidToken(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "valid token",
			id:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			want: true,
		},
		{
			name: "too short",
			id:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083",
			want: false,
		},
		{
			name: "too long",
			id:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc222222222222222222222222222222222222",
			want: false,
		},
		{
			name: "invalid character",
			id:   "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2$",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidToken(tt.id); got != tt.want {
				t.Errorf("IsValidToken: for %v  = %v, but want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestIsValidBlock(t *testing.T) {
	tests := []struct {
		name  string
		block string
		want  bool
	}{
		{
			name:  "valid block",
			block: "18319881",
			want:  true,
		},
		{
			name:  "non-numeric block",
			block: "abcdef",
			want:  false,
		},
		{
			name:  "negative block number",
			block: "-18319881",
			want:  false,
		},
		{
			name:  "zero block",
			block: "0",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidBlock(tt.block); got != tt.want {
				t.Errorf("IsValidBlock: for %v = %v, but want %v", tt.block, got, tt.want)
			}
		})
	}
}

func TestIsValidFirst(t *testing.T) {
	tests := []struct {
		name  string
		first int
		want  bool
	}{
		{
			name:  "valid first within range",
			first: 5,
			want:  true,
		},
		{
			name:  "valid first at minimum boundary",
			first: 1,
			want:  true,
		},
		{
			name:  "valid first at maximum boundary",
			first: 1000,
			want:  true,
		},
		{
			name:  "invalid first below minimum boundary",
			first: 0,
			want:  false,
		},
		{
			name:  "invalid first above maximum boundary",
			first: 1001,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidFirst(tt.first); got != tt.want {
				t.Errorf("IsValidFirst: for %v = %v, but want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsValidRange(t *testing.T) {
	now := time.Now().Unix()
	validPastTimestamp := strconv.FormatInt(now-100, 10)
	validFutureTimestamp := strconv.FormatInt(now+100, 10)
	invalidTimestamp := "invalidTimestamp"

	tests := []struct {
		name   string
		from   string
		to     string
		expect bool
	}{
		{
			name:   "valid range in the past",
			from:   validPastTimestamp,
			to:     strconv.FormatInt(now-50, 10),
			expect: true,
		},
		{
			name:   "invalid range with from equals to",
			from:   validPastTimestamp,
			to:     validPastTimestamp,
			expect: false,
		},
		{
			name:   "invalid range with from greater than to",
			from:   strconv.FormatInt(now-50, 10),
			to:     validPastTimestamp,
			expect: false,
		},
		{
			name:   "invalid range with to in the future",
			from:   validPastTimestamp,
			to:     validFutureTimestamp,
			expect: false,
		},
		{
			name:   "invalid from timestamp",
			from:   invalidTimestamp,
			to:     validPastTimestamp,
			expect: false,
		},
		{
			name:   "invalid to timestamp",
			from:   validPastTimestamp,
			to:     invalidTimestamp,
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRange(tt.from, tt.to); got != tt.expect {
				t.Errorf("IsValidRange: for %v and %v = %v, but want %v", tt.from, tt.to, got, tt.expect)
			}
		})
	}
}
