package matchupManager

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/samuelhegner/best-things/leaderboardManager"
	"github.com/samuelhegner/best-things/types"
)

type matchupManager struct {
	categories []types.Category
	client     *redis.Client
}

const checksumKey string = "checksum"
const matchupExpirationSeconds int64 = 60

func (m *matchupManager) init() {
	jsonFile, err := os.Open("./Files/data.json")

	if err != nil {
		return
	}

	byteValue, _ := io.ReadAll(jsonFile)
	var entries []types.SheetData
	json.Unmarshal(byteValue, &entries)

	data := make(map[string][]types.Card)

	for _, e := range entries {
		data[e.Category] = append(data[e.Category], types.Card(e.Name))
	}

	categories := make([]types.Category, len(data))

	i := 0

	for k := range data {
		categories[i] = types.Category{Name: k}
		i++
	}

	m.categories = categories
	m.initDB(data)
}

func (m *matchupManager) initDB(data map[string][]types.Card) {
	checksum := dataToChecksum(data)

	if !m.dbNeedsUpdate(checksum) {
		return
	}

	fmt.Println("Checksum not at parity. Updating db...")

	for _, c := range m.categories {
		m.createMemberSet(c.Name, data[c.Name])
	}

	m.setDbChecksum(checksum)
}

func (m *matchupManager) setDbChecksum(checksum string) {
	m.client.Set(checksumKey, checksum, 0)
}

func dataToChecksum(data map[string][]types.Card) string {
	json, _ := json.Marshal(data)

	hash := sha256.New()
	hash.Write(json)
	sum := hex.EncodeToString(hash.Sum(nil))
	return sum
}

func (m *matchupManager) dbNeedsUpdate(checksum string) bool {
	res, err := m.client.Get(checksumKey).Result()

	if err != nil {
		return true
	}

	return res != checksum
}

func (m *matchupManager) createMemberSet(category string, members []types.Card) {
	key := categorySetKey(category)

	names := make([]string, len(members))

	for i, m := range members {
		names[i] = string(m)
	}

	m.client.Del(key)
	m.client.SAdd(key, names)
}

func categorySetKey(category string) string {
	return "category_" + category
}

func NewMatchupManager() *matchupManager {
	redisUrl := os.Getenv("REDIS_URL")

	if redisUrl == "" {
		log.Fatal("Error loading .env file")
	}

	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}

	mm := matchupManager{
		categories: []types.Category{},
		client:     redis.NewClient(opt),
	}

	mm.init()

	return &mm
}

func (m *matchupManager) GetCategories() []types.Category {
	return m.categories
}

func (m *matchupManager) GetMatchup(category string) (types.Matchup, error) {

	if !m.hasCategory(category) {
		return types.Matchup{}, fmt.Errorf("category not available")
	}

	res, err := m.client.SRandMemberN(categorySetKey(category), 2).Result()

	if err != nil {
		return types.Matchup{}, fmt.Errorf("SRandMemberN failed")
	}

	expiration := time.Now().Unix() + matchupExpirationSeconds

	matchup := types.Matchup{
		Id:         uuid.New().String(),
		Category:   types.Category{Name: category},
		OptionOne:  types.Card(res[0]),
		OptionTwo:  types.Card(res[1]),
		Expiration: expiration,
	}

	m.client.SAdd(matchup.Id, []string{string(matchup.OptionOne), string(matchup.OptionTwo)})
	m.client.Expire(matchup.Id, time.Second*time.Duration(matchupExpirationSeconds))
	fmt.Println(res)

	return matchup, nil
}

func (m *matchupManager) GetCategoryBoards(category string) (types.CategoryBoards, error) {

	if !m.hasCategory(category) {
		return types.CategoryBoards{}, fmt.Errorf("category not available")
	}

	return leaderboardManager.GetLeaderboards(category, m.client), nil
}

func (m *matchupManager) hasCategory(name string) bool {

	for _, c := range m.categories {
		if c.Name == name {
			return true
		}
	}

	return false
}

func (m *matchupManager) SubmitMatchupResponse(guid string, winner string, category string) (bool, error) {
	if !m.hasCategory(category) {
		return false, fmt.Errorf("category not available")
	}

	isMember, err := m.client.SIsMember(guid, winner).Result()

	if err != nil {
		return false, err
	}

	if !isMember {
		return false, fmt.Errorf("matchup doesn't exist or name not in the matchup options")
	}

	_, err = m.client.Del(guid).Result()

	if err != nil {
		return false, fmt.Errorf("failed to del matchup entry: " + err.Error())
	}

	leaderboardManager.IncrementEntry(winner, category, m.client)
	return true, nil
}
