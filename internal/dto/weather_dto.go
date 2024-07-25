package dto

type WeatherDto struct {
	Location weatherLocationDto
	Current  weatherCurrentDto
}

type weatherCurrentDto struct {
	TempC float64 `json:"temp_c"`
}

type weatherLocationDto struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}
