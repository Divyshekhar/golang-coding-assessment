package routes

import (
	"github.com/Divyshekhar/golang-coding-assessment/controllers"
	"github.com/Divyshekhar/golang-coding-assessment/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(ctx *gin.Engine) {
	patientGroup := ctx.Group("/patient")
	{
		patientGroup.POST("/create", middleware.RequireAuth(), controllers.RegisterPatient)
		patientGroup.PATCH("/edit/:patient_id", middleware.RequireAuth(), controllers.EditPatient)
		patientGroup.DELETE("/delete/:patient_id", middleware.RequireAuth(), controllers.DeletePatient)
		patientGroup.GET("/all", middleware.RequireAuth(), controllers.GetPatient)
		patientGroup.POST("/create/notes/:patient_id", middleware.RequireAuth(), controllers.CreatePatientNotes)
		patientGroup.PATCH("/edit/notes/:patient_id", middleware.RequireAuth(), controllers.EditPatientNotes)
		patientGroup.GET("/notes/:patient_id", middleware.RequireAuth(), controllers.GetPatientNotes)
	}
}
