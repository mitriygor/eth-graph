package calc

import "testing"

func TestSumNumbers(t *testing.T) {
	tests := []struct {
		name      string
		numbers   []string
		want      string
		wantError bool
	}{
		{
			name:      "valid numbers",
			numbers:   []string{"123456789.123456789", "987654321.987654321"},
			want:      "1111111111.1111111100",
			wantError: false,
		},
		{
			name:      "valid numbers with many decimal places",
			numbers:   []string{"123456789.123456789123456789", "987654321.987654321987654321"},
			want:      "1111111111.1111111111",
			wantError: false,
		},
		{
			name:      "one invalid number",
			numbers:   []string{"123456789.123456789", "invalid"},
			want:      "",
			wantError: true,
		},
		{
			name:      "all valid numbers",
			numbers:   []string{"1.1", "2.2", "3.3"},
			want:      "6.6000000000",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SumNumbers(tt.numbers)

			if (err != nil) != tt.wantError {
				t.Errorf("SumNumbers: for %v error = %v, wantError %v", tt.numbers, err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("SumNumbers: for %v = %v, want %v", tt.numbers, got, tt.want)
			}
		})
	}
}
