package filemaker

import "encoding/xml"

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
	reader, err := s.getResult(query)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	fm := fmresultset{}
	if err := xml.NewDecoder(reader).Decode(&fm); err != nil {
		return nil, err
	}

	return getRecordsFromXml(fm), nil
}
