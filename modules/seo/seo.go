package seo

type Seo struct {
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Keyword     string `db:"keyword" json:"keyword"`
	Robot       string `db:"robot" json:"robot"`
	GoogleBot   string `db:"google_bot" json:"google_bot"`
}
