package rss

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalRSS(t *testing.T) {
	feed := Feed{
		Version: "2.0",
		Channels: []Channel{
			{
				Title:       "First channel",
				Link:        "https://example.com/first-channel",
				Description: "The first channel",
				Items: []Item{
					{
						GUID:        "1",
						PubDate:     Time(time.Date(2024, 12, 14, 12, 37, 0, 0, time.UTC)),
						Title:       "First item",
						Link:        "https://example.com/first-channel/first-item",
						Description: "The first item",
					},
					{
						GUID:        "2",
						PubDate:     Time(time.Date(2024, 12, 14, 12, 37, 0, 0, time.UTC)),
						Title:       "Second item",
						Link:        "https://example.com/first-channel/second-item",
						Description: "The second item",
					},
				},
			},
		},
	}

	expected := `<rss version="2.0">
	<channel>
		<title>First channel</title>
		<link>https://example.com/first-channel</link>
		<description>The first channel</description>
		<item>
			<guid>1</guid>
			<pubDate>Sat, 14 Dec 2024 12:37:00 UTC</pubDate>
			<title>First item</title>
			<link>https://example.com/first-channel/first-item</link>
			<description>The first item</description>
		</item>
		<item>
			<guid>2</guid>
			<pubDate>Sat, 14 Dec 2024 12:37:00 UTC</pubDate>
			<title>Second item</title>
			<link>https://example.com/first-channel/second-item</link>
			<description>The second item</description>
		</item>
	</channel>
</rss>`

	actual, err := xml.MarshalIndent(&feed, "", "\t")
	require.NoError(t, err)

	assert.Equal(t, expected, string(actual))
}

func TestUnmarshalXML(t *testing.T) {
	feed := `<rss version="2.0">
	<channel>
		<title>First channel</title>
		<link>https://example.com/first-channel</link>
		<description>The first channel</description>
		<item>
			<guid>1</guid>
			<pubDate>Sat, 14 Dec 2024 12:37:00 UTC</pubDate>
			<title>First item</title>
			<link>https://example.com/first-channel/first-item</link>
			<description>The first item</description>
		</item>
		<item>
			<guid>2</guid>
			<pubDate>Sat, 14 Dec 2024 12:37:00 UTC</pubDate>
			<title>Second item</title>
			<link>https://example.com/first-channel/second-item</link>
			<description>The second item</description>
		</item>
	</channel>
</rss>`

	expected := Feed{
		XMLName: xml.Name{
			Space: "",
			Local: "rss",
		},
		Version: "2.0",
		Channels: []Channel{
			{
				XMLName: xml.Name{
					Space: "",
					Local: "channel",
				},
				Title:       "First channel",
				Link:        "https://example.com/first-channel",
				Description: "The first channel",
				Items: []Item{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "item",
						},
						GUID:        "1",
						PubDate:     Time(time.Date(2024, 12, 14, 12, 37, 0, 0, time.UTC)),
						Title:       "First item",
						Link:        "https://example.com/first-channel/first-item",
						Description: "The first item",
					},
					{
						XMLName: xml.Name{
							Space: "",
							Local: "item",
						},
						GUID:        "2",
						PubDate:     Time(time.Date(2024, 12, 14, 12, 37, 0, 0, time.UTC)),
						Title:       "Second item",
						Link:        "https://example.com/first-channel/second-item",
						Description: "The second item",
					},
				},
			},
		},
	}

	var actual Feed
	err := xml.Unmarshal([]byte(feed), &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}
