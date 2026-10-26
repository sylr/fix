package statusorder

import (
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/sylr/fix2/fix"
	"github.com/sylr/fix2/fix/dictionary"

	"sylr.dev/fix/v2/cli/complete"
	"sylr.dev/fix/v2/cli/options"
	"sylr.dev/fix/v2/config"
	"sylr.dev/fix/v2/errors"
	"sylr.dev/fix/v2/fix50sp2"
	"sylr.dev/fix/v2/initiator"
	"sylr.dev/fix/v2/utils"
)

var (
	optionOrderID            string
	optionClOrdID            string
	optionSymbol             string
	optionSide               string
	partyIDOptions           *options.PartyIDOptions
	optionExecReportsTimeout time.Duration

	// statusReqID is the OrdStatusReqID of the request, sent if the
	// session's FIX version defines it.
	statusReqID string
)

var StatusOrderCmd = &cobra.Command{
	Use:               "order",
	Short:             "Order status",
	Long:              "Send an order status request after initiating a session with a FIX acceptor.",
	Args:              cobra.ExactArgs(0),
	ValidArgsFunction: cobra.NoFileCompletions,
	PersistentPreRunE: Validate,
	RunE:              Execute,
}

func init() {
	flags := StatusOrderCmd.Flags()
	flags.StringVar(&optionOrderID, "id", "", "Order id (required if clordid empty)")
	flags.StringVar(&optionClOrdID, "clordid", "", "Client order id of the order (required if id empty, and by FIX 4.x)")
	flags.StringVar(&optionClOrdID, "origclordid", "", "Client order id of the order")
	options.Must(flags.MarkDeprecated("origclordid", "use --clordid"))
	flags.StringVar(&optionSide, "side", "", "Order side (buy, sell ... etc)")
	flags.StringVar(&optionSymbol, "symbol", "", "Order symbol")

	partyIDOptions = options.NewPartyIDOptions(StatusOrderCmd)

	flags.DurationVar(&optionExecReportsTimeout, "exec-reports-timeout", 5*time.Second, "Fail if the execution report is not received within timeout (0s wait indefinitely)")

	options.Must(StatusOrderCmd.MarkFlagRequired("side"))
	options.Must(StatusOrderCmd.MarkFlagRequired("symbol"))

	options.Must(StatusOrderCmd.RegisterFlagCompletionFunc("side", complete.OrderSide))
	options.Must(StatusOrderCmd.RegisterFlagCompletionFunc("symbol", cobra.NoFileCompletions))
}

// Validate checks the options, then that the request fits the session,
// before connecting.
func Validate(cmd *cobra.Command, args []string) error {
	if err := validateOptions(); err != nil {
		return err
	}

	s, err := config.CurrentSession()
	if err != nil {
		return err
	}

	request, err := buildMessage(s)
	if err != nil {
		return err
	}

	return s.Check(request)
}

func validateOptions() error {
	if _, err := options.Side(optionSide); err != nil {
		return err
	}

	if len(optionSymbol) == 0 {
		return errors.ErrOptionsNoSymbolGiven
	}

	if len(optionOrderID) == 0 && len(optionClOrdID) == 0 {
		return fmt.Errorf("%w: --id or --clordid required", errors.ErrOptions)
	}

	statusReqID = uuid.NewString()

	return partyIDOptions.Validate()
}

func Execute(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	c, err := initiator.Connect(ctx, initiator.Options{
		MsgTypes: []string{string(fix50sp2.MsgTypeExecutionreport)},
	})
	if err != nil {
		return err
	}
	defer c.Close(true)

	request, err := buildMessage(c.Session)
	if err != nil {
		return err
	}

	if _, err := c.Send(ctx, request); err != nil {
		return err
	}

	h := &handler{out: cmd.OutOrStdout(), dictionary: c.Dictionary, logger: config.GetLogger(), byStatusReqID: definesStatusReqID(c.Session)}

	return c.Await(ctx, initiator.AwaitOptions{What: "execution report", Timeout: optionExecReportsTimeout}, h.handle)
}

// definesStatusReqID reports whether s's FIX version defines OrdStatusReqID
// in an OrderStatusRequest: FIX 4.3 does not.
func definesStatusReqID(s *config.Session) bool {
	return s.Defines(string(fix50sp2.MsgTypeOrderstatusrequest), fix50sp2.TagOrdStatusReqID)
}

func buildMessage(s *config.Session) (*fix50sp2.OrderStatusRequestBuilder, error) {
	if err := s.CheckSupported(); err != nil {
		return nil, err
	}

	// FIX 4.3 and 4.4 require ClOrdID.
	if len(optionClOrdID) == 0 && s.Requires(string(fix50sp2.MsgTypeOrderstatusrequest), fix50sp2.TagClOrdID) {
		return nil, fmt.Errorf("%w: %s requires --clordid", errors.ErrOptions, s.Version())
	}

	side, err := options.Side(optionSide)
	if err != nil {
		return nil, err
	}

	message := fix50sp2.BuildOrderStatusRequest(s.BeginString).
		SetSide(side).
		SetSymbol(optionSymbol)
	if len(optionClOrdID) > 0 {
		message.SetClOrdID(optionClOrdID)
	}
	if len(optionOrderID) > 0 {
		message.SetOrderID(optionOrderID)
	}
	if definesStatusReqID(s) {
		message.SetOrdStatusReqID(statusReqID)
	}
	for _, entry := range partyIDOptions.PartyEntries(s) {
		message.AddNoPartyIDsEntry(entry)
	}

	return message, nil
}

// handler waits for the execution report answering the request.
type handler struct {
	out        io.Writer
	dictionary *dictionary.Dict
	logger     *zerolog.Logger
	// byStatusReqID tells the report is matched by the request's
	// OrdStatusReqID, rather than by the order's ClOrdID or OrderID.
	byStatusReqID bool
}

// handle handles an execution report: it is done with the first that
// answers the request, and fails if it tells the order is unknown.
func (h *handler) handle(m *fix.Message) (bool, error) {
	if !h.matches(m) {
		h.logger.Debug().Msgf("Ignoring execution report of another request: %s", m.Get(fix50sp2.TagClOrdID))
		return false, initiator.ErrSkip
	}

	utils.WriteMessageBodyAsTable(h.out, m, h.dictionary)

	report := fix50sp2.WrapExecutionReport(m)
	ordStatus, ok := report.OrdStatus()
	if !ok {
		return true, fmt.Errorf("%w: missing OrdStatus", errors.ErrFix)
	}

	if reason, _ := report.OrdRejReason(); ordStatus == fix50sp2.OrdStatusRejected && reason == fix50sp2.OrdRejReasonUnknownOrder {
		return true, options.WithText(fmt.Errorf("%w: unknown order", errors.ErrFixOrderRejected), m)
	}

	return true, nil
}

// matches reports whether m answers the request: its OrdStatusReqID is the
// request's if sent and in m, otherwise its ClOrdID or OrderID is the
// order's.
func (h *handler) matches(m *fix.Message) bool {
	if h.byStatusReqID {
		if reqID := m.Get(fix50sp2.TagOrdStatusReqID); reqID != nil {
			return string(reqID) == statusReqID
		}
	}

	if len(optionClOrdID) > 0 {
		return string(m.Get(fix50sp2.TagClOrdID)) == optionClOrdID
	}

	return string(m.Get(fix50sp2.TagOrderID)) == optionOrderID
}
