package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	elUrl   = "https://developer.mozilla.org/en-US/docs/Web/HTML/Element"
	attrUrl = "https://developer.mozilla.org/en-US/docs/Web/HTML/Global_attributes"
	ariaUrl = "https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Attributes"
)

var void = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

var deprecated = map[string]bool{
	"acronym":   true,
	"big":       true,
	"center":    true,
	"dir":       true,
	"font":      true,
	"frame":     true,
	"frameset":  true,
	"image":     true,
	"marquee":   true,
	"menuitem":  true,
	"nobr":      true,
	"noembed":   true,
	"noframes":  true,
	"param":     true,
	"plaintext": true,
	"rb":        true,
	"rtc":       true,
	"strike":    true,
	"tt":        true,
	"xmp":       true,
}

var skipAttrs = map[string]bool{
	"data-*": true,
}

type Data struct {
	GlobalAttrs []Attr
	AriaAttrs   []Attr
	Elements    []Element
}

type Attr struct {
	Tag      string
	FuncName string
	Docs     string
	DocsLink string
}

type Element struct {
	Tag      string
	FuncName string
	DocsLink string
	Docs     string
	Void     bool
	Attrs    []Attr
}

func getPage(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getAttrs(url string, doc *goquery.Document) []Attr {
	attrs := []Attr{}
	found := map[string]bool{}

	doc.Find("dl dt").
		Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(
				cleanTag(
					s.Find("code").Text(),
				),
			)

			if name == "" || found[name] || skipAttrs[name] {
				return
			}

			found[name] = true

			docs := strings.TrimSpace(
				s.Next().Text(),
			)

			link := fmt.Sprintf("%s/%s", url, name)

			attrs = append(attrs, Attr{
				Tag: name,
				FuncName: toCamelCase(
					cleanName(name),
				),
				Docs:     docs,
				DocsLink: link,
			})
		})

	return attrs
}

func getAriaAttrs(doc *goquery.Document) []Attr {
	attrs := []Attr{}
	found := map[string]bool{}

	doc.Find("ol ul li a code").
		Each(func(i int, s *goquery.Selection) {
			name := strings.TrimSpace(
				s.Text(),
			)

			if name == "" || found[name] {
				return
			}

			found[name] = true

			attrs = append(attrs, Attr{
				Tag: name,
				FuncName: toCamelCase(
					cleanName(name),
				),
				DocsLink: fmt.Sprintf("%s/%s", ariaUrl, name),
			})
		})

	return attrs
}

func getElements(doc *goquery.Document) []Element {
	elements := []Element{}
	found := map[string]bool{}

	tr := doc.Find("tr")

	tr.
		Each(func(i int, s *goquery.Selection) {
			docs := s.Find("td").Eq(1).Text()

			s.Find("td").Eq(0).Find("a code").
				Each(func(i int, s *goquery.Selection) {
					name := strings.TrimSpace(
						cleanName(
							s.Text(),
						),
					)

					link := fmt.Sprintf("%s/%s", elUrl, name)

					if name == "" || found[name] || deprecated[name] {
						return
					}

					found[name] = true

					attrs := getAttrs(link, getPage(link))

					elements = append(elements, Element{
						Tag:      name,
						FuncName: toCamelCase(name),
						DocsLink: link,
						Docs:     docs,
						Void:     void[name],
						Attrs:    attrs,
					})
				})
		})

	return elements
}

func main() {
	data := Data{
		GlobalAttrs: getAttrs(attrUrl, getPage(attrUrl)),
		AriaAttrs:   getAriaAttrs(getPage(ariaUrl)),
		Elements:    getElements(getPage(elUrl)),
	}

	// b, err := json.MarshalIndent(data, "", "    ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = os.WriteFile("./gen/html.json", b, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	template, err := template.ParseFiles("./html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("elements.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = template.Execute(f, data)
	if err != nil {
		log.Fatal(err)
	}
}

func toCamelCase(s string) string {
	words := strings.Fields(
		cases.Title(language.English).String(s),
	)
	return strings.Join(words, "")
}

func cleanName(s string) string {
	replacer := strings.NewReplacer(
		"<", " ", ">", " ", "-", " ", "*", " ", ".", " ", "\"", " ", "(", " ", ")", " ",
	)
	return replacer.Replace(s)
}

func cleanTag(s string) string {
	replacer := strings.NewReplacer("\"", " ")
	return replacer.Replace(s)
}
