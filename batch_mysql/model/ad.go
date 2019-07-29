package model

type DebrisType struct {
	Id    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Sort  int    `json:"sort" db:"sort"`
}

type Ad struct {
	Id      int    `json:"id" db:"id"`
	Title   string `json:"title" db:"title"`
	AsId    int    `json:"as_id" db:"as_id"`
	Pic     string `json:"pic" db:"pic"`
	Url     string `json:"url" db:"url"`
	AddTime int64  `json:"add_time" db:"addtime"`
	Sort    int    `json:"sort" db:"sort"`
	Open    int    `json:"open" db:"open"`
	Content string `json:"content" db:"content"`
}
