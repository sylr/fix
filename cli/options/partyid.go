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
)

// PartyIDOptions are the flags of the Parties component of a message.
type PartyIDOptions struct {
	partyIDs              []string
	partyIDSources        []string
	partySubIDs           []string
	partySubIDTypes       []string
	partyRoles            []string
	partyRoleQualifiers   []string
	copyPartyIDFromConfig bool
}

// NewPartyIDOptions adds the party flags to command.
func NewPartyIDOptions(command *cobra.Command) *PartyIDOptions {
	opt := &PartyIDOptions{}

	flags := command.Flags()
	flags.StringSliceVar(&opt.partyIDs, "party-id", []string{}, "Party ids")
	flags.StringSliceVar(&opt.partyIDSources, "party-id-source", []string{}, "Party id sources, one per --party-id")
	flags.StringSliceVar(&opt.partyRoles, "party-role", []string{}, "Party roles, one per --party-id")
	flags.StringSliceVar(&opt.partyRoleQualifiers, "party-role-qualifier", []string{}, "Party role qualifiers, one per --party-id if given")
	flags.StringSliceVar(&opt.partySubIDs, "party-sub-ids", []string{}, "Party sub ids (space separated), one per --party-id if given")
	flags.StringSliceVar(&opt.partySubIDTypes, "party-sub-id-types", []string{}, "Party sub id types (space separated), one per --party-sub-ids")
	flags.BoolVar(&opt.copyPartyIDFromConfig, "copy-credentials-from-config", false, "Add the session's username as a customer account party")
	flags.BoolVar(&opt.copyPartyIDFromConfig, "copy-credentials-fom-config", false, "Add the session's username as a customer account party")
	Must(flags.MarkDeprecated("copy-credentials-fom-config", "use --copy-credentials-from-config"))

	Must(command.RegisterFlagCompletionFunc("party-id", cobra.NoFileCompletions))
	Must(command.RegisterFlagCompletionFunc("party-id-source", complete.OrderPartyIDSource))
	Must(command.RegisterFlagCompletionFunc("party-sub-ids", cobra.NoFileCompletions))
	Must(command.RegisterFlagCompletionFunc("party-sub-id-types", complete.OrderPartySubIDTypes))
	Must(command.RegisterFlagCompletionFunc("party-role", complete.OrderPartyIDRole))
	Must(command.RegisterFlagCompletionFunc("party-role-qualifier", complete.OrderPartyRoleQualifier))

	return opt
}

// Validate checks the party flags are consistent and their values known.
func (o PartyIDOptions) Validate() error {
	if len(o.partyIDs) != len(o.partyIDSources) || len(o.partyIDs) != len(o.partyRoles) {
		return fmt.Errorf("%w: you must provide the same number of --party-id, --party-id-source and --party-role", errors.ErrOptionsInconsistentValues)
	}

	if len(o.partyRoleQualifiers) > 0 && len(o.partyRoleQualifiers) != len(o.partyIDs) {
		return fmt.Errorf("%w: you must provide the same number of --party-id and --party-role-qualifier (%d != %d)", errors.ErrOptionsInconsistentValues, len(o.partyIDs), len(o.partyRoleQualifiers))
	}

	if len(o.partySubIDs) > 0 && len(o.partySubIDs) != len(o.partyIDs) {
		return fmt.Errorf("%w: you must provide the same number of --party-id and --party-sub-ids", errors.ErrOptionsInconsistentValues)
	}

	if len(o.partySubIDTypes) > 0 && len(o.partySubIDs) == 0 {
		return fmt.Errorf("%w: --party-sub-id-types needs --party-sub-ids", errors.ErrOptions)
	}

	if len(o.partySubIDTypes) > 0 && len(o.partySubIDTypes) != len(o.partySubIDs) {
		return fmt.Errorf("%w: you must provide the same number of --party-sub-ids and --party-sub-id-types", errors.ErrOptionsInconsistentValues)
	}

	for _, v := range o.partyIDSources {
		if _, ok := Lookup(dict.PartyIDSources, v); !ok {
			return fmt.Errorf("%w: `%s`", errors.ErrOptionOrderIDSourceUnknown, v)
		}
	}

	for _, v := range o.partyRoles {
		if _, ok := Lookup(dict.PartyRoles, v); !ok {
			return fmt.Errorf("%w: `%s`", errors.ErrOptionOrderRoleUnknown, v)
		}
	}

	for _, v := range o.partyRoleQualifiers {
		if _, ok := Lookup(dict.PartyRoleQualifiers, v); !ok {
			return fmt.Errorf("%w: `%s`", errors.ErrOptionOrderRoleQualifierUnknown, v)
		}
	}

	// Sub Parties
	for k, types := range o.partySubIDTypes {
		subIDTypes := strings.Split(types, " ")
		if len(subIDTypes) != len(strings.Split(o.partySubIDs[k], " ")) {
			return fmt.Errorf("%w: %s occurrence of --party-sub-ids and --party-sub-id-types do not have same number of elements (space separated)", errors.ErrOptionsInconsistentValues, humanize.Ordinal(k+1))
		}

		for _, v := range subIDTypes {
			if _, ok := Lookup(dict.PartySubIDTypes, v); !ok {
				return fmt.Errorf("%w: `%s`", errors.ErrOptionPartySubIDTypeUnknown, v)
			}
		}
	}

	return nil
}

// PartyEntries returns the NoPartyIDs entries of the flags, session telling
// the username --copy-credentials-from-config adds.
func (o PartyIDOptions) PartyEntries(session *config.Session) []*fix50sp2.NoPartyIDsEntryBuilder {
	entries := make([]*fix50sp2.NoPartyIDsEntryBuilder, 0, len(o.partyIDs)+1)

	for i := range o.partyIDs {
		source, _ := Lookup(dict.PartyIDSources, o.partyIDSources[i])
		role, _ := Lookup(dict.PartyRoles, o.partyRoles[i])
		party := fix50sp2.NewNoPartyIDsEntry().
			SetPartyID(o.partyIDs[i]).
			SetPartyIDSource(source).
			SetPartyRole(role)

		if len(o.partyRoleQualifiers) > 0 && len(o.partyRoleQualifiers[i]) > 0 {
			qualifier, _ := Lookup(dict.PartyRoleQualifiers, o.partyRoleQualifiers[i])
			party.SetPartyRoleQualifier(qualifier)
		}

		if len(o.partySubIDs) > 0 {
			subIDs := strings.Split(o.partySubIDs[i], " ")

			var subIDTypes []string
			if len(o.partySubIDTypes) > 0 {
				subIDTypes = strings.Split(o.partySubIDTypes[i], " ")
			}

			for k := range subIDs {
				subID := fix50sp2.NewNoPartySubIDsEntry().SetPartySubID(subIDs[k])
				if len(subIDTypes) > 0 {
					subIDType, _ := Lookup(dict.PartySubIDTypes, subIDTypes[k])
					subID.SetPartySubIDType(subIDType)
				}
				party.AddNoPartySubIDsEntry(subID)
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
