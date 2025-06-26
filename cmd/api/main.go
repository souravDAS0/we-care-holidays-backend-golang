package main

import (
	"log"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/bootstrap"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/server"

	// Import swagger docs - this is important!
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/docs"

	// Import packages containing types referenced in Swagger comments
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/dto"
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/presentation/http/dto"
)

//	@title			WeCare Holidays API
//	@version		1.0
//	@description	API Server for WeCare Holidays applications
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.wecareholidays.com/support
//	@contact.email	support@wecareholidays.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token

// Standard response model documentation
//	@Success	200	{object}	models.SwaggerStandardResponse{data=interface{}}
//	@Success	201	{object}	models.SwaggerStandardResponse{data=interface{}}
//	@Success	204	{object}	models.SwaggerStandardResponse
//	@Failure	400	{object}	models.SwaggerErrorResponse
//	@Failure	401	{object}	models.SwaggerErrorResponse
//	@Failure	403	{object}	models.SwaggerErrorResponse
//	@Failure	404	{object}	models.SwaggerErrorResponse
//	@Failure	409	{object}	models.SwaggerErrorResponse
//	@Failure	422	{object}	models.SwaggerErrorResponse
//	@Failure	429	{object}	models.SwaggerErrorResponse
//	@Failure	500	{object}	models.SwaggerErrorResponse
//	@Failure	503	{object}	models.SwaggerErrorResponse

func main() {

	appContainer := bootstrap.Bootstrap()
	// Load environment configs
	configs.LoadConfig()

	// Initialize logger
	logger.InitLogger(configs.AppConfig.Env)

	log.Println("Starting WeCare Holidays API server...")
	log.Println("Swagger UI should be available at /swagger/index.html")

	// Start the server
	if err := server.StartServer(appContainer); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
