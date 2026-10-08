package env

import "os"

type EnvKey string

func (key EnvKey) GetValue() string {
	return os.Getenv(string(key))
}

const (
	SqliteDbUri                  EnvKey = "SPOT_PRICING_SQLITE_DB_URI"
	SpotPricingApiUrl            EnvKey = "SPOT_PRICING_API_URL"
	ShellyRelayHost              EnvKey = "SPOT_PRICING_SHELLY_HOST"
	ElectricityBillMarginEurCent EnvKey = "SPOT_PRICING_ELECTRICITY_BILL_MARGIN"
	OilPricePerLiterEurCent      EnvKey = "SPOT_PRICING_OIL_PRICE"
)
