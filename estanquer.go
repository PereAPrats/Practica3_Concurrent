package main

import (
    "fmt"
    "log"
    "time"
    "math/rand"
    "github.com/streadway/amqp"
)

func main() {
    fmt.Println("Hola, som l'estanquqer il·legal")

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

    // Creació de les cues per a tabac i mistros
    tabacQueue, err := ch.QueueDeclare("tabac", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    mistrosQueue, err := ch.QueueDeclare("mistros", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Creació de la cua per rebre l'alerta de la policia
    alertQueue, err := ch.QueueDeclare("alerta", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Variable per mantenir el número de tabac
    var tabacCounter int
    var mistrosCounter int

    // Consumir missatges de les cues
    go func() {
        tabacMsgs, err := ch.Consume(tabacQueue.Name, "", true, false, false, false, nil)
        if err != nil {
            log.Fatal(err)
        }

        mistrosMsgs, err := ch.Consume(mistrosQueue.Name, "", true, false, false, false, nil)
        if err != nil {
            log.Fatal(err)
        }

        alertMsgs, err := ch.Consume(alertQueue.Name, "", true, false, false, false, nil)
        if err != nil {
            log.Fatal(err)
        }

        // Esperar i gestionar els missatges
        for {
            select {
            case <-tabacMsgs:  // Rebre peticions de tabac
                tabacCounter++  // Incrementar el número de tabac
                fmt.Printf("He posat el tabac %d damunt la taula\n", tabacCounter)
                ch.Publish("", "fumadorTabac", false, false, amqp.Publishing{
                    Body: []byte(fmt.Sprintf("%d", tabacCounter)),
                })
            case <-mistrosMsgs:  // Rebre peticions de mistros
                mistrosCounter++
                fmt.Printf("He posat el mistro %d damunt la taula\n", mistrosCounter)
                ch.Publish("", "fumadorMistros", false, false, amqp.Publishing{
                    Body: []byte(fmt.Sprintf("%d", mistrosCounter)),
                })
            case <-alertMsgs:  // Rebre l'alerta de la policia
                // Generar un temps d'espera aleatori
                waitTime := rand.Intn(1) + 1
                time.Sleep(time.Duration(waitTime) * time.Second)
                fmt.Println("Uyuyuy, la policia. Men vaig \n. . . M'enduc la taula")
                // Esborrar totes les cues
                ch.QueueDelete(tabacQueue.Name, false, false, false)
                ch.QueueDelete(mistrosQueue.Name, false, false, false)
                ch.QueueDelete(alertQueue.Name, false, false, false)
                return
            }
        }
    }()

    // Mantenir el programa en marxa mentre espera els missatges
    select {}
}

