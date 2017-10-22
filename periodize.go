package periodize

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"image"
	"image/gif"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strconv"
)

// An Issue of a periodical
type Issue struct {
	UniqueID    string
	Title       string
	Creator     string
	Publisher   string
	Subject     string
	Description string
	Date        string
	CoverImage  io.Reader
	Sections    []Section
}

// A Section of an Issue
type Section struct {
	Title    string
	Articles []Article
	// Set automatically by Issue.prepare()
	SectionID string
	PlayOrder int
	Href      string
}

// An Article appearing in an Issue of a periodical
type Article struct {
	Title   string
	Author  string
	Content string
	// Set automatically by Issue.prepare()
	PlayOrder int
	Href      string
}

const opfFilename = "content.opf"
const contentsFilename = "contents.html"
const navFilename = "nav-contents.ncx"
const coverFilename = "cover-image.gif"
const mobiFilename = "content.mobi"

var opfTemplate = template.Must(template.New("opf").Parse(opfTemplateString))
var contentsTemplate = template.Must(template.New("contents").Parse(contentsTemplateString))
var navTemplate = template.Must(template.New("nav").Parse(navTemplateString))

// GenerateMobi writes a .mobi format periodical file for an Issue to the provided Writer
func (iss *Issue) GenerateMobi(wr io.Writer) error {
	iss.prepare()
	dir, err := ioutil.TempDir(os.TempDir(), "Periodize_"+iss.UniqueID+"_")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	opfFile, err := os.Create(path.Join(dir, opfFilename))
	if err != nil {
		return fmt.Errorf("Failed to create OPF: %v", err)
	}
	defer opfFile.Close()

	err = iss.generateOpf(opfFile)
	if err != nil {
		return fmt.Errorf("Failed to generate OPF: %v", err)
	}

	contentsFile, err := os.Create(path.Join(dir, contentsFilename))
	if err != nil {
		return fmt.Errorf("Failed to create table of contents: %v", err)
	}
	defer contentsFile.Close()

	err = iss.generateTableOfContents(contentsFile)
	if err != nil {
		return fmt.Errorf("Failed to generate table of contents: %v", err)
	}

	navFile, err := os.Create(path.Join(dir, navFilename))
	if err != nil {
		return fmt.Errorf("Failed to create nav: %v", err)
	}
	defer navFile.Close()

	err = iss.generateNav(navFile)
	if err != nil {
		return fmt.Errorf("Failed to generate nav: %v", err)
	}

	for _, sect := range iss.Sections {
		for _, art := range sect.Articles {
			artFile, err := os.Create(path.Join(dir, art.Href))
			if err != nil {
				return fmt.Errorf("Failed to write article content: %v", err)
			}
			art.generateContent(artFile)
			err = artFile.Close()
			if err != nil {
				return fmt.Errorf("Failed to close article content: %f", err)
			}
		}
	}

	cover, err := os.Create(path.Join(dir, coverFilename))
	if err != nil {
		return fmt.Errorf("Failed to open cover image to write: %v", err)
	}
	defer cover.Close()

	_, err = io.Copy(cover, iss.CoverImage)
	if err != nil {
		return fmt.Errorf("Failed to copy cover image: %v", err)
	}

	cmdOut := new(bytes.Buffer)
	cmd := exec.Command("kindlegen", opfFilename, "-o", mobiFilename)
	cmd.Dir = dir
	cmd.Stdout = cmdOut
	cmd.Stderr = cmdOut
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("kindlegen returned an error: %v, %s", err, cmdOut.String())
	}

	mobi, err := os.Open(path.Join(dir, mobiFilename))
	if err != nil {
		return fmt.Errorf("Failed to open generated MOBI: %v", err)
	}
	defer mobi.Close()

	_, err = io.Copy(wr, mobi)
	if err != nil {
		return fmt.Errorf("Failed to write generated MOBI: %v", err)
	}

	return nil
}

func (iss *Issue) generateOpf(wr io.Writer) error {
	// html/template doesn't like templates with the XML header, so prepend it
	wr.Write([]byte(xml.Header))
	return opfTemplate.Execute(wr, iss)
}

func (iss *Issue) generateTableOfContents(wr io.Writer) error {
	return contentsTemplate.Execute(wr, iss)
}

func (iss *Issue) generateNav(wr io.Writer) error {
	// html/template doesn't like templates with the XML header, so prepend it
	wr.Write([]byte(xml.Header))
	return navTemplate.Execute(wr, iss)
}

func (art *Article) generateContent(wr io.Writer) {
	wr.Write([]byte(art.Content))
}

func (iss *Issue) prepare() {
	if iss.UniqueID == "" {
		iss.UniqueID = strconv.Itoa(rand.Int())
	}

	if iss.CoverImage == nil {
		iss.CoverImage = createBlankCover()
	}

	playOrder := 0
	for i := range iss.Sections {
		iss.Sections[i].prepare(i, playOrder)
		playOrder += len(iss.Sections[i].Articles)
	}
}

func (sect *Section) prepare(sectionIndex, playOrder int) {
	sect.SectionID = fmt.Sprintf("section-%d", sectionIndex)
	sect.PlayOrder = playOrder
	sect.Href = pad(playOrder) + ".html"

	for i := range sect.Articles {
		// The first article will have the same playOrder as the section that contains it
		sect.Articles[i].prepare(playOrder + i)
	}
}

func (art *Article) prepare(playOrder int) {
	art.PlayOrder = playOrder
	art.Href = pad(playOrder) + ".html"
}

func pad(number int) string {
	return fmt.Sprintf("%06d", number)
}

func createBlankCover() io.Reader {
	buf := new(bytes.Buffer)
	img := image.NewRGBA(image.Rect(0, 0, 600, 800))
	for i := range img.Pix {
		img.Pix[i] = 255
	}
	gif.Encode(buf, img, &gif.Options{NumColors: 1})
	return buf
}
