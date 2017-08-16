package main

import (
  "bytes"
  "errors"
  "log"
  "strconv"
  "github.com/valyala/fasthttp"
  "github.com/pquerna/ffjson/ffjson"
  "app/structs"
)

var ErrNilParam = errors.New("nil param")
var ErrEmptyEntity = errors.New("empty entity")

var dummyResponse = []byte("{}")
var nullRequest = []byte("\": null")

func ActionGetUserVisits(ctx *fasthttp.RequestCtx, bid []byte, v *fasthttp.Args) {
  visits, err := GetUserVisits(parseID(bid), v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSONUserVisits(ctx, visits)
}

func ActionGetLocationAvg(ctx *fasthttp.RequestCtx, bid []byte, v *fasthttp.Args) {
  avg, err := GetLocationAvg(parseID(bid), v)
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

func ActionUpdateLocation(ctx *fasthttp.RequestCtx, bid []byte) {
  e := &structs.Location{}
  if err := checkRequestLocation(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateLocationAsync(parseID(bid), e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateUser(ctx *fasthttp.RequestCtx, bid []byte) {
  e := &structs.User{}
  if err := checkRequestUser(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateUserAsync(parseID(bid), e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateVisit(ctx *fasthttp.RequestCtx, bid []byte) {
  e := &structs.Visit{}
  if err := checkRequestVisit(ctx, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateVisitAsync(parseID(bid), e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

// func checkRequest(ctx *fasthttp.RequestCtx, e interface{}) error {
//   if err := checkNils(ctx, ctx.PostBody()); err != nil {
//     return err
//   }
//   if err := ffjson.Unmarshal(ctx.PostBody(), e); err != nil {
//     return err
//   }
//   if e == nil {
//     return ErrEmptyEntity
//   }
//   return nil
// }

func checkRequestLocation(ctx *fasthttp.RequestCtx, e *structs.Location) error {
  if err := checkNils(ctx, ctx.PostBody()); err != nil {
    return err
  }
  if err := ffjson.Unmarshal(ctx.PostBody(), e); err != nil {
    return err
  }
  if e == nil {
    return ErrEmptyEntity
  }
  return nil
}

func checkRequestUser(ctx *fasthttp.RequestCtx, e *structs.User) error {
  if err := checkNils(ctx, ctx.PostBody()); err != nil {
    return err
  }
  if err := ffjson.Unmarshal(ctx.PostBody(), e); err != nil {
    return err
  }
  if e == nil {
    return ErrEmptyEntity
  }
  return nil
}

func checkRequestVisit(ctx *fasthttp.RequestCtx, e *structs.Visit) error {
  if err := checkNils(ctx, ctx.PostBody()); err != nil {
    return err
  }
  if err := ffjson.Unmarshal(ctx.PostBody(), e); err != nil {
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

func parseID(b []byte) uint32 {
  id64, err := strconv.ParseUint(string(b), 10, 32)
  if err != nil {
    log.Println(err)
    return 0
  }
  return uint32(id64)
}
