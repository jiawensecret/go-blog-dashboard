package request

import "time"

type PageInfo struct {
	Page     int `json:"page" form:"page" query:"page"`
	PageSize int `json:"page_size" form:"page_size" query:"page_size"`
}

type GetById struct {
	Id uint32 `json:"id" form:"id" validate:"required,int" query:"id"`
}

type GetByBigId struct {
	Id uint64 `json:"id" form:"id" validate:"required,int" query:"id"`
}

type IdsReq struct {
	Ids []uint32 `json:"ids" form:"ids"`
}

type TimeRange struct {
	StartTime time.Time
	EndTime   time.Time
}

func GetPageInfo(page int, pageSize int) (int, int) {
	if page == 0 {
		page = 1
	}
	if pageSize <= 10 {
		pageSize = 10
	}

	return page, pageSize
}
