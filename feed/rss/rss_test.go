// -----------------------------------------------------------------------------
// Copyright (c) 2025-present Detlef Stern
//
// This file is part of webs.
//
// webs is licensed under the latest version of the EUPL // (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2025-present Detlef Stern
// -----------------------------------------------------------------------------

package rss_test

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"t73f.de/r/webs/feed/rss"
)

func TestSimpleRSS(t *testing.T) {
	feed := rss.Feed{
		XMLName:        xml.Name{},
		Title:          "RSS Test",
		Link:           "https://r.t73f.de/webs/dir?ci=tip&name=feed",
		Description:    "Test Feed",
		Language:       "de",
		Copyright:      "none",
		ManagingEditor: "detlef@example.com",
		WebMaster:      "stern@example.com",
		PubDate:        rss.RFC822Date(time.Date(2025, time.January, 5, 16, 46, 17, 0, time.UTC)),
		LastBuildDate:  rss.RFC822Date(time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)),
		Generator:      "TestDriver",
		TTL:            60,
		Image:          nil,
		Items: []*rss.Item{
			{
				Title:       "Item One",
				Description: rss.CData{},
				Author:      "ds@example.com",
				Category:    []string{"test"},
				Link:        "https://example.com/one",
				GUID: &rss.GUID{
					IsPermaLink: false,
					Value:       "bla fasel",
				},
				PubDate: rss.RFC822Date(time.Date(2025, time.July, 15, 12, 12, 12, 12, time.UTC)),
				Source: &rss.Source{
					URL:   "http://example.com",
					Title: "Source",
				},
			},
		},
	}

	var sb strings.Builder
	err := feed.Write(&sb)
	if err != nil {
		t.Fatal(err)
	}
	got := sb.String()
	exp := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>RSS Test</title>
    <link>https://r.t73f.de/webs/dir?ci=tip&amp;name=feed</link>
    <description>Test Feed</description>
    <language>de</language>
    <copyright>none</copyright>
    <managingEditor>detlef@example.com</managingEditor>
    <webMaster>stern@example.com</webMaster>
    <pubDate>Sun, 05 Jan 2025 16:46:17 +0000</pubDate>
    <lastBuildDate>Thu, 01 Jan 1970 00:00:00 +0000</lastBuildDate>
    <generator>TestDriver</generator>
    <ttl>60</ttl>
    <item>
      <title>Item One</title>
      <description></description>
      <author>ds@example.com</author>
      <category>test</category>
      <link>https://example.com/one</link>
      <guid isPermaLink="false">bla fasel</guid>
      <pubDate>Tue, 15 Jul 2025 12:12:12 +0000</pubDate>
      <source url="http://example.com">Source</source>
    </item>
  </channel>
</rss>`
	if got != exp {
		t.Errorf("EXP: %s\nGOT: %s", exp, got)
	}
}
