package controllers

import (
	"net/http"
	"strconv"

	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/gin-gonic/gin"
)

func CreatePatientNotes(ctx *gin.Context) {
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID := userId.(uint)
	role, exist := ctx.Get("role")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	if role != "doctor" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only doctor can edit patient notes"})
		return
	}
	patientIdStr := ctx.Param("patient_id")
	patientId, err := strconv.ParseUint(patientIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}
	result := initializers.Db.Where("id = ?").First(patientId)
	if result.Error != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve record"})
		return
	}
	var body struct{
		Note string `json:"note"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	note := models.PatientNote{
		PatientID: uint(patientId),
		DoctorID: &userID,
		Note: body.Note,
	}
	txn := initializers.Db.Create(&note)
	if txn.Error != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create record",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully created",
		"patient_note": note,
	})
}
