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

func TestCapitalConcat() {
	log.Print("Connection" + " failed") // want "log message should not start with a capital letter"
	log.Print("Error " + "occurred") // want "log message should not start with a capital letter"

	// Valid concatenated strings
	log.Print("connection" + " established")
	log.Print("error " + "handled")
}

func TestOnlyEnglish() {
	log.Fatal("ошибка подключения к базе данных") // want "log message should contain only english symbols"

	// valid log calls
	log.Fatal("failed to connect to database")
}

func TestOnlyEnglishConcat() {
	log.Print("ошибка " + "соединения") // want "log message should contain only english symbols"
	log.Print("connection " + "失败") // want "log message should contain only english symbols"

	// Valid concatenated strings with English only
	log.Print("connection " + "established")
	log.Print("error " + "handled")
}

func TestDigitsAllowed() {
	log.Print("connected to port 8080")
	log.Print("server started on port 3000")
	log.Print("processing 123 items")
	slog.Info("retry attempt 5 of 10")
	log.Print("connection " + "established on port 5432")
}

func TestSpecialSymbols() {
	log.Print("connection failed!") // want "log message should not contain special symbols or emojis"
	log.Print("error: connection refused") // want "log message should not contain special symbols or emojis"
	log.Print("waiting; retry scheduled") // want "log message should not contain special symbols or emojis"
	log.Print("loading...") // want "log message should not contain special symbols or emojis"
	log.Print("success 🎉") // want "log message should not contain special symbols or emojis"

	// Valid messages without special symbols
	log.Print("connection failed")
	log.Print("error occurred")
	log.Print("waiting for retry")
}

func TestSpecialSymbolsConcat() {
	log.Print("connection " + "failed!") // want "log message should not contain special symbols or emojis"
	log.Print("error: " + "refused") // want "log message should not contain special symbols or emojis"

	// Valid concatenated strings
	log.Print("connection " + "failed")
	log.Print("error " + "refused")
}

func TestSensitive() {
	apiKey := "qwerty"
	log.Print("token" + apiKey) // want "log message should not contain sensitive data"

	// valid log calls
	slog.Debug("token validated")
	slog.Info("was" + "validated")
}

func TestSensitiveExtended() {
	log.Print("api_key is invalid") // want "log message should not contain sensitive data"
	log.Print("apikey value") // want "log message should not contain sensitive data"
	log.Print("secret key found") // want "log message should not contain sensitive data"
	log.Print("credential mismatch") // want "log message should not contain sensitive data"

	// Valid log calls without sensitive data
	log.Print("validation successful")
	log.Print("connection established")
}

func TestSensitiveConcat() {
	log.Print("api_key" + " invalid") // want "log message should not contain sensitive data"
	log.Print("secret " + "found") // want "log message should not contain sensitive data"

	// Valid concatenated strings
	log.Print("validation " + "successful")
	log.Print("connection " + "established")
}
