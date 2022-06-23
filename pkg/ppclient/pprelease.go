package ppclient

import "strings"

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
	return p.Phase == "Planning / Development / Testing"
}

func (p *ppRelease) IsZStream() bool {
	return strings.HasSuffix(p.GetVersion(), ".z")
}
