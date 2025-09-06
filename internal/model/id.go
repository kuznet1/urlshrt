package model

import (
	"fmt"
	"github.com/kuznet1/urlshrt/internal/errs"
	"net/http"
	"strconv"
)

type URLID uint64

func (urlId URLID) ID() uint64 {
	return uint64(urlId)
}

func (urlId URLID) String() string {
	return strconv.FormatUint(urlId.ID(), 36)
}

func (urlId URLID) AsURL(prefix string) string {
	return prefix + "/" + urlId.String()
}

func ParseURLID(id string) (URLID, error) {
	val, err := strconv.ParseUint(id, 36, 64)
	if err != nil {
		return 0, errs.NewHTTPError(fmt.Sprintf("unable to parse %q: it must be alphanumeric", id), http.StatusBadRequest)
	}
	return URLID(val), nil
}
