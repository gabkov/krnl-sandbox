// package main

// import (
// 	"context"
// 	"fmt"
// 	"reflect"
// 	//"runtime"
// 	"strings"
// 	//"log"
// 	"net"
// 	"sync"
// 	// "unicode"
// 	"net/http"
// 	"encoding/json"
// 	"sync/atomic"
// 	crand "crypto/rand"
// 	"time"
// 	"encoding/binary"
// 	"math/rand"
// 	"encoding/hex"
// 	"unicode"
// 	"github.com/ethereum/go-ethereum/p2p/netutil"
// 	"github.com/gorilla/websocket"
// 	"io"
// 	"github.com/ethereum/go-ethereum/log"
// 	"errors"
// )

// const MetadataApi = "rpc"

// var (
// 	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
// 	errorType        = reflect.TypeOf((*error)(nil)).Elem()
// 	subscriptionType = reflect.TypeOf(Subscription{})
// 	stringType       = reflect.TypeOf("")
// )

// var (
// 	// ErrNotificationsUnsupported is returned by the client when the connection doesn't
// 	// support notifications. You can use this error value to check for subscription
// 	// support like this:
// 	//
// 	//	sub, err := client.EthSubscribe(ctx, channel, "newHeads", true)
// 	//	if errors.Is(err, rpc.ErrNotificationsUnsupported) {
// 	//		// Server does not support subscriptions, fall back to polling.
// 	//	}
// 	//
// 	//ErrNotificationsUnsupported = notificationsUnsupportedError{}

// 	// ErrSubscriptionNotFound is returned when the notification for the given id is not found
// 	ErrSubscriptionNotFound = errors.New("subscription not found")
// )

// // CodecOption specifies which type of messages a codec supports.
// //
// // Deprecated: this option is no longer honored by Server.
// type CodecOption int

// // A Subscription is created by a notifier and tied to that notifier. The client can use
// // this subscription to wait for an unsubscribe request for the client, see Err().
// type Subscription struct {
// 	ID        ID
// 	namespace string
// 	err       chan error // closed on unsubscribe
// }

// type serviceRegistry struct {
// 	mu       sync.Mutex
// 	services map[string]service
// }

// type service struct {
// 	name          string               // name for service
// 	callbacks     map[string]*callback // registered handlers
// 	subscriptions map[string]*callback // available subscriptions/notifications
// }

// // callback is a method callback which was registered in the server
// type callback struct {
// 	fn          reflect.Value  // the function
// 	rcvr        reflect.Value  // receiver object of method, set if fn is method
// 	argTypes    []reflect.Type // input argument types
// 	hasCtx      bool           // method's first argument is a context (not included in argTypes)
// 	errPos      int            // err return idx, of -1 when method cannot return error
// 	isSubscribe bool           // true if this is a subscription callback
// }

// // ID defines a pseudo random number that is used to identify RPC subscriptions.
// type ID string

// type Server struct {
// 	services serviceRegistry
// 	idgen    func() ID

// 	mutex              sync.Mutex
// 	codecs             map[ServerCodec]struct{}
// 	run                atomic.Bool
// 	batchItemLimit     int
// 	batchResponseLimit int
// }

// // PeerInfo contains information about the remote end of the network connection.
// //
// // This is available within RPC method handlers through the context. Call
// // PeerInfoFromContext to get information about the client connection related to
// // the current method call.
// type PeerInfo struct {
// 	// Transport is name of the protocol used by the client.
// 	// This can be "http", "ws" or "ipc".
// 	Transport string

// 	// Address of client. This will usually contain the IP address and port.
// 	RemoteAddr string

// 	// Additional information for HTTP and WebSocket connections.
// 	HTTP struct {
// 		// Protocol version, i.e. "HTTP/1.1". This is not set for WebSocket.
// 		Version string
// 		// Header values sent by the client.
// 		UserAgent string
// 		Origin    string
// 		Host      string
// 	}
// }

// type jsonError struct {
// 	Code    int         `json:"code"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// }

