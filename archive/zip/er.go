package zip

func er(f func() error, oldErr *error) {
	err := f()
	if *oldErr == nil {
		*oldErr = err
	}
}
