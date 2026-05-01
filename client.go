package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	input := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	name, _ := input.ReadString('\n')
	name = strings.TrimSpace(name)

	conn.Write([]byte(name + "\n"))

	go func() {
		reader := bufio.NewReader(conn)

		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("disconnected from server")
				return
			}
			fmt.Println(msg)
		}
	}()

	for {
		fmt.Print("> ")

		text, err := input.ReadString('\n')
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}

		_, err = conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}
