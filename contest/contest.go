package contest

import (
	"math"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

// Version stores the git commit SHA from where the app got built. Ideally injected when building the app.
var Version = "not-specified"

// Result represents the contest result
type Result struct {
	Winners []string  `json:"winners"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

// Do executes a contest based on a candidates list and a number of winner to pick. The selection is done by shuffling the candidates.
func Do(candidates []string, pick int) *Result {
	if len(candidates) == 0 {
		logrus.Info("No candidates found meeting criteria (Follow + Retweet)")
		return nil
	}

	if pick == 0 {
		return nil
	}

	result := &Result{
		Version: Version,
		Time:    time.Now(),
	}

	logrus.WithField("candidates", candidates).Infoln("Found some candidates to win the contest")

	if pick > len(candidates) {
		result.Winners = candidates
		return result
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	result.Winners = candidates[:int(math.Min(float64(len(candidates)), float64(pick)))]
	return result
}
