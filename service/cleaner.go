package service

import (
	"log"
	"time"

	"newshub-twitter-service/dao"
)

func Clean() {
	timestamp := time.Now().Add(-((30 * 24) * time.Hour)) // mounth ago
	query := `delete from twitternews where CreatedAt <= ?`

	sources := []dao.TwitterSource{}

	if err := db.Exec(query, timestamp).Error; err != nil {
		log.Println("clean old twitternews error:", err)
	}

	if err := db.Select("Id").Find(&sources).Error; err != nil {
		log.Println("select Id for twitter sources error:", err)
		return
	}

	for _, source := range sources {
		cnt := int64(0)
		if err := db.Where(&dao.TwitterNews{SourceId: source.Id}).Count(&cnt).Error; err != nil {
			log.Printf("get count sources by id %d error: %s", source.Id, err)
			continue
		}

		if cnt == 0 {
			db.Where(&dao.TwitterSource{Id: source.Id}).Delete(&dao.TwitterSource{})
		}
	}
}
