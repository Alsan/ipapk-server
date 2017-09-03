package models

type VersionInfo []vInfo

type vInfo struct {
	Version string
	Builds  int
}

func (v VersionInfo) Len() int {
	return len(v)
}

func (v VersionInfo) Less(i, j int) bool {
	return v[i].Version > v[j].Version
}

func (v VersionInfo) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
