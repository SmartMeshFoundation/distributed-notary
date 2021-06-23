package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

//缺省主链接入点 spectrum
var MainChainEndpoint = "http://106.52.171.12:18003"

//侧链接入点 heco
var SideChainEndpoint = "http://106.52.171.12:8545"

//default port
var Port = 8081

func RestMain() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get("/pubkey2address/:pubkey", pubkey2Address),
		rest.Post("/generateTx", generateTx),
		rest.Post("/sendTx", sendTx),
		rest.Post("/querystatus", queryStatus),
		rest.Get("/generateSecret", generateSecret),
		rest.Post("/scPrepareLockin", scPrepareLockin),
		rest.Post("/mcPrepareLockout", mcPrepareLockout),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	http.Handle("/api/", http.StripPrefix("/api", api.MakeHandler()))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./app"))))
	log.Printf("port=%d,mainchain=%s,sidechain=%s", Port, MainChainEndpoint, SideChainEndpoint)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", Port), nil))
}
