# Periodize

[![GoDoc Reference](https://godoc.org/github.com/djcrock/periodize?status.svg)](http://godoc.org/github.com/djcrock/periodize)
[![Build Status](https://travis-ci.org/djcrock/periodize.svg?branch=master)](https://travis-ci.org/djcrock/periodize)
[![Go Report Card](https://goreportcard.com/badge/github.com/djcrock/periodize)](https://goreportcard.com/report/github.com/djcrock/periodize)

Go library to generate periodical .mobi files.

## Dependencies

Periodize requires that the [kindlegen](www.amazon.com/kindleformat/kindlegen) executable be available in `PATH`.

## Usage

```go
package main

import (
	"os"

	"github.com/djcrock/periodize"
)

func main() {
	iss := periodize.Issue{
		UniqueID:    "123",
		Title:       "My Periodical",
		Creator:     "djcrock",
		Publisher:   "djcrock",
		Subject:     "eBook Publishing",
		Description: "Demonstration of periodical publishing",
		Date:        "2017-10-21",
		Sections: []periodize.Section{
			{
				Title: "Section 1",
				Articles: []periodize.Article{
					{
						Title:   "Article 1-1",
						Author:  "djcrock",
						Content: "<body>Content 1</body>",
					},
				},
			},
			{
				Title: "Section 2",
				Articles: []periodize.Article{
					{
						Title:   "Article 2-1",
						Author:  "djcrock",
						Content: "<body>Content 2</body>",
					},
					{
						Title:   "Article 2-2",
						Author:  "djcrock",
						Content: "<body>Content 3</body>",
					},
				},
			},
		},
	}

	mobi, _ := os.Create("my_periodical.mobi")
	defer mobi.Close()

	// GenerateMobi accepts any io.Writer
	iss.GenerateMobi(mobi)
}
```
