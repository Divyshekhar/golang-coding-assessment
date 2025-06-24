package controllers

import (
	"net/http"
	"time"

	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/gin-gonic/gin"
)

func RegisterPatient(ctx *gin.Context) {
	var body struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Dob       string `json:"dob"`
		Gender    string `json:"gender"`
		Phone     string `json:"phone"`
		Email     string `json:"email"`
		Address   string `json:"address"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid inputs"})
		return
	}
	parsedDob, err := time.Parse("2006-01-02", body.Dob)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	patient := models.Patient{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		DOB:       parsedDob,
		Gender:    body.Gender,
		Phone:     body.Phone,
		Email:     body.Email,
		Address:   body.Address,
	}
	result := initializers.Db.Create(&patient)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create record"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Patient Registered",
		"patient": patient,
	})
}
