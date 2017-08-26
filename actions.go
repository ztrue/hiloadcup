package main

import "bytes"
import "github.com/valyala/fasthttp"
import "github.com/pquerna/ffjson/ffjson"
import "app/db"
import "app/structs"

var dummyResponse = []byte("{}")
var nullRequest = []byte("\": n")

func ActionGetLocation(ctx *fasthttp.RequestCtx, bid []byte) {
  e := db.GetLocation(string(bid))
  if e == nil {
    ResponseStatus(ctx, 404)
    return
  }
  ResponseJSONLocation(ctx, e)
}

func ActionGetUser(ctx *fasthttp.RequestCtx, bid []byte) {
  e := db.GetUser(string(bid))
  if e == nil {
    ResponseStatus(ctx, 404)
    return
  }
  ResponseJSONUser(ctx, e)
}

func ActionGetVisit(ctx *fasthttp.RequestCtx, bid []byte) {
  e := db.GetVisit(string(bid))
  if e == nil {
    ResponseStatus(ctx, 404)
    return
  }
  ResponseJSONVisit(ctx, e)
}

func ActionGetUserVisits(ctx *fasthttp.RequestCtx, bid []byte, v *fasthttp.Args) {
  visits, status := GetUserVisits(bid, v)
  if status != 200 {
    ResponseStatus(ctx, status)
    return
  }
  ResponseJSONUserVisits(ctx, visits)
}

func ActionGetLocationAvg(ctx *fasthttp.RequestCtx, bid []byte, v *fasthttp.Args) {
  avg, status := GetLocationAvg(bid, v)
  if status != 200 {
    ResponseStatus(ctx, status)
    return
  }
  ResponseJSONLocationAvg(ctx, avg)
}

func ActionNewLocation(ctx *fasthttp.RequestCtx) {
  e := &structs.LocationUp{}
  if checkRequestLocation(ctx.PostBody(), e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  if AddLocationAsync(e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionNewUser(ctx *fasthttp.RequestCtx) {
  e := &structs.UserUp{}
  if checkRequestUser(ctx.PostBody(), e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  if AddUserAsync(e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionNewVisit(ctx *fasthttp.RequestCtx) {
  e := &structs.VisitUp{}
  if checkRequestVisit(ctx.PostBody(), e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  if AddVisitAsync(e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateLocation(ctx *fasthttp.RequestCtx, bid []byte) {
  if !db.LocationExists(string(bid)) {
    ResponseStatus(ctx, 404)
    return
  }
  e := &structs.LocationUp{}
  if checkRequestLocation(ctx.PostBody(), e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  if UpdateLocationAsync(bid, e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateUser(ctx *fasthttp.RequestCtx, bid []byte) {
  if !db.UserExists(string(bid)) {
    ResponseStatus(ctx, 404)
    return
  }
  e := &structs.UserUp{}
  if checkRequestUser(ctx.PostBody(), e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  if UpdateUserAsync(bid, e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func ActionUpdateVisit(ctx *fasthttp.RequestCtx, bid []byte) {
  if !db.VisitExists(string(bid)) {
    ResponseStatus(ctx, 404)
    return
  }
  e := &structs.VisitUp{}
  if checkRequestVisit(ctx.PostBody(), e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  if UpdateVisitAsync(bid, e) != 200 {
    ResponseStatus(ctx, 400)
    return
  }
  ResponseBytes(ctx, dummyResponse)
}

func checkRequestLocation(body []byte, e *structs.LocationUp) int {
  if checkNils(body) != 200 {
    return 400
  }
  if ffjson.Unmarshal(body, e) != nil {
    return 400
  }
  if e == nil {
    return 400
  }
  return 200
}

func checkRequestUser(body []byte, e *structs.UserUp) int {
  if checkNils(body) != 200 {
    return 400
  }
  if ffjson.Unmarshal(body, e) != nil {
    return 400
  }
  if e == nil {
    return 400
  }
  return 200
}

func checkRequestVisit(body []byte, e *structs.VisitUp) int {
  if checkNils(body) != 200 {
    return 400
  }
  if ffjson.Unmarshal(body, e) != nil {
    return 400
  }
  if e == nil {
    return 400
  }
  return 200
}

func checkNils(body []byte) int {
  if bytes.Contains(body, nullRequest) {
    return 400
  }
  return 200
}
