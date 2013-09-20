package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"
)

const (
	DefaultPolicy = "<cross-domain-policy><allow-access-from domain=\"*\" to-ports=\"*\"/></cross-domain-policy>\x00"
)

type FspServer struct {
	listener   *net.TCPListener
	policyData []byte
}

func NewFspServer() *FspServer {
	fsp := &FspServer{}
	fsp.policyData = []byte(DefaultPolicy)
	return fsp
}

func (fs *FspServer) Startup() {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:843")
	if err != nil {
		log.Fatalln(err)
	}
	fs.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Listening on %s", addr.String())

	for {
		conn, err := fs.listener.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}
		go fs.handleConnection(conn)
	}
	fs.listener.Close()
}

func (fs *FspServer) LoadPolicy(src io.Reader) {
	data, err := ioutil.ReadAll(bufio.NewReader(src))
	if err != nil {
		log.Fatalln(err)
	}
	fs.policyData = append(data, []byte("\x00")...)
}

func (fs *FspServer) handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	select {
	case ok := <-fs.read(conn):
		if !ok {
			return
		}
	case <-fs.timeout(time.Second * 2):
		log.Println("Timeout reading from", conn.RemoteAddr().String())
		return
	}

	w := bufio.NewWriterSize(conn, len(fs.policyData))
	w.Write(fs.policyData)
	w.Flush()
	log.Println("Sent policy file to", conn.RemoteAddr().String())
}

func (fs *FspServer) read(conn *net.TCPConn) (c chan bool) {
	c = make(chan bool)
	go func() {
		r := bufio.NewReader(conn)
		_, err := r.ReadString('\x00')
		if err != nil {
			if err != io.EOF && !strings.HasSuffix(err.Error(), "use of closed network connection") {
				log.Println(err)
			}
			c <- false
		} else {
			c <- true
		}
	}()
	return
}

func (fs *FspServer) timeout(d time.Duration) (c chan bool) {
	c = make(chan bool)
	go func() {
		time.Sleep(d)
		c <- true
	}()
	return
}
