package main

import (
	hs "github.com/gabkov/krnl-node/httpserver"
	"github.com/gabkov/krnl-node/rpc"
	"github.com/gabkov/krnl-node/service"
	"net/http"
)

/*
Author: Gabor Kovacs | gabor.kovacs995@gmail.com | gabkov
*/

func main() {
	srv := rpc.NewServer()

	if err := srv.RegisterName("krnl", new(service.Krnl)); err != nil {
		panic(err)
	}

	if err := srv.RegisterName("eth", new(service.Eth)); err != nil {
		panic(err)
	}

	if err := srv.RegisterName("net", new(service.Net)); err != nil {
		panic(err)
	}

	httpsrv := hs.NewHttpServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}))
	defer httpsrv.Close()
	defer srv.Stop()
}
