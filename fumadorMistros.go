package main

import (
    "fmt"
    "log"
    "time"
    "math/rand"
    "github.com/streadway/amqp"
)

func main() {
    fmt.Println("Sóm fumador. Tenc tabac pero me fan falta mistros")

    conn, err := amqp.Dial("amqp://guest:guest@RabbitMQ:5672/")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    // Crear la cua per rebre mistros
    mistrosQueue, err := ch.QueueDeclare("fumadorMistros", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Demanar mistros
    ch.Publish("", "mistros", false, false, amqp.Publishing{
        Body: []byte("Petició de mistros"),
    })

    // Consumir missatges de la cua
    msgs, err := ch.Consume(mistrosQueue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Esperar resposta i gestionar el missatge
    for msg := range msgs {
        // Generar un temps d'espera aleatori
        waitTime := rand.Intn(4) + 2
        time.Sleep(time.Duration(waitTime) * time.Second)
        fmt.Printf("He agafat el mistro %s. Gràcies! \n. . . \nMe dones un altre mistro?\n", string(msg.Body))
        ch.Publish("", "mistros", false, false, amqp.Publishing{
            Body: []byte("Petició de mistro"),
        })
    }
}
