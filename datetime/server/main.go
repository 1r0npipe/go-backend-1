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

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			// continue
		}
		go sendMessage(conn)
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n\r"))
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func sendMessage(c net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter message >")
		str,_ := reader.ReadString('\n')
		_, err := io.WriteString(c, str)
		if err != nil {
			return
		}
		
	}
}