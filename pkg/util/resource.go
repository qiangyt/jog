package util

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/pkg/errors"
	statikFs "github.com/rakyll/statik/fs"

	_ "github.com/qiangyt/jog/res/statik"
)

const ResourceUrlPrefix = "resource://"

var _resourceFs http.FileSystem

func init() {
	var err error
	_resourceFs, err = statikFs.NewWithNamespace("res")
	if err != nil {
		panic(errors.Wrap(err, "failed to create resource file system"))
	}
}

func ResourceFs() http.FileSystem {
	return _resourceFs
}

var _ config.Source = (*StatikSourceT)(nil)

type ResourceT struct {
	path string
}

type Resource = *ResourceT

// NewResource return new a statik resource.
func NewResource(path string) Resource {
	return &ResourceT{path: path}
}

func (i Resource) NewKratoSource() StatikSource {
	return NewStatikSource(i)
}

func (i Resource) Path() string {
	return i.path
}

func (i Resource) Url() string {
	return filepath.Join(ResourceUrlPrefix, i.path)
}

func IsResourceUrl(url string) bool {
	return strings.HasPrefix(url, ResourceUrlPrefix)
}

func ResourcePath(url string) string {
	return url[len(ResourceUrlPrefix):]
}

func (i Resource) Open() http.File {
	r, err := ResourceFs().Open(i.Path())
	if err != nil {
		panic(errors.Wrapf(err, "failed to open resource: %s", i.Path()))
	}
	return r
}

func (i Resource) CopyToFile(targetDir string) {
	content := i.ReadBytes()

	targetPath := filepath.Join(targetDir, i.Path())
	MkdirAll(filepath.Dir(targetPath))

	ReplaceFile(targetPath, content)
}

func (i Resource) ReadBytes() []byte {
	f := i.Open()
	defer f.Close()

	return ReadAll(f)
}

func (i Resource) ReadString() string {
	return string(i.ReadBytes())
}

// ----------------------
type StatikSourceT struct {
	res Resource
}

type StatikSource = *StatikSourceT

// NewStatikSource return new a statik file source.
func NewStatikSource(res Resource) StatikSource {
	return &StatikSourceT{res: res}
}

func (i StatikSource) Resource() Resource {
	return i.res
}

func (i StatikSource) Load() (kvs []*config.KeyValue, err error) {
	res := i.Resource()

	data := res.ReadBytes()

	kv := &config.KeyValue{
		Key:    res.Url(),
		Format: filepath.Ext(res.Path()),
		Value:  []byte(data),
	}
	if err != nil {
		return nil, err
	}
	return []*config.KeyValue{kv}, nil
}

func (i StatikSource) Watch() (config.Watcher, error) {
	return nil, nil //newWatcher(f)
}
