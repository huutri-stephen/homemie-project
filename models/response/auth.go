package response

type LoginResponse struct {
    AccessToken string      `json:"access_token"`
    User        UserPayload `json:"user"`
}

type UserPayload struct {
    ID       int64 `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    UserType string `json:"user_type"`
    Status   string `json:"status"`
}