package a

import (
	"log"
	"log/slog"
)

func TestCapital() {
	log.Print("Capital letter") // want "log message should not start with a capital letter"
	slog.Error("Failed to connect to database") // want "log message should not start with a capital letter"

	// Valid log calls
	log.Println("valid log message")
	log.Print("formatted message value")
}

func TestOnlyEnglish() {
	log.Fatal("ошибка подключения к базе данных") // want "log message should contain only english symbols"

	// valid log calls
	log.Fatal("failed to connect to database")
}

func TestSensitive() {
	apiKey := "qwerty"
	log.Print("token" + apiKey) // want "log message should not contain sensitive data"

	// valid log calls
	slog.Debug("token validated")
	slog.Info("was" + "validated")
}
