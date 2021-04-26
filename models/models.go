package models

type Login struct {
	Email     string `json:"email"`
	User_name string `json:"user_name"`
	Password  string `json:"password"`
}

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Year   string `json:"year"`
	WebURL string `json:"web_url"`
	ImgURL string `json:"img_url"`
}
