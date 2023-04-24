package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
)

func readEgaugeData(egaugeJwT string, deviceName string) (map[string]interface{}, error) {
	auth := http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", egaugeJwT)},
	}
	config := &http.Request{
		Method: "GET",
		Header: auth,
		URL: &url.URL{
			Scheme: "https",
			Host:   fmt.Sprintf("%s.d.egauge.net", deviceName),
			Path:   "/api/register",
		},
		Form: url.Values{
			"time": []string{"now"},
			"rate": []string{""},
		},
	}
	response, err := http.DefaultClient.Do(config)
	if err != nil {
		fmt.Printf("%v can't get data from egauge\n", err)
		return nil, err
	}
	defer response.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		fmt.Printf("%v can't parse response from egauge\n", err)
		return nil, err
	}
	return data, nil
}

func egaugeLogin(deviceName string, usr string, pwd string) (string, error) {
	unauthorizedResponse, err := http.Get(fmt.Sprintf("https://%s.d.egauge.net/api/auth/unauthorized"))
	if err != nil {
		fmt.Printf("%v can't connect to egauge\n", err)
		return "", err
	}
	defer unauthorizedResponse.Body.Close()

	var unauthorizedData struct {
		Rlm string `json:"rlm"`
		Nnc string `json:"nnc"`
	}
	err = json.NewDecoder(unauthorizedResponse.Body).Decode(&unauthorizedData)
	if err != nil {
		fmt.Printf("%v can't parse unauthorized response from egauge\n", err)
		return "", err
	}

	rlm := unauthorizedData.Rlm
	cnnc := fmt.Sprintf("%d", rand.Intn(999999999))

	hasher := md5.New()
	hasher.Write([]byte(fmt.Sprintf("%s:%s:%s", usr, rlm, pwd)))
	ha1 := hex.EncodeToString(hasher.Sum(nil))

	hasher.Reset()
	hasher.Write([]byte(fmt.Sprintf("%s:%s:%s", ha1, unauthorizedData.Nnc, cnnc)))
	ha2 := hex.EncodeToString(hasher.Sum(nil))

	loginData := map[string]interface{}{
		"rlm":  rlm,
		"usr":  usr,
		"cnnc": cnnc,
		"nnc":  unauthorizedData.Nnc,
		"hash": ha2,
	}

	requestBody, err := json.Marshal(loginData)
	if err != nil {
		fmt.Printf("%v can't serialize login data for egauge\n", err)
		return "", err
	}

	response, err := http.Post(fmt.Sprintf("https://%s.d.egauge.net/api/auth/login", deviceName), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("%v can't login to egauge\n", err)
		return "", err
	}
	defer response.Body.Close()

	var responseData struct {
		Jwt string `json:"jwt"`
	}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		fmt.Printf("%v can't parse login")
	}
}
