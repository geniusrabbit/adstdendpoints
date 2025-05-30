package dynamic

//easyjson:json
type tracker struct {
	Clicks      []string `json:"clicks,omitempty"`
	Impressions []string `json:"impressions,omitempty"`
	Views       []string `json:"views,omitempty"`
}

type assetThumb struct {
	Path   string `json:"path"`
	Type   string `json:"type,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

//easyjson:json
type asset struct {
	Name   string       `json:"name,omitempty"`
	Path   string       `json:"path"`
	Type   string       `json:"type,omitempty"`
	Width  int          `json:"width,omitempty"`
	Height int          `json:"height,omitempty"`
	Thumbs []assetThumb `json:"thumbs,omitempty"`
}

//easyjson:json
type item struct {
	ID         any            `json:"id"`
	Type       string         `json:"type"`
	URL        string         `json:"url,omitempty"`
	Content    string         `json:"content,omitempty"`
	ContentURL string         `json:"content_url,omitempty"`
	Fields     map[string]any `json:"fields,omitempty"`
	Assets     []asset        `json:"assets,omitempty"`
	Tracker    tracker        `json:"tracker"`
	Debug      any            `json:"debug,omitempty"`
}

//easyjson:json
type group struct {
	ID            string  `json:"id"`
	CustomTracker tracker `json:"custom_tracker,omitempty"`
	Items         []*item `json:"items"`
}

func (g *group) addItem(i *item) *group {
	g.Items = append(g.Items, i)
	return g
}

// Response object description
//
//easyjson:json
type Response struct {
	Version       string   `json:"version"`
	CustomTracker tracker  `json:"custom_tracker,omitempty"`
	Groups        []*group `json:"groups,omitempty"`
	Debug         any      `json:"debug,omitempty"`
}

func (r *Response) getGroupOrCreate(groupID string) *group {
	for _, g := range r.Groups {
		if g.ID == groupID {
			return g
		}
	}
	g := &group{ID: groupID}
	r.Groups = append(r.Groups, g)
	return g
}
