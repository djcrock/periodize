package periodize

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestGenerateMobi(t *testing.T) {
	cover, err := os.Open("test_samples/cover-image.gif")
	if err != nil {
		t.Errorf("Failed to open test cover image: %v", err)
	}
	iss := Issue{
		UniqueID:    "TestID",
		Title:       "Test Title",
		Creator:     "Test Creator",
		Publisher:   "Test Publisher",
		Subject:     "Test Subject",
		Description: "Test Description",
		Date:        "2017-10-17",
		CoverImage:  cover,
		Sections: []Section{
			{
				Title: "Test Section 1",
				Articles: []Article{
					{Title: "Test Article 1-1", Author: "djcrock1", Content: "<body>Content 1</body>"},
					{Title: "Test Article 1-2", Author: "djcrock2", Content: "<body>Content 2</body>"},
				},
			},
			{
				Title: "Test Section 2",
				Articles: []Article{
					{Title: "Test Article 2-1", Author: "djcrock3", Content: "<body>Content 3</body>"},
					{Title: "Test Article 2-2", Author: "djcrock4", Content: "<body>Content 4</body>"},
				},
			},
		},
	}

	mobi := new(bytes.Buffer)
	err = iss.GenerateMobi(mobi)
	if err != nil {
		t.Errorf("Failed to generate MOBI: %v", err)
	}
}

func TestGenerateOpf(t *testing.T) {
	opf := new(bytes.Buffer)
	iss := Issue{
		UniqueID:    "TestID",
		Title:       "Test Title",
		Creator:     "Test Creator",
		Publisher:   "Test Publisher",
		Subject:     "Test Subject",
		Description: "Test Description",
		Date:        "2017-10-17",
		Sections: []Section{
			{
				Title: "Test Section 1",
				Articles: []Article{
					{Title: "Test Article 1-1"},
					{Title: "Test Article 1-2"},
				},
			},
			{
				Title: "Test Section 2",
				Articles: []Article{
					{Title: "Test Article 2-1"},
					{Title: "Test Article 2-2"},
				},
			},
		},
	}

	iss.prepare()

	if err := iss.generateOpf(opf); err != nil {
		t.Errorf("Error occurred while generating OPF: %v", err)
	}

	opfExpected, err := ioutil.ReadFile("test_samples/content.opf")
	if err != nil {
		t.Errorf("Error occurred while reading sample OPF: %v", err)
	}

	if !bytes.Equal(opf.Bytes(), opfExpected) {
		t.Errorf("Actual OPF did not match expected \nActual\n%s\n\n===\nExpected\n%s", opf.String(), string(opfExpected))
	}
}

func TestGenerateTableOfContents(t *testing.T) {
	contents := new(bytes.Buffer)
	iss := Issue{
		UniqueID:    "TestID",
		Title:       "Test Title",
		Creator:     "Test Creator",
		Publisher:   "Test Publisher",
		Subject:     "Test Subject",
		Description: "Test Description",
		Date:        "2017-10-17",
		Sections: []Section{
			{
				Title: "Test Section 1",
				Articles: []Article{
					{Title: "Test Article 1-1"},
					{Title: "Test Article 1-2"},
				},
			},
			{
				Title: "Test Section 2",
				Articles: []Article{
					{Title: "Test Article 2-1"},
					{Title: "Test Article 2-2"},
				},
			},
		},
	}

	iss.prepare()

	if err := iss.generateTableOfContents(contents); err != nil {
		t.Errorf("Error occurred while generating table of contents: %v", err)
	}

	contentsExpected, err := ioutil.ReadFile("test_samples/contents.html")
	if err != nil {
		t.Errorf("Error occurred while reading sample table of contents: %v", err)
	}

	if !bytes.Equal(contents.Bytes(), contentsExpected) {
		t.Errorf(
			"Actual table of contents did not match expected \nActual\n%s\n\n===\nExpected\n%s",
			contents.String(),
			string(contentsExpected),
		)
	}
}

func TestGenerateNav(t *testing.T) {
	nav := new(bytes.Buffer)
	iss := Issue{
		UniqueID:    "TestID",
		Title:       "Test Title",
		Creator:     "Test Creator",
		Publisher:   "Test Publisher",
		Subject:     "Test Subject",
		Description: "Test Description",
		Date:        "2017-10-17",
		Sections: []Section{
			{
				Title: "Test Section 1",
				Articles: []Article{
					{Title: "Test Article 1-1", Author: "djcrock1"},
					{Title: "Test Article 1-2", Author: "djcrock2"},
				},
			},
			{
				Title: "Test Section 2",
				Articles: []Article{
					{Title: "Test Article 2-1", Author: "djcrock3"},
					{Title: "Test Article 2-2", Author: "djcrock4"},
				},
			},
		},
	}

	iss.prepare()

	if err := iss.generateNav(nav); err != nil {
		t.Errorf("Error occurred while generating nav: %v", err)
	}

	navExpected, err := ioutil.ReadFile("test_samples/nav-contents.ncx")
	if err != nil {
		t.Errorf("Error occurred while reading sample nav: %v", err)
	}

	if !bytes.Equal(nav.Bytes(), navExpected) {
		t.Errorf(
			"Actual nav did not match expected \nActual\n%s\n\n===\nExpected\n%s",
			nav.String(),
			string(navExpected),
		)
	}
}

