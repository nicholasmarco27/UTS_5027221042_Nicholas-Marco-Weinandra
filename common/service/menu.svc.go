package service

import (
	"context"
	"log"
	"github.com/nicholasmarco27/UTS_5027221042_Nicholas-Marco-Weinandra/common/genproto/menulist"
	"github.com/nicholasmarco27/UTS_5027221042_Nicholas-Marco-Weinandra/common/model"
	"github.com/nicholasmarco27/UTS_5027221042_Nicholas-Marco-Weinandra/common/repository"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type MenuService struct {
	repo *repository.MenuRepository
}

func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{
		repo: repo,
	}
}

// buat menu
func (s *MenuService) CreateMenu(ctx context.Context, tm *menulist.Menu) (*menulist.Menu, error) {
	log.Printf("CreateMenu(%v) \n", tm)

	newMenu := &model.Menu{
		Title:       tm.Title,
		Description: tm.Description,
	}

	menu, err := s.repo.Save(newMenu)
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return s.toMenu(&menu), nil
}

// dapatin semua menu
func (s *MenuService) ListMenus(ctx context.Context, e *empty.Empty) (*menulist.MenuList, error) {
	log.Printf("ListMenus() \n")

	var totas []*menulist.Menu
	Menus, err := s.repo.FindAll()
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	for _, u := range Menus {
		totas = append(totas, s.toMenu(&u))
	}

	MenuList := &menulist.MenuList{
		List: totas,
	}

	return MenuList, nil
}

// update menu
func (s *MenuService) UpdateMenu(ctx context.Context, tm *menulist.Menu) (*menulist.Menu, error) {
	log.Printf("UpdateMenu(%v) \n", tm)

	if tm.Id == "" {
		return nil, status.Error(codes.FailedPrecondition, "UpdateMenu must provide menuID")
	}

	menuID, err := primitive.ObjectIDFromHex(tm.Id)
	if err != nil {
		log.Printf("Invalid MenuID(%s) \n", tm.Id)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	updateMenu := &model.Menu{
		ID:          menuID,
		Title:       tm.Title,
		Description: tm.Description,
	}

	menu, err := s.repo.Update(updateMenu)
	if err != nil {
		log.Printf("Fail UpdateMenu %v \n", err)
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	return s.toMenu(&menu), nil
}

// delete menu
func (s *MenuService) DeleteMenu(ctx context.Context, id *wrappers.StringValue) (*wrappers.BoolValue, error) {
	log.Printf("DeleteMenu(%s) \n", id.GetValue())

	deleted, err := s.repo.Delete(id.GetValue())
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}

	return &wrapperspb.BoolValue{Value: deleted}, nil
}

// map menu ke toMenu
func (s *MenuService) toMenu(u *model.Menu) *menulist.Menu {
	tota := &menulist.Menu{
		Id:          u.ID.Hex(),
		Title:       u.Title,
		Description: u.Description,
	}
	return tota
}
