

# libDatabox
`import "./"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [func ChkErr(err error)](#ChkErr)
* [func ChkErrFatal(err error)](#ChkErrFatal)
* [func Debug(msg string)](#Debug)
* [func Err(msg string)](#Err)
* [func GetStoreURLFromDsHref(href string) (string, error)](#GetStoreURLFromDsHref)
* [func Info(msg string)](#Info)
* [func NewDataboxHTTPsAPI() *http.Client](#NewDataboxHTTPsAPI)
* [func NewDataboxHTTPsAPIWithPaths(cmRootCaPath string) *http.Client](#NewDataboxHTTPsAPIWithPaths)
* [func Warn(msg string)](#Warn)
* [type AggregationType](#AggregationType)
* [type ArbiterClient](#ArbiterClient)
  * [func NewArbiterClient(arbiterTokenPath string, databoxRequest *http.Client, arbiterURL string) ArbiterClient](#NewArbiterClient)
  * [func (arb *ArbiterClient) GetRootDataSourceCatalogue() (HypercatRoot, error)](#ArbiterClient.GetRootDataSourceCatalogue)
  * [func (arb *ArbiterClient) GrantComponentPermission()](#ArbiterClient.GrantComponentPermission)
  * [func (arb *ArbiterClient) GrantContainerPermissions(permissions ContainerPermissions) error](#ArbiterClient.GrantContainerPermissions)
  * [func (arb *ArbiterClient) InvalidateCache(href string, method string)](#ArbiterClient.InvalidateCache)
  * [func (arb *ArbiterClient) RegesterDataboxComponent(componentName, componenttype string)](#ArbiterClient.RegesterDataboxComponent)
  * [func (arb *ArbiterClient) RemoveDataboxComponent()](#ArbiterClient.RemoveDataboxComponent)
  * [func (arb *ArbiterClient) RequestToken(href string, method string) ([]byte, error)](#ArbiterClient.RequestToken)
  * [func (arb *ArbiterClient) RevokeComponentPermission()](#ArbiterClient.RevokeComponentPermission)
  * [func (arb *ArbiterClient) UpdateArbiter(name string, tokenString string, databoxType DataboxType) error](#ArbiterClient.UpdateArbiter)
* [type ContainerManagerOptions](#ContainerManagerOptions)
* [type ContainerPermissions](#ContainerPermissions)
* [type CoreNetworkClient](#CoreNetworkClient)
  * [func NewCoreNetworkClient(containerManagerKeyPath string, request *http.Client) CoreNetworkClient](#NewCoreNetworkClient)
  * [func (cnc CoreNetworkClient) ConnectEndpoints(serviceName string, peers []string) error](#CoreNetworkClient.ConnectEndpoints)
  * [func (cnc CoreNetworkClient) DisconnectEndpoints(serviceName string, netConfig PostNetworkConfig) error](#CoreNetworkClient.DisconnectEndpoints)
  * [func (cnc CoreNetworkClient) NetworkOfService(service swarm.Service, serviceName string) (PostNetworkConfig, error)](#CoreNetworkClient.NetworkOfService)
  * [func (cnc CoreNetworkClient) PostUninstall(name string, netConfig PostNetworkConfig) error](#CoreNetworkClient.PostUninstall)
  * [func (cnc CoreNetworkClient) PreConfig(localContainerName string, sla SLA) NetworkConfig](#CoreNetworkClient.PreConfig)
  * [func (cnc CoreNetworkClient) RegisterPrivileged() error](#CoreNetworkClient.RegisterPrivileged)
  * [func (cnc CoreNetworkClient) ServiceRestart(serviceName string, oldIP string, newIP string) error](#CoreNetworkClient.ServiceRestart)
* [type CoreStoreClient](#CoreStoreClient)
  * [func NewCoreStoreClient(databoxRequest *http.Client, arbiterClient *ArbiterClient, serverKeyPath string, storeEndPoint string, enableLogging bool) *CoreStoreClient](#NewCoreStoreClient)
  * [func (csc *CoreStoreClient) GetStoreDataSourceCatalogue(href string) (HypercatRoot, error)](#CoreStoreClient.GetStoreDataSourceCatalogue)
  * [func (csc *CoreStoreClient) HypercatToDataSourceMetadata(hypercatDataSourceDescription string) (DataSourceMetadata, string, error)](#CoreStoreClient.HypercatToDataSourceMetadata)
  * [func (csc *CoreStoreClient) RegisterDatasource(metadata DataSourceMetadata) error](#CoreStoreClient.RegisterDatasource)
* [type DataSource](#DataSource)
* [type DataSourceMetadata](#DataSourceMetadata)
* [type DataboxType](#DataboxType)
* [type ExportWhitelist](#ExportWhitelist)
* [type ExternalWhitelist](#ExternalWhitelist)
* [type Filter](#Filter)
* [type FilterType](#FilterType)
* [type HypercatItem](#HypercatItem)
* [type HypercatRoot](#HypercatRoot)
* [type KVStore](#KVStore)
  * [func (kvj *KVStore) Delete(dataSourceID string, key string) error](#KVStore.Delete)
  * [func (kvj *KVStore) DeleteAll(dataSourceID string) error](#KVStore.DeleteAll)
  * [func (kvj *KVStore) ListKeys(dataSourceID string) ([]string, error)](#KVStore.ListKeys)
  * [func (kvj *KVStore) Observe(dataSourceID string) (&lt;-chan []byte, error)](#KVStore.Observe)
  * [func (kvj *KVStore) ObserveKey(dataSourceID string, key string) (&lt;-chan []byte, error)](#KVStore.ObserveKey)
  * [func (kvj *KVStore) Read(dataSourceID string, key string) ([]byte, error)](#KVStore.Read)
  * [func (kvj *KVStore) Write(dataSourceID string, key string, payload []byte) error](#KVStore.Write)
* [type LogEntries](#LogEntries)
* [type Logger](#Logger)
  * [func New(store *CoreStoreClient, outputDebugLogs bool) (*Logger, error)](#New)
  * [func (l Logger) ChkErr(err error)](#Logger.ChkErr)
  * [func (l Logger) Debug(msg string)](#Logger.Debug)
  * [func (l Logger) Err(msg string)](#Logger.Err)
  * [func (l Logger) GetLastNLogEntries(n int) Logs](#Logger.GetLastNLogEntries)
  * [func (l Logger) GetLastNLogEntriesRaw(n int) []byte](#Logger.GetLastNLogEntriesRaw)
  * [func (l Logger) Info(msg string)](#Logger.Info)
  * [func (l Logger) Warn(msg string)](#Logger.Warn)
* [type Logs](#Logs)
* [type Macaroon](#Macaroon)
* [type Manifest](#Manifest)
* [type NetworkConfig](#NetworkConfig)
* [type Package](#Package)
* [type PostNetworkConfig](#PostNetworkConfig)
* [type RelValPair](#RelValPair)
* [type RelValPairBool](#RelValPairBool)
* [type Repository](#Repository)
* [type ResourceRequirements](#ResourceRequirements)
* [type Route](#Route)
* [type SLA](#SLA)
* [type StoreContentType](#StoreContentType)
* [type StoreType](#StoreType)
* [type TSBlobStore](#TSBlobStore)
  * [func (tbs *TSBlobStore) Earliest(dataSourceID string) ([]byte, error)](#TSBlobStore.Earliest)
  * [func (tbs *TSBlobStore) FirstN(dataSourceID string, n int) ([]byte, error)](#TSBlobStore.FirstN)
  * [func (tbs *TSBlobStore) LastN(dataSourceID string, n int) ([]byte, error)](#TSBlobStore.LastN)
  * [func (tbs *TSBlobStore) Latest(dataSourceID string) ([]byte, error)](#TSBlobStore.Latest)
  * [func (tbs *TSBlobStore) Length(dataSourceID string) (int, error)](#TSBlobStore.Length)
  * [func (tbs *TSBlobStore) Observe(dataSourceID string) (&lt;-chan []byte, error)](#TSBlobStore.Observe)
  * [func (tbs *TSBlobStore) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64) ([]byte, error)](#TSBlobStore.Range)
  * [func (tbs *TSBlobStore) Since(dataSourceID string, sinceTimeStamp int64) ([]byte, error)](#TSBlobStore.Since)
  * [func (tbs *TSBlobStore) Write(dataSourceID string, payload []byte) error](#TSBlobStore.Write)
  * [func (tbs *TSBlobStore) WriteAt(dataSourceID string, timstamp int64, payload []byte) error](#TSBlobStore.WriteAt)
* [type TSStore](#TSStore)
  * [func (tsc TSStore) Earliest(dataSourceID string) ([]byte, error)](#TSStore.Earliest)
  * [func (tsc TSStore) FirstN(dataSourceID string, n int, opt TimeSeriesQueryOptions) ([]byte, error)](#TSStore.FirstN)
  * [func (tsc TSStore) LastN(dataSourceID string, n int, opt TimeSeriesQueryOptions) ([]byte, error)](#TSStore.LastN)
  * [func (tsc TSStore) Latest(dataSourceID string) ([]byte, error)](#TSStore.Latest)
  * [func (tsc TSStore) Length(dataSourceID string) (int, error)](#TSStore.Length)
  * [func (tsc TSStore) Observe(dataSourceID string) (&lt;-chan []byte, error)](#TSStore.Observe)
  * [func (tsc TSStore) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64, opt TimeSeriesQueryOptions) ([]byte, error)](#TSStore.Range)
  * [func (tsc TSStore) Since(dataSourceID string, sinceTimeStamp int64, opt TimeSeriesQueryOptions) ([]byte, error)](#TSStore.Since)
  * [func (tsc TSStore) Write(dataSourceID string, payload []byte) error](#TSStore.Write)
  * [func (tsc TSStore) WriteAt(dataSourceID string, timstamp int64, payload []byte) error](#TSStore.WriteAt)
* [type TimeSeriesQueryOptions](#TimeSeriesQueryOptions)


#### <a name="pkg-files">Package files</a>
[arbiterClient.go](/src/target/arbiterClient.go) [coreNetworkClient.go](/src/target/coreNetworkClient.go) [coreStoreClient.go](/src/target/coreStoreClient.go) [coreStoreKV.go](/src/target/coreStoreKV.go) [coreStoreTS.go](/src/target/coreStoreTS.go) [coreStoreTSBlob.go](/src/target/coreStoreTSBlob.go) [databoxRequest.go](/src/target/databoxRequest.go) [databoxlog.go](/src/target/databoxlog.go) [export.go](/src/target/export.go) [types.go](/src/target/types.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    Equals            FilterType      = "equals"
    Contains          FilterType      = "contains"
    Sum               AggregationType = "sum"
    Count             AggregationType = "count"
    Min               AggregationType = "min"
    Max               AggregationType = "max"
    Mean              AggregationType = "mean"
    Median            AggregationType = "median"
    StandardDeviation AggregationType = "sd"
)
```
Allowed values for FilterType and AggregationFunction




