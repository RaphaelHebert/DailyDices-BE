package model

type User struct {
	UID      string   `json:"Uid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserPublic struct {
	UID      string   `json:"Uid"`
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
	Score []Dice
	UID string
}

type Scores map[string]Score



