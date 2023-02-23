package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"goprox/proxy"
	"goprox/util/logging"
	"net"
)

func main() {
	// parsing command line arguments
	pflag.StringP("host", "h", "127.0.0.1", "host to listen on")
	pflag.StringP("port", "p", "12345", "port to listen on")
	pflag.Parse()

	// setup logging
	logger := logging.LoggerSetup()
	defer logger.Sync()
	sugar := logger.Sugar()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		sugar.Fatalln("Error when parsing flags: ", err)
	}

	host := viper.GetString("host")
	port := viper.GetString("port")

	listener, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		sugar.Fatalf("Error when trying to listen on %s:%s\n", host, port)
	}

	defer listener.Close()
	sugar.Infof("Goprox server started on %s:%s\n", host, port)

	for {
		conn, err := listener.Accept()
		sugar.Info("Accepted connection from " + conn.RemoteAddr().String())
		if err != nil {
			sugar.Warn("Error during connection: ", err)
			continue
		}

		go proxy.HandleClient(conn, sugar)
	}
}
