package ppclient

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type PpRelease interface {
	GetVersion() string
	GetShortname() string
	GetPhase() string
	IsUnsupported() bool
	IsMaintained() bool
	IsDeveloped() bool
	IsZStream() bool
}

type PpReleaseList []PpRelease

type ByVersion []PpRelease

type ppRelease struct {
	Shortname string `json:"shortname,omitempty"`
	Phase     string `json:"phase_display,omitempty"`
}

func NewPpRelease(shortname string, phase string) PpRelease {
	return &ppRelease{
		Shortname: shortname,
		Phase:     phase,
	}
}

func (p *ppRelease) GetVersion() string {
	s := strings.SplitN(p.Shortname, "-", 2)
	if len(s) < 2 {
		return ""
	}
	return s[1]
}

func (p *ppRelease) GetShortname() string {
	s := strings.SplitN(p.Shortname, "-", 2)
	if len(s) < 2 {
		return ""
	}
	return s[0]
}

func (p *ppRelease) GetPhase() string {
	return p.Phase
}

func (p *ppRelease) IsUnsupported() bool {
	return p.Phase == "Unsupported"
}

func (p *ppRelease) IsMaintained() bool {
	return p.Phase == "Maintenance"
}

func (p *ppRelease) IsDeveloped() bool {
	return (p.Phase == "Planning / Development / Testing" || p.Phase == "CI / CD")
}

func (p *ppRelease) IsZStream() bool {
	return strings.HasSuffix(p.GetVersion(), ".z")
}

func (l *PpReleaseList) GetLatestVersion(includeZStream bool) string {
	var i int
	sort.Sort(ByVersion(*l))
	for i = len(*l) - 1; i >= 0; i-- {
		if !includeZStream && (*l)[i].IsZStream() {
			continue
		}
		if !(*l)[i].IsDeveloped() {
			// Skip "Concept" releases
			continue
		}
		break
	}
	return (*l)[i].GetVersion()
}

func (r ByVersion) Len() int {
	return len(r)
}

func (r ByVersion) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func normalizeVersionArray(arr *[]string) error {
	if len(*arr) == 2 {
		*arr = append(*arr, "0")
	}
	for idx, num := range *arr {
		if num != "z" {
			intVal, err := strconv.Atoi((*arr)[idx])
			if err != nil {
				return err
			}
			(*arr)[idx] = fmt.Sprintf("%3d", intVal)
		}
	}
	return nil
}

func (r ByVersion) Less(i, j int) bool {
	versionArr1 := strings.SplitN(r[i].GetVersion(), ".", 3)
	versionArr2 := strings.SplitN(r[j].GetVersion(), ".", 3)
	if len(versionArr1) < 2 || len(versionArr1) > 3 {
		// Should never happen
		return false
	}
	err := normalizeVersionArray(&versionArr1)
	if err != nil {
		return false
	}
	err = normalizeVersionArray(&versionArr2)
	if err != nil {
		return false
	}
	return (versionArr1[0] < versionArr2[0]) ||
		((versionArr1[0] == versionArr2[0]) && (versionArr1[1] < versionArr2[1])) ||
		((versionArr1[0] == versionArr2[0] && (versionArr1[1] == versionArr2[1])) && (versionArr1[2] < versionArr2[2]))
}
