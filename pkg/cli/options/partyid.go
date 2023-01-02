package options

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/spf13/cobra"

	"sylr.dev/fix/config"
	"sylr.dev/fix/pkg/cli/complete"
	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/utils"
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

func (o PartyIdOptions) EnrichMessageBody(messageBody *quickfix.Body, session config.Session) {
	NewNoPartySubIDsRepeatingGroup := func() *quickfix.RepeatingGroup {
		return quickfix.NewRepeatingGroup(
			tag.NoPartySubIDs,
			quickfix.GroupTemplate{
				quickfix.GroupElement(tag.PartySubID),
				quickfix.GroupElement(tag.PartySubIDType),
			},
		)
	}
	parties := quickfix.NewRepeatingGroup(
		tag.NoPartyIDs,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.PartyID),
			quickfix.GroupElement(tag.PartyIDSource),
			quickfix.GroupElement(tag.PartyRole),
			NewNoPartySubIDsRepeatingGroup(),
		},
	)

	for i := range o.partyIDs {
		party := parties.Add()

		party.Set(field.NewPartyID(o.partyIDs[i]))
		party.Set(field.NewPartyIDSource(enum.PartyIDSource(dict.PartyIDSources[strings.ToUpper(o.partyIDSources[i])])))
		party.Set(field.NewPartyRole(enum.PartyRole(dict.PartyRoles[strings.ToUpper(o.partyRoles[i])])))

		// Role Qualifier
		if len(o.partyRoleQualifiers) > 0 && len(o.partyRoleQualifiers[i]) > 0 {
			party.Set(field.NewPartyRoleQualifier(utils.Must(strconv.Atoi(string(dict.PartyRoleQualifiers[strings.ToUpper(o.partyRoleQualifiers[i])])))))
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
				subIDs := NewNoPartySubIDsRepeatingGroup()
				for k := range partySubIDs {
					subID := subIDs.Add()
					mustAdd := false
					if len(partySubIDs) > 0 {
						mustAdd = true
						subID.Set(field.NewPartySubID(partySubIDs[k]))
					}
					if len(partySubIDTypes) > 0 {
						mustAdd = true
						subID.Set(field.NewPartySubIDType(enum.PartySubIDType(dict.PartySubIDTypes[strings.ToUpper(partySubIDTypes[k])])))
					}
					if mustAdd {
						party.SetGroup(subIDs)
					}
				}
			}
		}
	}

	if o.copyPartyIDFromConfig {
		party := parties.Add()
		party.Set(field.NewPartyID(session.Username))
		party.Set(field.NewPartyRole(enum.PartyRole_CUSTOMER_ACCOUNT))
	}

	messageBody.SetGroup(parties)
}
