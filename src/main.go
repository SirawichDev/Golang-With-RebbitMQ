package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	go client()
	go server()

	var a string
	fmt.Scanln(&a)
}
func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, //queue string
		"", //consumer string
		true, //AutoAck bool
		false,//exclusive bool, 
		false,//noLocal bool, 
		true,//noWait bool, 
		nil)//args amqp.Table)
	failOnError(err,"Failed to register a cosumer")
	for msg := range msgs{
		log.Printf("Received message with message: %s",msg.Body)
	}
}
func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hellow Me "),
	}
	for {
		ch.Publish("", q.Name, false, false, msg)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open Channel")
	q, err := ch.QueueDeclare(
		"Hellow",
		false,
		false,
		false,
		false,
		nil)
	failOnError(err, "Failed to declare a queue")
	return conn, ch, &q
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalln(fmt.Sprintf("%s: %s", msg, err))
	}
}
