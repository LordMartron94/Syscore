//go:build !linux && !windows && !darwin

package platform

/*
DisplayAPIDetect is unsupported on this GOOS and always returns an error.
*/
func DisplayAPIDetect() (DisplayAPI, error) {
	return DisplayAPIInvalid, displayAPIUnsupportedError()
}

/*
DisplayAPIIsAvailable always returns false on unsupported GOOS values.
*/
func DisplayAPIIsAvailable(_ DisplayAPI) bool {
	return false
}

/*
DisplayAPIListAvailable returns nil on unsupported GOOS values.
*/
func DisplayAPIListAvailable() []DisplayAPI {
	return nil
}
