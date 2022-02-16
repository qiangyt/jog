package res

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	_io "github.com/qiangyt/jog/pkg/io"
	statikFs "github.com/rakyll/statik/fs"

	_ "github.com/qiangyt/jog/res/statik"
)

const UrlPrefix = "res:"

var _fs http.FileSystem

func init() {
	var err error
	_fs, err = statikFs.NewWithNamespace("res")
	if err != nil {
		panic(errors.Wrap(err, "failed to create resource file system"))
	}
}

func Fs() http.FileSystem {
	return _fs
}

type ResourceT struct {
	path string
}

type Resource = *ResourceT

func New(path string) Resource {
	return &ResourceT{path: path}
}

func (i Resource) NewKratoSource() StatikSource {
	return NewStatikSource(i)
}

func (i Resource) Path() string {
	return i.path
}

func (i Resource) Url() string {
	return filepath.Join(UrlPrefix, i.path)
}

func IsResourceUrl(url string) bool {
	return strings.HasPrefix(url, UrlPrefix)
}

func ResourcePath(url string) string {
	return url[len(UrlPrefix):]
}

func (i Resource) Open() http.File {
	r, err := Fs().Open(i.Path())
	if err != nil {
		panic(errors.Wrapf(err, "failed to open resource: %s", i.Path()))
	}
	return r
}

func (i Resource) CopyToFile(targetDir string) {
	content := i.ReadBytes()

	targetPath := filepath.Join(targetDir, i.Path())
	_io.MkdirAll(filepath.Dir(targetPath))

	_io.ReplaceFile(targetPath, content)
}

func (i Resource) ReadBytes() []byte {
	f := i.Open()
	defer f.Close()

	return _io.ReadAll(f)
}

func (i Resource) ReadString() string {
	return string(i.ReadBytes())
}
