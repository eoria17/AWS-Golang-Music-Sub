package models

type Login struct {
	Email     string
	User_name string
	Password  string
}

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Year   string `json:"year"`
	WebURL string `json:"web_url"`
	ImgURL string `json:"img_url"`
}
