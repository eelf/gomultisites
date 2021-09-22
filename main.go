package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	log.Println(os.Args, os.Args[1])

	configFile, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	configDecoder := json.NewDecoder(configFile)
	type SiteType struct {
		Port string `json:"port"`
		Key string `json:"key"`
		Cert string `json:"cert"`
	}
	var config struct{
		CertBase string `json:"cert_base"`
		Sites map[string]SiteType `json:"sites"`
	}
	err = configDecoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{}
	m := map[string]http.Handler{}
	for host, siteConfig := range config.Sites {
		m[host] = httputil.NewSingleHostReverseProxy(&url.URL{
			Scheme: "http",
			Host:   "127.0.0.1:" + siteConfig.Port,
		})
		if len(siteConfig.Key) != 0 && len(siteConfig.Cert) != 0 {
			cert, err := tls.LoadX509KeyPair(config.CertBase + siteConfig.Cert, config.CertBase + siteConfig.Key)
			if err != nil {
				log.Fatal(err)
			}
			cfg.Certificates = append(cfg.Certificates, cert)
		}
	}

	ha := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p, ok := m[r.Host]; ok {
			p.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	hh := http.Server{
		Addr: ":80",
		Handler: ha,
	}
	hs := http.Server{
		Addr:      ":443",
		Handler:   ha,
		TLSConfig: cfg,
	}

	signalCh := make(chan os.Signal)
	go func() {
		<-signalCh
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		go func() {
			err := hh.Shutdown(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}()
		err := hs.Shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()
	signal.Notify(signalCh, os.Signal(syscall.SIGTERM))

	if len(cfg.Certificates) != 0 {
		go func() {
			log.Fatal(hs.ListenAndServeTLS("", ""))
		}()
	}
	log.Println("hh listenAndServe")
	log.Fatal(hh.ListenAndServe())
}
