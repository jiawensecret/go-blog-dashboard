package request

type MenuStore struct {
	Name        string `json:"name" query:"name" form:"name" validate:"required"`
	MenuSort    uint8  `json:"menu_sort" query:"menu_sort" form:"menu_sort"`
	MenuStatus  uint8  `json:"menu_status" query:"menu_status" form:"menu_status"`
	Route       string `json:"route" query:"route" form:"route"`
	Code        string `json:"code" query:"code" form:"code" validate:"required"`
	Description string `json:"description" query:"description" form:"description"`
	ParentId    uint32 `json:"parent_id" query:"parent_id" form:"parent_id"`
}

type MenuUpdate struct {
	Name        string `json:"name" query:"name" form:"name" validate:"required"`
	MenuSort    uint8  `json:"menu_sort" query:"menu_sort" form:"menu_sort"`
	MenuStatus  uint8  `json:"menu_status" query:"menu_status" form:"menu_status"`
	Route       string `json:"route" query:"route" form:"route"`
	Code        string `json:"code" query:"code" form:"code" validate:"required"`
	Description string `json:"description" query:"description" form:"description"`
	ParentId    uint32 `json:"parent_id" query:"parent_id" form:"parent_id"`
}
