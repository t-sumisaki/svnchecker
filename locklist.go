package svnchecker

type LockInfo struct {
	User        string
	Path        string
	LastChanged string
}

func GetLockName(path string) string {
	return ""
}
