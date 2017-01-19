package filemaker

import (
	"encoding/xml"
	"errors"
	"net/url"
)

func (s *Server) Listlayouts(database string) (Records, error) {
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-layoutnames"
	return s.getQueryResult(query)
}

func (s *Server) getQueryResult(query string) (Records, error) {
	reader, err := s.getResult(query)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	fm := fmresultset{}
	if err := xml.NewDecoder(reader).Decode(&fm); err != nil {
		return nil, err
	}

	if fm.Error.Code != "0" {
		return nil, errors.New("Filemaker error " + fm.Error.Code)
	}
	return getRecordsFromXml(fm), nil
}
