package main

import (
    "fmt"
    "log"
    "time"
    "math/rand"
    "github.com/streadway/amqp"
)

func main() {
    fmt.Println("Sóm fumador. Tenc mistros però me fa falta tabac")

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

    // Crear la cua per rebre tabac
    tabacQueue, err := ch.QueueDeclare("fumadorTabac", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Demanar tabac
    ch.Publish("", "tabac", false, false, amqp.Publishing{
        Body: []byte("Petició de tabac"),
    })

    // Consumir missatges de la cua
    msgs, err := ch.Consume(tabacQueue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Esperar resposta i gestionar el missatge
    for msg := range msgs {
        // Generar un temps d'espera aleatori
        waitTime := rand.Intn(4) + 2
        time.Sleep(time.Duration(waitTime) * time.Second)
        fmt.Printf("He agafat el tabac %s. Gràcies! \n. . . \nMe dones mes tabac?\n", string(msg.Body))
        ch.Publish("", "tabac", false, false, amqp.Publishing{
            Body: []byte("Petició de tabac"),
        })
    }
}