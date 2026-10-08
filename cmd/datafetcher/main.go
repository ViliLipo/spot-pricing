package main

import (
	"log"
	"time"
	"github.com/ViliLipo/spotpricing/internal/api-client"
	"github.com/ViliLipo/spotpricing/internal/data-access"
	"github.com/ViliLipo/spotpricing/internal/env"
)


func main() {
	dbUri := env.SqliteDbUri.GetValue()
	spotPricingApiUrl := env.SpotPricingApiUrl.GetValue()
	log.Println(spotPricingApiUrl)
	dataaccess.InitDatabase(dbUri)
	latestPriceStartDate := dataaccess.GetLastRowStartDate(dbUri)
	config := apiclient.DownloadConfiguration{SpotPricingEndpointUrl: spotPricingApiUrl}
	prices := apiclient.DownloadSpotPricingData(config)
	convertedPrices := convertSpotPrices(prices)
	filteredPrices := filterPrices(convertedPrices, latestPriceStartDate)
	dataaccess.StoreSpotPricingDataPoints(filteredPrices, dbUri)
	log.Printf("Successfully stored %d rows of spotpricing information to database.\n", len(filteredPrices))
}

func filterPrices(prices []dataaccess.SpotPricingDataPointRow, latestPriceStartDate int64) []dataaccess.SpotPricingDataPointRow {
	var filteredPrices []dataaccess.SpotPricingDataPointRow

	for _, price := range prices {
		if price.StartDate > latestPriceStartDate {
			filteredPrices = append(filteredPrices, price)
		}
	}
	return filteredPrices
}

func convertSpotPrices(apiPrices []apiclient.SpotPricingDataPoint) []dataaccess.SpotPricingDataPointRow {
	var datarows []dataaccess.SpotPricingDataPointRow
	for _, item := range apiPrices {
		price := item.Price
		startDate, err := dateStrToUnixStamp(item.StartDate)
		if err != nil {
			log.Fatalln("Failed to parse startDate " + item.StartDate + " as time.")
			break
		}
		endDate, err := dateStrToUnixStamp(item.EndDate)
		if err != nil {
			log.Fatalln("Failed to parse endDate " + item.EndDate + " as time.")
			break
		}
		datarows = append(datarows, dataaccess.SpotPricingDataPointRow{
			Price:     price,
			StartDate: startDate,
			EndDate:   endDate,
		})
	}
	return datarows
}

func dateStrToUnixStamp(timeStr string) (int64, error) {
	d, err := time.Parse("2006-01-02T15:04:05.000Z", timeStr)
	if err != nil {
		return 0, err
	}
	return d.Unix(), nil
}
