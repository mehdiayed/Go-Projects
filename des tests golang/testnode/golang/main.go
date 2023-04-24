package main

import (
	"fmt"
	"net/http"
	"time"
)

func readEgaugeData(egaugeJwT string, deviceName string) (string, error) {
	auth := http.Header{}
	auth.Set("Authorization", "Bearer "+egaugeJwT)

	config := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := config.Get(fmt.Sprintf("https://%s.d.egauge.net/api/register?time=now&rate=", deviceName))
	if err != nil {
		return "", fmt.Errorf("%v can't get data from egauge", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	if err != nil {
		return "", fmt.Errorf("%v can't read data from egauge", err)
	}

	return string(buf[:n]), nil
}
