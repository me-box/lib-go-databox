package libDatabox

import (
	"encoding/json"
	"net/url"
	"strings"
)

// DefaultHTTPSCertPath is the defaut loaction where apps and drivers can find the https certivicats needed to offer a secure UI
const DefaultHTTPSCertPath = "/run/secrets/DATABOX.pem"

//DefaultHTTPSRootCertPath contins the Public key of this databoxes Root certificate needed to verify requests to other components (used in )
const DefaultHTTPSRootCertPath = "/run/secrets/DATABOX_ROOT_CA"

const DefaultArbiterKeyPath = "/run/secrets/ARBITER_TOKEN"
const DefaultStorePublicKeyPath = "/run/secrets/ZMQ_PUBLIC_KEY"

const DefaultArbiterURI = "tcp://arbiter:4444"

// HypercatToDataSourceMetadata is a helper function to convert the hypercat description of a datasource to a DataSourceMetadata instance
// Also returns the store url for this data source.
func HypercatToDataSourceMetadata(hypercatDataSourceDescription string) (DataSourceMetadata, string, error) {
	dm := DataSourceMetadata{}

	hc := HypercatItem{}
	err := json.Unmarshal([]byte(hypercatDataSourceDescription), &hc)
	if err != nil {
		return dm, "", err
	}

	for _, pair := range hc.ItemMetadata {
		vals := pair.(map[string]interface{})
		if vals["rel"].(string) == "urn:X-hypercat:rels:hasDescription:en" {
			dm.Description = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-hypercat:rels:isContentType" {
			dm.ContentType = StoreContentType(vals["val"].(string))
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasVendor" {
			dm.Vendor = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasType" {
			dm.DataSourceType = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasDatasourceid" {
			dm.DataSourceID = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasStoreType" {
			st := vals["val"].(string)
			switch st {
			case "kv":
				dm.StoreType = StoreTypeKV
				break
			case "ts":
				dm.StoreType = StoreTypeTS
				break
			case "ts/blob":
				dm.StoreType = StoreTypeTSBlob
				break
			default:
				//some old SLAs will not have this most use TSBlob
				//TODO CHECK THIS AND BE NOISY
				dm.StoreType = StoreTypeTSBlob
			}
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:isActuator" {
			dm.IsActuator = vals["val"].(bool)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:isFunc" {
			dm.IsActuator = vals["val"].(bool)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasLocation" {
			dm.Location = vals["val"].(string)
			continue
		}
		if vals["rel"].(string) == "urn:X-databox:rels:hasUnit" {
			dm.Unit = vals["val"].(string)
			continue
		}

	}

	url, getStoreURLErr := GetStoreURLFromDsHref(hc.Href)

	return dm, url, getStoreURLErr
}

func IsActuator(dsm DataSource) bool {

	for _, item := range dsm.Hypercat.ItemMetadata {
		switch item.(type) {
		case RelValPairBool:
			if item.(RelValPairBool).Rel == "urn:X-databox:rels:isActuator" && item.(RelValPairBool).Val == true {
				return true
			}
		case RelValPair:
			if item.(RelValPair).Rel == "urn:X-databox:rels:isActuator" && strings.ToLower(item.(RelValPair).Val) == "true" {
				return true
			}
		case interface{}:
			if item.(map[string]interface{})["Val"] != nil && item.(map[string]interface{})["Rel"] != nil {
				val := item.(map[string]interface{})["Val"].(bool)
				rel := item.(map[string]interface{})["Rel"]
				if rel == "urn:X-databox:rels:isActuator" && val == true {
					return true
				}
			}
			if item.(map[string]interface{})["val"] != nil && item.(map[string]interface{})["rel"] != nil {
				val := strings.ToLower(item.(map[string]interface{})["val"].(string))
				rel := item.(map[string]interface{})["rel"]
				if rel == "urn:X-databox:rels:isActuator" && val == "true" {
					return true
				}
			}
		}
	}

	return false
}

func IsFunc(dsm DataSource) bool {

	for _, item := range dsm.Hypercat.ItemMetadata {
		switch item.(type) {
		case RelValPairBool:
			if item.(RelValPairBool).Rel == "urn:X-databox:rels:isFunc" && item.(RelValPairBool).Val == true {
				return true
			}
		case RelValPair:
			if item.(RelValPair).Rel == "urn:X-databox:rels:isFunc" && strings.ToLower(item.(RelValPair).Val) == "true" {
				return true
			}
		case interface{}:
			if item.(map[string]interface{})["Val"] != nil && item.(map[string]interface{})["Rel"] != nil {
				val := item.(map[string]interface{})["Val"].(bool)
				rel := item.(map[string]interface{})["Rel"]
				if rel == "urn:X-databox:rels:isFunc" && val == true {
					return true
				}
			}
			if item.(map[string]interface{})["val"] != nil && item.(map[string]interface{})["rel"] != nil {
				val := strings.ToLower(item.(map[string]interface{})["val"].(string))
				rel := item.(map[string]interface{})["rel"]
				if rel == "urn:X-databox:rels:isFunc" && val == "true" {
					return true
				}
			}
		}
	}

	return false
}

// GetStoreURLFromDsHref extracts the base store url from the href provied in the hypercat descriptions.
func GetStoreURLFromDsHref(href string) (string, error) {

	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil

}

//GetHttpsCredentials Returns a string containing the HTTPS credentials to pass to https server when offering an https server.
//These are read form /run/secrets/DATABOX.pem and are generated by the container-manger at run time.
func GetHttpsCredentials() string {
	return DefaultHTTPSCertPath
}
