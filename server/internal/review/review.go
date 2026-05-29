package review

type Review struct {
	Id          int
	BookId      int
	UserId      int
	UserOpinion float32
	Ratings     []Rating
}

type Rating struct {
	Name  string
	Value float32
}

type ReviewRules struct {
	RatingRules []RatingRule `json:"rating_rules"`
}

type RatingRule struct {
	Name         string  `json:"name"`
	ValueCeiling float32 `json:"value_ceiling"`
	Contribute   bool    `json:"contribute"`
}
