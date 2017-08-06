package sftp

import (
	"os"
	"path/filepath"
	"strings"
	"sort"
	"io"
	"log"
	"strconv"
)

type vfs struct {
	files []string
	pathMap map[string][]string
}

func VfsHandler(matchingPaths []string) Handlers {
	virtualFileSystem := &vfs{}

	sort.Strings(matchingPaths)


	virtualFileSystem.files = matchingPaths
	virtualFileSystem.pathMap = MakePathMap(matchingPaths)


	return Handlers{
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
		virtualFileSystem,
	}
}


func (fs *vfs) Fileread(r Request) (io.ReaderAt, error) {
	log.Printf("vfshandler.go->Fileread: r=%+v", r)
	foundFile := fetch(fs, r.Filepath)
	// log.Println("foundFile: ", foundFile)
	if(foundFile == "") {
		return nil, os.ErrInvalid
	}

	f, err := os.Open(foundFile)

	if err != nil {
		// log.Println("Could not open file", foundFile, err)
		return nil, os.ErrInvalid
	}
	return f, nil
}

func (fs *vfs) Filewrite(r Request) (io.WriterAt, error) {
	log.Printf("vfshandler.go->Filewrite: r=%+v", r)
	return nil, os.ErrInvalid
}

func (fs *vfs) Filecmd(r Request) error {
	log.Printf("vfshandler.go->Filecmd: r=%+v", r)
	return os.ErrInvalid
}

func (fs *vfs) Fileinfo(r Request) ([]os.FileInfo, error) {
	log.Printf("vfshandler.go->FileInfo: r=%+v", r)
	requestedPath := filepath.ToSlash(r.Filepath)


	switch r.Method {
	case "List":
		var err error
		batch_size := 10
		current_offset := 0
		if token := r.LsNext(); token != "" {
			current_offset, err = strconv.Atoi(token)
			if err != nil {
				return nil, os.ErrInvalid
			}
		}

		ordered_names, ok := fs.pathMap[requestedPath]
		if ! ok {
			// log.Println("did not find pathMapping for requestedPath", requestedPath)
			return nil, os.ErrInvalid
		}
		// log.Println("pathMapping for " + requestedPath + " contains: ", len(ordered_names))
		sort.Sort(sort.StringSlice(ordered_names))
		list := make([]os.FileInfo, len(ordered_names))
		for i, fileName := range ordered_names {

			// if runtime.GOOS != "windows" {
			//if err := syscall.Stat(file, &stat); err != nil {
			//	panic(err)
			//}

			stat, err := os.Stat(fileName)
			if err != nil {
				 log.Println("  List => Could not stat file", fileName, err)
				continue
			}

			log.Printf("  List => stat.Name()=%s, stat.IsDir()=%v", stat.Name(), stat.IsDir())

			list[i] = stat
			// log.Println("Stat for file " + fileName + ": isDir=>",stat.IsDir(), "size=>", stat.Size())
		}

		if len(list) < current_offset {
			return nil, io.EOF
		}

		new_offset := current_offset + batch_size
		if new_offset > len(list) {
			new_offset = len(list)
		}
		r.LsSave(strconv.Itoa(new_offset))
		return list[current_offset:new_offset], nil

	case "Stat":
		// log.Println("Stat filepath: ", requestedPath)
		foundFile := fetch(fs, requestedPath)
		if foundFile != "" {
			// log.Println("foundFile: ", foundFile)
			stat, err := os.Stat(foundFile)
			if err != nil {
				log.Println("  Stat => Could not stat file", foundFile, err)
				return nil, os.ErrInvalid
			}
			log.Printf("  Stat => stat.Name()=%s, stat.IsDir()=%v", stat.Name(), stat.IsDir())
			return []os.FileInfo{stat}, nil
		}
		// log.Println("Could not 'fetch' file for " + requestedPath)
		return nil, os.ErrInvalid
	}
	return nil, os.ErrInvalid
}


func fetch(fs *vfs, requestedPath string) string {
	// log.Println("fetch requestedPath:  ", requestedPath)

	key := filepath.ToSlash(filepath.Dir(requestedPath))
	// log.Println("mapping key:  ", key)

	ordered_names, ok := fs.pathMap[key]

	if ok == false {
		// log.Println("did not find key in pathMap", key)
		return ""
	}

	var foundFile = ""

	for _, b := range ordered_names {
		if b == requestedPath || b == strings.TrimLeft(requestedPath, "/") {
			foundFile = filepath.ToSlash(b)
			break
		}
	}
	return foundFile
}


func MakePathMap(matchingPaths []string) map[string][]string {
	pathMap := make(map[string][]string)

	sort.Strings(matchingPaths)

	//if val, ok := dict["foo"]; ok {
	//	//do something here
	//}

	for _, path := range matchingPaths {
		key, parentPath := normalizePathMapItem(path)

		for  {
			// println("append: ", key, " => ", path)
			pathMap[key] = append(pathMap[key], path)
			path = parentPath
			//println("before => key:", key, "parentPath:", parentPath)
			key, parentPath = normalizePathMapItem(parentPath)
			//println("after  => key:", key, "parentPath:", parentPath)
			_, ok := pathMap[key]

			//println("is present?", key, ok)
			if ok {
				break
			}
		}
	}



	return pathMap
}

func normalizePathMapItem(path string) (string, string) {
	parentPath := filepath.ToSlash(filepath.Dir(path))
	key := parentPath
	if parentPath == "." {
		key = "/"
	}

	firstChar := string([]rune(key)[0])
	if firstChar != "/" {
		key = "/" + key
	}
	return key, parentPath
}
