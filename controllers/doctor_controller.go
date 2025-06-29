package controllers

import (
	"net/http"
	"strconv"

	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/Divyshekhar/golang-coding-assessment/utils"
	"github.com/gin-gonic/gin"
)

func CreatePatientNotes(ctx *gin.Context) {
	role, exist := ctx.Get("role")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	roleVal := role.(string)
	if roleVal != "doctor" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only doctor can edit patient notes"})
		return
	}
	user, ok := utils.GetUserAndCheckRole(ctx, "doctor")
	if !ok {
		return
	}

	patientIdStr := ctx.Param("patient_id")
	patientId, err := strconv.ParseUint(patientIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}
	var patient models.Patient
	result := initializers.Db.First(&patient, patientId)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve record"})
		return
	}
	var body struct {
		Note string `json:"note"`
	}
	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	note := models.PatientNote{
		PatientID: uint(patientId),
		DoctorID:  &user.ID,
		Note:      body.Note,
	}
	txn := initializers.Db.Create(&note)
	if txn.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create record",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Successfully created",
		"patient_note": note,
	})
}

func EditPatientNotes(ctx *gin.Context) {
	role, exist := ctx.Get("role")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	roleVal := role.(string)
	if roleVal != "doctor" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only doctors are allowed to edit the patient's note"})
		return
	}
	user, ok := utils.GetUserAndCheckRole(ctx, "doctor")
	if !ok {
		return
	}
	patientIdStr := ctx.Param("patient_id")
	patientId, err := strconv.ParseUint(patientIdStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no patient id found"})
		return
	}
	var note models.PatientNote
	if err := initializers.Db.
		Where("patient_id = ? AND doctor_id = ?", patientId, user.ID).
		Order("created_at DESC").
		First(&note).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "patient note not found"})
		return
	}

	if note.DoctorID != nil && *note.DoctorID != user.ID {

		ctx.JSON(http.StatusForbidden, gin.H{"error": "unauthorized to edit this note"})
		return
	}
	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	safeUpdates := map[string]interface{}{}
	if val, ok := updates["note"]; ok {
		safeUpdates["note"] = val
	}

	if len(safeUpdates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no valid fields to update"})
		return
	}
	if err := initializers.Db.Model(&note).Updates(safeUpdates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update note"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Patient note updated successfully",
		"note":    note,
	})
}

func GetPatientNotes(ctx *gin.Context) {
	role, exist := ctx.Get("role")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}
	roleVal := role.(string)
	if roleVal != "doctor" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only doctors can retrieve the patient notes"})
		return
	}
	user, ok := utils.GetUserAndCheckRole(ctx, "doctor")
	if !ok {
		return
	}
	patientIdUint, err := strconv.ParseUint(ctx.Param("patient_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid patient id"})
		return
	}
	var patient models.PatientNote
	txn := initializers.Db.Where("patient_id = ?", patientIdUint).
		Where("doctor_id = ?", user.ID).
		First(&patient)
	if txn.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve the patient notes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Fetched successful",
		"patient data": patient,
	})
}
