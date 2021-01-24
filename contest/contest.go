package contest

import (
	"math"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

// Do executes a contest based on a candidates list and a number of winner to pick. The selection is done by shuffling the candidates.
func Do(candidates []string, pick int) []string {
	if len(candidates) == 0 {
		logrus.Info("No candidates found meeting criteria (Follow + Retweet)")
		return nil
	}

	logrus.WithField("candidates", candidates).Infoln("Found some candidates to win the contest")

	if pick > len(candidates) {
		return candidates
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })

	return candidates[:int(math.Min(float64(len(candidates)), float64(pick)))]
}
