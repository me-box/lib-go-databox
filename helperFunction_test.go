package libDatabox

import (
	"testing"
)

func TestIsActuator(t *testing.T) {

	dsm := DataSource{
		Type:          "databox:container-manager:data",
		Required:      true,
		Name:          "container-manager:data",
		Clientid:      "CM_DATA",
		Granularities: []string{},
		Hypercat: HypercatItem{
			ItemMetadata: []interface{}{
				RelValPair{
					Rel: "urn:X-databox:rels:test",
					Val: "testing",
				},
				RelValPairBool{
					Rel: "urn:X-databox:rels:isActuator",
					Val: true,
				},
				RelValPair{
					Rel: "urn:X-databox:rels:hasDatasourceid",
					Val: "data",
				},
			},
			Href: "tcp://container-manager:5555/kv/data",
		},
	}

	if !IsActuator(dsm) {
		t.Errorf("TestIsActuator failed expected True but got %t", IsActuator(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPair{
				Rel: "urn:X-databox:rels:test",
				Val: "testing",
			},
			RelValPair{
				Rel: "urn:X-databox:rels:hasDatasourceid",
				Val: "data",
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if IsActuator(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsActuator(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPair{
				Rel: "urn:X-databox:rels:test",
				Val: "testing",
			},
			RelValPair{
				Rel: "urn:X-databox:rels:hasDatasourceid",
				Val: "data",
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if IsActuator(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsActuator(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{},
		Href:         "tcp://container-manager:5555/kv/data",
	}

	if IsActuator(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsActuator(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPairBool{
				Rel: "urn:X-databox:rels:isActuator",
				Val: false,
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if IsActuator(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsActuator(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPairBool{
				Rel: "urn:X-databox:rels:isActuator",
				Val: true,
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if !IsActuator(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsActuator(dsm))
	}
}

func TestIsFunc(t *testing.T) {

	dsm := DataSource{
		Type:          "databox:container-manager:data",
		Required:      true,
		Name:          "container-manager:data",
		Clientid:      "CM_DATA",
		Granularities: []string{},
		Hypercat: HypercatItem{
			ItemMetadata: []interface{}{
				RelValPair{
					Rel: "urn:X-databox:rels:test",
					Val: "testing",
				},
				RelValPairBool{
					Rel: "urn:X-databox:rels:isFunc",
					Val: true,
				},
				RelValPair{
					Rel: "urn:X-databox:rels:hasDatasourceid",
					Val: "data",
				},
			},
			Href: "tcp://container-manager:5555/kv/data",
		},
	}

	if !IsFunc(dsm) {
		t.Errorf("TestIsActuator failed expected True but got %t", IsFunc(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPair{
				Rel: "urn:X-databox:rels:test",
				Val: "testing",
			},
			RelValPair{
				Rel: "urn:X-databox:rels:hasDatasourceid",
				Val: "data",
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if IsFunc(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsFunc(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPair{
				Rel: "urn:X-databox:rels:test",
				Val: "testing",
			},
			RelValPair{
				Rel: "urn:X-databox:rels:hasDatasourceid",
				Val: "data",
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if IsFunc(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsFunc(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{},
		Href:         "tcp://container-manager:5555/kv/data",
	}

	if IsFunc(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsFunc(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPairBool{
				Rel: "urn:X-databox:rels:isFunc",
				Val: false,
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if IsFunc(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsFunc(dsm))
	}

	dsm.Hypercat = HypercatItem{
		ItemMetadata: []interface{}{
			RelValPairBool{
				Rel: "urn:X-databox:rels:isFunc",
				Val: true,
			},
		},
		Href: "tcp://container-manager:5555/kv/data",
	}

	if !IsFunc(dsm) {
		t.Errorf("TestIsActuator failed expected False but got %t", IsFunc(dsm))
	}
}
