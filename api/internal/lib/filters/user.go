package filters

type UsersFilters struct {
	Limit      *int
	Offset     *int
	Surname    *string
	Name       *string
	Patronymic *string
	Address    *string
}
