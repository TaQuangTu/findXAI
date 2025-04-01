package system

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
