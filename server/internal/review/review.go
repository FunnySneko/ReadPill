// █▀▀▀▀▀█ ▀ █ ▄ ▀ ▄ █▀▀▀▀▀█    PROJECT: ReadPill
// █ ███ █  ▀█▄▄█  ▀ █ ███ █    AUTHOR:  FunnySneko
// █ ▀▀▀ █  ▄ ▄  █▀▀ █ ▀▀▀ █
// ▀▀▀▀▀▀▀ █▄▀ █▄▀ ▀ ▀▀▀▀▀▀▀    © 2026
// ▀█▄█▀ ▀▄▄▀▄█▀▄█▄  ▀▄▄▄▄▄▀
// ▀▀▄█▀█▀▀ ▀ ██▀██  ▄█▄█▄ █
// █▄█▀ █▀▀▄▀█▀  ▄▄█▀▀▀▄▀▄▀
// █▀▀▄▄▀▀▄▀▀▀ ▀ ▀█▀█  ████▀
// ▀ ▀  ▀▀▀█▄ ▀▀█▀ █▀▀▀█▄█▄▀
// █▀▀▀▀▀█  ██▀▄▄ ▄█ ▀ █▄▀ ▀
// █ ███ █ █▀█▄▀▀▀████▀▀▀▄█▄
// █ ▀▀▀ █ ▄▀▀█ ▀▀▄ ▀█▀█▀███
// ▀▀▀▀▀▀▀ ▀▀▀▀▀   ▀▀   ▀  ▀

package review

type Review struct {
	Id     int
	BookId int
	UserId int
	// contribute rating is an average value of all ratings that have contribute rule set to true transformed to scale 0 to 1 and is not shown to user
	ContributeRating float32
	// user opinion is how strongly the rating deviates from user's average rating value (the way i see it, i guess ???) and is not shown to user
	UserOpinion float32
	// user bias is how strongly the user rating deviates from given book's average rating (i don't think i use this term right ?) and is not shown to user
	UserBias float32
	Ratings  []Rating
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
