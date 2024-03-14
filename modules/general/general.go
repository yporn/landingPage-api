package general

import "github.com/yporn/sirarom-backend/modules/entities"

type General struct {
	Id int `db:"id" json:"id"`
	Tel string `db:"tel" json:"tel"`
	Email string `db:"email" json:"email"`
	LinkFacebook string `db:"link_facebook" json:"link_facebook"`
	LinkInstagram string `db:"link_instagram" json:"link_instagram"`
	LinkTwitter string `db:"link_twitter" json:"link_twitter"`
	LinkTikTok string `db:"link_tiktok" json:"link_tiktok"`
	LinkLine string `db:"link_line" json:"link_line"`
	LinkWebsite string `db:"link_website" json:"link_website"`
	Images      []*entities.Image `json:"images"`
	FileName string `db:"filename" json:"filename"`
	Url      string `db:"url" json:"url"`
}
