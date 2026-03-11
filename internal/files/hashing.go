package files

import (
	"archive/zip"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nwaples/rardecode"
	"github.com/saracen/go7z"
)

type FileHashes struct {
	MD5   string
	SHA1  string
	CRC32 string
}

// HashFile calcula los hashes de un archivo directo
func HashFile(path string) (*FileHashes, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return calculateHashes(f)
}

func calculateHashes(r io.Reader) (*FileHashes, error) {
	hMD5 := md5.New()
	hSHA1 := sha1.New()
	hCRC32 := crc32.NewIEEE()

	mw := io.MultiWriter(hMD5, hSHA1, hCRC32)

	if _, err := io.Copy(mw, r); err != nil {
		return nil, err
	}

	return &FileHashes{
		MD5:   hex.EncodeToString(hMD5.Sum(nil)),
		SHA1:  hex.EncodeToString(hSHA1.Sum(nil)),
		CRC32: hex.EncodeToString(hCRC32.Sum(nil)),
	}, nil
}

// HashCompressedFile calcula los hashes de un archivo dentro de un contenedor (zip, rar, 7z)
func HashCompressedFile(containerPath, internalFile string) (*FileHashes, error) {
	ext := strings.ToLower(filepath.Ext(containerPath))
	switch ext {
	case ".zip":
		return hashZip(containerPath, internalFile)
	case ".rar":
		return hashRar(containerPath, internalFile)
	case ".7z":
		return hash7z(containerPath, internalFile)
	}
	return nil, fmt.Errorf("formato no soportado: %s", ext)
}

func hashZip(path, internalFile string) (*FileHashes, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == internalFile {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			return calculateHashes(rc)
		}
	}
	return nil, fmt.Errorf("archivo %s no encontrado en zip", internalFile)
}

func hashRar(path, internalFile string) (*FileHashes, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rr, err := rardecode.NewReader(f, "")
	if err != nil {
		return nil, err
	}

	for {
		header, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if header.Name == internalFile {
			return calculateHashes(rr)
		}
	}
	return nil, fmt.Errorf("archivo %s no encontrado en rar", internalFile)
}

func hash7z(path, internalFile string) (*FileHashes, error) {
	sz, err := go7z.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer sz.Close()

	for {
		hdr, err := sz.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Name == internalFile {
			// go7z implementa io.Reader
			return calculateHashes(sz)
		}
	}
	return nil, fmt.Errorf("archivo %s no encontrado en 7z", internalFile)
}
