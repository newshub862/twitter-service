package service

import "newshub-twitter-service/dao"

var config = dao.Config{}

func Setup(cfg dao.Config) {
	config = cfg

	go listener()
}
