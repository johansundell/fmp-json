package filemaker

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (s *Server) GetContainerField(database, layout, recid, field string) ([]byte, string, string, error) {

	record, err := s.GetRow(database, layout, recid)
	if err != nil {
		return nil, "", "", err
	}
	//log.Println("test", record)
	r := record[field]
	filename := ""
	if r != nil {
		//log.Println(r.v[13:strings.Index(r.v, "?")])
		filename = r.v[13:strings.Index(r.v, "?")]
	}
	//filename := "me.jpg"
	query := s.host + "/fmi/xml/cnt/" + filename + "?-db=" + url.QueryEscape(database) + "&-lay=" + url.QueryEscape(layout) + "&-recid=" + url.QueryEscape(recid) + "&-field=" + url.QueryEscape(field) + "(1)"
	//log.Println(query)
	client := &http.Client{}
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, "", "", err
	}
	req.SetBasicAuth(s.username, s.password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", "", nil
	}
	defer resp.Body.Close()
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", err
	}
	//ioutil.WriteFile("test.jpeg", buffer, 0664)
	contentType := ""
	if len(resp.Header["Content-Type"]) > 0 {
		contentType = resp.Header["Content-Type"][0]
	}

	return buffer, contentType, filename, nil
}
