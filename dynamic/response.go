package dynamic

//easyjson:json
type adSourceInfo struct {
	ID          string         `json:"id"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Domain      string         `json:"domain,omitempty"`
	IconURL     string         `json:"icon_url,omitempty"`
	LogoURL     string         `json:"logo_url,omitempty"`
	URL         string         `json:"url,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

//easyjson:json
type itemMetaActionInfo struct {
	Type        string `json:"type,omitempty"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

//easyjson:json
type tracker struct {
	Clicks      []string `json:"clicks,omitempty"`
	Impressions []string `json:"impressions,omitempty"`
	Views       []string `json:"views,omitempty"`
}

type adAssetThumb struct {
	Path   string `json:"path"`
	Type   string `json:"type,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

//easyjson:json
type adAsset struct {
	Name     string         `json:"name,omitempty"`
	Path     string         `json:"path"`
	Type     string         `json:"type,omitempty"`
	Width    int            `json:"width,omitempty"`
	Height   int            `json:"height,omitempty"`
	Duration int            `json:"duration,omitempty"` // Duration in seconds for video assets
	Thumbs   []adAssetThumb `json:"thumbs,omitempty"`
}

//easyjson:json
type itemMetaAdvertiserInfo struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	AboutURL   string `json:"about_url,omitempty"`
	ContactURL string `json:"contact_url,omitempty"`
	PrivacyURL string `json:"privacy_url,omitempty"`
	TermsURL   string `json:"terms_url,omitempty"`
}

//easyjson:json
type itemMetaAdInfo struct {
	ID          string `json:"id,omitempty"`
	CampaignID  string `json:"campaign_id,omitempty"`
	AdSourceID  string `json:"adsource_id,omitempty"`
	Description string `json:"description,omitempty"`
	MinAge      int    `json:"min_age,omitempty"`
	AboutURL    string `json:"about_url,omitempty"`
	ContactURL  string `json:"contact_url,omitempty"`
	PrivacyURL  string `json:"privacy_url,omitempty"`
	TermsURL    string `json:"terms_url,omitempty"`
}

//easyjson:json
type itemMetaInfo struct {
	Advertiser *itemMetaAdvertiserInfo `json:"advertiser,omitempty"`
	Ad         *itemMetaAdInfo         `json:"ad,omitempty"`
	Actions    []*itemMetaActionInfo   `json:"actions,omitempty"`
}

func (m *itemMetaInfo) addAction(actionType, title, description, url string) {
	m.Actions = append(m.Actions, &itemMetaActionInfo{
		Type:        actionType,
		Title:       title,
		Description: description,
		URL:         url,
	})
}

//easyjson:json
type item struct {
	ID       any            `json:"id"`
	Type     string         `json:"type"`
	URL      string         `json:"url,omitempty"`
	Fields   map[string]any `json:"fields,omitempty"`
	Assets   []adAsset      `json:"assets,omitempty"`
	Tracker  *tracker       `json:"tracker"`
	Metadata map[string]any `json:"metadata,omitempty"`
	AdInfo   *itemMetaInfo  `json:"adinfo,omitempty"`
	Debug    any            `json:"debug,omitempty"`
}

//easyjson:json
type group struct {
	ID            string  `json:"id"` // ID of the placement on the site (adzone, slot, unit, etc.)
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
	Version       string          `json:"version"`
	CustomTracker tracker         `json:"custom_tracker,omitempty"`
	Groups        []*group        `json:"groups,omitempty"`
	AdSources     []*adSourceInfo `json:"adsources,omitempty"`
	Debug         any             `json:"debug,omitempty"`
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

func (r *Response) hasSource(sourceID string) bool {
	for _, s := range r.AdSources {
		if s.ID == sourceID {
			return true
		}
	}
	return false
}
