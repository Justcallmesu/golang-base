package images

type ImageResolution int

const (
	SMALL  = 640
	MEDIUM = 1280
	BIG    = 1920
)

type ImageResolutionWriteReturn struct {
	Resolution ImageResolution
	path       string
}
