package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Divyshekhar/golang-coding-assessment/controllers"
	"github.com/Divyshekhar/golang-coding-assessment/initializers"
	"github.com/Divyshekhar/golang-coding-assessment/middleware"
	"github.com/Divyshekhar/golang-coding-assessment/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDb()

	// Lookup user by ID only
	var user models.User
	result := initializers.Db.FirstOrCreate(&user, models.User{ID: 1})

	// If user created OR mismatch on critical fields, update them
	if result.RowsAffected == 1 || user.Email != "test@receptionist.com" || user.Role != "receptionist" {
		user.Name = "Test Receptionist"
		user.Email = "test@receptionist.com"
		user.Password = "dummyhash" // You may replace this with hashed password if needed
		user.Role = "receptionist"
		initializers.Db.Save(&user)
	}
}

func getTestJWTToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": 1,
		"role":    "receptionist",
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenStr
}

func TestRegisterPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(middleware.RequireAuth())
	router.POST("/patient/create", controllers.RegisterPatient)

	body := map[string]string{
		"firstName": "Test",
		"lastName":  "Name",
		"dob":       "2025-06-25",
		"gender":    "male",
		"phone":     "0123456789",
		"email":     "test@test.com",
		"address":   "test address",
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, "/patient/create", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: getTestJWTToken(),
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Response:", resp.Body.String())
}

func TestDeletePatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(middleware.RequireAuth())
	router.DELETE("/patient/delete/:patient_id", controllers.DeletePatient)

	patientID := "1"
	req, err := http.NewRequest(http.MethodDelete, "/patient/delete/"+patientID, nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: getTestJWTToken(),
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Delete Patient Response:", resp.Body.String())
}

func TestEditPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.PATCH("/patient/edit/:patient_id", controllers.EditPatient)
	registeredBy := uint(1)
	initializers.Db.Create(&models.Patient{
		ID:           1,
		FirstName:    "Test",
		LastName:     "Name",
		DOB:          time.Date(2025, 6, 25, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
		Phone:        "0123456789",
		Email:        "test@test.com",
		Address:      "test address",
		RegisteredBy: &registeredBy,
	})

	body := map[string]string{
		"phone": "1177771111",
		"name":  "Edited Test",
	}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	patientID := "1"
	req, err := http.NewRequest(http.MethodPatch, "/patient/edit/"+patientID, bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: getTestJWTToken(),
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Edit Patient Response:", resp.Body.String())
}
func TestGetPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.GET("/patient/all", controllers.GetPatient)
	req, err := http.NewRequest(http.MethodGet, "/patient/all", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: getTestJWTToken(),
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Get Patient Response:", resp.Body.String())

}
