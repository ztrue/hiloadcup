package main

import (
  "encoding/json"
  "errors"
  "log"
  "github.com/valyala/fasthttp"
)

var ErrNilParam = errors.New("nil param")
var ErrEmptyEntity = errors.New("empty entity")

// func ActionGetLocation(ctx *fasthttp.RequestCtx, id uint32) {
//   l := GetLocation(id)
//   if l == nil {
//     ResponseStatus(ctx, 404)
//     return
//   }
//   ResponseJSON(ctx, l)
// }
//
// func ActionGetUser(ctx *fasthttp.RequestCtx, id uint32) {
//   u := GetUser(id)
//   if u == nil {
//     ResponseStatus(ctx, 404)
//     return
//   }
//   ResponseJSON(ctx, u)
// }
//
// func ActionGetVisit(ctx *fasthttp.RequestCtx, id uint32) {
//   v := GetVisit(id)
//   if v == nil {
//     ResponseStatus(ctx, 404)
//     return
//   }
//   ResponseJSON(ctx, v)
// }

type UserVisitsResponse struct {
  Visits []UserVisit `json:"visits"`
}

func ActionGetUserVisits(ctx *fasthttp.RequestCtx, userID uint32, v *fasthttp.Args) {
  visits, err := GetUserVisits(userID, v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, UserVisitsResponse{visits})
}

type LocationAvgResponse struct {
  Avg float32 `json:"avg"`
}

func ActionGetLocationAvg(ctx *fasthttp.RequestCtx, id uint32, v *fasthttp.Args) {
  avg, err := GetLocationAvg(id, v)
  if err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, LocationAvgResponse{avg})
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
  ResponseJSON(ctx, DummyResponse{})
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
  ResponseJSON(ctx, DummyResponse{})
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
  ResponseJSON(ctx, DummyResponse{})
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
  ResponseJSON(ctx, DummyResponse{})
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
  ResponseJSON(ctx, DummyResponse{})
}

func ActionUpdateVisit(ctx *fasthttp.RequestCtx, id uint32) {
  var e *Visit
  if err := checkRequest(ctx, &e); err != nil {
    log.Println(id, err)
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateVisit(id, e); err != nil {
    log.Println(id, err)
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, DummyResponse{})
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
