package db

import "strconv"

func IDToStr(id uint32) string {
  return strconv.FormatUint(uint64(id), 10)
}
