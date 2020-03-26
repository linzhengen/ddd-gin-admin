package schema

import "time"

// User ...
type User struct {
	RecordID  string    `json:"record_id"`
	UserName  string    `json:"user_name" binding:"required"`
	RealName  string    `json:"real_name" binding:"required"`
	Password  string    `json:"password"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    int       `json:"status" binding:"required,max=2,min=1"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"created_at"`
}

// CleanSecure clean secure info.
func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

// UserQueryParam user search params.
type UserQueryParam struct {
	UserName     string
	LikeUserName string
	LikeRealName string
	Status       int
}

// UserQueryOptions user search options
type UserQueryOptions struct {
	PageParam *PaginationParam
}

// UserQueryResult search result.
type UserQueryResult struct {
	Data       Users
	PageResult *PaginationResult
}

// Users 用户对象列表
type Users []*User

// ToUserShows to user list.
func (a Users) ToUserShows() UserShows {
	list := make(UserShows, len(a))

	for i, item := range a {
		showItem := &UserShow{
			RecordID:  item.RecordID,
			RealName:  item.RealName,
			UserName:  item.UserName,
			Email:     item.Email,
			Phone:     item.Phone,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		}
		list[i] = showItem
	}

	return list
}

// UserShow show user.
type UserShow struct {
	RecordID  string    `json:"record_id"`
	UserName  string    `json:"user_name"`
	RealName  string    `json:"real_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// UserShows show user list.
type UserShows []*UserShow

// UserShowQueryResult show query result
type UserShowQueryResult struct {
	Data       UserShows
	PageResult *PaginationResult
}
