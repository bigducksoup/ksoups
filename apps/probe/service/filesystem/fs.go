package filesystem

import "apps/probe/common"

type FileSystemService struct {
	cache *common.LRUCache[*common.File]
}

func NewFileSystemService() *FileSystemService {
	return &FileSystemService{
		cache: common.NewLRUCache[*common.File](10),
	}
}
