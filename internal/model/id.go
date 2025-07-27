package model

import (
	"errors"
	"strconv"
)

var ErrIDParsing = errors.New("id parsing error")

type URLID uint64

func (urlId URLID) ID() uint64 {
	return uint64(urlId)
}

func (urlId URLID) String() string {
	return strconv.FormatUint(urlId.ID(), 36)
}

func ParseURLID(id string) (URLID, error) {
	val, err := strconv.ParseUint(id, 36, 64)
	if err != nil {
		return 0, ErrIDParsing
	}
	return URLID(val), nil
}
