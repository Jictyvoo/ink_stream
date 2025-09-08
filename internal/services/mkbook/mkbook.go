package mkbook

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-shiori/go-epub"

	"github.com/Jictyvoo/ink_stream/internal/services/imgprocessor"
	"github.com/Jictyvoo/ink_stream/internal/services/mkbook/tmplepub"
)

type EpubMounter struct {
	epub          *epub.Epub
	styleLocation string
	outDir        string
	coverInfo     struct{ location, name string }
}

func NewEpubMounter(outputDirectory string) (*EpubMounter, error) {
	title := filepath.Base(outputDirectory)
	if title == "." {
		asSha := sha256.Sum256([]byte(outputDirectory))
		title = string(asSha[:16])
	}
	e, err := epub.NewEpub(title)
	if err != nil {
		return nil, err
	}

	epubMounter := &EpubMounter{epub: e, outDir: outputDirectory}
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

func (em *EpubMounter) Handler(filename string, callback imgprocessor.WriterCallback) error {
	var buf bytes.Buffer
	if err := callback(&buf); err != nil {
		return fmt.Errorf("error while processing file %s: %w", filename, err)
	}
	location, err := writeBinaryFile(filename, buf.Bytes(), em.epub.AddImage)
	if err != nil {
		return fmt.Errorf("error while writing file `%s` to epub: %w", filename, err)
	}

	if em.coverInfo.name == "" || filename < em.coverInfo.name {
		em.coverInfo.name = filename
		em.coverInfo.location = location
	}
	// Add an image page to the EPUB using the written filename as the image source
	pageData := tmplepub.ImageData{ImageSrc: location}
	return em.AddImagePage(pageData, filename, filename)
}

func (em *EpubMounter) AddImagePage(
	pageData tmplepub.ImageData,
	sectionTitle, fileName string,
) error {
	// If viewport dimensions were not provided, use image dimensions if available
	if pageData.ViewportWidth == 0 && pageData.ImageWidth > 0 {
		pageData.ViewportWidth = pageData.ImageWidth
	}
	if pageData.ViewportHeight == 0 && pageData.ImageHeight > 0 {
		pageData.ViewportHeight = pageData.ImageHeight
	}

	if err := em.epub.SetCover(em.coverInfo.location, ""); err != nil {
		return err
	}

	tmpl := tmplepub.EpubImagePage()
	var buf strings.Builder
	if err := tmpl.Execute(&buf, pageData); err != nil {
		return err
	}
	// Add the rendered XHTML body as a section; go-epub will wrap it with full XHTML
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
