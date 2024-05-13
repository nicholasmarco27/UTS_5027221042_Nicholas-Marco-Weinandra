package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//
const MenuCollection = "menu"

//
type Menu struct {
	ID			primitive.ObjectID	`bson:"_id,omitempty"`
	Title		string				`bson:"title"`
	Description	string        		`bson:"description"`
}