## <a name="ChkErr">func</a> [ChkErr](/src/target/databoxlog.go?s=1841:1863#L85)
``` go
func ChkErr(err error)
```


## <a name="ChkErrFatal">func</a> [ChkErrFatal](/src/target/databoxlog.go?s=1916:1943#L92)
``` go
func ChkErrFatal(err error)
```


## <a name="Debug">func</a> [Debug](/src/target/databoxlog.go?s=2601:2623#L129)
``` go
func Debug(msg string)
```


## <a name="Err">func</a> [Err](/src/target/databoxlog.go?s=2373:2393#L117)
``` go
func Err(msg string)
```


## <a name="GetStoreURLFromDsHref">func</a> [GetStoreURLFromDsHref](/src/target/coreStoreClient.go?s=6785:6840#L206)
``` go
func GetStoreURLFromDsHref(href string) (string, error)
```
GetStoreURLFromDsHref extracts the base store url from the href provied in the hypercat descriptions.



## <a name="Info">func</a> [Info](/src/target/databoxlog.go?s=2172:2193#L105)
``` go
func Info(msg string)
```


## <a name="NewDataboxHTTPsAPI">func</a> [NewDataboxHTTPsAPI](/src/target/databoxRequest.go?s=108:146#L3)
``` go
func NewDataboxHTTPsAPI() *http.Client
```


