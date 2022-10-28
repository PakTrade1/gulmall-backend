package controllers

import (
	"example/apies/models"
	"example/apies/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ColorController struct {
	ColorService services.ColorService
}

func New(colorservice services.ColorService) ColorController {
	return ColorController{
		ColorService: colorservice,
	}
}

func (uc *ColorController) CreateColor(ctx *gin.Context) {
	var color models.Color
	if err := ctx.ShouldBindJSON(&color); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.ColorService.CreateColor(&color)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (uc *ColorController) GetAll(ctx *gin.Context) {
	color, err := uc.ColorService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, color)
}

func (uc *ColorController) RegisterColorRouts(rg *gin.RouterGroup) {
	colorroute := rg.Group("/color")
	colorroute.POST("/create", uc.CreateColor)
	//colorroute.GET("/get/:name", uc.GetColor)
	colorroute.GET("/getall", uc.GetAll)

}
