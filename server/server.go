package main

import (
	"log"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func (d *serverDaemon) startServer() {

	certCacheDir := "/etc/letsencrypt"

	m := &autocert.Manager{
		Cache:      autocert.DirCache(certCacheDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("trader.ballresin.com", "fleetcmdr.com"),
	}

	d.hs.Addr = ":2213"
	d.hs.Handler = d.router
	d.hs.ReadTimeout = 5 * time.Minute
	d.hs.WriteTimeout = 5 * time.Minute
	d.hs.MaxHeaderBytes = 1 << 20
	d.hs.TLSConfig = m.TLSConfig()

	log.Fatal(d.hs.ListenAndServeTLS("", ""))
}
