package laplace

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Politician struct {
	Id             int32     `json:"id"`
	PoliticianName string    `json:"politicianName"`
	TotalHoldings  int32     `json:"totalHoldings"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

type Holding struct {
	PoliticianName string    `json:"politicianName"`
	Symbol         string    `json:"symbol"`
	Company        string    `json:"company"`
	Holding        string    `json:"holding"`
	Allocation     string    `json:"allocation"`
	LastUpdated    time.Time `json:"lastUpdated"`
}

type HoldingShort struct {
	Symbol     string `json:"symbol"`
	Company    string `json:"company"`
	Holding    string `json:"holding"`
	Allocation string `json:"allocation"`
}

type TopHolding struct {
	Symbol      string                 `json:"symbol"`
	Company     string                 `json:"company"`
	Politicians []TopHoldingPolitician `json:"politicians"`
	Count       int32                  `json:"count"`
}

type TopHoldingPolitician struct {
	Name       string `json:"name"`
	Holding    string `json:"holding"`
	Allocation string `json:"allocation"`
}

type PoliticianDetail struct {
	Id            int32          `json:"id"`
	Name          string         `json:"name"`
	Holdings      []HoldingShort `json:"holdings"`
	TotalHoldings int32          `json:"totalHoldings"`
	LastUpdated   time.Time      `json:"lastUpdated"`
}

func (c *Client) GetAllPoliticians(ctx context.Context) ([]Politician, error) {
	endpoint := fmt.Sprintf("%s/api/v1/politician", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequest[[]Politician](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetPoliticianHoldingsBySymbol(ctx context.Context, symbol string) ([]Holding, error) {
	endpoint := fmt.Sprintf("%s/api/v1/holding/%s", c.baseUrl, symbol)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequest[[]Holding](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetAllTopHoldings(ctx context.Context) ([]TopHolding, error) {
	endpoint := fmt.Sprintf("%s/api/v1/top-holding", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequest[[]TopHolding](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GetPoliticianDetail(ctx context.Context, id int) (PoliticianDetail, error) {
	endpoint := fmt.Sprintf("%s/api/v1/politician/%d", c.baseUrl, id)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return PoliticianDetail{}, err
	}

	resp, err := sendRequest[PoliticianDetail](ctx, c, req)
	if err != nil {
		return PoliticianDetail{}, err
	}

	return resp, nil
}
