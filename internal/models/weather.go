package models

import "time"

type WeatherData struct {
	Location    string    `json:"location"`
	Temperature float64   `json:"temperature"`
	Condition   string    `json:"condition"`
	Icon        string    `json:"icon"`
	LastUpdated time.Time `json:"last_updated"`
}
