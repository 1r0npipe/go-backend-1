package main

import (
	"io"
	"log"
	"net"
	"os"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
      
      buf := make([]byte, 256) // создаем буфер
      for {
          _, err = conn.Read(buf)
          if err == io.EOF {
              break
          }
          io.WriteString(os.Stdout, fmt.Sprintf("Custom output! %s", string(buf))) // выводим измененное сообщение сервера в консоль
      }
}

