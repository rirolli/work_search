// Package xlsx scrive file .xlsx minimali (SpreadsheetML) usando solo la
// libreria standard: testo, numeri, stili fissi, link, riga bloccata e filtro.
package xlsx

import (
	"archive/zip"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Stili disponibili (indici in cellXfs di styles.xml).
const (
	Default = iota
	Header
	Wrap
	Link
	High
	Medium
	Low
	Bold
	New
)

// Cell è una cella. Se Num != nil viene scritta come numero.
type Cell struct {
	Text  string
	Num   *float64
	Style int
	URL   string // se valorizzato, la cella diventa un collegamento
}

// T crea una cella di testo.
func T(s string, style int) Cell { return Cell{Text: s, Style: style} }

// N crea una cella numerica.
func N(v float64, style int) Cell { return Cell{Num: &v, Style: style} }

// L crea un collegamento.
func L(text, url string) Cell { return Cell{Text: text, URL: url, Style: Link} }

// Sheet è un foglio di lavoro.
type Sheet struct {
	Name       string
	Rows       [][]Cell
	Widths     []float64
	FreezeRow  bool // blocca la prima riga
	AutoFilter bool // filtro sulla prima riga
}

// Col converte un indice di colonna (0-based) in lettere: 0→A, 26→AA.
func Col(i int) string {
	s := ""
	for i >= 0 {
		s = string(rune('A'+i%26)) + s
		i = i/26 - 1
	}
	return s
}

func esc(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == '&':
			b.WriteString("&amp;")
		case r == '<':
			b.WriteString("&lt;")
		case r == '>':
			b.WriteString("&gt;")
		case r == '"':
			b.WriteString("&quot;")
		case r == '\t' || r == '\n' || r == '\r':
			b.WriteRune(r)
		case r < 0x20 || r == 0xFFFE || r == 0xFFFF:
			// carattere non valido in XML: scartato
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func clip(s string) string {
	const max = 32000 // limite Excel: 32767 caratteri per cella
	if len(s) <= max {
		return s
	}
	r := []rune(s)
	if len(r) > max {
		r = r[:max]
	}
	return string(r)
}

// sanitizeName rende valido il nome di un foglio (max 31 caratteri, niente []:*?/\).
func sanitizeName(s string) string {
	s = strings.Map(func(r rune) rune {
		if strings.ContainsRune(`[]:*?/\`, r) {
			return '-'
		}
		return r
	}, s)
	if r := []rune(s); len(r) > 31 {
		s = string(r[:31])
	}
	return s
}

// Write crea il file xlsx.
func Write(path string, sheets []Sheet) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(f)
	add := func(name, content string) error {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(content))
		return err
	}

	var ct, wbSheets, wbRels, defined strings.Builder
	ct.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
<Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>
<Override PartName="/docProps/core.xml" ContentType="application/vnd.openxmlformats-package.core-properties+xml"/>
<Override PartName="/docProps/app.xml" ContentType="application/vnd.openxmlformats-officedocument.extended-properties+xml"/>
`)
	for i, sh := range sheets {
		n := i + 1
		fmt.Fprintf(&ct, `<Override PartName="/xl/worksheets/sheet%d.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>`+"\n", n)
		fmt.Fprintf(&wbSheets, `<sheet name="%s" sheetId="%d" r:id="rId%d"/>`, esc(sanitizeName(sh.Name)), n, n)
		fmt.Fprintf(&wbRels, `<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet%d.xml"/>`, n, n)
		if sh.AutoFilter && len(sh.Rows) > 0 {
			fmt.Fprintf(&defined, `<definedName name="_xlnm._FilterDatabase" localSheetId="%d" hidden="1">'%s'!$A$1:$%s$%d</definedName>`,
				i, esc(strings.ReplaceAll(sanitizeName(sh.Name), "'", "''")), Col(len(sh.Rows[0])-1), len(sh.Rows))
		}
	}
	ct.WriteString(`</Types>`)
	styleRel := len(sheets) + 1
	fmt.Fprintf(&wbRels, `<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>`, styleRel)

	definedNames := ""
	if defined.Len() > 0 {
		definedNames = "<definedNames>" + defined.String() + "</definedNames>"
	}

	files := []struct{ name, content string }{
		{"[Content_Types].xml", ct.String()},
		{"_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
</Relationships>`},
		{"docProps/core.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<cp:coreProperties xmlns:cp="http://schemas.openxmlformats.org/package/2006/metadata/core-properties" xmlns:dc="http://purl.org/dc/elements/1.1/"><dc:title>Ricerca lavoro</dc:title><dc:creator>work_search</dc:creator></cp:coreProperties>`},
		{"docProps/app.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Properties xmlns="http://schemas.openxmlformats.org/officeDocument/2006/extended-properties"><Application>work_search</Application></Properties>`},
		{"xl/workbook.xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets>` + wbSheets.String() + `</sheets>` + definedNames + `</workbook>`},
		{"xl/_rels/workbook.xml.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + wbRels.String() + `</Relationships>`},
		{"xl/styles.xml", stylesXML},
	}
	for _, fl := range files {
		if err := add(fl.name, fl.content); err != nil {
			f.Close()
			return err
		}
	}
	for i, sh := range sheets {
		body, rels := sheetXML(sh)
		if err := add(fmt.Sprintf("xl/worksheets/sheet%d.xml", i+1), body); err != nil {
			f.Close()
			return err
		}
		if rels != "" {
			if err := add(fmt.Sprintf("xl/worksheets/_rels/sheet%d.xml.rels", i+1), rels); err != nil {
				f.Close()
				return err
			}
		}
	}
	if err := zw.Close(); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func sheetXML(sh Sheet) (string, string) {
	var b, links, rels strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">`)
	if sh.FreezeRow {
		b.WriteString(`<sheetViews><sheetView workbookViewId="0"><pane ySplit="1" topLeftCell="A2" activePane="bottomLeft" state="frozen"/><selection pane="bottomLeft" activeCell="A2" sqref="A2"/></sheetView></sheetViews>`)
	}
	if len(sh.Widths) > 0 {
		b.WriteString("<cols>")
		for i, w := range sh.Widths {
			fmt.Fprintf(&b, `<col min="%d" max="%d" width="%.1f" customWidth="1"/>`, i+1, i+1, w)
		}
		b.WriteString("</cols>")
	}
	b.WriteString("<sheetData>")
	linkID := 0
	for r, row := range sh.Rows {
		fmt.Fprintf(&b, `<row r="%d">`, r+1)
		for c, cell := range row {
			ref := Col(c) + strconv.Itoa(r+1)
			if cell.Num != nil {
				fmt.Fprintf(&b, `<c r="%s" s="%d"><v>%s</v></c>`, ref, cell.Style, strconv.FormatFloat(*cell.Num, 'f', -1, 64))
				continue
			}
			if cell.Text == "" && cell.Style == Default {
				continue
			}
			fmt.Fprintf(&b, `<c r="%s" s="%d" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, ref, cell.Style, esc(clip(cell.Text)))
			if cell.URL != "" && len(cell.URL) < 2000 {
				linkID++
				fmt.Fprintf(&links, `<hyperlink ref="%s" r:id="rId%d"/>`, ref, linkID)
				fmt.Fprintf(&rels, `<Relationship Id="rId%d" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="%s" TargetMode="External"/>`, linkID, esc(cell.URL))
			}
		}
		b.WriteString("</row>")
	}
	b.WriteString("</sheetData>")
	if sh.AutoFilter && len(sh.Rows) > 0 {
		fmt.Fprintf(&b, `<autoFilter ref="A1:%s%d"/>`, Col(len(sh.Rows[0])-1), len(sh.Rows))
	}
	if links.Len() > 0 {
		b.WriteString("<hyperlinks>" + links.String() + "</hyperlinks>")
	}
	b.WriteString(`<pageMargins left="0.5" right="0.5" top="0.75" bottom="0.75" header="0.3" footer="0.3"/></worksheet>`)
	relXML := ""
	if rels.Len() > 0 {
		relXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">` + rels.String() + `</Relationships>`
	}
	return b.String(), relXML
}

// Ordine di cellXfs = costanti Default..New.
const stylesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
<fonts count="4">
<font><sz val="10"/><name val="Arial"/><family val="2"/></font>
<font><b/><sz val="10"/><color rgb="FFFFFFFF"/><name val="Arial"/><family val="2"/></font>
<font><u/><sz val="10"/><color rgb="FF0563C1"/><name val="Arial"/><family val="2"/></font>
<font><b/><sz val="10"/><name val="Arial"/><family val="2"/></font>
</fonts>
<fills count="7">
<fill><patternFill patternType="none"/></fill>
<fill><patternFill patternType="gray125"/></fill>
<fill><patternFill patternType="solid"><fgColor rgb="FF1F4E79"/><bgColor indexed="64"/></patternFill></fill>
<fill><patternFill patternType="solid"><fgColor rgb="FFC6EFCE"/><bgColor indexed="64"/></patternFill></fill>
<fill><patternFill patternType="solid"><fgColor rgb="FFFFEB9C"/><bgColor indexed="64"/></patternFill></fill>
<fill><patternFill patternType="solid"><fgColor rgb="FFEDEDED"/><bgColor indexed="64"/></patternFill></fill>
<fill><patternFill patternType="solid"><fgColor rgb="FFDDEBF7"/><bgColor indexed="64"/></patternFill></fill>
</fills>
<borders count="1"><border><left/><right/><top/><bottom/><diagonal/></border></borders>
<cellStyleXfs count="1"><xf numFmtId="0" fontId="0" fillId="0" borderId="0"/></cellStyleXfs>
<cellXfs count="9">
<xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0" applyAlignment="1"><alignment vertical="top"/></xf>
<xf numFmtId="0" fontId="1" fillId="2" borderId="0" xfId="0" applyFont="1" applyFill="1" applyAlignment="1"><alignment vertical="center" wrapText="1"/></xf>
<xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0" applyAlignment="1"><alignment vertical="top" wrapText="1"/></xf>
<xf numFmtId="0" fontId="2" fillId="0" borderId="0" xfId="0" applyFont="1" applyAlignment="1"><alignment vertical="top"/></xf>
<xf numFmtId="0" fontId="3" fillId="3" borderId="0" xfId="0" applyFont="1" applyFill="1" applyAlignment="1"><alignment vertical="top" horizontal="center"/></xf>
<xf numFmtId="0" fontId="3" fillId="4" borderId="0" xfId="0" applyFont="1" applyFill="1" applyAlignment="1"><alignment vertical="top" horizontal="center"/></xf>
<xf numFmtId="0" fontId="0" fillId="5" borderId="0" xfId="0" applyFill="1" applyAlignment="1"><alignment vertical="top" horizontal="center"/></xf>
<xf numFmtId="0" fontId="3" fillId="0" borderId="0" xfId="0" applyFont="1" applyAlignment="1"><alignment vertical="top"/></xf>
<xf numFmtId="0" fontId="3" fillId="6" borderId="0" xfId="0" applyFont="1" applyFill="1" applyAlignment="1"><alignment vertical="top" horizontal="center"/></xf>
</cellXfs>
<cellStyles count="1"><cellStyle name="Normal" xfId="0" builtinId="0"/></cellStyles>
</styleSheet>`
