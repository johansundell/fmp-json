package filemaker

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type fmresultset struct {
	Version         string            `xml:"version,attr"`
	Error           fmError           `xml:"error"`
	FieldDefinition []fieldDefinition `xml:"metadata>field-definition"`
	Resultset       resultset         `xml:"resultset"`
	Xmlns           string            `xml:"xmlns,attr"`
	Product         product           `xml:"product"`
	Datasource      datasource        `xml:"datasource"`
}

type record struct {
	RecordId string  `xml:"record-id,attr"`
	ModId    string  `xml:"mod-id,attr"`
	Field    []field `xml:"field"`
}
type field struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"data"`
}
type fmError struct {
	Code string `xml:"code,attr"`
}
type product struct {
	Build   string `xml:"build,attr"`
	Name    string `xml:"name,attr"`
	Version string `xml:"version,attr"`
}
type datasource struct {
	Layout          string `xml:"layout,attr"`
	Table           string `xml:"table,attr"`
	TimeFormat      string `xml:"time-format,attr"`
	TimestampFormat string `xml:"timestamp-format,attr"`
	TotalCount      string `xml:"total-count,attr"`
	Database        string `xml:"database,attr"`
	DateFormat      string `xml:"date-format,attr"`
}
type fieldDefinition struct {
	Type          string `xml:"type,attr"`
	AutoEnter     string `xml:"auto-enter,attr"`
	MaxRepeat     string `xml:"max-repeat,attr"`
	TimeOfDay     string `xml:"time-of-day,attr"`
	Global        string `xml:"global,attr"`
	NotEmpty      string `xml:"not-empty,attr"`
	NumericOnly   string `xml:"numeric-only,attr"`
	Result        string `xml:"result,attr"`
	FourDigitYear string `xml:"four-digit-year,attr"`
	Name          string `xml:"name,attr"`
}
type resultset struct {
	Count     string   `xml:"count,attr"`
	FetchSize string   `xml:"fetch-size,attr"`
	Record    []record `xml:"record"`
}

type Server struct {
	host, username, password string
}

type Records []Record

type Record map[string]*internal

func (r Record) Add(key string, value string) {
	r[key] = &internal{v: value}
}

func (r internal) String() string {
	return r.v
}

type FileMakerType int

const (
	Unknown FileMakerType = iota
	FileMakerNumber
	FileMakerString
	FileMakerDate
	FileMakerTimestamp
	FileMakerContainer
)

type FileMakerFieldInfo struct {
	Type   FileMakerType
	Format string
}

type internal struct {
	v      string
	Type   FileMakerType
	Format string
}

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

func getRecordsFromXml(fm fmresultset) Records {
	test := make(map[string]fieldDefinition)
	for _, v := range fm.FieldDefinition {
		test[v.Name] = v
	}
	records := make(Records, 0)
	for _, v := range fm.Resultset.Record {
		row := make(Record)
		row["recid"] = &internal{v: v.RecordId, Type: FileMakerNumber}
		for _, r := range v.Field {
			row[r.Name] = &internal{v: r.Value}
			switch test[r.Name].Result {
			case "number":
				row[r.Name].Type = FileMakerNumber
			case "text":
				row[r.Name].Type = FileMakerString
			case "date":
				row[r.Name].Type = FileMakerDate
				replacer := strings.NewReplacer("yyyy", "2006", "MM", "01", "dd", "02")
				dateLayout := replacer.Replace(fm.Datasource.DateFormat)
				row[r.Name].Format = dateLayout
			case "timestamp":
				row[r.Name].Type = FileMakerTimestamp
				replacer := strings.NewReplacer("yyyy", "2006", "MM", "01", "dd", "02", "HH", "15", "mm", "04", "ss", "05")
				dateLayout := replacer.Replace(fm.Datasource.TimestampFormat)
				row[r.Name].Format = dateLayout
			case "container":
				if r.Value != "" {
					row[r.Name].Type = FileMakerContainer
				}
			}
			//log.Println(test[r.Name].Result, row[r.Name])
		}
		records = append(records, row)
	}
	return records
}

func (s *Server) getResult(query string) (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.username, s.password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
