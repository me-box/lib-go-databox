package libDatabox

//
//
// OBSERVE RESPONSE
//
//

type JsonObserveResponse struct {
	TimestampMS  int64
	DataSourceID string
	Key          string
	Json         []byte
}

type TextObserveResponse struct {
	TimestampMS  int64
	DataSourceID string
	Key          string
	Text         string
}

type BinaryObserveResponse struct {
	TimestampMS  int64
	DataSourceID string
	Key          string
	Data         []byte
}

//
//
// HYPERCAT TYPES
//
//

type DataSourceMetadata struct {
	Description    string
	ContentType    string
	Vendor         string
	DataSourceType string
	DataSourceID   string
	StoreType      string
	IsActuator     bool
	Unit           string
	Location       string
}

type relValPair struct {
	Rel string `json:"rel"`
	Val string `json:"val"`
}

type relValPairBool struct {
	Rel string `json:"rel"`
	Val bool   `json:"val"`
}

type hypercat struct {
	ItemMetadata []interface{} `json:"item-metadata"`
	Href         string        `json:"href"`
}
