package l2db

type ResponseData struct {
	User  string
	Spins []SpinData
}

type SpinData struct {
	WinnerCategory string `json:"category"`
	WinnerSector   string `json:"sector"`
}
