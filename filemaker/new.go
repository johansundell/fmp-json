package filemaker

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

func (s *Server) NewRow(database, layout string, data map[string]string) (Record, error) {
	fmt.Println(data)
	//return nil, nil
	client := &http.Client{}
	//recidPart := "&-recid=" + recid
	searchType := "&-new"
	url := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + searchType
	for k, v := range data {
		url += "&" + k + "=" + v
	}
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

	record := make(map[string]string)
	for _, v := range fm.Resultset.Record {
		row := make(map[string]string)
		row["recid"] = v.RecordId
		for _, r := range v.Field {
			row[r.Name] = r.Value
		}
		record = row
	}
	//fmt.Println(record)
	return record, nil
	//return s.GetRow(database, layout, recid)
}
