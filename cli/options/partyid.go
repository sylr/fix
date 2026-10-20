package options

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/utils"
)

type PartyIdOptions struct {
	partyIDs              []string
	partyIDSources        []string
	partySubIDs           []string
	partySubIDTypes       []string
	partyRoles            []string
	partyRoleQualifiers   []string
	copyPartyIDFromConfig bool
}

func NewPartyIdOptions(command *cobra.Command) *PartyIdOptions {
	opt := &PartyIdOptions{}

	command.Flags().StringSliceVar(&opt.partyIDs, "party-id", []string{}, "Order party ids")
	command.Flags().StringSliceVar(&opt.partyIDSources, "party-id-source", []string{}, "Order party id sources")
	command.Flags().StringSliceVar(&opt.partyRoles, "party-role", []string{}, "Order party roles")
	command.Flags().StringSliceVar(&opt.partyRoleQualifiers, "party-role-qualifier", []string{}, "Order party role qualifiers")
	command.Flags().StringSliceVar(&opt.partySubIDs, "party-sub-ids", []string{}, "Order party sub ids (space separated)")
	command.Flags().StringSliceVar(&opt.partySubIDTypes, "party-sub-id-types", []string{}, "Order party sub id types (space separated)")
	command.Flags().BoolVar(&opt.copyPartyIDFromConfig, "copy-credentials-fom-config", false, "Copy credentials from config in party id fields")

	command.RegisterFlagCompletionFunc("party-id", cobra.NoFileCompletions)
	command.RegisterFlagCompletionFunc("party-id-source", complete.OrderPartyIDSource)
	command.RegisterFlagCompletionFunc("party-sub-ids", cobra.NoFileCompletions)
	command.RegisterFlagCompletionFunc("party-sub-id-types", complete.OrderPartySubIDTypes)
	command.RegisterFlagCompletionFunc("party-role", complete.OrderPartyIDRole)
	command.RegisterFlagCompletionFunc("party-role-qualifier", complete.OrderPartyRoleQualifier)

	return opt
}

func (o PartyIdOptions) Validate() error {
	// Parties
	if len(o.partyIDs) != len(o.partyIDSources) || len(o.partyIDs) != len(o.partyRoles) || len(o.partyIDSources) != len(o.partyRoles) {
		return fmt.Errorf("%v: you must provide the same number of --party-id, --party-id-source, --party-sub-id and --party-role", errors.OptionsInconsistentValues)
	}

	if len(o.partyRoleQualifiers) > 0 && len(o.partyRoleQualifiers) != len(o.partyIDs) {
		return fmt.Errorf("%v: you must provide the same number of --party-id and --party-role-qualifier (%d != %d), %#v", errors.OptionsInconsistentValues, len(o.partyIDs), len(o.partyRoleQualifiers), o.partyRoleQualifiers)
	}

	if len(o.partySubIDs) > 0 && len(o.partySubIDs) != len(o.partyIDs) {
		return fmt.Errorf("%v: you must provide the same number of --party-id and --party-sub-ids", errors.OptionsInconsistentValues)
	}

	if len(o.partySubIDTypes) > 0 && len(o.partySubIDTypes) != len(o.partyIDs) {
		return fmt.Errorf("%v: you must provide the same number of --party-id and --party-sub-id-types", errors.OptionsInconsistentValues)
	}

	sources := utils.PrettyOptionValues(dict.PartyIDSources)
	for k := range o.partyIDSources {
		if utils.Search(sources, strings.ToLower(o.partyIDSources[k])) < 0 {
			return fmt.Errorf("%w: `%s`", errors.OptionOrderIDSourceUnknown, o.partyIDSources[k])
		}
	}

	roles := utils.PrettyOptionValues(dict.PartyRoles)
	for k := range o.partyRoles {
		if utils.Search(roles, strings.ToLower(o.partyRoles[k])) < 0 {
			return fmt.Errorf("%w: `%s`", errors.OptionOrderRoleUnknown, o.partyRoles[k])
		}
	}

	roleQualifiers := utils.PrettyOptionValues(dict.PartyRoleQualifiers)
	for k := range o.partyRoleQualifiers {
		if utils.Search(roleQualifiers, strings.ToLower(o.partyRoleQualifiers[k])) < 0 {
			return fmt.Errorf("%w: `%s`", errors.OptionOrderRoleQualifierUnknown, o.partyRoleQualifiers[k])
		}
	}

	// Sub Parties
	partySubIDTypes := utils.PrettyOptionValues(dict.PartySubIDTypes)
	for k := range o.partySubIDs {
		var subIDs, subIDTypes []string

		if len(o.partySubIDs) > 0 {
			subIDs = strings.Split(o.partySubIDs[k], " ")
		}

		if len(o.partySubIDTypes) > 0 {
			subIDTypes = strings.Split(o.partySubIDTypes[k], " ")
		}

		if len(subIDs) > 0 && len(subIDTypes) > 0 && len(subIDs) != len(subIDTypes) {
			return fmt.Errorf("%v: %s occurence of --party-sub-ids and --party-sub-id-types do not have same number of elements (space separated)", errors.OptionsInconsistentValues, humanize.Ordinal(k))
		}

		if len(subIDTypes) > 0 {
			for kk := range subIDTypes {
				if utils.Search(partySubIDTypes, strings.ToLower(subIDTypes[kk])) < 0 {
					return fmt.Errorf("%w: `%s`", errors.OptionPartySubIDTypeUnknown, subIDTypes[kk])
				}
			}
		}
	}

	return nil
}

