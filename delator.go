package main

import (
    "fmt"
    "log"
    amqp "github.com/streadway/amqp"
)

func main() {
    fmt.Println("No som fumador. ALERTA! Que ve la policia")

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

    // Creem un Fanout exchange per comunicar-se amb els fumadors
    err = ch.ExchangeDeclare(
        "alerta",   // nom de l'exchange
        "fanout",   // tipus de l'exchange
        true,       // durable
        false,      // auto-deletes
        false,      // internal
        false,      // no-wait
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Enviar un missatge d'alerta
    ch.Publish("alerta", "", false, false, amqp.Publishing{
        Body: []byte("Anem, que ve la policia!"),
    })
}
