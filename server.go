package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]string)
var mu sync.Mutex

func broadcast(msg string) {
	mu.Lock()
	conns := make([]net.Conn, 0, len(clients))
	for c := range clients {
		conns = append(conns, c)
	}
	mu.Unlock()

	for _, c := range conns {
		_, err := c.Write([]byte(msg))
		if err != nil {
			fmt.Println("Write error:", err)
		}
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New client connected:", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read name:", err)
		return
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = "Oakheart"
	}
	conn.Write([]byte("Здарова окхарт\n"))

	mu.Lock()
	clients[conn] = name
	mu.Unlock()

	broadcast(fmt.Sprintf("%s joined the chat\n", name))

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())

			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			broadcast(fmt.Sprintf("%s left the chat\n", name))
			return
		}

		fullMsg := fmt.Sprintf("[%s]: %s", name, msg)
		broadcast(fullMsg)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConn(conn)
	}

	select {}
}
