package monitor

import (
	"github.com/dghubble/go-twitter/twitter"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type MonitorClient struct {
	client *twitter.Client
	stream *twitter.Stream
}

func NewTwitterClient(consumerKey string, consumerSecret string) MonitorClient {
	config := &clientcredentials.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(oauth2.NoContext)
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

func (c *MonitorClient) StopStream() {
	if c.stream != nil {
		c.stream.Stop()
	}
}

func (c *MonitorClient) StartMonitor() {
	log.Println("Starting monitor...")
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		log.Println(tweet)
	}

	// filterParams := &twitter.StreamFilterParams{
	// 	Follow:         []string{"@DuvalMagic"},
	// 	StallWarnings: twitter.Bool(true),
	// }

	// c.stream, err = c.client.Streams.Filter(filterParams)
	var err error
	c.stream, err = c.client.Streams.Sample(&twitter.StreamSampleParams{
		StallWarnings: twitter.Bool(true),
	})
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(c.stream.Messages)

	log.Println("Waiting for stream end...")
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	log.Println("Stopping stream...")
	c.StopStream()
}
