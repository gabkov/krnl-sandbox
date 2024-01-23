package main

import (
	"net/http"
	hs "github.com/gabkov/krnl-node/httpserver"
	"github.com/gabkov/krnl-node/rpc"
	"github.com/gabkov/krnl-node/service"
)

func main() {
	srv := rpc.NewServer()

	if err := srv.RegisterName("krnl", new(service.Krnl)); err != nil {
		panic(err)
	}

	httpsrv := hs.NewHttpServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}))
	defer httpsrv.Close()
	defer srv.Stop()
}
