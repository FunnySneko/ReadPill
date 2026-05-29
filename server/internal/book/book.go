package book

type Book struct {
	Id            int
	Title         string
	AuthorId      int
	AuthorName    string
	Description   string
	YearOfRelease int
	Tags          []BookTag
	CoverImageURL string
}

type BookFilter struct {
	ID             []int
	Title          *string
	AuthorID       []int
	AuthorFullName *string
	Description    *string
	YearOfRelease  map[int]int
	Tags           []int
}

type BookTag struct {
	ID   int
	Name string
}
