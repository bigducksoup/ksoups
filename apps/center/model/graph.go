package model

import "encoding/json"

type GraphData []Element

func UnmarshalGraphData(data []byte) (GraphData, error) {
	var r GraphData
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *GraphData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Element struct {
	Shape    string    `json:"shape"`
	Attrs    *Attrs    `json:"attrs,omitempty"`
	ID       string    `json:"id"`
	Markup   []Markup  `json:"markup,omitempty"`
	ZIndex   int64     `json:"zIndex"`
	Router   *Router   `json:"router,omitempty"`
	Source   *Source   `json:"source,omitempty"`
	Target   *Source   `json:"target,omitempty"`
	Position *Position `json:"position,omitempty"`
	Size     *Size     `json:"size,omitempty"`
	View     *string   `json:"view,omitempty"`
	Ports    *Ports    `json:"ports,omitempty"`
	Data     *Data     `json:"data,omitempty"`
}

type Attrs struct {
	Line     Line     `json:"line"`
	P1       P1       `json:"p1"`
	P2       P2       `json:"p2"`
	Sign     Sign     `json:"sign"`
	SignText SignText `json:"signText"`
	C2       C2       `json:"c2"`
}

type C2 struct {
	R                 int64   `json:"r"`
	Stroke            string  `json:"stroke"`
	Fill              string  `json:"fill"`
	AtConnectionRatio float64 `json:"atConnectionRatio"`
	StrokeWidth       int64   `json:"strokeWidth"`
	Cursor            string  `json:"cursor"`
	Event             string  `json:"event"`
}

type Line struct {
	Stroke          string `json:"stroke"`
	TargetMarker    string `json:"targetMarker"`
	StrokeDasharray string `json:"strokeDasharray"`
}

type P1 struct {
	Connection     bool   `json:"connection"`
	Fill           string `json:"fill"`
	Stroke         string `json:"stroke"`
	StrokeWidth    int64  `json:"strokeWidth"`
	StrokeLinejoin string `json:"strokeLinejoin"`
}

type P2 struct {
	Connection     bool         `json:"connection"`
	Fill           string       `json:"fill"`
	Stroke         string       `json:"stroke"`
	StrokeWidth    int64        `json:"strokeWidth"`
	PointerEvents  string       `json:"pointerEvents"`
	StrokeLinejoin string       `json:"strokeLinejoin"`
	TargetMarker   TargetMarker `json:"targetMarker"`
}

type TargetMarker struct {
	TagName     string `json:"tagName"`
	Fill        string `json:"fill"`
	Stroke      string `json:"stroke"`
	StrokeWidth int64  `json:"strokeWidth"`
	D           string `json:"d"`
}

type Sign struct {
	X                  int64  `json:"x"`
	Y                  int64  `json:"y"`
	Width              int64  `json:"width"`
	Height             int64  `json:"height"`
	Stroke             string `json:"stroke"`
	Fill               string `json:"fill"`
	AtConnectionLength int64  `json:"atConnectionLength"`
	StrokeWidth        int64  `json:"strokeWidth"`
	Event              string `json:"event"`
	Cursor             string `json:"cursor"`
}

type SignText struct {
	AtConnectionLength int64  `json:"atConnectionLength"`
	TextAnchor         string `json:"textAnchor"`
	TextVerticalAnchor string `json:"textVerticalAnchor"`
	Text               string `json:"text"`
	Event              string `json:"event"`
	Cursor             string `json:"cursor"`
	FontSize           int64  `json:"fontSize"`
}

type Data struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Proto Proto  `json:"proto"`
	Root  *bool  `json:"root,omitempty"`
}

type Proto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int64  `json:"type"`
	CreateTime  string `json:"createTime"`
	Timeout     int64  `json:"timeout"`
	JustRun     bool   `json:"justRun"`
	Payload     string `json:"payload"`
	ProbeID     string `json:"probeId"`
}

type Markup struct {
	TagName  string `json:"tagName"`
	Selector string `json:"selector"`
}

type Ports struct {
	Groups Groups `json:"groups"`
	Items  []Item `json:"items"`
}

type Groups struct {
	Top    Bottom `json:"top"`
	Bottom Bottom `json:"bottom"`
}

type Bottom struct {
	Position string      `json:"position"`
	Attrs    BottomAttrs `json:"attrs"`
}

type BottomAttrs struct {
	Circle Circle `json:"circle"`
}

type Circle struct {
	R           int64  `json:"r"`
	Magnet      bool   `json:"magnet"`
	Stroke      string `json:"stroke"`
	StrokeWidth int64  `json:"strokeWidth"`
	Fill        string `json:"fill"`
}

type Item struct {
	ID    string `json:"id"`
	Group string `json:"group"`
}

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Router struct {
	Name string `json:"name"`
	Args Args   `json:"args"`
}

type Args struct {
}

type Size struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

type Source struct {
	Cell string `json:"cell"`
	Port string `json:"port"`
}
