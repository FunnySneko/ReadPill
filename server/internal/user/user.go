package user

type User struct {
	ID                      int
	Username                string
	Email                   string
	PasswordHash            string
	AverageContributeRating float32
	AverageBias             float32
}
