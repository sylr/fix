package dict

type FIX41 struct {
}

func (f FIX41) Version() string {
	return "FIX.4.1"
}

func (f FIX41) TagName(tag int) (string, bool) {
	switch tag {
	case 1:
		return "Account", true
	case 2:
		return "AdvId", true
	case 3:
		return "AdvRefID", true
	case 4:
		return "AdvSide", true
	case 5:
		return "AdvTransType", true
	case 6:
		return "AvgPx", true
	case 7:
		return "BeginSeqNo", true
	case 8:
		return "BeginString", true
	case 9:
		return "BodyLength", true
	case 10:
		return "CheckSum", true
	case 11:
		return "ClOrdID", true
	case 12:
		return "Commission", true
	case 13:
		return "CommType", true
	case 14:
		return "CumQty", true
	case 15:
		return "Currency", true
	case 16:
		return "EndSeqNo", true
	case 17:
		return "ExecID", true
	case 18:
		return "ExecInst", true
	case 19:
		return "ExecRefID", true
	case 20:
		return "ExecTransType", true
	case 21:
		return "HandlInst", true
	case 22:
		return "IDSource", true
	case 23:
		return "IOIid", true
	case 24:
		return "IOIOthSvc", true
	case 25:
		return "IOIQltyInd", true
	case 26:
		return "IOIRefID", true
	case 27:
		return "IOIShares", true
	case 28:
		return "IOITransType", true
	case 29:
		return "LastCapacity", true
	case 30:
		return "LastMkt", true
	case 31:
		return "LastPx", true
	case 32:
		return "LastShares", true
	case 33:
		return "LinesOfText", true
	case 34:
		return "MsgSeqNum", true
	case 35:
		return "MsgType", true
	case 36:
		return "NewSeqNo", true
	case 37:
		return "OrderID", true
	case 38:
		return "OrderQty", true
	case 39:
		return "OrdStatus", true
	case 40:
		return "OrdType", true
	case 41:
		return "OrigClOrdID", true
	case 42:
		return "OrigTime", true
	case 43:
		return "PossDupFlag", true
	case 44:
		return "Price", true
	case 45:
		return "RefSeqNum", true
	case 46:
		return "RelatdSym", true
	case 47:
		return "Rule80A", true
	case 48:
		return "SecurityID", true
	case 49:
		return "SenderCompID", true
	case 50:
		return "SenderSubID", true
	case 52:
		return "SendingTime", true
	case 53:
		return "Shares", true
	case 54:
		return "Side", true
	case 55:
		return "Symbol", true
	case 56:
		return "TargetCompID", true
	case 57:
		return "TargetSubID", true
	case 58:
		return "Text", true
	case 59:
		return "TimeInForce", true
	case 60:
		return "TransactTime", true
	case 61:
		return "Urgency", true
	case 62:
		return "ValidUntilTime", true
	case 63:
		return "SettlmntTyp", true
	case 64:
		return "FutSettDate", true
	case 65:
		return "SymbolSfx", true
	case 66:
		return "ListID", true
	case 67:
		return "ListSeqNo", true
	case 68:
		return "ListNoOrds", true
	case 69:
		return "ListExecInst", true
	case 70:
		return "AllocID", true
	case 71:
		return "AllocTransType", true
	case 72:
		return "RefAllocID", true
	case 73:
		return "NoOrders", true
	case 74:
		return "AvgPrxPrecision", true
	case 75:
		return "TradeDate", true
	case 76:
		return "ExecBroker", true
	case 77:
		return "OpenClose", true
	case 78:
		return "NoAllocs", true
	case 79:
		return "AllocAccount", true
	case 80:
		return "AllocShares", true
	case 81:
		return "ProcessCode", true
	case 82:
		return "NoRpts", true
	case 83:
		return "RptSeq", true
	case 84:
		return "CxlQty", true
	case 87:
		return "AllocStatus", true
	case 88:
		return "AllocRejCode", true
	case 89:
		return "Signature", true
	case 90:
		return "SecureDataLen", true
	case 91:
		return "SecureData", true
	case 92:
		return "BrokerOfCredit", true
	case 93:
		return "SignatureLength", true
	case 94:
		return "EmailType", true
	case 95:
		return "RawDataLength", true
	case 96:
		return "RawData", true
	case 97:
		return "PossResend", true
	case 98:
		return "EncryptMethod", true
	case 99:
		return "StopPx", true
	case 100:
		return "ExDestination", true
	case 102:
		return "CxlRejReason", true
	case 103:
		return "OrdRejReason", true
	case 104:
		return "IOIQualifier", true
	case 105:
		return "WaveNo", true
	case 106:
		return "Issuer", true
	case 107:
		return "SecurityDesc", true
	case 108:
		return "HeartBtInt", true
	case 109:
		return "ClientID", true
	case 110:
		return "MinQty", true
	case 111:
		return "MaxFloor", true
	case 112:
		return "TestReqID", true
	case 113:
		return "ReportToExch", true
	case 114:
		return "LocateReqd", true
	case 115:
		return "OnBehalfOfCompID", true
	case 116:
		return "OnBehalfOfSubID", true
	case 117:
		return "QuoteID", true
	case 118:
		return "NetMoney", true
	case 119:
		return "SettlCurrAmt", true
	case 120:
		return "SettlCurrency", true
	case 121:
		return "ForexReq", true
	case 122:
		return "OrigSendingTime", true
	case 123:
		return "GapFillFlag", true
	case 124:
		return "NoExecs", true
	case 126:
		return "ExpireTime", true
	case 127:
		return "DKReason", true
	case 128:
		return "DeliverToCompID", true
	case 129:
		return "DeliverToSubID", true
	case 130:
		return "IOINaturalFlag", true
	case 131:
		return "QuoteReqID", true
	case 132:
		return "BidPx", true
	case 133:
		return "OfferPx", true
	case 134:
		return "BidSize", true
	case 135:
		return "OfferSize", true
	case 136:
		return "NoMiscFees", true
	case 137:
		return "MiscFeeAmt", true
	case 138:
		return "MiscFeeCurr", true
	case 139:
		return "MiscFeeType", true
	case 140:
		return "PrevClosePx", true
	case 141:
		return "ResetSeqNumFlag", true
	case 142:
		return "SenderLocationID", true
	case 143:
		return "TargetLocationID", true
	case 144:
		return "OnBehalfOfLocationID", true
	case 145:
		return "DeliverToLocationID", true
	case 146:
		return "NoRelatedSym", true
	case 147:
		return "Subject", true
	case 148:
		return "Headline", true
	case 149:
		return "URLLink", true
	case 150:
		return "ExecType", true
	case 151:
		return "LeavesQty", true
	case 152:
		return "CashOrderQty", true
	case 153:
		return "AllocAvgPx", true
	case 154:
		return "AllocNetMoney", true
	case 155:
		return "SettlCurrFxRate", true
	case 156:
		return "SettlCurrFxRateCalc", true
	case 157:
		return "NumDaysInterest", true
	case 158:
		return "AccruedInterestRate", true
	case 159:
		return "AccruedInterestAmt", true
	case 160:
		return "SettlInstMode", true
	case 161:
		return "AllocText", true
	case 162:
		return "SettlInstID", true
	case 163:
		return "SettlInstTransType", true
	case 164:
		return "EmailThreadID", true
	case 165:
		return "SettlInstSource", true
	case 166:
		return "SettlLocation", true
	case 167:
		return "SecurityType", true
	case 168:
		return "EffectiveTime", true
	case 169:
		return "StandInstDbType", true
	case 170:
		return "StandInstDbName", true
	case 171:
		return "StandInstDbID", true
	case 172:
		return "SettlDeliveryType", true
	case 173:
		return "SettlDepositoryCode", true
	case 174:
		return "SettlBrkrCode", true
	case 175:
		return "SettlInstCode", true
	case 176:
		return "SecuritySettlAgentName", true
	case 177:
		return "SecuritySettlAgentCode", true
	case 178:
		return "SecuritySettlAgentAcctNum", true
	case 179:
		return "SecuritySettlAgentAcctName", true
	case 180:
		return "SecuritySettlAgentContactName", true
	case 181:
		return "SecuritySettlAgentContactPhone", true
	case 182:
		return "CashSettlAgentName", true
	case 183:
		return "CashSettlAgentCode", true
	case 184:
		return "CashSettlAgentAcctNum", true
	case 185:
		return "CashSettlAgentAcctName", true
	case 186:
		return "CashSettlAgentContactName", true
	case 187:
		return "CashSettlAgentContactPhone", true
	case 188:
		return "BidSpotRate", true
	case 189:
		return "BidForwardPoints", true
	case 190:
		return "OfferSpotRate", true
	case 191:
		return "OfferForwardPoints", true
	case 192:
		return "OrderQty2", true
	case 193:
		return "FutSettDate2", true
	case 194:
		return "LastSpotRate", true
	case 195:
		return "LastForwardPoints", true
	case 196:
		return "AllocLinkID", true
	case 197:
		return "AllocLinkType", true
	case 198:
		return "SecondaryOrderID", true
	case 199:
		return "NoIOIQualifiers", true
	case 200:
		return "MaturityMonthYear", true
	case 201:
		return "PutOrCall", true
	case 202:
		return "StrikePrice", true
	case 203:
		return "CoveredOrUncovered", true
	case 204:
		return "CustomerOrFirm", true
	case 205:
		return "MaturityDay", true
	case 206:
		return "OptAttribute", true
	case 207:
		return "SecurityExchange", true
	case 208:
		return "NotifyBrokerOfCredit", true
	case 209:
		return "AllocHandlInst", true
	case 210:
		return "MaxShow", true
	case 211:
		return "PegDifference", true
	default:
		return "", false
	}
}

