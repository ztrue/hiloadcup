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
  for {
    c, err := ln.Accept()
    if err != nil {
      return err
    }
    go func(c net.Conn) {
      defer c.Close()
      handler(c)
    }(c)
  }
}
