package s5relay

import (
	// _ "net/http/pprof"
	"bytes"
	"encoding/hex"
	"io"
	"net"
	L "socks5relay/mylog"
)

// _ "net/http/pprof"

type s5relay struct {
	frontend string
	backend  string
	username string
	password string
	sig      bytes.Buffer
}

func NewS5relay(frontend, backend, username, password string) *s5relay {
	s5 := &s5relay{
		frontend: frontend,
		backend:  backend,
		username: username,
		password: password,
	}

	sig := bytes.Buffer{}
	sig.WriteByte(byte(len(s5.username)))
	sig.WriteString(s5.username)
	sig.WriteByte(byte(len(s5.password)))
	sig.WriteString(s5.password)
	L.Debugf("sig: \n%s", hex.Dump(sig.Bytes()))

	s5.sig = sig
	return s5
}

func (s5 *s5relay) Run() {
	lner, err := net.Listen("tcp", s5.frontend)
	if err != nil {
		L.Fatalf("listen failed:%v", err)
	}
	defer lner.Close()

	for {
		conn, err := lner.Accept()
		if err != nil {
			L.Errorf("accept failed:%v", err)
			continue
		}
		L.Debugf("new conn from %s", conn.RemoteAddr().String())
		go s5.Relay(conn)
	}
}

func (s5 *s5relay) Relay(conn net.Conn) {
	defer conn.Close()
	bconn, err := net.Dial("tcp", s5.backend)
	if err != nil {
		L.Errorf("dial failed:%v", err)
		return
	}
	defer bconn.Close()

	go func() {
		io.Copy(conn, bconn)
	}()

	// io.Copy(bconn, conn)

	buf := make([]byte, 4096)
	nbytes := 0
	sigdone := false
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		if bytes.Contains(buf[:n], s5.sig.Bytes()) {
			L.Debug("found sig!")
			sigdone = true
		}

		nbytes += n

		if nbytes >= 100 && !sigdone {
			L.Errorf("can't find sig in first 100 bytes")
			break
		}

		_, err = bconn.Write(buf[:n])
		if err != nil {
			break
		}
	}
}