func (f FIX41) ValueName(tag int, value string) (string, bool) {
	switch tag {
	case 4:
		switch value {
		case "B":
			return `Buy`, true
		case "S":
			return `Sell`, true
		case "T":
			return `Trade`, true
		case "X":
			return `Cross`, true
		}

	case 5:
		switch value {
		case "C":
			return `Cancel`, true
		case "N":
			return `New`, true
		case "R":
			return `Replace`, true
		}

	case 13:
		switch value {
		case "1":
			return `per share`, true
		case "2":
			return `percentage`, true
		case "3":
			return `absolute`, true
		}

	case 18:
		switch value {
		case "0":
			return `Stay on offerside`, true
		case "1":
			return `Not held`, true
		case "2":
			return `Work`, true
		case "3":
			return `Go along`, true
		case "4":
			return `Over the day`, true
		case "5":
			return `Held`, true
		case "6":
			return `Participate don't initiate`, true
		case "7":
			return `Strict scale`, true
		case "8":
			return `Try to scale`, true
		case "9":
			return `Stay on bidside`, true
		case "A":
			return `No cross (cross is forbidden)`, true
		case "B":
			return `OK to cross`, true
		case "C":
			return `Call first`, true
		case "D":
			return `Percent of volume (indicates that the sender does not want to be all of the volume on the floor vs. a specific percentage)`, true
		case "E":
			return `Do not increase - DNI`, true
		case "F":
			return `Do not reduce - DNR`, true
		case "G":
			return `All or none - AON`, true
		case "I":
			return `Institutions only`, true
		case "L":
			return `Last peg (last sale)`, true
		case "M":
			return `Mid-price peg (midprice of inside quote)`, true
		case "N":
			return `Non-negotiable`, true
		case "O":
			return `Opening peg`, true
		case "P":
			return `Market peg`, true
		case "R":
			return `Primary peg (primary market - buy at bid/sell at offer)`, true
		case "S":
			return `Suspend`, true
		case "U":
			return `Customer Display Instruction (Rule11Ac1-1/4)`, true
		case "V":
			return `Netting (for Forex)`, true
		}

	case 20:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Cancel`, true
		case "2":
			return `Correct`, true
		case "3":
			return `Status`, true
		}

	case 21:
		switch value {
		case "1":
			return `Automated execution order, private, no Broker intervention`, true
		case "2":
			return `Automated execution order, public, Broker intervention OK`, true
		case "3":
			return `Manual order, best execution`, true
		}

	case 22:
		switch value {
		case "1":
			return `CUSIP`, true
		case "2":
			return `SEDOL`, true
		case "3":
			return `QUIK`, true
		case "4":
			return `ISIN number`, true
		case "5":
			return `RIC code`, true
		case "6":
			return `ISO Currency Code`, true
		case "7":
			return `ISO Country Code`, true
		}

	case 24:
		switch value {
		case "A":
			return `Autex`, true
		case "B":
			return `Bridge`, true
		}

	case 25:
		switch value {
		case "H":
			return `High`, true
		case "L":
			return `Low`, true
		case "M":
			return `Medium`, true
		}

	case 27:
		switch value {
		case "L":
			return `Large`, true
		case "M":
			return `Medium`, true
		case "S":
			return `Small`, true
		}

	case 28:
		switch value {
		case "C":
			return `Cancel`, true
		case "N":
			return `New`, true
		case "R":
			return `Replace`, true
		}

	case 29:
		switch value {
		case "1":
			return `Agent`, true
		case "2":
			return `Cross as agent`, true
		case "3":
			return `Cross as principal`, true
		case "4":
			return `Principal`, true
		}

	case 35:
		switch value {
		case "0":
			return `Heartbeat`, true
		case "1":
			return `Test Request`, true
		case "2":
			return `Resend Request`, true
		case "3":
			return `Reject`, true
		case "4":
			return `Sequence Reset`, true
		case "5":
			return `Logout`, true
		case "6":
			return `Indication of Interest`, true
		case "7":
			return `Advertisement`, true
		case "8":
			return `Execution Report`, true
		case "9":
			return `Order Cancel Reject`, true
		case "A":
			return `Logon`, true
		case "B":
			return `News`, true
		case "C":
			return `Email`, true
		case "D":
			return `Order - Single`, true
		case "E":
			return `Order - List`, true
		case "F":
			return `Order Cancel Request`, true
		case "G":
			return `Order Cancel/Replace Request`, true
		case "H":
			return `Order Status Request`, true
		case "J":
			return `Allocation`, true
		case "K":
			return `List Cancel Request`, true
		case "L":
			return `List Execute`, true
		case "M":
			return `List Status Request`, true
		case "N":
			return `List Status`, true
		case "P":
			return `Allocation ACK`, true
		case "Q":
			return `Dont Know Trade (DK)`, true
		case "R":
			return `Quote Request`, true
		case "S":
			return `Quote`, true
		case "T":
			return `Settlement Instructions`, true
		}

	case 39:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Partially filled`, true
		case "2":
			return `Filled`, true
		case "3":
			return `Done for day`, true
		case "4":
			return `Canceled`, true
		case "5":
			return `Replaced`, true
		case "6":
			return `Pending Cancel/Replace`, true
		case "7":
			return `Stopped`, true
		case "8":
			return `Rejected`, true
		case "9":
			return `Suspended`, true
		case "A":
			return `Pending New`, true
		case "B":
			return `Calculated`, true
		case "C":
			return `Expired`, true
		}

	case 40:
		switch value {
		case "1":
			return `Market`, true
		case "2":
			return `Limit`, true
		case "3":
			return `Stop`, true
		case "4":
			return `Stop limit`, true
		case "5":
			return `Market on close`, true
		case "6":
			return `With or without`, true
		case "7":
			return `Limit or better`, true
		case "8":
			return `Limit with or without`, true
		case "9":
			return `On basis`, true
		case "A":
			return `On close`, true
		case "B":
			return `Limit on close`, true
		case "C":
			return `Forex - Market`, true
		case "D":
			return `Previously quoted`, true
		case "E":
			return `Previously indicated`, true
		case "F":
			return `Forex - Limit`, true
		case "G":
			return `Forex - Swap`, true
		case "H":
			return `Forex - Previously Quoted`, true
		case "P":
			return `Pegged (requires ExecInst = L, R, M, P or O)`, true
		}

	case 43:
		switch value {
		case "N":
			return `Original transmission`, true
		case "Y":
			return `Possible duplicate`, true
		}

	case 47:
		switch value {
		case "A":
			return `Agency single order`, true
		case "B":
			return `Short exempt transaction (refer to A type)`, true
		case "C":
			return `Program Order, non-index arb, for Member firm/org`, true
		case "D":
			return `Program Order, index arb, for Member firm/org`, true
		case "E":
			return `Registered Equity Market Maker trades`, true
		case "F":
			return `Short exempt transaction (refer to W type)`, true
		case "H":
			return `Short exempt transaction (refer to I type)`, true
		case "I":
			return `Individual Investor, single order`, true
		case "J":
			return `Program Order, index arb, for individual customer`, true
		case "K":
			return `Program Order, non-index arb, for individual customer`, true
		case "L":
			return `Short exempt transaction for member competing market-maker affiliated with the firm clearing the trade (refer to P and O types)`, true
		case "M":
			return `Program Order, index arb, for other member`, true
		case "N":
			return `Program Order, non-index arb, for other member`, true
		case "O":
			return `Competing dealer trades`, true
		case "P":
			return `Principal`, true
		case "R":
			return `Competing dealer trades`, true
		case "S":
			return `Specialist trades`, true
		case "T":
			return `Competing dealer trades`, true
		case "U":
			return `Program Order, index arb, for other agency`, true
		case "W":
			return `All other orders as agent for other member`, true
		case "X":
			return `Short exempt transaction for member competing market-maker not affiliated with the firm clearing the trade (refer to W and T types)`, true
		case "Y":
			return `Program Order, non-index arb, for other agency`, true
		case "Z":
			return `Short exempt transaction for non-member competing market-maker (refer to A and R types)`, true
		}

	case 54:
		switch value {
		case "1":
			return `Buy`, true
		case "2":
			return `Sell`, true
		case "3":
			return `Buy minus`, true
		case "4":
			return `Sell plus`, true
		case "5":
			return `Sell short`, true
		case "6":
			return `Sell short exempt`, true
		case "7":
			return `Undisclosed (valid for IOI message only)`, true
		case "8":
			return `Cross (orders where counterparty is an exchange, valid for all messages except IOIs)`, true
		}

	case 59:
		switch value {
		case "0":
			return `Day`, true
		case "1":
			return `Good Till Cancel (GTC)`, true
		case "2":
			return `At the Opening (OPG)`, true
		case "3":
			return `Immediate or Cancel (OC)`, true
		case "4":
			return `Fill or Kill (FOK)`, true
		case "5":
			return `Good Till Crossing (GTX)`, true
		case "6":
			return `Good Till Date`, true
		}

	case 61:
		switch value {
		case "0":
			return `Normal`, true
		case "1":
			return `Flash`, true
		case "2":
			return `Background`, true
		}

	case 63:
		switch value {
		case "0":
			return `Regular`, true
		case "1":
			return `Cash`, true
		case "2":
			return `Next Day`, true
		case "3":
			return `T+2`, true
		case "4":
			return `T+3`, true
		case "5":
			return `T+4`, true
		case "6":
			return `Future`, true
		case "7":
			return `When Issued`, true
		case "8":
			return `Sellers Option`, true
		case "9":
			return `T+ 5`, true
		}

	case 71:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Replace`, true
		case "2":
			return `Cancel`, true
		case "3":
			return `Preliminary (without MiscFees and NetMoney)`, true
		case "4":
			return `Calculated (includes MiscFees and NetMoney)`, true
		}

	case 77:
		switch value {
		case "C":
			return `Close`, true
		case "O":
			return `Open`, true
		}

	case 81:
		switch value {
		case "0":
			return `regular`, true
		case "1":
			return `soft dollar`, true
		case "2":
			return `step-in`, true
		case "3":
			return `step-out`, true
		case "4":
			return `soft-dollar step-in`, true
		case "5":
			return `soft-dollar step-out`, true
		case "6":
			return `plan sponsor`, true
		}

	case 87:
		switch value {
		case "0":
			return `accepted (successfully processed)`, true
		case "1":
			return `rejected`, true
		case "2":
			return `partial accept`, true
		case "3":
			return `received (received, not yet processed)`, true
		}

	case 88:
		switch value {
		case "0":
			return `unknown account(s)`, true
		case "1":
			return `incorrect quantity`, true
		case "2":
			return `incorrect average price`, true
		case "3":
			return `unknown executing broker mnemonic`, true
		case "4":
			return `commission difference`, true
		case "5":
			return `unknown OrderID`, true
		case "6":
			return `unknown ListID`, true
		case "7":
			return `other`, true
		}

	case 94:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Reply`, true
		case "2":
			return `Admin Reply`, true
		}

	case 97:
		switch value {
		case "N":
			return `Original transmission`, true
		case "Y":
			return `Possible resend`, true
		}

	case 98:
		switch value {
		case "0":
			return `None / other`, true
		case "1":
			return `PKCS (proprietary)`, true
		case "2":
			return `DES (ECB mode)`, true
		case "3":
			return `PKCS/DES (proprietary)`, true
		case "4":
			return `PGP/DES (defunct)`, true
		case "5":
			return `PGP/DES-MD5 (see app note on FIX web site)`, true
		case "6":
			return `PEM/DES-MD5 (see app note on FIX web site)`, true
		}

	case 102:
		switch value {
		case "0":
			return `Too late to cancel`, true
		case "1":
			return `Unknown order`, true
		}

	case 103:
		switch value {
		case "0":
			return `Broker option`, true
		case "1":
			return `Unknown symbol`, true
		case "2":
			return `Exchange closed`, true
		case "3":
			return `Order exceeds limit`, true
		case "4":
			return `Too late to enter`, true
		case "5":
			return `Unknown Order`, true
		case "6":
			return `Duplicate Order`, true
		}

	case 104:
		switch value {
		case "A":
			return `All or none`, true
		case "C":
			return `At the close`, true
		case "I":
			return `In touch with`, true
		case "L":
			return `Limit`, true
		case "M":
			return `More behind`, true
		case "O":
			return `At the open`, true
		case "P":
			return `Taking a position`, true
		case "Q":
			return `At the Market (previously called Current Quote)`, true
		case "S":
			return `Portfolio show-n`, true
		case "T":
			return `Through the day`, true
		case "V":
			return `Versus`, true
		case "W":
			return `Indication - Working away`, true
		case "X":
			return `Crossing opportunity`, true
		case "Y":
			return `At the Midpoint`, true
		case "Z":
			return `Pre-open`, true
		}

	case 113:
		switch value {
		case "N":
			return `Indicates that party sending message will report trade`, true
		case "Y":
			return `Indicates that party receiving message must report trade`, true
		}

	case 114:
		switch value {
		case "N":
			return `Indicates the broker is not required to locate`, true
		case "Y":
			return `Indicates the broker is responsible for locating the stock`, true
		}

	case 121:
		switch value {
		case "N":
			return `Do not execute Forex after security trade`, true
		case "Y":
			return `Execute Forex after security trade`, true
		}

	case 123:
		switch value {
		case "N":
			return `Sequence Reset, ignore MsgSeqNum`, true
		case "Y":
			return `Gap Fill message, MsgSeqNum field valid`, true
		}

	case 127:
		switch value {
		case "A":
			return `Unknown symbol`, true
		case "B":
			return `Wrong side`, true
		case "C":
			return `Quantity exceeds order`, true
		case "D":
			return `No matching order`, true
		case "E":
			return `Price exceeds limit`, true
		case "Z":
			return `Other`, true
		}

	case 130:
		switch value {
		case "N":
			return `Not natural`, true
		case "Y":
			return `Natural`, true
		}

	case 139:
		switch value {
		case "1":
			return `Regulatory (e.g. SEC)`, true
		case "2":
			return `Tax`, true
		case "3":
			return `Local Commission`, true
		case "4":
			return `Exchange Fees`, true
		case "5":
			return `Stamp`, true
		case "6":
			return `Levy`, true
		case "7":
			return `Other`, true
		case "8":
			return `Markup`, true
		}

	case 141:
		switch value {
		case "N":
			return `No`, true
		case "Y":
			return `Yes, reset sequence numbers`, true
		}

	case 150:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Partial fill`, true
		case "2":
			return `Fill`, true
		case "3":
			return `Done for day`, true
		case "4":
			return `Cancelled`, true
		case "5":
			return `Replace`, true
		case "6":
			return `Pending Cancel/Replace`, true
		case "7":
			return `Stopped`, true
		case "8":
			return `Rejected`, true
		case "9":
			return `Suspended`, true
		case "A":
			return `Pending New`, true
		case "B":
			return `Calculated`, true
		case "C":
			return `Expired`, true
		}

	case 160:
		switch value {
		case "0":
			return `Default`, true
		case "1":
			return `Standing Instructions Provided`, true
		case "2":
			return `Specific Allocation Account Overriding`, true
		case "3":
			return `Specific Allocation Account Standing`, true
		}

	case 163:
		switch value {
		case "C":
			return `Cancel`, true
		case "N":
			return `New`, true
		case "R":
			return `Replace`, true
		}

	case 165:
		switch value {
		case "1":
			return `Brokers Instructions`, true
		case "2":
			return `Institutions Instructions`, true
		}

	case 166:
		switch value {
		case "CED":
			return `CEDEL`, true
		case "DTC":
			return `Depository Trust Company`, true
		case "EUR":
			return `Euroclear`, true
		case "FED":
			return `Federal Book Entry`, true
		case "ISO Country Code":
			return `Local Market Settle Location`, true
		case "PNY":
			return `Physical`, true
		case "PTC":
			return `Participant Trust Company`, true
		}

	case 167:
		switch value {
		case "BA":
			return `Bankers Acceptance`, true
		case "CD":
			return `Certificate Of Deposit`, true
		case "CMO":
			return `Collateralize Mortgage Obligation`, true
		case "CORP":
			return `Corporate Bond`, true
		case "CP":
			return `Commercial Paper`, true
		case "CPP":
			return `Corporate Private Placement`, true
		case "CS":
			return `Common Stock`, true
		case "FHA":
			return `Federal Housing Authority`, true
		case "FHL":
			return `Federal Home Loan`, true
		case "FN":
			return `Federal National Mortgage Association`, true
		case "FOR":
			return `Foreign Exchange Contract`, true
		case "FUT":
			return `Future`, true
		case "GN":
			return `Government National Mortgage Association`, true
		case "GOVT":
			return `Treasuries + Agency Debenture`, true
		case "MF":
			return `Mutual Fund`, true
		case "MIO":
			return `Mortgage Interest Only`, true
		case "MPO":
			return `Mortgage Principle Only`, true
		case "MPP":
			return `Mortgage Private Placement`, true
		case "MPT":
			return `Miscellaneous Pass-Thru`, true
		case "MUNI":
			return `Municipal Bond`, true
		case "NONE":
			return `No ISITC Security Type`, true
		case "OPT":
			return `Option`, true
		case "PS":
			return `Preferred Stock`, true
		case "RP":
			return `Repurchase Agreement`, true
		case "RVRP":
			return `Reverse Repurchase Agreement`, true
		case "SL":
			return `Student Loan Marketing Association`, true
		case "TD":
			return `Time Deposit`, true
		case "USTB":
			return `US Treasury Bill`, true
		case "WAR":
			return `Warrant`, true
		case "ZOO":
			return `Cats, Tigers & Lions (a real code: US Treasury Receipts)`, true
		}

	case 169:
		switch value {
		case "0":
			return `Other`, true
		case "1":
			return `DTC SID`, true
		case "2":
			return `Thomson ALERT`, true
		case "3":
			return `A Global Custodian (StandInstDbName must be provided)`, true
		}

	case 197:
		switch value {
		case "0":
			return `F/X Netting`, true
		case "1":
			return `F/X Swap`, true
		}

	case 201:
		switch value {
		case "0":
			return `Put`, true
		case "1":
			return `Call`, true
		}

	case 203:
		switch value {
		case "0":
			return `Covered`, true
		case "1":
			return `Uncovered`, true
		}

	case 204:
		switch value {
		case "0":
			return `Customer`, true
		case "1":
			return `Firm`, true
		}

	case 208:
		switch value {
		case "N":
			return `Details should not be communicated`, true
		case "Y":
			return `Details should be communicated`, true
		}

	case 209:
		switch value {
		case "1":
			return `Match`, true
		case "2":
			return `Forward`, true
		case "3":
			return `Forward and Match`, true
		}

	}
	return "", false
}
