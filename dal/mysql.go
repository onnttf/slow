package dal

import (
	"fmt"
	"slow/config"
	"slow/logger"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	orm_logger "gorm.io/gorm/logger"
)

var (
	once sync.Once
	Slow *gorm.DB
)

// Initialize initializes the database connection and configures the logger.
// It should be called once at the application's start.
func Initialize() error {
	// Construct DSN using provided credentials and configuration
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.Database.User,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Name,
	)
	var err error
	once.Do(func() {
		// Configure custom GORM logger
		newLogger := orm_logger.New(
			logger.Get(),
			orm_logger.Config{
				SlowThreshold:             time.Second,     // Slow SQL threshold
				LogLevel:                  orm_logger.Info, // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for the logger
				ParameterizedQueries:      false,           // Include params in the SQL log
				Colorful:                  false,           // Disable color for log
			},
		)

		// Initialize GORM connection
		Slow, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			logger.Get().Err(err).Msg("connect database fail")
		}
	})
	return err // Return error to be handled by caller
}

func GenerateModel() {
	//g := gen.NewGenerator(gen.Config{
	//	OutPath:           "../query",
	//	FieldWithIndexTag: true,
	//	Mode:              gen.WithDefaultQuery,
	//})
	//
	//g.UseDB(Slow) // reuse your gorm db
	//
	//// Generate basic type-safe DAO API for struct `model.User` following conventions
	//g.GenerateAllTable()
	//// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	//
	//// Generate the code
	//g.Execute()
}