// // A value of this type can a JSON-RPC request, notification, successful response or
// // error response. Which one it is depends on the fields.
// type jsonrpcMessage struct {
// 	Version string          `json:"jsonrpc,omitempty"`
// 	ID      json.RawMessage `json:"id,omitempty"`
// 	Method  string          `json:"method,omitempty"`
// 	Params  json.RawMessage `json:"params,omitempty"`
// 	Error   *jsonError      `json:"error,omitempty"`
// 	Result  json.RawMessage `json:"result,omitempty"`
// }

// // ServerCodec implements reading, parsing and writing RPC messages for the server side of
// // a RPC session. Implementations must be go-routine safe since the codec can be called in
// // multiple go-routines concurrently.
// type ServerCodec interface {
// 	peerInfo() PeerInfo
// 	readBatch() (msgs []*jsonrpcMessage, isBatch bool, err error)
// 	close()

// 	jsonWriter
// }

// // jsonWriter can write JSON messages to its underlying connection.
// // Implementations must be safe for concurrent use.
// type jsonWriter interface {
// 	// writeJSON writes a message to the connection.
// 	writeJSON(ctx context.Context, msg interface{}, isError bool) error

// 	// Closed returns a channel which is closed when the connection is closed.
// 	closed() <-chan interface{}
// 	// RemoteAddr returns the peer address of the connection.
// 	remoteAddr() string
// }

// // randomIDGenerator returns a function generates a random IDs.
// func randomIDGenerator() func() ID {
// 	var buf = make([]byte, 8)
// 	var seed int64
// 	if _, err := crand.Read(buf); err == nil {
// 		seed = int64(binary.BigEndian.Uint64(buf))
// 	} else {
// 		seed = int64(time.Now().Nanosecond())
// 	}

// 	var (
// 		mu  sync.Mutex
// 		rng = rand.New(rand.NewSource(seed))
// 	)
// 	return func() ID {
// 		mu.Lock()
// 		defer mu.Unlock()
// 		id := make([]byte, 16)
// 		rng.Read(id)
// 		return encodeID(id)
// 	}
// }

// func encodeID(b []byte) ID {
// 	id := hex.EncodeToString(b)
// 	id = strings.TrimLeft(id, "0")
// 	if id == "" {
// 		id = "0" // ID's are RPC quantities, no leading zero's and 0 is 0x0.
// 	}
// 	return ID("0x" + id)
// }

// // RPCService gives meta information about the server.
// // e.g. gives information about the loaded modules.
// type RPCService struct {
// 	server *Server
// }

// // RegisterName creates a service for the given receiver type under the given name. When no
// // methods on the given receiver match the criteria to be either a RPC method or a
// // subscription an error is returned. Otherwise a new service is created and added to the
// // service collection this server provides to clients.
// func (s *Server) RegisterName(name string, receiver interface{}) error {
// 	return s.services.registerName(name, receiver)
// }

// func (r *serviceRegistry) registerName(name string, rcvr interface{}) error {
// 	rcvrVal := reflect.ValueOf(rcvr)
// 	if name == "" {
// 		return fmt.Errorf("no service name for type %s", rcvrVal.Type().String())
// 	}
// 	callbacks := suitableCallbacks(rcvrVal)
// 	if len(callbacks) == 0 {
// 		return fmt.Errorf("service %T doesn't have any suitable methods/subscriptions to expose", rcvr)
// 	}

// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if r.services == nil {
// 		r.services = make(map[string]service)
// 	}
// 	svc, ok := r.services[name]
// 	if !ok {
// 		svc = service{
// 			name:          name,
// 			callbacks:     make(map[string]*callback),
// 			subscriptions: make(map[string]*callback),
// 		}
// 		r.services[name] = svc
// 	}
// 	for name, cb := range callbacks {
// 		if cb.isSubscribe {
// 			svc.subscriptions[name] = cb
// 		} else {
// 			svc.callbacks[name] = cb
// 		}
// 	}
// 	return nil
// }

