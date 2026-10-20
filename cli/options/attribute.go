package options

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/utils"
)

type AttributeOptions struct {
	types  []string
	values []string
}

func NewAttributeOptions(command *cobra.Command) *AttributeOptions {
	opt := &AttributeOptions{}

	command.Flags().StringSliceVar(&opt.types, "attribute-type", []string{}, "Order attribute types")
	command.Flags().StringSliceVar(&opt.values, "attribute-value", []string{}, "Order attribute value")

	command.RegisterFlagCompletionFunc("attribute-type", complete.OrderAttributeType)
	command.RegisterFlagCompletionFunc("attribute-value", cobra.NoFileCompletions)

	return opt
}

func (o AttributeOptions) Validate() error {
	// Attributes
	if len(o.types) != len(o.values) &&
		len(o.values) > 0 {
		return fmt.Errorf("%v: you must provide the same number of --attribute-type and --attribute-values", errors.OptionsInconsistentValues)
	}

	attributeTypes := utils.PrettyOptionValues(dict.OrderAttributeTypes)
	for k := range o.types {
		if utils.Search(attributeTypes, strings.ToLower(o.types[k])) < 0 {
			return fmt.Errorf("%w: `%s`", errors.OptionOrderAttributeTypeUnkonwn, o.types[k])
		}
	}

	return nil
}

func (o AttributeOptions) AttributeEntries() []*fix50sp2.NoOrderAttributesEntryBuilder {
	if len(o.types) == 0 {
		return nil
	}

	entries := make([]*fix50sp2.NoOrderAttributesEntryBuilder, 0, len(o.types))

	for i := range o.types {
		entry := fix50sp2.NewNoOrderAttributesEntry().
			SetOrderAttributeType(dict.OrderAttributeTypes[strings.ToUpper(o.types[i])])

		if len(o.values) > 0 {
			if len(o.values[i]) > 0 {
				entry.SetOrderAttributeValue(o.values[i])
			}
		}

		entries = append(entries, entry)
	}

	return entries
}
