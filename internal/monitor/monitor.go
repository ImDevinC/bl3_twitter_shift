package monitor

import (
	"regexp"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type MonitorClient struct {
	client *twitter.Client
	stream *twitter.Stream
}

func NewTwitterClient(consumerKey string, consumerSecret string, accessToken string, accessSecret string) MonitorClient {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return MonitorClient{client, nil}
}

func (c *MonitorClient) SearchTweets(query string) (*twitter.Search, error) {
	tweets, _, err := c.client.Search.Tweets(&twitter.SearchTweetParams{
		Query: query,
	})
	if err != nil {
		return &twitter.Search{}, err
	}
	return tweets, nil
}

func (c *MonitorClient) StopMonitor() {
	if c.stream != nil {
		c.stream.Stop()
	}
}

func (c *MonitorClient) StartMonitor() error {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		processTweet(tweet.Text)
	}

	var err error
	filterParams := &twitter.StreamFilterParams{
		Follow:        []string{"8369072"},
		StallWarnings: twitter.Bool(true),
	}

	c.stream, err = c.client.Streams.Filter(filterParams)
	if err != nil {
		return err
	}

	go demux.HandleChan(c.stream.Messages)

	return nil
}

func processTweet(input string) (bool, []string) {
	codes := []string{}
	found := false

	re := regexp.MustCompile(`(\w{5}\-){4}\w{5}`)
	matches := re.FindAll([]byte(input), -1)
	if matches != nil {
		found = true
		for _, code := range matches {
			codes = append(codes, string(code))
		}
	}

	return found, codes
}
