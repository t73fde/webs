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

// Package rss assists in building a RSS 2.0 feed.
//
// Based on RSS 2.0.11 standard: https://www.rssboard.org/rss-specification
//
// Currently, not all channel and item elements are supported.
package rss

import (
	"encoding/xml"
	"io"
	"time"
)

type header struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Feed    *Feed
}

// Feed is the main structure for a RSS feed.
type Feed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Description    string   `xml:"description"`
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"`
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	TTL            int      `xml:"ttl,omitempty"`
	Image          *Image
	Items          []*Item
}

// Image is the structure of an image that can be displayed with the feed.
type Image struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

// Item is the structure of a feed item.
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description CData    `xml:"description"`
	Author      string   `xml:"author,omitempty"`
	Category    []string `xml:"category"`
	Link        string   `xml:"link"`
	GUID        *GUID
	PubDate     string `xml:"pubDate"`
	Source      *Source
}

// GUID is a string that uniquely identifies an item.
// It may be a URL to the item that can be opened in a web browser (permalink).
type GUID struct {
	XMLName     xml.Name `xml:"guid"`
	IsPermaLink bool     `xml:"isPermaLink,attr"`
	Value       string   `xml:",chardata"`
}

// Source is the name and URL of the RSS feed, where the item originates.
type Source struct {
	XMLName xml.Name `xml:"source"`
	URL     string   `xml:"url,attr"`
	Title   string   `xml:",chardata"`
}

// CData is a helper structure to tell a RSS parser that it should not analyze the string.
type CData struct {
	Data string `xml:",cdata"`
}

// RFC822Date returns the time as a RFC822 encoded string.
func RFC822Date(t time.Time) string {
	// The W3C feed validator complains on Go's RFC822 encoder.
	// The RFC1123 encoder seems to legible.
	return t.Format(time.RFC1123Z)
}

// Write the feed as XML.
func (rss *Feed) Write(w io.Writer) error {
	hd := header{Version: "2.0", Feed: rss}
	_, err := io.WriteString(w, xml.Header)
	if err == nil {
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		err = enc.Encode(hd)
	}
	return err
}
