package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
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

func broadcaster() {
	clients := make(map[client]bool)
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

func handleConn(c net.Conn) {
	defer c.Close()
	ch := make(chan string)
	go clientWriter(c, ch)
	entering <- ch
	for {
		ch <- time.Now().Format("15:04:05")
		time.Sleep(1 * time.Second)
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func sendMessage() chan string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("send message > ")
		msg, _, err := reader.ReadLine()
		if err != nil {
			reader.Reset(os.Stdin)
			continue
		}
		messages <- string(msg)
	}
}