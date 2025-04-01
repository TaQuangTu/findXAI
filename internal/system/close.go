package system

type RootCloser func() error

var (
	rootClosers = make([]RootCloser, 0)
)

func RegisterRootCloser(closer RootCloser) {
	rootClosers = append(rootClosers, closer)
}

func SafeClose() {
	for _, closer := range rootClosers {
		_ = closer()
	}
}