// // suitableCallbacks iterates over the methods of the given type. It determines if a method
// // satisfies the criteria for a RPC callback or a subscription callback and adds it to the
// // collection of callbacks. See server documentation for a summary of these criteria.
// func suitableCallbacks(receiver reflect.Value) map[string]*callback {
// 	typ := receiver.Type()
// 	callbacks := make(map[string]*callback)
// 	for m := 0; m < typ.NumMethod(); m++ {
// 		method := typ.Method(m)
// 		if method.PkgPath != "" {
// 			continue // method not exported
// 		}
// 		cb := newCallback(receiver, method.Func)
// 		if cb == nil {
// 			continue // function invalid
// 		}
// 		name := formatName(method.Name)
// 		callbacks[name] = cb
// 	}
// 	return callbacks
// }

// // makeArgTypes composes the argTypes list.
// func (c *callback) makeArgTypes() {
// 	fntype := c.fn.Type()
// 	// Skip receiver and context.Context parameter (if present).
// 	firstArg := 0
// 	if c.rcvr.IsValid() {
// 		firstArg++
// 	}
// 	if fntype.NumIn() > firstArg && fntype.In(firstArg) == contextType {
// 		c.hasCtx = true
// 		firstArg++
// 	}
// 	// Add all remaining parameters.
// 	c.argTypes = make([]reflect.Type, fntype.NumIn()-firstArg)
// 	for i := firstArg; i < fntype.NumIn(); i++ {
// 		c.argTypes[i-firstArg] = fntype.In(i)
// 	}
// }


// // newCallback turns fn (a function) into a callback object. It returns nil if the function
// // is unsuitable as an RPC callback.
// func newCallback(receiver, fn reflect.Value) *callback {
// 	fntype := fn.Type()
// 	c := &callback{fn: fn, rcvr: receiver, errPos: -1, isSubscribe: isPubSub(fntype)}
// 	// Determine parameter types. They must all be exported or builtin types.
// 	c.makeArgTypes()

// 	// Verify return types. The function must return at most one error
// 	// and/or one other non-error value.
// 	outs := make([]reflect.Type, fntype.NumOut())
// 	for i := 0; i < fntype.NumOut(); i++ {
// 		outs[i] = fntype.Out(i)
// 	}
// 	if len(outs) > 2 {
// 		return nil
// 	}
// 	// If an error is returned, it must be the last returned value.
// 	switch {
// 	case len(outs) == 1 && isErrorType(outs[0]):
// 		c.errPos = 0
// 	case len(outs) == 2:
// 		if isErrorType(outs[0]) || !isErrorType(outs[1]) {
// 			return nil
// 		}
// 		c.errPos = 1
// 	}
// 	return c
// }

// // Does t satisfy the error interface?
// func isErrorType(t reflect.Type) bool {
// 	return t.Implements(errorType)
// }

// // Is t Subscription or *Subscription?
// func isSubscriptionType(t reflect.Type) bool {
// 	for t.Kind() == reflect.Ptr {
// 		t = t.Elem()
// 	}
// 	return t == subscriptionType
// }

// // isPubSub tests whether the given method has as as first argument a context.Context and
// // returns the pair (Subscription, error).
// func isPubSub(methodType reflect.Type) bool {
// 	// numIn(0) is the receiver type
// 	if methodType.NumIn() < 2 || methodType.NumOut() != 2 {
// 		return false
// 	}
// 	return methodType.In(1) == contextType &&
// 		isSubscriptionType(methodType.Out(0)) &&
// 		isErrorType(methodType.Out(1))
// }

// // formatName converts to first character of name to lowercase.
// func formatName(name string) string {
// 	ret := []rune(name)
// 	if len(ret) > 0 {
// 		ret[0] = unicode.ToLower(ret[0])
// 	}
// 	return string(ret)
// }

// func sequentialIDGenerator() func() ID {
// 	var (
// 		mu      sync.Mutex
// 		counter uint64
// 	)
// 	return func() ID {
// 		mu.Lock()
// 		defer mu.Unlock()
// 		counter++
// 		id := make([]byte, 8)
// 		binary.BigEndian.PutUint64(id, counter)
// 		return encodeID(id)
// 	}
// }

