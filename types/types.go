package types

type BoardResult struct {
	Results map[Card]int
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
