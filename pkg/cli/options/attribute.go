package options

import (
	"fmt"
	"strings"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/spf13/cobra"

	"sylr.dev/fix/pkg/cli/complete"
	"sylr.dev/fix/pkg/dict"
	"sylr.dev/fix/pkg/errors"
	"sylr.dev/fix/pkg/utils"
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

func (o AttributeOptions) EnrichMessageBody(messageBody *quickfix.Body) {
	// Attributes
	attributes := quickfix.NewRepeatingGroup(
		tag.NoOrderAttributes,
		quickfix.GroupTemplate{
			quickfix.GroupElement(tag.OrderAttributeType),
			quickfix.GroupElement(tag.OrderAttributeValue),
		},
	)

	for i := range o.types {
		attribute := attributes.Add()
		attribute.Set(field.NewOrderAttributeType(enum.OrderAttributeType(dict.OrderAttributeTypes[strings.ToUpper(o.types[i])])))

		if len(o.values) > 0 {
			if len(o.values[i]) > 0 {
				attribute.SetString(tag.OrderAttributeValue, o.values[i])
			}
		}
	}

	messageBody.SetGroup(attributes)
}
