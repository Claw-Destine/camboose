package components

import "log/slog"

const paramCurrentProjec = "currentProject"

func appendQueryParams(basePath string, params ...string) string {
	if len(params)%2 != 0 {
		slog.Warn("List of params for this method should be even")
	}
	ret := basePath
	separator := "?"
	i := 0
	for {
		if i+2 > len(params) {
			break
		}
		ret = ret + separator + params[i] + "=" + params[i+1]
		separator = "&"
		i = i + 2
	}
	return ret
}
