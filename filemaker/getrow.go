package filemaker

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
)

func (s *Server) GetRow(database, layout string, recid string) (Record, error) {
	client := &http.Client{}
	search := "&-recid=" + recid
	searchType := "&-find"
	url := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + search + searchType
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.username, s.password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fm := fmresultset{}
	if err := xml.NewDecoder(resp.Body).Decode(&fm); err != nil {
		return nil, err
	}
	if fm.Error.Code != "0" {
		return nil, errors.New("Filemaker error " + fm.Error.Code)
	}

	record := make(map[string]string)
	for _, v := range fm.Resultset.Record {
		row := make(map[string]string)
		row["recid"] = v.RecordId
		for _, r := range v.Field {
			row[r.Name] = r.Value
		}
		record = row
	}
	return record, nil
}
