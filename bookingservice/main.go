package main

import (
	"bookingservice/configs"
	"bookingservice/routes"
	"bookingservice/services"

	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	//kafkaProducer, err := service.NewKafkaProducer()
	//if err != nil {
	//	log.Fatalf("Error creating Kafka producer: %v", err)
	//}
	//
	//kafkaConsumer, err := service.NewKafkaConsumer()
	//if err != nil {
	//	log.Fatalf("Error creating Kafka consumer: %v", err)
	//}
	//
	//kafkaConsumer.ConsumeMessages()

	r := gin.Default()

	db := configs.InitDB()

	rdb := configs.ConnectRedis()

	ts := services.NewTicketService(db, rdb)

	routes.RegisterRouters(r, ts)

	if err := r.Run(":8084"); err != nil {
		log.Fatal("Failed to start the server: ", err)
	}

}
