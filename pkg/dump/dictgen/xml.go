package dictgen

type XMLDoc struct {
	Version    string          `xml:"version,attr"`
	Messages   []*XMLMessage   `xml:"message"`
	Components []*XMLComponent `xml:"component"`
	Fields     []*XMLField     `xml:"field"`
	Enums      []*XMLEnum      `xml:"enum"`
}

type XMLElement struct {
	FIXML     string `xml:"fixml,attr"`
	Name      string `xml:"name,attr"`
	NotReqXML int    `xml:"not_req_xml,attr"`
}

type XMLMessage struct {
	XMLElement
	FIX        string             `xml:"fix,attr"`
	Cat        string             `xml:"cat,attr"`
	Sec        string             `xml:"sec,attr"`
	Components []*XMLComponentRef `xml:"component"`
	Fields     []*XMLFieldRef     `xml:"field"`
}

type XMLComponent struct {
	XMLElement
	Cat  string `xml:"cat,attr"`
	Type string `xml:"aid,type"`
}

type XMLField struct {
	XMLElement
	AID   int       `xml:"aid,attr"`
	Type  string    `xml:"type,attr"`
	EID   int       `xml:"eid,attr"`
	Alias *XMLAlias `xml:"alias"`
}

type XMLAlias struct {
	FIXML string `xml:"fixml,attr"`
	Cat   string `xml:"cat,attr"`
}

type XMLFieldRef struct {
	AID int `xml:"aid,attr"`
}

type XMLComponentRef struct {
	Name string `xml:"name,attr"`
}

type XMLEnum struct {
	EID   int        `xml:"eid,attr"`
	Items []*XMLItem `xml:"item"`
}

type XMLItem struct {
	Value string `xml:"value,attr"`
	Text  string `xml:",chardata"`
}
