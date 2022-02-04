package res

import (
	"github.com/go-kratos/kratos/v2/config"

	_ "github.com/qiangyt/jog/res/statik"
)

type DummyWatcher struct {
}

func (i *DummyWatcher) Next() ([]*config.KeyValue, error) {
	return nil, nil
}

func (i *DummyWatcher) Stop() error {
	return nil
}

func (i StatikSource) Watch() (config.Watcher, error) {
	return &DummyWatcher{}, nil
}
