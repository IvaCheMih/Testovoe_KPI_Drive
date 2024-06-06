package main

import (
	_ "github.com/IvaCheMih/Testovoe_KPI_Drive/src/docs"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/services"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"log"
)

var rabbitQueueService facts.RabbitQueueService
var rabbitConsumerService facts.RabbitConsumerService
var factsHandlers facts.FactHandlers

func Init() {
	// сохраняем переменные окружения (URL Rabbit и внешних api):

	configurationService, err := services.CreateConfigurationService()
	if err != nil {
		panic(err)
	}

	// создаём структуры сервисов и хэндлеров:

	rabbitQueueService = facts.CreateRabbitQueueService(configurationService)

	factsServices := facts.CreateFactsService(&rabbitQueueService, configurationService)

	rabbitConsumerService = facts.CreateRabbitConsumerService(&factsServices, configurationService)

	factsHandlers = facts.CreateFactHandlers(&factsServices)

	// запускаем сервис, который ждёт сообщений из rabbit:

	go facts.StartRabbitConsumer(&rabbitConsumerService)
}

// @title 						Fiber Swagger Example API
// @version 					2.0
// @description 				This is a sample server.
// @termsOfService 				http://swagger.io/terms/

// @contact.name				API Support
// @contact.url 				http://www.swagger.io/support
// @contact.email				support@swagger.io

// @license.name 				Apache 2.0
// @license.url 				http://www.apache.org/licenses/LICENSE-2.0.html

// @host 						localhost:8080
// @BasePath 					/
// @schemes 					http
//
//	@securityDefinitions.apiKey JWT
//	@in                         header
//	@name                       Authorization
//	@description                JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {

	// запускаем приложение:

	Init()

	server := fiber.New()

	// хэндлер для swagger:

	server.Get("/swagger/*", swagger.HandlerDefault)

	// два основых хэндлера:

	server.Post("/CreateFact", factsHandlers.CreateFact)

	server.Post("/GetFact", factsHandlers.GetFact)

	if err := server.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
