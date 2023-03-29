package config

import (
	"os"
	"path"
	"strconv"
)

type Config struct {
	MaxImageSize int64
	MaxPdfSize   int64
	FilesPath    string
	Port         string
	Host         string
}

const (
	filesPathKey         = "FILES_PATH"
	portKey              = "FILE_SERVICE_PORT"
	maxImgSizeKey        = "MAX_IMAGE_SIZE"
	maxPdfSizeKey        = "MAX_PDF_SIZE"
	defaultMaxSize       = 2 * 1024 * 1024
	defaultMaxPdfSizeKey = 3 * 1024
	defaultPort          = "8000"
)

func NewConfig() (c Config) {
	if p, set := os.LookupEnv(filesPathKey); set && p != "" {
		c.FilesPath = p
	} else {
		c.FilesPath = path.Join("srv", "files")
	}

	if port, set := os.LookupEnv(portKey); set && port != "" {
		c.Port = port
	} else {
		c.Port = defaultPort
	}

	if ms, set := os.LookupEnv(maxImgSizeKey); set && ms != "" {
		max, err := strconv.ParseInt(ms, 10, 0)
		if err != nil {
			c.MaxImageSize = defaultMaxSize
		} else {
			c.MaxImageSize = max
		}
	} else {
		c.MaxImageSize = defaultMaxSize
	}

	if ms, set := os.LookupEnv(maxPdfSizeKey); set && ms != "" {
		max, err := strconv.ParseInt(ms, 10, 0)
		if err != nil {
			c.MaxPdfSize = defaultMaxPdfSizeKey
		} else {
			c.MaxPdfSize = max
		}
	} else {
		c.MaxPdfSize = defaultMaxPdfSizeKey
	}

	return
}
