package utils

import (
	"ecommerce-api/config"
	"ecommerce-api/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {

	var dsn string

	if os.Getenv("DATABASE_URL") != "" {
		dsn = os.Getenv("DATABASE_URL")
	} else {
		dbConfig := config.GetDBConfig()
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s ",
			dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port,
		)
	}

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos", err)
	}

	log.Println("Conexión a la base de datos exitosa")

}

// AutoMigrateDB crea las tablas en la base de datos

func AutoMigrateDB() {

	if os.Getenv("DB_AUTOMIGRATE") == "true" {
		DB.AutoMigrate(&models.User{})
		log.Println("Tablas creadas exitosamente")
	} else {
		log.Println("Tablas no creadas")
	}
}
