package api

import (
	"net/http"
	"testing"

	"log"

	"time"

	"strconv"

	"sync"

	"github.com/ant0ine/go-json-rest/rest"
	"golang.org/x/net/websocket"
)

var listen = "127.0.0.1:20000"
var mServer *sync.Map
var mClient *sync.Map

func TestWs(t *testing.T) {
	mServer = new(sync.Map)
	mClient = new(sync.Map)
	// 1. 启动服务端
	go startServer()
	time.Sleep(time.Second)
	// 2. 启动客户端
	go startClient(5)
	time.Sleep(100000 * time.Second)
	// 4. 回写
	temp, _ := mClient.Load(0)
	ws0 := temp.(*websocket.Conn)
	msg0 := "100"
	log.Printf("write %s to %p", msg0, ws0)
	if _, err2 := ws0.Write([]byte(msg0)); err2 != nil {
		log.Fatal(err2)
	}
	temp, _ = mServer.Load(0)
	ws1 := temp.(*websocket.Conn)
	msg1 := "200"
	log.Printf("write %s to %p", msg1, ws1)
	if _, err2 := ws1.Write([]byte(msg1)); err2 != nil {
		log.Fatal(err2)
	}
	time.Sleep(time.Second)
	// 5. 断开
	//time.Sleep(5 * time.Second)
	//log.Printf("close %p", ws0)
	//ws0.Close()
	//time.Sleep(5 * time.Second)
	//log.Printf("close %p", ws1)
	//ws1.Close()
	//
	//time.Sleep(5 * time.Second)

}

func startServer() {

	wsHandler := websocket.Handler(func(ws *websocket.Conn) {
		var msg = make([]byte, 512)
		var n int
		var err error
		if n, err = ws.Read(msg); err != nil {
			panic(err)
		}
		log.Printf("Server Received: %s at %p", msg[:n], ws)
		i, err := strconv.Atoi(string(msg[:n]))
		if err != nil {
			panic(err)
		}
		mServer.Store(i, ws)
		for {
			var msg = make([]byte, 512)
			var n int
			var err error
			if n, err = ws.Read(msg); err != nil {
				panic(err)
			}
			log.Printf("Gorouting Server Received: %s at %p", msg[:n], ws)
		}
	})

	router, err := rest.MakeRouter(
		rest.Get("/", func(w rest.ResponseWriter, r *rest.Request) {
			wsHandler.ServeHTTP(w.(http.ResponseWriter), r.Request)
		}),
	)
	if err != nil {
		panic(err)
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.SetApp(router)
	err = http.ListenAndServe("0.0.0.0:20000", api.MakeHandler())
	if err != nil {
		panic(err)
	}
}

func startClient(num int) {
	origin := "http://192.168.122.1:20000/"
	url := "ws://127.0.0.1:20000"
	for i := 0; i < num; i++ {
		msg := strconv.Itoa(i)
		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			panic(err)
		}
		if _, err2 := ws.Write([]byte(msg)); err2 != nil {
			log.Fatal(err2)
		}
		log.Printf("Send: %s at %p", msg, ws)
		mClient.Store(i, ws)
		time.Sleep(time.Second)
		go func() {
			for {
				var msg = make([]byte, 512)
				var n int
				var err error
				if n, err = ws.Read(msg); err != nil {
					panic(err)
				}
				log.Printf("Gorouting Client Received: %s at %p", msg[:n], ws)
			}
		}()
	}
	wsInterface, _ := mClient.Load(0)
	ws0 := wsInterface.(*websocket.Conn)
	for {
		time.Sleep(time.Second)
		ws0.Write([]byte("text"))
	}
}
