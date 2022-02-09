package main

import (
	"flag"
	"os"

	// _ "net/http/pprof"
	L "socks5relay/mylog"
	"socks5relay/s5relay"
)

var frontend = flag.String("frontend", "0.0.0.0:1080", "listen addr/port")
var backend = flag.String("backend", "", "backend addr/port (e.g 1.2.3.4:1080)")
var username = flag.String("username", "", "socks5 username")
var password = flag.String("password", "", "socks5 password")

func main() {
	flag.Parse()
	L.MylogInit("/tmp/socks5relay.log", true, "debug")
	if *backend == "" || *username == "" || *password == "" {
		flag.Usage()
		os.Exit(0)
	}
	s5 := s5relay.NewS5relay(*frontend, *backend, *username, *password)
	s5.Run()
}
