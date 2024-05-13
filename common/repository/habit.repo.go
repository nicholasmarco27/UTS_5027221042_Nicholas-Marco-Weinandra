package repository

import (
	"context"
	"log"
	"time"

	"github.com/GabriellaErlinda/UTS_5027221018_Gabriella-Erlinda/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HabitRepository struct {
	db  *mongo.Database
	col *mongo.Collection
}

func NewHabitRepo(db *mongo.Database) *HabitRepository {
	return &HabitRepository {
		db: db,
		col: db.Collection(model.HabitCollection),
	}
}

// Save habit
func (r *HabitRepository) Save(u *model.Habit) (model.Habit, error) {
	log.Printf("Save(%v) \n", u)
	ctx, cancel := timeoutContext()
	defer cancel()

	var habit model.Habit
	res, err := r.col.InsertOne(ctx, u)
	if err != nil {
		log.Println(err)
		return habit, err
	}

	err = r.col.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&habit)
	if err != nil {
		log.Println(err)
		return habit, err
	}

	return habit, nil
}

// find all habits
func (r *HabitRepository) FindAll() ([]model.Habit, error) {
	log.Println("FindAll()")
	ctx, cancel := timeoutContext()
	defer cancel()

	var habits []model.Habit
	cur, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return habits, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var habit model.Habit
		err := cur.Decode(&habit)
		if err != nil {
			log.Println(err)
		}
		habits = append(habits, habit)
	}

	if err := cur.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return habits, nil
}

// update habit 
func (r *HabitRepository) Update(u *model.Habit) (model.Habit, error) {
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

	var habit model.Habit
	err := r.col.FindOneAndUpdate(ctx, filter, update).Decode(&habit)
	if err != nil {
		log.Printf("ERR 115 %v", err)
		return habit, err
	}

	return habit, nil
}

// delete habit by id
func (r  *HabitRepository) Delete(id string) (bool, error) {
	log.Printf("Delete(%s) \n", id)
	ctx, cancel := timeoutContext()
	defer cancel()

	var habit model.Habit
	oid, _ := primitive.ObjectIDFromHex(id)
	err := r.col.FindOneAndDelete(ctx, bson.M{"_id": oid}).Decode(&habit)
	if err != nil {
		log.Printf("Fail to delete habit: %v \n", err)
		return false, err
	}
	log.Printf("Deleted_habit(%v) \n", habit)
	return true, nil
}

// buat timeout
func timeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(60)*time.Second)
}