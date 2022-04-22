package dict

import (
	"testing"
)

func TestPartyIDSourcesMapsSync(t *testing.T) {
	for k, v := range PartyIDSources {
		if rv, ok := PartyIDSourcesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.PartyIDSourcesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.PartyIDSources[%s] = `%s` AND dict.PartyIDSourcesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range PartyIDSourcesReversed {
		if rv, ok := PartyIDSources[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.PartyIDSources", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.PartyIDSourcesReversed[`%s`] = `%s` AND dict.PartyIDSources[`%s`] = `%s`", k, v, v, rv)
		}
	}
}

func TestPartyRolesMapsSync(t *testing.T) {
	for k, v := range PartyRoles {
		if rv, ok := PartyRolesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.PartyRolesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.PartyRoles[%s] = `%s` AND dict.PartyRolesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range PartyRolesReversed {
		if rv, ok := PartyRoles[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.PartyRoles", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.PartyRolesReversed[`%s`] = `%s` AND dict.PartyRoles[`%s`] = `%s`", k, v, v, rv)
		}
	}
}

func TestPartySubIDTypesMapsSync(t *testing.T) {
	for k, v := range PartySubIDTypes {
		if rv, ok := PartySubIDTypesReversed[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.PartySubIDTypesReversed", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.PartySubIDTypes[%s] = `%s` AND dict.PartySubIDTypesReversed[`%s`] = `%s`", k, v, v, rv)
		}
	}

	for k, v := range PartySubIDTypesReversed {
		if rv, ok := PartySubIDTypes[v]; !ok {
			t.Errorf("maps not in sync: `%s` key is missing in dict.PartySubIDTypes", v)
		} else if k != rv {
			t.Errorf("maps not in sync: dict.PartySubIDTypesReversed[`%s`] = `%s` AND dict.PartySubIDTypes[`%s`] = `%s`", k, v, v, rv)
		}
	}
}
