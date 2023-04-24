package database

type Measurements struct {
	ID        int    `influx:"id"`
	Name      string `influx:"name"`
	TimeStamp int64  `influx:"timestamp"`
	Value     int    `influx:"value"`
}
