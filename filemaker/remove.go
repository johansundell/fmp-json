package filemaker

import (
	"fmt"
	"net/http"
)

func (s *Server) DeleteRow(database, layout string, recid string) error {
	client := &http.Client{}
	recidPart := "&-recid=" + recid
	searchType := "&-delete"
	url := s.host + "/fmi/xml/fmresultset.xml?-db=" + database + "&-lay=" + layout + recidPart + searchType

	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(s.username, s.password)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
