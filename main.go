package main

import (
  "log"
)

var archivePath = "/tmp/data/data.zip"
var dataPath = "/tmp/unzip"

var httpAddr = ":80"

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)
  log.Println("IMPORT")
  if err := Import(archivePath, dataPath); err != nil {
    log.Fatal(err)
  }
  log.Println("CACHE")
  PrepareCache()
  log.Println("SERVE")
  if err := Serve(httpAddr); err != nil {
    log.Fatal(err)
  }
}

// TODO Add indexes for param queries
// TODO Easier age check
// TODO Check why there is a first long request
// TODO Use ffjson
// TODO Cache UserVisits when POST
// TODO Fix countries index
// TODO Unzip using system util
