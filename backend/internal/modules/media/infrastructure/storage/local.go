package storage

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"webtranshome/internal/modules/media/domain/asset"
)

const defaultMaxBytes int64 = 2 * 1024 * 1024

type Local struct {
	directory string
	maxBytes  int64
}

func NewLocal(directory string, maxBytes int64) *Local {
	if strings.TrimSpace(directory) == "" {
		directory = "./uploads"
	}
	if maxBytes <= 0 {
		maxBytes = defaultMaxBytes
	}
	return &Local{directory: directory, maxBytes: maxBytes}
}

func (s *Local) Save(originalName string, data []byte) (asset.StoredFile, error) {
	if int64(len(data)) == 0 || int64(len(data)) > s.maxBytes {
		return asset.StoredFile{}, asset.ErrInvalidImage
	}
	mime, width, height, ext, ok := inspect(data)
	if !ok || width < 64 || height < 64 {
		return asset.StoredFile{}, asset.ErrInvalidImage
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return asset.StoredFile{}, err
	}
	filename := hex.EncodeToString(random[:]) + ext
	if err := os.MkdirAll(s.directory, 0o750); err != nil {
		return asset.StoredFile{}, err
	}
	if err := os.WriteFile(filepath.Join(s.directory, filename), data, 0o640); err != nil {
		return asset.StoredFile{}, err
	}
	return asset.StoredFile{Path: filename, MIME: mime, Size: int64(len(data)), Width: width, Height: height}, nil
}

func (s *Local) Delete(filename string) error {
	if filepath.Base(filename) != filename || filename == "." || filename == ".." {
		return asset.ErrInvalidAsset
	}
	err := os.Remove(filepath.Join(s.directory, filename))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *Local) Open(filename string) (string, error) {
	if filepath.Base(filename) != filename || filename == "." || filename == ".." {
		return "", asset.ErrInvalidAsset
	}
	path := filepath.Join(s.directory, filename)
	if _, err := os.Stat(path); err != nil {
		return "", err
	}
	return path, nil
}

func inspect(data []byte) (mime string, width, height int, ext string, ok bool) {
	switch {
	case len(data) >= 24 && string(data[:8]) == "\x89PNG\r\n\x1a\n":
		return "image/png", int(uint32(data[16])<<24 | uint32(data[17])<<16 | uint32(data[18])<<8 | uint32(data[19])), int(uint32(data[20])<<24 | uint32(data[21])<<16 | uint32(data[22])<<8 | uint32(data[23])), ".png", true
	case len(data) >= 4 && data[0] == 0xff && data[1] == 0xd8 && data[2] == 0xff:
		w, h, valid := jpegDimensions(data)
		return "image/jpeg", w, h, ".jpg", valid
	case len(data) >= 30 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP":
		w, h, valid := webpDimensions(data)
		return "image/webp", w, h, ".webp", valid
	default:
		return "", 0, 0, "", false
	}
}

func jpegDimensions(data []byte) (int, int, bool) {
	for i := 2; i+9 < len(data); {
		if data[i] != 0xff {
			i++
			continue
		}
		marker := data[i+1]
		if marker == 0xd8 || marker == 0xd9 || marker == 0x01 || marker >= 0xd0 && marker <= 0xd7 {
			i += 2
			continue
		}
		if i+4 > len(data) {
			break
		}
		length := int(data[i+2])<<8 | int(data[i+3])
		if length < 2 || i+2+length > len(data) {
			break
		}
		if marker >= 0xc0 && marker <= 0xc3 || marker >= 0xc5 && marker <= 0xc7 || marker >= 0xc9 && marker <= 0xcb || marker >= 0xcd && marker <= 0xcf {
			return int(data[i+7])<<8 | int(data[i+8]), int(data[i+5])<<8 | int(data[i+6]), true
		}
		i += 2 + length
	}
	return 0, 0, false
}

func webpDimensions(data []byte) (int, int, bool) {
	switch string(data[12:16]) {
	case "VP8X":
		if len(data) < 30 {
			return 0, 0, false
		}
		w := 1 + int(data[24]) + (int(data[25]) << 8) + (int(data[26]) << 16)
		h := 1 + int(data[27]) + (int(data[28]) << 8) + (int(data[29]) << 16)
		return w, h, w > 0 && h > 0
	case "VP8L":
		if len(data) < 25 || data[20] != 0x2f {
			return 0, 0, false
		}
		w := 1 + int(data[21]) + (int(data[22]&0x3f) << 8)
		h := 1 + int(data[22]>>6) + (int(data[23]) << 2) + (int(data[24]&0x0f) << 10)
		return w, h, w > 0 && h > 0
	case "VP8 ":
		if len(data) < 30 || data[23] != 0x9d || data[24] != 0x01 || data[25] != 0x2a {
			return 0, 0, false
		}
		w := int(data[26]) | int(data[27]&0x3f)<<8
		h := int(data[28]) | int(data[29]&0x3f)<<8
		return w, h, w > 0 && h > 0
	default:
		return 0, 0, false
	}
}
