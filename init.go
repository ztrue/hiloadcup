package main

import (
  "encoding/json"
  "io/ioutil"
  "log"
)

func Import(src, dst string) error {
  // List(filepath.Dir(src))
  files, err := Unzip(src, dst)
  if err != nil {
    return err
  }
  // List(dst)
  // log.Println(files)
  for _, filename := range files {
    // go func(filename string) {
      if err := Parse(filename); err != nil {
        log.Fatal(err)
      }
    // }(filename)
  }
  // time.Sleep(3 * time.Second)
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
  Save(payload)
  return nil
}

func Save(payload Payload) {
  for _, e := range payload.Locations {
    if err := AddLocation(e); err != nil {
      log.Println(err)
    }
  }
  for _, e := range payload.Users {
    if err := AddUser(e); err != nil {
      log.Println(err)
    }
  }
  for _, e := range payload.Visits {
    if err := AddVisit(e); err != nil {
      log.Println(err)
    }
  }
}
