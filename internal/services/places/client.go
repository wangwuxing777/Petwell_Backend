package places

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const searchNearbyURL = "https://places.googleapis.com/v1/places:searchNearby"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type Place struct {
	Name             string   `json:"name"`
	Address          string   `json:"address"`
	Lat              float64  `json:"lat"`
	Lng              float64  `json:"lng"`
	Status           string   `json:"businessStatus"`
	OpenNow          *bool    `json:"openNow,omitempty"`
	Rating           *float64 `json:"rating,omitempty"`
	UserRatingsTotal *int     `json:"userRatingsTotal,omitempty"`
}

type openingHoursInfo struct {
	OpenNow bool `json:"openNow"`
}

type googlePlace struct {
	DisplayName struct {
		Text string `json:"text"`
	} `json:"displayName"`
	FormattedAddress string `json:"formattedAddress"`
	Location         struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	BusinessStatus      string            `json:"businessStatus"`
	RegularOpeningHours *openingHoursInfo `json:"regularOpeningHours,omitempty"`
	CurrentOpeningHours *openingHoursInfo `json:"currentOpeningHours,omitempty"`
	Rating              *float64          `json:"rating,omitempty"`
	UserRatingCount     *int              `json:"userRatingCount,omitempty"`
}

type searchResponse struct {
	Places []googlePlace `json:"places"`
}

func (c *Client) SearchNearbyVets(lat, lng, radius float64) ([]Place, error) {
	reqBody := map[string]interface{}{
		"includedTypes": []string{"veterinary_care"},
		"locationRestriction": map[string]interface{}{
			"circle": map[string]interface{}{
				"center": map[string]interface{}{
					"latitude":  lat,
					"longitude": lng,
				},
				"radius": radius,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", searchNearbyURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "places.displayName,places.formattedAddress,places.location,places.businessStatus,places.regularOpeningHours,places.currentOpeningHours,places.rating,places.userRatingCount")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google places api returned status: %d", resp.StatusCode)
	}

	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var places []Place
	for _, p := range searchResp.Places {
		// Prefer currentOpeningHours (real-time) over regularOpeningHours (scheduled)
		var isOpen *bool
		if p.CurrentOpeningHours != nil {
			isOpen = &p.CurrentOpeningHours.OpenNow
		} else if p.RegularOpeningHours != nil {
			isOpen = &p.RegularOpeningHours.OpenNow
		}

		places = append(places, Place{
			Name:             p.DisplayName.Text,
			Address:          p.FormattedAddress,
			Lat:              p.Location.Latitude,
			Lng:              p.Location.Longitude,
			Status:           p.BusinessStatus,
			OpenNow:          isOpen,
			Rating:           p.Rating,
			UserRatingsTotal: p.UserRatingCount,
		})
	}

	return places, nil
}

const searchTextURL = "https://places.googleapis.com/v1/places:searchText"

func (c *Client) SearchTextVets(query string) ([]Place, error) {
	reqBody := map[string]interface{}{
		"textQuery":    query,
		"includedType": "veterinary_care",
		"locationBias": map[string]interface{}{
			"circle": map[string]interface{}{
				"center": map[string]interface{}{
					"latitude":  22.3193, // Center of Hong Kong
					"longitude": 114.1694,
				},
				"radius": 50000.0, // 50km radius to cover all of HK
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", searchTextURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)
	req.Header.Set("X-Goog-FieldMask", "places.displayName,places.formattedAddress,places.location,places.businessStatus,places.regularOpeningHours,places.currentOpeningHours,places.rating,places.userRatingCount")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google places api returned status: %d", resp.StatusCode)
	}

	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var places []Place
	for _, p := range searchResp.Places {
		// Prefer currentOpeningHours (real-time) over regularOpeningHours (scheduled)
		var isOpen *bool
		if p.CurrentOpeningHours != nil {
			isOpen = &p.CurrentOpeningHours.OpenNow
		} else if p.RegularOpeningHours != nil {
			isOpen = &p.RegularOpeningHours.OpenNow
		}

		places = append(places, Place{
			Name:             p.DisplayName.Text,
			Address:          p.FormattedAddress,
			Lat:              p.Location.Latitude,
			Lng:              p.Location.Longitude,
			Status:           p.BusinessStatus,
			OpenNow:          isOpen,
			Rating:           p.Rating,
			UserRatingsTotal: p.UserRatingCount,
		})
	}

	return places, nil
}
