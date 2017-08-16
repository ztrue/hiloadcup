package main

import (
  "io/ioutil"
  "log"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

func Import(src, dst string) error {
  files, err := Unzip(src, dst)
  if err != nil {
    return err
  }
  for _, filename := range files {
    if err := Parse(filename); err != nil {
      log.Fatal(err)
    }
  }
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
  payload := &structs.Payload{}
  if err := ffjson.Unmarshal(data, payload); err != nil {
    return err
  }
  Save(payload)
  return nil
}

func Save(payload *structs.Payload) {
  for _, e := range payload.Locations {
    AddLocation(e)
  }
  for _, e := range payload.Users {
    AddUser(e)
  }
  for _, e := range payload.Visits {
    AddVisit(e)
  }
}
