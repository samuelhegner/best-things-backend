package leaderboardManager

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/samuelhegner/best-things/types"
)

func incrementEntry(member string, category string, redis *redis.Client) {
	total, yearly, monthly, daily := getDynamicBoardNames(category)
	redis.ZIncrBy(total, 1, member)
	redis.ZIncrBy(yearly, 1, member)
	redis.ZIncrBy(monthly, 1, member)
	redis.ZIncrBy(daily, 1, member)
}

func GetLeaderboards(category string, redis *redis.Client) types.CategoryBoards {
	total, yearly, monthly, daily := getDynamicBoardNames(category)
	tr := getBoardResult(total, redis)
	yr := getBoardResult(yearly, redis)
	mr := getBoardResult(monthly, redis)
	dr := getBoardResult(daily, redis)

	return types.CategoryBoards{
		Total: tr,
		Year:  yr,
		Month: mr,
		Day:   dr,
	}
}

func getBoardResult(key string, redis *redis.Client) types.BoardResult {
	tr, err := redis.ZRevRangeWithScores(key, 0, 5).Result()
	r := types.BoardResult{}

	if err != nil || len(tr) < 1 {
		return r
	}

	for _, z := range tr {
		r.Results[z.Member.(types.Card)] = int(z.Score)
	}

	return r
}

// Returns the Total, Yearly, monthly and daily board names for the provided category
func getDynamicBoardNames(categoryName string) (string, string, string, string) {
	year, month, day := getDateNumbers()

	total := getBoardKey("Total", categoryName)
	yearly := getBoardKey("Yearly", categoryName) + "-" + strconv.Itoa(year)
	monthly := getBoardKey("Monthly", categoryName) + "-" + strconv.Itoa(year) + "-" + strconv.Itoa(month)
	daily := getBoardKey("Daily", categoryName) + "-" + strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(day)

	return total, yearly, monthly, daily
}

func getDateNumbers() (int, int, int) {
	t := time.Now()
	year := t.Year()
	month := t.Month()
	day := t.Day()

	return year, int(month), day
}

func getBoardKey(name string, category string) string {
	builder := strings.Builder{}

	builder.WriteString(category)
	builder.WriteString("-")
	builder.WriteString(name)

	return builder.String()
}