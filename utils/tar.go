package utils

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

func Compressor(names []string, outPath string) error {
	filemaps := make(map[string]string, 0)
	for _, name := range names {
		filemaps[name] = filepath.Base(name)
	}

	files, err := archiver.FilesFromDisk(nil, filemaps)
	if err != nil {
		return err
	}

	outName := filepath.Join(outPath, TimeFormat()+".tar.gz")
	out, err := os.Create(outName)
	if err != nil {
		return err
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{
			CompressionLevel: 9,
		},
		Archival: archiver.Tar{},
	}

	err = format.Archive(context.Background(), out, files)
	return err
}
