package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type EgaugeData struct {
	// Define the structure of the response data
}

// type params struct {
// 	time, rate string
// }

// func readEgaugeData(egaugeJwT string, deviceName string) (EgaugeData, error) {
// 	// Define the function to read the eGauge data
// 	authHeader := http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", egaugeJwT)}}
// 	params := map[string]string{"time": "now", "rate": ""}

// 	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s.d.egauge.net/api/register", deviceName), nil)
// 	if err != nil {
// 		return EgaugeData{}, err
// 	}
// 	req.Header = authHeader
// 	req.URL.RawQuery = params
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return EgaugeData{}, err
// 	}
// 	defer resp.Body.Close()
// 	var egaugeData EgaugeData
// 	err = json.NewDecoder(resp.Body).Decode(&egaugeData)
// 	if err != nil {
// 		return EgaugeData{}, err
// 	}
// 	return egaugeData, nil
// }
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

//--------

func egaugeLogin(deviceName string, usr string, pwd string) string {
	unauthorizedURL := fmt.Sprintf("https://%s.d.egauge.net/api/auth/unauthorized", deviceName)

	// Send GET request to unauthorized endpoint and retrieve response
	client := &http.Client{}
	req, _ := http.NewRequest("GET", unauthorizedURL, nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error() + " can't connect to egauge")
		return "disconnected"
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	// Extract necessary data from response
	var rlm string
	var nnc string
	if strings.Contains(bodyString, "rlm") {
		rlmIndex := strings.Index(bodyString, "rlm") + 5
		rlm = bodyString[rlmIndex : rlmIndex+32]
	}
	if strings.Contains(bodyString, "nnc") {
		nncIndex := strings.Index(bodyString, "nnc") + 5
		nnc = bodyString[nncIndex : nncIndex+32]
	}

	// Generate random cnnc
	cnnc := fmt.Sprintf("%d", rand.Intn(999999999))

	// Calculate necessary hashes
	h1 := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", usr, rlm, pwd)))
	h1str := hex.EncodeToString(h1[:])
	h2 := md5.Sum([]byte(fmt.Sprintf("%s:%s:%s", h1str, nnc, cnnc)))
	h2str := hex.EncodeToString(h2[:])

	// Send POST request to login endpoint with necessary data
	loginURL := fmt.Sprintf("https://%s.d.egauge.net/api/auth/login", deviceName)
	data := fmt.Sprintf(`{"rlm":"%s","usr":"%s","cnnc":"%s","nnc":"%s","hash":"%s"}`, rlm, usr, cnnc, nnc, h2str)
	req, _ = http.NewRequest("POST", loginURL, strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error() + " can't connect to egauge")
		return "disconnected"
	}

	defer resp.Body.Close()
	bodyBytes, _ = ioutil.ReadAll(resp.Body)
	bodyString = string(bodyBytes)

	// Extract JWT from response and return it
	var jwt string
	if strings.Contains(bodyString, "jwt") {
		jwtIndex := strings.Index(bodyString, "jwt") + 5
		jwt = bodyString[jwtIndex : jwtIndex+401]
	}
	return jwt
}

const Device string = "egauge67897"
const user string = "owner"
const password string = "000000"

func main() {
	// await w asynchrone

	// c := make(chan string)
	// c := make(chan string)

	jwt := egaugeLogin(Device, user, password)
	// c <- jwt
	data, err := readEgaugeData(jwt, Device)

	if err != nil {
		fmt.Println("err:", err)
	}

	fmt.Println(data)

}
