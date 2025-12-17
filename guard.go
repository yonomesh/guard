package guard

type Config struct {
	Logging *Logging
}

// App is a thing that Guard Runs
type App interface {
	Start()
	Stop()
}

var CustomVersion string = "v0.0.0"

func Version() (simple, full string) {
	return "v0.0.1", "v0.0.1"
}
