package gamethrive

import (
	"fmt"
	"net/url"
)

func mustParse(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(fmt.Sprintf("Must Parse: %s", err.Error()))
	}
	return u
}
