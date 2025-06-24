package initializers

import (
	"fmt"

	"github.com/Divyshekhar/golang-coding-assessment/models"
)

func init() {
	LoadEnv()
	ConnectDb()
}

func Migrate() {
	err := Db.AutoMigrate(
		&models.User{},
		&models.Patient{},
		&models.PatientNote{},
	)
	if err != nil {
		panic("Migration failed: " + err.Error())
	}
	fmt.Println("Database migrated successfully")
}
