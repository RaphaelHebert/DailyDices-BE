package model

type User struct {
	UID      string   `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PublicUser struct {
	UID      string   `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserPublic struct {
	UID      string   `json:"uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLogin struct {
	UserPublic
	password string
}

type Users map[string]User

type Dice int

type Score struct{
	Score []Dice `json:"score"`
	UID string  `json:"uid"`// uid of the score, to be replace by date
}

type Scores map[string][]Score



