package review

type Review struct {
	Id          int
	BookId      int
	UserId      int
	UserOpinion float32
	Ratings     []Rating
}

type Rating struct {
	Name         string
	Value        float32
	ValueCeiling float32
	Contribute   bool
}

type ReviewRules struct {
	ShowcaseRatingScale float32      `json:"showcase_rating_scale"`
	RatingRules         []RatingRule `json:"rating_rules"`
}

type RatingRule struct {
	Name         string  `json:"name"`
	ValueCeiling float32 `json:"value_ceiling"`
	Contribute   bool    `json:"contribute"`
}
