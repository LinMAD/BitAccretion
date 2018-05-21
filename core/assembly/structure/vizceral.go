package structure

type (
	// VClassType tag to indicate health type of node
	VClassType string
	// VLayout type of layout to draw position of nodes (ring, chain etc)
	VLayout string
	// VRenderer type of graph lvl
	VRenderer string
)

var (
	// VNormal tag
	VNormal VClassType = "normal"
	// VDanger tag
	VDanger VClassType = "danger"
	// VWarning tag
	VWarning VClassType = "warning"

	// VGlobalRenderer main graph style
	VGlobalRenderer VRenderer = "global"
	// VRegionRenderer main graph style
	VRegionRenderer VRenderer = "region"
	// VFocusedRenderer main graph style
	VFocusedRenderer VRenderer = "focused"
	// VFocusedChildRenderer main graph style
	VFocusedChildRenderer VRenderer = "focusedChild"

	// VDNSLayout connection style
	VDNSLayout VLayout = "dns"
	// VLTRTreeLayout connection style
	VLTRTreeLayout VLayout = "ltrTree"
	// VRingCenterLayout connection style
	VRingCenterLayout VLayout = "ringCenter"
	// VRingLayout connection style
	VRingLayout VLayout = "ring"
)

// VRegionGraph it's a global structure, it can be wrapped if has more than one decenter for example
type VRegionGraph struct {
	Renderer    VRenderer         `json:"renderer"`
	Name        string            `json:"name"`
	MaxVolume   float64           `json:"maxVolume,omitempty"`
	Nodes       []VNode           `json:"nodes,omitempty"`
	Connections []VNodeConnection `json:"connections,omitempty"`
}

// VNode structure it's keeps all nodes and connections of them
type VNode struct {
	Renderer      VRenderer         `json:"renderer,omitempty"`
	Name          string            `json:"name,omitempty"`
	DisplayName   string            `json:"displayName,omitempty"`
	MaxVolume     float64           `json:"maxVolume,omitempty"`
	Updated       int64             `json:"updated,omitempty"`
	Nodes         []VNode           `json:"nodes,omitempty"`
	Connections   []VNodeConnection `json:"connections,omitempty"`
	Notices       []VNotice         `json:"notices,omitempty"`
	Class         VClassType        `json:"class,omitempty"`
	Metadata      VMeta             `json:"metadata,omitempty"`
	EntryNode     string            `json:"entryNode,omitempty"`
	Layout        VLayout           `json:"layout,omitempty"`
	SystemDetails interface{}       `json:"-"`
}

// VNodeConnection represents connection between all nodes
type VNodeConnection struct {
	// Source name represents node name
	Source string `json:"source,omitempty"`
	// Target name represents connection to the source node
	Target  string        `json:"target,omitempty"`
	Metrics VMetricLevels `json:"metrics,omitempty"`
	Notices []VNotice     `json:"notices,omitempty"`
	Class   VClassType    `json:"class,omitempty"`
}

// VMetricLevels node's request metrics
type VMetricLevels struct {
	Danger  float32 `json:"danger"`
	Warning float32 `json:"warning"`
	Normal  float32 `json:"normal"`
}

// VMeta data an external fields can be passed
type VMeta struct {
	Streaming int `json:"streaming,omitempty"`
}

// VNotice some external notice data
type VNotice struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
	Link     string `json:"link,omitempty"`
	Severity int    `json:"severity,omitempty"`
}
