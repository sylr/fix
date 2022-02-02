package dict

type FIX43 struct {
}

func (f FIX43) Version() string {
	return "FIX.4.3"
}

func (f FIX43) TagName(tag int) (string, bool) {
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
	case 21:
		return "HandlInst", true
	case 22:
		return "SecurityIDSource", true
	case 23:
		return "IOIid", true
	case 25:
		return "IOIQltyInd", true
	case 26:
		return "IOIRefID", true
	case 27:
		return "IOIQty", true
	case 28:
		return "IOITransType", true
	case 29:
		return "LastCapacity", true
	case 30:
		return "LastMkt", true
	case 31:
		return "LastPx", true
	case 32:
		return "LastQty", true
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
		return "Quantity", true
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
		return "TotNoOrders", true
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
	case 77:
		return "PositionEffect", true
	case 78:
		return "NoAllocs", true
	case 79:
		return "AllocAccount", true
	case 80:
		return "AllocQty", true
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
	case 106:
		return "Issuer", true
	case 107:
		return "SecurityDesc", true
	case 108:
		return "HeartBtInt", true
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
	case 202:
		return "StrikePrice", true
	case 203:
		return "CoveredOrUncovered", true
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
	case 212:
		return "XmlDataLen", true
	case 213:
		return "XmlData", true
	case 214:
		return "SettlInstRefID", true
	case 215:
		return "NoRoutingIDs", true
	case 216:
		return "RoutingType", true
	case 217:
		return "RoutingID", true
	case 218:
		return "Spread", true
	case 219:
		return "Benchmark", true
	case 220:
		return "BenchmarkCurveCurrency", true
	case 221:
		return "BenchmarkCurveName", true
	case 222:
		return "BenchmarkCurvePoint", true
	case 223:
		return "CouponRate", true
	case 224:
		return "CouponPaymentDate", true
	case 225:
		return "IssueDate", true
	case 226:
		return "RepurchaseTerm", true
	case 227:
		return "RepurchaseRate", true
	case 228:
		return "Factor", true
	case 229:
		return "TradeOriginationDate", true
	case 230:
		return "ExDate", true
	case 231:
		return "ContractMultiplier", true
	case 232:
		return "NoStipulations", true
	case 233:
		return "StipulationType", true
	case 234:
		return "StipulationValue", true
	case 235:
		return "YieldType", true
	case 236:
		return "Yield", true
	case 237:
		return "TotalTakedown", true
	case 238:
		return "Concession", true
	case 239:
		return "RepoCollateralSecurityType", true
	case 240:
		return "RedemptionDate", true
	case 241:
		return "UnderlyingCouponPaymentDate", true
	case 242:
		return "UnderlyingIssueDate", true
	case 243:
		return "UnderlyingRepoCollateralSecurityType", true
	case 244:
		return "UnderlyingRepurchaseTerm", true
	case 245:
		return "UnderlyingRepurchaseRate", true
	case 246:
		return "UnderlyingFactor", true
	case 247:
		return "UnderlyingRedemptionDate", true
	case 248:
		return "LegCouponPaymentDate", true
	case 249:
		return "LegIssueDate", true
	case 250:
		return "LegRepoCollateralSecurityType", true
	case 251:
		return "LegRepurchaseTerm", true
	case 252:
		return "LegRepurchaseRate", true
	case 253:
		return "LegFactor", true
	case 254:
		return "LegRedemptionDate", true
	case 255:
		return "CreditRating", true
	case 256:
		return "UnderlyingCreditRating", true
	case 257:
		return "LegCreditRating", true
	case 258:
		return "TradedFlatSwitch", true
	case 259:
		return "BasisFeatureDate", true
	case 260:
		return "BasisFeaturePrice", true
	case 262:
		return "MDReqID", true
	case 263:
		return "SubscriptionRequestType", true
	case 264:
		return "MarketDepth", true
	case 265:
		return "MDUpdateType", true
	case 266:
		return "AggregatedBook", true
	case 267:
		return "NoMDEntryTypes", true
	case 268:
		return "NoMDEntries", true
	case 269:
		return "MDEntryType", true
	case 270:
		return "MDEntryPx", true
	case 271:
		return "MDEntrySize", true
	case 272:
		return "MDEntryDate", true
	case 273:
		return "MDEntryTime", true
	case 274:
		return "TickDirection", true
	case 275:
		return "MDMkt", true
	case 276:
		return "QuoteCondition", true
	case 277:
		return "TradeCondition", true
	case 278:
		return "MDEntryID", true
	case 279:
		return "MDUpdateAction", true
	case 280:
		return "MDEntryRefID", true
	case 281:
		return "MDReqRejReason", true
	case 282:
		return "MDEntryOriginator", true
	case 283:
		return "LocationID", true
	case 284:
		return "DeskID", true
	case 285:
		return "DeleteReason", true
	case 286:
		return "OpenCloseSettleFlag", true
	case 287:
		return "SellerDays", true
	case 288:
		return "MDEntryBuyer", true
	case 289:
		return "MDEntrySeller", true
	case 290:
		return "MDEntryPositionNo", true
	case 291:
		return "FinancialStatus", true
	case 292:
		return "CorporateAction", true
	case 293:
		return "DefBidSize", true
	case 294:
		return "DefOfferSize", true
	case 295:
		return "NoQuoteEntries", true
	case 296:
		return "NoQuoteSets", true
	case 297:
		return "QuoteStatus", true
	case 298:
		return "QuoteCancelType", true
	case 299:
		return "QuoteEntryID", true
	case 300:
		return "QuoteRejectReason", true
	case 301:
		return "QuoteResponseLevel", true
	case 302:
		return "QuoteSetID", true
	case 303:
		return "QuoteRequestType", true
	case 304:
		return "TotQuoteEntries", true
	case 305:
		return "UnderlyingSecurityIDSource", true
	case 306:
		return "UnderlyingIssuer", true
	case 307:
		return "UnderlyingSecurityDesc", true
	case 308:
		return "UnderlyingSecurityExchange", true
	case 309:
		return "UnderlyingSecurityID", true
	case 310:
		return "UnderlyingSecurityType", true
	case 311:
		return "UnderlyingSymbol", true
	case 312:
		return "UnderlyingSymbolSfx", true
	case 313:
		return "UnderlyingMaturityMonthYear", true
	case 315:
		return "UnderlyingPutOrCall", true
	case 316:
		return "UnderlyingStrikePrice", true
	case 317:
		return "UnderlyingOptAttribute", true
	case 320:
		return "SecurityReqID", true
	case 321:
		return "SecurityRequestType", true
	case 322:
		return "SecurityResponseID", true
	case 323:
		return "SecurityResponseType", true
	case 324:
		return "SecurityStatusReqID", true
	case 325:
		return "UnsolicitedIndicator", true
	case 326:
		return "SecurityTradingStatus", true
	case 327:
		return "HaltReason", true
	case 328:
		return "InViewOfCommon", true
	case 329:
		return "DueToRelated", true
	case 330:
		return "BuyVolume", true
	case 331:
		return "SellVolume", true
	case 332:
		return "HighPx", true
	case 333:
		return "LowPx", true
	case 334:
		return "Adjustment", true
	case 335:
		return "TradSesReqID", true
	case 336:
		return "TradingSessionID", true
	case 337:
		return "ContraTrader", true
	case 338:
		return "TradSesMethod", true
	case 339:
		return "TradSesMode", true
	case 340:
		return "TradSesStatus", true
	case 341:
		return "TradSesStartTime", true
	case 342:
		return "TradSesOpenTime", true
	case 343:
		return "TradSesPreCloseTime", true
	case 344:
		return "TradSesCloseTime", true
	case 345:
		return "TradSesEndTime", true
	case 346:
		return "NumberOfOrders", true
	case 347:
		return "MessageEncoding", true
	case 348:
		return "EncodedIssuerLen", true
	case 349:
		return "EncodedIssuer", true
	case 350:
		return "EncodedSecurityDescLen", true
	case 351:
		return "EncodedSecurityDesc", true
	case 352:
		return "EncodedListExecInstLen", true
	case 353:
		return "EncodedListExecInst", true
	case 354:
		return "EncodedTextLen", true
	case 355:
		return "EncodedText", true
	case 356:
		return "EncodedSubjectLen", true
	case 357:
		return "EncodedSubject", true
	case 358:
		return "EncodedHeadlineLen", true
	case 359:
		return "EncodedHeadline", true
	case 360:
		return "EncodedAllocTextLen", true
	case 361:
		return "EncodedAllocText", true
	case 362:
		return "EncodedUnderlyingIssuerLen", true
	case 363:
		return "EncodedUnderlyingIssuer", true
	case 364:
		return "EncodedUnderlyingSecurityDescLen", true
	case 365:
		return "EncodedUnderlyingSecurityDesc", true
	case 366:
		return "AllocPrice", true
	case 367:
		return "QuoteSetValidUntilTime", true
	case 368:
		return "QuoteEntryRejectReason", true
	case 369:
		return "LastMsgSeqNumProcessed", true
	case 370:
		return "OnBehalfOfSendingTime", true
	case 371:
		return "RefTagID", true
	case 372:
		return "RefMsgType", true
	case 373:
		return "SessionRejectReason", true
	case 374:
		return "BidRequestTransType", true
	case 375:
		return "ContraBroker", true
	case 376:
		return "ComplianceID", true
	case 377:
		return "SolicitedFlag", true
	case 378:
		return "ExecRestatementReason", true
	case 379:
		return "BusinessRejectRefID", true
	case 380:
		return "BusinessRejectReason", true
	case 381:
		return "GrossTradeAmt", true
	case 382:
		return "NoContraBrokers", true
	case 383:
		return "MaxMessageSize", true
	case 384:
		return "NoMsgTypes", true
	case 385:
		return "MsgDirection", true
	case 386:
		return "NoTradingSessions", true
	case 387:
		return "TotalVolumeTraded", true
	case 388:
		return "DiscretionInst", true
	case 389:
		return "DiscretionOffset", true
	case 390:
		return "BidID", true
	case 391:
		return "ClientBidID", true
	case 392:
		return "ListName", true
	case 393:
		return "TotalNumSecurities", true
	case 394:
		return "BidType", true
	case 395:
		return "NumTickets", true
	case 396:
		return "SideValue1", true
	case 397:
		return "SideValue2", true
	case 398:
		return "NoBidDescriptors", true
	case 399:
		return "BidDescriptorType", true
	case 400:
		return "BidDescriptor", true
	case 401:
		return "SideValueInd", true
	case 402:
		return "LiquidityPctLow", true
	case 403:
		return "LiquidityPctHigh", true
	case 404:
		return "LiquidityValue", true
	case 405:
		return "EFPTrackingError", true
	case 406:
		return "FairValue", true
	case 407:
		return "OutsideIndexPct", true
	case 408:
		return "ValueOfFutures", true
	case 409:
		return "LiquidityIndType", true
	case 410:
		return "WtAverageLiquidity", true
	case 411:
		return "ExchangeForPhysical", true
	case 412:
		return "OutMainCntryUIndex", true
	case 413:
		return "CrossPercent", true
	case 414:
		return "ProgRptReqs", true
	case 415:
		return "ProgPeriodInterval", true
	case 416:
		return "IncTaxInd", true
	case 417:
		return "NumBidders", true
	case 418:
		return "TradeType", true
	case 419:
		return "BasisPxType", true
	case 420:
		return "NoBidComponents", true
	case 421:
		return "Country", true
	case 422:
		return "TotNoStrikes", true
	case 423:
		return "PriceType", true
	case 424:
		return "DayOrderQty", true
	case 425:
		return "DayCumQty", true
	case 426:
		return "DayAvgPx", true
	case 427:
		return "GTBookingInst", true
	case 428:
		return "NoStrikes", true
	case 429:
		return "ListStatusType", true
	case 430:
		return "NetGrossInd", true
	case 431:
		return "ListOrderStatus", true
	case 432:
		return "ExpireDate", true
	case 433:
		return "ListExecInstType", true
	case 434:
		return "CxlRejResponseTo", true
	case 435:
		return "UnderlyingCouponRate", true
	case 436:
		return "UnderlyingContractMultiplier", true
	case 437:
		return "ContraTradeQty", true
	case 438:
		return "ContraTradeTime", true
	case 441:
		return "LiquidityNumSecurities", true
	case 442:
		return "MultiLegReportingType", true
	case 443:
		return "StrikeTime", true
	case 444:
		return "ListStatusText", true
	case 445:
		return "EncodedListStatusTextLen", true
	case 446:
		return "EncodedListStatusText", true
	case 447:
		return "PartyIDSource", true
	case 448:
		return "PartyID", true
	case 449:
		return "TotalVolumeTradedDate", true
	case 450:
		return "TotalVolumeTradedTime", true
	case 451:
		return "NetChgPrevDay", true
	case 452:
		return "PartyRole", true
	case 453:
		return "NoPartyIDs", true
	case 454:
		return "NoSecurityAltID", true
	case 455:
		return "SecurityAltID", true
	case 456:
		return "SecurityAltIDSource", true
	case 457:
		return "NoUnderlyingSecurityAltID", true
	case 458:
		return "UnderlyingSecurityAltID", true
	case 459:
		return "UnderlyingSecurityAltIDSource", true
	case 460:
		return "Product", true
	case 461:
		return "CFICode", true
	case 462:
		return "UnderlyingProduct", true
	case 463:
		return "UnderlyingCFICode", true
	case 464:
		return "TestMessageIndicator", true
	case 465:
		return "QuantityType", true
	case 466:
		return "BookingRefID", true
	case 467:
		return "IndividualAllocID", true
	case 468:
		return "RoundingDirection", true
	case 469:
		return "RoundingModulus", true
	case 470:
		return "CountryOfIssue", true
	case 471:
		return "StateOrProvinceOfIssue", true
	case 472:
		return "LocaleOfIssue", true
	case 473:
		return "NoRegistDtls", true
	case 474:
		return "MailingDtls", true
	case 475:
		return "InvestorCountryOfResidence", true
	case 476:
		return "PaymentRef", true
	case 477:
		return "DistribPaymentMethod", true
	case 478:
		return "CashDistribCurr", true
	case 479:
		return "CommCurrency", true
	case 480:
		return "CancellationRights", true
	case 481:
		return "MoneyLaunderingStatus", true
	case 482:
		return "MailingInst", true
	case 483:
		return "TransBkdTime", true
	case 484:
		return "ExecPriceType", true
	case 485:
		return "ExecPriceAdjustment", true
	case 486:
		return "DateOfBirth", true
	case 487:
		return "TradeReportTransType", true
	case 488:
		return "CardHolderName", true
	case 489:
		return "CardNumber", true
	case 490:
		return "CardExpDate", true
	case 491:
		return "CardIssNo", true
	case 492:
		return "PaymentMethod", true
	case 493:
		return "RegistAcctType", true
	case 494:
		return "Designation", true
	case 495:
		return "TaxAdvantageType", true
	case 496:
		return "RegistRejReasonText", true
	case 497:
		return "FundRenewWaiv", true
	case 498:
		return "CashDistribAgentName", true
	case 499:
		return "CashDistribAgentCode", true
	case 500:
		return "CashDistribAgentAcctNumber", true
	case 501:
		return "CashDistribPayRef", true
	case 503:
		return "CardStartDate", true
	case 504:
		return "PaymentDate", true
	case 505:
		return "PaymentRemitterID", true
	case 506:
		return "RegistStatus", true
	case 507:
		return "RegistRejReasonCode", true
	case 508:
		return "RegistRefID", true
	case 509:
		return "RegistDetls", true
	case 510:
		return "NoDistribInsts", true
	case 511:
		return "RegistEmail", true
	case 512:
		return "DistribPercentage", true
	case 513:
		return "RegistID", true
	case 514:
		return "RegistTransType", true
	case 515:
		return "ExecValuationPoint", true
	case 516:
		return "OrderPercent", true
	case 517:
		return "OwnershipType", true
	case 518:
		return "NoContAmts", true
	case 519:
		return "ContAmtType", true
	case 520:
		return "ContAmtValue", true
	case 521:
		return "ContAmtCurr", true
	case 522:
		return "OwnerType", true
	case 523:
		return "PartySubID", true
	case 524:
		return "NestedPartyID", true
	case 525:
		return "NestedPartyIDSource", true
	case 526:
		return "SecondaryClOrdID", true
	case 527:
		return "SecondaryExecID", true
	case 528:
		return "OrderCapacity", true
	case 529:
		return "OrderRestrictions", true
	case 530:
		return "MassCancelRequestType", true
	case 531:
		return "MassCancelResponse", true
	case 532:
		return "MassCancelRejectReason", true
	case 533:
		return "TotalAffectedOrders", true
	case 534:
		return "NoAffectedOrders", true
	case 535:
		return "AffectedOrderID", true
	case 536:
		return "AffectedSecondaryOrderID", true
	case 537:
		return "QuoteType", true
	case 538:
		return "NestedPartyRole", true
	case 539:
		return "NoNestedPartyIDs", true
	case 540:
		return "TotalAccruedInterestAmt", true
	case 541:
		return "MaturityDate", true
	case 542:
		return "UnderlyingMaturityDate", true
	case 543:
		return "InstrRegistry", true
	case 544:
		return "CashMargin", true
	case 545:
		return "NestedPartySubID", true
	case 546:
		return "Scope", true
	case 547:
		return "MDImplicitDelete", true
	case 548:
		return "CrossID", true
	case 549:
		return "CrossType", true
	case 550:
		return "CrossPrioritization", true
	case 551:
		return "OrigCrossID", true
	case 552:
		return "NoSides", true
	case 553:
		return "Username", true
	case 554:
		return "Password", true
	case 555:
		return "NoLegs", true
	case 556:
		return "LegCurrency", true
	case 557:
		return "TotalNumSecurityTypes", true
	case 558:
		return "NoSecurityTypes", true
	case 559:
		return "SecurityListRequestType", true
	case 560:
		return "SecurityRequestResult", true
	case 561:
		return "RoundLot", true
	case 562:
		return "MinTradeVol", true
	case 563:
		return "MultiLegRptTypeReq", true
	case 564:
		return "LegPositionEffect", true
	case 565:
		return "LegCoveredOrUncovered", true
	case 566:
		return "LegPrice", true
	case 567:
		return "TradSesStatusRejReason", true
	case 568:
		return "TradeRequestID", true
	case 569:
		return "TradeRequestType", true
	case 570:
		return "PreviouslyReported", true
	case 571:
		return "TradeReportID", true
	case 572:
		return "TradeReportRefID", true
	case 573:
		return "MatchStatus", true
	case 574:
		return "MatchType", true
	case 575:
		return "OddLot", true
	case 576:
		return "NoClearingInstructions", true
	case 577:
		return "ClearingInstruction", true
	case 578:
		return "TradeInputSource", true
	case 579:
		return "TradeInputDevice", true
	case 580:
		return "NoDates", true
	case 581:
		return "AccountType", true
	case 582:
		return "CustOrderCapacity", true
	case 583:
		return "ClOrdLinkID", true
	case 584:
		return "MassStatusReqID", true
	case 585:
		return "MassStatusReqType", true
	case 586:
		return "OrigOrdModTime", true
	case 587:
		return "LegSettlmntTyp", true
	case 588:
		return "LegFutSettDate", true
	case 589:
		return "DayBookingInst", true
	case 590:
		return "BookingUnit", true
	case 591:
		return "PreallocMethod", true
	case 592:
		return "UnderlyingCountryOfIssue", true
	case 593:
		return "UnderlyingStateOrProvinceOfIssue", true
	case 594:
		return "UnderlyingLocaleOfIssue", true
	case 595:
		return "UnderlyingInstrRegistry", true
	case 596:
		return "LegCountryOfIssue", true
	case 597:
		return "LegStateOrProvinceOfIssue", true
	case 598:
		return "LegLocaleOfIssue", true
	case 599:
		return "LegInstrRegistry", true
	case 600:
		return "LegSymbol", true
	case 601:
		return "LegSymbolSfx", true
	case 602:
		return "LegSecurityID", true
	case 603:
		return "LegSecurityIDSource", true
	case 604:
		return "NoLegSecurityAltID", true
	case 605:
		return "LegSecurityAltID", true
	case 606:
		return "LegSecurityAltIDSource", true
	case 607:
		return "LegProduct", true
	case 608:
		return "LegCFICode", true
	case 609:
		return "LegSecurityType", true
	case 610:
		return "LegMaturityMonthYear", true
	case 611:
		return "LegMaturityDate", true
	case 612:
		return "LegStrikePrice", true
	case 613:
		return "LegOptAttribute", true
	case 614:
		return "LegContractMultiplier", true
	case 615:
		return "LegCouponRate", true
	case 616:
		return "LegSecurityExchange", true
	case 617:
		return "LegIssuer", true
	case 618:
		return "EncodedLegIssuerLen", true
	case 619:
		return "EncodedLegIssuer", true
	case 620:
		return "LegSecurityDesc", true
	case 621:
		return "EncodedLegSecurityDescLen", true
	case 622:
		return "EncodedLegSecurityDesc", true
	case 623:
		return "LegRatioQty", true
	case 624:
		return "LegSide", true
	case 625:
		return "TradingSessionSubID", true
	case 626:
		return "AllocType", true
	case 627:
		return "NoHops", true
	case 628:
		return "HopCompID", true
	case 629:
		return "HopSendingTime", true
	case 630:
		return "HopRefID", true
	case 631:
		return "MidPx", true
	case 632:
		return "BidYield", true
	case 633:
		return "MidYield", true
	case 634:
		return "OfferYield", true
	case 635:
		return "ClearingFeeIndicator", true
	case 636:
		return "WorkingIndicator", true
	case 637:
		return "LegLastPx", true
	case 638:
		return "PriorityIndicator", true
	case 639:
		return "PriceImprovement", true
	case 640:
		return "Price2", true
	case 641:
		return "LastForwardPoints2", true
	case 642:
		return "BidForwardPoints2", true
	case 643:
		return "OfferForwardPoints2", true
	case 644:
		return "RFQReqID", true
	case 645:
		return "MktBidPx", true
	case 646:
		return "MktOfferPx", true
	case 647:
		return "MinBidSize", true
	case 648:
		return "MinOfferSize", true
	case 649:
		return "QuoteStatusReqID", true
	case 650:
		return "LegalConfirm", true
	case 651:
		return "UnderlyingLastPx", true
	case 652:
		return "UnderlyingLastQty", true
	case 654:
		return "LegRefID", true
	case 655:
		return "ContraLegRefID", true
	case 656:
		return "SettlCurrBidFxRate", true
	case 657:
		return "SettlCurrOfferFxRate", true
	case 658:
		return "QuoteRequestRejectReason", true
	case 659:
		return "SideComplianceID", true
	default:
		return "", false
	}
}

