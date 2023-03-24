package repository

type FactRepository interface {
	AddTweet(t tweet) (int, error)
	Tweets() ([]tweet, error)
}
