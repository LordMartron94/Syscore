//go:build darwin

package platform

/*
DisplayAPIDetect always returns DisplayAPIAppkit on macOS.
*/
func DisplayAPIDetect() (DisplayAPI, error) {
	return DisplayAPIAppkit, nil
}

/*
DisplayAPIIsAvailable returns true only for DisplayAPIAppkit on macOS.
*/
func DisplayAPIIsAvailable(api DisplayAPI) bool {
	return api == DisplayAPIAppkit
}

/*
DisplayAPIListAvailable returns a single-element slice containing DisplayAPIAppkit.
*/
func DisplayAPIListAvailable() []DisplayAPI {
	return []DisplayAPI{DisplayAPIAppkit}
}
