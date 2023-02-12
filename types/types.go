package types

type BoardResult struct {
	Results []BoardEntry
}

type BoardEntry struct {
	Member string
	Score  int
}

type CategoryBoards struct {
	Total BoardResult
	Year  BoardResult
	Month BoardResult
	Day   BoardResult
}

type Card string

type Category struct {
	Name string
}

type Matchup struct {
	Id        string
	OptionOne Card
	OptionTwo Card
	Category
	Expiration int64
}

type SheetData struct {
	Name     string
	Image    string
	Category string
}

type MatchupSubmit struct {
	Guid     string `json:"guid" binding:"required"`
	Winner   string `json:"winner" binding:"required"`
	Category string `json:"category" binding:"required"`
}
