package response

import (
	"time"

	"github.com/linzhengen/ddd-gin-admin/app/domain/user"

	"github.com/linzhengen/ddd-gin-admin/app/domain/pagination"

	"github.com/linzhengen/ddd-gin-admin/pkg/util/json"
	"github.com/linzhengen/ddd-gin-admin/pkg/util/structure"
)

type User struct {
	ID        string    `json:"id"`                                    // ID
	UserName  string    `json:"user_name" binding:"required"`          // User Name
	RealName  string    `json:"real_name" binding:"required"`          // Real Name
	Password  string    `json:"password"`                              // Password
	Phone     string    `json:"phone"`                                 // Phone
	Email     string    `json:"email"`                                 // Email
	Status    int       `json:"status" binding:"required,max=2,min=1"` // Status(1:enable 2:disable)
	Creator   string    `json:"creator"`                               // Creator
	CreatedAt time.Time `json:"created_at"`                            // CreatedAt
	UserRoles UserRoles `json:"user_roles" binding:"required,gt=0"`    // UserRoles
}

func UserFromDomain(user *user.User) *User {
	item := new(User)
	structure.Copy(user, item)
	return item
}

func (a *User) String() string {
	return json.MarshalToString(a)
}

func (a *User) CleanSecure() *User {
	a.Password = ""
	return a
}

type UserQueryResult struct {
	Data       Users
	PageResult *pagination.Pagination
}

func (a UserQueryResult) ToShowResult(mUserRoles map[string]UserRoles, mRoles map[string]*Role) *UserShowQueryResult {
	return &UserShowQueryResult{
		PageResult: a.PageResult,
		Data:       a.Data.ToUserShows(mUserRoles, mRoles),
	}
}

type Users []*User

func UsersFromDomain(users user.Users) Users {
	list := make([]*User, len(users))
	for i, item := range users {
		ts := new(User)
		structure.Copy(item, ts)
		list[i] = ts
	}
	return list
}

func (a Users) ToIDs() []string {
	idList := make([]string, len(a))
	for i, item := range a {
		idList[i] = item.ID
	}
	return idList
}

func (a Users) ToUserShows(mUserRoles map[string]UserRoles, mRoles map[string]*Role) UserShows {
	list := make(UserShows, len(a))
	for i, item := range a {
		showItem := new(UserShow)
		structure.Copy(item, showItem)
		for _, roleID := range mUserRoles[item.ID].ToRoleIDs() {
			if v, ok := mRoles[roleID]; ok {
				showItem.Roles = append(showItem.Roles, v)
			}
		}
		list[i] = showItem
	}

	return list
}

// ----------------------------------------UserRole--------------------------------------

type UserRole struct {
	ID     string `json:"id"`      // ID
	UserID string `json:"user_id"` // User ID
	RoleID string `json:"role_id"` // Role ID
}

type UserRoleQueryResult struct {
	Data       UserRoles
	PageResult *pagination.Pagination
}

type UserRoles []*UserRole

func (a UserRoles) ToMap() map[string]*UserRole {
	m := make(map[string]*UserRole)
	for _, item := range a {
		m[item.RoleID] = item
	}
	return m
}

func (a UserRoles) ToRoleIDs() []string {
	list := make([]string, len(a))
	for i, item := range a {
		list[i] = item.RoleID
	}
	return list
}

func (a UserRoles) ToUserIDMap() map[string]UserRoles {
	m := make(map[string]UserRoles)
	for _, item := range a {
		m[item.UserID] = append(m[item.UserID], item)
	}
	return m
}

// ----------------------------------------UserShow--------------------------------------

type UserShow struct {
	ID        string    `json:"id"`         // ID
	UserName  string    `json:"user_name"`  // User Name
	RealName  string    `json:"real_name"`  // Real Name
	Phone     string    `json:"phone"`      // Phone
	Email     string    `json:"email"`      // Email
	Status    int       `json:"status"`     // Status(1:enable 2:disable)
	CreatedAt time.Time `json:"created_at"` // CreatedAt
	Roles     []*Role   `json:"roles"`      // Roles
}

type UserShows []*UserShow

type UserShowQueryResult struct {
	Data       UserShows
	PageResult *pagination.Pagination
}
