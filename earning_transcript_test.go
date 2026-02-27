package laplace

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EarningTranscriptTestSuite struct {
	*ClientTestSuite
}

func TestEarningTranscript(t *testing.T) {
	suite.Run(t, &EarningTranscriptTestSuite{
		NewClientTestSuite(),
	})
}

func TestFlexibleTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantY   int // expected year, 0 means zero time
	}{
		{"RFC3339", `"2024-03-15T10:30:00Z"`, false, 2024},
		{"date only", `"2024-03-15"`, false, 2024},
		{"datetime no tz", `"2024-03-15T10:30:00"`, false, 2024},
		{"unix timestamp", `1710489000`, false, 2024},
		{"null", `null`, false, 0},
		{"empty string", `""`, false, 0},
		{"invalid", `"not-a-date"`, true, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ft FlexibleTime
			err := json.Unmarshal([]byte(tc.input), &ft)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tc.wantY == 0 {
				if !ft.Time.IsZero() {
					t.Fatalf("expected zero time, got %v", ft.Time)
				}
			} else {
				if ft.Time.Year() != tc.wantY {
					t.Fatalf("expected year %d, got %d", tc.wantY, ft.Time.Year())
				}
			}
		})
	}
}

func TestFlexibleTime_MarshalJSON(t *testing.T) {
	var ft FlexibleTime
	b, err := json.Marshal(ft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(b) != "null" {
		t.Fatalf("expected null for zero time, got %s", b)
	}
}

func (s *EarningTranscriptTestSuite) TestGetEarningsTranscriptList() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetEarningsTranscriptList(ctx, RegionUs, "AAPL")
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Greater(len(resp), 0)

	item := resp[0]
	s.Require().NotEmpty(item.Symbol)
	s.Require().Greater(item.Year, 0)
	s.Require().Greater(item.Quarter, 0)
}
