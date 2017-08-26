package main

import "log"
import "runtime"
import "app/db"

func main() {
  log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile | log.LUTC)

  LogMemory("PREPARE")
  db.PrepareLocations()
  db.PrepareUsers()
  db.PrepareVisits()
  db.PrepareUserVisits()
  db.PrepareLocationVisits()

  LogMemory("IMPORT")
  if err := Import("/tmp/unzip"); err != nil {
    log.Fatal(err)
  }

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
// TODO Fix countries index
// TODO Parallel operation if possible
// TODO POST workers pool
