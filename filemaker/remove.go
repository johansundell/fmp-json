package filemaker

import (
	"encoding/xml"
	"errors"
)

func (s *Server) DeleteRow(database, layout string, recid string) error {
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + "&-recid=" + recid + "&-delete"
	reader, err := s.getResult(query)
	if err != nil {
		return err
	}
	defer reader.Close()

	fm := fmresultset{}
	if err := xml.NewDecoder(reader).Decode(&fm); err != nil {
		return err
	}
	if fm.Error.Code != "0" {
		return errors.New("Filemaker error " + fm.Error.Code)
	}

	return nil
}
