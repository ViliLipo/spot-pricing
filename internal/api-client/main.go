package apiclient

import (
	"encoding/json"
	"log"
	"net/http"
)

type DownloadConfiguration struct {
	SpotPricingEndpointUrl string
}

type SpotPricingDataPoint struct {
	// EUR CENTS PER KILOWATT HOUR
	Price      float64 `json:"price"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
}

type SpotPricingDataResponse struct {
	Prices []SpotPricingDataPoint `json:"prices"`
}

// Fetches spot pricing data from url defined by configuration
// If the request or JSON decoding fails returns an empty array
// The caller should handle the retry logic.
func DownloadSpotPricingData(configuration DownloadConfiguration) []SpotPricingDataPoint {
	resp, err := http.Get(configuration.SpotPricingEndpointUrl)

	if err != nil {
		log.Fatalln(err)
		return make([]SpotPricingDataPoint, 0)
	}

	var data SpotPricingDataResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)

	if err != nil {
		log.Fatalln(err)
		return make([]SpotPricingDataPoint, 0)
	}

	return data.Prices
}
