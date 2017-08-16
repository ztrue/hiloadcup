package main

import (
  "log"
)

var dataPath = "/tmp/unzip"

var httpAddr = ":80"

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)
  log.Println("IMPORT")
  if err := Import(dataPath); err != nil {
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
// TODO Cache UserVisits when POST
// TODO Fix countries index
