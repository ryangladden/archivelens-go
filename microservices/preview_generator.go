package microservices

import (
	// 	"fmt"
	// 	"os/exec"
	// 	"path"
	// 	"strconv"

	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/ryangladden/archivelens-go/storage"
)

type PreviewGenerator struct {
	storageManager *storage.StorageManager
}

func NewPreviewGenerator(storageManager *storage.StorageManager) *PreviewGenerator {
	return &PreviewGenerator{
		storageManager: storageManager,
	}
}

func (p *PreviewGenerator) GeneratePreview(id string, filename string) (int, error) {
	tmpFile, err := p.storageManager.CreateTempFile(id, "original", filename)
	if err != nil {
		return 0, err
	}

	tmpDir, err := p.storageManager.CreateTempDir(id, "preview")
	if err != nil {
		return 0, err
	}

	output := filepath.Join(tmpDir, "preview")
	var pages int

	if filepath.Ext(filename) == ".pdf" {
		pages, err = p.magickPreviewPDF(tmpFile, output, id)
	} else {
		pages, err = p.magickPreviewIMG(tmpFile, output, id)
	}
	if err != nil {
		return 0, err
	}

	return pages, nil
}

func (p *PreviewGenerator) magickPreviewIMG(input string, output string, id string) (int, error) {

	output = fmt.Sprintf("%s.png")
	cmd := exec.Command(
		"magick",
		input,
		"-resize",
		"600x",
		"-background",
		"white",
		"-flatten",
		"-crop",
		"600x370+0+0",
		output,
	)

	log.Debug().Msg(cmd.String())

	err := cmd.Run()
	if err != nil {
		log.Error().Err(err).Msgf("ImageMagick failed to generate preview")
		return 0, err
	}

	key := filepath.Join("/documents", id, "preview", "preview-001.png")
	err = p.storageManager.UploadLocalFile(output, key)
	if err != nil {
		return 0, err
	}

	os.RemoveAll(filepath.Join("/tmp", id))

	return 1, nil
}

func (p *PreviewGenerator) magickPreviewPDF(input string, output string, id string) (int, error) {

	pages, err := getPageNumber(input)
	if err != nil {
		return 0, err
	}
	var numberFormat string
	if pages < 10 {
		numberFormat = "%d"
	} else if pages < 100 {
		numberFormat = "%02d"
	} else {
		numberFormat = "%03d"
	}

	cmd := exec.Command(
		"pdftoppm",
		"-png",
		input,
		output,
	)

	log.Debug().Msg(cmd.String())

	err = cmd.Run()
	if err != nil {
		log.Error().Err(err).Msg("Poppler failed to convert PDF to PNG")
	}

	if pages > 1 {
		// filename := fmt.Sprintf("preview-%s.png", numberFormat)
		for page := 1; page <= pages; page++ {
			// number := fmt.Sprintf()
			currentPage := fmt.Sprintf("/tmp/%s/preview/preview-"+numberFormat+".png", id, page)
			key := fmt.Sprintf("/documents/%s/preview/preview-%03d.png", id, page)
			err = p.storageManager.UploadLocalFile(currentPage, key)
		}
	} else {
		tmpFile := filepath.Join(output, "preview-1.png")
		key := fmt.Sprintf("/documents/%s/preview/preview-001.png", id)
		err = p.storageManager.UploadLocalFile(tmpFile, key)
	}

	if err != nil {
		log.Error().Err(err).Msgf("Failed to upload temp file(s) %s-*", output)
		return 0, err
	}

	return pages, nil
}

func getPageNumber(input string) (int, error) {
	cmd := exec.Command(
		"pdfinfo",
		input,
	)

	log.Debug().Msg(cmd.String())

	out, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to run pdfinfo on %s", input)
		return 0, err
	}
	log.Debug().Msg(string(out))

	pagesIndex := strings.Index(string(out), "Pages:")
	endIndex := strings.Index(string(out)[pagesIndex+1:], "\n")
	log.Debug().Msgf("Indexes for page number: %d to %d", pagesIndex+6, endIndex+pagesIndex+1)
	pages := strings.TrimSpace(string(out[pagesIndex+6 : endIndex+pagesIndex+1]))
	log.Debug().Msgf("Extracted pages: %s", pages)

	pageNumber, err := strconv.Atoi(pages)
	if err != nil {
		log.Error().Err(err).Msgf("Unable to get number of pages from %s", input)
	}

	return pageNumber, nil
}
