package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, os.Kill)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			go handle(conn)
		}
	}()

	<-done
	fmt.Println("Exiting...")
}

func handle(conn net.Conn) error {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)

		switch text {
		case "PING":
			if _, err := conn.Write([]byte("+PING\r\n")); err != nil {
				// TODO: log error
				continue
			}
		}
	}

	return nil
}
