package main

import (
  "encoding/json"
  "errors"
  "github.com/valyala/fasthttp"
)

var ErrNilParam = errors.New("nil param")
var ErrEmptyEntity = errors.New("empty entity")

var dummyResponse = []byte("{}")

func ActionGetUserVisits(ctx *fasthttp.RequestCtx, id uint32, v *fasthttp.Args) {
  visits, err := GetUserVisits(id, v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, visits)
}

func ActionGetLocationAvg(ctx *fasthttp.RequestCtx, id uint32, v *fasthttp.Args) {
  avg, err := GetLocationAvg(id, v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, avg)
}

type DummyResponse struct {}

func ActionNewLocation(ctx *fasthttp.RequestCtx) {
  var e *Location
  if err := checkRequest(ctx, &e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddLocation(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionNewUser(ctx *fasthttp.RequestCtx) {
  var e *User
  if err := checkRequest(ctx, &e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddUser(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionNewVisit(ctx *fasthttp.RequestCtx) {
  var e *Visit
  if err := checkRequest(ctx, &e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddVisit(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateLocation(ctx *fasthttp.RequestCtx, id uint32) {
  var e *Location
  if err := checkRequest(ctx, &e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateLocation(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateUser(ctx *fasthttp.RequestCtx, id uint32) {
  var e *User
  if err := checkRequest(ctx, &e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateUser(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateVisit(ctx *fasthttp.RequestCtx, id uint32) {
  var e *Visit
  if err := checkRequest(ctx, &e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateVisit(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func checkRequest(ctx *fasthttp.RequestCtx, e interface{}) error {
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    return err
  }
  if err := json.Unmarshal(body, e); err != nil {
    return err
  }
  if e == nil {
    return ErrEmptyEntity
  }
  return nil
}

func checkNils(ctx *fasthttp.RequestCtx, body []byte) error {
  var m map[string]interface{}
  if err := json.Unmarshal(body, &m); err != nil {
    return err
  }
  for _, v := range m {
    if v == nil {
      return ErrNilParam
    }
  }
  return nil
}
