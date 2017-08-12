package main

// TODO "github.com/valyala/fasthttp"

import (
  "log"
)

var archivePath = "/tmp/data/data.zip"
var dataPath = "/tmp/unzip"

var httpAddr = ":80"

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)
  if err := Import(archivePath, dataPath); err != nil {
    log.Fatal(err)
  }
  if err := Serve(httpAddr); err != nil {
    log.Fatal(err)
  }
}

// TODO Fix 200 != 400 on POST (both new and update)
// TODO Use redis
// TODO Save data async, return nil error after validation
// TODO Return 200 for POST and 400 for GET after 100ms anyway
// TODO Cache plain JSON for each GET /<entity>/<id>
// TODO flusher, ok := w.(http.Flusher); flusher.Flush()
// TODO sync.Pool
// TODO Pointers
// TODO Avoid conversion between []byte and string
// TODO Check how it works without param queries
// TODO 1 GET => 1 POST => 1 GET => disable param queries after 20 secs
// TODO Add indexes for param queries


// FIXME Check if query param is null for POST
