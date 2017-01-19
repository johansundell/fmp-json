package filemaker

import (
	"encoding/xml"
	"errors"
	"net/url"
)

func (s *Server) NewRow(database, layout string, data map[string]string) (Record, error) {
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + "&-new"
	/*for k, v := range data {
		query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
	}*/
	query, err := s.getFormatedDataQuery(query, database, layout, data)
	if err != nil {
		return nil, err
	}

	buffer, err := s.getResult(query)
	if err != nil {
		return nil, err
	}
	defer buffer.Close()

	fm := fmresultset{}
	if err := xml.NewDecoder(buffer).Decode(&fm); err != nil {
		return nil, err
	}
	if fm.Error.Code != "0" {
		return nil, errors.New("Filemaker error " + fm.Error.Code)
	}

	return getRecordsFromXml(fm)[0], nil
}
