package services

import (
	"context"
	"errors"
	"example/apies/models"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ColorServiceImp struct {
	colorcollection *mongo.Collection
	ctx             context.Context
}

func NewColorService(colorcollection *mongo.Collection, ctx context.Context) ColorService {
	return &ColorServiceImp{
		colorcollection: colorcollection,
		ctx:             ctx,
	}
}

func (u *ColorServiceImp) CreateColor(color *models.Color) error {
	_, err := u.colorcollection.InsertOne(u.ctx, color)
	return err
}

//	func (u *ColorServiceImp) GetColor(name *string) (*models.Color, error) {
//		var color *models.Color
//		query := bson.D{bson.E{Key: "color_name", Value: name}}
//		err := u.colorcollection.FindOne(u.ctx, query).Decode(&color)
//		return color, err
//	}
func (u *ColorServiceImp) GetAll() ([]*models.Color, error) {

	fmt.Println("getall called")
	var colors []*models.Color
	cursor, err := u.colorcollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var color models.Color
		err := cursor.Decode(&color)
		if err != nil {
			return nil, err
		}
		colors = append(colors, &color)

	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(u.ctx)
	if len(colors) == 0 {
		return nil, errors.New("not found")
	}
	return nil, nil
}
