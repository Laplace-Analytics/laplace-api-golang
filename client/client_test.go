package client

import (
	"context"
	"net/http"
	"testing"

	"finfree.co/laplace/utilities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type LaplaceClientTestSuite struct {
	*utilities.ClientTestSuite
}

func TestLaplaceClient(t *testing.T) {
	suite.Run(t, &LaplaceClientTestSuite{
		utilities.NewClientTestSuite(),
	})
}

func (s *LaplaceClientTestSuite) TestClient() {
	client := NewClient(s.Config, logrus.New())

	req, err := http.NewRequest("GET", s.Config.BaseURL+"/api/v1/industry", nil)
	s.Require().NoError(err)
	q := req.URL.Query()
	q.Add("region", string(RegionTr))
	q.Add("locale", string(LocaleTr))
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[any](context.Background(), client, req)
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
}
