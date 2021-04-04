package service

import (
	"log"
	"time"

	"newshub-twitter-service/dao"
)

var (
	updChan   = make(chan dao.User, 4)
	closeChan = make(chan bool, 1)
)

func Update() {
	users := getUsers()

	for _, user := range users {
		updChan <- user
	}
}

func Close() {
	closeChan <- true
}

func listener() {
	for {
		select {
		case user := <-updChan:
			go saveNews(user)
		case <-closeChan:
			log.Println("updater stopped")
			close(updChan)
			close(closeChan)

			return
		}
	}
}

func saveNews(user dao.User) {
	if user.TwitterScreenName == "" {
		return
	}

	sources := getSources(user.Id)
	tweetIds := getLastPostId(sources)
	tweets := getNews(user.TwitterScreenName)

	for _, tweet := range tweets {
		if _, ok := sources[tweet.Source.Url]; !ok {
			addSource(user.Id, tweet)
			sources = getSources(user.Id)
			tweetIds = getLastPostId(sources)
		}
		if _, ok := tweetIds[tweet.Id]; !ok {
			addNews(sources, tweet)
			tweetIds[tweet.Id] = true
		}
	}
}

func addSource(userId int64, data dao.TweetJson) {
	source := data.ToDbSource()
	source.UserId = userId

	if err := db.Save(&source).Error; err != nil {
		log.Println("save source err:", err)
	}
}

func addNews(sources map[string]dao.TwitterSource, data dao.TweetJson) {
	source := sources[data.Source.Url]
	news := data.ToDbNews()
	news.UserId = source.UserId
	news.SourceId = source.Id
	news.CreatedAt = time.Now().Unix()

	if err := db.Save(&news).Error; err != nil {
		log.Println("save news err:", err)
	}
}

func getLastPostId(sources map[string]dao.TwitterSource) map[int64]bool {
	ids := make(map[int64]bool)

	for _, source := range sources {
		news := []dao.TwitterNews{}

		err := db.Where(&dao.TwitterNews{SourceId: source.Id}).
			Order("Id desc").
			Select("TweetId").
			Find(&news).
			Error
		if err != nil {
			log.Printf("get TweeId for source %d error: %s", source.Id, err)
			return nil
		}

		for _, item := range news {
			ids[item.TweetId] = true
		}
	}

	return ids
}

func getSources(userId int64) map[string]dao.TwitterSource {
	sources := []dao.TwitterSource{}
	result := map[string]dao.TwitterSource{}

	err := db.Where(&dao.TwitterSource{UserId: userId}).Find(&sources).Error
	if err != nil {
		log.Printf("get sources for user %d error: %s", userId, err)
	}

	for _, source := range sources {
		result[source.Url] = source
	}

	return result
}

func getUsers() []dao.User {
	users := []dao.User{}
	ids := []int64{}

	err := db.Model(&dao.Settings{}).
		Where(&dao.Settings{TwitterEnabled: true}).
		Pluck("UserId", &ids).
		Error
	if err != nil {
		log.Println("get users ids for get twitter news error:", err)
		return nil
	}

	if len(ids) == 0 {
		return nil
	}

	err = db.Where("Id in (?)", ids).Find(&users).Error
	if err != nil {
		log.Println("get tweet users error:", err)
	}

	return users
}
