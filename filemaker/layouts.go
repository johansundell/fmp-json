package filemaker

import (
	"encoding/xml"
	"errors"
	"net/url"
	"strings"
)

func (s *Server) Listlayouts(database string) (Records, error) {
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-layoutnames"
	return s.getQueryResult(query)
}

func (s *Server) getLayoutFields(database, layout string) (map[string]*FileMakerFieldInfo, error) {
	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + "&-view"
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
	fmDef := make(map[string]*FileMakerFieldInfo)
	for _, v := range fm.FieldDefinition {
		name := strings.ToLower(v.Name)
		switch v.Result {
		case "date":
			r := strings.NewReplacer("yyyy", "2006", "MM", "01", "dd", "02")
			dateLayout := r.Replace(fm.Datasource.DateFormat)
			fmDef[name] = &FileMakerFieldInfo{Type: FileMakerDate, Format: dateLayout}
		default:
			fmDef[name] = &FileMakerFieldInfo{Type: FileMakerString}
		}
	}
	return fmDef, nil
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
