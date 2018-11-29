package filemaker

import (
	"net/url"
)

func (s *Server) RunScript(database, layout, script, param string) (Records, int, error) {
	if param != "" {
		param = "&-script.param=" + url.QueryEscape(param)
	}
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + "&-findany" + "&-script=" + url.QueryEscape(script) + param
	//log.Println(query)
	return s.getQueryResult(query)
}
