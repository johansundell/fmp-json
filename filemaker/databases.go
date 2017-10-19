package filemaker

func (s *Server) ListDatabases() (Records, int, error) {
	query := s.host + "/fmi/xml/fmresultset.xml?-dbnames"
	return s.getQueryResult(query)
}
