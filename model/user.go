package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UID      string `json:"uid" bson:"_id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
	IsAdmin  bool `json:"isadmin" bson:"isadmin"`
	Scores 	[]Score `json:"scores" bson:"scores"`
	Password string `json:"password"`
}

type PublicUser struct {
	UID      string `json:"uid" bson:"_id,omitempty"`
	Username string `json:"username"`
	IsAdmin  bool `json:"isadmin" bson:"isadmin"`
	Scores []Score `json:"scores" bson:"scores"`
	Email    string `json:"email"`
}

type LoginInput struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type Users map[string]User

type Dice int

type Score struct{
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Date  int64              `json:"date" bson:"date"`
	Score []Dice             `json:"score"`
}

type Scores []Score



