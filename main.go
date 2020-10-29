package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"github.com/sidhant/IMDBApiGo/connectionbuilder"	
	"github.com/sidhant/IMDBApiGo/models"	
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	//Init Router
	r := mux.NewRouter()

  	// arrange our route
	  r.HandleFunc("/api/movies", getMoviesAll).Methods("GET")
	  r.HandleFunc("/api/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/api/movies", createMovie).Methods("POST")
	r.HandleFunc("/api/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/api/movies/{id}", deleteMovie).Methods("DELETE")	
	// set our port address
	log.Fatal(http.ListenAndServe(":8000", r))

}

func getMoviesAll(w http.ResponseWriter, r *http.Request) {
	var collection = connectionbuilder.ConnectDB()
	w.Header().Set("Content-Type", "application/json")
	var MyMovieAll []models.MyMovie

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		connectionbuilder.GetError(err, w)
		return
	}

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var currMovie models.MyMovie
		err := cur.Decode(&currMovie) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}
		MyMovieAll = append(MyMovieAll, currMovie)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(MyMovieAll) // encode similar to serialize process.
}


func getMovie(w http.ResponseWriter, r *http.Request) {
	var collection = connectionbuilder.ConnectDB()
	var currMovie models.MyMovie
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&currMovie)

	if err != nil {
		connectionbuilder.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(currMovie)

}



func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var currMovie models.MyMovie
	var collection = connectionbuilder.ConnectDB()
	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&currMovie)

	// insert our book model.
	result, err := collection.InsertOne(context.TODO(), currMovie)
	if err != nil {
		connectionbuilder.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}




func deleteMovie(w http.ResponseWriter, r *http.Request) {
	// Set header
	w.Header().Set("Content-Type", "application/json")
	// get params
	var params = mux.Vars(r)
	var collection = connectionbuilder.ConnectDB()
	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])
	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		connectionbuilder.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var collection = connectionbuilder.ConnectDB()
	var currMovie models.MyMovie
	
	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&currMovie)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"popularity", currMovie.Popularity},
			{"director", currMovie.Director},
			{"genre",currMovie.Genre},
			{"ImdbScore",currMovie.ImdbScore},
			{"Name",currMovie.Name}},}}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&currMovie)

	if err != nil {
		connectionbuilder.GetError(err, w)
		return
	}

	currMovie.ID = id

	json.NewEncoder(w).Encode(currMovie)
}
