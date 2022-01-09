package main

import (
	"log"
	"net/http"
)

func main() {
	tlsCfg, e := generateTLSConfig()
	if e != nil {
		return
	}
	addr := ":8080"
	l, e := newListener(addr, tlsCfg)
	if e != nil {
		return
	}
	s := &http.Server{
		Addr:      addr,
		TLSConfig: tlsCfg,
	}
	// ハンドラをセット
	http.HandleFunc("/", helloHandler)
	log.Fatal(s.Serve(l))
}
