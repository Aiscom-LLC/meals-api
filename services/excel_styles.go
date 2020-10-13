package services

type Alignment struct {
	Horizontal      string `json:"horizontal"`
	Indent          int    `json:"indent"`
	JustifyLastLine bool   `json:"justify_last_line"`
	ReadingOrder    uint64 `json:"reading_order"`
	RelativeIndent  int    `json:"relative_indent"`
	ShrinkToFit     bool   `json:"shrink_to_fit"`
	TextRotation    int    `json:"text_rotation"`
	Vertical        string `json:"vertical"`
	WrapText        bool   `json:"wrap_text"`
}

type Border struct {
	Type  string `json:"type"`
	Color string `json:"color"`
	Style int    `json:"style"`
}

type Font struct {
	Bold      bool    `json:"bold"`
	Italic    bool    `json:"italic"`
	Underline string  `json:"underline"`
	Family    string  `json:"family"`
	Size      float64 `json:"size"`
	Strike    bool    `json:"strike"`
	Color     string  `json:"color"`
}

type Fill struct {
	Type    string   `json:"type"`
	Pattern int      `json:"pattern"`
	Color   []string `json:"color"`
	Shading int      `json:"shading"`
}

type Protection struct {
	Hidden bool `json:"hidden"`
	Locked bool `json:"locked"`
}

// Style struct
type Style struct {
	Border        []Border    `json:"border"`
	Fill          Fill        `json:"fill"`
	Font          *Font       `json:"font"`
	Alignment     *Alignment  `json:"alignment"`
	Protection    *Protection `json:"protection"`
	NumFmt        int         `json:"number_format"`
	DecimalPlaces int         `json:"decimal_places"`
	CustomNumFmt  *string     `json:"custom_number_format"`
	Lang          string      `json:"lang"`
	NegRed        bool        `json:"negred"`
}
