//go:build windows

package platform

/*
DisplayAPIDetect always returns DisplayAPIWin32 on Windows.
*/
func DisplayAPIDetect() (DisplayAPI, error) {
	return DisplayAPIWin32, nil
}

/*
DisplayAPIIsAvailable returns true only for DisplayAPIWin32 on Windows.
*/
func DisplayAPIIsAvailable(api DisplayAPI) bool {
	return api == DisplayAPIWin32
}

/*
DisplayAPIListAvailable returns a single-element slice containing DisplayAPIWin32.
*/
func DisplayAPIListAvailable() []DisplayAPI {
	return []DisplayAPI{DisplayAPIWin32}
}
