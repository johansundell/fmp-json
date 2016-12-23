package filemaker

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
)

type Server struct {
	host, username, password string
}

type Records []map[string]string

//type Records []Record
type Record map[string]string

type SearchOperator string

const (
	MoreThan SearchOperator = "gt"
	LessThan SearchOperator = "lt"
	Equal    SearchOperator = "eq"
)

type SearchParam struct {
	Name, Value string
	Op          SearchOperator
}

func (sp SearchParam) String() string {
	return "&" + url.QueryEscape(sp.Name) + "=" + url.QueryEscape(sp.Value) + "&" + url.QueryEscape(sp.Name) + ".op=" + string(sp.Op)
}

func NewServer(host, username, password string) Server {
	return Server{host: host, username: username, password: password}
}

func (s *Server) Get(database, layout string, params []SearchParam) (Records, error) {
	client := &http.Client{}
	search := ""
	for _, v := range params {
		search += v.String()
	}
	//fmt.Println(search)
	// "http://fmh-iwp12.no-ip.info/fmi/xml/fmresultset.xml?-db=FMServer_Sample&-lay=test&-findall"
	searchType := "&-findall"
	if search != "" {
		searchType = "&-find"
	}
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

	records := make([]map[string]string, 0)
	for _, v := range fm.Resultset.Record {
		row := make(map[string]string)
		row["recid"] = v.RecordId
		for _, r := range v.Field {
			row[r.Name] = r.Value
		}
		records = append(records, row)
	}
	return records, nil
}
