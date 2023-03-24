package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server struct {
	FactRepository repository.FactRepository
}

type TweetMemoryRepository struct {
	TweetsList []tweet
}

func (r *TweetMemoryRepository) AddTweet(t tweet) (int, error) {
	r.TweetsList = append(r.TweetsList, t)
	return len(r.TweetsList), nil
}

func (r *TweetMemoryRepository) Tweets() ([]tweet, error) {
	return r.TweetsList, nil
}

type tweet struct {
	Message  string `json:"message"`
	Location string `json:"location"`
}

type tweetsList struct {
	Tweets []tweet `json:"tweets"`
}

type tweetData struct {
	ID int
}

func (s *Server) ListTweets(w http.ResponseWriter, r *http.Request) {
	tweets, err := s.TweetRepository.Tweets()
	if err != nil {
		log.Println("Failed to get tweets:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tl := tweetsList{
		Tweets: tweets,
	}

	response, err := json.Marshal(tl)
	if err != nil {
		log.Println("Failed to marshal tweet data for response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	fmt.Printf("%+v\n", response)
	return
}

func (s *Server) AddTweet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read body:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	t := tweet{}

	if err := json.Unmarshal(body, &t); err != nil {
		log.Println("Failed to unmarshal tweet payload:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := s.TweetRepository.AddTweet(t)
	if err != nil {
		log.Println("Failed to store tweet:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tweetData := tweetData{
	 	ID: id,
	 }

	response, err := json.Marshal(tweetData)
	if err != nil {
		log.Println("Failed to marshal tweet data for response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(response)
	fmt.Printf("Tweet: `%s` from %s\n", t.Message, t.Location)
	return
}