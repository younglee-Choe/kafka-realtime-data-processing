package structures

type BlockData struct {
	SourceData []byte `json:"sourcedata"`
	Length     int    `json:"length"`
}

type StateData struct {
	Counter     int
	EncodedData [][]byte
	Length      int
	IsLength    bool
}
