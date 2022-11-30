package users

type User struct {
}

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CreateChatRequest struct {
	Token         string   `json:"token"`
	ChatMemberIDs []uint32 `json:"chat_member_ids"`
	ChatName      string   `json:"name"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	Token string `json:"token"`
}

type UserInfo struct {
	UserID    uint32   `json:"user_id"`
	Username  string   `json:"username"`
	ChatIDs   []uint32 `json:"chat_ids"`
	ChatNames []string `json:"chat_names"`
}
