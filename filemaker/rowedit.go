package filemaker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
)

func (s *Server) EditRow(database, layout string, recid string, data map[string]string) (Record, error) {
	fmt.Println(data)

	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + "&-recid=" + recid + "&-edit"
	for k, v := range data {
		query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
	}

	buffer, err := s.getResult(query)
	if err != nil {
		return nil, err
	}
	fm := fmresultset{}
	if err := xml.NewDecoder(buffer).Decode(&fm); err != nil {
		return nil, err
	}
	if fm.Error.Code != "0" {
		return nil, errors.New("Filemaker error " + fm.Error.Code)
	}

	return s.GetRow(database, layout, recid)
}
