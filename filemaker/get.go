package filemaker

func (s *Server) Get(database, layout string, params []SearchParam) (Records, error) {
	search := ""
	for _, v := range params {
		search += v.String()
	}

	searchType := "&-findall"
	if search != "" {
		searchType = "&-find"
	}
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + search + searchType
	return s.getQueryResult(query)
}
