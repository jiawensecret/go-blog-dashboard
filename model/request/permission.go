package request

type PermissionList struct {
	Code string `json:"code" query:"code"`
	Name string `json:"name" query:"name"`
}

type PermissionStore struct {
	Name        string `json:"name" query:"name" form:"name" validate:"required"`
	MenuId      uint32 `json:"menu_id" query:"menu_id" form:"menu_id" validate:"required"`
	Code        string `json:"code" query:"code" form:"code" validate:"required"`
	Flag        string `json:"flag" query:"flag" form:"flag"`
	Description string `json:"description" query:"description" form:"description"`
}

type PermissionUpdate struct {
	Name        string `json:"name" query:"name" form:"name" validate:"required"`
	MenuId      uint32 `json:"menu_id" query:"menu_id" form:"menu_id" validate:"required"`
	Code        string `json:"code" query:"code" form:"code" validate:"required"`
	Flag        string `json:"flag" query:"flag" form:"flag"`
	Description string `json:"description" query:"description" form:"description"`
}

type PermissionAssRoute struct {
	RouteIds []uint32 `json:"route_ids" form:"route_ids" query:"route_ids"`
}
