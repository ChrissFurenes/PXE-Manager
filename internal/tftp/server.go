package tftp

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"PXE-Manager/internal/config"
	tftp "github.com/pin/tftp/v3"
)

func Start(cfg *config.Config) error {
	root := filepath.Clean(cfg.TFTP.RootDir)

	s := tftp.NewServer(
		func(filename string, rf io.ReaderFrom) error {
			cleanName := filepath.Clean("/" + filename)
			fullPath := filepath.Join(root, cleanName)

			if !strings.HasPrefix(filepath.Clean(fullPath), root) {
				return os.ErrPermission
			}

			f, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			defer f.Close()

			log.Printf("[TFTP] serving %s", filename)
			_, err = rf.ReadFrom(f)
			return err
		},
		nil,
	)

	log.Printf("[TFTP] listening on %s", cfg.TFTP.ListenAddr)
	return s.ListenAndServe(cfg.TFTP.ListenAddr)
}
