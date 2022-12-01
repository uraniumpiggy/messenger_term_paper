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
	ChatMemberNames []string `json:"chat_member_names"`
	ChatName        string   `json:"name"`
}

type UserInfo struct {
	UserID    uint32   `json:"user_id"`
	Username  string   `json:"username"`
	ChatIDs   []uint32 `json:"chat_ids"`
	ChatNames []string `json:"chat_names"`
}

type UserWithToken struct {
	UserID    uint32   `json:"user_id"`
	Username  string   `json:"username"`
	ChatIDs   []uint32 `json:"chat_ids"`
	ChatNames []string `json:"chat_names"`
	Token     string   `json:"token"`
}
