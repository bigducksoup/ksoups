package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var validFileSuffix = map[string]bool{
	".txt":    true,
	".md":     true,
	".go":     true,
	".js":     true,
	".java":   true,
	".cpp":    true,
	".php":    true,
	".py":     true,
	".sh":     true,
	".c":      true,
	".h":      true,
	".rs":     true,
	".ts":     true,
	".tsx":    true,
	".jsx":    true,
	".vue":    true,
	".swift":  true,
	".kt":     true,
	".kts":    true,
	".rb":     true,
	".cs":     true,
	".clj":    true,
	".groovy": true,
	".config": true,
	".json":   true,
	".xml":    true,
	".yml":    true,
	".yaml":   true,
	".csv":    true,
	".ini":    true,
	".toml":   true,
	".sql":    true,
	".conf":   true,
	".html":   true,
	"csv":     true,
	".css":    true,
}

func checkFilePathValid(path string) (os.FileInfo, bool) {

	fileInfo, err := os.Stat(path)

	if err != nil {
		return nil, false
	}

	if fileInfo.IsDir() {
		return nil, false
	}

	if fileInfo.Size() > 1024*1024*20 {
		return nil, false
	}

	if !strings.Contains(path, ".log") {
		//get file extension name
		ext := strings.ToLower(filepath.Ext(path))

		if condition, ok := validFileSuffix[ext]; !ok || !condition {
			return nil, false
		}
	}

	return fileInfo, true

}

type File struct {
	Path    string
	Size    int64
	Content []string
	lock    sync.Mutex
}

func NewFile(path string) (*File, error) {

	//check if path and file type is valid
	fileInfo, ok := checkFilePathValid(path)

	if !ok {
		return nil, errors.New("path is not valid")
	}

	//read whole file from disk to memory
	ct, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	content := strings.Split(string(ct), "\n")

	file := &File{
		Path:    path,
		Size:    fileInfo.Size(),
		Content: content,
		lock:    sync.Mutex{},
	}

	return file, nil

}

// 获取文件内容
func (f *File) GetContent() string {

	res := strings.Join(f.Content, "\n")

	return res
}

// ReplaceLine replaces a line in the File.
func (f *File) ReplaceLine(index int, content string) error {
	//加锁修改
	f.lock.Lock()

	//释放锁
	defer f.lock.Unlock()

	if index < 0 || index > len(f.Content) {
		return errors.New("index out of range")
	}

	f.Content[index] = content

	return nil
}

// RemoveLine removes a line from the File.
func (f *File) RemoveLine(index int) error {

	err := f.RemoveLines(index, index+1)

	return err
}

// RemoveLines removes lines from the File.
// e.g.
// file.RemoveLines(2,5) will remove line 2,3,4
func (f *File) RemoveLines(start, end int) error {

	f.lock.Lock()

	defer f.lock.Unlock()

	if start < 0 || start > len(f.Content) || end < 0 || end > len(f.Content) {
		return errors.New("start index out of range")
	}

	f.Content = append(f.Content[:start], f.Content[end:]...)

	return nil

}

// InsertLines inserts the given content before line:index in the File.
// e.g.
//
//		file content is {
//			"line1",
//	 		"line2"
//	 		"line3"
//	 		"line4"
//		}
//		file.InsertLines(2, []string{"linex", "liney"})
//		the result is {
//			"line1",
//	 		"line2"
//			"linex"
//			"liney"
//	 		"line3"
//	 		"line4"
//	}
func (f *File) InsertLines(index int, content []string) error {

	f.lock.Lock()
	defer f.lock.Unlock()

	// Check if the index is out of range
	if index < 0 || index > len(f.Content) {
		return errors.New("index out of range")
	}

	// Create a new slice to hold the updated content
	updatedContent := make([]string, len(f.Content)+len(content))

	// Copy the lines before the insertion index
	copy(updatedContent, f.Content[:index])

	// Copy the inserted lines
	copy(updatedContent[index:], content)

	// Copy the lines after the insertion index
	copy(updatedContent[index+len(content):], f.Content[index:])

	// Update the content of the file
	f.Content = updatedContent

	return nil

}

// InsertLine inserts a line of content at the specified index in the File.
//
// Parameters:
// - index: the index at which to insert the line.
// - content: the content of the line to be inserted.
//
// Returns:
// - error: an error if the insertion fails.
func (f *File) InsertLine(index int, content string) error {

	err := f.InsertLines(index, []string{content})
	return err
}

// AppendLines appends the given content to the File's content slice.
//
// It takes in a parameter named content, which is a slice of strings representing the lines to be appended.
// The function does not return anything.
func (f *File) AppendLines(content []string) error {

	f.lock.Lock()

	defer f.lock.Unlock()

	merged := append(f.Content, content...)

	f.Content = merged

	return nil

}

// GetSize returns the size of the File.
//
// It calculates the size of the File by summing the length of each content
// slice element and adding the length of the content slice itself which represents the count of '\n'.
//
// Return:
//
//	int64 - The size of the File.
func (f *File) GetSize() int64 {

	s := 0

	for _, v := range f.Content {
		s += len(v)
	}

	size := len(f.Content) + s

	return int64(size)

}

// 文件落盘
func (f *File) Flush() error {

	f.lock.Lock()

	defer f.lock.Unlock()

	res := strings.Join(f.Content, "\n")

	err := os.WriteFile(f.Path, []byte(res), 0644)

	if err != nil {
		return err
	}
	return nil
}
