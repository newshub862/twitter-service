package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"newshub-twitter-service/dao"
)

var (
	httpClient = http.Client{
		Timeout: 1 * time.Minute,
		Transport: &http.Transport{
			Dial:                (&net.Dialer{Timeout: 5 * time.Second}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
)

func getNews(name string) []dao.TweetJson {
	token, err := getToken()
	if err != nil {
		return nil
	}

	news := []dao.TweetJson{}
	friends := getFriends(name, token)

	for _, friend := range friends {
		requestUrl := "https://api.twitter.com/1.1/statuses/user_timeline.json?include_rts=1&exclude_replies=1&screen_name=" + friend

		req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
		if err != nil {
			log.Printf("create request for %s error: %s", friend, err)
			continue
		}

		req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

		response, err := httpClient.Do(req)
		if err != nil {
			log.Println("get news error:", err)
			continue
		}

		items := []dao.TweetJson{}
		if err := json.NewDecoder(response.Body).Decode(&items); err != nil {
			log.Println("decode news error:", err)
			continue
		}

		news = append(news, items...)
	}

	return news
}

func getFriends(name string, token dao.TwitterToken) []string {
	names := []string{}
	requestUrl := fmt.Sprintf("https://api.twitter.com/1.1/friends/list.json?cursor=-1&screen_name=%s&skip_status=true&include_user_entities=false", name)

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		log.Printf("create request for get friends for %s error: %s", name, err)
		return nil
	}

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	res, err := httpClient.Do(req)
	if err != nil {
		return names
	}

	friends := dao.TwitterFriends{}
	if err := json.NewDecoder(res.Body).Decode(&friends); err != nil {
		log.Println("decode friends err:", err)
		return names
	}

	for _, friend := range friends.Users {
		names = append(names, friend.ScreenName)
	}

	return names
}

func getToken() (dao.TwitterToken, error) {
	authString := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.ApiKey, config.ClientSecret)))

	val := url.Values{
		"grant_type": []string{"client_credentials"},
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.twitter.com/oauth2/token", strings.NewReader(val.Encode()))
	if err != nil {
		log.Println("create request for get token error:", err)
		return dao.TwitterToken{}, err
	}

	req.Header.Set("Authorization", "Basic "+authString)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	response, err := httpClient.Do(req)
	if err != nil {
		log.Println("get token err:", err)
		return dao.TwitterToken{}, err
	}

	data := dao.TwitterToken{}
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		log.Println("decode token error:", err)
		return dao.TwitterToken{}, err
	}

	return data, nil
}
