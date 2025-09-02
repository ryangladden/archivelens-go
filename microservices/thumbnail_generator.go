package microservices

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

func (dw *DocumentWorker) GenerateThumb(id string, filename string) error {
	original, err := dw.storageManager.CreateTempFile(id, "original", filename)
	if err != nil {
		return err
	}

	dest, err := dw.storageManager.CreateTempDir(id, "thumb")
	if err != nil {
		return err
	}

	thumb := filepath.Join(dest, "thumb.webp")

	log.Debug().Msgf("Converting %s to: %s", original, thumb)

	err = magickThumbnail(original, thumb)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/documents/%s/thumb.webp", id)
	err = dw.storageManager.UploadLocalFile(thumb, key)
	if err != nil {
		return err
	}

	os.RemoveAll(filepath.Join("/tmp", id))

	return nil
}

func magickThumbnail(input string, output string) error {

	page0 := input + "[0]"

	cmd := exec.Command(
		"magick",
		page0,
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
		log.Error().Err(err).Msgf("ImageMagick failed to generate thumbnail")
		return err
	}

	return nil
}
