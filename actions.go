package main

import (
  "bytes"
  "errors"
  "github.com/valyala/fasthttp"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

var ErrNilParam = errors.New("nil param")
var ErrEmptyEntity = errors.New("empty entity")

var dummyResponse = []byte("{}")
var nullRequest = []byte("\": null")

func ActionGetUserVisits(ctx *fasthttp.RequestCtx, id uint32, v *fasthttp.Args) {
  visits, err := GetUserVisits(id, v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSONUserVisits(ctx, visits)
}

func ActionGetLocationAvg(ctx *fasthttp.RequestCtx, id uint32, v *fasthttp.Args) {
  avg, err := GetLocationAvg(id, v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSONLocationAvg(ctx, avg)
}

func ActionNewLocation(ctx *fasthttp.RequestCtx) {
  e := &structs.Location{}
  if err := checkRequestLocation(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddLocationAsync(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionNewUser(ctx *fasthttp.RequestCtx) {
  e := &structs.User{}
  if err := checkRequestUser(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddUserAsync(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionNewVisit(ctx *fasthttp.RequestCtx) {
  e := &structs.Visit{}
  if err := checkRequestVisit(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddVisitAsync(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateLocation(ctx *fasthttp.RequestCtx, id uint32) {
  e := &structs.Location{}
  if err := checkRequestLocation(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateLocationAsync(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateUser(ctx *fasthttp.RequestCtx, id uint32) {
  e := &structs.User{}
  if err := checkRequestUser(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateUserAsync(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateVisit(ctx *fasthttp.RequestCtx, id uint32) {
  e := &structs.Visit{}
  if err := checkRequestVisit(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateVisitAsync(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

// func checkRequest(ctx *fasthttp.RequestCtx, e interface{}) error {
//   body := ctx.PostBody()
//   if err := checkNils(ctx, body); err != nil {
//     return err
//   }
//   if err := ffjson.Unmarshal(body, e); err != nil {
//     return err
//   }
//   if e == nil {
//     return ErrEmptyEntity
//   }
//   return nil
// }

func checkRequestLocation(ctx *fasthttp.RequestCtx, e *structs.Location) error {
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    return err
  }
  if err := ffjson.Unmarshal(body, e); err != nil {
    return err
  }
  if e == nil {
    return ErrEmptyEntity
  }
  return nil
}

func checkRequestUser(ctx *fasthttp.RequestCtx, e *structs.User) error {
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    return err
  }
  if err := ffjson.Unmarshal(body, e); err != nil {
    return err
  }
  if e == nil {
    return ErrEmptyEntity
  }
  return nil
}

func checkRequestVisit(ctx *fasthttp.RequestCtx, e *structs.Visit) error {
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    return err
  }
  if err := ffjson.Unmarshal(body, e); err != nil {
    return err
  }
  if e == nil {
    return ErrEmptyEntity
  }
  return nil
}

func checkNils(ctx *fasthttp.RequestCtx, body []byte) error {
  if bytes.Contains(body, nullRequest) {
    return ErrNilParam
  }
  return nil

  // var m map[string]interface{}
  // if err := ffjson.Unmarshal(body, &m); err != nil {
  //   return err
  // }
  // for _, v := range m {
  //   if v == nil {
  //     return ErrNilParam
  //   }
  // }
  // return nil
}
