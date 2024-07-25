package dto

type AddressDto struct {
	Cep        string `json:"cep"`
	Localidade string `json:"localidade"`
	Error      string `json:"erro"`
}
