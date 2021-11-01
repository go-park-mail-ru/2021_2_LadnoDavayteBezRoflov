package closer

import "backendServer/pkg/logger"

type Closer struct {
	logger *logger.Logger
}

func CreateCloser(logger *logger.Logger) *Closer {
	return &Closer{logger: logger}
}

func (closer *Closer) Close(closeFunc func() error) {
	err := closeFunc()
	if err != nil {
		closer.logger.Error(err)
	}
}