func (o PartyIdOptions) PartyEntries(session *config.Session) []*fix50sp2.NoPartyIDsEntryBuilder {
	entries := make([]*fix50sp2.NoPartyIDsEntryBuilder, 0, len(o.partyIDs)+1)

	for i := range o.partyIDs {
		party := fix50sp2.NewNoPartyIDsEntry().
			SetPartyID(o.partyIDs[i]).
			SetPartyIDSource(dict.PartyIDSources[strings.ToUpper(o.partyIDSources[i])]).
			SetPartyRole(dict.PartyRoles[strings.ToUpper(o.partyRoles[i])])

		// Role Qualifier
		if len(o.partyRoleQualifiers) > 0 && len(o.partyRoleQualifiers[i]) > 0 {
			party.SetPartyRoleQualifier(dict.PartyRoleQualifiers[strings.ToUpper(o.partyRoleQualifiers[i])])
		}

		if len(o.partySubIDs) == len(o.partyIDs) || len(o.partySubIDTypes) == len(o.partyIDs) {
			var partySubIDs, partySubIDTypes []string

			if len(o.partySubIDs) > 0 {
				partySubIDs = strings.Split(o.partySubIDs[i], " ")
			}

			if len(o.partySubIDTypes) > 0 {
				partySubIDTypes = strings.Split(o.partySubIDTypes[i], " ")
			}

			if len(partySubIDs) > 0 || len(partySubIDTypes) > 0 {
				for k := range partySubIDs {
					subID := fix50sp2.NewNoPartySubIDsEntry()
					if len(partySubIDs) > 0 {
						subID.SetPartySubID(partySubIDs[k])
					}
					if len(partySubIDTypes) > 0 {
						subID.SetPartySubIDType(dict.PartySubIDTypes[strings.ToUpper(partySubIDTypes[k])])
					}
					party.AddNoPartySubIDsEntry(subID)
				}
			}
		}

		entries = append(entries, party)
	}

	if o.copyPartyIDFromConfig {
		entries = append(entries, fix50sp2.NewNoPartyIDsEntry().
			SetPartyID(session.Username).
			SetPartyRole(fix50sp2.PartyRoleCustomerAccount))
	}

	return entries
}
