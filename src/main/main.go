package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"os"
	"spy-cat/src/config"
	"spy-cat/src/db"
	"spy-cat/src/restful"
)

func main() {
	err := runRestApi()
	log.Fatal(err)
}

func init() {
	conf := os.Getenv("CVE_CONF_PATH")
	if conf == "" {
		exception(fmt.Errorf("CVE_CONF_PATH is required"))
	}
	co, err := config.Fetch(conf)
	exception(err)

	err = db.Init(&co.DB)
	exception(err)

}

func runRestApi() error {
	serv := &fasthttp.Server{
		Handler: restful.Handler(),
	}

	return serv.ListenAndServe(fmt.Sprintf("localhost:%v", 8081))
}

func exception(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
