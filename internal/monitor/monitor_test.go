package monitor

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewTwitterClient(t *testing.T) {
	client := NewTwitterClient(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"))
	tweets, err := client.SearchTweets("@DuvalMagic")
	assert.Nil(t, err, err)
	assert.NotEmpty(t, tweets.Statuses, "Tweets shouldn't be empty")
}

func TestTwitterStream(t *testing.T) {
	client := NewTwitterClient(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"))
	client.StartMonitor()
}