// func NewServer() *Server {
// 	server := &Server{
// 		idgen:  randomIDGenerator(),
// 		codecs: make(map[ServerCodec]struct{}),
// 	}
// 	server.run.Store(true)
// 	// Register the default service providing meta information about the RPC service such
// 	// as the services and methods it offers.
// 	rpcService := &RPCService{server}
// 	server.RegisterName(MetadataApi, rpcService)
// 	return server
// }

// func newTestServer() *Server {
// 	server := NewServer()
// 	server.idgen = sequentialIDGenerator()
// 	if err := server.RegisterName("test", new(testService)); err != nil {
// 		panic(err)
// 	}

// 	return server
// }

// type testService struct{}

// // Stop stops reading new requests, waits for stopPendingRequestTimeout to allow pending
// // requests to finish, then closes all codecs which will cancel pending requests and
// // subscriptions.
// func (s *Server) Stop() {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()

// 	if s.run.CompareAndSwap(true, false) {
// 		log.Warn("RPC server shutting down")
// 		for codec := range s.codecs {
// 			codec.close()
// 		}
// 	}
// }

// func (s *Server) trackCodec(codec ServerCodec) bool {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()

// 	if !s.run.Load() {
// 		return false // Don't serve if server is stopped.
// 	}
// 	s.codecs[codec] = struct{}{}
// 	return true
// }

// func (s *Server) untrackCodec(codec ServerCodec) {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()

// 	delete(s.codecs, codec)
// }

// // A HTTPAuth function is called by the client whenever a HTTP request is sent.
// // The function must be safe for concurrent use.
// //
// // Usually, HTTPAuth functions will call h.Set("authorization", "...") to add
// // auth information to the request.
// type HTTPAuth func(h http.Header) error

// type clientConfig struct {
// 	// HTTP settings
// 	httpClient  *http.Client
// 	httpHeaders http.Header
// 	httpAuth    HTTPAuth

// 	// WebSocket options
// 	wsDialer           *websocket.Dialer
// 	wsMessageSizeLimit *int64 // wsMessageSizeLimit nil = default, 0 = no limit

// 	// RPC handler options
// 	idgen              func() ID
// 	batchItemLimit     int
// 	batchResponseLimit int
// }

// type reconnectFunc func(context.Context) (ServerCodec, error)


// type readOp struct {
// 	msgs  []*jsonrpcMessage
// 	batch bool
// }

// // requestOp represents a pending request. This is used for both batch and non-batch
// // requests.
// type requestOp struct {
// 	ids         []json.RawMessage
// 	err         error
// 	resp        chan []*jsonrpcMessage // the response goes here
// 	sub         *ClientSubscription    // set for Subscribe requests.
// 	hadResponse bool                   // true when the request was responded to
// }

// // ClientSubscription is a subscription established through the Client's Subscribe or
// // EthSubscribe methods.
// type ClientSubscription struct {
// 	client    *Client
// 	etype     reflect.Type
// 	channel   reflect.Value
// 	namespace string
// 	subid     string

// 	// The in channel receives notification values from client dispatcher.
// 	in chan json.RawMessage

// 	// The error channel receives the error from the forwarding loop.
// 	// It is closed by Unsubscribe.
// 	err     chan error
// 	errOnce sync.Once

// 	// Closing of the subscription is requested by sending on 'quit'. This is handled by
// 	// the forwarding loop, which closes 'forwardDone' when it has stopped sending to
// 	// sub.channel. Finally, 'unsubDone' is closed after unsubscribing on the server side.
// 	quit        chan error
// 	forwardDone chan struct{}
// 	unsubDone   chan struct{}
// }

// // Client represents a connection to an RPC server.
// type Client struct {
// 	idgen    func() ID // for subscriptions
// 	isHTTP   bool      // connection type: http, ws or ipc
// 	services *serviceRegistry

// 	idCounter atomic.Uint32

// 	// This function, if non-nil, is called when the connection is lost.
// 	reconnectFunc reconnectFunc

