package main

import (
	"io"
	"log"
	"net"
	"time"
	"fmt"
	"bufio"
	"os"
)
type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
	clients  = make(map[client]bool)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go sendMessage()
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		select {
		case msg := <- messages:
			_, err := io.WriteString(c, msg)
			if err != nil {
				return
			}
		default:
			_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
			if err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}


func sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message >")
		str,_ := reader.ReadString('\n')
		messages <- str
	}
}

func broadcaster() {

	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}