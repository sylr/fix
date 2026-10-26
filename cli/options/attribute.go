package options

import (
	"fmt"

	"github.com/spf13/cobra"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
)

// AttributeOptions are the flags of the NoOrderAttributes group of an order.
type AttributeOptions struct {
	types  []string
	values []string
}

// NewAttributeOptions adds the order attribute flags to command.
func NewAttributeOptions(command *cobra.Command) *AttributeOptions {
	opt := &AttributeOptions{}

	command.Flags().StringSliceVar(&opt.types, "attribute-type", []string{}, "Order attribute types")
	command.Flags().StringSliceVar(&opt.values, "attribute-value", []string{}, "Order attribute values, one per --attribute-type if given")

	Must(command.RegisterFlagCompletionFunc("attribute-type", complete.OrderAttributeType))
	Must(command.RegisterFlagCompletionFunc("attribute-value", cobra.NoFileCompletions))

	return opt
}

// Validate checks the attribute flags are consistent and their types known.
func (o AttributeOptions) Validate() error {
	if len(o.values) > 0 && len(o.types) != len(o.values) {
		return fmt.Errorf("%w: you must provide the same number of --attribute-type and --attribute-value", errors.ErrOptionsInconsistentValues)
	}

	for _, v := range o.types {
		if _, ok := Lookup(dict.OrderAttributeTypes, v); !ok {
			return fmt.Errorf("%w: `%s`", errors.ErrOptionOrderAttributeTypeUnknown, v)
		}
	}

	return nil
}

// AttributeEntries returns the NoOrderAttributes entries of the flags.
func (o AttributeOptions) AttributeEntries() []*fix50sp2.NoOrderAttributesEntryBuilder {
	entries := make([]*fix50sp2.NoOrderAttributesEntryBuilder, 0, len(o.types))

	for i, v := range o.types {
		attributeType, _ := Lookup(dict.OrderAttributeTypes, v)
		entry := fix50sp2.NewNoOrderAttributesEntry().SetOrderAttributeType(attributeType)

		if len(o.values) > 0 && len(o.values[i]) > 0 {
			entry.SetOrderAttributeValue(o.values[i])
		}

		entries = append(entries, entry)
	}

	return entries
}
