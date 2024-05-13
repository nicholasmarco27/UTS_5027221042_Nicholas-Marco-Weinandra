package service

import (
	"context"
	"log"
	"github.com/GabriellaErlinda/UTS_5027221018_Gabriella-Erlinda/common/genproto/habittracker"
	"github.com/GabriellaErlinda/UTS_5027221018_Gabriella-Erlinda/common/model"
	"github.com/GabriellaErlinda/UTS_5027221018_Gabriella-Erlinda/common/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type HabitService struct {
	repo *repository.HabitRepository
}

func NewHabitService(repo *repository.HabitRepository) *HabitService {
	return &HabitService{
		repo: repo,
	}
}

// buat habit
func (s *HabitService) CreateHabit(ctx context.Context, tm *habittracker.Habit) (*habittracker.Habit, error) {
	log.Printf("CreateHabit(%v) \n", tm)

	newHabit := &model.Habit{
		Title:       tm.Title,
		Description: tm.Description,
	}

	habit, err := s.repo.Save(newHabit)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return s.toHabit(&habit), nil
}

// dapatin semua habit
func (s *HabitService) ListHabits(ctx context.Context, e *empty.Empty) (*habittracker.HabitList, error) {
	log.Printf("ListHabits() \n")

	var totas []*habittracker.Habit
	Habits, err := s.repo.FindAll()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	for _, u := range Habits {
		totas = append(totas, s.toHabit(&u))
	}

	HabitList := &habittracker.HabitList{
		List: totas,
	}

	return HabitList, nil
}

// update habit
func (s *HabitService) UpdateHabit(ctx context.Context, tm *habittracker.Habit) (*habittracker.Habit, error) {
	log.Printf("UpdateHabit(%v) \n", tm)

	if tm.Id == "" {
		return nil, status.Error(codes.FailedPrecondition, "UpdateHabit must provide habitID")
	}

	habitID, err := primitive.ObjectIDFromHex(tm.Id)
	if err != nil {
		log.Printf("Invalid HabitID(%s) \n", tm.Id)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	updateHabit := &model.Habit{
		ID:          habitID,
		Title:       tm.Title,
		Description: tm.Description,
	}

	habit, err := s.repo.Update(updateHabit)
	if err != nil {
		log.Printf("Fail UpdateHabit %v \n", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return s.toHabit(&habit), nil
}

// delete habit
func (s *HabitService) DeleteHabit(ctx context.Context, id *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("DeleteHabit(%s) \n", id.GetValue())

	deleted, err := s.repo.Delete(id.GetValue())
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &wrapperspb.BoolValue{Value: deleted}, nil
}

// map habit ke toHabit
func (s *HabitService) toHabit(u *model.Habit) *habittracker.Habit {
	tota := &habittracker.Habit{
		Id:          u.ID.Hex(),
		Title:       u.Title,
		Description: u.Description,
	}
	return tota
}
