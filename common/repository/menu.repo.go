package repository

import (
	"context"
	"log"
	"time"

	"github.com/nicholasmarco27/UTS_5027221042_Nicholas-Marco-Weinandra/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

func NewMenuRepo(db *mongo.Database) *MenuRepository {
	return &MenuRepository {
		db: db,
		col: db.Collection(model.MenuCollection),
	}
}

// Save menu
func (r *MenuRepository) Save(u *model.Menu) (model.Menu, error) {
	log.Printf("Save(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	var menu model.Menu
	res, err := r.col.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return menu, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&menu)
	if err != nil {
		log.Println(err)
		return menu, err
	}

	return menu, nil
}

// find all menus
func (r *MenuRepository) FindAll() ([]model.Menu, error) {
	log.Println("FindAll()")
	ctx, cancel := timeoutContext()
	defer cancel()

	var menus []model.Menu
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return menus, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var menu model.Menu
		err := cur.Decode(&menu)
		if err != nil {
			log.Println(err)
		}
		menus = append(menus, menu)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return menus, nil
}

// update menu 
func (r *MenuRepository) Update(u *model.Menu) (model.Menu, error) {
	log.Printf("Update(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	filter := bson.M{"_id": u.ID}
	update := bson.M{
		"$set": bson.M{
			"title":  u.Title,
			"description": u.Description,
		},
	}

	var menu model.Menu
	err := r.col.FindOneAndUpdate(ctx, filter, update).Decode(&menu)
	if err != nil {
		log.Printf("ERR 115 %v", err)
		return menu, err
	}

	return menu, nil
}

// delete menu by id
func (r  *MenuRepository) Delete(id string) (bool, error) {
	log.Printf("Delete(%s) \n", id)
	ctx, cancel := timeoutContext()
	defer cancel()

	var menu model.Menu
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.col.FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&menu)
	if err != nil {
		log.Printf("Fail to delete menu: %v \n", err)
		return false, err
	}
	log.Printf("Deleted_menu(%v) \n", menu)
	return true, nil
}

// buat timeout
func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
}