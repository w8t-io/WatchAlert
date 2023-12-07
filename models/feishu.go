package models

// FeiShuMsg 飞书
type FeiShuMsg struct {
	MsgType string `json:"msg_type"`
	Card    Cards  `json:"card"`
}

type Cards struct {
	Config   Configs    `json:"config"`
	Elements []Elements `json:"elements"`
	Header   Headers    `json:"header"`
	Tag      string     `json:"tag"`
}

type Actions struct {
	Tag      string      `json:"tag,omitempty"`
	Text     ActionsText `json:"text,omitempty"`
	Type     string      `json:"type,omitempty"`
	Value    interface{} `json:"value,omitempty"`
	Confirm  Confirms    `json:"confirm,omitempty"`
	URL      string      `json:"url,omitempty"`
	MultiURL *MultiURLs  `json:"multi_url,omitempty"`
}

type MultiURLs struct {
	URL        string `json:"url,omitempty"`
	AndroidURL string `json:"android_url,omitempty"`
	IosURL     string `json:"ios_url,omitempty"`
	PcURL      string `json:"pc_url,omitempty"`
}

type Confirms struct {
	Title Titles `json:"title"`
	Text  Texts  `json:"text"`
}

type ActionsText struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Configs struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

type Elements struct {
	Tag            string             `json:"tag"`
	FlexMode       string             `json:"flexMode"`
	BackgroupStyle string             `json:"background_style"`
	Text           Texts              `json:"text"`
	Actions        []Actions          `json:"actions"`
	Columns        []Columns          `json:"columns"`
	Elements       []ElementsElements `json:"elements"`
}

type ElementsElements struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
}

type Columns struct {
	Tag           string            `json:"tag"`
	Width         string            `json:"width"`
	Weight        int64             `json:"weight"`
	VerticalAlign string            `json:"vertical_align"`
	Elements      []ColumnsElements `json:"elements"`
}

type ColumnsElements struct {
	Tag  string `json:"tag"`
	Text Texts  `json:"text"`
}

type Texts struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Headers struct {
	Template string `json:"template"`
	Title    Titles `json:"title"`
}

type Titles struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

// CardInfo 飞书回传
type CardInfo struct {
	OpenID        string         `json:"open_id"`
	UserID        string         `json:"user_id"`
	OpenMessageID string         `json:"open_message_id"`
	OpenChatID    string         `json:"open_chat_id"`
	TenantKey     string         `json:"tenant_key"`
	Token         string         `json:"token"`
	Action        CardInfoAction `json:"action"`
}

type CardInfoAction struct {
	Value interface{} `json:"value"`
	Tag   string      `json:"tag"`
}

// FeiShuUserInfo 飞书用户信息
type FeiShuUserInfo struct {
	Data Data `json:"data"`
}

type Data struct {
	User User `json:"user"`
}

type User struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
}
