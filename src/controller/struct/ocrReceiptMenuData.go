package _struct

type OcrReceiptDataList struct {
	RestId      string
	ReceiptId   string
	State       string
	StateDetail string
	BizNum      string
	AprvDt      string
	AprvNo      string
	TotalAmt    string
	List        []OcrReceiptData
}

type OcrReceiptData struct {
	ReceiptMenuId string
	MenuNm        string
	MenuPrice     string
	MenuEa        string
	MenuAmt       string
}
