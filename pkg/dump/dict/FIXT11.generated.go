package dict

type FIXT11 struct {
}

func (f FIXT11) Version() string {
	return "FIXT.1.1"
}

func (f FIXT11) TagName(tag int) (string, bool) {
	switch tag {
	case 7:
		return "BeginSeqNo", true
	case 8:
		return "BeginString", true
	case 9:
		return "BodyLength", true
	case 10:
		return "CheckSum", true
	case 16:
		return "EndSeqNo", true
	case 34:
		return "MsgSeqNum", true
	case 35:
		return "MsgType", true
	case 36:
		return "NewSeqNo", true
	case 43:
		return "PossDupFlag", true
	case 45:
		return "RefSeqNum", true
	case 49:
		return "SenderCompID", true
	case 50:
		return "SenderSubID", true
	case 52:
		return "SendingTime", true
	case 56:
		return "TargetCompID", true
	case 57:
		return "TargetSubID", true
	case 58:
		return "Text", true
	case 89:
		return "Signature", true
	case 90:
		return "SecureDataLen", true
	case 91:
		return "SecureData", true
	case 93:
		return "SignatureLength", true
	case 95:
		return "RawDataLength", true
	case 96:
		return "RawData", true
	case 97:
		return "PossResend", true
	case 98:
		return "EncryptMethod", true
	case 108:
		return "HeartBtInt", true
	case 112:
		return "TestReqID", true
	case 115:
		return "OnBehalfOfCompID", true
	case 116:
		return "OnBehalfOfSubID", true
	case 122:
		return "OrigSendingTime", true
	case 123:
		return "GapFillFlag", true
	case 128:
		return "DeliverToCompID", true
	case 129:
		return "DeliverToSubID", true
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
	case 212:
		return "XmlDataLen", true
	case 213:
		return "XmlData", true
	case 347:
		return "MessageEncoding", true
	case 354:
		return "EncodedTextLen", true
	case 355:
		return "EncodedText", true
	case 369:
		return "LastMsgSeqNumProcessed", true
	case 371:
		return "RefTagID", true
	case 372:
		return "RefMsgType", true
	case 373:
		return "SessionRejectReason", true
	case 383:
		return "MaxMessageSize", true
	case 464:
		return "TestMessageIndicator", true
	case 553:
		return "Username", true
	case 554:
		return "Password", true
	case 627:
		return "NoHops", true
	case 628:
		return "HopCompID", true
	case 629:
		return "HopSendingTime", true
	case 630:
		return "HopRefID", true
	case 789:
		return "NextExpectedMsgSeqNum", true
	case 925:
		return "NewPassword", true
	case 1128:
		return "ApplVerID", true
	case 1129:
		return "CstmApplVerID", true
	case 1130:
		return "RefApplVerID", true
	case 1131:
		return "RefCstmApplVerID", true
	case 1137:
		return "DefaultApplVerID", true
	case 1156:
		return "ApplExtID", true
	case 1400:
		return "EncryptedPasswordMethod", true
	case 1401:
		return "EncryptedPasswordLen", true
	case 1402:
		return "EncryptedPassword", true
	case 1403:
		return "EncryptedNewPasswordLen", true
	case 1404:
		return "EncryptedNewPassword", true
	case 1406:
		return "RefApplExtID", true
	case 1407:
		return "DefaultApplExtID", true
	case 1408:
		return "DefaultCstmApplVerID", true
	case 1409:
		return "SessionStatus", true
	default:
		return "", false
	}
}

func (f FIXT11) ValueName(tag int, value string) (string, bool) {
	switch tag {
	case 43:
		switch value {
		case "N":
			return `Original transmission`, true
		case "Y":
			return `Possible duplicate`, true
		}

	case 97:
		switch value {
		case "N":
			return `Original Transmission`, true
		case "Y":
			return `Possible Resend`, true
		}

	case 98:
		switch value {
		case "0":
			return `None / Other`, true
		case "1":
			return `PKCS (Proprietary)`, true
		case "2":
			return `DES (ECB Mode)`, true
		case "3":
			return `PKCS / DES (Proprietary)`, true
		case "4":
			return `PGP / DES (Defunct)`, true
		case "5":
			return `PGP / DES-MD5 (See app note on FIX web site)`, true
		case "6":
			return `PEM / DES-MD5 (see app note on FIX web site)`, true
		}

	case 123:
		switch value {
		case "N":
			return `Sequence Reset, Ignore Msg Seq Num (N/A For FIXML - Not Used)`, true
		case "Y":
			return `Gap Fill Message, Msg Seq Num Field Valid`, true
		}

	case 141:
		switch value {
		case "N":
			return `No`, true
		case "Y":
			return `Yes, reset sequence numbers`, true
		}

	case 373:
		switch value {
		case "0":
			return `Invalid Tag Number`, true
		case "1":
			return `Required Tag Missing`, true
		case "2":
			return `Tag not defined for this message type`, true
		case "3":
			return `Undefined tag`, true
		case "4":
			return `Tag specified without a value`, true
		case "5":
			return `Value is incorrect (out of range) for this tag`, true
		case "6":
			return `Incorrect data format for value`, true
		case "7":
			return `Decryption problem`, true
		case "8":
			return `Signature problem`, true
		case "9":
			return `CompID problem`, true
		case "10":
			return `SendingTime Accuracy Problem`, true
		case "11":
			return `Invalid MsgType`, true
		case "12":
			return `XML Validation Error`, true
		case "13":
			return `Tag appears more than once`, true
		case "14":
			return `Tag specified out of required order`, true
		case "15":
			return `Repeating group fields out of order`, true
		case "16":
			return `Incorrect NumInGroup count for repeating group`, true
		case "17":
			return `Non "Data" value includes field delimiter (<SOH> character)`, true
		case "18":
			return `Invalid/Unsupported Application Version`, true
		case "99":
			return `Other`, true
		}

	case 464:
		switch value {
		case "N":
			return `Fales (Production)`, true
		case "Y":
			return `True (Test)`, true
		}

	case 1128:
		switch value {
		case "0":
			return `FIX27`, true
		case "1":
			return `FIX30`, true
		case "2":
			return `FIX40`, true
		case "3":
			return `FIX41`, true
		case "4":
			return `FIX42`, true
		case "5":
			return `FIX43`, true
		case "6":
			return `FIX44`, true
		case "7":
			return `FIX50`, true
		case "8":
			return `FIX50SP1`, true
		}

	case 1409:
		switch value {
		case "0":
			return `Session active`, true
		case "1":
			return `Session password changed`, true
		case "2":
			return `Session password due to expire`, true
		case "3":
			return `New session password does not comply with policy`, true
		case "4":
			return `Session logout complete`, true
		case "5":
			return `Invalid username or password`, true
		case "6":
			return `Account locked`, true
		case "7":
			return `Logons are not allowed at this time`, true
		case "8":
			return `Password expired`, true
		}

	}
	return "", false
}
