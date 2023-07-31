package files

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	// Host is the URL to the server where the files are hosted.
	Host string
	// Folders contains folder and the amount of files they contain.
	Folders map[string]int
	// RootFolder is the root folder of the media files.
	// It's only used to generate a correct URL to the media being served.
	RootFolder string
}

func NewClient(host, rootFolder string) *Client {
	return &Client{
		Host:       host,
		Folders:    map[string]int{},
		RootFolder: rootFolder,
	}
}

func (c *Client) Init(path string) error {
	if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if err := c.upsertFolder(path); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

// HandleFolder will update the `Client.Folders` variable depending on what has happened with the supplied path.
//
// A path with no extension will be treated as a folder.
func (c *Client) HandleFolder(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		// File is not found, and must therefore be deleted
		if errors.Is(err, fs.ErrNotExist) {
			if filepath.Ext(path) == "" {
				if err := c.removeFolder(path); err != nil {
					return err
				}
			} else {
				if err := c.upsertFolder(filepath.Dir(path)); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	} else {
		if !fi.IsDir() {
			path = filepath.Dir(path)
		}

		if err := c.upsertFolder(path); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) upsertFolder(path string) error {
	log.WithField("path", path).Info("Upserting folder")

	dirEntries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	files := []string{}
	for _, dir := range dirEntries {
		if !dir.IsDir() {
			files = append(files, dir.Name())
		}
	}

	files = toVaildFiles(files)

	if len(files) > 0 {
		c.Folders[path] = len(files)
	} else {
		delete(c.Folders, path)
	}

	return nil
}

func (c *Client) removeFolder(path string) error {
	log.WithField("path", path).Info("Removing folder")

	newFolders := map[string]int{}

	for f, i := range c.Folders {
		if !strings.HasPrefix(f, path) {
			newFolders[f] = i
		} else {
			log.WithField("path", f).Debug("Removing subfolder")
		}
	}

	c.Folders = newFolders

	return nil
}

func (c *Client) RandomFile() (*File, error) {
	if len(c.Folders) == 0 {
		return nil, errors.New("no folders with valid files")
	}

	rand.New(rand.NewSource(time.Now().Unix()))
	imageAmount := 0
	randomMap := map[int]string{}
	keys := make([]int, 0)

	for f, amount := range c.Folders {
		randomMap[imageAmount] = f
		imageAmount += amount
		keys = append(keys, imageAmount)
	}

	sort.Ints(keys)

	fileNo := rand.Intn(imageAmount)
	foundKeyNo := 0

	for _, amount := range keys {
		if fileNo < amount {
			break
		}
		foundKeyNo = amount
	}

	dirName := randomMap[foundKeyNo]

	dirEntries, err := os.ReadDir(dirName)
	if err != nil {
		return nil, err
	}

	fileName := dirEntries[fileNo-foundKeyNo]

	fileLink := fmt.Sprintf("%s%s", c.Host, strings.TrimPrefix(dirName, c.RootFolder))
	fileRelativePath := fmt.Sprintf("%s/%s", strings.TrimPrefix(dirName, c.RootFolder), fileName.Name())

	file := &File{
		Title: fileRelativePath,
		Url:   fmt.Sprintf("%s/%s", fileLink, fileName.Name()),
	}

	log.WithField("file", fileRelativePath).Info("Found file")

	return file, nil
}
