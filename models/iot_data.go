package models

// IoTData represents the structure of IoT data
type IoTData struct {
	Device         string  `json:"device"`
	Timestamp      string  `json:"timestamp"`
	ProVer         int     `json:"pro_ver"`
	MinorVer       int     `json:"minor_ver"`
	SN             int64   `json:"sn"`
	Model          string  `json:"model"`
	TYield         float64 `json:"tyield"`
	DYield         float64 `json:"dyield"`
	PF             float64 `json:"pf"`
	PMax           int     `json:"pmax"`
	PAC            int     `json:"pac"`
	SAC            int     `json:"sac"`
	UAB            int     `json:"uab"`
	UBC            int     `json:"ubc"`
	UCA            int     `json:"uca"`
	IA             int     `json:"ia"`
	IB             int     `json:"ib"`
	IC             int     `json:"ic"`
	Freq           int     `json:"freq"`
	TMod           float64 `json:"tmod"`
	TAmb           float64 `json:"tamb"`
	Mode           string  `json:"mode"`
	QAC            int     `json:"qac"`
	BusCapacitance float64 `json:"bus_capacitance"`
	ACCapacitance  float64 `json:"ac_capacitance"`
	PDC            float64 `json:"pdc"`
	PMaxLim        float64 `json:"pmax_lim"`
	SMaxLim        float64 `json:"smax_lim"`
	IsSent         bool    `json:"is_sent"`
	RegTimestamp   string  `json:"reg_timestamp"`
}

// SingleIoTDataResponse represents the response structure for single IoT data
type SingleIoTDataResponse struct {
	Data IoTData `json:"data"`
}
