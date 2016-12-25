package filemaker

import (
	"encoding/xml"
	"errors"
)

func (s *Server) GetRow(database, layout string, recid string) (Record, error) {
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + "&-recid=" + recid + "&-find"

	reader, err := s.getResult(query)
	defer reader.Close()
	if err != nil {
		return nil, err
	}

	fm := fmresultset{}
	if err := xml.NewDecoder(reader).Decode(&fm); err != nil {
		return nil, err
	}
	if fm.Error.Code != "0" {
		return nil, errors.New("Filemaker error " + fm.Error.Code)
	}
	return getRecordsFromXml(fm)[0], nil
}