## <a name="NewDataboxHTTPsAPIWithPaths">func</a> [NewDataboxHTTPsAPIWithPaths](/src/target/databoxRequest.go?s=254:320#L8)
``` go
func NewDataboxHTTPsAPIWithPaths(cmRootCaPath string) *http.Client
```


## <a name="Warn">func</a> [Warn](/src/target/databoxlog.go?s=2271:2292#L111)
``` go
func Warn(msg string)
```



## <a name="AggregationType">type</a> [AggregationType](/src/target/coreStoreTS.go?s=70:97#L1)
``` go
type AggregationType string
```









## <a name="ArbiterClient">type</a> [ArbiterClient](/src/target/arbiterClient.go?s=171:342#L8)
``` go
type ArbiterClient struct {
    // contains filtered or unexported fields
}
```






### <a name="NewArbiterClient">func</a> [NewArbiterClient](/src/target/arbiterClient.go?s=344:452#L16)
``` go
func NewArbiterClient(arbiterTokenPath string, databoxRequest *http.Client, arbiterURL string) ArbiterClient
```




### <a name="ArbiterClient.GetRootDataSourceCatalogue">func</a> (\*ArbiterClient) [GetRootDataSourceCatalogue](/src/target/arbiterClient.go?s=890:966#L36)
``` go
func (arb *ArbiterClient) GetRootDataSourceCatalogue() (HypercatRoot, error)
```



### <a name="ArbiterClient.GrantComponentPermission">func</a> (\*ArbiterClient) [GrantComponentPermission](/src/target/arbiterClient.go?s=5684:5736#L227)
``` go
func (arb *ArbiterClient) GrantComponentPermission()
```



### <a name="ArbiterClient.GrantContainerPermissions">func</a> (\*ArbiterClient) [GrantContainerPermissions](/src/target/arbiterClient.go?s=2418:2509#L100)
``` go
func (arb *ArbiterClient) GrantContainerPermissions(permissions ContainerPermissions) error
```



### <a name="ArbiterClient.InvalidateCache">func</a> (\*ArbiterClient) [InvalidateCache](/src/target/arbiterClient.go?s=4764:4833#L189)
``` go
func (arb *ArbiterClient) InvalidateCache(href string, method string)
```



