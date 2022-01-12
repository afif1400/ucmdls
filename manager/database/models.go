package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type Lab struct {
	ID      primitive.ObjectID   `json:"id" bson:"_id"`
	Name    string               `json:"name" bson:"name"`
	Branch  string               `json:"branch" bson:"branch"`
	Systems []primitive.ObjectID `json:"systems" bson:"systems"`
}

type System struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name,omitempty"`
	UserName  string             `json:"user_name" bson:"user_name,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	IpAddress string             `json:"ip_address" bson:"ip_address,omitempty"`
	ROM       string             `json:"rom" bson:"rom,omitempty"`
	RAM       string             `json:"ram" bson:"ram,omitempty"`
	OS        string             `json:"os" bson:"os,omitempty"`
}
