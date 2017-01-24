package filemaker

import (
	"net/url"
	"strconv"
)

func (s *Server) Get(database, layout string, params []SearchParam, start, stop int) (Records, error) {
	search := ""
	for _, v := range params {
		search += v.String()
	}
	limit := "&-max=all"
	if start != -1 && stop != -1 {
		limit = "&-max=" + strconv.Itoa(stop-start) + "&-skip=" + strconv.Itoa(start)
	}

	searchType := "&-findall"
	if search != "" {
		searchType = "&-find"
	}
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + limit + search + searchType
	//log.Println(query)
	return s.getQueryResult(query)
}
