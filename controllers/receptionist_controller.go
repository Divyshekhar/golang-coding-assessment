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
	userId, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userID := userId.(uint)
	role, exist := ctx.Get("role")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated with correct role"})
		return
	}
	roleVal := role.(string)
	if roleVal != "receptionist" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only receptionist can register the patient"})
		return
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
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		DOB:          parsedDob,
		Gender:       body.Gender,
		Phone:        body.Phone,
		Email:        body.Email,
		Address:      body.Address,
		RegisteredBy: &userID,
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

func EditPatient(ctx *gin.Context) {
	id := ctx.Param("patient_id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no patient id found"})
		return
	}
	_, exist := ctx.Get("user_id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}


	role, exist := ctx.Get("role")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no role specified"})
		return
	}
	roleVal := role.(string)
	if roleVal != "receptionist" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "only receptionist are edit the patient information"})
	}

	var patient models.Patient
	result := initializers.Db.Where("id = ?", id).First(&patient)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find the patient with the given id"})
		return
	}

	var updates map[string]interface{}
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}


	allowed := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"dob":        true,
		"gender":     true,
		"phone":      true,
		"email":      true,
		"address":    true,
	}

	safeUpdates := map[string]interface{}{}
	for key, value := range updates {
		if allowed[key] {
			if key == "dob" {

				if strDob, ok := value.(string); ok {
					parsedDob, err := time.Parse("2006-01-02", strDob)
					if err != nil {
						ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DOB format. Use YYYY-MM-DD"})
						return
					}
					safeUpdates["dob"] = parsedDob
				}
			} else {
				safeUpdates[key] = value
			}
		}
	}

	if len(safeUpdates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields to update"})
		return
	}

	if err := initializers.Db.Model(&patient).Updates(safeUpdates).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update patient"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Patient updated successfully",
		"patient": patient,
	})

}
