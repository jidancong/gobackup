package helper

import (
	"context"
	"gobackup/utils"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v4"
)

type Archive struct {
	path string
}

func NewArchive(path string) *Archive {
	return &Archive{path: path}
}

func (a *Archive) Compressor(names []string) (string, error) {
	filemaps := make(map[string]string, 0)
	for _, v := range names {
		fileKey := filepath.Join(a.path, v)
		filemaps[fileKey] = ""
	}

	files, err := archiver.FilesFromDisk(nil, filemaps)
	if err != nil {
		return "", err
	}

	outName := filepath.Join(a.path, utils.TimeFormat()+".tar.gz")
	out, err := os.Create(outName)
	if err != nil {
		return "", err
	}
	defer out.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{
			CompressionLevel: 9,
		},
		Archival: archiver.Tar{},
	}

	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return "", err
	}

	return outName, nil
}

func (a *Archive) DeCompressor(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	format := archiver.Rar{}
	err = format.Extract(context.Background(), file, nil, func(ctx context.Context, f archiver.File) error {
		extractedFilePath := filepath.Join(a.path, f.NameInArchive)
		os.MkdirAll(filepath.Dir(extractedFilePath), 0755)

		af, err := f.Open()
		if err != nil {
			return err
		}
		defer af.Close()

		if !f.IsDir() {
			out, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				f.Mode(),
			)
			if err != nil {
				return err
			}
			defer out.Close()
			_, err = io.Copy(out, af)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
