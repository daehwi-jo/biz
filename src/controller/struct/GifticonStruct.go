package _struct

type GifticonState struct {
	ResultCode string
	ResultData GifticonStateData
	ResultMsg  string
}

type GifticonStateData struct {
	CouponInfo struct {
		CpNo       string
		CpStatus   string
		ExpireDate string
		ItemDesc   string
		ItemImg    string
		ItemName   string
		ItemPrice  string
		OrdNo      string
		RestId     string
	}
}
