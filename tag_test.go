package mapper

import (
	"errors"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		upperBound int
		wantStart  int
		wantEnd    int
		wantErr    error
	}{
		{
			name:       "valid tag with positive range",
			tag:        "2,5",
			upperBound: 10,
			wantStart:  2,
			wantEnd:    5,
			wantErr:    nil,
		},
		{
			name:       "invalid upper bound",
			tag:        "2,5",
			upperBound: 0,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInvalidUpperBound,
		},
		{
			name:       "valid tag with end capped by upperBound",
			tag:        "1,20",
			upperBound: 15,
			wantStart:  1,
			wantEnd:    15,
			wantErr:    nil,
		},
		{
			name:       "invalid empty tag",
			tag:        "",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagEmpty,
		},
		{
			name:       "invalid tag with non-numeric start",
			tag:        "a,5",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInvalidRangeValues,
		},
		{
			name:       "invalid tag with non-numeric end",
			tag:        "2,b",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInvalidRangeValues,
		},
		{
			name:       "ineffectual range: start equals end",
			tag:        "3,3",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInefectualRange,
		},
		{
			name:       "ineffectual range: start greater than end",
			tag:        "5,2",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    0,
			wantErr:    ErrTagInefectualRange,
		},
		{
			name:       "valid tag with -1 end (upperBound is respected)",
			tag:        "0,-1",
			upperBound: 10,
			wantStart:  0,
			wantEnd:    10,
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd, err := parseTag(tt.tag, tt.upperBound)

			// Check for expected errors
			if err != nil && tt.wantErr == nil {
				t.Errorf("expected no error, got %v", err)
			}
			if err == nil && tt.wantErr != nil {
				t.Errorf("expected error %v, got none", tt.wantErr)
			}
			if err != nil && tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}

			// Check the start and end values
			if gotStart != tt.wantStart {
				t.Errorf("expected start %d, got %d", tt.wantStart, gotStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("expected end %d, got %d", tt.wantEnd, gotEnd)
			}
		})
	}
}
