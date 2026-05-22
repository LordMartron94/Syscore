//go:build darwin

package platform

func DisplayAPIDetect() (DisplayAPI, error) {
	return DisplayAPIAppkit, nil
}

func DisplayAPIIsAvailable(api DisplayAPI) bool {
	return api == DisplayAPIAppkit
}

func DisplayAPIListAvailable() []DisplayAPI {
	return []DisplayAPI{DisplayAPIAppkit}
}