// 	// config fields
// 	batchItemLimit       int
// 	batchResponseMaxSize int

// 	// writeConn is used for writing to the connection on the caller's goroutine. It should
// 	// only be accessed outside of dispatch, with the write lock held. The write lock is
// 	// taken by sending on reqInit and released by sending on reqSent.
// 	writeConn jsonWriter

// 	// for dispatch
// 	close       chan struct{}
// 	closing     chan struct{}    // closed when client is quitting
// 	didClose    chan struct{}    // closed when client quits
// 	reconnected chan ServerCodec // where write/reconnect sends the new connection
// 	readOp      chan readOp      // read messages
// 	readErr     chan error       // errors from read
// 	reqInit     chan *requestOp  // register response IDs, takes write lock
// 	reqSent     chan error       // signals write completion, releases write lock
// 	reqTimeout  chan *requestOp  // removes response IDs when call timeout expires
// }

// type httpConn struct {
// 	client    *http.Client
// 	url       string
// 	closeOnce sync.Once
// 	closeCh   chan interface{}
// 	mu        sync.Mutex // protects headers
// 	headers   http.Header
// 	auth      HTTPAuth
// }

// // httpConn implements ServerCodec, but it is treated specially by Client
// // and some methods don't work. The panic() stubs here exist to ensure
// // this special treatment is correct.

// func (hc *httpConn) writeJSON(context.Context, interface{}, bool) error {
// 	panic("writeJSON called on httpConn")
// }

// func (hc *httpConn) peerInfo() PeerInfo {
// 	panic("peerInfo called on httpConn")
// }

// func (hc *httpConn) remoteAddr() string {
// 	return hc.url
// }

// func (hc *httpConn) readBatch() ([]*jsonrpcMessage, bool, error) {
// 	<-hc.closeCh
// 	return nil, false, io.EOF
// }

// func (hc *httpConn) close() {
// 	hc.closeOnce.Do(func() { close(hc.closeCh) })
// }

// func (hc *httpConn) closed() <-chan interface{} {
// 	return hc.closeCh
// }

// func initClient(conn ServerCodec, services *serviceRegistry, cfg *clientConfig) *Client {
// 	_, isHTTP := conn.(*httpConn)
// 	c := &Client{
// 		isHTTP:               isHTTP,
// 		services:             services,
// 		idgen:                cfg.idgen,
// 		batchItemLimit:       cfg.batchItemLimit,
// 		batchResponseMaxSize: cfg.batchResponseLimit,
// 		writeConn:            conn,
// 		close:                make(chan struct{}),
// 		closing:              make(chan struct{}),
// 		didClose:             make(chan struct{}),
// 		reconnected:          make(chan ServerCodec),
// 		readOp:               make(chan readOp),
// 		readErr:              make(chan error),
// 		reqInit:              make(chan *requestOp),
// 		reqSent:              make(chan error, 1),
// 		reqTimeout:           make(chan *requestOp),
// 	}

// 	// Set defaults.
// 	if c.idgen == nil {
// 		c.idgen = randomIDGenerator()
// 	}

// 	// Launch the main loop.
// 	if !isHTTP {
// 		go c.dispatch(conn)
// 	}
// 	return c
// }

// type handler struct {
// 	reg                  *serviceRegistry
// 	unsubscribeCb        *callback
// 	idgen                func() ID                      // subscription ID generator
// 	respWait             map[string]*requestOp          // active client requests
// 	clientSubs           map[string]*ClientSubscription // active client subscriptions
// 	callWG               sync.WaitGroup                 // pending call goroutines
// 	rootCtx              context.Context                // canceled by close()
// 	cancelRoot           func()                         // cancel function for rootCtx
// 	conn                 jsonWriter                     // where responses will be sent
// 	log                  log.Logger
// 	allowSubscribe       bool
// 	batchRequestLimit    int
// 	batchResponseMaxSize int

// 	subLock    sync.Mutex
// 	serverSubs map[ID]*Subscription
// }

