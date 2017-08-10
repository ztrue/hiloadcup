package main

import (
  "net/http"
)

func Serve(addr string) error {
  http.HandleFunc("/", Handle)
  return http.ListenAndServe(addr, nil)
}

func Handle(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte(""))
}
