package app

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	MySQL *gorm.DB
	Redis *redis.Client
}

func NewApplication() Application {
	app := &Application{}

	mysql, err := NewMySqlDatabase()
	if err != nil {
		return Application{}
	}
	app.MySQL = mysql

	redisDB, err := NewRedisInMemoryDatabase()
	if err != nil {
		return Application{}
	}
	app.Redis = redisDB

	return *app
}

func (app *Application) Close() {
	err := CloseMySqlDatabase(app.MySQL)
	if err != nil {
		fmt.Println("Cannot close MySQL database connection")
	}
	err = CloseRedisInMemoryDatabase(app.Redis)
	if err != nil {
		fmt.Println("Cannot close Redis database connection")
	}
	fmt.Println("Closing database connections")
}

func StartServer() {
	application := NewApplication()
	defer application.Close()

	if err := CreateIndexes(application.MySQL); err != nil {
		panic(err)
	}

	ginServer := NewServer()

	// Setup middleware in order of execution
	ginServer.CorsMiddleware()
	ginServer.SecurityMiddleware()
	ginServer.RateLimitMiddleware()
	ginServer.RouteHandler(application.MySQL, application.Redis)

	// Start the server
	ginServer.StartServer()
}
