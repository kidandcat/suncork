package main

import (
	"encoding/json"
	"os"
)

func getConfig() Config {
	file, e := os.Open("config.json")
	var f Config
	if err(e) {
		return f
	}
	decoder := json.NewDecoder(file)
	decoder.Decode(&f)
	return f
}

type ShippingOption struct {
	Lang        string `json:"lang"`
	ID          string `json:"id"`
	Label       string `json:"label"`
	DeliverTime string `json:"deliverTime"`
	Amount      int    `json:"amount"`
	OneLine     string `json:"oneLine"`
}

type Config struct {
	ShippingOptions []ShippingOption
}
