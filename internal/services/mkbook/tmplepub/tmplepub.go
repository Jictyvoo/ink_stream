package tmplepub

import (
	_ "embed"
	"fmt"
	htmlTmpl "html/template"
	"text/template"
)

var (
	//go:embed epub_image_template.gohtml
	epubImageTemplate string

	//go:embed image_page_style.css
	epubImageStyleTemplate string
)

var loadedTemplates struct {
	epubImagePage  *htmlTmpl.Template
	epubImageStyle *template.Template
}

func init() {
	var err error
	loadedTemplates.epubImagePage, err = htmlTmpl.New("epub_image").Parse(epubImageTemplate)
	if err != nil {
		panic(fmt.Errorf("failed to parse epub_image template: %w", err))
	}

	if loadedTemplates.epubImageStyle, err = template.New("epub_image_style").Parse(epubImageStyleTemplate); err != nil {
		panic(fmt.Errorf("failed to parse epub_image_style template: %w", err))
	}
}

func EpubImagePage() *htmlTmpl.Template {
	return loadedTemplates.epubImagePage
}

func EpubImageStyle() *template.Template {
	return loadedTemplates.epubImageStyle
}
