package filemaker

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
	//Field    [][]Field `xml:"field"`
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
