package unleash

import (
	"net/http"

	"github.com/Unleash/unleash-client-go/v3"
)

var (
	enableFeatureToogle bool
)

func RegisterFeatureToggle(appName, apiToken, apiUrl string) error {
	enableFeatureToogle = true
	err := unleash.Initialize(
		unleash.WithListener(unleash.DefaultStorage{}),
		// unleash.WithListener(unleash.DebugListener{}),
		unleash.WithAppName(appName),
		unleash.WithUrl(apiUrl),
		unleash.WithCustomHeaders(http.Header{
			"Authorization": {apiToken},
		}),
	)
	if err != nil {
		return err
	}
	return nil
}

func IsEnabled(key string) bool {
	if !enableFeatureToogle {
		return true
	}
	return unleash.IsEnabled(key)
}

func MustEnabled(key string) {
	if !enableFeatureToogle {
		return
	}
	if !unleash.IsEnabled(key) {
		panic("Underdevelopment")
	}
}
