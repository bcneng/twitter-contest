package twitter

import (
	"context"
	"errors"

	"github.com/bcneng/twitter-contest/contest"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
)

// Credentials contains all needed Twitter credentials.
type Credentials struct {
	APIKey       string
	APIKeySecret string
}

// Contest runs a contest on Twitter.
func Contest(creds Credentials, tweetID, pick int, account string) (*contest.Result, error) {
	config := &clientcredentials.Config{
		ClientID:     creds.APIKey,
		ClientSecret: creds.APIKeySecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := config.Client(context.Background())
	client := twitter.NewClient(httpClient)

	var accountID int64
	if account != "" {
		logrus.WithField("account", account).Debugln("Looking for Twitter account")
		user, _, err := client.Users.Show(&twitter.UserShowParams{
			ScreenName:      account,
			IncludeEntities: twitter.Bool(false),
		})

		if err != nil {
			return nil, err
		}

		accountID = user.ID
	} else {
		logrus.Debugln("Twitter account defaults to Tweet author")
		t, _, err := client.Statuses.Show(int64(tweetID), &twitter.StatusShowParams{
			TrimUser:        twitter.Bool(true),
			IncludeEntities: twitter.Bool(false),
		})
		if err != nil {
			return nil, err
		}

		accountID = t.User.ID
	}

	if accountID == 0 {
		return nil, errors.New("account not found")
	}

	logrus.Debugln("Fetching followers for the specified account")
	followers, err := followers(client, accountID)
	if err != nil {
		return nil, err
	}

	logrus.WithField("tweet_id", tweetID).Debugln("Fetching retweets")
	rs, _, err := client.Statuses.Retweets(int64(tweetID), nil)
	if err != nil {
		return nil, err
	}

	if len(rs) == 0 {
		logrus.Debugln("No retweets found")
		return nil, nil
	}

	logrus.WithField("retweets", len(rs)).Debugln("Retweets found")
	var candidates []string
	for _, r := range rs {
		if _, ok := followers[r.User.ID]; ok {
			candidates = append(candidates, r.User.ScreenName)
		}
	}

	return contest.Do(candidates, pick), nil
}

func followers(client *twitter.Client, userID int64) (map[int64]struct{}, error) {
	var result *twitter.FollowerIDs
	var err error
	followers := make(map[int64]struct{})
	nextCursor := int64(-1)

	for {
		result, _, err = client.Followers.IDs(&twitter.FollowerIDParams{
			UserID: userID,
			Cursor: nextCursor,
		})
		if err != nil {
			return nil, err
		}

		for _, id := range result.IDs {
			followers[id] = struct{}{}
		}

		if result.NextCursor == 0 {
			break
		}

		nextCursor = result.NextCursor
	}

	return followers, nil
}
