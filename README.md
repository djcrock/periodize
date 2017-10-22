# Periodize
Go library to generate periodical .mobi files.

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

	iss.GenerateMobi(mobi)
}
```
