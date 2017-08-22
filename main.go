package main

import (
  "log"
)

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

  log.Println("IMPORT")
  if err := Import("/tmp/unzip"); err != nil {
    log.Fatal(err)
  }

  log.Println("CACHE")
  PrepareCache()

  log.Println("SERVE")
  log.Fatal(Serve(":80"))
}

// TODO Add indexes for param queries
// TODO Easier age check
// TODO Check why there is a first long request
// TODO Fix countries index
// TODO Parallel operation if possible
