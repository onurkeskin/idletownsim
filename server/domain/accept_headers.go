// The MIT License (MIT)

// Copyright (c) 2015 Hafiz Ismail

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.


package domain

import (
	"regexp"
	"strconv"
	"strings"
)

type MediaTypeParams map[string]string
type MediaType struct {
	String     string          `json:"string"`
	Type       string          `json:"type"`
	Tree       string          `json:"tree"`
	SubType    string          `json:"subtype"`
	Suffix     string          `json:"suffix"`
	Parameters MediaTypeParams `json:"parameters"`
}
type AcceptHeader struct {
	MediaType     MediaType `json:"media_type"`
	QualityFactor float64   `json:"quality_factor"`
}

type AcceptHeaders []AcceptHeader

// mediaTypeRegExp match ((type)/(subtype)((+)(suffix))?)
var mediaTypeRegExp = regexp.MustCompile(`^([\w\*\-]+)\/([\w\*\.\-]+)((\+)(\w+))?`)

func NewAcceptHeadersFromString(str string) AcceptHeaders {

	var headers AcceptHeaders

	// parses `application/vnd.api+json;q=0.8;version=1.0` into MediaType type
	parseMediaType := func(str string) MediaType {

		var mediaType MediaType

		str = strings.Replace(str, " ", "", -1)
		tokens := strings.Split(str, ";")

		mediaType.String = tokens[0]

		// if params exists, parse params; else params is nil
		paramsTokens := tokens[1:]
		if len(paramsTokens) > 0 && paramsTokens[0] != "" {
			mediaType.Parameters = map[string]string{}
			for _, paramsToken := range paramsTokens {
				p := strings.Split(paramsToken, "=")
				if len(p) == 1 {
					mediaType.Parameters[p[0]] = ""
				}
				if len(p) > 1 {
					mediaType.Parameters[p[0]] = p[1]
				}
			}
		}

		// match ((type)/(subtype)((+)(suffix))?)
		match := mediaTypeRegExp.FindStringSubmatch(mediaType.String)
		if len(match) == 0 {
			return mediaType
		}

		// successful match results len() always 6
		mediaType.Type = match[1]
		mediaType.SubType = match[2]
		mediaType.Suffix = match[5]

		// parse [tree .] sub-type
		treeStr := strings.Split(mediaType.SubType, ".")
		if len(treeStr) > 1 && treeStr[0] != "" {
			mediaType.Tree = treeStr[0]
			mediaType.SubType = strings.Join(treeStr[1:], ".")
		}

		return mediaType
	}

	str = strings.Replace(str, " ", "", -1)
	mediaTypes := strings.Split(str, ",")
	for _, mediaTypeStr := range mediaTypes {
		mediaType := parseMediaType(mediaTypeStr)
		header := AcceptHeader{mediaType, 1}
		if len(mediaType.Parameters["q"]) > 0 {
			q, err := strconv.ParseFloat(mediaType.Parameters["q"], 64)
			if err != nil {
				q = 1 // default `q` value
			}
			header.QualityFactor = q
		}
		headers = append(headers, header)
	}
	return headers
}
