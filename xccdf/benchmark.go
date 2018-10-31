package xccdf

import (
	"encoding/xml"
)

type Benchmark struct {
	XMLName     xml.Name  `xml:"Benchmark"`
	Id          string    `xml:"id,attr"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	ReleaseInfo string    `xml:"plain-text"`
	Version     string    `xml:"version"`
	Profiles    []Profile `xml:"Profile"`
	Groups      []Group   `xml:"Group"`
}

type Profile struct {
	XMLName xml.Name `xml:"Profile"`
	Id      string   `xml:"id,attr"`
	Title   string   `xml:"title"`
	Selects []Select `xml:"select"`
}

type Select struct {
	XMLName  xml.Name `xml:"select"`
	IdRef    string   `xml:"idref,attr"`
	Selected string   `xml:"selected,attr"`
}

type Group struct {
	XMLName xml.Name `xml:"Group"`
	Id      string   `xml:"id,attr"`
	Title   string   `xml:"title"`
	Rules   []Rule   `xml:"Rule"`
}

type Rule struct {
	XMLName      xml.Name        `xml:"Rule"`
	Id           string          `xml:"id,attr"`
	SeverityAttr string          `xml:"severity,attr"`
	Version      string          `xml:"version"`
	Title        string          `xml:"title"`
	Description  RuleDescription `xml:"description"`
	FixText      string          `xml:"fixtext"`
	CheckContent string          `xml:"check>check-content"`
}

func (r *Rule) Severity() string {
	switch r.SeverityAttr {
	case "low":
		return "CAT III"
	case "medium":
		return "CAT II"
	case "high":
		return "CAT I"
	}
	return "UNKNOWN"
}

type RuleDescription struct {
	XMLName xml.Name `xml:"description"`
	XML     string   `xml:",chardata"`
	attr    struct {
		VulnDiscussion string `xml:"VulnDiscussion"`
	}
}

func (r *RuleDescription) Discussion() string {
	if r.attr.VulnDiscussion == "" {
		xml.Unmarshal([]byte("<xml>"+r.XML+"</xml>"), &r.attr)
	}
	return r.attr.VulnDiscussion
}
