package main

// TODO "github.com/valyala/fasthttp"

import (
  "log"
)

var archivePath = "/tmp/data/data.zip"
var dataPath = "/tmp/unzip"

var httpAddr = ":80"

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC)
  if err := Import(archivePath, dataPath); err != nil {
    log.Fatal(err)
  }
  Prepare()
  if err := Serve(httpAddr); err != nil {
    log.Fatal(err)
  }
}

// TODO Use redis
// TODO Save data async, return nil error after validation
// TODO Return 200 for POST and 400 for GET after 100ms anyway
