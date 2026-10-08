package dataaccess

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type SpotPricingDataPointRow struct {
	// EUR CENTS PER KILOWATT HOUR
	Price float64
	// StartDate as 64bit unix timestamp
	StartDate int64
	// EndDate as 64bit unix timestamp
	EndDate int64
}

func InitDatabase(dbUrl string) {
	db, err := sql.Open("sqlite", dbUrl)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
		return
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS spot_pricing (pk INTEGER PRIMARY KEY ASC, price REAL, start_date INTEGER, end_date INTEGER )")
	if err != nil {
		log.Fatalln(err)
		return
	}
	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS spot_pricing_start_date on spot_pricing(start_date)")
	if err != nil {
		log.Fatalln(err)
		return
	}
	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS spot_pricing_end_date on spot_pricing(end_date)")
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func StoreSpotPricingDataPoints(spotPricingDataPoints []SpotPricingDataPointRow, dbUrl string) {
	db, err := sql.Open("sqlite", dbUrl)
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
		return
	}
	tx, _ := db.Begin()
	for _, dataPoint := range spotPricingDataPoints {
		tx.Exec(`INSERT INTO spot_pricing (price, start_date, end_date) VALUES (?,?,?)`, dataPoint.Price, dataPoint.StartDate, dataPoint.EndDate)
	}
	tx.Commit()
}

func GetCurrentPrice(timeStamp int64, dbUrl string) (float64, error) {
	db, err := sql.Open("sqlite", dbUrl)
	defer db.Close()
	if err != nil {
		return 0.0, err
	}
	var price float64
	err = db.QueryRow("SELECT price FROM spot_pricing WHERE start_date < $1 and end_date > $1 ORDER BY pk DESC LIMIT 1",
		timeStamp).Scan(&price)
	if err != nil {
		return 0.0, err
	}
	return price, nil
}

// Read the startdate for price that is least in the past or most in the future...
func GetLastRowStartDate(dbUrl string) int64 {
	db, err := sql.Open("sqlite", dbUrl)
	defer db.Close()
	if err != nil {
		return 0
	}
	var row int64
	err = db.QueryRow("SELECT start_date FROM spot_pricing ORDER BY start_date DESC LIMIT 1").Scan(&row)
	if err != nil {
		return 0
	}
	return row
}
