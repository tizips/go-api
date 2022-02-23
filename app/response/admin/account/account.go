package account

type ToAccountByInformationResponse struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Mobile   string `json:"mobile"`
}

type ToAccountByModuleResponse struct {
	Id   uint   `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}
