package main

import (
	"log"
	"os"
	"strconv"
	"time"
	"github.com/ViliLipo/spotpricing/internal/data-access"
	"github.com/ViliLipo/spotpricing/internal/env"
	shellyhttp "github.com/ViliLipo/spotpricing/internal/shelly_http"
)

const (
	BURNER_EFFICIENCY_PERCENTAGE     = 85.0
	OIL_ENERGY_DENSITY_PER_LITER_KWH = 10.4
)

func calculateThreshold(oilPricePerLiterEurCents, margin float64,) float64 {
	threshold := oilPricePerLiterEurCents / (OIL_ENERGY_DENSITY_PER_LITER_KWH * (BURNER_EFFICIENCY_PERCENTAGE / 100)) - margin
	return threshold
}

func main() {
	oilPricePerLiter, err := strconv.ParseFloat(env.OilPricePerLiterEurCent.GetValue(), 64)
	if err != nil {
		log.Fatal("Failed to parse oilPrice")
		os.Exit(1)
	}
	ebMargin, err := strconv.ParseFloat(env.ElectricityBillMarginEurCent.GetValue(), 64)
	if err != nil {
		log.Fatal("Failed to parse electricity bill margin.")
		os.Exit(1)
	}
	dbUri := env.SqliteDbUri.GetValue()
	currentTime := time.Now().UTC()
	// Lets add some margin to accomodate some clock drifting
	// We are sending the price for next 15 minutes.
	adjustedTime := currentTime.Add(time.Second * 15).UTC()
	currentUnix := adjustedTime.Unix()
	currentPrice, err := dataaccess.GetCurrentPrice(currentUnix, dbUri)
	if err != nil {
		log.Fatalln("Failed to read current price from database.")
	}
	log.Printf("Read current price %f.\n", currentPrice)
	threshold := calculateThreshold(oilPricePerLiter, ebMargin)
	log.Printf("Calculated threshold %f,\n", threshold)
	setting := currentPrice > threshold
	shellyhttp.SetRelay(setting)
}
