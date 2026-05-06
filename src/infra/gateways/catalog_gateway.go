package gateways

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CatalogProductInfo struct {
	ID         string `json:"id"`
	Category   string `json:"category"`
	Popularity int    `json:"popularity"`
}

type CatalogGateway struct {
	baseURL string
	client  *http.Client
}

func NewCatalogGateway(baseURL string) *CatalogGateway {
	return &CatalogGateway{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (g *CatalogGateway) GetProductInfo(ctx context.Context, productID string) (*CatalogProductInfo, error) {
	url := fmt.Sprintf("%s/catalog/products/%s", g.baseURL, productID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog gateway returned status %d", resp.StatusCode)
	}

	var info CatalogProductInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}
