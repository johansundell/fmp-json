package filemaker

import (
	"fmt"
	"net/http"
)

func (s *Server) EditRow(database, layout string, recid string, data map[string]string) (Record, error) {
	fmt.Println(data)
	//return nil, nil
	client := &http.Client{}
	recidPart := "&-recid=" + recid
	searchType := "&-edit"
	url := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + recidPart + searchType
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

	return s.GetRow(database, layout, recid)
}
