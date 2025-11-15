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
type itemMetaAdvertiserInfo struct {
	ID         uint64 `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	AboutURL   string `json:"about_url,omitempty"`
	ContactURL string `json:"contact_url,omitempty"`
	PrivacyURL string `json:"privacy_url,omitempty"`
	TermsURL   string `json:"terms_url,omitempty"`
}

//easyjson:json
type itemMetaAdInfo struct {
	ID          uint64 `json:"id,omitempty"`
	CampaignID  uint64 `json:"campaign_id,omitempty"`
	Description string `json:"description,omitempty"`
	MinAge      int    `json:"min_age,omitempty"`
	AboutURL    string `json:"about_url,omitempty"`
	ContactURL  string `json:"contact_url,omitempty"`
	PrivacyURL  string `json:"privacy_url,omitempty"`
	TermsURL    string `json:"terms_url,omitempty"`
}

//easyjson:json
type itemMetaInfoHide struct {
	Type            string            `json:"type,omitempty"`               // 'cookie', 'url_param', 'script'
	HideAdURLType   string            `json:"hide_ad_url_type,omitempty"`   // 'get', 'post', 'pixel'
	HideAdURL       string            `json:"hide_ad_url,omitempty"`        // URL where user can hide the ad
	HideAdURLParams map[string]string `json:"hide_ad_url_params,omitempty"` // Additional params for hide ad URL
	Name            string            `json:"name,omitempty"`               // name of cookie or url param
	ScriptURL       string            `json:"script_url,omitempty"`         // URL of script which perform hiding
	Script          string            `json:"script,omitempty"`             // Executable script which perform hiding
}

//easyjson:json
type itemMetaInfo struct {
	Advertiser *itemMetaAdvertiserInfo `json:"advertiser,omitempty"`
	Ad         *itemMetaAdInfo         `json:"ad,omitempty"`
	// ComplaintURL it's URL where user can send complaint about the Ad
	ComplaintURL string `json:"complaint_url,omitempty"`
	// Hide it's info about hiding the Ad
	Hide *itemMetaInfoHide `json:"hide,omitempty"`
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
	Meta       *itemMetaInfo  `json:"meta,omitempty"`
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
