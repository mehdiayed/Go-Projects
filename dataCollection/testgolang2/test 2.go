package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

func egaugeLogin(deviceName string, usr string, pwd string) (string, error) {
	unauthResp, err := http.Get(fmt.Sprintf("https://%s.d.egauge.net/api/auth/unauthorized", deviceName))
	if err != nil {
		return "", fmt.Errorf("%v can't connect to egauge", err)
	}
	defer unauthResp.Body.Close()

	unauthRespData := make(map[string]string)
	if err = unmarshalJSON(unauthResp.Body, &unauthRespData); err != nil {
		return "", fmt.Errorf("error decoding JSON: %v", err)
	}

	rlm := unauthRespData["rlm"]
	nnc := unauthRespData["nnc"]
	cnnc := fmt.Sprintf("%d", time.Now().UnixNano())

	ha1 := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", usr, rlm, pwd)))
	ha2 := md5.Sum([]byte(fmt.Sprintf("%x:%s:%s", ha1, nnc, cnnc)))

	loginBody := map[string]string{
		"rlm":  rlm,
		"usr":  usr,
		"cnnc": cnnc,
		"nnc":  nnc,
		"hash": hex.EncodeToString(ha2[:]),
	}

	config := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := config.Post(fmt.Sprintf("https://%s.d.egauge.net/api/auth/login", deviceName), "application/json", marshalJSON(loginBody))
	if err != nil {
		return "", fmt.Errorf("%v can't connect to egauge", err)
	}
	defer resp.Body.Close()

	loginRespData := make(map[string]string)
	if err = unmarshalJSON(resp.Body, &loginRespData); err != nil {
		return "", fmt.Errorf("error decoding JSON: %v", err)
	}

	return loginRespData["jwt"], nil
}

func main() {
	deviceName := "egauge67897"
	usr := "owner"
	pwd := "000000"

	jwt, err := egaugeLogin(deviceName, usr, pwd)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	data, err := readEgaugeData(jwt, deviceName)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("%s\n", data)
}

// helper functions for marshalling and unmarshalling JSON data
func marshalJSON(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)
	return jsonData
}
