package mkbook

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-shiori/go-epub"

	"github.com/Jictyvoo/ink_stream/internal/services/mkbook/tmplepub"
)

type EpubMounter struct {
	epub          *epub.Epub
	styleLocation string
	outDir        string
}

func NewEpubMounter(title string) (*EpubMounter, error) {
	e, err := epub.NewEpub(title)
	if err != nil {
		return nil, err
	}

	epubMounter := &EpubMounter{epub: e, outDir: title}
	err = epubMounter.registerMainCSS()

	return epubMounter, err
}

func (em *EpubMounter) registerMainCSS() (err error) {
	// `data:text/plain;charset=utf-8;base64,aGV5YQ==`
	var buffer bytes.Buffer
	if err = tmplepub.EpubImageStyle().Execute(&buffer, nil); err != nil {
		return err
	}

	em.styleLocation, err = writeBinaryFile("style.css", buffer.Bytes(), em.epub.AddCSS)
	return err
}

func (em *EpubMounter) Handler(filename string, callback func(writer io.Writer) error) error {
	var buf bytes.Buffer
	if err := callback(&buf); err != nil {
		return fmt.Errorf("error while processing file %s: %w", filename, err)
	}
	location, err := writeBinaryFile(filename, buf.Bytes(), em.epub.AddImage)
	if err != nil {
		return fmt.Errorf("error while writing file `%s` to epub: %w", filename, err)
	}

	// Add an image page to the EPUB using the written filename as the image source
	pageData := tmplepub.ImageData{ImageSrc: location}
	return em.AddImagePage(pageData, filename, filename)
}

func (em *EpubMounter) AddImagePage(
	pageData tmplepub.ImageData,
	sectionTitle, fileName string,
) error {
	tmpl := tmplepub.EpubImagePage()
	var buf strings.Builder
	if err := tmpl.Execute(&buf, pageData); err != nil {
		return err
	}
	// Add the rendered XHTML as a section to the EPUB
	_, err := em.epub.AddSection(buf.String(), sectionTitle, fileName, em.styleLocation)
	return err
}

func (em *EpubMounter) Flush() error {
	file, err := os.Create(em.outDir + ".epub")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = em.epub.WriteTo(file)
	return err
}