func (f FIX43) ValueName(tag int, value string) (string, bool) {
	switch tag {
	case 4:
		switch value {
		case "B":
			return `Buy`, true
		case "S":
			return `Sell`, true
		case "X":
			return `Cross`, true
		case "T":
			return `Trade`, true
		}

	case 5:
		switch value {
		case "N":
			return `New`, true
		case "C":
			return `Cancel`, true
		case "R":
			return `Replace`, true
		}

	case 13:
		switch value {
		case "6":
			return `per bond`, true
		case "1":
			return `per share`, true
		case "2":
			return `percentage`, true
		case "3":
			return `absolute`, true
		case "5":
			return `(for CIV buy orders) percentage waived enhanced units`, true
		case "4":
			return `(for CIV buy orders) percentage waived cash discount`, true
		}

	case 18:
		switch value {
		case "Y":
			return `TryToStop`, true
		case "M":
			return `MidPrcPeg`, true
		case "P":
			return `MarkPeg`, true
		case "Q":
			return `CancelOnSysFail`, true
		case "R":
			return `PrimPeg`, true
		case "S":
			return `Suspend`, true
		case "U":
			return `CustDispInst`, true
		case "V":
			return `Netting`, true
		case "W":
			return `PegVWAP`, true
		case "X":
			return `TradeAlong`, true
		case "D":
			return `PercVol`, true
		case "0":
			return `StayOffer`, true
		case "2":
			return `Work`, true
		case "4":
			return `OverDay`, true
		case "5":
			return `Held`, true
		case "6":
			return `PartNotInit`, true
		case "7":
			return `StrictScale`, true
		case "8":
			return `TryToScale`, true
		case "9":
			return `StayBid`, true
		case "A":
			return `NoCross`, true
		case "O":
			return `OpenPeg`, true
		case "C":
			return `CallFirst`, true
		case "N":
			return `NonNego`, true
		case "E":
			return `DNI`, true
		case "F":
			return `DNR`, true
		case "G":
			return `AON`, true
		case "H":
			return `RestateOnSysFail`, true
		case "I":
			return `InstitOnly`, true
		case "J":
			return `RestateOnTradingHalt`, true
		case "K":
			return `CancelOnTradingHalt`, true
		case "L":
			return `LastPeg`, true
		case "3":
			return `GoAlong`, true
		case "B":
			return `OkCross`, true
		case "1":
			return `NotHeld`, true
		}

	case 20:
		switch value {
		case "1":
			return `Cancel`, true
		case "0":
			return `New`, true
		case "3":
			return `Status`, true
		case "2":
			return `Correct`, true
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
		case "E":
			return `Sicovam`, true
		case "2":
			return `SEDOL`, true
		case "1":
			return `CUSIP`, true
		case "3":
			return `QUIK`, true
		case "F":
			return `Belgian`, true
		case "D":
			return `Valoren`, true
		case "C":
			return `Dutch`, true
		case "B":
			return `Wertpapier`, true
		case "A":
			return `Bloomberg Symbol`, true
		case "9":
			return `Consolidated Tape Association (CTA) Symbol (SIAC CTS/CQS line format)`, true
		case "8":
			return `Exchange Symbol`, true
		case "7":
			return `ISO Country Code`, true
		case "6":
			return `ISO Currency Code`, true
		case "5":
			return `RIC code`, true
		case "4":
			return `ISIN number`, true
		case "G":
			return `Common (Clearstream and Euroclear)`, true
		}

	case 25:
		switch value {
		case "M":
			return `Medium`, true
		case "H":
			return `High`, true
		case "L":
			return `Low`, true
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
		case "4":
			return `Principal`, true
		case "3":
			return `Cross as principal`, true
		case "1":
			return `Agent`, true
		case "2":
			return `Cross as agent`, true
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
		case "a":
			return `Quote Status Request`, true
		case "A":
			return `Logon`, true
		case "AA":
			return `Derivative Security List`, true
		case "AB":
			return `New Order - Multileg`, true
		case "AC":
			return `Multileg Order Cancel/Replace (a.k.a. Multileg Order Modification Request)`, true
		case "AD":
			return `Trade Capture Report Request`, true
		case "AE":
			return `Trade Capture Report`, true
		case "AF":
			return `Order Mass Status Request`, true
		case "AG":
			return `Quote Request Reject`, true
		case "AH":
			return `RFQ Request`, true
		case "AI":
			return `Quote Status Report`, true
		case "b":
			return `Mass Quote Acknowledgement`, true
		case "B":
			return `News`, true
		case "c":
			return `Security Definition Request`, true
		case "C":
			return `Email`, true
		case "d":
			return `Security Definition`, true
		case "D":
			return `Order Single`, true
		case "e":
			return `Security Status Request`, true
		case "E":
			return `Order List`, true
		case "f":
			return `Security Status`, true
		case "F":
			return `Order Cancel Request`, true
		case "G":
			return `Order Cancel/Replace Request`, true
		case "g":
			return `Trading Session Status Request`, true
		case "h":
			return `Trading Session Status`, true
		case "H":
			return `Order Status Request`, true
		case "i":
			return `Mass Quote`, true
		case "j":
			return `Business Message Reject`, true
		case "J":
			return `Allocation`, true
		case "K":
			return `List Cancel Request`, true
		case "k":
			return `Bid Request`, true
		case "l":
			return `Bid Response (lowercase L)`, true
		case "L":
			return `List Execute`, true
		case "m":
			return `List Strike Price`, true
		case "M":
			return `List Status Request`, true
		case "N":
			return `List Status`, true
		case "n":
			return `XML message (e.g. non-FIX MsgType)`, true
		case "o":
			return `Registration Instructions`, true
		case "P":
			return `Allocation ACK`, true
		case "p":
			return `Registration Instructions Response`, true
		case "q":
			return `Order Mass Cancel Request`, true
		case "Q":
			return `Dont Know Trade (DK)`, true
		case "r":
			return `Order Mass Cancel Report`, true
		case "R":
			return `Quote Request`, true
		case "s":
			return `New Order - Cross`, true
		case "S":
			return `Quote`, true
		case "t":
			return `Cross Order Cancel/Replace Request (a.k.a. Cross Order Modification Request)`, true
		case "T":
			return `Settlement Instructions`, true
		case "u":
			return `Cross Order Cancel Request`, true
		case "v":
			return `Security Type Request`, true
		case "V":
			return `Market Data Request`, true
		case "w":
			return `Security Types`, true
		case "W":
			return `Market Data-Snapshot/Full Refresh`, true
		case "x":
			return `Security List Request`, true
		case "X":
			return `Market Data-Incremental Refresh`, true
		case "y":
			return `Security List`, true
		case "Y":
			return `Market Data Request Reject`, true
		case "z":
			return `Derivative Security List Request`, true
		case "Z":
			return `Quote Cancel`, true
		}

	case 39:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Partially filled`, true
		case "5":
			return `Replaced (Removed/Replaced)`, true
		case "2":
			return `Filled`, true
		case "6":
			return `Pending Cancel (e.g. result of Order Cancel Request)`, true
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
		case "D":
			return `Accepted for bidding`, true
		case "E":
			return `Pending Replace (e.g. result of Order Cancel/Replace Request)`, true
		case "3":
			return `Done for day`, true
		case "4":
			return `Canceled`, true
		}

	case 40:
		switch value {
		case "D":
			return `Previously quoted`, true
		case "2":
			return `Limit`, true
		case "3":
			return `Stop`, true
		case "4":
			return `Stop limit`, true
		case "5":
			return `Market on close (Deprecated)`, true
		case "6":
			return `With or without`, true
		case "7":
			return `Limit or better`, true
		case "8":
			return `Limit with or without`, true
		case "9":
			return `On basis`, true
		case "A":
			return `On close (Deprecated)`, true
		case "1":
			return `Market`, true
		case "C":
			return `Forex - Market (Deprecated)`, true
		case "F":
			return `Forex - Limit (Deprecated)`, true
		case "E":
			return `Previously indicated`, true
		case "G":
			return `Forex - Swap`, true
		case "I":
			return `Funari (Limit Day Order with unexecuted portion handled as Market On Close. e.g. Japan)`, true
		case "J":
			return `Market If Touched (MIT)`, true
		case "K":
			return `Market with Leftover as Limit (market order then unexecuted quantity becomes limit order at last price)`, true
		case "L":
			return `Previous Fund Valuation Point (Historic pricing) (for CIV)`, true
		case "M":
			return `Next Fund Valuation Point (Forward pricing) (for CIV)`, true
		case "P":
			return `Pegged`, true
		case "B":
			return `Limit on close (Deprecated)`, true
		case "H":
			return `Forex - Previously Quoted (Deprecated)`, true
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
		case "N":
			return `Program Order, non-index arb, for other member`, true
		case "B":
			return `Short exempt transaction (refer to A type)`, true
		case "D":
			return `Program Order, index arb, for Member firm/org`, true
		case "E":
			return `Short Exempt Transaction for Principal (was incorrectly identified in the FIX spec as Registered Equity Market Maker trades)`, true
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
		case "M":
			return `Program Order, index arb, for other member`, true
		case "A":
			return `Agency single order`, true
		case "O":
			return `Proprietary transactions for competing market-maker that is affiliated with the clearing member (was incorrectly identified in the FIX spec as Competing dealer trades)`, true
		case "P":
			return `Principal`, true
		case "R":
			return `Transactions for the account of a non-member competing market maker (was incorrectly identified in the FIX spec as Competing dealer trades)`, true
		case "S":
			return `Specialist trades`, true
		case "T":
			return `Transactions for the account of an unaffiliated members competing market maker (was incorrectly identified in the FIX spec as Competing dealer trades)`, true
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
		case "L":
			return `Short exempt transaction for member competing market-maker affiliated with the firm clearing the trade (refer to P and O types)`, true
		case "C":
			return `Program Order, non-index arb, for Member firm/org`, true
		}

	case 54:
		switch value {
		case "6":
			return `Sell short exempt`, true
		case "B":
			return `As Defined (for use with multileg instruments)`, true
		case "C":
			return `Opposite (for use with multileg instruments)`, true
		case "8":
			return `Cross (orders where counterparty is an exchange, valid for all messages except IOIs)`, true
		case "9":
			return `Cross short`, true
		case "1":
			return `Buy`, true
		case "2":
			return `Sell`, true
		case "3":
			return `Buy minus`, true
		case "4":
			return `Sell plus`, true
		case "A":
			return `Cross short exempt`, true
		case "5":
			return `Sell short`, true
		case "7":
			return `Undisclosed (valid for IOI and List Order messages only)`, true
		}

	case 59:
		switch value {
		case "7":
			return `At the Close`, true
		case "0":
			return `Day`, true
		case "1":
			return `Good Till Cancel (GTC)`, true
		case "2":
			return `At the Opening (OPG)`, true
		case "3":
			return `Immediate or Cancel (IOC)`, true
		case "4":
			return `Fill or Kill (FOK)`, true
		case "5":
			return `Good Till Crossing (GTX)`, true
		case "6":
			return `Good Till Date`, true
		}

	case 61:
		switch value {
		case "1":
			return `Flash`, true
		case "2":
			return `Background`, true
		case "0":
			return `Normal`, true
		}

	case 63:
		switch value {
		case "5":
			return `T+4`, true
		case "A":
			return `T+1`, true
		case "6":
			return `Future`, true
		case "3":
			return `T+2`, true
		case "2":
			return `Next Day`, true
		case "8":
			return `Sellers Option`, true
		case "1":
			return `Cash`, true
		case "7":
			return `When And If Issued`, true
		case "0":
			return `Regular`, true
		case "9":
			return `T+ 5`, true
		case "4":
			return `T+3`, true
		}

	case 71:
		switch value {
		case "5":
			return `Calculated without Preliminary (sent unsolicited by broker, includes MiscFees and NetMoney) (Removed/Replaced)`, true
		case "4":
			return `Calculated (includes MiscFees and NetMoney) (Removed/Replaced)`, true
		case "3":
			return `Preliminary (without MiscFees and NetMoney) (Removed/Replaced)`, true
		case "2":
			return `Cancel`, true
		case "1":
			return `Replace`, true
		case "0":
			return `New`, true
		}

	case 77:
		switch value {
		case "F":
			return `FIFO`, true
		case "R":
			return `Rolled`, true
		case "C":
			return `Close`, true
		case "O":
			return `Open`, true
		}

	case 81:
		switch value {
		case "6":
			return `plan sponsor`, true
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
		}

	case 87:
		switch value {
		case "1":
			return `rejected`, true
		case "2":
			return `partial accept`, true
		case "3":
			return `received (received, not yet processed)`, true
		case "0":
			return `accepted (successfully processed)`, true
		}

	case 88:
		switch value {
		case "0":
			return `unknown account(s)`, true
		case "6":
			return `unknown ListID`, true
		case "3":
			return `unknown executing broker mnemonic`, true
		case "5":
			return `unknown OrderID`, true
		case "7":
			return `other`, true
		case "4":
			return `commission difference`, true
		case "1":
			return `incorrect quantity`, true
		case "2":
			return `incorrect average price`, true
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
		case "2":
			return `DES (ECB mode)`, true
		case "6":
			return `PEM/DES-MD5 (see app note on FIX web site)`, true
		case "5":
			return `PGP/DES-MD5 (see app note on FIX web site)`, true
		case "3":
			return `PKCS/DES (proprietary)`, true
		case "0":
			return `None / other`, true
		case "1":
			return `PKCS (proprietary)`, true
		case "4":
			return `PGP/DES (defunct)`, true
		}

	case 102:
		switch value {
		case "1":
			return `Unknown order`, true
		case "0":
			return `Too late to cancel`, true
		case "6":
			return `Duplicate ClOrdID received`, true
		case "5":
			return `OrigOrdModTime did not match last TransactTime of order`, true
		case "4":
			return `Unable to process Order Mass Cancel Request`, true
		case "3":
			return `Order already in Pending Cancel or Pending Replace status`, true
		case "2":
			return `Broker / Exchange Option`, true
		}

	case 103:
		switch value {
		case "2":
			return `Exchange closed`, true
		case "1":
			return `Unknown symbol`, true
		case "3":
			return `Order exceeds limit`, true
		case "4":
			return `Too late to enter`, true
		case "5":
			return `Unknown Order`, true
		case "7":
			return `Duplicate of a verbally communicated order`, true
		case "9":
			return `Trade Along required`, true
		case "10":
			return `Invalid Investor ID`, true
		case "6":
			return `Duplicate Order (e.g. dupe ClOrdID)`, true
		case "11":
			return `Unsupported order characteristic`, true
		case "12":
			return `Surveillence Option`, true
		case "0":
			return `Broker / Exchange option`, true
		case "8":
			return `Stale Order`, true
		}

	case 104:
		switch value {
		case "O":
			return `At the open`, true
		case "X":
			return `Crossing opportunity`, true
		case "W":
			return `Indication - Working away`, true
		case "V":
			return `Versus`, true
		case "T":
			return `Through the day`, true
		case "S":
			return `Portfolio shown`, true
		case "R":
			return `Ready to trade`, true
		case "A":
			return `All or none`, true
		case "P":
			return `Taking a position`, true
		case "M":
			return `More behind`, true
		case "L":
			return `Limit`, true
		case "I":
			return `In touch with`, true
		case "D":
			return `VWAP (Volume Weighted Avg Price)`, true
		case "C":
			return `At the close (around/not held to close)`, true
		case "B":
			return `Market On Close (MOC) (held to close)`, true
		case "Q":
			return `At the Market (previously called Current Quote)`, true
		case "Y":
			return `At the Midpoint`, true
		case "Z":
			return `Pre-open`, true
		}

	case 113:
		switch value {
		case "Y":
			return `Indicates that party receiving message must report trade`, true
		case "N":
			return `Indicates that party sending message will report trade`, true
		}

	case 114:
		switch value {
		case "Y":
			return `Indicates the broker is responsible for locating the stock`, true
		case "N":
			return `Indicates the broker is not required to locate`, true
		}

	case 121:
		switch value {
		case "Y":
			return `Execute Forex after security trade`, true
		case "N":
			return `Do not execute Forex after security trade`, true
		}

	case 123:
		switch value {
		case "Y":
			return `Gap Fill message, MsgSeqNum field valid`, true
		case "N":
			return `Sequence Reset, ignore MsgSeqNum`, true
		}

	case 127:
		switch value {
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
		case "A":
			return `Unknown symbol`, true
		}

	case 130:
		switch value {
		case "Y":
			return `Natural`, true
		case "N":
			return `Not natural`, true
		}

	case 139:
		switch value {
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
		case "9":
			return `Consumption Tax`, true
		case "1":
			return `Regulatory (e.g. SEC)`, true
		case "2":
			return `Tax`, true
		}

	case 141:
		switch value {
		case "Y":
			return `Yes, reset sequence numbers`, true
		case "N":
			return `No`, true
		}

	case 150:
		switch value {
		case "6":
			return `Pending Cancel (e.g. result of Order Cancel Request)`, true
		case "0":
			return `New`, true
		case "1":
			return `Partial fill (Replaced)`, true
		case "2":
			return `Fill (Replaced)`, true
		case "4":
			return `Canceled`, true
		case "5":
			return `Replace`, true
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
		case "D":
			return `Restated (ExecutionRpt sent unsolicited by sellside, with ExecRestatementReason set)`, true
		case "E":
			return `Pending Replace (e.g. result of Order Cancel/Replace Request)`, true
		case "F":
			return `Trade (partial fill or fill)`, true
		case "G":
			return `Trade Correct (formerly an ExecTransType)`, true
		case "H":
			return `Trade Cancel (formerly an ExecTransType)`, true
		case "I":
			return `Order Status (formerly an ExecTransType)`, true
		case "3":
			return `Done for day`, true
		case "7":
			return `Stopped`, true
		}

	case 156:
		switch value {
		case "D":
			return `Divide`, true
		case "M":
			return `Multiply`, true
		}

	case 160:
		switch value {
		case "0":
			return `Default`, true
		case "4":
			return `Specific Order for a single account (for CIV)`, true
		case "3":
			return `Specific Allocation Account Standing`, true
		case "1":
			return `Standing Instructions Provided`, true
		case "2":
			return `Specific Allocation Account Overriding`, true
		}

	case 163:
		switch value {
		case "N":
			return `New`, true
		case "R":
			return `Replace`, true
		case "C":
			return `Cancel`, true
		}

	case 165:
		switch value {
		case "2":
			return `Institutions Instructions`, true
		case "3":
			return `Investor (e.g. CIV use)`, true
		case "1":
			return `Brokers Instructions`, true
		}

	case 166:
		switch value {
		case "FED":
			return `Federal Book Entry`, true
		case "ISO Country Code":
			return `Local Market Settle Location`, true
		case "PNY":
			return `Physical`, true
		case "EUR":
			return `Euroclear`, true
		case "DTC":
			return `Depository Trust Company`, true
		case "CED":
			return `CEDEL`, true
		case "PTC":
			return `Participant Trust Company`, true
		}

	case 167:
		switch value {
		case "CP":
			return `Commercial Paper`, true
		case "VRDN":
			return `Variable Rate Demand Note`, true
		case "PZFJ":
			return `Plazos Fijos`, true
		case "PN":
			return `Promissory Note`, true
		case "ONITE":
			return `Overnight`, true
		case "MTN":
			return `Medium Term Notes`, true
		case "TECP":
			return `Tax Exempt Commercial Paper`, true
		case "AMENDED":
			return `Amended & Restated`, true
		case "BRIDGE":
			return `Bridge Loan`, true
		case "LOFC":
			return `Letter of Credit`, true
		case "SWING":
			return `Swing Line Facility`, true
		case "DINP":
			return `Debtor in Possession`, true
		case "DEFLTED":
			return `Defaulted`, true
		case "WITHDRN":
			return `Withdrawn`, true
		case "LQN":
			return `Liquidity Note`, true
		case "MATURED":
			return `Matured`, true
		case "DN":
			return `Deposit Notes`, true
		case "RETIRED":
			return `Retired`, true
		case "BA":
			return `Bankers Acceptance`, true
		case "BN":
			return `Bank Notes`, true
		case "BOX":
			return `Bill of Exchanges`, true
		case "CD":
			return `Certificate of Deposit`, true
		case "CL":
			return `Call Loans`, true
		case "REPLACD":
			return `Replaced`, true
		case "MT":
			return `Mandatory Tender`, true
		case "RVLVTRM":
			return `Revolver/Term Loan`, true
		case "MPP":
			return `Mortgage Private Placement`, true
		case "STN":
			return `Short Term Loan Note`, true
		case "MPT":
			return `Miscellaneous Pass-through`, true
		case "TBA":
			return `To be Announced`, true
		case "AN":
			return `Other Anticipation Notes BAN, GAN, etc.`, true
		case "MIO":
			return `Mortgage Interest Only`, true
		case "COFP":
			return `Certificate of Participation`, true
		case "MBS":
			return `Mortgage-backed Securities`, true
		case "REV":
			return `Revenue Bonds`, true
		case "SPCLA":
			return `Special Assessment`, true
		case "SPCLO":
			return `Special Obligation`, true
		case "SPCLT":
			return `Special Tax`, true
		case "TAN":
			return `Tax Anticipation Note`, true
		case "TAXA":
			return `Tax Allocation`, true
		case "COFO":
			return `Certificate of Obligation`, true
		case "TD":
			return `Time Deposit`, true
		case "GO":
			return `General Obligation Bonds`, true
		case "?":
			return `Wildcard entry (used on Security Definition Request message)`, true
		case "WAR":
			return `Warrant`, true
		case "MF":
			return `Mutual Fund (i.e. any kind of open-ended Collective Investment Vehicle)`, true
		case "MLEG":
			return `Multi-leg instrument (e.g. options strategy or futures spread. CFICode can be used to identify if options-based, futures-based, etc.)`, true
		case "TRAN":
			return `Tax & Revenue Anticipation Note`, true
		case "MPO":
			return `Mortgage Principal Only`, true
		case "RP":
			return `Repurchase Agreement`, true
		case "NONE":
			return `No Security Type`, true
		case "XCN":
			return `Extended Comm Note`, true
		case "POOL":
			return `Agency Pools`, true
		case "ABS":
			return `Asset-backed Securities`, true
		case "CMBS":
			return `Corp. Mortgage-backed Securities`, true
		case "CMO":
			return `Collateralized Mortgage Obligation`, true
		case "IET":
			return `IOETTE Mortgage`, true
		case "RVRP":
			return `Reverse Repurchase Agreement`, true
		case "FOR":
			return `Foreign Exchange Contract`, true
		case "RAN":
			return `Revenue Anticipation Note`, true
		case "RVLV":
			return `Revolver Loan`, true
		case "FAC":
			return `Federal Agency Coupon`, true
		case "FADN":
			return `Federal Agency Discount Note`, true
		case "PEF":
			return `Private Export Funding`, true
		case "CORP":
			return `Corporate Bond`, true
		case "CPP":
			return `Corporate Private Placement`, true
		case "CB":
			return `Convertible Bond`, true
		case "DUAL":
			return `Dual Currency`, true
		case "XLINKD":
			return `Indexed Linked`, true
		case "YANK":
			return `Yankee Corporate Bond`, true
		case "CS":
			return `Common Stock`, true
		case "PS":
			return `Preferred Stock`, true
		case "BRADY":
			return `Brady Bond`, true
		case "TBOND":
			return `US Treasury Bond`, true
		case "TINT":
			return `Interest strip from any bond or note`, true
		case "TIPS":
			return `Treasury Inflation Protected Securities`, true
		case "TCAL":
			return `Principal strip of a callable bond or note`, true
		case "TPRN":
			return `Principal strip from a non-callable bond or note`, true
		case "UST":
			return `US Treasury Note/Bond`, true
		case "USTB":
			return `US Treasury Bill`, true
		case "TERM":
			return `Term Loan`, true
		case "STRUCT":
			return `Structured Notes`, true
		}

	case 169:
		switch value {
		case "0":
			return `Other`, true
		case "1":
			return `DTC SID`, true
		case "3":
			return `A Global Custodian (StandInstDbName must be provided)`, true
		case "2":
			return `Thomson ALERT`, true
		}

	case 172:
		switch value {
		case "1":
			return `Free`, true
		case "0":
			return `Versus Payment`, true
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
		case "1":
			return `Uncovered`, true
		case "0":
			return `Covered`, true
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
		case "3":
			return `Forward and Match`, true
		case "2":
			return `Forward`, true
		case "1":
			return `Match`, true
		}

	case 216:
		switch value {
		case "1":
			return `Target Firm`, true
		case "2":
			return `Target List`, true
		case "3":
			return `Block Firm`, true
		case "4":
			return `Block List`, true
		}

	case 219:
		switch value {
		case "5":
			return `OLD-10`, true
		case "1":
			return `CURVE`, true
		case "2":
			return `5-YR`, true
		case "4":
			return `10-YR`, true
		case "6":
			return `30-YR`, true
		case "7":
			return `OLD-30`, true
		case "8":
			return `3-MO-LIBOR`, true
		case "9":
			return `6-MO-LIBOR`, true
		case "3":
			return `OLD-5`, true
		}

	case 221:
		switch value {
		case "SWAP":
			return `SWAP`, true
		case "LIBID":
			return `LIBID`, true
		case "OTHER":
			return `OTHER`, true
		case "Treasury":
			return `Treasury`, true
		case "Euribor":
			return `Euribor`, true
		case "Pfandbriefe":
			return `Pfandbriefe`, true
		case "FutureSWAP":
			return `FutureSWAP`, true
		case "MuniAAA":
			return `MuniAAA`, true
		case "LIBOR":
			return `LIBOR (London Inter-Bank Offers)`, true
		}

	case 233:
		switch value {
		case "ABS":
			return `Absolute Prepayment Speed`, true
		case "WALA":
			return `Weighted Average Loan Age (value in months)`, true
		case "WAM":
			return `Weighted Average Maturity (value in months)`, true
		case "CPR":
			return `Constant Prepayment Rate`, true
		case "HEP":
			return `final CPR of Home Equity Prepayment Curve`, true
		case "WAL":
			return `Weighted Average Life (value in months)`, true
		case "MHP":
			return `% of Manufactured Housing Prepayment Curve`, true
		case "SMM":
			return `Single Monthly Mortality`, true
		case "MPR":
			return `Monthly Prepayment Rate`, true
		case "PSA":
			return `% of BMA Prepayment Curve`, true
		case "PPC":
			return `% of Prospectus Prepayment Curve`, true
		case "CPP":
			return `Constant Prepayment Penalty`, true
		case "LOTVAR":
			return `Lot Variance (value in percent maximum over- or under-allocation allowed)`, true
		case "CPY":
			return `Constant Prepayment Yield`, true
		case "WAC":
			return `Weighted Average Coupon (value in percent)`, true
		case "ISSUE":
			return `Year of Issue`, true
		case "MAT":
			return `Maturity Year`, true
		case "PIECES":
			return `Number of Pieces`, true
		case "PMAX":
			return `Pools Maximum`, true
		case "PPM":
			return `Pools per Million`, true
		case "PPL":
			return `Pools per Lot`, true
		case "PPT":
			return `Pools per Trade`, true
		case "PROD":
			return `Production Year`, true
		case "TRDVAR":
			return `Trade Variance (value in percent maximum over- or under-allocation allowed)`, true
		case "GEOG":
			return `Geographics`, true
		}

	case 235:
		switch value {
		case "TRUE":
			return `True Yield The yield calculated with coupon dates moved from a weekend or holiday to the next valid settlement date.`, true
		case "PREVCLOSE":
			return `Previous Close Yield The yield of a bond based on the closing price 1 day ago.`, true
		case "LONGEST":
			return `Yield to Longest Average (Sinking Fund Bonds) The yield assuming only mandatory sinks are taken. This results in a slower paydown of debt; the yield is then calculated to the final payment date.`, true
		case "LONGAVGLIFE":
			return `Yield to Longest Average Life The yield assuming only mandatory sinks are taken. This results in a lower paydown of debt; the yield is then calculated to the final payment date.`, true
		case "MATURITY":
			return `Yield to Maturity The yield of a bond to its maturity date.`, true
		case "MARK":
			return `Mark To Market Yield An adjustment in the valuation of a securities portfolio to reflect the current market values of the respective securities in the portfolio.`, true
		case "OPENAVG":
			return `Open Average Yield The average yield of the respective securities in the portfolio.`, true
		case "PUT":
			return `Yield to Next Put The yield to the date at which the bond holder can next put the bond to the issuer.`, true
		case "PROCEEDS":
			return `Proceeds Yield The CD equivalent yield when the remaining time to maturity is less than two years.`, true
		case "SEMIANNUAL":
			return `Semi-annual Yield The yield of a bond whose coupon payments are reinvested semi-annually`, true
		case "SHORTAVGLIFE":
			return `Yield to Shortest Average Life same as AVGLIFE above.`, true
		case "SHORTEST":
			return `Yield to Shortest Average (Sinking Fund Bonds) The yield assuming that all sinks (mandatory and voluntary) are taken. This results in a faster paydown of debt; the yield is then calculated to the final payment date.`, true
		case "SIMPLE":
			return `Simple Yield The yield of a bond assuming no reinvestment of coupon payments. (Act/360 day count)`, true
		case "TENDER":
			return `Yield to Tender Date The yield on a Municipal bond to its mandatory tender date.`, true
		case "VALUE1/32":
			return `Yield Value Of 1/32 The amount that the yield will change for a 1/32nd change in price.`, true
		case "WORST":
			return `Yield To Worst Convention The lowest yield to all possible redemption date scenarios.`, true
		case "TAXEQUIV":
			return `Tax Equivalent Yield The after tax yield grossed up by the maximum federal tax rate of 39.6%. For comparison to taxable yields.`, true
		case "ANNUAL":
			return `Annual Yield The annual interest or dividend income an investment earns, expressed as a percentage of the investments total value.`, true
		case "LASTYEAR":
			return `Closing Yield Most Recent Year The yield of a bond based on the closing price as of the most recent years end.`, true
		case "NEXTREFUND":
			return `Yield To Next Refund (Sinking Fund Bonds) Yield assuming all bonds are redeemed at the next refund date at the redemption price.`, true
		case "AFTERTAX":
			return `After Tax Yield (Municipals) The yield on the bond net of any tax consequences from holding the bond. The discount on municipal securities can be subject to both capital gains taxes and ordinary income taxes. Calculated from dollar price.`, true
		case "ATISSUE":
			return `Yield At Issue (Municipals) The yield of the bond offered on the issue date.`, true
		case "AVGLIFE":
			return `Yield To Average Life The yield assuming that all sinks (mandatory and voluntary) are taken at par. This results in a faster paydown of debt; the yield is then calculated to the average life date.`, true
		case "AVGMATURITY":
			return `Yield To Average Maturity The yield achieved by substituting a bond's average maturity for the issue's final maturity date.`, true
		case "BOOK":
			return `Book Yield The yield of a security calculated by using its book value instead of the current market price. This term is typically used in the US domestic market.`, true
		case "CALL":
			return `Yield to Next Call The yield of a bond to the next possible call date.`, true
		case "CHANGE":
			return `Yield Change Since Close The change in the yield since the previous day's closing yield.`, true
		case "COMPOUND":
			return `Compound Yield The yield of certain Japanese bonds based on its price. Certain Japanese bonds have irregular first or last coupons, and the yield is calculated compound for these irregular periods.`, true
		case "CURRENT":
			return `Current Yield Annual interest on a bond divided by the market value. The actual income rate of return as opposed to the coupon rate expressed as a percentage.`, true
		case "GROSS":
			return `True Gross Yield Yield calculated using the price including accrued interest, where coupon dates are moved from holidays and weekends to the next trading day.`, true
		case "GOVTEQUIV":
			return `Government Equivalent Yield Ask yield based on semi-annual coupons compounding in all periods and actual/actual calendar.`, true
		case "INFLATION":
			return `Yield with Inflation Assumption Based on price, the return an investor would require on a normal bond that would make the real return equal to that of the inflation-indexed bond, assuming a constant inflation rate.`, true
		case "INVERSEFLOATER":
			return `Inverse Floater Bond Yield Inverse floater semi-annual bond equivalent rate.`, true
		case "LASTQUARTER":
			return `Closing Yield Most Recent Quarter The yield of a bond based on the closing price as of the most recent quarters end.`, true
		case "LASTCLOSE":
			return `Most Recent Closing Yield The last available yield stored in history, computed using price.`, true
		case "LASTMONTH":
			return `Closing Yield Most Recent Month The yield of a bond based on the closing price as of the most recent month's end.`, true
		case "CLOSE":
			return `Closing Yield The yield of a bond based on the closing price.`, true
		}

	case 258:
		switch value {
		case "N":
			return `Not Traded Flat`, true
		case "Y":
			return `Traded Flat`, true
		}

	case 263:
		switch value {
		case "1":
			return `Snapshot + Updates (Subscribe)`, true
		case "2":
			return `Disable previous Snapshot + Update Request (Unsubscribe)`, true
		case "0":
			return `Snapshot`, true
		}

	case 265:
		switch value {
		case "0":
			return `Full Refresh`, true
		case "1":
			return `Incremental Refresh`, true
		}

	case 266:
		switch value {
		case "Y":
			return `one book entry per side per price`, true
		case "N":
			return `Multiple entries per side per price allowed`, true
		}

	case 269:
		switch value {
		case "7":
			return `Trading Session High Price`, true
		case "1":
			return `Offer`, true
		case "A":
			return `Imbalance`, true
		case "9":
			return `Trading Session VWAP Price`, true
		case "8":
			return `Trading Session Low Price`, true
		case "5":
			return `Closing Price`, true
		case "4":
			return `Opening Price`, true
		case "0":
			return `Bid`, true
		case "2":
			return `Trade`, true
		case "3":
			return `Index Value`, true
		case "6":
			return `Settlement Price`, true
		}

	case 274:
		switch value {
		case "0":
			return `Plus Tick`, true
		case "1":
			return `Zero-Plus Tick`, true
		case "2":
			return `Minus Tick`, true
		case "3":
			return `Zero-Minus Tick`, true
		}

	case 276:
		switch value {
		case "E":
			return `Locked`, true
		case "I":
			return `Non-Firm`, true
		case "H":
			return `Fast Trading`, true
		case "F":
			return `Crossed`, true
		case "D":
			return `Consolidated Best`, true
		case "C":
			return `Exchange Best`, true
		case "B":
			return `Closed / Inactive`, true
		case "A":
			return `Open / Active`, true
		case "G":
			return `Depth`, true
		}

	case 277:
		switch value {
		case "J":
			return `Next Day Trade (next day clearing)`, true
		case "K":
			return `Opened (late report of opened trade)`, true
		case "L":
			return `Seller`, true
		case "B":
			return `Average Price Trade`, true
		case "M":
			return `Sold (out of sequence)`, true
		case "H":
			return `Rule 155 Trade (Amex)`, true
		case "N":
			return `Stopped Stock (guarantee of price but does not execute the order)`, true
		case "P":
			return `Imbalance More Buyers (Cannot be used in combination with Q)`, true
		case "Q":
			return `Imbalance More Sellers (Cannot be used in combination with P)`, true
		case "R":
			return `Opening Price`, true
		case "I":
			return `Sold Last (late reporting)`, true
		case "A":
			return `Cash (only) Market`, true
		case "C":
			return `Cash Trade (same day clearing)`, true
		case "E":
			return `Opening / Reopening Trade Detail`, true
		case "F":
			return `Intraday Trade Detail`, true
		case "G":
			return `Rule 127 Trade (NYSE)`, true
		case "D":
			return `Next Day (only) Market`, true
		}

	case 279:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Change`, true
		case "2":
			return `Delete`, true
		}

	case 281:
		switch value {
		case "7":
			return `Unsupported AggregatedBook`, true
		case "1":
			return `Duplicate MDReqID`, true
		case "C":
			return `Unsupported MDImplicitDelete`, true
		case "B":
			return `Unsupported OpenCloseSettleFlag`, true
		case "A":
			return `Unsupported Scope`, true
		case "9":
			return `Unsupported TradingSessionID`, true
		case "8":
			return `Unsupported MDEntryType`, true
		case "6":
			return `Unsupported MDUpdateType`, true
		case "5":
			return `Unsupported MarketDepth`, true
		case "4":
			return `Unsupported SubscriptionRequestType`, true
		case "2":
			return `Insufficient Bandwidth`, true
		case "0":
			return `Unknown symbol`, true
		case "3":
			return `Insufficient Permissions`, true
		}

	case 285:
		switch value {
		case "0":
			return `Cancelation / Trade Bust`, true
		case "1":
			return `Error`, true
		}

	case 286:
		switch value {
		case "1":
			return `Session Open / Close / Settlement price`, true
		case "2":
			return `Delivery Settlement price`, true
		case "3":
			return `Expected price`, true
		case "4":
			return `Price from previous business day`, true
		case "0":
			return `Daily Open / Close / Settlement price`, true
		}

	case 291:
		switch value {
		case "1":
			return `Bankrupt`, true
		case "2":
			return `Pending delisting`, true
		}

	case 292:
		switch value {
		case "B":
			return `Ex-Distribution`, true
		case "E":
			return `Ex-Interest`, true
		case "C":
			return `Ex-Rights`, true
		case "A":
			return `Ex-Dividend`, true
		case "D":
			return `New`, true
		}

	case 297:
		switch value {
		case "6":
			return `Removed from Market`, true
		case "1":
			return `Canceled for Symbol(s)`, true
		case "10":
			return `Pending`, true
		case "9":
			return `Quote Not Found`, true
		case "8":
			return `Query`, true
		case "7":
			return `Expired`, true
		case "5":
			return `Rejected`, true
		case "4":
			return `Canceled All`, true
		case "3":
			return `Canceled for Underlying`, true
		case "2":
			return `Canceled for Security Type(s)`, true
		case "0":
			return `Accepted`, true
		}

	case 298:
		switch value {
		case "4":
			return `Cancel All Quotes`, true
		case "2":
			return `Cancel for Security Type(s)`, true
		case "1":
			return `Cancel for Symbol(s)`, true
		case "3":
			return `Cancel for Underlying Symbol`, true
		}

	case 300:
		switch value {
		case "9":
			return `Not authorized to quote security`, true
		case "1":
			return `Unknown symbol (Security)`, true
		case "2":
			return `Exchange (Security) closed`, true
		case "3":
			return `Quote Request exceeds limit`, true
		case "4":
			return `Too late to enter`, true
		case "5":
			return `Unknown Quote`, true
		case "6":
			return `Duplicate Quote`, true
		case "7":
			return `Invalid bid/ask spread`, true
		case "8":
			return `Invalid price`, true
		}

	case 301:
		switch value {
		case "1":
			return `Acknowledge only negative or erroneous quotes`, true
		case "0":
			return `No Acknowledgement (Default)`, true
		case "2":
			return `Acknowledge each quote messages`, true
		}

	case 303:
		switch value {
		case "2":
			return `Automatic`, true
		case "1":
			return `Manual`, true
		}

	case 321:
		switch value {
		case "0":
			return `Request Security identity and specifications`, true
		case "1":
			return `Request Security identity for the specifications provided (Name of the security is not supplied)`, true
		case "2":
			return `Request List Security Types`, true
		case "3":
			return `Request List Securities (Can be qualified with Symbol, SecurityType, TradingSessionID, SecurityExchange. If provided then only list Securities for the specific type)`, true
		}

	case 323:
		switch value {
		case "5":
			return `Reject security proposal`, true
		case "1":
			return `Accept security proposal as is`, true
		case "6":
			return `Can not match selection criteria`, true
		case "2":
			return `Accept security proposal with revisions as indicated in the message`, true
		case "4":
			return `List of securities returned per request`, true
		case "3":
			return `List of security types returned per request`, true
		}

	case 325:
		switch value {
		case "Y":
			return `Message is being sent unsolicited`, true
		case "N":
			return `Message is being sent as a result of a prior request`, true
		}

	case 326:
		switch value {
		case "20":
			return `Unknown or Invalid`, true
		case "13":
			return `No Market On Close Imbalance`, true
		case "14":
			return `ITS Pre-Opening`, true
		case "15":
			return `New Price Indication`, true
		case "16":
			return `Trade Dissemination Time`, true
		case "17":
			return `Ready to trade (start of session)`, true
		case "19":
			return `Not Traded on this Market`, true
		case "22":
			return `Opening Rotation`, true
		case "21":
			return `Pre-Open`, true
		case "12":
			return `No Market Imbalance`, true
		case "18":
			return `Not Available for trading (end of session)`, true
		case "10":
			return `Market On Close Imbalance Sell`, true
		case "9":
			return `Market On Close Imbalance Buy`, true
		case "8":
			return `Market Imbalance Sell`, true
		case "7":
			return `Market Imbalance Buy`, true
		case "6":
			return `Trading Range Indication`, true
		case "5":
			return `Price Indication`, true
		case "4":
			return `No Open/No Resume`, true
		case "3":
			return `Resume`, true
		case "1":
			return `Opening Delay`, true
		case "2":
			return `Trading Halt`, true
		case "23":
			return `Fast Market`, true
		}

	case 327:
		switch value {
		case "X":
			return `Equipment Changeover`, true
		case "M":
			return `Additional Information`, true
		case "E":
			return `Order Influx`, true
		case "P":
			return `News Pending`, true
		case "I":
			return `Order Imbalance`, true
		case "D":
			return `News Dissemination`, true
		}

	case 328:
		switch value {
		case "Y":
			return `Halt was due to common stock being halted`, true
		case "N":
			return `Halt was not related to a halt of the common stock`, true
		}

	case 329:
		switch value {
		case "Y":
			return `Halt was due to related security being halted`, true
		case "N":
			return `Halt was not related to a halt of the related security`, true
		}

	case 334:
		switch value {
		case "1":
			return `Cancel`, true
		case "2":
			return `Error`, true
		case "3":
			return `Correction`, true
		}

	case 338:
		switch value {
		case "3":
			return `Two Party`, true
		case "1":
			return `Electronic`, true
		case "2":
			return `Open Outcry`, true
		}

	case 339:
		switch value {
		case "3":
			return `Production`, true
		case "1":
			return `Testing`, true
		case "2":
			return `Simulated`, true
		}

	case 340:
		switch value {
		case "5":
			return `Pre-Close`, true
		case "6":
			return `Request Rejected`, true
		case "4":
			return `Pre-Open`, true
		case "3":
			return `Closed`, true
		case "2":
			return `Open`, true
		case "1":
			return `Halted`, true
		case "0":
			return `Unknown`, true
		}

	case 347:
		switch value {
		case "UTF-8":
			return `(for using Unicode)`, true
		case "ISO-2022-JP":
			return `(for using JIS)`, true
		case "EUC-JP":
			return `(for using EUC)`, true
		case "Shift_JIS":
			return `(for using SJIS)`, true
		}

	case 373:
		switch value {
		case "12":
			return `XML Validation error`, true
		case "17":
			return `Non data value includes field delimiter (SOH character)`, true
		case "16":
			return `Incorrect NumInGroup count for repeating group`, true
		case "15":
			return `Repeating group fields out of order`, true
		case "14":
			return `Tag specified out of required order`, true
		case "11":
			return `Invalid MsgType`, true
		case "0":
			return `Invalid tag number`, true
		case "9":
			return `CompID problem`, true
		case "8":
			return `Signature problem`, true
		case "7":
			return `Decryption problem`, true
		case "6":
			return `Incorrect data format for value`, true
		case "5":
			return `Value is incorrect (out of range) for this tag`, true
		case "4":
			return `Tag specified without a value`, true
		case "3":
			return `Undefined Tag`, true
		case "10":
			return `SendingTime accuracy problem`, true
		case "13":
			return `Tag appears more than once`, true
		case "2":
			return `Tag not defined for this message type`, true
		case "1":
			return `Required tag missing`, true
		}

	case 374:
		switch value {
		case "N":
			return `New`, true
		case "C":
			return `Cancel`, true
		}

	case 377:
		switch value {
		case "N":
			return `Was not solicited`, true
		case "Y":
			return `Was solcitied`, true
		}

	case 378:
		switch value {
		case "7":
			return `Cancel on System Failure`, true
		case "0":
			return `GT Corporate action`, true
		case "8":
			return `Market (Exchange) Option`, true
		case "6":
			return `Cancel on Trading Halt`, true
		case "5":
			return `Partial decline of OrderQty (e.g. exchange-initiated partial cancel)`, true
		case "4":
			return `Broker option`, true
		case "3":
			return `Repricing of order`, true
		case "1":
			return `GT renewal / restatement (no corporate action)`, true
		case "2":
			return `Verbal change`, true
		}

	case 380:
		switch value {
		case "3":
			return `Unsupported Message Type`, true
		case "7":
			return `DeliverTo firm not available at this time`, true
		case "4":
			return `Application not available`, true
		case "6":
			return `Not authorized`, true
		case "0":
			return `Other`, true
		case "5":
			return `Conditionally Required Field Missing`, true
		case "1":
			return `Unkown ID`, true
		case "2":
			return `Unknown Security`, true
		}

	case 385:
		switch value {
		case "S":
			return `Send`, true
		case "R":
			return `Receive`, true
		}

	case 388:
		switch value {
		case "0":
			return `Related to displayed price`, true
		case "1":
			return `Related to market price`, true
		case "2":
			return `Related to primary price`, true
		case "3":
			return `Related to local primary price`, true
		case "4":
			return `Related to midpoint price`, true
		case "5":
			return `Related to last trade price`, true
		}

	case 394:
		switch value {
		case "1":
			return `Non Disclosed Style (e.g. US/European)`, true
		case "2":
			return `Disclosed Style (e.g. Japanese)`, true
		case "3":
			return `No Bidding Process`, true
		}

	case 399:
		switch value {
		case "3":
			return `Index`, true
		case "2":
			return `Country`, true
		case "1":
			return `Sector`, true
		}

	case 401:
		switch value {
		case "1":
			return `SideValue1`, true
		case "2":
			return `SideValue 2`, true
		}

	case 409:
		switch value {
		case "3":
			return `Normal Market Size`, true
		case "4":
			return `Other`, true
		case "2":
			return `20 day moving average`, true
		case "1":
			return `5day moving average`, true
		}

	case 411:
		switch value {
		case "N":
			return `False`, true
		case "Y":
			return `True`, true
		}

	case 414:
		switch value {
		case "3":
			return `Real-time execution reports (to be discouraged)`, true
		case "2":
			return `SellSide periodically sends status using ListStatus. Period optionally specified in ProgressPeriod`, true
		case "1":
			return `BuySide explicitly requests status using StatusRequest (Default) The sell-side firm can however, send a DONE status List Status Response in an unsolicited fashion`, true
		}

	case 416:
		switch value {
		case "2":
			return `Gross`, true
		case "1":
			return `Net`, true
		}

	case 418:
		switch value {
		case "G":
			return `VWAP Guarantee`, true
		case "A":
			return `Agency`, true
		case "J":
			return `Guaranteed Close`, true
		case "R":
			return `Risk Trade`, true
		}

	case 419:
		switch value {
		case "8":
			return `VWAP through an afternoon session`, true
		case "D":
			return `Open`, true
		case "Z":
			return `Others`, true
		case "C":
			return `Strike`, true
		case "B":
			return `VWAP through an afternoon session except YORI (an opening auction)`, true
		case "9":
			return `VWAP through a day except YORI (an opening auction)`, true
		case "7":
			return `VWAP through a morning session`, true
		case "6":
			return `VWAP through a day`, true
		case "5":
			return `SQ`, true
		case "4":
			return `Current price`, true
		case "3":
			return `Closing Price`, true
		case "2":
			return `Closing Price at morning session`, true
		case "A":
			return `VWAP through a morning session except YORI (an opening auction)`, true
		}

	case 423:
		switch value {
		case "3":
			return `Fixed Amount (absolute value)`, true
		case "1":
			return `Percentage`, true
		case "4":
			return `discount - percentage points below par`, true
		case "6":
			return `basis points relative to benchmark`, true
		case "7":
			return `TED price (see "Volume 1 - Glossary")`, true
		case "8":
			return `TED yield (see "Volume 1 - Glossary")`, true
		case "5":
			return `premium - percentage points over par`, true
		case "2":
			return `per share (e.g. cents per share)`, true
		}

	case 427:
		switch value {
		case "0":
			return `book out all trades on day of execution`, true
		case "2":
			return `accumulate until verbally notified otherwise`, true
		case "1":
			return `accumulate executions until order is filled or expires`, true
		}

	case 429:
		switch value {
		case "6":
			return `Alert`, true
		case "4":
			return `ExecStarted`, true
		case "3":
			return `Timed`, true
		case "2":
			return `Response`, true
		case "1":
			return `Ack`, true
		case "5":
			return `AllDone`, true
		}

	case 430:
		switch value {
		case "1":
			return `Net`, true
		case "2":
			return `Gross`, true
		}

	case 431:
		switch value {
		case "4":
			return `Canceling`, true
		case "3":
			return `Executing`, true
		case "7":
			return `Reject`, true
		case "6":
			return `All Done`, true
		case "5":
			return `Alert`, true
		case "2":
			return `ReceivedForExecution`, true
		case "1":
			return `InBiddingProcess`, true
		}

	case 433:
		switch value {
		case "5":
			return `Exchange/switch CIV order Buy driven, cash withdraw (i.e. additional cash will not be provided to fulfil the order)`, true
		case "4":
			return `Exchange/switch CIV order Buy driven, cash top-up (i.e. additional cash will be provided to fulfil the order)`, true
		case "2":
			return `Wait for Execute Instruction (e.g. a List Execute message or phone call before proceeding with execution of the list)`, true
		case "1":
			return `Immediate`, true
		case "3":
			return `Exchange/switch CIV order Sell driven`, true
		}

	case 434:
		switch value {
		case "2":
			return `Order Cancel/Replace Request`, true
		case "1":
			return `Order Cancel Request`, true
		}

	case 442:
		switch value {
		case "1":
			return `Single Security (default if not specified)`, true
		case "2":
			return `Individual leg of a multi-leg security`, true
		case "3":
			return `Multi-leg security`, true
		}

	case 447:
		switch value {
		case "5":
			return `Chinese B Share (Shezhen and Shanghai)`, true
		case "8":
			return `US Employer Identification Number`, true
		case "A":
			return `Australian Tax File Number`, true
		case "9":
			return `Australian Business Number`, true
		case "E":
			return `ISO Country Code`, true
		case "B":
			return `BIC (Bank Identification CodeSwift managed) code (ISO 9362 - See Appendix 6-B)`, true
		case "7":
			return `US Social Security Number`, true
		case "D":
			return `Proprietary/Custom code`, true
		case "F":
			return `Settlement Entity Location (note if Local Market Settlement use E = ISO Country Code) (see Appendix 6-G for valid values)`, true
		case "1":
			return `Korean Investor ID`, true
		case "2":
			return `Taiwanese Qualified Foreign Investor ID QFII / FID`, true
		case "3":
			return `Taiwanese Trading Account`, true
		case "4":
			return `Malaysian Central Depository (MCD) number`, true
		case "6":
			return `UK National Insurance or Pension Number`, true
		case "C":
			return `Generally accepted market participant identifier (e.g. NASD mnemonic)`, true
		}

	case 452:
		switch value {
		case "15":
			return `Correspondant Clearing Firm`, true
		case "3":
			return `Client ID (formerly FIX 4.2 ClientID)`, true
		case "20":
			return `Underlying Contra Firm`, true
		case "19":
			return `Sponsoring Firm`, true
		case "18":
			return `Contra Clearing Firm`, true
		case "17":
			return `Contra Firm`, true
		case "16":
			return `Executing System`, true
		case "7":
			return `Entering Firm`, true
		case "1":
			return `Executing Firm (formerly FIX 4.2 ExecBroker)`, true
		case "2":
			return `Broker of Credit (formerly FIX 4.2 BrokerOfCredit)`, true
		case "5":
			return `Investor ID`, true
		case "6":
			return `Introducing Firm`, true
		case "14":
			return `Giveup Clearing Firm (firm to which trade is given up)`, true
		case "8":
			return `Locate/Lending Firm (for short-sales)`, true
		case "9":
			return `Fund manager Client ID (for CIV)`, true
		case "10":
			return `Settlement Location (formerly FIX 4.2 SettlLocation)`, true
		case "11":
			return `Order Origination Trader (associated with Order Origination Firm e.g. trader who initiates/submits the order)`, true
		case "12":
			return `Executing Trader (associated with Executing Firm - actually executes)`, true
		case "13":
			return `Order Origination Firm (e.g. buyside firm)`, true
		case "4":
			return `Clearing Firm (formerly FIX 4.2 ClearingFirm)`, true
		}

	case 460:
		switch value {
		case "8":
			return `LOAN`, true
		case "12":
			return `OTHER`, true
		case "11":
			return `MUNICIPAL`, true
		case "1":
			return `AGENCY`, true
		case "3":
			return `CORPORATE`, true
		case "4":
			return `CURRENCY`, true
		case "2":
			return `COMMODITY`, true
		case "6":
			return `GOVERNMENT`, true
		case "10":
			return `MORTGAGE`, true
		case "7":
			return `INDEX`, true
		case "9":
			return `MONEYMARKET`, true
		case "5":
			return `EQUITY`, true
		}

	case 464:
		switch value {
		case "Y":
			return `True (Test)`, true
		case "N":
			return `False (Production)`, true
		}

	case 465:
		switch value {
		case "6":
			return `CONTRACTS`, true
		case "7":
			return `OTHER`, true
		case "5":
			return `CURRENCY`, true
		case "4":
			return `ORIGINALFACE`, true
		case "3":
			return `CURRENTFACE`, true
		case "2":
			return `BONDS`, true
		case "1":
			return `SHARES`, true
		case "8":
			return `PAR (see Volume 1 Glossary)`, true
		}

	case 468:
		switch value {
		case "0":
			return `Round to nearest`, true
		case "1":
			return `Round down`, true
		case "2":
			return `Round up`, true
		}

	case 480:
		switch value {
		case "M":
			return `No waiver agreement`, true
		case "N":
			return `No execution only`, true
		case "Y":
			return `Yes`, true
		case "O":
			return `No institutional.`, true
		}

	case 481:
		switch value {
		case "3":
			return `Exempt Authorised Credit or Financial Institution.`, true
		case "2":
			return `Exempt Client Money Type Exemption`, true
		case "1":
			return `Exempt Below The Limit`, true
		case "Y":
			return `Passed`, true
		case "N":
			return `Not checked`, true
		}

	case 484:
		switch value {
		case "S":
			return `Single price`, true
		case "Q":
			return `Offer price minus adjustment amount`, true
		case "P":
			return `Offer price minus adjustment %`, true
		case "O":
			return `Offer price`, true
		case "E":
			return `Creation price plus adjustment amount`, true
		case "D":
			return `Creation price plus adjustment %`, true
		case "C":
			return `Creation price`, true
		case "B":
			return `Bid price`, true
		}

	case 487:
		switch value {
		case "N":
			return `New`, true
		case "R":
			return `Replace`, true
		case "C":
			return `Cancel`, true
		}

	case 492:
		switch value {
		case "14":
			return `BPAY`, true
		case "13":
			return `ACH Credit`, true
		case "12":
			return `ACH Debit`, true
		case "11":
			return `Credit Card`, true
		case "10":
			return `Direct Credit (BECS)`, true
		case "9":
			return `Direct Debit (BECS)`, true
		case "8":
			return `Debit Card`, true
		case "7":
			return `FedWire`, true
		case "15":
			return `High Value Clearing System (HVACS)`, true
		case "3":
			return `Euroclear`, true
		case "6":
			return `Telegraphic Transfer`, true
		case "4":
			return `Clearstream`, true
		case "1":
			return `CREST`, true
		case "2":
			return `NSCC`, true
		case "5":
			return `Cheque`, true
		}

	case 495:
		switch value {
		case "19":
			return `Profit Sharing Plan (US)`, true
		case "11":
			return `Employer - prior year (US)`, true
		case "12":
			return `Employer – current year (US)`, true
		case "13":
			return `Non-fund prototype IRA (US)`, true
		case "14":
			return `Non-fund qualified plan (US)`, true
		case "15":
			return `Defined contribution plan (US)`, true
		case "10":
			return `Employee – current year (US)`, true
		case "17":
			return `Individual Retirement Account – Rollover (US)`, true
		case "5":
			return `Mini Insurance ISA (UK)`, true
		case "16":
			return `Individual Retirement Account (US)`, true
		case "9":
			return `Employee - prior year (US)`, true
		case "8":
			return `Asset transfer (US)`, true
		case "21":
			return `Self-Directed IRA (US)`, true
		case "6":
			return `Current year payment (US)`, true
		case "20":
			return `401K (US)`, true
		case "4":
			return `Mini Stocks and Shares ISA (UK)`, true
		case "3":
			return `Mini Cash ISA (UK)`, true
		case "2":
			return `TESSA (UK)`, true
		case "1":
			return `Maxi ISA (UK)`, true
		case "0":
			return `None/Not Applicable (default)`, true
		case "7":
			return `Prior year payment (US)`, true
		case "23":
			return `457 (US)`, true
		case "24":
			return `Roth IRA (fund prototype) (US)`, true
		case "25":
			return `Roth IRA (non-prototype) (US)`, true
		case "26":
			return `Roth Conversion IRA (fund prototype) (US)`, true
		case "27":
			return `Roth Conversion IRA (non-prototype) (US)`, true
		case "28":
			return `Education IRA (fund prototype) (US)`, true
		case "29":
			return `Education IRA (non-prototype) (US)`, true
		case "18":
			return `KEOGH (US)`, true
		case "22":
			return `403(b) (US)`, true
		}

	case 497:
		switch value {
		case "N":
			return `No`, true
		case "Y":
			return `Yes`, true
		}

	case 506:
		switch value {
		case "A":
			return `Accept`, true
		case "N":
			return `Reminder`, true
		case "R":
			return `Reject`, true
		case "H":
			return `Held`, true
		}

	case 507:
		switch value {
		case "13":
			return `Invalid/unacceptable NoDistribInstns`, true
		case "17":
			return `Invalid/unacceptable Cash Distrib Agent Code`, true
		case "16":
			return `Invalid/unacceptable Cash Distrib Agent Acct Name`, true
		case "4":
			return `Invalid/unacceptable No Reg Detls`, true
		case "15":
			return `Invalid/unacceptable Distrib Payment Method`, true
		case "14":
			return `Invalid/unacceptable Distrib Percentage`, true
		case "3":
			return `Invalid/unacceptable Ownership Type`, true
		case "2":
			return `Invalid/unacceptable Tax Exempt Type`, true
		case "12":
			return `Invalid/unacceptable Investor Country Of Residence`, true
		case "11":
			return `Invalid/unacceptable Date of Birth`, true
		case "10":
			return `Invalid/unacceptable Investor ID Source`, true
		case "9":
			return `Invalid/unacceptable Investor ID`, true
		case "8":
			return `Invalid/unacceptable Mailing Inst`, true
		case "7":
			return `Invalid/unacceptable Mailing Dtls`, true
		case "5":
			return `Invalid/unacceptable Reg Seq No`, true
		case "1":
			return `Invalid/unacceptable Account Type`, true
		case "18":
			return `Invalid/unacceptable Cash Distrib Agent Acct Num`, true
		case "6":
			return `Invalid/unacceptable Reg Dtls`, true
		}

	case 514:
		switch value {
		case "2":
			return `Cancel`, true
		case "0":
			return `New`, true
		case "1":
			return `Replace`, true
		}

	case 519:
		switch value {
		case "15":
			return `Net Settlement Amount`, true
		case "1":
			return `Commission Amount (actual)`, true
		case "2":
			return `Commission % (actual)`, true
		case "3":
			return `Initial Charge Amount`, true
		case "4":
			return `Initial Charge %`, true
		case "5":
			return `Discount Amount`, true
		case "6":
			return `Discount %`, true
		case "7":
			return `Dilution Levy Amount`, true
		case "8":
			return `Dilution Levy %`, true
		case "9":
			return `Exit Charge Amount`, true
		case "10":
			return `Exit Charge %`, true
		case "11":
			return `Fund-based Renewal Commission % (a.k.a. Trail commission)`, true
		case "12":
			return `Projected Fund Value (i.e. for investments intended to realise or exceed a specific future value)`, true
		case "14":
			return `Fund-based Renewal Commission Amount (based on Projected Fund value)`, true
		case "13":
			return `Fund-based Renewal Commission Amount (based on Order value)`, true
		}

	case 522:
		switch value {
		case "5":
			return `Company Trustee`, true
		case "13":
			return `Nominee`, true
		case "12":
			return `Corporate Body`, true
		case "11":
			return `Non-Profit Organization`, true
		case "10":
			return `Networking Sub-Account`, true
		case "9":
			return `Fiduciaries`, true
		case "8":
			return `Trusts`, true
		case "6":
			return `Pension Plan`, true
		case "4":
			return `Individual Trustee`, true
		case "2":
			return `Public Company`, true
		case "3":
			return `Private Company`, true
		case "1":
			return `Individual Investor`, true
		case "7":
			return `Custodian Under Gifts to Minors Act`, true
		}

	case 528:
		switch value {
		case "R":
			return `Riskless Principal`, true
		case "I":
			return `Individual`, true
		case "P":
			return `Principal (Note for CMS purposes, Principal includes Proprietary)`, true
		case "W":
			return `Agent for Other Member`, true
		case "A":
			return `Agency`, true
		case "G":
			return `Proprietary`, true
		}

	case 529:
		switch value {
		case "7":
			return `Foreign Entity (of foreign governmnet or regulatory jurisdiction)`, true
		case "A":
			return `Riskless Arbitrage`, true
		case "1":
			return `Program Trade`, true
		case "8":
			return `External Market Participant`, true
		case "6":
			return `Acting as Market Maker or Specialist in the underlying security of a derivative security`, true
		case "5":
			return `Acting as Market Maker or Specialist in the security`, true
		case "3":
			return `Non-Index Arbitrage`, true
		case "2":
			return `Index Arbitrage`, true
		case "4":
			return `Competing Market Maker`, true
		case "9":
			return `External Inter-connected Market Linkage`, true
		}

	case 530:
		switch value {
		case "1":
			return `Cancel orders for a security`, true
		case "7":
			return `Cancel all orders`, true
		case "6":
			return `Cancel orders for a trading session`, true
		case "5":
			return `Cancel orders for a SecurityType`, true
		case "4":
			return `Cancel orders for a CFICode`, true
		case "2":
			return `Cancel orders for an Underlying security`, true
		case "3":
			return `Cancel orders for a Product`, true
		}

	case 531:
		switch value {
		case "6":
			return `Cancel orders for a trading session`, true
		case "0":
			return `Cancel Request Rejected -- See MassCancelRejectReason (532)`, true
		case "7":
			return `Cancel all orders`, true
		case "3":
			return `Cancel orders for a Product`, true
		case "5":
			return `Cancel orders for a SecurityType`, true
		case "4":
			return `Cancel orders for a CFICode`, true
		case "1":
			return `Cancel orders for a security`, true
		case "2":
			return `Cancel orders for an Underlying security`, true
		}

	case 532:
		switch value {
		case "2":
			return `Invalid or unknown underlying`, true
		case "6":
			return `Invalid or unknown trading session`, true
		case "5":
			return `Invalid or unknown Security Type`, true
		case "3":
			return `Invalid or unknown Product`, true
		case "1":
			return `Invalid or unknown Security`, true
		case "0":
			return `Mass Cancel Not Supported`, true
		case "4":
			return `Invalid or unknown CFICode`, true
		}

	case 537:
		switch value {
		case "0":
			return `Indicative`, true
		case "1":
			return `Tradeable`, true
		case "2":
			return `Restricted Tradeable`, true
		}

	case 544:
		switch value {
		case "2":
			return `Margin Open`, true
		case "3":
			return `Margin Close`, true
		case "1":
			return `Cash`, true
		}

	case 546:
		switch value {
		case "1":
			return `Local (Exchange, ECN, ATS)`, true
		case "2":
			return `National`, true
		case "3":
			return `Global`, true
		}

	case 547:
		switch value {
		case "Y":
			return `Client has responsibility for implicitly deleting bids or offers falling outside the MarketDepth of the request.`, true
		case "N":
			return `Server must send an explicit delete for bids or offers falling outside the requested MarketDepth of the request.`, true
		}

	case 549:
		switch value {
		case "1":
			return `Cross Trade which is executed completely or not. Both sides are treated in the same manner. This is equivalent to an All or None.`, true
		case "2":
			return `Cross Trade which is executed partially and the rest is cancelled. One side is fully executed, the other side is partially executed with the remainder being cancelled. This is equivalent to an Immediate or Cancel on the other side.`, true
		case "3":
			return `Cross trade which is partially executed with the unfilled portions remaining active. One side of the cross is fully executed (as denoted with the CrossPrioritization field), but the unfilled portion remains active.`, true
		case "4":
			return `Cross trade is executed with existing orders with the same price.`, true
		}

	case 550:
		switch value {
		case "2":
			return `Sellside prioritized`, true
		case "0":
			return `None`, true
		case "1":
			return `Buyside prioritized`, true
		}

	case 552:
		switch value {
		case "1":
			return `one side`, true
		case "2":
			return `both sides`, true
		}

	case 559:
		switch value {
		case "1":
			return `SecurityType and/or CFICode`, true
		case "2":
			return `Product`, true
		case "3":
			return `TradingSessionID`, true
		case "4":
			return `All Securities`, true
		case "0":
			return `Symbol`, true
		}

	case 560:
		switch value {
		case "4":
			return `Instrument data temporarily unavailable`, true
		case "0":
			return `Valid request`, true
		case "1":
			return `Invalid or unsupported request`, true
		case "5":
			return `Request for instrument data not supported`, true
		case "3":
			return `Not authorized to retrieve instrument data`, true
		case "2":
			return `No instruments found that match selection criteria`, true
		}

	case 567:
		switch value {
		case "1":
			return `Unknown or invalid TradingSessionID`, true
		}

	case 569:
		switch value {
		case "4":
			return `Advisories that match criteria`, true
		case "3":
			return `Unreported trades that match criteria`, true
		case "2":
			return `Unmatched trades that match criteria`, true
		case "1":
			return `Matched trades matching Criteria provided on request (parties, order id, instrument, input source, etc.)`, true
		case "0":
			return `All trades`, true
		}

	case 570:
		switch value {
		case "N":
			return `not reported to counterparty`, true
		case "Y":
			return `previously reported to counterparty`, true
		}

	case 573:
		switch value {
		case "0":
			return `compared, matched or affirmed`, true
		case "1":
			return `uncompared, unmatched, or unaffirmed`, true
		case "2":
			return `advisory or alert`, true
		}

	case 574:
		switch value {
		case "S5":
			return `Summarized Match using A1 to A5 exact match criteria except quantity is summarized`, true
		case "M1":
			return `ACT M1 Match / Exact Match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator minus badges and times`, true
		case "M6":
			return `ACT M6 Match`, true
		case "M5":
			return `ACT Default After M2`, true
		case "M3":
			return `ACT Accepted Trade`, true
		case "S2":
			return `Summarized Match using A1 to A5 exact match criteria except quantity is summarized`, true
		case "S3":
			return `Summarized Match using A1 to A5 exact match criteria except quantity is summarized`, true
		case "S4":
			return `Summarized Match using A1 to A5 exact match criteria except quantity is summarized`, true
		case "M2":
			return `ACT M2 Match / Summarized Match minus badges and times`, true
		case "A2":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus four badges`, true
		case "A3":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus two badges and execution time (within two-minute window)`, true
		case "A4":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and`, true
		case "AQ":
			return `Compared records resulting from stamped advisories or specialist`, true
		case "MT":
			return `Non-ACT / OCS Locked In`, true
		case "M4":
			return `ACT Default Trade`, true
		case "A1":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus four badges and execution time (within two-minute window)`, true
		case "S1":
			return `Summarized Match using A1 to A5 exact match criteria except quantity is summarized`, true
		case "A5":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus execution time (within two-minute window)`, true
		}

	case 577:
		switch value {
		case "8":
			return `Manual mode (pre-posting and/or pre-giveup)`, true
		case "5":
			return `multilateral netting`, true
		case "9":
			return `Automatic posting mode (trade posting to the position account number specified)`, true
		case "2":
			return `bilateral netting only`, true
		case "6":
			return `clear against central counterparty`, true
		case "10":
			return `Automatic give-up mode (trade give-up to the give-up destination number specified)`, true
		case "4":
			return `special trade`, true
		case "3":
			return `ex clearing`, true
		case "0":
			return `process normally`, true
		case "7":
			return `exclude from central counterparty`, true
		case "1":
			return `exclude from all netting`, true
		}

	case 581:
		switch value {
		case "3":
			return `House Trader`, true
		case "7":
			return `Account is house trader and is cross margined`, true
		case "6":
			return `Account is carried on non-customer side of books and is cross margined`, true
		case "4":
			return `Floor Trader`, true
		case "2":
			return `Account is carried on non-Customer Side of books`, true
		case "1":
			return `Account is carried on customer Side of Books`, true
		case "8":
			return `Joint Backoffice Account (JBO)`, true
		}

	case 585:
		switch value {
		case "1":
			return `Status for orders for a security`, true
		case "2":
			return `Status for orders for an Underlying security`, true
		case "3":
			return `Status for orders for a Product`, true
		case "4":
			return `Status for orders for a CFICode`, true
		case "5":
			return `Status for orders for a SecurityType`, true
		case "6":
			return `Status for orders for a trading session`, true
		case "8":
			return `Status for orders for a PartyID`, true
		case "7":
			return `Status for all orders`, true
		}

	case 589:
		switch value {
		case "0":
			return `Can trigger booking without reference to the order initiator ("auto")`, true
		case "1":
			return `Speak with order initiator before booking ("speak first")`, true
		}

	case 590:
		switch value {
		case "1":
			return `Aggregate partial executions on this order, and book one trade per order`, true
		case "2":
			return `Aggregate executions for this symbol, side, and settlement date`, true
		case "0":
			return `Each partial execution is a bookable unit`, true
		}

	case 591:
		switch value {
		case "0":
			return `Pro-rata`, true
		case "1":
			return `Do not pro-rata = discuss first`, true
		}

	case 626:
		switch value {
		case "6":
			return `Buyside Ready-To-Book - Combined Set of Orders`, true
		case "2":
			return `Buyside Preliminary (without MiscFees and NetMoney)`, true
		case "3":
			return `Sellside Calculated Using Preliminary (includes MiscFees and NetMoney)`, true
		case "5":
			return `Buyside Ready-To-Book - Single Order`, true
		case "1":
			return `Buyside Calculated (includes MiscFees and NetMoney)`, true
		case "4":
			return `Sellside Calculated Without Preliminary (sent unsolicited by sellside, includes MiscFees and NetMoney)`, true
		}

	case 635:
		switch value {
		case "H":
			return `106.H and 106.J Firms`, true
		case "5":
			return `5th year delegate trading for his own account`, true
		case "4":
			return `4th year delegate trading for his own account`, true
		case "3":
			return `3rd year delegate trading for his own account`, true
		case "2":
			return `2nd year delegate trading for his own account`, true
		case "1":
			return `1st year delegate trading for his own account`, true
		case "M":
			return `All other ownership types`, true
		case "I":
			return `GIM, IDEM and COM Membership Interest Holders`, true
		case "9":
			return `6th year and beyond delegate trading for his own account`, true
		case "F":
			return `Full and Associate Member trading for own account and as floor`, true
		case "E":
			return `Equity Member and Clearing Member`, true
		case "C":
			return `Non-member and Customer`, true
		case "B":
			return `CBOE Member`, true
		case "L":
			return `Lessee and 106.F Employees`, true
		}

	case 636:
		switch value {
		case "N":
			return `Order has been accepted but not yet in a working state`, true
		case "Y":
			return `Order is currently being worked`, true
		}

	case 638:
		switch value {
		case "0":
			return `Priority Unchanged`, true
		case "1":
			return `Lost Priority as result of order change`, true
		}

	case 650:
		switch value {
		case "Y":
			return `Legal confirm`, true
		case "N":
			return `Does not constitute a legal confirm`, true
		}

	case 653:
		switch value {
		case "2":
			return `Rejected`, true
		case "3":
			return `Unauthorized request`, true
		case "1":
			return `Approved (Accepted)`, true
		case "0":
			return `Pending Approval`, true
		case "4":
			return `Invalid definition request`, true
		}

	case 658:
		switch value {
		case "1":
			return `Unknown symbol (Security)`, true
		case "2":
			return `Exchange (Security) closed`, true
		case "3":
			return `Quote Request exceeds limit`, true
		case "4":
			return `Too late to enter`, true
		case "5":
			return `Invalid price`, true
		case "6":
			return `Not authorized to request quote`, true
		}

	}
	return "", false
}
