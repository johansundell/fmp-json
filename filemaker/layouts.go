package filemaker

import (
	"encoding/xml"
	"errors"
	"net/url"
	"strconv"
	"strings"
)

func (s *Server) Listlayouts(database string) (Records, int, error) {
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
		case "timestamp":
			r := strings.NewReplacer("yyyy", "2006", "MM", "01", "dd", "02", "HH", "15", "mm", "04", "ss", "05")
			dateLayout := r.Replace(fm.Datasource.TimestampFormat)
			fmDef[name] = &FileMakerFieldInfo{Type: FileMakerDate, Format: dateLayout}
		default:
			fmDef[name] = &FileMakerFieldInfo{Type: FileMakerString}
		}
	}
	return fmDef, nil
}

func (s *Server) getQueryResult(query string) (Records, int, error) {
	reader, err := s.getResult(query)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	fm := fmresultset{}
	if err := xml.NewDecoder(reader).Decode(&fm); err != nil {
		return nil, 0, err
	}

	if fm.Error.Code != "0" {
		return nil, 0, errors.New("Filemaker error " + fm.Error.Code)
	}
	found, err := strconv.Atoi(fm.Resultset.Count)
	if err != nil {
		return nil, 0, err
	}
	return getRecordsFromXml(fm), found, nil
}
