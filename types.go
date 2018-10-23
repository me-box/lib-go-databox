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
	CoreUIImage           string
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
	Arch                  string //current architecture used to chose the correct docker images "" for x86 or "arm64v8" for arm64v8 ;-)
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

type ExternalWhitelist struct {
	Urls        []string `json:"urls"`
	Description string   `json:"description"`
}

type ExportWhitelist struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type DataSource struct {
	Type               string       `json:"type"`
	Required           bool         `json:"required"`
	Name               string       `json:"name"`
	Clientid           string       `json:"clientid"`
	Granularities      []string     `json:"granularities"`
	Hypercat           HypercatItem `json:"hypercat"`
	Min                int          `json:"min"`
	Max                int          `json:"max"`
	NotifyOfNewSources bool         `json:"allow-notification-of-new-sources"`
}

type Manifest struct {
	ManifestVersion      int                  `json:"manifest-version"`      //
	Name                 string               `json:"name"`                  //
	DockerImage          string               `json:"docker-image"`          // Optional: docker image name  e.g my-cool-app (-amd64 or -amd64v8 will be added depending on the platform) defaults to Name from above
	DockerRegistry       string               `json:"docker-registry"`       // Optional: docker registry e.g myDockerRegistry defaults to datboxsystems
	DockerImageTag       string               `json:"docker-image-tag"`      // Optional: docker image tag e.g latest or v0.5.1 etc defaults to the running version of databox
	DataboxType          DataboxType          `json:"databox-type"`          //
	Version              string               `json:"version"`               // this is databox version e.g 0.3.1
	Description          string               `json:"description"`           // free text description
	Author               string               `json:"author"`                // Tosh Brown <Anthony.Brown@nottingham.ac.uk>
	License              string               `json:"license"`               // Software licence
	Tags                 []string             `json:"tags"`                  // search tags
	Homepage             string               `json:"homepage"`              // homepage url
	Repository           Repository           `json:"repository"`            //  git repo where the core can be found
	DataSources          []DataSource         `json:"datasources"`           //
	ExportWhitelists     []ExportWhitelist    `json:"export-whitelist"`      // "export-whitelist": [{"url": "https://export.amar.io/","description": "Exports the data to amar.io"}],
	ExternalWhitelist    []ExternalWhitelist  `json:"external-whitelist"`    // "external-whitelist": [{"urls": ["https://api.twitter.com/"],"description": "reason displayed to user for requiring access"}]
	ResourceRequirements ResourceRequirements `json:"resource-requirements"` //this is where you can requests a store "store":"core-store" is the only valid option for now.
	DisplayName          string               `json:"displayName"`
	StoreURL             string               `json:"storeUrl"`
}

type SLA struct {
	ManifestVersion      int                  `json:"manifest-version"`      //
	Name                 string               `json:"name"`                  // container name  e.g core-store
	DockerImage          string               `json:"docker-image"`          // Optional: docker image name  e.g my-cool-app (-amd64 or -amd64v8 will be added depending on the platform) defaults to Name from above
	DockerRegistry       string               `json:"docker-registry"`       // Optional: docker registry e.g myDockerRegistry defaults to datboxsystems
	DockerImageTag       string               `json:"docker-image-tag"`      // Optional: docker image tag e.g latest or v0.5.1 etc defaults to the running version of databox
	DataboxType          DataboxType          `json:"databox-type"`          // App or driver
	Version              string               `json:"version"`               // this is databox version e.g 0.3.1
	Description          string               `json:"description"`           // free text description
	Author               string               `json:"author"`                // Tosh Brown <Anthony.Brown@nottingham.ac.uk>
	License              string               `json:"license"`               // Software licence
	Tags                 []string             `json:"tags"`                  // search tags
	Homepage             string               `json:"homepage"`              // homepage url
	Repository           Repository           `json:"repository"`            // git repo where the core can be found
	Datasources          []DataSource         `json:"datasources"`           //
	ExportWhitelists     []ExportWhitelist    `json:"export-whitelist"`      //
	ExternalWhitelist    []ExternalWhitelist  `json:"external-whitelist"`    //
	ResourceRequirements ResourceRequirements `json:"resource-requirements"` //this is where you can requests a store "store":"core-store" is the only valid option for now.
	DisplayName          string               `json:"displayName"`           //
	StoreURL             string               `json:"storeUrl"`              //
}

type ResourceRequirements struct {
	Store string `json:"store"`
}

type DataSourceMetadata struct {
	Description    string           //required
	ContentType    StoreContentType //required
	Vendor         string           //required
	DataSourceType string           //required
	DataSourceID   string           //required
	StoreType      StoreType        //required
	IsActuator     bool
	IsFunc         bool
	Unit           string
	Location       string
}

type StoreType string

const StoreTypeTS StoreType = "ts"
const StoreTypeTSBlob StoreType = "ts/blob"
const StoreTypeKV StoreType = "kv"
const StoreTypeFunc StoreType = "notification/request"

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

//
//
// OBSERVE RESPONSE
//
//
type ObserveResponse struct {
	TimestampMS  int64
	DataSourceID string
	Key          string
	Data         []byte
}

type NotifyResponse struct {
	TimestampMS int64
	ContentType StoreContentType
	Data        []byte
}
