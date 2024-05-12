package db

import (
	"github.com/RaphaelHebert/DailyDices-BE/model"
	"github.com/google/uuid"
)

var MockUUID = uuid.NewString()

var UsersList = model.Users{
		MockUUID: {
			Username:  "Joe",
			Email:    "joe@mymail.com",
			Password: "someHash",
		},
		// Add more users as needed
	}

var Scores = model.Scores{
	"scoreId234deq": model.Score{
		Score: []model.Dice{1, 3, 5},
		UID: "23wfsw4rfser34",
	},
	// Add more scores as needed
}