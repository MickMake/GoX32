package embed

import (
	"embed"
	"errors"
	"fmt"
	"github.com/MickMake/GoX32/Only"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)


type Storage embed.FS

//go:embed *.json
var Files embed.FS


func SaveFiles(baseDir string, force bool) error {
	var err error

	for range Only.Once {
		// fmt.Println("Embedded files:")
		// fsys := fs.FS(files)
		//
		// fsRoot, _ := fs.Sub(fsys, ".")
		// fmt.Printf("%v\n", fsRoot)
		// spew.Dump(fsRoot)
		//
		// dir, _ := fs.ReadDir(fsys, ".")
		// fmt.Printf("Dir: %v\n", dir)
		// spew.Dump(dir)
		//
		//
		// dir, _ := fs.ReadDir(files, ".")
		// for _, file := range dir {
		// 	fmt.Println(file.Name())
		// }

		if !force && !IsEmpty(baseDir) {
			break
		}

		var storedRoot fs.FS
		storedRoot, err = fs.Sub(Files, ".")
		if err != nil {
			break
		}

		storedFiles := GetStoredFiles()
		fmt.Printf("Creating %d files\n", len(storedFiles))
		for _, storedFile := range storedFiles {
			fn := filepath.Join(baseDir, storedFile)
			fmt.Printf("\t%s\n", fn)

			var data []byte
			data, err = ReadFile(storedRoot, storedFile)
			if err != nil {
				break
			}

			err = WriteFile(fn, data, 0644)
			if err != nil {
				break
			}
		}
	}

	return err
}


// GetStoredFiles -
func GetStoredFiles() []string {
	var ret []string

	for range Only.Once {
		files, err := Files.ReadDir(".")
		if err != nil {
			break
		}
		for _, file := range files {
			ret = append(ret, file.Name())
		}
	}

	return ret
}

// IsEmpty -
func IsEmpty(dir string) bool {
	var ok bool

	for range Only.Once {
		ok = false

		if !DirExists(dir) {
			err := os.MkdirAll(dir, 0700)
			if err != nil {
				break
			}
			ok = true
			break
		}

		f, _ := DirectoryRead(dir, ".*.json")
		if len(f) == 0 {
			ok = true
			break
		}
	}

	return ok
}

// DirExists -
func DirExists(fn string) bool {
	var ok bool

	for range Only.Once {
		if fn == "" {
			ok = false
			break
		}

		f, err := os.Stat(fn)
		if os.IsNotExist(err) {
			ok = false
			break
		}
		if err != nil {
			ok = false
			break
		}
		if f.IsDir() {
			ok = true
			break
		}
	}

	return ok
}

// DirectoryRead -
func DirectoryRead(dir string, glob string) ([]string, error) {
	var ret []string
	var err error

	for range Only.Once {
		if dir == "" {
			err = errors.New("empty dir")
			break
		}

		var f []fs.FileInfo
		f, err = ioutil.ReadDir(dir)
		if err != nil {
			break
		}

		var re = regexp.MustCompile(glob)
		for _, file := range f {
			if file.IsDir() {
				continue
			}
			if file.Size() == 0 {
				continue
			}
			if !re.MatchString(file.Name()) {
				continue
			}
			ret = append(ret, file.Name())
			// fmt.Println(file.Name(), file.IsDir())
		}
	}

	return ret, err
}

// WriteFile Saves data to a file path.
func WriteFile(fn string, data []byte, perm os.FileMode) error {
	var err error

	for range Only.Once {
		if fn == "" {
			err = errors.New("empty file")
			break
		}

		var f *os.File
		f, err = os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			err = errors.New(fmt.Sprintf("Unable to write to file %s - %v", fn, err))
			break
		}

		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer f.Close()

		_, err = f.Write(data)
	}

	return err
}

// ReadFile Read data from a file path.
func ReadFile(fsRoot fs.FS, fn string) ([]byte, error) {
	var data []byte
	var err error

	for range Only.Once {
		if fn == "" {
			err = errors.New("empty file")
			break
		}

		var f fs.File
		f, err = fsRoot.Open(fn)
		if err != nil {
			break
		}

		data, err = ioutil.ReadAll(f)
		if err != nil {
			break
		}
		_ = f.Close()
	}

	return data, err
}
