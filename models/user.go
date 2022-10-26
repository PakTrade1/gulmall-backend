package models

type Color struct {
	Id     string `json:"_id" bson:"color_id"`
	CssHex string `json:"cssHex" bson:"color_csshex"`
	Name   Name   `json:"name" bson:"color_name"`
}
type Name struct {
	En string `json:"en" bson:"color_en"`
	Ar string `json:"ar" bson:"color_ar"`
}
