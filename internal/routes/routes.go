package routes

import (
	"locntp-user-counter/internal/controllers"
	"locntp-user-counter/internal/repositories"
	"locntp-user-counter/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

type ApplicationContainer struct {
	// Repositories
	UserRepository *repositories.UserRepository

	// Services
	UserService *services.UserService

	// Controllers
	UserController *controllers.UserController
}

func (app *ApplicationContainer) SetupRepositories(db *gorm.DB, cache *redis.Client) {
	app.UserRepository = repositories.NewUserRepository(db, cache)
}

func (app *ApplicationContainer) SetupServices() {
	app.UserService = services.NewUserService(app.UserRepository)
}

func (app *ApplicationContainer) SetupControllers() {
	app.UserController = controllers.NewUserController(app.UserService)
}

var appContainer *ApplicationContainer

// NewApplicationContainer initializes the application container with repositories, services, and controllers
func GetApplicationContainer(db *gorm.DB, cache *redis.Client) *ApplicationContainer {
	// Initialize the application container only once
	if appContainer != nil {
		return appContainer
	}

	appContainer := &ApplicationContainer{}
	appContainer.SetupRepositories(db, cache)
	appContainer.SetupServices()
	appContainer.SetupControllers()
	return appContainer
}

// SetupRoutes sets up the routes and the corresponding handlers
func SetupRouter(db *gorm.DB, cache *redis.Client, gin *gin.Engine) *gin.Engine {
	// Swagger routes
	//gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRouter := gin.Group("")

	// Setup the v1 routes
	v1 := publicRouter.Group("/api/v1")

	// Public routes
	{
		// Setup the user routes
		NewUserRouters(db, cache, v1)

	}

	return gin
}

// NewUserRouters sets up the routes and the corresponding handlers
func NewUserRouters(db *gorm.DB, cache *redis.Client, group *gin.RouterGroup) {
	// Initialize the application container
	appContainer := GetApplicationContainer(db, cache)
	uc := appContainer.UserController

	// Create a new group for the user routes
	userGroup := group.Group("/users")
	publicGroup := userGroup.Group("")
	// Public Routes
	{
		publicGroup.POST("", uc.CreateUser)
		publicGroup.GET("/:id/count", uc.GetUser)
		publicGroup.PUT("/:id/increment", uc.IncrementCounter)
		publicGroup.DELETE("/:id", uc.DeleteUser)
	}
}