// // unsubscribe is the callback function for all *_unsubscribe calls.
// func (h *handler) unsubscribe(ctx context.Context, id ID) (bool, error) {
// 	h.subLock.Lock()
// 	defer h.subLock.Unlock()

// 	s := h.serverSubs[id]
// 	if s == nil {
// 		return false, ErrSubscriptionNotFound
// 	}
// 	close(s.err)
// 	delete(h.serverSubs, id)
// 	return true, nil
// }

// func newHandler(connCtx context.Context, conn jsonWriter, idgen func() ID, reg *serviceRegistry, batchRequestLimit, batchResponseMaxSize int) *handler {
// 	rootCtx, cancelRoot := context.WithCancel(connCtx)
// 	h := &handler{
// 		reg:                  reg,
// 		idgen:                idgen,
// 		conn:                 conn,
// 		respWait:             make(map[string]*requestOp),
// 		clientSubs:           make(map[string]*ClientSubscription),
// 		rootCtx:              rootCtx,
// 		cancelRoot:           cancelRoot,
// 		allowSubscribe:       true,
// 		serverSubs:           make(map[ID]*Subscription),
// 		log:                  log.Root(),
// 		batchRequestLimit:    batchRequestLimit,
// 		batchResponseMaxSize: batchResponseMaxSize,
// 	}
// 	if conn.remoteAddr() != "" {
// 		h.log = h.log.New("conn", conn.remoteAddr())
// 	}
// 	h.unsubscribeCb = newCallback(reflect.Value{}, reflect.ValueOf(h.unsubscribe))
// 	return h
// }

// type clientContextKey struct{}

// type clientConn struct {
// 	codec   ServerCodec
// 	handler *handler
// }

// type peerInfoContextKey struct{}

// func (c *Client) newClientConn(conn ServerCodec) *clientConn {
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, clientContextKey{}, c)
// 	ctx = context.WithValue(ctx, peerInfoContextKey{}, conn.peerInfo())
// 	handler := newHandler(ctx, conn, c.idgen, c.services, c.batchItemLimit, c.batchResponseMaxSize)
// 	return &clientConn{conn, handler}
// }
// // close cancels all requests except for inflightReq and waits for
// // call goroutines to shut down.
// func (h *handler) close(err error, inflightReq *requestOp) {
// 	h.cancelAllRequests(err, inflightReq)
// 	h.callWG.Wait()
// 	h.cancelRoot()
// 	h.cancelServerSubscriptions(err)
// }

// func (cc *clientConn) close(err error, inflightReq *requestOp) {
// 	cc.handler.close(err, inflightReq)
// 	cc.codec.close()
// }

// // dispatch is the main loop of the client.
// // It sends read messages to waiting calls to Call and BatchCall
// // and subscription notifications to registered subscriptions.
// func (c *Client) dispatch(codec ServerCodec) {
// 	var (
// 		lastOp      *requestOp  // tracks last send operation
// 		reqInitLock = c.reqInit // nil while the send lock is held
// 		conn        = c.newClientConn(codec)
// 		reading     = true
// 	)
// 	defer func() {
// 		close(c.closing)
// 		if reading {
// 			conn.close(ErrClientQuit, nil)
// 			c.drainRead()
// 		}
// 		close(c.didClose)
// 	}()

// 	// Spawn the initial read loop.
// 	go c.read(codec)

// 	for {
// 		select {
// 		case <-c.close:
// 			return

// 		// Read path:
// 		case op := <-c.readOp:
// 			if op.batch {
// 				conn.handler.handleBatch(op.msgs)
// 			} else {
// 				conn.handler.handleMsg(op.msgs[0])
// 			}

// 		case err := <-c.readErr:
// 			conn.handler.log.Debug("RPC connection read error", "err", err)
// 			conn.close(err, lastOp)
// 			reading = false

