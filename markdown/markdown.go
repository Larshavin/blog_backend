package markdown

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
)

type folderInfo struct {
	Name   string `json:"name"`
	Matter matter `json:"matter"`
}

type matter struct {
	Title   string   `yaml:"title"`
	Tags    []string `yaml:"Tags,omitempty"`
	Content string   `yaml:"content,omitempty"`
	Date    string   `yaml:"date,omitempty"`
}

func FindFolderList(path string) ([]folderInfo, error) {
	var folders []folderInfo

	// Walk the directory starting from the given path
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		// Check if it's a directory and not the root directory itself
		if err == nil && info.IsDir() && filePath != path {
			// Read the markdown file
			content, err := os.ReadFile(filePath + "/" + info.Name() + ".md")
			if err != nil {
				fmt.Println(err)
			}

			var matter matter
			_, err = frontmatter.Parse(strings.NewReader(string(content)), &matter)
			if err != nil {
				fmt.Println(err)
			}

			modTime := info.ModTime().Format("2006-01-02T15:04:05-07:00")
			if matter.Date == "" {
				matter.Date = modTime
			}

			// Add the folder information to the list
			folders = append(folders, folderInfo{
				Name:   filePath,
				Matter: matter,
			})
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort the folders by modification time (latest first)
	sort.SliceStable(folders, func(i, j int) bool {
		time1, err := time.Parse(time.RFC3339, folders[i].Matter.Date)
		if err != nil {
			fmt.Println("Error parsing time1:", err)
		}
		time2, err := time.Parse(time.RFC3339, folders[j].Matter.Date)
		if err != nil {
			fmt.Println("Error parsing time1:", err)
		}
		return time1.After(time2)
	})

	return folders, nil
}
