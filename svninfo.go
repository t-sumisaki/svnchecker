package svnchecker

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type SVNInfo struct {
	Path                string
	Name                string
	WorkingCopyRootPath string
	URL                 string
	RelativeURL         string
	RepositoryRoot      string
	RepositoryUUID      string
	Revision            string
	NodeKind            string
	Schedule            string
	LastChangedAuthor   string
	LastChangedRev      string
	LastChangedDate     string
	Checksum            string
	LockToken           string
	LockOwner           string
}

func (t SVNInfo) IsLocked() bool {
	return t.LockOwner != ""
}

const (
	SVN_Path                = "Path"
	SVN_Name                = "Name"
	SVN_WorkingCopyRootPath = "Working Copy Root Path"
	SVN_URL                 = "URL"
	SVN_RelativeURL         = "Relative URL"
	SVN_RepositoryRoot      = "Repository Root"
	SVN_RepositoryUUID      = "Repository UUID"
	SVN_Revision            = "Revision"
	SVN_NodeKind            = "Node Kind"
	SVN_Schedule            = "Schedule"
	SVN_LastChangedAuthor   = "Last Changed Author"
	SVN_LastChangedRev      = "Last Changed Rev"
	SVN_LastChangedDate     = "Last Changed Date"
	SVN_Checksum            = "Checksum"
	SVN_LockToken           = "Lock Token"
	SVN_LockOwner           = "Lock Owner"
)

func parseSVNInfo(src string) (*SVNInfo, error) {
	info := &SVNInfo{}
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")

	for _, l := range lines {
		elems := strings.SplitN(l, ": ", 2)

		if len(elems) != 2 {
			continue
		}

		key, value := elems[0], strings.TrimSpace(elems[1])

		switch key {
		case SVN_Path:
			info.Path = value
		case SVN_WorkingCopyRootPath:
			info.WorkingCopyRootPath = value
		case SVN_URL:
			info.URL = value
		case SVN_RelativeURL:
			info.RelativeURL = value
		case SVN_RepositoryRoot:
			info.RepositoryRoot = value
		case SVN_RepositoryUUID:
			info.RepositoryUUID = value
		case SVN_Revision:
			info.Revision = value
		case SVN_NodeKind:
			info.NodeKind = value
		case SVN_Schedule:
			info.Schedule = value
		case SVN_LastChangedAuthor:
			info.LastChangedAuthor = value
		case SVN_LastChangedRev:
			info.LastChangedRev = value
		case info.LastChangedDate:
			info.LastChangedDate = value
		}
	}

	return info, nil

}

func GetInfo(path string) (*SVNInfo, error) {

	out, err := exec.Command("svn", "info", path).Output()
	if err != nil {
		return nil, fmt.Errorf("svn info command error: %w", err)
	}

	buf := bytes.NewBuffer(out)
	b, err := io.ReadAll(transform.NewReader(buf, japanese.ShiftJIS.NewDecoder()))

	if err != nil {
		return nil, fmt.Errorf("svn info result read error: %w", err)
	}

	info, err := parseSVNInfo(string(b))

	if err != nil {
		return nil, fmt.Errorf("svn info result parse error: %w", err)
	}

	return info, nil

}