// 		// Reconnect:
// 		case newcodec := <-c.reconnected:
// 			log.Debug("RPC client reconnected", "reading", reading, "conn", newcodec.remoteAddr())
// 			if reading {
// 				// Wait for the previous read loop to exit. This is a rare case which
// 				// happens if this loop isn't notified in time after the connection breaks.
// 				// In those cases the caller will notice first and reconnect. Closing the
// 				// handler terminates all waiting requests (closing op.resp) except for
// 				// lastOp, which will be transferred to the new handler.
// 				conn.close(errClientReconnected, lastOp)
// 				c.drainRead()
// 			}
// 			go c.read(newcodec)
// 			reading = true
// 			conn = c.newClientConn(newcodec)
// 			// Re-register the in-flight request on the new handler
// 			// because that's where it will be sent.
// 			conn.handler.addRequestOp(lastOp)

// 		// Send path:
// 		case op := <-reqInitLock:
// 			// Stop listening for further requests until the current one has been sent.
// 			reqInitLock = nil
// 			lastOp = op
// 			conn.handler.addRequestOp(op)

// 		case err := <-c.reqSent:
// 			if err != nil {
// 				// Remove response handlers for the last send. When the read loop
// 				// goes down, it will signal all other current operations.
// 				conn.handler.removeRequestOp(lastOp)
// 			}
// 			// Let the next request in.
// 			reqInitLock = c.reqInit
// 			lastOp = nil

// 		case op := <-c.reqTimeout:
// 			conn.handler.removeRequestOp(op)
// 		}
// 	}
// }

// // ServeCodec reads incoming requests from codec, calls the appropriate callback and writes
// // the response back using the given codec. It will block until the codec is closed or the
// // server is stopped. In either case the codec is closed.
// //
// // Note that codec options are no longer supported.
// func (s *Server) ServeCodec(codec ServerCodec, options CodecOption) {
// 	defer codec.close()

// 	if !s.trackCodec(codec) {
// 		return
// 	}
// 	defer s.untrackCodec(codec)

// 	cfg := &clientConfig{
// 		idgen:              s.idgen,
// 		batchItemLimit:     s.batchItemLimit,
// 		batchResponseLimit: s.batchResponseLimit,
// 	}
// 	c := initClient(codec, &s.services, cfg)
// 	<-codec.closed()
// 	c.Close()
// }

// // ServeListener accepts connections on l, serving JSON-RPC on them.
// func (s *Server) ServeListener(l net.Listener) error {
// 	for {
// 		conn, err := l.Accept()
// 		if netutil.IsTemporaryError(err) {
// 			log.Println("RPC accept error", "err", err)
// 			continue
// 		} else if err != nil {
// 			return err
// 		}
// 		log.Println("Accepted RPC connection", "conn", conn.RemoteAddr())
// 		go s.ServeCodec(NewCodec(conn), 0)
// 	}
// }

// func main() {
// 	server := newTestServer()
// 	defer server.Stop()

// 	listener, err := net.Listen("tcp", "127.0.0.1:0")
// 	if err != nil {
// 		log.Fatal("can't listen:", err)
// 	}
// 	defer listener.Close()
// 	go server.ServeListener(listener)

// 	var (
// 		request  = `{"jsonrpc":"2.0","id":1,"method":"rpc_modules"}` + "\n"
// 		wantResp = `{"jsonrpc":"2.0","id":1,"result":{"nftest":"1.0","rpc":"1.0","test":"1.0"}}` + "\n"
// 		deadline = time.Now().Add(10 * time.Second)
// 	)
// 	for i := 0; i < 20; i++ {
// 		conn, err := net.Dial("tcp", listener.Addr().String())
// 		if err != nil {
// 			t.Fatal("can't dial:", err)
// 		}

// 		conn.SetDeadline(deadline)
// 		// Write the request, then half-close the connection so the server stops reading.
// 		conn.Write([]byte(request))
// 		conn.(*net.TCPConn).CloseWrite()
// 		// Now try to get the response.
// 		buf := make([]byte, 2000)
// 		n, err := conn.Read(buf)
// 		conn.Close()

// 		if err != nil {
// 			t.Fatal("read error:", err)
// 		}
// 		if !bytes.Equal(buf[:n], []byte(wantResp)) {
// 			t.Fatalf("wrong response: %s", buf[:n])
// 		}
// 	}
// }