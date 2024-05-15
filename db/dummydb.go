package db

import (
	"github.com/RaphaelHebert/DailyDices-BE/model"
)


var UsersList = model.Users{
		"mockedUID": {
			Username:  "Joe",
			Email:    "joe@mymail.com",
			Password: "someHash",
			UID: "mockedUID",
		},
		// Add more users as needed
	}

var Scores = model.Scores{
	"mockedUID": {
		model.Score{
			Score: []model.Dice{1, 3, 5},
			UID: "23wfsw4rfser34",
		},
		model.Score{
			Score: []model.Dice{3, 3, 6},
			UID: "dfs4wrfsfge5",
		},
	},
		// Add more scores as needed
	}