package godotenv

import (
	"cmp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Environment      string //development,staging,production
	HTTPPort         string
	ReadCoinInterval int64
	DatabaseHost     string
}

func NewEnv() *Env {
	e := &Env{}
	e.Load()
	return e
}

func (e *Env) Load() {
	godotenv.Load(".env") // using .env file is not mandatory
	e.HTTPPort = cmp.Or(os.Getenv("HTTP_PORT"), "8080")
	e.Environment = cmp.Or(os.Getenv("ENVIRONMENT"), "development")
	e.DatabaseHost = os.Getenv("DATABASE_HOST")
	readCoinInterval, err := strconv.ParseInt(
		cmp.Or(os.Getenv("READ_COIN_INTERVAL"), "60"),
		10,
		64)

	if err != nil {
		readCoinInterval = 60
	}
	e.ReadCoinInterval = readCoinInterval
}
