package request

type UserList struct {
	Ids         []uint32 `json:"ids" form:"ids" query:"ids"`
	Username    string   `json:"username" form:"username" query:"username"`
	Name        string   `json:"name" form:"name" query:"name"`
	Tel         string   `json:"tel" form:"tel" query:"tel"`
	Email       string   `json:"email" form:"email" query:"email"`
	UserStatus  int8     `json:"user_status" form:"user_status" query:"user_status"`
	Permissions []uint32 `json:"permissions"`
	PageInfo
}

type UserStore struct {
	Username      string   `json:"username" form:"username" validate:"required,min=2,max=25" query:"username"`
	Name          string   `json:"name" form:"name" validate:"required,min=2,max=25" query:"name"`
	Tel           string   `json:"tel" form:"tel" validate:"required,min=8,max=11" query:"tel"`
	Email         string   `json:"email" form:"email" query:"email"`
	Position      string   `json:"position" form:"position" validate:"required" query:"position"`
	DepartmentIds []uint64 `json:"department_ids" form:"department_ids" query:"department_ids"`
}

type UserUpdate struct {
	Username      string   `json:"username" form:"username" validate:"min=2,max=25" query:"username"`
	Name          string   `json:"name" form:"name" validate:"required,min=2,max=25" query:"name"`
	Tel           string   `json:"tel" form:"tel" validate:"min=8,max=11" query:"tel"`
	Email         string   `json:"email" form:"email" query:"email"`
	Position      string   `json:"position" form:"position" validate:"required" query:"position"`
	DepartmentIds []uint64 `json:"department_ids" form:"department_ids" query:"department_ids"`
}

type UserChangeStatus struct {
	UserStatus int8 `json:"user_status" form:"user_status"`
}

type UserPermissionStore struct {
	UserIds         []uint32         `json:"user_ids" form:"user_ids"`
	PermissionIds   []uint32         `json:"permission_ids" form:"permission_ids"`
	DataPermissions []DataPermission `json:"data_permissions" form:"data_permissions"`
}

type UserPermissionUpdate struct {
	PermissionIds   []uint32         `json:"permission_ids" form:"permission_ids"`
	DataPermissions []DataPermission `json:"data_permissions" form:"data_permissions"`
}

type DataPermission struct {
	MenuId      uint32 `json:"menu_id" form:"menu_id"`
	ExpiredTime int64  `json:"expired_time" form:"expired_time"`
	DataLevel   uint32 `json:"data_level" form:"data_level"`
}

type CopyUser struct {
	User UserStore `json:"user" form:"user"`
	UserPermissionUpdate
}

type UserInfo struct {
	Mobile string `json:"mobile" form:"mobile" query:"mobile" validate:"required"`
}

type UserPermissionDelete struct {
	PermissionIds []uint32 `json:"permission_ids" form:"permission_ids"`
	UserIds       []uint32 `json:"user_ids" form:"user_ids"`
}
