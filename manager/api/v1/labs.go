package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/afif1400/ucmdls/manager/database"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// HandleLabs is a handler function for /api/manager/labs
func HandleLabs(dbc *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var labs []database.Lab

		// bson.M{}
		res, err := dbc.Find(context.TODO(), bson.M{})

		if err != nil {
			database.Error(err, w)
			return
		}

		defer res.Close(context.TODO())

		for res.Next(context.TODO()) {

			var lab database.Lab

			err := res.Decode(&lab)

			if err != nil {
				database.Error(err, w)
				return
			}

			labs = append(labs, lab)
		}

		if err := res.Err(); err != nil {
			database.Error(err, w)
			return
		}

		json.NewEncoder(w).Encode(labs)

	}
}

// HandleLab is a handler for /api/manager/labs/{labID}
func HandleLab(dbc *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// get the labID from the url
		vars := mux.Vars(r)
		labID := vars["labID"]

		// get the lab from the database
		var lab database.Lab

		err := dbc.FindOne(context.TODO(), bson.M{"labID": labID}).Decode(&lab)

		if err != nil {
			database.Error(err, w)
			return
		}

		json.NewEncoder(w).Encode(lab)
	}
}

// HandleLabCreate is a handler for /api/manager/labs/create
func HandleLabCreate(dbc *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var lab database.Lab

		_ = json.NewDecoder(r.Body).Decode(&lab)

		// insert the lab into the database
		result, err := dbc.InsertOne(context.TODO(), lab)

		if err != nil {
			database.Error(err, w)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

// HandleLabAddSystems is a handler for /api/manager/labs/{labID}/systems
func HandleLabAddSystems(dbc *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// get the labID from the url
		vars := mux.Vars(r)
		labID := vars["labID"]

		// get the lab from the database
		var lab database.Lab

		err := dbc.FindOne(context.TODO(), bson.M{"labID": labID}).Decode(&lab)

		if err != nil {
			database.Error(err, w)
			return
		}

		// get the systems from the body
		var systems []database.System

		_ = json.NewDecoder(r.Body).Decode(&systems)

		// insert systems to database one by one and append to lab.Systems
		for _, system := range systems {
			result, err := dbc.InsertOne(context.TODO(), system)

			if err != nil {
				database.Error(err, w)
				return
			}

			lab.Systems = append(lab.Systems, result.InsertedID.(primitive.ObjectID))
		}

		json.NewEncoder(w).Encode(lab)

	}
}
