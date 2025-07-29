package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CapitalIncreaseTestSuite struct {
	*ClientTestSuite
}

func TestCapitalIncrease(t *testing.T) {
	suite.Run(t, &CapitalIncreaseTestSuite{
		NewClientTestSuite(),
	})
}

func (s *CapitalIncreaseTestSuite) TestGetAllCapitalIncreases() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetAllCapitalIncreases(ctx, 1, 10, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	s.Require().GreaterOrEqual(resp.RecordCount, 0)
	s.Require().NotEmpty(resp.Items)

	item := resp.Items[0]
	s.Require().NotZero(item.ID)
	s.Require().NotEmpty(item.Symbol)
	s.Require().NotEmpty(item.SpecifiedCurrency)
	s.Require().Equal("TRY", item.SpecifiedCurrency)
	s.Require().NotEmpty(item.CurrentCapital)
	s.Require().NotEmpty(item.TargetCapital)
	s.Require().NotNil(item.BoardDecisionDate)
	s.Require().NotEmpty(item.BoardDecisionDate)
	s.Require().NotEmpty(item.Types)
	s.Require().Contains([]string{"bonus", "rights", "external", "bonus_dividend"}, item.Types[0])
}

func (s *CapitalIncreaseTestSuite) TestGetCapitalIncreasesForInstrument() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetCapitalIncreasesForInstrument(ctx, "SASA", 1, 10, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)

	hasApprovedRecord := false
	for _, item := range resp.Items {
		s.Require().Equal("SASA.E", item.Symbol)
		s.Require().NotEmpty(item.CurrentCapital)
		s.Require().NotEmpty(item.TargetCapital)
		s.Require().Equal("TRY", item.SpecifiedCurrency)
		
		if item.SpkApplicationResult != nil && *item.SpkApplicationResult == "onay" {
			hasApprovedRecord = true
			s.Require().NotNil(item.SpkApplicationDate)
			s.Require().NotNil(item.SpkApprovalDate)
			s.Require().True(item.SpkApplicationDate.Before(*item.SpkApprovalDate) || 
				item.SpkApplicationDate.Equal(*item.SpkApprovalDate))
		}
	}
	
	s.Require().True(hasApprovedRecord)
}

func (s *CapitalIncreaseTestSuite) TestGetActiveRightsForInstrument() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetActiveRightsForInstrument(ctx, "SASA", "2024-07-20", RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
}