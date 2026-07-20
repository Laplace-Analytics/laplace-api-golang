package laplace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HalkaArzTestSuite struct {
	*ClientTestSuite
}

func TestHalkaArz(t *testing.T) {
	suite.Run(t, &HalkaArzTestSuite{
		NewClientTestSuite(),
	})
}

func (s *HalkaArzTestSuite) TestGetAllHalkaArz() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	resp, err := client.GetAllHalkaArz(ctx, 1, 10, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().GreaterOrEqual(resp.RecordCount, 0)

	if len(resp.Items) > 0 {
		item := resp.Items[0]
		s.Require().NotZero(item.ID)
		s.Require().NotEmpty(item.CompanyName)
		s.Require().NotEmpty(item.Currency)
		s.Require().NotEmpty(item.Status)
	}
}

func (s *HalkaArzTestSuite) TestGetHalkaArzByID() {
	client := newTestClient(s.Config)
	ctx := context.Background()

	all, err := client.GetAllHalkaArz(ctx, 1, 10, RegionTr)
	s.Require().NoError(err)
	s.Require().NotNil(all)

	if len(all.Items) == 0 {
		s.T().Skip("no IPO offerings available to fetch by id")
	}

	id := all.Items[0].ID
	resp, err := client.GetHalkaArzByID(ctx, id)
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(id, resp.ID)
	s.Require().NotEmpty(resp.CompanyName)
}
