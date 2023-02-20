package image

// Image is the docker image client interface
type Image interface {
	Build(cfg *BuildConfig) error
	Pull(cfg *PullConfig) error
	Push(cfg *PushConfig) error
}

type image struct {
}

// New creates a docker image client
func New() Image {
	return &image{}
}

func (i *image) Build(cfg *BuildConfig) error {
	return Build(cfg)
}

func (i *image) Pull(cfg *PullConfig) error {
	return Pull(cfg)
}

func (i *image) Push(cfg *PushConfig) error {
	return Push(cfg)
}
