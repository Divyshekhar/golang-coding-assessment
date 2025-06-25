package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var doctorUser models.User
var patient models.Patient
var createdNoteID uint

func init() {
	initializers.LoadEnv()
	initializers.ConnectDb()

	// Create or fetch doctor
	doctorUser = models.User{
		Name:     "Test Doctor",
		Email:    "test@doctor.com",
		Password: "dummyhash",
		Role:     "doctor",
	}
	initializers.Db.Where("email = ?", doctorUser.Email).FirstOrCreate(&doctorUser)

	// Create or fetch patient
	patient = models.Patient{
		ID:           1001,
		FirstName:    "John",
		LastName:     "Doe",
		DOB:          time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		Gender:       "male",
		Phone:        "1234567890",
		Email:        "johndoe@example.com",
		Address:      "123 Main St",
		RegisteredBy: &doctorUser.ID,
	}
	initializers.Db.Where("id = ?", patient.ID).FirstOrCreate(&patient)

	// Clean old notes
	initializers.Db.Where("patient_id = ?", patient.ID).Delete(&models.PatientNote{})
}

func GetTestJWTTokenDoctor() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": doctorUser.ID,
		"role":    doctorUser.Role,
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

	body := map[string]string{
		"note": "Patient showed improvement during the week.",
	}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, "/patient/create/notes/1001", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: GetTestJWTTokenDoctor(),
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Log raw response
	t.Log("Raw response:", resp.Body.String())

	// Parse response JSON
	var response map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	pnRaw, exists := response["patient_note"]
	assert.True(t, exists, "Expected key 'patient_note' in response")

	pnMap, ok := pnRaw.(map[string]interface{})
	assert.True(t, ok, "Expected 'patient_note' to be an object")

	idRaw, ok := pnMap["ID"]
	assert.True(t, ok, "Expected 'ID' field in patient_note")

	floatID, ok := idRaw.(float64)
	assert.True(t, ok, "Expected 'ID' to be a number")

	createdNoteID = uint(floatID)
	t.Logf("Successfully created patient note with ID: %d", createdNoteID)
}

func TestEditPatientNotes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.PATCH("/patient/edit/notes/:patient_id", controllers.EditPatientNotes)

	// Step 1: Create a patient note to update
	note := models.PatientNote{
		PatientID: 1001,
		DoctorID:  func(u uint) *uint { return &u }(3),
		Note:      "Initial note",
		CreatedAt: time.Now(),
	}
	if err := initializers.Db.Create(&note).Error; err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	body := map[string]string{
		"note": "updated doctor note",
	}
	jsonBody, _ := json.Marshal(body)

	// Step 2: Create the PATCH request using the correct ID
	url := fmt.Sprintf("/patient/edit/notes/%d", note.ID)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: GetTestJWTTokenDoctor(),
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	t.Log("Edit Patient Notes Response:", resp.Body.String())
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGetPatientNotes(t *testing.T) {
	TestCreatePatientNote(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(middleware.RequireAuth())
	router.GET("/patient/notes/:patient_id", controllers.GetPatientNotes)

	req, err := http.NewRequest(http.MethodGet, "/patient/notes/1001", nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{
		Name:  "jwt_token",
		Value: GetTestJWTTokenDoctor(),
	})

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
	t.Log("Get Patient Notes Response:", resp.Body.String())
}
