package monitor

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const sampleTweet = `
Sunday, Bloody Sunday. SHiFT code for a Golden Key for Borderlands 3: 

CZCTJ-CZ59T-HC35W-T3BJB-ZTZJC

Active until 10pm CST.

Redeem in game or via http://shift.gearboxsoftware.com  - Collect Golden Key from in-game mail in “Social” menu. 

Happy Vault Hunting!
`
const sampleCode = "CZCTJ-CZ59T-HC35W-T3BJB-ZTZJC"

func TestNewTwitterClient(t *testing.T) {
	client := NewTwitterClient(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	tweets, err := client.SearchTweets("@DuvalMagic")
	assert.Nil(t, err, err)
	assert.NotEmpty(t, tweets.Statuses, "Tweets shouldn't be empty")
}

func TestTwitterStream(t *testing.T) {
	client := NewTwitterClient(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"), os.Getenv("ACCESS_TOKEN"), os.Getenv("ACCESS_SECRET"))
	err := client.StartMonitor()
	assert.Nil(t, err, err)
	time.Sleep(10 * time.Second)
	client.StopMonitor()
}

func TestTweetProcessing(t *testing.T) {
	found, codes := processTweet(sampleTweet)
	assert.True(t, found, "Found should be true")
	assert.Len(t, codes, 1, "Should only be one code")
	assert.Equal(t, "CZCTJ-CZ59T-HC35W-T3BJB-ZTZJC", codes[0], "Code %s does not match %s", codes[0], sampleCode)
}
