package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables de entorno del archivo .env

func LoadEnv() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}
	log.Println("Archivo .env cargado correctamente")

}

// DBConfig es una estructura que contiene los datos de configuración de la base de datos

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDBConfig obtiene la configuración de la base de datos

func GetDBConfig() DBConfig {

	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

//  GetJWTSecret obtiene la clave secreta para firmar el token JWT

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("La clave secreta de JWT no está configurada")
	}
	return secret
}
