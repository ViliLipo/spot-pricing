package shellyhttp

import (
	"log"
	"net/http"
	"github.com/ViliLipo/spotpricing/internal/env"
)

func SetRelay(state bool) {
	shellyHost := env.ShellyRelayHost.GetValue()
	var val string
	if state {
		val = "true"
	} else {
		val = "false"
	}
	res, err := http.Get("http://" + shellyHost + "/rpc/Switch.Set?id=0&on=" + val)
	if err != nil {
		log.Fatal(err.Error())
	}
	if res.StatusCode != 200 {
		log.Println("Non 200 status was returned, " + res.Status)
	} else {
		log.Println("Set relay to state: " + val +".")
	}
}
