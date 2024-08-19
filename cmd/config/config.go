package config

import (
	"flag"
	"os"
)

func Flags() (string, string, string) {
	// Определение флагов
	address := flag.String("a", "localhost:8080", "адрес запуска HTTP-сервера")

	addressbonus := flag.String("r", "", "адрес системы расчёта начислений")

	db := flag.String("d", "postgresql://postgres:190603@localhost:5432/postgres?sslmode=disable", "адрес для бд")

	// Парсинг флагов
	flag.Parse()
	if envAddress := os.Getenv("RUN_ADDRESS"); envAddress != "" {
		*address = envAddress
	}

	if envaddressbonus := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envaddressbonus != "" {
		*addressbonus = envaddressbonus
	}

	if envDB := os.Getenv("DATABASE_URI"); envDB != "" {
		*db = envDB
	}

	return *address, *addressbonus, *db
}
