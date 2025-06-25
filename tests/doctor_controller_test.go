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

	var user models.User
	err := initializers.Db.Where("email = ?", "test@doctor.com").First(&user).Error
	if err != nil {
		user = models.User{
			Name:     "Test Doctor",
			Email:    "test@doctor.com",
			Password: "dummyhash",
			Role:     "doctor",
		}
		if err := initializers.Db.Create(&user).Error; err != nil {
			panic("Failed to insert doctor: " + err.Error())
		}
	}
}

func GetTestJWTTokenDoctor() string {
	var user models.User
	err := initializers.Db.Where("email = ?", "test@doctor.com").First(&user).Error
	if err != nil {
		panic("Doctor not found for token generation")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    "doctor",
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenStr
}
func TestCreatePatientNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.POST("/patient/create/notes/:patient_id", controllers.CreatePatientNotes)
	registerdBy := uint(1)
	patient := models.Patient{
		ID:           1,
		FirstName:    "John",
		LastName:     "Doe",
		DOB:          time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
		Phone:        "1234567890",
		Email:        "johndoe@example.com",
		Address:      "123 Main St",
		RegisteredBy: &registerdBy,
	}
	initializers.Db.FirstOrCreate(&patient)

	body := map[string]string{
		"note": "Patient showed improvement during the week.",
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, "/patient/create/notes/1", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	token := GetTestJWTTokenDoctor()
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: token,
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Create Note Response:", resp.Body.String())
}
func TestEditPatientNotes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.PATCH("/patient/edit/notes/:patient_id", controllers.EditPatientNotes)
	body := map[string]string{
		"note": "updated doctors note",
	}
	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPatch, "/patient/edit/notes/1", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: GetTestJWTTokenDoctor(),
	})
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Edit Patient Notes Response:", resp.Body.String())
}
func TestGetPatientNotes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.GET("/patient/notes/:patient_id", controllers.GetPatientNotes)
	req, err := http.NewRequest(http.MethodGet, "/patient/notes/1", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: GetTestJWTTokenDoctor(),
	})
	resp := httptest.NewRecorder()
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Get Patient Notes Response:", resp.Body.String())
}
