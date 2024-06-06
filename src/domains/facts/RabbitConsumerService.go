package facts

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts/models"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitConsumerService struct {
	factsService *FactsService
	config       services.ConfigurationService
	Connect      *amqp.Connection
	Channel      *amqp.Channel
	Queue        amqp.Queue
	Ctx          context.Context
	Cancel       context.CancelFunc
}

func CreateRabbitConsumerService(factsService *FactsService, config services.ConfigurationService) RabbitConsumerService {
	var rabbitConsumerService = RabbitConsumerService{
		factsService: factsService,
		config:       config,
	}
	var err error

	rabbitConsumerService.Connect, err = amqp.Dial(config.RabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	rabbitConsumerService.Channel, err = rabbitConsumerService.Connect.Channel()
	failOnError(err, "Failed to open a channel")

	rabbitConsumerService.Queue, err = rabbitConsumerService.Channel.QueueDeclare(
		"fact_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	rabbitConsumerService.Ctx, rabbitConsumerService.Cancel = context.WithTimeout(context.Background(), 5*time.Second)

	return rabbitConsumerService
}

func StartRabbitConsumer(c *RabbitConsumerService) {
	c.Consumer()
}

func (c *RabbitConsumerService) Consumer() {

	// принимаем сообщения из очереди

	defer c.Connect.Close()

	defer c.Channel.Close()

	err := c.Channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := c.Channel.Consume(
		c.Queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	// обрабатываем сообщение и конвертируем его в факт

	var forever chan struct{}

	go func() {
		for message := range msgs {
			log.Printf("Received a message: %s \n", message.Body)
			fmt.Println()

			fact, er := UnmarshalFact(message.Body)
			if er != nil {
				fmt.Println(er)
			}

			// передаём факт в сервис для отправки
			er = c.factsService.SaveOne(fact)
			if er != nil {
				fmt.Println(er)
			}

			er = message.Ack(true)
			if er != nil {
				fmt.Println(er)
			}
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func UnmarshalFact(body []byte) (models.Fact, error) {
	var fact models.Fact

	err := json.Unmarshal(body, &fact)

	return fact, err
}
