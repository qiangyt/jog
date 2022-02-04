package res

import (
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"

	_ "github.com/qiangyt/jog/res/statik"
)

var _ config.Source = (*StatikSourceT)(nil)

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
		Format: filepath.Ext(res.Path())[1:],
		Value:  []byte(data),
	}
	if err != nil {
		return nil, err
	}
	return []*config.KeyValue{kv}, nil
}
