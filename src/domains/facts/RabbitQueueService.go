package facts

import (
	"context"
	"encoding/json"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/facts/models"
	"github.com/IvaCheMih/Testovoe_KPI_Drive/src/domains/services"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitQueueService struct {
	config  services.ConfigurationService
	Connect *amqp.Connection
	Channel *amqp.Channel
	Queue   *amqp.Queue
	Ctx     *context.Context
	Cancel  *context.CancelFunc
}

func CreateRabbitQueueService(config services.ConfigurationService) RabbitQueueService {
	var rabbitQueueService RabbitQueueService
	rabbitQueueService.config = config
	var err error

	//Создаем подключение к RabbitMQ
	rabbitQueueService.Connect, err = amqp.Dial(config.RabbitURL)

	failOnError(err, "Failed to connect to RabbitMQ")

	rabbitQueueService.Channel, err = rabbitQueueService.Connect.Channel()
	failOnError(err, "Failed to open a channel")

	queue, err := rabbitQueueService.Channel.QueueDeclare(
		"fact_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	rabbitQueueService.Queue = &queue

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rabbitQueueService.Ctx, rabbitQueueService.Cancel = &ctx, &cancel

	return rabbitQueueService
}

func (r *RabbitQueueService) RabbitPublish(fact models.Fact) error {

	// конвертируем факт в json
	body, err := MarshalFact(fact)
	failOnError(err, "Failed to marshal fact")

	// отправляем факт в очередь rabbit
	err = r.Channel.PublishWithContext(*r.Ctx,
		"",           // exchange
		r.Queue.Name, // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})

	failOnError(err, "Failed to publish a message")

	// печатаем отправленное сообщение
	log.Printf(" [x] Sent %s \n", body)

	return err
}

func MarshalFact(fact models.Fact) ([]byte, error) {
	body, err := json.Marshal(fact)

	return body, err
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
