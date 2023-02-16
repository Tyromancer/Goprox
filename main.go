package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func handleClient(conn net.Conn) {
	// TODO: send client status codes

	defer conn.Close()
	clientReader := bufio.NewReader(conn)
	req, err := http.ReadRequest(clientReader)
	if err != nil {
		log.Println("Error when handling request: ", err)
		return
	}

	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}
	serverConn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Println("Error connecting to server: ", err)
		return
	}
	defer serverConn.Close()
	err = req.Write(serverConn)
	if err != nil {
		log.Println("Error sending request to server: ", err)
		return
	}

	switch req.Method {
	case "GET":
		log.Println("Handling GET request to " + serverConn.RemoteAddr().String())
		serverReader := bufio.NewReader(serverConn)
		res, err := http.ReadResponse(serverReader, req)
		if err != nil {
			log.Println("Error reading server response: ", err)
			return
		}
		defer res.Body.Close()

		err = res.Write(conn)
		if err != nil {
			log.Println("Error forwarding response to client: ", err)
			return
		}

	case "POST":
		log.Println("Handling POST request to " + serverConn.RemoteAddr().String())
		serverReader := bufio.NewReader(serverConn)
		res, err := http.ReadResponse(serverReader, req)
		if err != nil {
			log.Println("Error reading server response: ", err)
			return
		}
		defer res.Body.Close()

		err = res.Write(conn)
		if err != nil {
			log.Println("Error forwarding response to client: ", err)
			return
		}
	case "CONNECT":
		log.Println("Handling CONNECT request to " + serverConn.RemoteAddr().String())
		conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		go io.Copy(conn, serverConn)
		io.Copy(serverConn, conn)
	}

}

func main() {
	// TODO: use cobra to let the user config port number
	listener, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer listener.Close()
	fmt.Println("Server started on ", ":12345")

	for {
		conn, err := listener.Accept()
		log.Println("Accepted connection from " + conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("Error during connection: ", err)
			continue
		}

		go handleClient(conn)
	}
}
