package database

import (
	"context"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var (
	client   influxdb.Client //Client thenya hia interface en fin de conte tajam testi biha barcha 7ajet
	writeAPI api.WriteAPI    // WriteAPI interface zeda
)

//initialisation mta3 database

func Init() (influxdb.Client, error) {

	var err error
	//NewClientWithOptions mich yibda yisna3 fil connexion mta3o
	client = influxdb.NewClientWithOptions(InfluxDBURL, InfluxDBToken, influxdb.DefaultOptions().SetBatchSize(100)) // kol fil config  // l DefaultOptions().SetBatchSize(100) : hia 9a dech bech tijbid data , m3a DefaultOptions() 9adech bech tzid 3liha
	// 9adech min point mich yakhouha min request wa7da
	// batch hia sequance mta3 des données

	writeAPI = client.WriteAPI(InfluxDBDatabase, "") // isna3li l BDD

	_, err = client.Ready(context.Background())
	//wa9t ma ysir init() fil main y3adi houwa context lil package y3adi context mta3 l main , itha ken l client 'ready'
	//iwali l context maykharajich error or itha keni not ready l context bech ikharaj nil w bech ikharaj error w mich ikharaj problem
	if err != nil {
		return nil, err
	}

	return client, nil
}

func Close() {
	client.Close()
}
