package main

import (
  "encoding/json"
  "errors"
  "github.com/valyala/fasthttp"
)

var ErrNilParam = errors.New("nil param")

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

func ActionUpdateLocation(ctx *fasthttp.RequestCtx, id uint32) {
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  e := new(Location)
  if err := json.Unmarshal(body, e); err != nil {
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
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  e := new(User)
  if err := json.Unmarshal(body, e); err != nil {
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
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  e := new(Visit)
  if err := json.Unmarshal(body, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := UpdateVisit(id, e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, DummyResponse{})
}

func ActionNewLocation(ctx *fasthttp.RequestCtx) {
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  e := new(Location)
  if err := json.Unmarshal(body, e); err != nil {
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
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  e := new(User)
  if err := json.Unmarshal(body, e); err != nil {
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
  body := ctx.PostBody()
  if err := checkNils(ctx, body); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  e := new(Visit)
  if err := json.Unmarshal(body, e); err != nil {
    ResponseStatus(ctx, 400)
    return
  }
  if err := AddVisit(e); err != nil {
    ResponseError(ctx, err)
    return
  }
  ResponseJSON(ctx, DummyResponse{})
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
