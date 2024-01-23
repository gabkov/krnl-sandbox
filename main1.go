package main

import (
	//"context"
	//"encoding/binary"
	//"errors"
	//"strings"
	//"sync"
	"time"
	"log"
	"net"
	"bytes"

	"github.com/gabkov/krnl-node/rpc"
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

func newTestServer() *rpc.Server {
	server := rpc.NewServer()
	//server.idgen = sequentialIDGenerator()
	if err := server.RegisterName("test", new(testService)); err != nil {
		panic(err)
	}
	if err := server.RegisterName("nftest", new(notificationTestService)); err != nil {
		panic(err)
	}
	return server
}

func main() {
	server := newTestServer()
	defer server.Stop()

	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("can't listen:", err)
	}
	log.Println(listener.Addr())
	defer listener.Close()
	go server.ServeListener(listener)

	var (
		request  = `{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n"
		wantResp = `{"jsonrpc":"2.0","id":1,"result":{"nftest":"1.0","rpc":"1.0","test":"1.0"}}` + "\n"
		deadline = time.Now().Add(10 * time.Second)
	)

	for i := 0; i < 20; i++ {
		conn, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			log.Fatal("can't dial:", err)
		}

		conn.SetDeadline(deadline)
		// Write the request, then half-close the connection so the server stops reading.
		conn.Write([]byte(request))
		conn.(*net.TCPConn).CloseWrite()
		// Now try to get the response.
		buf := make([]byte, 2000)
		n, err := conn.Read(buf)
		conn.Close()

		if err != nil {
			log.Fatal("read error:", err)
		}
		//log.Println("LOL")
		//log.Println(buf[:n])
		if !bytes.Equal(buf[:n], []byte(wantResp)) {
			log.Fatalf("wrong response: %s", buf[:n])
		}
	}
}
