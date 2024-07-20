package dto

import "fmt"

type LocaleInput struct {
	Cep string `json:"cep"`
}

type LocaleOutput struct {
	Localidade string `json:"localidade"`
}

func (o LocaleOutput) String() string {
	return fmt.Sprintf("{ Localidade: %s }", o.Localidade)
}

type ErrorOutput struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
