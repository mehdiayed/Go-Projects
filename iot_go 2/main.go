package main

import (
	"examlpe/influx-connection/database"
	"fmt"
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

func main() {

	client, err := database.Init()

	if err != nil {
		log.Fatal(err)
		return
	}

	defer database.Close() // bech akhir 7aja titsakkar

	writeAPI := influxdb.Client.WriteAPI(client, "iot", "measurements_business_event")

	for i := 0; i < 10; i++ {

		m := database.Measurements{
			ID:        i,
			Name:      fmt.Sprintf("measurement_%d", i),
			TimeStamp: time.Now().UnixNano(),
			Value:     i * 12,
		}

		p := influxdb.NewPointWithMeasurement("device_1").
			AddTag("id", fmt.Sprintf("%d", m.ID)).
			AddTag("name", m.Name).
			AddField("value", m.Value).
			SetTime(time.Unix(0, m.TimeStamp))

		writeAPI.WritePoint(p)
	}

}
