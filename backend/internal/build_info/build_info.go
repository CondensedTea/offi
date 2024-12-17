package build_info

import (
	"debug/buildinfo"
	"os"
)

var Version = "unknown"

func init() {
	f, err := os.Executable()
	if err != nil {
		return
	}

	info, err := buildinfo.ReadFile(f)
	if err != nil {
		return
	}

	for _, v := range info.Settings {
		if v.Key == "vcs.revision" && len(v.Value) > 7 {
			Version = v.Value[0:7]
			return
		}
	}

	return
}
