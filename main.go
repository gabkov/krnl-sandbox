package main

import (
	"log"
	"net/http"
	"github.com/gabkov/krnl-node/rpc"
	hs "github.com/gabkov/krnl-node/httpserver"
)

type testService struct{}

type echoArgs struct {
	S string
}

type echoResult struct {
	String string
	Int    int
	Args   *echoArgs
}

func (s *testService) Null() any {
	log.Println("kis faszocska")
	return nil
}

func (s *testService) Echo(str string, i int, args *echoArgs) echoResult {
	return echoResult{str, i, args}
}

type notificationTestService struct {
	unsubscribed            chan string
	gotHangSubscriptionReq  chan struct{}
	unblockHangSubscription chan struct{}
}

func (s *notificationTestService) Echo(i int) int {
	return i
}

func (s *notificationTestService) Unsubscribe(subid string) {
	if s.unsubscribed != nil {
		s.unsubscribed <- subid
	}
}


func main() {
	srv := rpc.NewServer()
	if err := srv.RegisterName("test", new(testService)); err != nil {
		panic(err)
	}
	if err := srv.RegisterName("nftest", new(notificationTestService)); err != nil {
		panic(err)
	}
	
	httpsrv := hs.NewHttpServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}))
	defer httpsrv.Close()
	defer srv.Stop()
}
