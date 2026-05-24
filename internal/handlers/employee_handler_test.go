package handlers

import "testing"

func TestParseHiredAt(t *testing.T) {
	tests := []struct {
		name        string
		input       *string
		wantErr     bool
		wantNilDate bool
	}{
		{
			name:        "valid date",
			input:       stringPtr("2026-05-24"),
			wantErr:     false,
			wantNilDate: false,
		},
		{
			name:        "nil date",
			input:       nil,
			wantErr:     false,
			wantNilDate: true,
		},
		{
			name:        "invalid date format",
			input:       stringPtr("20-25-2503"),
			wantErr:     true,
			wantNilDate: true,
		},
		{
			name:        "empty date",
			input:       stringPtr(""),
			wantErr:     true,
			wantNilDate: true,
		},
		{
			name:        "date with spaces",
			input:       stringPtr(" 2026-05-24 "),
			wantErr:     false,
			wantNilDate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, err := parseHiredAt(tt.input)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if tt.wantNilDate && gotDate != nil {
				t.Fatal("expected nil date")
			}

			if !tt.wantNilDate && gotDate == nil {
				t.Fatal("expected parsed date, got nil")
			}
		})
	}
}

func stringPtr(value string) *string {
	return &value
}
