package proxy

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
)

func HandleClient(conn net.Conn, sugar *zap.SugaredLogger) {
	// TODO: send client status codes

	defer conn.Close()
	clientReader := bufio.NewReader(conn)
	req, err := http.ReadRequest(clientReader)
	if err != nil {
		sugar.Warn("Error when handling request: ", err)
		return
	}

	sugar.Info("Got request: ", req)

	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		port = "80"
	}
	serverConn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		sugar.Warn("Error connecting to server: ", err)
		return
	}
	defer serverConn.Close()
	err = req.Write(serverConn)
	if err != nil {
		sugar.Warn("Error sending request to server: ", err)
		return
	}

	switch req.Method {
	case "GET":
		sugar.Debug("Handling GET request to " + serverConn.RemoteAddr().String())
		serverReader := bufio.NewReader(serverConn)
		res, err := http.ReadResponse(serverReader, req)
		if err != nil {
			sugar.Warn("Error reading server response: ", err)
			return
		}
		defer res.Body.Close()

		err = res.Write(conn)
		if err != nil {
			sugar.Warn("Error forwarding response to client: ", err)
			return
		}

	case "POST":
		sugar.Debug("Handling POST request to " + serverConn.RemoteAddr().String())
		serverReader := bufio.NewReader(serverConn)
		res, err := http.ReadResponse(serverReader, req)
		if err != nil {
			sugar.Warn("Error reading server response: ", err)
			return
		}
		defer res.Body.Close()

		err = res.Write(conn)
		if err != nil {
			sugar.Warn("Error forwarding response to client: ", err)
			return
		}
	case "CONNECT":
		sugar.Debug("Handling CONNECT request to " + serverConn.RemoteAddr().String())
		conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		go io.Copy(conn, serverConn)
		io.Copy(serverConn, conn)
	}

}
