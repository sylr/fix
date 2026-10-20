package listsecurity

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix/session"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/dict"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionType string
)

var ListSecurityCmd = &cobra.Command{
	Use:               "security",
	Aliases:           []string{"securities"},
	Short:             "List securities",
	Long:              "Send a securitylist FIX Message after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: utils.MakePersistentPreRunE(Validate),
	RunE:              Execute,
}

func init() {
	ListSecurityCmd.Flags().StringVar(&optionType, "type", "symbol", "Securities type (symbol, product ... etc)")

	ListSecurityCmd.RegisterFlagCompletionFunc("type", complete.SecurityListRequestType)
}

func Validate(cmd *cobra.Command, args []string) error {
	types := utils.PrettyOptionValues(dict.SecurityListRequestTypes)
	search := utils.Search(types, strings.ToLower(optionType))
	if search < 0 {
		return fmt.Errorf("unknown security type")
	}

	return nil
}

func Execute(cmd *cobra.Command, args []string) error {
	c, err := initiator.Connect(cmd.Context(), initiator.Options{
		MsgTypes: []string{string(fix50sp2.MsgTypeSecuritylist)},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	// Prepare securitylist
	securitylist, err := buildMessage(c)
	if err != nil {
		return err
	}

	// Send the securitylist request
	if _, err = c.Send(cmd.Context(), securitylist); err != nil {
		return err
	}

	// Wait for the response
	select {
	case <-time.After(config.Timeout(0)):
		return errors.ResponseTimeout
	case <-c.Done():
		return errors.FixLogout
	case responseMessage := <-c.Messages():
		utils.WriteMessageBodyAsTable(os.Stdout, responseMessage, c.Dictionary)
	}

	return nil
}

func buildMessage(c *initiator.Client) (session.Outgoing, error) {
	etype, err := dict.SecurityListRequestTypeStringToEnum(optionType)
	if err != nil {
		return nil, err
	}

	if !c.Session.IsFIX50SP2() {
		return nil, errors.FixVersionNotImplemented
	}

	return fix50sp2.BuildSecurityListRequest(c.BeginString()).
		SetSecurityReqID(string(dict.SecurityRequestTypeSymbol)).
		SetSecurityListRequestType(etype), nil
}
