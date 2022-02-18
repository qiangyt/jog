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

func NewResourceWithPath(path string) Resource {
	return &ResourceT{path: path}
}

func NewResourceWithUrl(url string) Resource {
	return NewResourceWithPath(ResourcePath(url))
}

func (i Resource) NewKratoSource() StatikSource {
	return NewStatikSource(i)
}

func (me Resource) Path() string {
	return me.path
}

func (me Resource) Url() string {
	return filepath.Join(UrlPrefix, me.path)
}

func IsResourceUrl(url string) bool {
	return strings.HasPrefix(url, UrlPrefix)
}

func ResourcePath(url string) string {
	return url[len(UrlPrefix):]
}

func (me Resource) Open() http.File {
	r, err := Fs().Open(me.Path())
	if err != nil {
		panic(errors.Wrapf(err, "failed to open resource: %s", me.Path()))
	}
	return r
}

func (me Resource) CopyToFile(targetDir string) {
	content := me.ReadBytes()

	targetPath := filepath.Join(targetDir, me.Path())
	_io.MkdirAll(filepath.Dir(targetPath))

	_io.ReplaceFile(targetPath, content)
}

func (me Resource) ReadBytes() []byte {
	f := me.Open()
	defer f.Close()

	return _io.ReadAll(f)
}

func (me Resource) ReadString() string {
	return string(me.ReadBytes())
}
