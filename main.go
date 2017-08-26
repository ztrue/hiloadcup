package main

import (
  "log"
  "runtime"
)

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

  LogMemory("IMPORT")
  // log.Println("IMPORT")
  if err := Import("/tmp/unzip"); err != nil {
    log.Fatal(err)
  }

  // LogMemory("CACHE")
  // // log.Println("CACHE")
  // PrepareCache()

  LogMemory("SERVE")
  // log.Println("SERVE")
  log.Fatal(Serve(":80"))
}

func LogMemory(mark string) {
  log.Println(mark)
  return
  var m runtime.MemStats
  runtime.ReadMemStats(&m)
  log.Printf(
    "%s | HeapSys = %d | StackSys = %d | MSpanSys = %d | OtherSys = %d | Sys = %d | NumGC = %d\n",
    mark,
    m.HeapSys / 1024,
    m.StackSys / 1024,
    m.MSpanSys / 1024,
    m.OtherSys / 1024,
    m.Sys / 1024,
    m.NumGC,
  )
}

// TODO Add indexes for param queries
// TODO Easier age check
// TODO Check why there is a first long request
// TODO Cache UserVisits when POST
// TODO Fix countries index
// TODO Parallel operation if possible
