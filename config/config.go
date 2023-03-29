package config

import (
	"os"
	"path"
	"strconv"
)

type Config struct {
	MaxFileSize int64
	FilesPath   string
	Port        string
	Host        string
}

const (
	filesPathKey   = "FILES_PATH"
	portKey        = "FILE_SERVICE_PORT"
	maxImgSizeKey  = "MAX_FILE_SIZE"
	defaultMaxSize = 2 * 1024 * 1024
	defaultPort    = "8000"
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
			c.MaxFileSize = defaultMaxSize
		} else {
			c.MaxFileSize = max
		}
	} else {
		c.MaxFileSize = defaultMaxSize
	}

	return
}
