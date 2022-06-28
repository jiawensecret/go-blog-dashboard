package request

type RouteList struct {
	Name string `json:"name" query:"name"`
	PageInfo
}

type RouteStore struct {
	Name   string `json:"name" query:"name" form:"name" validate:"required"`
	Method string `json:"method" query:"method" form:"method" validate:"required"`
	Path   string `json:"path" query:"path" form:"path"`
}

type RouteUpdate struct {
	Name   string `json:"name" query:"name" form:"name" validate:"required"`
	Method string `json:"method" query:"method" form:"method" validate:"required"`
	Path   string `json:"path" query:"path" form:"path"`
}
