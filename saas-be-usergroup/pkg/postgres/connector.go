package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func Connect() (*sql.DB, error) {
	url := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", viper.GetString("postgres.host"), viper.GetString("postgres.user"), viper.GetString("postgres.password"), viper.GetString("postgres.database"), viper.GetString("postgres.port"))
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
