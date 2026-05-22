//go:build !linux && !windows && !darwin

package platform

func DisplayAPIDetect() (DisplayAPI, error) {
	return DisplayAPIInvalid, displayAPIUnsupportedError()
}

func DisplayAPIIsAvailable(_ DisplayAPI) bool {
	return false
}

func DisplayAPIListAvailable() []DisplayAPI {
	return nil
}
