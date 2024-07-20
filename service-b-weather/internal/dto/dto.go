package dto

import "fmt"

type LocaleInput struct {
	Cep string `json:"cep"`
}

type LocaleOutput struct {
	Localidade string `json:"localidade"`
}

type WeatherOutput struct {
	Current CurrentWeather `json:"current"`
}

type CurrentWeather struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type ResultOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func (o ResultOutput) String() string {
	return fmt.Sprintf("{ TempC: %f, TempF: %f, TempK: %f }", o.TempC, o.TempF, o.TempK)
}

type ErrorOutput struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
