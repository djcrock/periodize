package periodize

const opfTemplateString = `<package xmlns="http://www.idpf.org/2007/opf" version="2.0" unique-identifier="{{.UniqueID}}">
  <metadata>
    <meta content="cover-image" name="cover"/>
    <dc-metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
      <dc:title>{{.Title}}</dc:title>
      <dc:language>en-us</dc:language>
      <dc:creator>{{.Creator}}</dc:creator>
      <dc:publisher>{{.Publisher}}</dc:publisher>
      <dc:subject>{{.Subject}}</dc:subject>
      <dc:date>{{.Date}}</dc:date>
      <dc:description>{{.Description}}</dc:description>
    </dc-metadata>
    <x-metadata>
      <output content-type="application/x-mobipocket-subscription-magazine" encoding="utf-8"/>
    </x-metadata>
  </metadata>
  <manifest>
    <item href="cover-image.gif" media-type="image/gif" id="cover-image"/>
    <item href="contents.html" media-type="application/xhtml+xml" id="contents"/>
    <item href="nav-contents.ncx" media-type="application/x-dtbncx+xml" id="nav-contents"/>
{{range .Sections}}
{{range .Articles}}
    <item href="{{.Href}}" media-type="application/xhtml+xml" id="{{.PlayOrder}}"/>
{{end}}
{{end}}
  </manifest>
  <spine toc="nav-contents">
    <itemref idref="contents"/>
{{range .Sections}}
{{range .Articles}}
        <itemref idref="{{.PlayOrder}}"/>
{{end}}
{{end}}
  </spine>
  <guide>
    <reference href="contents.html" type="toc" title="Table of Contents"/>
  </guide>
</package>
`

const contentsTemplateString = `<html>
  <head>
    <meta content="text/html; charset=utf-8" http-equiv="Content-Type"/>
    <title>Table of Contents</title>
  </head>
  <body>
    <h1>Contents</h1>
{{range .Sections}}
    <h4>{{.Title}}</h4>
    <ul>
{{range .Articles}}
      <li>
        <a href="{{.Href}}">{{.Title}}</a>
      </li>
{{end}}
    </ul>
{{end}}
  </body>
</html>
`

const navTemplateString = `<!DOCTYPE ncx PUBLIC "-//NISO//DTD ncx 2005-1//EN" "http://www.daisy.org/z3986/2005/ncx-2005-1.dtd">
<ncx xmlns:mbp="http://mobipocket.com/ns/mbp" xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" xml:lang="en-US">
  <head>
    <meta content="Template" name="dtb:uid"/>
    <meta content="2" name="dtb:depth"/>
    <meta content="0" name="dtb:totalPageCount"/>
    <meta content="0" name="dtb:maxPageNumber"/>
  </head>
  <docTitle>
    <text>ncx:Title</text>
  </docTitle>
  <docAuthor>
    <text>ncs:Author</text>
  </docAuthor>
  <navMap>
    <navPoint playOrder="0" class="periodical" id="periodical">
      <mbp:meta-img src="masthead.gif" name="mastheadImage"/>
      <navLabel>
        <text>Table of Contents</text>
      </navLabel>
      <content src="contents.html"/>
{{range .Sections}}
      <navPoint playOrder="{{.PlayOrder}}" class="section" id="{{.SectionID}}">
        <navLabel>
          <text>{{.Title}}</text>
        </navLabel>
        <content src="{{.Href}}"/>
{{range .Articles}}
        <navPoint playOrder="{{.PlayOrder}}" class="article" id="{{.PlayOrder}}">
          <navLabel>
            <text>{{.Title}}</text>
          </navLabel>
          <content src="{{.Href}}"/>
          <mbp:meta name="description">{{.Title}}</mbp:meta>
          <mbp:meta name="author">{{.Author}}</mbp:meta>
        </navPoint>
{{end}}
      </navPoint>
{{end}}
    </navPoint>
  </navMap>
</ncx>
`
