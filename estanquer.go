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

    //Connexió amd RabbitMQ 
    //(Canviar localhost al nom del contenedor de docker en cas de usar-ho)
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

    //Suscripció a la cua de l'alerta
    err = ch.QueueBind(alertQueue.Name, "", "alerta", false, nil)

    // Variable per mantenir el número de tabac i de mistros
    var tabacCounter int
    var mistrosCounter int
    var partir = false

    // Consumir missatges de les cues
    go func() {
        //Missatges dels fumadorTabac.go
        tabacMsgs, err := ch.Consume(tabacQueue.Name, "", true, false, false, false, nil)
        if err != nil {
            log.Fatal(err)
        }

        //Missatges dels fumadorMistros.go
        mistrosMsgs, err := ch.Consume(mistrosQueue.Name, "", true, false, false, false, nil)
        if err != nil {
            log.Fatal(err)
        }
        //Alerta del delator
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
                //Posar tabac damunt la taula
                ch.Publish("", "fumadorTabac", false, false, amqp.Publishing{
                    Body: []byte(fmt.Sprintf("%d", tabacCounter)),
                })
            case <-mistrosMsgs:  // Rebre peticions de mistros
            mistrosCounter++    // Incrementar numero de mistros
                fmt.Printf("He posat el mistro %d damunt la taula\n", mistrosCounter)
                // Posar mistro damunt la taula
                ch.Publish("", "fumadorMistros", false, false, amqp.Publishing{
                    Body: []byte(fmt.Sprintf("%d", mistrosCounter)),
                })
            case <-alertMsgs:  // Rebre l'alerta del delator
                // Generar un temps d'espera aleatori
                waitTime := rand.Intn(1) + 1
                time.Sleep(time.Duration(waitTime) * time.Second)
                fmt.Println("Uyuyuy, la policia. Men vaig \n. . . M'enduc la taula")
                // Esborrar totes les cues
                ch.QueueDelete(tabacQueue.Name, false, false, false)
                ch.QueueDelete(mistrosQueue.Name, false, false, false)
                ch.QueueDelete(alertQueue.Name, false, false, false)
                ch.QueueDelete("alertaFumadorTabac", false, false, false)
                ch.QueueDelete("alertaFumadorMistros", false, false, false)
                ch.QueueDelete("fumadorTabac", false, false, false)
                ch.QueueDelete("fumadorMistros", false, false, false)
                partir = true
                return
            }
        }
    }()

    // Mantenir el programa en marxa fins que vengui la policia
    for !partir{
        
    }
}