### <a name="ArbiterClient.RegesterDataboxComponent">func</a> (\*ArbiterClient) [RegesterDataboxComponent](/src/target/arbiterClient.go?s=5487:5574#L219)
``` go
func (arb *ArbiterClient) RegesterDataboxComponent(componentName, componenttype string)
```



### <a name="ArbiterClient.RemoveDataboxComponent">func</a> (\*ArbiterClient) [RemoveDataboxComponent](/src/target/arbiterClient.go?s=5604:5654#L223)
``` go
func (arb *ArbiterClient) RemoveDataboxComponent()
```



### <a name="ArbiterClient.RequestToken">func</a> (\*ArbiterClient) [RequestToken](/src/target/arbiterClient.go?s=4022:4104#L157)
``` go
func (arb *ArbiterClient) RequestToken(href string, method string) ([]byte, error)
```



### <a name="ArbiterClient.RevokeComponentPermission">func</a> (\*ArbiterClient) [RevokeComponentPermission](/src/target/arbiterClient.go?s=5766:5819#L231)
``` go
func (arb *ArbiterClient) RevokeComponentPermission()
```



### <a name="ArbiterClient.UpdateArbiter">func</a> (\*ArbiterClient) [UpdateArbiter](/src/target/arbiterClient.go?s=1346:1449#L53)
``` go
func (arb *ArbiterClient) UpdateArbiter(name string, tokenString string, databoxType DataboxType) error
```



## <a name="ContainerManagerOptions">type</a> [ContainerManagerOptions](/src/target/types.go?s=89:664#L1)
``` go
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
    InternalIP            string
    ExternalIP            string
    HostPath              string
}
```
ContainerManagerOptions is used to configure the Container Manager










## <a name="ContainerPermissions">type</a> [ContainerPermissions](/src/target/arbiterClient.go?s=2280:2416#L94)
``` go
type ContainerPermissions struct {
    Name    string   `json:"name"`
    Route   Route    `json:"route"`
    Caveats []string `json:"caveats"`
}
```









## <a name="CoreNetworkClient">type</a> [CoreNetworkClient](/src/target/coreNetworkClient.go?s=393:488#L13)
``` go
type CoreNetworkClient struct {
    CM_KEY string
    // contains filtered or unexported fields
}
```






### <a name="NewCoreNetworkClient">func</a> [NewCoreNetworkClient](/src/target/coreNetworkClient.go?s=636:733#L29)
``` go
func NewCoreNetworkClient(containerManagerKeyPath string, request *http.Client) CoreNetworkClient
```




### <a name="CoreNetworkClient.ConnectEndpoints">func</a> (CoreNetworkClient) [ConnectEndpoints](/src/target/coreNetworkClient.go?s=5847:5934#L210)
``` go
func (cnc CoreNetworkClient) ConnectEndpoints(serviceName string, peers []string) error
```



### <a name="CoreNetworkClient.DisconnectEndpoints">func</a> (CoreNetworkClient) [DisconnectEndpoints](/src/target/coreNetworkClient.go?s=6218:6321#L227)
``` go
func (cnc CoreNetworkClient) DisconnectEndpoints(serviceName string, netConfig PostNetworkConfig) error
```



### <a name="CoreNetworkClient.NetworkOfService">func</a> (CoreNetworkClient) [NetworkOfService](/src/target/coreNetworkClient.go?s=3560:3675#L132)
``` go
func (cnc CoreNetworkClient) NetworkOfService(service swarm.Service, serviceName string) (PostNetworkConfig, error)
```



### <a name="CoreNetworkClient.PostUninstall">func</a> (CoreNetworkClient) [PostUninstall](/src/target/coreNetworkClient.go?s=4851:4941#L177)
``` go
func (cnc CoreNetworkClient) PostUninstall(name string, netConfig PostNetworkConfig) error
```



### <a name="CoreNetworkClient.PreConfig">func</a> (CoreNetworkClient) [PreConfig](/src/target/coreNetworkClient.go?s=1116:1204#L49)
``` go
func (cnc CoreNetworkClient) PreConfig(localContainerName string, sla SLA) NetworkConfig
```



### <a name="CoreNetworkClient.RegisterPrivileged">func</a> (CoreNetworkClient) [RegisterPrivileged](/src/target/coreNetworkClient.go?s=6616:6671#L244)
``` go
func (cnc CoreNetworkClient) RegisterPrivileged() error
```



### <a name="CoreNetworkClient.ServiceRestart">func</a> (CoreNetworkClient) [ServiceRestart](/src/target/coreNetworkClient.go?s=6884:6981#L256)
``` go
func (cnc CoreNetworkClient) ServiceRestart(serviceName string, oldIP string, newIP string) error
```



