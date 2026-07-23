package models

type Users struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	CategoryGroup string `json:"categoryGroup"`
	Username      string `json:"username"`
	Password      string `json:"-"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Name          string `json:"name" binding:"required"`
	CategoryGroup string `json:"categoryGroup" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

func (u *Users) IsAdmin() bool {
	return u.CategoryGroup == "IT"
}
