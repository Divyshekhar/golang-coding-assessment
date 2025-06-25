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
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	initializers.LoadEnv()
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
	token := getTestJWTToken()
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
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
	patient_id := "1"
	req, err := http.NewRequest(http.MethodDelete, "/patient/delete/"+patient_id, nil)
	assert.NoError(t, err)
	token := getTestJWTToken()
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
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
	body := map[string]string{
		"phone": "1177771111",
		"name":  "Edited Test",
	}
	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)
	patient_id := "1"
	req, err := http.NewRequest(http.MethodPatch, "/patient/edit/"+patient_id, bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	token := getTestJWTToken()
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Edit Patient Response:", resp.Body.String())

}
