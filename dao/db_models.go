package dao

// User struct for acess to users table
type User struct {
	Id                int64    `gorm:"column:Id;primaryKey;AUTO_INCREMENT"`
	Name              string   `gorm:"column:Name"`
	Password          string   `gorm:"column:Password" json:"-"`
	VkLogin           string   `gorm:"column:VkLogin"`
	VkPassword        string   `gorm:"column:VkPassword"`
	TwitterScreenName string   `gorm:"column:TwitterScreenName"`
	VkNewsEnabled     bool     `gorm:"column:VkNewsEnabled"`
	Settings          Settings `gorm:"ForeignKey:UserId"`
}

// TableName return table name for User struct
func (User) TableName() string {
	return "users"
}

// Settings struct for acess to settings table
type Settings struct {
	Id                   int64 `gorm:"column:Id;primary_key;AUTO_INCREMENT"`
	UserId               int64 `gorm:"column:UserId;index"`
	UnreadOnly           bool  `gorm:"column:UnreadOnly"`
	MarkSameRead         bool  `gorm:"column:MarkSameRead"`
	RssEnabled           bool  `gorm:"column:RssEnabled"`
	VkNewsEnabled        bool  `gorm:"column:VkNewsEnabled"`
	TwitterEnabled       bool  `gorm:"column:TwitterEnabled"`
	TwitterSimpleVersion bool  `gorm:"column:TwitterSimpleVersion"`
	ShowPreviewButton    bool  `gorm:"column:ShowPreviewButton"`
	ShowTabButton        bool  `gorm:"column:ShowTabButton"`
	ShowReadButton       bool  `gorm:"column:ShowReadButton"`
	ShowLinkButton       bool  `gorm:"column:ShowLinkButton"`
	ShowBookmarkButton   bool  `gorm:"column:ShowBookmarkButton"`
}

// TableName return table name for Settings struct
func (Settings) TableName() string {
	return "settings"
}

// TwitterNews struct for acess to twitternews table
type TwitterNews struct {
	Id          int64  `gorm:"column:Id;primary_key;AUTO_INCREMENT"`
	TweetId     int64  `gorm:"column:TweetId;index"`
	UserId      int64  `gorm:"column:UserId;index"`
	SourceId    int64  `gorm:"column:SourceId;index"`
	CreatedAt   int64  `gorm:"column:CreatedAt"`
	Text        string `gorm:"column:Text"`
	ExpandedUrl string `gorm:"column:ExpandedUrl"`
	Image       string `gorm:"column:Image"`
}

// TableName return table name for TwitterNews struct
func (TwitterNews) TableName() string {
	return "twitternews"
}

// TwitterSource struct for acess to twittersource table
type TwitterSource struct {
	Id         int64  `gorm:"column:Id;primary_key;AUTO_INCREMENT"`
	UserId     int64  `gorm:"column:UserId;index"`
	Name       string `gorm:"column:Name"`
	ScreenName string `gorm:"column:ScreenName"`
	Url        string `gorm:"column:Url"`
	Image      string `gorm:"column:Image"`
}

// TableName return table name for TwitterSource struct
func (TwitterSource) TableName() string {
	return "twittersource"
}
