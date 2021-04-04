package dao

// Config struct for app configuration
type Config struct {
	Driver           string `json:"driver"`
	ConnectionString string `json:"connection_string"`
	DbHost           string `json:"db_host"`
	DbName           string `json:"db_name"`
	DbUser           string `json:"db_user"`
	DbPassword       string `json:"db_password"`
	DbPort           int    `json:"db_port"`
	ClientSecret     string `json:"client_secret"`
	ApiKey           string `json:"api_key"`
	UpdateMinutes    int    `json:"update_minutes"`
}

// TwitterFriends struct for user friends in Twitter account
type TwitterFriends struct {
	Users []TwitterUser `json:"users"`
}

// TwitterUser name for user from twitter api
type TwitterUser struct {
	ScreenName string `json:"screen_name"`
}

// TweetJson struct for tweets
type TweetJson struct {
	Id       int64             `json:"id"`
	Text     string            `json:"text"`
	Source   TweetSourceJson   `json:"user"`
	Entities TweetEntitiesJson `json:"entities"`
}

// TweetSourceJson struct for twitter sources for news
type TweetSourceJson struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Url        string `json:"url"`
	Image      string `json:"profile_image_url_https"`
}

// TweetEntitiesJson struct for additional information in tweet
type TweetEntitiesJson struct {
	Urls  []TweetEntitiesUrlsJson  `json:"urls"`
	Media []TweetEntitiesMediaJson `json:"media"`
}

// TweetEntitiesUrlsJson struct for external url in tweet
type TweetEntitiesUrlsJson struct {
	ExpandedUrl string `json:"expanded_url"`
}

type TweetEntitiesMediaJson struct {
	ExpandedUrl   string `json:"expanded_url"`
	MediaUrlHttps string `json:"media_url_https"`
}

// TwitterToken towen pair for auth in twitter
type TwitterToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// ToDbSource convers to DB struct
func (t *TweetJson) ToDbSource() TwitterSource {
	source := TwitterSource{
		Name:       t.Source.Name,
		ScreenName: t.Source.ScreenName,
		Url:        t.Source.Url,
		Image:      t.Source.Image,
	}

	return source
}

// ToDbSource convers to DB struct
func (t *TweetJson) ToDbNews() TwitterNews {
	news := TwitterNews{
		Text:    t.Text,
		TweetId: t.Id,
	}
	for _, urlsJson := range t.Entities.Urls {
		if urlsJson.ExpandedUrl != "" {
			news.ExpandedUrl = urlsJson.ExpandedUrl
			break
		}
	}
	for _, media := range t.Entities.Media {
		if media.MediaUrlHttps != "" {
			news.Image = media.MediaUrlHttps
			break
		}
	}

	return news
}
