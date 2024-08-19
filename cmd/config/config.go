package config

import (
	"flag"
	"os"
)

type Cfg struct {
	ServerAddress  string
	DBAddress      *string
	AccrualAddress *string
}

var config Cfg

func init() {
	config.ServerAddress = *flag.String("a", "localhost:8080", "server address")
	config.DBAddress = flag.String("d", "", "data base connection address")
	config.AccrualAddress = flag.String("r", "", "accrual system server address")
}
func Flags() (string, string, string) {
	// Определение флагов

	// Парсинг флагов
	flag.Parse()

	dbAddressEnv := os.Getenv("DATABASE_URI")

	if dbAddressEnv != "" {
		config.DBAddress = &dbAddressEnv
	}
	serverAddressEnv := os.Getenv("RUN_ADDRESS")
	if serverAddressEnv != "" {
		config.ServerAddress = serverAddressEnv
	}
	accrualEnv := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if accrualEnv != "" {
		config.AccrualAddress = &accrualEnv
	}
	if *config.DBAddress == "" || *config.AccrualAddress == "" || config.ServerAddress == "" {
		panic("invalid config")
	}
	return config.ServerAddress, *config.AccrualAddress, *config.DBAddress

}
