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

    // Conexió amb RabbitMQ
    //(Canviar localhost al nom del contenedor de docker en cas de usar-ho)
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

    // Crear la cua per rebre alertes
    alertQueue, err := ch.QueueDeclare("alertaFumadorMistros", false, false, false, false, nil)
    if err != nil{
        log.Fatal(err)
    }

    // Suscripció a la cua de les alertes
    err = ch.QueueBind(alertQueue.Name, "", "alerta", false, nil)
    if err != nil{
        log.Fatal(err)
    }

    // Demanar mistros
    ch.Publish("", "mistros", false, false, amqp.Publishing{
        Body: []byte("Petició de mistros"),
    })

    // Consumir missatges de les cues
    msgs, err := ch.Consume(mistrosQueue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    alertMsgs, err := ch.Consume(alertQueue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    for{
        select{
        // Diferenciar entre mistros i alertes
        case msg := <- msgs:
            // Generar un temps d'espera aleatori
            waitTime := rand.Intn(4) + 2
            time.Sleep(time.Duration(waitTime) * time.Second)
            fmt.Printf("He agafat el mistro %s. Gràcies! \n. . . \nMe dones un altre mistro?\n", string(msg.Body))
            // Demanar mes mistros
            ch.Publish("", "mistros", false, false, amqp.Publishing{
                Body: []byte("Petició de mistro"),
            })
        case <- alertMsgs:
            fmt.Println("Anem, que ve la policia!")
            return
        }
    }
}
