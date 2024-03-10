package entities

type Image struct {
	Id       int `db:"id" json:"id"`
	FileName string `db:"filename" json:"filename"`
	Url      string `db:"url" json:"url"`
}