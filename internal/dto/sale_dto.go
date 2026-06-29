package dto

import "CotoChallenge/internal/domain"

type ModelPercentage struct {
	Vehicle    domain.VehicleType `json:"vehicle"`
	Percentage float64            `json:"percentage"`
}

type CenterModelPercentage struct {
	Center string            `json:"center"`
	Models []ModelPercentage `json:"models"`
}