## <a name="CoreStoreClient">type</a> [CoreStoreClient](/src/target/coreStoreClient.go?s=153:461#L5)
``` go
type CoreStoreClient struct {
    ZestC      zest.ZestClient
    Arbiter    *ArbiterClient
    Request    *http.Client
    ZEndpoint  string
    DEndpoint  string
    KVJSON     *KVStore
    KVText     *KVStore
    KVBin      *KVStore
    TSBlobJSON *TSBlobStore
    TSBlobText *TSBlobStore
    TSBlobBin  *TSBlobStore
    TSJSON     *TSStore
}
```






### <a name="NewCoreStoreClient">func</a> [NewCoreStoreClient](/src/target/coreStoreClient.go?s=463:626#L20)
``` go
func NewCoreStoreClient(databoxRequest *http.Client, arbiterClient *ArbiterClient, serverKeyPath string, storeEndPoint string, enableLogging bool) *CoreStoreClient
```




### <a name="CoreStoreClient.GetStoreDataSourceCatalogue">func</a> (\*CoreStoreClient) [GetStoreDataSourceCatalogue](/src/target/coreStoreClient.go?s=1598:1688#L50)
``` go
func (csc *CoreStoreClient) GetStoreDataSourceCatalogue(href string) (HypercatRoot, error)
```



### <a name="CoreStoreClient.HypercatToDataSourceMetadata">func</a> (\*CoreStoreClient) [HypercatToDataSourceMetadata](/src/target/coreStoreClient.go?s=4812:4942#L135)
``` go
func (csc *CoreStoreClient) HypercatToDataSourceMetadata(hypercatDataSourceDescription string) (DataSourceMetadata, string, error)
```
HypercatToDataSourceMetadata is a helper function to convert the hypercat description of a datasource to a DataSourceMetadata instance
Also returns the store url for this data source.




### <a name="CoreStoreClient.RegisterDatasource">func</a> (\*CoreStoreClient) [RegisterDatasource](/src/target/coreStoreClient.go?s=2299:2380#L75)
``` go
func (csc *CoreStoreClient) RegisterDatasource(metadata DataSourceMetadata) error
```
RegisterDatasource is used by apps and drivers to register datasource in stores they
own.




## <a name="DataSource">type</a> [DataSource](/src/target/types.go?s=1203:1502#L44)
``` go
type DataSource struct {
    Type          string       `json:"type"`
    Required      bool         `json:"required"`
    Name          string       `json:"name"`
    Clientid      string       `json:"clientid"`
    Granularities []string     `json:"granularities"`
    Hypercat      HypercatItem `json:"hypercat"`
}
```









## <a name="DataSourceMetadata">type</a> [DataSourceMetadata](/src/target/types.go?s=4147:4389#L97)
``` go
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
```









## <a name="DataboxType">type</a> [DataboxType](/src/target/types.go?s=666:689#L15)
``` go
type DataboxType string
```

``` go
const (
    DataboxTypeApp    DataboxType = "app"
    DataboxTypeDriver DataboxType = "driver"
)
```









## <a name="ExportWhitelist">type</a> [ExportWhitelist](/src/target/types.go?s=1096:1201#L39)
``` go
type ExportWhitelist struct {
    Url         string `json:"url"`
    Description string `json:"description"`
}
```









## <a name="ExternalWhitelist">type</a> [ExternalWhitelist](/src/target/types.go?s=982:1094#L34)
``` go
type ExternalWhitelist struct {
    Urls        []string `json:"urls"`
    Description string   `json:"description"`
}
```









## <a name="Filter">type</a> [Filter](/src/target/coreStoreTS.go?s=692:775#L17)
``` go
type Filter struct {
    TagName    string
    FilterType FilterType
    Value      string
}
```
Filter types to hold the required data to apply the filtering functions of the structured json API










## <a name="FilterType">type</a> [FilterType](/src/target/coreStoreTS.go?s=99:121#L1)
``` go
type FilterType string
```









## <a name="HypercatItem">type</a> [HypercatItem](/src/target/types.go?s=5009:5130#L136)
``` go
type HypercatItem struct {
    ItemMetadata []interface{} `json:"item-metadata"`
    Href         string        `json:"href"`
}
```









## <a name="HypercatRoot">type</a> [HypercatRoot](/src/target/types.go?s=4868:5007#L131)
``` go
type HypercatRoot struct {
    CatalogueMetadata []RelValPair   `json:"catalogue-metadata"`
    Items             []HypercatItem `json:"items"`
}
```









## <a name="KVStore">type</a> [KVStore](/src/target/coreStoreKV.go?s=59:142#L1)
``` go
type KVStore struct {
    // contains filtered or unexported fields
}
```









### <a name="KVStore.Delete">func</a> (\*KVStore) [Delete](/src/target/coreStoreKV.go?s=894:959#L30)
``` go
func (kvj *KVStore) Delete(dataSourceID string, key string) error
```
Delete deletes data under the key.




### <a name="KVStore.DeleteAll">func</a> (\*KVStore) [DeleteAll](/src/target/coreStoreKV.go?s=1117:1173#L39)
``` go
func (kvj *KVStore) DeleteAll(dataSourceID string) error
```
DeleteAll deletes all keys and data from the datasource.




