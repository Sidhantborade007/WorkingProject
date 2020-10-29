package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Creating Struct for movies
type MyMovie struct {
	ID     primitive.ObjectID `json:"_id" bson:"_id"`
	Popularity int64    `json:"popularity" json:"popularity" `
	Director        string   `json:"director" json:"director"`
	Genre           []string `json:"genre" json:"genre"`
	ImdbScore       float64  `json:"imdb_score" json:"imdb_score"`
	Name            string   `json:"name" json:"name"`
}

