package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type KeyInsightsTestSuite struct {
	*ClientTestSuite
}

func TestKeyInsights(t *testing.T) {
	suite.Run(t, &KeyInsightsTestSuite{
		NewClientTestSuite(),
	})
}

func (s *KeyInsightsTestSuite) TestGetKeyInsights() {
	client := newTestClient(s.Config)

	ctx := context.Background()

	resp, err := client.GetKeyInsights(ctx, "TUPRS", RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	
	s.Require().Equal("TUPRS", resp.Symbol)
	s.Require().NotEmpty(resp.Insight)
}