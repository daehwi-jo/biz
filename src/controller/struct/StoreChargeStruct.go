package _struct

type StoreChargeObj struct {
	RestId     string
	AllUseYn   string
	ChargeList []chargeList
}

type chargeList struct {
	SeqNo  string
	Amt    string
	AddAmt string
	UseYn  string
}
