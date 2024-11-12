package cbxr

import (
	"fmt"
	"io"
	"iter"
	"strconv"
	"strings"

	pdfApi "github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

type PDFExtractor struct {
	fileReader FileContentStream
	pdfCtx     *model.Context
}

func NewPDFExtractor(fileReader FileContentStream) (*PDFExtractor, error) {
	pdfConf := model.NewDefaultConfiguration()
	ctx, err := pdfApi.ReadValidateAndOptimize(fileReader, pdfConf)
	if err != nil {
		return nil, err
	}
	return &PDFExtractor{fileReader: fileReader, pdfCtx: ctx}, nil
}

func (e *PDFExtractor) FileSeq() iter.Seq2[FileName, FileResult] {
	return func(yield func(FileName, FileResult) bool) {
		pageNrs := e.pdfCtx.PageCount

		for i := 0; i <= pageNrs; i++ {
			mm, err := pdfcpu.ExtractPageImages(e.pdfCtx, i, false)
			if err != nil {
				yield("", FileResult{Error: err})
				return
			}

			// singleImgPerPage := len(mm) == 1
			maxPageDigits := len(strconv.Itoa(pageNrs))
			for subindex, img := range mm {
				var result FileResult
				result.Data, result.Error = io.ReadAll(img.Reader)

				// fmt.Printf("%06d", 42)  // 0 padding, prints '000042'
				paddingFormatter := "%0" + strconv.Itoa(maxPageDigits+1) + "d"
				pageIndex := fmt.Sprintf(paddingFormatter, i)

				filename := strings.Join(
					[]string{pageIndex, fmt.Sprintf(paddingFormatter, subindex), img.Name},
					"_",
				)
				if !yield(FileName(filename), result) {
					return
				}
			}
		}

		return
	}
}
