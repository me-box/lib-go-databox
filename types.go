package libDatabox

//ContainerManagerOptions is used to configure the Container Manager
type ContainerManagerOptions struct {
	Version               string
	SwarmAdvertiseAddress string
	DefaultRegistryHost   string
	DefaultRegistry       string
	DefaultAppStore       string
	DefaultStoreImage     string
	ContainerManagerImage string
	ArbiterImage          string
	CoreNetworkImage      string
	CoreNetworkRelayImage string
	AppServerImage        string
	ExportServiceImage    string
	EnableDebugLogging    bool
	ClearSLAs             bool
	OverridePasword       string
	Hostname              string
	InternalIPs           []string
	ExternalIP            string
	HostPath              string
}

type DataboxType string

const (
	DataboxTypeApp    DataboxType = "app"
	DataboxTypeDriver DataboxType = "driver"
	DataboxTypeStore  DataboxType = "store"
)

type Macaroon string

type Repository struct {
	Type string `json:"Type"`
	Url  string `json:"url"`
}

type Package struct {
	//todo flesh this out
	// if they are going we still need datasources
}

type ExternalWhitelist struct {
	Urls        []string `json:"urls"`
	Description string   `json:"description"`
}

type ExportWhitelist struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type DataSource struct {
	Type          string       `json:"type"`
	Required      bool         `json:"required"`
	Name          string       `json:"name"`
	Clientid      string       `json:"clientid"`
	Granularities []string     `json:"granularities"`
	Hypercat      HypercatItem `json:"hypercat"`
}

type Manifest struct {
	ManifestVersion      int                  `json:"manifest-version"` //
	Name                 string               `json:"name"`
	DataboxType          DataboxType          `json:"databox-type"`
	Version              string               `json:"version"`     //this is databox version e.g 0.3.1
	Description          string               `json:"description"` // free text description
	Author               string               `json:"author"`      //Tosh Brown <Anthony.Brown@nottingham.ac.uk>
	License              string               `json:"license"`     //Software licence
	Tags                 []string             `json:"tags"`        //search tags
	Homepage             string               `json:"homepage"`    //homepage url
	Repository           Repository           `json:"repository"`
	Packages             []Package            `json:"packages"`
	ExportWhitelists     []ExportWhitelist    `json:"export-whitelist"`
	ExternalWhitelist    []ExternalWhitelist  `json:"external-whitelist"`
	ResourceRequirements ResourceRequirements `json:"resource-requirements"`
	DisplayName          string               `json:"displayName"`
	StoreURL             string               `json:"storeUrl"`
}

type SLA struct {
	ManifestVersion      int                  `json:"manifest-version"` //
	Name                 string               `json:"name"`
	DataboxType          DataboxType          `json:"databox-type"`
	Version              string               `json:"version"`     //this is databox version e.g 0.3.1
	Description          string               `json:"description"` // free text description
	Author               string               `json:"author"`      //Tosh Brown <Anthony.Brown@nottingham.ac.uk>
	License              string               `json:"license"`     //Software licence
	Tags                 []string             `json:"tags"`        //search tags
	Homepage             string               `json:"homepage"`    //homepage url
	Repository           Repository           `json:"repository"`
	Packages             []Package            `json:"packages"`
	AllowedCombinations  []string             `json:"allowed-combinations"`
	Datasources          []DataSource         `json:"datasources"`
	ExportWhitelists     []ExportWhitelist    `json:"export-whitelist"`
	ExternalWhitelist    []ExternalWhitelist  `json:"external-whitelist"`
	ResourceRequirements ResourceRequirements `json:"resource-requirements"`
	DisplayName          string               `json:"displayName"`
	StoreURL             string               `json:"storeUrl"`
}

type ResourceRequirements struct {
	Store string `json:"store"`
}

type DataSourceMetadata struct {
	Description    string
	ContentType    string
	Vendor         string
	DataSourceType string
	DataSourceID   string
	StoreType      StoreType
	IsActuator     bool
	Unit           string
	Location       string
}

type StoreType string

const StoreTypeTS StoreType = "ts"
const StoreTypeTSBlob StoreType = "ts/blob"
const StoreTypeKV StoreType = "kv"

type StoreContentType string

const ContentTypeJSON StoreContentType = "JSON"
const ContentTypeTEXT StoreContentType = "TEXT"
const ContentTypeBINARY StoreContentType = "BINARY"

type RelValPair struct {
	Rel string `json:"rel"`
	Val string `json:"val"`
}

type RelValPairBool struct {
	Rel string `json:"rel"`
	Val bool   `json:"val"`
}

type HypercatRoot struct {
	CatalogueMetadata []RelValPair   `json:"catalogue-metadata"`
	Items             []HypercatItem `json:"items"`
}

type HypercatItem struct {
	ItemMetadata []interface{} `json:"item-metadata"`
	Href         string        `json:"href"`
}
