package buildinfo

import "runtime/debug"

var BuildInfo *debug.BuildInfo

func init() {
	var ok bool
	BuildInfo, ok = debug.ReadBuildInfo()
	if !ok {
		panic("cannot read build info")
	}
}

func DepsVersion(path string) string {
	for _, d := range BuildInfo.Deps {
		if d.Path == path {
			return d.Version
		}
	}
	return ""
}
