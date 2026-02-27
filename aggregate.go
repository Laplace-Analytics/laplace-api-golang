package laplace

import (
	"context"
	"fmt"
	"net/http"
)

type CollectionPriceGraph struct {
	PreviousClose float64          `json:"previous_close"`
	Graph         []PriceDataPoint `json:"graph"`
}

type AggregatePricePeriod string

const (
	AggregatePricePeriodOneDay     AggregatePricePeriod = "1G"
	AggregatePricePeriodOneWeek    AggregatePricePeriod = "1H"
	AggregatePricePeriodOneMonth   AggregatePricePeriod = "1A"
	AggregatePricePeriodThreeMonth AggregatePricePeriod = "3A"
	AggregatePricePeriodOneYear    AggregatePricePeriod = "1Y"
	AggregatePricePeriodTwoYear    AggregatePricePeriod = "2Y"
	AggregatePricePeriodThreeYear  AggregatePricePeriod = "3Y"
	AggregatePricePeriodFiveYear   AggregatePricePeriod = "5Y"
)

// GetAggregateGraph retrieves the aggregate price graph for a sector, industry, or collection.
func (c *Client) GetAggregateGraph(ctx context.Context, period AggregatePricePeriod, region Region, sectorId, industryId, collectionId string) (CollectionPriceGraph, error) {
	endpoint := fmt.Sprintf("%s/api/v1/aggregate/graph", c.baseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return CollectionPriceGraph{}, err
	}

	q := req.URL.Query()
	q.Add("region", string(region))
	q.Add("period", string(period))
	if sectorId != "" {
		q.Add("sectorId", sectorId)
	}
	if industryId != "" {
		q.Add("industryId", industryId)
	}
	if collectionId != "" {
		q.Add("collectionId", collectionId)
	}
	req.URL.RawQuery = q.Encode()

	res, err := sendRequest[CollectionPriceGraph](ctx, c, req)
	if err != nil {
		return CollectionPriceGraph{}, err
	}

	return res, nil
}
