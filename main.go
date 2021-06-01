// Copyright 2015 sms-api-server authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// HTTP API for sending SMS via SMPP.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	_ "net/http/pprof"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/go-web/httplog"

	"github.com/engagespark/sms-api-server/apiserver"
)

// Version of this server.
var Version = "v1.2.2"

type Opts struct {
	ListenAddr        string
	APIPrefix         string
	PublicDir         string
	Log               bool
	LogTS             bool
	CAFile            string
	CertFile          string
	KeyFile           string
	SMPPAddr          string
	ClientTLS         bool
	ClientTLSInsecure bool
	ShowVersion       bool
}

func main() {
	o := ParseOpts()
	if o.ShowVersion {
		fmt.Println("sms-api-server", Version)
		os.Exit(0)
	}
	tx := &smpp.Transceiver{
		Addr:   o.SMPPAddr,
		User:   os.Getenv("SMPP_USER"),
		Passwd: os.Getenv("SMPP_PASSWD"),
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)
	go func() {
		<-exit
		tx.Close()
		os.Exit(0)
	}()
	if o.ClientTLS {
		host, _, _ := net.SplitHostPort(tx.Addr)
		tx.TLS = &tls.Config{
			ServerName: host,
		}
		if o.ClientTLSInsecure {
			tx.TLS.InsecureSkipVerify = true
		}
	}
	api := &apiserver.Handler{Prefix: o.APIPrefix, Tx: tx}
	conn := api.Register(http.DefaultServeMux)
	go func() {
		for c := range conn {
			m := fmt.Sprintf("SMPP connection status to %s: %s",
				o.SMPPAddr, c.Status())
			if err := c.Error(); err != nil {
				m = fmt.Sprintf("%s (%v)", m, err)
			}
			log.Println(m)
		}
	}()
	if o.PublicDir != "" {
		fs := http.FileServer(http.Dir(o.PublicDir))
		http.Handle("/", http.StripPrefix(o.APIPrefix, fs))
	}
	mux := http.Handler(http.DefaultServeMux)
	if o.Log {
		var l *log.Logger
		if o.LogTS {
			l = log.New(os.Stderr, "", log.LstdFlags)
		} else {
			l = log.New(os.Stderr, "", 0)
		}
		mux = httplog.ApacheCombinedFormat(l)(mux.ServeHTTP)
	}
	err := ListenAndServe(o, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseOpts() *Opts {
	o := &Opts{ListenAddr: ":8080", SMPPAddr: "localhost:2775", LogTS: true}
	flag.StringVar(&o.ListenAddr, "http", o.ListenAddr, "host:port to listen on for http or https")
	flag.StringVar(&o.APIPrefix, "prefix", o.APIPrefix, "prefix for http(s) endpoints")
	flag.StringVar(&o.PublicDir, "public", o.PublicDir, "public dir to serve under \"/\", optional")
	flag.BoolVar(&o.Log, "log", o.Log, "log http requests")
	flag.BoolVar(&o.LogTS, "log-timestamp", o.LogTS, "add timestamp to logs")
	flag.StringVar(&o.CAFile, "ca", o.CAFile, "x509 CA certificate file (for client auth)")
	flag.StringVar(&o.CertFile, "cert", o.CertFile, "x509 certificate file for https server")
	flag.StringVar(&o.KeyFile, "key", o.KeyFile, "x509 key file for https server")
	flag.StringVar(&o.SMPPAddr, "smpp", o.SMPPAddr, "host:port of the SMSC to connect to via SMPP v3.4")
	flag.BoolVar(&o.ClientTLS, "tls", o.ClientTLS, "connect to SMSC using TLS")
	flag.BoolVar(&o.ClientTLSInsecure, "precaire", o.ClientTLSInsecure, "disable TLS checks for client connection")
	flag.BoolVar(&o.ShowVersion, "version", o.ShowVersion, "show version and exit")
	flag.Usage = func() {
		fmt.Printf("Usage: [env] %s [options]\n", os.Args[0])
		fmt.Printf("Environment variables:\n")
		fmt.Printf(" SMPP_USER: username for smpp client connection\n")
		fmt.Printf(" SMPP_PASSWD: password for smpp client connection\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	return o
}

func ListenAndServe(o *Opts, f http.Handler) error {
	s := &http.Server{Addr: o.ListenAddr, Handler: f}
	if o.CertFile == "" || o.KeyFile == "" {
		return s.ListenAndServe()
	}
	if o.CAFile != "" {
		b, err := ioutil.ReadFile(o.CAFile)
		if err != nil {
			return err
		}
		cp := x509.NewCertPool()
		cp.AppendCertsFromPEM(b)
		s.TLSConfig = &tls.Config{
			ClientCAs:  cp,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	}
	return s.ListenAndServeTLS(o.CertFile, o.KeyFile)
}
