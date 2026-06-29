package repository

type ErrNoRows struct {
	message string
}

func (e ErrNoRows) Error() string {
	return "no rows"
}