### <a name="KVStore.ListKeys">func</a> (\*KVStore) [ListKeys](/src/target/coreStoreKV.go?s=1327:1394#L48)
``` go
func (kvj *KVStore) ListKeys(dataSourceID string) ([]string, error)
```
ListKeys returns an array of key registed under the dataSourceID




### <a name="KVStore.Observe">func</a> (\*KVStore) [Observe](/src/target/coreStoreKV.go?s=1731:1802#L67)
``` go
func (kvj *KVStore) Observe(dataSourceID string) (<-chan []byte, error)
```



### <a name="KVStore.ObserveKey">func</a> (\*KVStore) [ObserveKey](/src/target/coreStoreKV.go?s=1896:1982#L75)
``` go
func (kvj *KVStore) ObserveKey(dataSourceID string, key string) (<-chan []byte, error)
```



### <a name="KVStore.Read">func</a> (\*KVStore) [Read](/src/target/coreStoreKV.go?s=687:760#L21)
``` go
func (kvj *KVStore) Read(dataSourceID string, key string) ([]byte, error)
```
Read will read the vale store at under tha key
return data is a  object of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="KVStore.Write">func</a> (\*KVStore) [Write](/src/target/coreStoreKV.go?s=353:433#L11)
``` go
func (kvj *KVStore) Write(dataSourceID string, key string, payload []byte) error
```
Write Write will add data to the key value data store.




## <a name="LogEntries">type</a> [LogEntries](/src/target/databoxlog.go?s=126:205#L4)
``` go
type LogEntries struct {
    Msg  string `json:"msg"`
    Type string `json:"type"`
}
```









## <a name="Logger">type</a> [Logger](/src/target/databoxlog.go?s=78:124#L1)
``` go
type Logger struct {
    Store *CoreStoreClient
}
```






### <a name="New">func</a> [New](/src/target/databoxlog.go?s=250:321#L13)
``` go
func New(store *CoreStoreClient, outputDebugLogs bool) (*Logger, error)
```




### <a name="Logger.ChkErr">func</a> (Logger) [ChkErr](/src/target/databoxlog.go?s=1417:1450#L58)
``` go
func (l Logger) ChkErr(err error)
```



### <a name="Logger.Debug">func</a> (Logger) [Debug](/src/target/databoxlog.go?s=1247:1280#L52)
``` go
func (l Logger) Debug(msg string)
```



### <a name="Logger.Err">func</a> (Logger) [Err](/src/target/databoxlog.go?s=1082:1113#L47)
``` go
func (l Logger) Err(msg string)
```



### <a name="Logger.GetLastNLogEntries">func</a> (Logger) [GetLastNLogEntries](/src/target/databoxlog.go?s=1524:1570#L67)
``` go
func (l Logger) GetLastNLogEntries(n int) Logs
```



### <a name="Logger.GetLastNLogEntriesRaw">func</a> (Logger) [GetLastNLogEntriesRaw](/src/target/databoxlog.go?s=1702:1753#L77)
``` go
func (l Logger) GetLastNLogEntriesRaw(n int) []byte
```



### <a name="Logger.Info">func</a> (Logger) [Info](/src/target/databoxlog.go?s=750:782#L37)
``` go
func (l Logger) Info(msg string)
```



### <a name="Logger.Warn">func</a> (Logger) [Warn](/src/target/databoxlog.go?s=916:948#L42)
``` go
func (l Logger) Warn(msg string)
```



## <a name="Logs">type</a> [Logs](/src/target/databoxlog.go?s=207:229#L9)
``` go
type Logs []LogEntries
```









## <a name="Macaroon">type</a> [Macaroon](/src/target/types.go?s=783:803#L22)
``` go
type Macaroon string
```









## <a name="Manifest">type</a> [Manifest](/src/target/types.go?s=1504:2724#L53)
``` go
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
```









## <a name="NetworkConfig">type</a> [NetworkConfig](/src/target/coreNetworkClient.go?s=490:559#L19)
``` go
type NetworkConfig struct {
    NetworkName string
    DNS         string
}
```









## <a name="Package">type</a> [Package](/src/target/types.go?s=886:980#L29)
``` go
type Package struct {
}
```









## <a name="PostNetworkConfig">type</a> [PostNetworkConfig](/src/target/coreNetworkClient.go?s=561:634#L24)
``` go
type PostNetworkConfig struct {
    NetworkName string
    IPv4Address string
}
```









## <a name="RelValPair">type</a> [RelValPair](/src/target/types.go?s=4708:4784#L121)
``` go
type RelValPair struct {
    Rel string `json:"rel"`
    Val string `json:"val"`
}
```









## <a name="RelValPairBool">type</a> [RelValPairBool](/src/target/types.go?s=4786:4866#L126)
``` go
type RelValPairBool struct {
    Rel string `json:"rel"`
    Val bool   `json:"val"`
}
```









## <a name="Repository">type</a> [Repository](/src/target/types.go?s=805:884#L24)
``` go
type Repository struct {
    Type string `json:"Type"`
    Url  string `json:"url"`
}
```









## <a name="ResourceRequirements">type</a> [ResourceRequirements](/src/target/types.go?s=4080:4145#L93)
``` go
type ResourceRequirements struct {
    Store string `json:"store"`
}
```









## <a name="Route">type</a> [Route](/src/target/arbiterClient.go?s=2166:2278#L88)
``` go
type Route struct {
    Target string `json:"target"`
    Path   string `json:"path"`
    Method string `json:"method"`
}
```









## <a name="SLA">type</a> [SLA](/src/target/types.go?s=2726:4078#L72)
``` go
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
```









## <a name="StoreContentType">type</a> [StoreContentType](/src/target/types.go?s=4529:4557#L115)
``` go
type StoreContentType string
```

``` go
const ContentTypeBINARY StoreContentType = "BINARY"
```

``` go
const ContentTypeJSON StoreContentType = "JSON"
```

``` go
const ContentTypeTEXT StoreContentType = "TEXT"
```









## <a name="StoreType">type</a> [StoreType](/src/target/types.go?s=4391:4412#L109)
``` go
type StoreType string
```

``` go
const StoreTypeKV StoreType = "kv"
```

``` go
const StoreTypeTS StoreType = "ts"
```

``` go
const StoreTypeTSBlob StoreType = "ts/blob"
```









## <a name="TSBlobStore">type</a> [TSBlobStore](/src/target/coreStoreTSBlob.go?s=70:157#L1)
``` go
type TSBlobStore struct {
    // contains filtered or unexported fields
}
```









### <a name="TSBlobStore.Earliest">func</a> (\*TSBlobStore) [Earliest](/src/target/coreStoreTSBlob.go?s=1875:1944#L57)
``` go
func (tbs *TSBlobStore) Earliest(dataSourceID string) ([]byte, error)
```
Earliest will retrieve the first entry stored at the requested datasource ID
return data is a byte array contingin  of the format
{"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSBlobStore.FirstN">func</a> (\*TSBlobStore) [FirstN](/src/target/coreStoreTSBlob.go?s=2630:2704#L79)
``` go
func (tbs *TSBlobStore) FirstN(dataSourceID string, n int) ([]byte, error)
```
FirstN will retrieve the first N entries stored at the requested datasource ID
return data is a byte array contingin  of the format
{"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSBlobStore.LastN">func</a> (\*TSBlobStore) [LastN](/src/target/coreStoreTSBlob.go?s=2242:2315#L68)
``` go
func (tbs *TSBlobStore) LastN(dataSourceID string, n int) ([]byte, error)
```
LastN will retrieve the last N entries stored at the requested datasource ID
return data is a byte array contingin  of the format
{"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSBlobStore.Latest">func</a> (\*TSBlobStore) [Latest](/src/target/coreStoreTSBlob.go?s=1512:1579#L46)
``` go
func (tbs *TSBlobStore) Latest(dataSourceID string) ([]byte, error)
```
TSBlobLatest will retrieve the last entry stored at the requested datasource ID
return data is a byte array contingin  of the format
{"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSBlobStore.Length">func</a> (\*TSBlobStore) [Length](/src/target/coreStoreTSBlob.go?s=3834:3898#L110)
``` go
func (tbs *TSBlobStore) Length(dataSourceID string) (int, error)
```
TSBlobLength returns then number of items stored in the timeseries




### <a name="TSBlobStore.Observe">func</a> (\*TSBlobStore) [Observe](/src/target/coreStoreTSBlob.go?s=4424:4499#L134)
``` go
func (tbs *TSBlobStore) Observe(dataSourceID string) (<-chan []byte, error)
```
Observe allows you to get notifications when a new value is written by a driver
the returned chan receives chan []byte continuing json of the
form {"TimestampMS":213123123,"Json":byte[]}




### <a name="TSBlobStore.Range">func</a> (\*TSBlobStore) [Range](/src/target/coreStoreTSBlob.go?s=3476:3582#L101)
``` go
func (tbs *TSBlobStore) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64) ([]byte, error)
```
Range will retrieve all entries between  formTimeStamp and toTimeStamp timestamp in ms since unix epoch
return data is a byte array contingin  of the format
{"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSBlobStore.Since">func</a> (\*TSBlobStore) [Since](/src/target/coreStoreTSBlob.go?s=3025:3113#L90)
``` go
func (tbs *TSBlobStore) Since(dataSourceID string, sinceTimeStamp int64) ([]byte, error)
```
Since will retrieve all entries since the requested timestamp (ms since unix epoch)
return data is a byte array contingin  of the format
{"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSBlobStore.Write">func</a> (\*TSBlobStore) [Write](/src/target/coreStoreTSBlob.go?s=439:511#L12)
``` go
func (tbs *TSBlobStore) Write(dataSourceID string, payload []byte) error
```
Write will add data to the times series data store. Data will be time stamped at insertion (format ms since 1970)




### <a name="TSBlobStore.WriteAt">func</a> (\*TSBlobStore) [WriteAt](/src/target/coreStoreTSBlob.go?s=772:862#L22)
``` go
func (tbs *TSBlobStore) WriteAt(dataSourceID string, timstamp int64, payload []byte) error
```
WriteAt will add data to the times series data store. Data will be time stamped with the timstamp provided in the
timstamp paramiter (format ms since 1970)




## <a name="TSStore">type</a> [TSStore](/src/target/coreStoreTS.go?s=965:1048#L29)
``` go
type TSStore struct {
    // contains filtered or unexported fields
}
```









### <a name="TSStore.Earliest">func</a> (TSStore) [Earliest](/src/target/coreStoreTS.go?s=2659:2723#L82)
``` go
func (tsc TSStore) Earliest(dataSourceID string) ([]byte, error)
```
Earliest will retrieve the first entry stored at the requested datasource ID
return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSStore.FirstN">func</a> (TSStore) [FirstN](/src/target/coreStoreTS.go?s=3443:3540#L102)
``` go
func (tsc TSStore) FirstN(dataSourceID string, n int, opt TimeSeriesQueryOptions) ([]byte, error)
```
FirstN will retrieve the first N entries stored at the requested datasource ID
return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSStore.LastN">func</a> (TSStore) [LastN](/src/target/coreStoreTS.go?s=3014:3110#L92)
``` go
func (tsc TSStore) LastN(dataSourceID string, n int, opt TimeSeriesQueryOptions) ([]byte, error)
```
LastN will retrieve the last N entries stored at the requested datasource ID
return data is an array of JSON objects of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSStore.Latest">func</a> (TSStore) [Latest](/src/target/coreStoreTS.go?s=2319:2381#L72)
``` go
func (tsc TSStore) Latest(dataSourceID string) ([]byte, error)
```
Latest will retrieve the last entry stored at the requested datasource ID
return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSStore.Length">func</a> (TSStore) [Length](/src/target/coreStoreTS.go?s=4748:4807#L131)
``` go
func (tsc TSStore) Length(dataSourceID string) (int, error)
```
Length retruns the number of records stored for that dataSourceID




### <a name="TSStore.Observe">func</a> (TSStore) [Observe](/src/target/coreStoreTS.go?s=5133:5203#L153)
``` go
func (tsc TSStore) Observe(dataSourceID string) (<-chan []byte, error)
```



### <a name="TSStore.Range">func</a> (TSStore) [Range](/src/target/coreStoreTS.go?s=4348:4477#L122)
``` go
func (tsc TSStore) Range(dataSourceID string, formTimeStamp int64, toTimeStamp int64, opt TimeSeriesQueryOptions) ([]byte, error)
```
Range will retrieve all entries between  formTimeStamp and toTimeStamp timestamp in ms since unix epoch
return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSStore.Since">func</a> (TSStore) [Since](/src/target/coreStoreTS.go?s=3867:3978#L112)
``` go
func (tsc TSStore) Since(dataSourceID string, sinceTimeStamp int64, opt TimeSeriesQueryOptions) ([]byte, error)
```
Since will retrieve all entries since the requested timestamp (ms since unix epoch)
return data is a JSON object of the format {"timestamp":213123123,"data":[data-written-by-driver]}




### <a name="TSStore.Write">func</a> (TSStore) [Write](/src/target/coreStoreTS.go?s=1318:1385#L42)
``` go
func (tsc TSStore) Write(dataSourceID string, payload []byte) error
```
Write will add data to the times series data store. Data will be time stamped at insertion (format ms since 1970)




### <a name="TSStore.WriteAt">func</a> (TSStore) [WriteAt](/src/target/coreStoreTS.go?s=1641:1726#L52)
``` go
func (tsc TSStore) WriteAt(dataSourceID string, timstamp int64, payload []byte) error
```
WriteAt will add data to the times series data store. Data will be time stamped with the timstamp provided in the
timstamp paramiter (format ms since 1970)




## <a name="TimeSeriesQueryOptions">type</a> [TimeSeriesQueryOptions](/src/target/coreStoreTS.go?s=859:963#L24)
``` go
type TimeSeriesQueryOptions struct {
    AggregationFunction AggregationType
    Filter              *Filter
}
```
TimeSeriesQueryOptions describes the query options for the structured json API














## Development of databox was supported by the following funding
```
EP/N028260/1, Databox: Privacy-Aware Infrastructure for Managing Personal Data
EP/N028260/2, Databox: Privacy-Aware Infrastructure for Managing Personal Data
EP/N014243/1, Future Everyday Interaction with the Autonomous Internet of Things
EP/M001636/1, Privacy-by-Design: Building Accountability into the Internet of Things EP/M02315X/1, From Human Data to Personal Experience
```
