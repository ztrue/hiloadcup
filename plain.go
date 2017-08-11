package main

import (
  "sync"
)

var m = &sync.Mutex{}

var data = map[string][]byte{}

func CacheSet(key string, val []byte) {
  m.Lock()
  data[key] = val
  m.Unlock()
}

func CacheGet(key string) ([]byte, bool) {
  val, ok := data[key]
  return val, ok
}
