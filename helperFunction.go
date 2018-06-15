package libDatabox

import (
	"encoding/json"
	"net/url"
)

// DefaultHTTPSCertPath is the defaut loaction where apps and drivers can find the https certivicats needed to offer a secure UI
const DefaultHTTPSCertPath = "/run/secrets/DATABOX.pem"

//DefaultHTTPSRootCertPath contins the Public key of this databoxes Root certificate needed to verify requests to other components (used in )
const DefaultHTTPSRootCertPath = "/run/secrets/DATABOX_ROOT_CA"

const DefaultArbiterKeyPath = "/run/secrets/ARBITER_TOKEN"
const DefaultStorePublicKeyPath = "/run/secrets/ZMQ_PUBLIC_KEY"

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
			dm.ContentType = vals["val"].(string)
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

// GetStoreURLFromDsHref extracts the base store url from the href provied in the hypercat descriptions.
func GetStoreURLFromDsHref(href string) (string, error) {

	u, err := url.Parse(href)
	if err != nil {
		return "", err
	}

	return u.Scheme + "://" + u.Host, nil

}
