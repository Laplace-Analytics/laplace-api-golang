package laplace

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type HalkaArz struct {
	ID                   int        `json:"id"`
	CompanyName          string     `json:"companyName"`
	Symbol               *string    `json:"symbol"`
	InstrumentID         *int       `json:"instrumentId"`
	PriceMin             *float64   `json:"priceMin"`
	PriceMax             *float64   `json:"priceMax"`
	DemandStartDate      *time.Time `json:"demandStartDate"`
	DemandEndDate        *time.Time `json:"demandEndDate"`
	FirstTradingDate     *time.Time `json:"firstTradingDate"`
	SharesOffered        *float64   `json:"sharesOffered"`
	OfferingSize         *float64   `json:"offeringSize"`
	OfferingType         *string    `json:"offeringType"`
	ConsortiumLeader     *string    `json:"consortiumLeader"`
	AdditionalShares     *float64   `json:"additionalShares"`
	DistributionMethod   *string    `json:"distributionMethod"`
	FreeFloatRate        *float64   `json:"freeFloatRate"`
	IntendedMarket       *string    `json:"intendedMarket"`
	Sector               *string    `json:"sector"`
	MaxLotPerInvestor    *float64   `json:"maxLotPerInvestor"`
	Currency             string     `json:"currency"`
	RelatedDisclosureIDs []int      `json:"relatedDisclosureIds"`
	Reviewed             bool       `json:"reviewed"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	Status               string     `json:"status"`
	IsFixedPrice         bool       `json:"isFixedPrice"`
}

// GetAllHalkaArz retrieves all IPO (halka arz) offerings with pagination.
func (c *Client) GetAllHalkaArz(ctx context.Context, page int, pageSize int, region Region) (*PaginatedResponse[HalkaArz], error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/ipo/all", c.baseUrl), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	q.Add("region", string(region))
	req.URL.RawQuery = q.Encode()

	resp, err := sendRequest[PaginatedResponse[HalkaArz]](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetHalkaArzByID fetches a single IPO (halka arz) offering by its ID.
func (c *Client) GetHalkaArzByID(ctx context.Context, id int) (*HalkaArz, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/ipo/%d", c.baseUrl, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := sendRequest[HalkaArz](ctx, c, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
