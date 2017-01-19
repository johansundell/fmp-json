package filemaker

import (
	"encoding/xml"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"
)

func (s *Server) EditRow(database, layout string, recid string, data map[string]string) (Record, error) {

	query := s.host + "/fmi/xml/fmresultset.xml?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + "&-recid=" + recid + "&-edit"
	/*for k, v := range data {
		fm := fmDef[k]
		if fm != nil {
			switch fm.Type {
			case FileMakerDate:
				if val, err := time.Parse(time.RFC3339, v); err == nil {
					query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(val.Format(fm.Format))
				} else {
					log.Println("Date format error", val)
				}
			default:
				query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
			}
		}
	}*/
	//log.Println(query)
	query, err := s.getFormatedDataQuery(query, database, layout, data)
	if err != nil {
		return nil, err
	}
	//log.Println(query)
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

func (s *Server) getFormatedDataQuery(query, database, layout string, data map[string]string) (string, error) {
	fmDef, err := s.getLayoutFields(database, layout)
	if err != nil {
		return "", err
	}
	//log.Println(fmDef)
	for k, v := range data {
		fm := fmDef[strings.ToLower(k)]
		if fm != nil {
			switch fm.Type {
			case FileMakerDate:
				if val, err := time.Parse(time.RFC3339, v); err == nil {
					query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(val.Format(fm.Format))
				} else {
					log.Println("Date format error", val)
				}
			default:
				query += "&" + url.QueryEscape(k) + "=" + url.QueryEscape(v)
			}
		}
	}
	return query, nil
}
