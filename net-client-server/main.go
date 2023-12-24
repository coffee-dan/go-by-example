package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	buf := make([]byte, 13)
	conn.Read(buf)
	fmt.Printf("%s", buf)
	conn.Close()
}

func startServer() {
	lnr, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer lnr.Close()

	for {
		conn, err := lnr.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func startClient() {
	var err error

	rdr := bufio.NewReader(os.Stdin)
	var msg string

	for {
		fmt.Printf("Huh?: ")
		msg, err = rdr.ReadString('\n')
		if err != nil {
			panic(err)
		}
		msg = strings.TrimSpace(msg)

		if len(msg) > 0 {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				panic(err)
			}

			fmt.Fprint(conn, msg)
		}
	}
}

func main() {
	rdr := bufio.NewReader(os.Stdin)
	fmt.Print("What?: ")
	cmd, err := rdr.ReadString('\n')
	if err != nil {
		panic(err)
	}

	cmd = strings.TrimSpace(cmd)

	switch cmd {
	case "server", "s":
		startServer()
	case "client", "c":
		startClient()
	}
}