func TestGenerateArticle(t *testing.T) {
	artOut := new(bytes.Buffer)
	art := Article{
		Content: "<body>Test Content</body>",
	}
	art.generateContent(artOut)
	if artOut.String() != art.Content {
		t.Errorf(
			"Actual article content did not match expected \nActual\n%s\n\n====\nExpected\n%s",
			artOut.String(),
			art.Content,
		)
	}
}

func TestPad(t *testing.T) {
	expected := []struct {
		number int
		padded string
	}{
		{1, "000001"},
		{12, "000012"},
		{123, "000123"},
		{1234, "001234"},
		{12345, "012345"},
		{123456, "123456"},
	}

	for _, exp := range expected {
		if result := pad(exp.number); result != exp.padded {
			t.Errorf("pad(%d) returned %s, expected %s", exp.number, result, exp.padded)
		}
	}
}

func TestPrepareIssue(t *testing.T) {
	iss := Issue{
		Sections: []Section{
			{
				Title: "Test Section 1",
				Articles: []Article{
					{Title: "Test Article 1-1"},
					{Title: "Test Article 1-2"},
				},
			},
			{
				Title: "Test Section 2",
				Articles: []Article{
					{Title: "Test Article 2-1"},
					{Title: "Test Article 2-2"},
				},
			},
		},
	}

	iss.prepare()

	if iss.UniqueID == "" {
		t.Error("UniqueID must not be empty")
	}

	verifySectionID(t, iss.Sections[0], "section-0")
	verifySectionID(t, iss.Sections[1], "section-1")
	verifySectionPlayOrder(t, iss.Sections[0], iss.Sections[0].Articles[0].PlayOrder)
	verifySectionPlayOrder(t, iss.Sections[1], iss.Sections[1].Articles[0].PlayOrder)
	verifySectionHref(t, iss.Sections[0], iss.Sections[0].Articles[0].Href)
	verifySectionHref(t, iss.Sections[1], iss.Sections[1].Articles[0].Href)

	verifyArticlePlayOrder(t, iss.Sections[0].Articles[0], 0)
	verifyArticlePlayOrder(t, iss.Sections[0].Articles[1], 1)
	verifyArticlePlayOrder(t, iss.Sections[1].Articles[0], 2)
	verifyArticlePlayOrder(t, iss.Sections[1].Articles[1], 3)
	verifyArticleHref(t, iss.Sections[0].Articles[0], "000000.html")
	verifyArticleHref(t, iss.Sections[0].Articles[1], "000001.html")
	verifyArticleHref(t, iss.Sections[1].Articles[0], "000002.html")
	verifyArticleHref(t, iss.Sections[1].Articles[1], "000003.html")

}

func TestPrepareSection(t *testing.T) {
	sect := Section{
		Title: "Test Section",
		Articles: []Article{
			{
				Title:   "Test Title 1",
				Author:  "Test Author 1",
				Content: "Test Content 1",
			},
			{
				Title:   "Test Title 2",
				Author:  "Test Author 2",
				Content: "Test Content 2",
			},
		},
	}

	sect.prepare(2, 123)
	verifySectionID(t, sect, "section-2")
	verifySectionPlayOrder(t, sect, 123)
	verifySectionHref(t, sect, "000123.html")

	verifyArticlePlayOrder(t, sect.Articles[0], 123)
	verifyArticleHref(t, sect.Articles[0], "000123.html")

	verifyArticlePlayOrder(t, sect.Articles[1], 124)
	verifyArticleHref(t, sect.Articles[1], "000124.html")
}

func TestPrepareArticle(t *testing.T) {
	art := Article{
		Title:   "Test Title",
		Author:  "Test Author",
		Content: "Test Content",
	}
	art.prepare(0)
	verifyArticlePlayOrder(t, art, 0)
	verifyArticleHref(t, art, "000000.html")

	art.prepare(123)
	verifyArticlePlayOrder(t, art, 123)
	verifyArticleHref(t, art, "000123.html")
}

func verifySectionID(t *testing.T, sect Section, sectionID string) {
	if sect.SectionID != sectionID {
		t.Errorf("Section SectionID %v did not match expected value %v", sect.SectionID, sectionID)
	}
}

func verifySectionPlayOrder(t *testing.T, sect Section, playOrder int) {
	if sect.PlayOrder != playOrder {
		t.Errorf("Section PlayOrder %v did not match expected value %v", sect.PlayOrder, playOrder)
	}
}

func verifySectionHref(t *testing.T, sect Section, href string) {
	if sect.Href != href {
		t.Errorf("Section Href %v did not match expected value %v", sect.Href, href)
	}
}

func verifyArticlePlayOrder(t *testing.T, art Article, playOrder int) {
	if art.PlayOrder != playOrder {
		t.Errorf("Article PlayOrder %v did not match expected value %v", art.PlayOrder, playOrder)
	}
}

func verifyArticleHref(t *testing.T, art Article, href string) {
	if art.Href != href {
		t.Errorf("Article Href %v did not match expected value %v", art.Href, href)
	}
}
