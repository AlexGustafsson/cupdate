package rss

import (
	"encoding/xml"
	"time"
)

var rfc2822 = "Mon, 02 Jan 2006 15:04:05 MST"

// Feed is an RSS feed item.
type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`

	Channels []Channel `xml:"channel"`
}

// Channel is an RSS Channel item.
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Items       []Item   `xml:"item"`
}

// Item is a feed item.
type Item struct {
	XMLName     xml.Name `xml:"item"`
	GUID        string   `xml:"guid"`
	PubDate     Time     `xml:"pubDate"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

// Time represents a RFC2822 time, as used by RSS.
type Time time.Time

func (t Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(time.Time(t).Format(rfc2822), start)
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string
	err := d.DecodeElement(&value, &start)
	if err != nil {
		return err
	}

	time, err := time.Parse(rfc2822, value)
	if err != nil {
		return err
	}

	*t = Time(time)
	return nil
}
