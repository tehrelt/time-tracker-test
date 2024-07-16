package domain

type User struct {
	Id             string `json:"id" db:"id"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
	Patronymic     string `json:"patronymic" db:"patronymic"`
	Address        string `json:"address" db:"address"`
	PassportSerie  string `json:"passportSerie" db:"passport_serie"`
	PassportNumber string `json:"passportNumber" db:"passport_number"`
}
