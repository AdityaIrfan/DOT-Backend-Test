package postgres

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	url := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", viper.GetString("postgres.host"), viper.GetString("postgres.user"), viper.GetString("postgres.password"), viper.GetString("postgres.database"), viper.GetString("postgres.port"))
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		PrepareStmt:          true,
		DisableAutomaticPing: false,
	})
	fmt.Println(url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
