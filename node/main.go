package main

import (
	"log"
	"net/http"

	hs "github.com/gabkov/krnl-node/httpserver"
	"github.com/gabkov/krnl-node/rpc"
	"github.com/gabkov/krnl-node/service"
	"github.com/joho/godotenv"
)

/*
Author: Gabor Kovacs | gabor.kovacs995@gmail.com | gabkov
*/

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// most of the rpc server setup is lifted over from geth
	srv := rpc.NewServer()

	// register krnl namespace and rpc calls
	if err := srv.RegisterName("krnl", new(service.Krnl)); err != nil {
		panic(err)
	}

	// register eth namespace and rpc calls
	if err := srv.RegisterName("eth", new(service.Eth)); err != nil {
		panic(err)
	}

	// register net namespace and rpc calls
	if err := srv.RegisterName("net", new(service.Net)); err != nil {
		panic(err)
	}

	log.Println("starting krnl node")
	// starting the http server so we can accept incomming rpc requests
	// similarly to the rpc setup, the http server setup is lifted over from geth
	httpsrv := hs.NewHttpServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}))

	defer httpsrv.Close()
	defer srv.Stop()
}
