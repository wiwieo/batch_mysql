package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"wiwieo/batch_mysql/server"
)

func main() {
	flag.Parse()
	s := server.NewSrv()
	http.HandleFunc("/add-ad", s.AdHandle)
	http.HandleFunc("/add-dt", s.DebrisTypeHandle)
	http.HandleFunc("/test", s.TestHandle)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		panic(err)
	}
}
