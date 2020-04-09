package rest

import (
	"errors"
	"math"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

const (
	selfRel  string = "self"
	firstRel string = "first"
	prevRel  string = "prev"
	nextRel  string = "next"
	lastRel  string = "last"
)

func getIntFromQuery(i url.Values, n string) int {
	if n, err := strconv.Atoi(i.Get(n)); err == nil {
		return n
	}

	return 0
}

func getListParam(i url.Values, minP, defaultL int) *listParam {
	o := 0
	p := getIntFromQuery(i, "page")
	l := getIntFromQuery(i, "limit")

	if p < minP {
		p = minP
	}

	if l <= 0 {
		l = defaultL
	}

	if p > minP {
		o = (p - minP) * l
	}

	return &listParam{
		Page:   p,
		Offset: o,
		Limit:  l,
	}
}

func isUUID(id string, v int) (*uuid.UUID, bool) {
	i, err := uuid.Parse(id)

	if err != nil || i.Version() != uuid.Version(v) {
		return nil, false
	}

	return &i, true
}

func toCollectionResp(u *url.URL, page, limit, first int, coll *Collection) *CollectionResponse {
	return &CollectionResponse{
		Meta:  Meta{Total: coll.Total},
		Links: getHypertextCtrl(u, page, limit, coll.Total, first),
		Items: coll.Items,
	}
}

type getLinkFunc func(u *url.URL, limit int) (*Link, error)

func getHypertextCtrl(u *url.URL, page, limit, total, first int) []Link {
	last := first
	if total != 0 && limit != 0 {
		last = int(math.Ceil(float64(total) / float64(limit)))
	}

	m := []getLinkFunc{
		getPageLink(page, selfRel),
		getPageLink(first, firstRel),
		getPrevLink(page, last, first, prevRel),
		getNextLink(page, last, first, nextRel),
		getLastLink(last, lastRel),
	}

	links := make([]Link, 0)
	for _, f := range m {
		if link, err := f(u, limit); err == nil {
			links = append(links, *link)
		}
	}

	return links
}

func getNextLink(page, last, first int, rel string) getLinkFunc {
	return func(u *url.URL, limit int) (*Link, error) {
		if page < last && page >= first {
			return getLink(u, page+1, limit, rel)
		}

		return nil, errors.New("out of range")
	}
}

func getPrevLink(page, last, first int, rel string) getLinkFunc {
	return func(u *url.URL, limit int) (*Link, error) {
		if page > first && page <= last {
			return getLink(u, page-1, limit, rel)
		}

		return nil, errors.New("out of range")
	}
}

func getLastLink(last int, rel string) getLinkFunc {
	return func(u *url.URL, limit int) (*Link, error) {
		return getLink(u, last, limit, rel)
	}
}

func getPageLink(page int, rel string) getLinkFunc {
	return func(u *url.URL, limit int) (*Link, error) {
		return getLink(u, page, limit, rel)
	}
}

func getLink(u *url.URL, page, limit int, rel string) (*Link, error) {
	tmp := *u
	tmp.RawQuery = getRawQuery(u.Query(), page, limit)

	return &Link{Href: tmp.RequestURI(), Rel: rel}, nil
}

func getRawQuery(v url.Values, page, limit int) string {
	v.Set("page", strconv.Itoa(page))
	v.Set("limit", strconv.Itoa(limit))

	return v.Encode()
}
