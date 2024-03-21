package LPUCache

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"
)

var fileLock sync.Mutex

// Store stores a given data interface at the provided location
func Store(path string, data interface{}) error {
	fileLock.Lock()
	defer fileLock.Unlock()
	tmpPath := path + ".part"
	file, err := os.Create(tmpPath)
	defer os.RemoveAll(tmpPath)
	if err != nil {
		return err
	}
	defer file.Close()
	dataBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, bytes.NewReader(dataBytes))
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// Fetch fetches the data from the given path location and puts it into the given interface
func Fetch(path string, data interface{}) error {
	// Ensure File exists
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Read from file into
	fileLock.Lock()
	defer fileLock.Unlock()
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(data)
}
