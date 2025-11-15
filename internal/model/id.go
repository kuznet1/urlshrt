package model

import (
	"fmt"
	"github.com/kuznet1/urlshrt/internal/errs"
	"net/http"
	"strconv"
)

// URLID is the compact base-36 identifier of a shortened URL.

type URLID uint64

// ID returns the numeric value of the identifier.
func (urlId URLID) ID() uint64 {
	return uint64(urlId)
}

// String returns the base-36 textual representation of the identifier.
func (urlId URLID) String() string {
	return strconv.FormatUint(urlId.ID(), 36)
}

// AsURL builds an absolute short URL using the given base prefix.
func (urlId URLID) AsURL(prefix string) string {
	return prefix + "/" + urlId.String()
}

// ParseURLID parses a base-36 textual identifier into URLID.
// It returns an HTTPError with 400 status for invalid input.
func ParseURLID(id string) (URLID, error) {
	val, err := strconv.ParseUint(id, 36, 64)
	if err != nil {
		return 0, errs.NewHTTPError(fmt.Sprintf("unable to parse %q: it must be alphanumeric", id), http.StatusBadRequest)
	}
	return URLID(val), nil
}
