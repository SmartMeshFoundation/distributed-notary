package rest

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

//缺省主链接入点
var MainChainEndpoint = "http://127.0.0.1:8545"

//侧链接入点
var SideChainEndpoint = "http://127.0.0.1:8545"

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

	log.Fatal(http.ListenAndServe(":8080", nil))
}
