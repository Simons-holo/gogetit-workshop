package scraper

type WebsiteMetadata struct {
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Keywords       []string          `json:"keywords"`
	CanonicalURL   string            `json:"canonical_url"`
	OGTitle        string            `json:"og_title"`
	OGDescription  string            `json:"og_description"`
	OGImage        string            `json:"og_image"`
	OGURL          string            `json:"og_url"`
	OGType         string            `json:"og_type"`
	TwitterCard    string            `json:"twitter_card"`
	TwitterTitle   string            `json:"twitter_title"`
	TwitterImage   string            `json:"twitter_image"`
	Links          []LinkMetadata    `json:"links"`
	Images         []ImageMetadata   `json:"images"`
	MetaTags       map[string]string `json:"meta_tags"`
	StructuredData []string          `json:"structured_data"`
	ResponseTime   int64             `json:"response_time_ms"`
	StatusCode     int               `json:"status_code"`
	ContentLength  int64             `json:"content_length"`
	ContentType    string            `json:"content_type"`
}

type LinkMetadata struct {
	Href       string `json:"href"`
	Text       string `json:"text"`
	Rel        string `json:"rel"`
	Target     string `json:"target"`
	IsExternal bool   `json:"is_external"`
}

type ImageMetadata struct {
	Src    string `json:"src"`
	Alt    string `json:"alt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func NewWebsiteMetadata() *WebsiteMetadata {
	return &WebsiteMetadata{
		Keywords:       []string{},
		Links:          []LinkMetadata{},
		Images:         []ImageMetadata{},
		MetaTags:       make(map[string]string),
		StructuredData: []string{},
	}
}
