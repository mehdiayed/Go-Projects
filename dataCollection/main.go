package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

func main() {
	// Connect to Kafka broker
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, config)
	if err != nil {
		fmt.Println("Error creating Kafka producer:", err)
		return
	}
	defer producer.Close()

	// Connect to eGauge register
	conn, err := net.Dial("tcp", "192.168.1.10:80")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Send request to eGauge
	fmt.Fprintf(conn, "GET /cgi-bin/egauge?inst&csv HTTP/1.0\r\n\r\n")

	// Read response
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Date") {
			// Parse CSV data
			reader := csv.NewReader(strings.NewReader(line))
			reader.TrimLeadingSpace = true
			records, err := reader.ReadAll()
			if err != nil {
				fmt.Println("Error parsing CSV:", err)
				return
			}

			// Send data to Kafka
			for _, record := range records {
				key := sarama.StringEncoder(record[0])
				value := sarama.StringEncoder(record[1])
				_, _, err := producer.SendMessage(&sarama.ProducerMessage{
					Topic: "egauge",
					Key:   key,
					Value: value,
				})
				if err != nil {
					fmt.Println("Error sending message to Kafka:", err)
				}
			}
