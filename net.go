package main

import (
  "net"
)

func ListenAndServe(addr string, handler func(net.Conn)) error {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
  defer ln.Close()
	return serve(ln, handler)
}

func serve(ln net.Listener, handler func(net.Conn)) error {
  jobs := make(chan net.Conn, 100)

  for i := 0; i < 4; i++ {
    go worker(jobs, handler)
  }

  for {
    c, err := ln.Accept()
    if err != nil {
      return err
    }
    jobs <- c
  }
}

func worker(jobs chan net.Conn, handler func(net.Conn)) {
  for c := range jobs {
    handler(c)
  }
}
