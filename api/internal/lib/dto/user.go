package dto

type AddUserDto struct {
	PassportSerie  int
	PassportNumber int
}

type UserInfoDto struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

type SaveUserDto struct {
	*AddUserDto
	*UserInfoDto
}
