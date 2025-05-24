package domain

const (
	VisibilityPublic  = "Public"
	VisibilityPrivate = "Private"
)

type Playlist struct {
	Name        string
	Description string
	Visibility  string
}

func (p *Playlist) IsPublic() bool {
	return p.Visibility == VisibilityPublic
}
