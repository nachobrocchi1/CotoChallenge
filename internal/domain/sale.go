package domain

import "time"

type Sale struct {
	Id      string      `json:"id"`
	Vehicle VehicleType `json:"vehicle"`
	Price   float64     `json:"price"`
	Date    time.Time   `json:"date"`
	Center  string      `json:"center"`
}
