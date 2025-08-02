package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
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
			b := make([]byte, 1024)

			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}

			n, err := conn.Read(b)
			if err != nil {
				// TODO: log error
				continue
			}

			go handleCommand(conn, b[:n])
		}
	}()

	<-done
	fmt.Println("Exiting...")
}

func handleCommand(conn net.Conn, cmd []byte) error {
	switch string(cmd) {
	case "PING":
		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			// TODO: log error
			// TODO: write error back
			return err
		}
	}

	return nil
}
