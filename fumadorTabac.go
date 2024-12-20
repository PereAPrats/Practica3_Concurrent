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

    // Crear la cua per rebre tabac
    tabacQueue, err := ch.QueueDeclare("fumadorTabac", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Crear la cua per rebre alertes
    alertQueue, err := ch.QueueDeclare("alertaFumadorTabac", false, false, false, false, nil)
    if err != nil{
        log.Fatal(err)
    }

    // Suscripció a la cua de les alertes
    err = ch.QueueBind(alertQueue.Name, "", "alerta", false, nil)
    if err != nil{
        log.Fatal(err)
    }

    // Demanar tabac
    ch.Publish("", "tabac", false, false, amqp.Publishing{
        Body: []byte("Petició de tabac"),
    })

    // Consumir missatges de les cues
    msgs, err := ch.Consume(tabacQueue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    alertMsgs, err := ch.Consume(alertQueue.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    for{
        select{
        // Diferenciar entre tabac i aletres
        case msg := <- msgs:
            // Generar un temps d'espera aleatori
            waitTime := rand.Intn(4) + 2
            time.Sleep(time.Duration(waitTime) * time.Second)
            fmt.Printf("He agafat el tabac %s. Gràcies! \n. . . \nMe dones mes tabac?\n", string(msg.Body))
            // Demanar mes tabac
            ch.Publish("", "tabac", false, false, amqp.Publishing{
                Body: []byte("Petició de tabac"),
            })
        case <- alertMsgs:
            fmt.Println("Anem, que ve la policia!")
            return
        }
    }
}
