//go:build windows

package platform

func DisplayAPIDetect() (DisplayAPI, error) {
	return DisplayAPIWin32, nil
}

func DisplayAPIIsAvailable(api DisplayAPI) bool {
	return api == DisplayAPIWin32
}

func DisplayAPIListAvailable() []DisplayAPI {
	return []DisplayAPI{DisplayAPIWin32}
}
