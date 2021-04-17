package persistence

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitialisePersistence creates a new database connection using config variables.
//
// The available variables are:
// db.username - username
// db.password - password
// db.host - host
// db.database - database
// db.port - port (numeric)
func InitialisePersistence(config *gorm.Config, models ...interface{}) (*gorm.DB, error) {

	var (
		username = viper.GetString("db.username")
		password = viper.GetString("db.password")
		host     = viper.GetString("db.host")
		port     = viper.GetInt("db.port")
		database = viper.GetString("db.database")
	)
	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, username, database, password)
	dialector := postgres.Open(dsn)

	// format dsn based on above values
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}

	// migrate models
	err = db.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}

	return db, nil
}

func InitialisePersistenceFlags() {
	// database flags
	pflag.String("db.host", "localhost", "database hostname")
	pflag.Int("db.port", 3306, "database port")
	pflag.String("db.username", "root", "database username")
	pflag.String("db.password", "", "database password")
	pflag.String("db.database", "default", "database name")
}
