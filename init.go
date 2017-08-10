package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "path/filepath"
  "time"
)

func Import(src, dst string) error {
  List(filepath.Dir(src))
  files, err := Unzip(src, dst)
  if err != nil {
    return err
  }
  List(dst)
  log.Println(files)
  for _, filename := range files {
    // go func(filename string) {
      if err := Parse(filename); err != nil {
        log.Fatal(err)
      }
    // }(filename)
  }
  time.Sleep(3 * time.Second)
  // log.Println(GetUsers())
  return nil
}

func List(dir string) {
  log.Println("FILES: " + dir)
  files, _ := ioutil.ReadDir(dir)
  for _, f := range files {
    log.Println(f.Name())
  }
}

func Parse(filename string) error {
  data, err := ioutil.ReadFile(filename)
  if err != nil {
    return err
  }
  payload := Payload{}
  if err := json.Unmarshal(data, &payload); err != nil {
    return err
  }
  if err := Save(payload); err != nil {
    return err
  }
  return nil
}

func Save(payload Payload) error {
  for _, l := range payload.Locations {
    AddLocation(l)
  }
  for _, u := range payload.Users {
    AddUser(u)
  }
  for _, v := range payload.Visits {
    AddVisit(v)
  }
  return nil
}
