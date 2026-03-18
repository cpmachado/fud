package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
)

func main() {
	port := 1337
	msg := "This is up"
	versionFlag := false

	flag.IntVar(&port, "port", port, "port to serve")
	flag.StringVar(&msg, "msg", msg, "message to show")
	flag.BoolVar(&versionFlag, "v", versionFlag, "display version")
	flag.Parse()

	if port < 0 || port > 60999 {
		log.Fatalf("invalid port: %d", port)
	}

	if versionFlag {
		log.Printf("fud-%s", version())
		os.Exit(0)
	}

	log.Printf("port: %d", port)
	log.Printf("msg: %q", msg)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("call from ip: %s", r.RemoteAddr)
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Fatal(err)
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	var ips []string

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Fatal(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip.String())
		}
	}

	log.Println("Listening in:")
	for _, ip := range ips {
		log.Printf("- http://%s:%d\n", ip, port)
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	version := "unknown"

	if ok {
		version = info.Main.Version
	}
	return version
}
