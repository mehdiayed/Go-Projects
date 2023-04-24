package main

import (
	"examlpe/influx-connection/database" // haka ya3ni jbidna les ficher ali fi database lkol
	"fmt"
	"log"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

func main() {

	client, err := database.Init()

	if err != nil {
		log.Fatal(err)
		return
	}

	defer database.Close() // bech akhir 7aja titsakkar
	// measurements_business_event  hia l backet mta3na
	writeAPI := influxdb.Client.WriteAPI(client, "iot", "measurements_business_event")

	for i := 0; i < 10; i++ {

		m := database.Measurements{
			ID:                i,
			TypeMesure:        "electricité",
			DesignationMesure: "DesignationMesure",
			ClasseMesure:      "ClasseMesure",
			StatuMesure:       "StatuMesure",
			CategorieMesure:   "CategorieMesure",
			MinimumMesure:     "MinimumMesure",
			Target:            "Target",
		}

		p := influxdb.NewPointWithMeasurement("device_1").
			AddTag("id", fmt.Sprintf("%d", m.ID)). //addtag hia lkey mta3ha string wil value string
			AddField("TypeMesure", m.TypeMesure).
			AddField("DesignationMesure", m.DesignationMesure).
			AddField("ClasseMesure", m.ClasseMesure).
			AddField("StatuMesure", m.StatuMesure).
			AddField("CategorieMesure", m.CategorieMesure).
			AddField("MinimumMesure", m.MinimumMesure).
			AddField("Target", m.Target)
		writeAPI.WritePoint(p)
	}

}
