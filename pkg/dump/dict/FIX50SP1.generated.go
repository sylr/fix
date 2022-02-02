package dict

type FIX50SP1 struct {
}

func (f FIX50SP1) Version() string {
	return "FIX.5.0SP1"
}

func (f FIX50SP1) TagName(tag int) (string, bool) {
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
		return "IOIID", true
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
		return "NoLinesOfText", true
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
		return "SettlType", true
	case 64:
		return "SettlDate", true
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
		return "AvgPxPrecision", true
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
	case 85:
		return "NoDlvyInst", true
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
		return "SettlDate2", true
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
		return "PegOffsetValue", true
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
		return "OpenCloseSettlFlag", true
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
		return "TotNoQuoteEntries", true
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
	case 318:
		return "UnderlyingCurrency", true
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
		return "DiscretionOffsetValue", true
	case 390:
		return "BidID", true
	case 391:
		return "ClientBidID", true
	case 392:
		return "ListName", true
	case 393:
		return "TotNoRelatedSym", true
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
		return "BidTradeType", true
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
		return "CardIssNum", true
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
	case 502:
		return "CashDistribAgentAcctName", true
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
		return "RegistDtls", true
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
		return "TotNoSecurityTypes", true
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
		return "LegSettlType", true
	case 588:
		return "LegSettlDate", true
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
	case 660:
		return "AcctIDSource", true
	case 661:
		return "AllocAcctIDSource", true
	case 662:
		return "BenchmarkPrice", true
	case 663:
		return "BenchmarkPriceType", true
	case 664:
		return "ConfirmID", true
	case 665:
		return "ConfirmStatus", true
	case 666:
		return "ConfirmTransType", true
	case 667:
		return "ContractSettlMonth", true
	case 668:
		return "DeliveryForm", true
	case 669:
		return "LastParPx", true
	case 670:
		return "NoLegAllocs", true
	case 671:
		return "LegAllocAccount", true
	case 672:
		return "LegIndividualAllocID", true
	case 673:
		return "LegAllocQty", true
	case 674:
		return "LegAllocAcctIDSource", true
	case 675:
		return "LegSettlCurrency", true
	case 676:
		return "LegBenchmarkCurveCurrency", true
	case 677:
		return "LegBenchmarkCurveName", true
	case 678:
		return "LegBenchmarkCurvePoint", true
	case 679:
		return "LegBenchmarkPrice", true
	case 680:
		return "LegBenchmarkPriceType", true
	case 681:
		return "LegBidPx", true
	case 682:
		return "LegIOIQty", true
	case 683:
		return "NoLegStipulations", true
	case 684:
		return "LegOfferPx", true
	case 685:
		return "LegOrderQty", true
	case 686:
		return "LegPriceType", true
	case 687:
		return "LegQty", true
	case 688:
		return "LegStipulationType", true
	case 689:
		return "LegStipulationValue", true
	case 690:
		return "LegSwapType", true
	case 691:
		return "Pool", true
	case 692:
		return "QuotePriceType", true
	case 693:
		return "QuoteRespID", true
	case 694:
		return "QuoteRespType", true
	case 695:
		return "QuoteQualifier", true
	case 696:
		return "YieldRedemptionDate", true
	case 697:
		return "YieldRedemptionPrice", true
	case 698:
		return "YieldRedemptionPriceType", true
	case 699:
		return "BenchmarkSecurityID", true
	case 700:
		return "ReversalIndicator", true
	case 701:
		return "YieldCalcDate", true
	case 702:
		return "NoPositions", true
	case 703:
		return "PosType", true
	case 704:
		return "LongQty", true
	case 705:
		return "ShortQty", true
	case 706:
		return "PosQtyStatus", true
	case 707:
		return "PosAmtType", true
	case 708:
		return "PosAmt", true
	case 709:
		return "PosTransType", true
	case 710:
		return "PosReqID", true
	case 711:
		return "NoUnderlyings", true
	case 712:
		return "PosMaintAction", true
	case 713:
		return "OrigPosReqRefID", true
	case 714:
		return "PosMaintRptRefID", true
	case 715:
		return "ClearingBusinessDate", true
	case 716:
		return "SettlSessID", true
	case 717:
		return "SettlSessSubID", true
	case 718:
		return "AdjustmentType", true
	case 719:
		return "ContraryInstructionIndicator", true
	case 720:
		return "PriorSpreadIndicator", true
	case 721:
		return "PosMaintRptID", true
	case 722:
		return "PosMaintStatus", true
	case 723:
		return "PosMaintResult", true
	case 724:
		return "PosReqType", true
	case 725:
		return "ResponseTransportType", true
	case 726:
		return "ResponseDestination", true
	case 727:
		return "TotalNumPosReports", true
	case 728:
		return "PosReqResult", true
	case 729:
		return "PosReqStatus", true
	case 730:
		return "SettlPrice", true
	case 731:
		return "SettlPriceType", true
	case 732:
		return "UnderlyingSettlPrice", true
	case 733:
		return "UnderlyingSettlPriceType", true
	case 734:
		return "PriorSettlPrice", true
	case 735:
		return "NoQuoteQualifiers", true
	case 736:
		return "AllocSettlCurrency", true
	case 737:
		return "AllocSettlCurrAmt", true
	case 738:
		return "InterestAtMaturity", true
	case 739:
		return "LegDatedDate", true
	case 740:
		return "LegPool", true
	case 741:
		return "AllocInterestAtMaturity", true
	case 742:
		return "AllocAccruedInterestAmt", true
	case 743:
		return "DeliveryDate", true
	case 744:
		return "AssignmentMethod", true
	case 745:
		return "AssignmentUnit", true
	case 746:
		return "OpenInterest", true
	case 747:
		return "ExerciseMethod", true
	case 748:
		return "TotNumTradeReports", true
	case 749:
		return "TradeRequestResult", true
	case 750:
		return "TradeRequestStatus", true
	case 751:
		return "TradeReportRejectReason", true
	case 752:
		return "SideMultiLegReportingType", true
	case 753:
		return "NoPosAmt", true
	case 754:
		return "AutoAcceptIndicator", true
	case 755:
		return "AllocReportID", true
	case 756:
		return "NoNested2PartyIDs", true
	case 757:
		return "Nested2PartyID", true
	case 758:
		return "Nested2PartyIDSource", true
	case 759:
		return "Nested2PartyRole", true
	case 760:
		return "Nested2PartySubID", true
	case 761:
		return "BenchmarkSecurityIDSource", true
	case 762:
		return "SecuritySubType", true
	case 763:
		return "UnderlyingSecuritySubType", true
	case 764:
		return "LegSecuritySubType", true
	case 765:
		return "AllowableOneSidednessPct", true
	case 766:
		return "AllowableOneSidednessValue", true
	case 767:
		return "AllowableOneSidednessCurr", true
	case 768:
		return "NoTrdRegTimestamps", true
	case 769:
		return "TrdRegTimestamp", true
	case 770:
		return "TrdRegTimestampType", true
	case 771:
		return "TrdRegTimestampOrigin", true
	case 772:
		return "ConfirmRefID", true
	case 773:
		return "ConfirmType", true
	case 774:
		return "ConfirmRejReason", true
	case 775:
		return "BookingType", true
	case 776:
		return "IndividualAllocRejCode", true
	case 777:
		return "SettlInstMsgID", true
	case 778:
		return "NoSettlInst", true
	case 779:
		return "LastUpdateTime", true
	case 780:
		return "AllocSettlInstType", true
	case 781:
		return "NoSettlPartyIDs", true
	case 782:
		return "SettlPartyID", true
	case 783:
		return "SettlPartyIDSource", true
	case 784:
		return "SettlPartyRole", true
	case 785:
		return "SettlPartySubID", true
	case 786:
		return "SettlPartySubIDType", true
	case 787:
		return "DlvyInstType", true
	case 788:
		return "TerminationType", true
	case 789:
		return "NextExpectedMsgSeqNum", true
	case 790:
		return "OrdStatusReqID", true
	case 791:
		return "SettlInstReqID", true
	case 792:
		return "SettlInstReqRejCode", true
	case 793:
		return "SecondaryAllocID", true
	case 794:
		return "AllocReportType", true
	case 795:
		return "AllocReportRefID", true
	case 796:
		return "AllocCancReplaceReason", true
	case 797:
		return "CopyMsgIndicator", true
	case 798:
		return "AllocAccountType", true
	case 799:
		return "OrderAvgPx", true
	case 800:
		return "OrderBookingQty", true
	case 801:
		return "NoSettlPartySubIDs", true
	case 802:
		return "NoPartySubIDs", true
	case 803:
		return "PartySubIDType", true
	case 804:
		return "NoNestedPartySubIDs", true
	case 805:
		return "NestedPartySubIDType", true
	case 806:
		return "NoNested2PartySubIDs", true
	case 807:
		return "Nested2PartySubIDType", true
	case 808:
		return "AllocIntermedReqType", true
	case 810:
		return "UnderlyingPx", true
	case 811:
		return "PriceDelta", true
	case 812:
		return "ApplQueueMax", true
	case 813:
		return "ApplQueueDepth", true
	case 814:
		return "ApplQueueResolution", true
	case 815:
		return "ApplQueueAction", true
	case 816:
		return "NoAltMDSource", true
	case 817:
		return "AltMDSourceID", true
	case 818:
		return "SecondaryTradeReportID", true
	case 819:
		return "AvgPxIndicator", true
	case 820:
		return "TradeLinkID", true
	case 821:
		return "OrderInputDevice", true
	case 822:
		return "UnderlyingTradingSessionID", true
	case 823:
		return "UnderlyingTradingSessionSubID", true
	case 824:
		return "TradeLegRefID", true
	case 825:
		return "ExchangeRule", true
	case 826:
		return "TradeAllocIndicator", true
	case 827:
		return "ExpirationCycle", true
	case 828:
		return "TrdType", true
	case 829:
		return "TrdSubType", true
	case 830:
		return "TransferReason", true
	case 832:
		return "TotNumAssignmentReports", true
	case 833:
		return "AsgnRptID", true
	case 834:
		return "ThresholdAmount", true
	case 835:
		return "PegMoveType", true
	case 836:
		return "PegOffsetType", true
	case 837:
		return "PegLimitType", true
	case 838:
		return "PegRoundDirection", true
	case 839:
		return "PeggedPrice", true
	case 840:
		return "PegScope", true
	case 841:
		return "DiscretionMoveType", true
	case 842:
		return "DiscretionOffsetType", true
	case 843:
		return "DiscretionLimitType", true
	case 844:
		return "DiscretionRoundDirection", true
	case 845:
		return "DiscretionPrice", true
	case 846:
		return "DiscretionScope", true
	case 847:
		return "TargetStrategy", true
	case 848:
		return "TargetStrategyParameters", true
	case 849:
		return "ParticipationRate", true
	case 850:
		return "TargetStrategyPerformance", true
	case 851:
		return "LastLiquidityInd", true
	case 852:
		return "PublishTrdIndicator", true
	case 853:
		return "ShortSaleReason", true
	case 854:
		return "QtyType", true
	case 855:
		return "SecondaryTrdType", true
	case 856:
		return "TradeReportType", true
	case 857:
		return "AllocNoOrdersType", true
	case 858:
		return "SharedCommission", true
	case 859:
		return "ConfirmReqID", true
	case 860:
		return "AvgParPx", true
	case 861:
		return "ReportedPx", true
	case 862:
		return "NoCapacities", true
	case 863:
		return "OrderCapacityQty", true
	case 864:
		return "NoEvents", true
	case 865:
		return "EventType", true
	case 866:
		return "EventDate", true
	case 867:
		return "EventPx", true
	case 868:
		return "EventText", true
	case 869:
		return "PctAtRisk", true
	case 870:
		return "NoInstrAttrib", true
	case 871:
		return "InstrAttribType", true
	case 872:
		return "InstrAttribValue", true
	case 873:
		return "DatedDate", true
	case 874:
		return "InterestAccrualDate", true
	case 875:
		return "CPProgram", true
	case 876:
		return "CPRegType", true
	case 877:
		return "UnderlyingCPProgram", true
	case 878:
		return "UnderlyingCPRegType", true
	case 879:
		return "UnderlyingQty", true
	case 880:
		return "TrdMatchID", true
	case 881:
		return "SecondaryTradeReportRefID", true
	case 882:
		return "UnderlyingDirtyPrice", true
	case 883:
		return "UnderlyingEndPrice", true
	case 884:
		return "UnderlyingStartValue", true
	case 885:
		return "UnderlyingCurrentValue", true
	case 886:
		return "UnderlyingEndValue", true
	case 887:
		return "NoUnderlyingStips", true
	case 888:
		return "UnderlyingStipType", true
	case 889:
		return "UnderlyingStipValue", true
	case 890:
		return "MaturityNetMoney", true
	case 891:
		return "MiscFeeBasis", true
	case 892:
		return "TotNoAllocs", true
	case 893:
		return "LastFragment", true
	case 894:
		return "CollReqID", true
	case 895:
		return "CollAsgnReason", true
	case 896:
		return "CollInquiryQualifier", true
	case 897:
		return "NoTrades", true
	case 898:
		return "MarginRatio", true
	case 899:
		return "MarginExcess", true
	case 900:
		return "TotalNetValue", true
	case 901:
		return "CashOutstanding", true
	case 902:
		return "CollAsgnID", true
	case 903:
		return "CollAsgnTransType", true
	case 904:
		return "CollRespID", true
	case 905:
		return "CollAsgnRespType", true
	case 906:
		return "CollAsgnRejectReason", true
	case 907:
		return "CollAsgnRefID", true
	case 908:
		return "CollRptID", true
	case 909:
		return "CollInquiryID", true
	case 910:
		return "CollStatus", true
	case 911:
		return "TotNumReports", true
	case 912:
		return "LastRptRequested", true
	case 913:
		return "AgreementDesc", true
	case 914:
		return "AgreementID", true
	case 915:
		return "AgreementDate", true
	case 916:
		return "StartDate", true
	case 917:
		return "EndDate", true
	case 918:
		return "AgreementCurrency", true
	case 919:
		return "DeliveryType", true
	case 920:
		return "EndAccruedInterestAmt", true
	case 921:
		return "StartCash", true
	case 922:
		return "EndCash", true
	case 923:
		return "UserRequestID", true
	case 924:
		return "UserRequestType", true
	case 925:
		return "NewPassword", true
	case 926:
		return "UserStatus", true
	case 927:
		return "UserStatusText", true
	case 928:
		return "StatusValue", true
	case 929:
		return "StatusText", true
	case 930:
		return "RefCompID", true
	case 931:
		return "RefSubID", true
	case 932:
		return "NetworkResponseID", true
	case 933:
		return "NetworkRequestID", true
	case 934:
		return "LastNetworkResponseID", true
	case 935:
		return "NetworkRequestType", true
	case 936:
		return "NoCompIDs", true
	case 937:
		return "NetworkStatusResponseType", true
	case 938:
		return "NoCollInquiryQualifier", true
	case 939:
		return "TrdRptStatus", true
	case 940:
		return "AffirmStatus", true
	case 941:
		return "UnderlyingStrikeCurrency", true
	case 942:
		return "LegStrikeCurrency", true
	case 943:
		return "TimeBracket", true
	case 944:
		return "CollAction", true
	case 945:
		return "CollInquiryStatus", true
	case 946:
		return "CollInquiryResult", true
	case 947:
		return "StrikeCurrency", true
	case 948:
		return "NoNested3PartyIDs", true
	case 949:
		return "Nested3PartyID", true
	case 950:
		return "Nested3PartyIDSource", true
	case 951:
		return "Nested3PartyRole", true
	case 952:
		return "NoNested3PartySubIDs", true
	case 953:
		return "Nested3PartySubID", true
	case 954:
		return "Nested3PartySubIDType", true
	case 955:
		return "LegContractSettlMonth", true
	case 956:
		return "LegInterestAccrualDate", true
	case 957:
		return "NoStrategyParameters", true
	case 958:
		return "StrategyParameterName", true
	case 959:
		return "StrategyParameterType", true
	case 960:
		return "StrategyParameterValue", true
	case 961:
		return "HostCrossID", true
	case 962:
		return "SideTimeInForce", true
	case 963:
		return "MDReportID", true
	case 964:
		return "SecurityReportID", true
	case 965:
		return "SecurityStatus", true
	case 966:
		return "SettleOnOpenFlag", true
	case 967:
		return "StrikeMultiplier", true
	case 968:
		return "StrikeValue", true
	case 969:
		return "MinPriceIncrement", true
	case 970:
		return "PositionLimit", true
	case 971:
		return "NTPositionLimit", true
	case 972:
		return "UnderlyingAllocationPercent", true
	case 973:
		return "UnderlyingCashAmount", true
	case 974:
		return "UnderlyingCashType", true
	case 975:
		return "UnderlyingSettlementType", true
	case 976:
		return "QuantityDate", true
	case 977:
		return "ContIntRptID", true
	case 978:
		return "LateIndicator", true
	case 979:
		return "InputSource", true
	case 980:
		return "SecurityUpdateAction", true
	case 981:
		return "NoExpiration", true
	case 982:
		return "ExpirationQtyType", true
	case 983:
		return "ExpQty", true
	case 984:
		return "NoUnderlyingAmounts", true
	case 985:
		return "UnderlyingPayAmount", true
	case 986:
		return "UnderlyingCollectAmount", true
	case 987:
		return "UnderlyingSettlementDate", true
	case 988:
		return "UnderlyingSettlementStatus", true
	case 989:
		return "SecondaryIndividualAllocID", true
	case 990:
		return "LegReportID", true
	case 991:
		return "RndPx", true
	case 992:
		return "IndividualAllocType", true
	case 993:
		return "AllocCustomerCapacity", true
	case 994:
		return "TierCode", true
	case 996:
		return "UnitOfMeasure", true
	case 997:
		return "TimeUnit", true
	case 998:
		return "UnderlyingUnitOfMeasure", true
	case 999:
		return "LegUnitOfMeasure", true
	case 1000:
		return "UnderlyingTimeUnit", true
	case 1001:
		return "LegTimeUnit", true
	case 1002:
		return "AllocMethod", true
	case 1003:
		return "TradeID", true
	case 1005:
		return "SideTradeReportID", true
	case 1006:
		return "SideFillStationCd", true
	case 1007:
		return "SideReasonCd", true
	case 1008:
		return "SideTrdSubTyp", true
	case 1009:
		return "SideQty", true
	case 1011:
		return "MessageEventSource", true
	case 1012:
		return "SideTrdRegTimestamp", true
	case 1013:
		return "SideTrdRegTimestampType", true
	case 1014:
		return "SideTrdRegTimestampSrc", true
	case 1015:
		return "AsOfIndicator", true
	case 1016:
		return "NoSideTrdRegTS", true
	case 1017:
		return "LegOptionRatio", true
	case 1018:
		return "NoInstrumentParties", true
	case 1019:
		return "InstrumentPartyID", true
	case 1020:
		return "TradeVolume", true
	case 1021:
		return "MDBookType", true
	case 1022:
		return "MDFeedType", true
	case 1023:
		return "MDPriceLevel", true
	case 1024:
		return "MDOriginType", true
	case 1025:
		return "FirstPx", true
	case 1026:
		return "MDEntrySpotRate", true
	case 1027:
		return "MDEntryForwardPoints", true
	case 1028:
		return "ManualOrderIndicator", true
	case 1029:
		return "CustDirectedOrder", true
	case 1030:
		return "ReceivedDeptID", true
	case 1031:
		return "CustOrderHandlingInst", true
	case 1032:
		return "OrderHandlingInstSource", true
	case 1033:
		return "DeskType", true
	case 1034:
		return "DeskTypeSource", true
	case 1035:
		return "DeskOrderHandlingInst", true
	case 1036:
		return "ExecAckStatus", true
	case 1037:
		return "UnderlyingDeliveryAmount", true
	case 1038:
		return "UnderlyingCapValue", true
	case 1039:
		return "UnderlyingSettlMethod", true
	case 1040:
		return "SecondaryTradeID", true
	case 1041:
		return "FirmTradeID", true
	case 1042:
		return "SecondaryFirmTradeID", true
	case 1043:
		return "CollApplType", true
	case 1044:
		return "UnderlyingAdjustedQuantity", true
	case 1045:
		return "UnderlyingFXRate", true
	case 1046:
		return "UnderlyingFXRateCalc", true
	case 1047:
		return "AllocPositionEffect", true
	case 1048:
		return "DealingCapacity", true
	case 1049:
		return "InstrmtAssignmentMethod", true
	case 1050:
		return "InstrumentPartyIDSource", true
	case 1051:
		return "InstrumentPartyRole", true
	case 1052:
		return "NoInstrumentPartySubIDs", true
	case 1053:
		return "InstrumentPartySubID", true
	case 1054:
		return "InstrumentPartySubIDType", true
	case 1055:
		return "PositionCurrency", true
	case 1056:
		return "CalculatedCcyLastQty", true
	case 1057:
		return "AggressorIndicator", true
	case 1058:
		return "NoUndlyInstrumentParties", true
	case 1059:
		return "UndlyInstrumentPartyID", true
	case 1060:
		return "UndlyInstrumentPartyIDSource", true
	case 1061:
		return "UndlyInstrumentPartyRole", true
	case 1062:
		return "NoUndlyInstrumentPartySubIDs", true
	case 1063:
		return "UndlyInstrumentPartySubID", true
	case 1064:
		return "UndlyInstrumentPartySubIDType", true
	case 1065:
		return "BidSwapPoints", true
	case 1066:
		return "OfferSwapPoints", true
	case 1067:
		return "LegBidForwardPoints", true
	case 1068:
		return "LegOfferForwardPoints", true
	case 1069:
		return "SwapPoints", true
	case 1070:
		return "MDQuoteType", true
	case 1071:
		return "LastSwapPoints", true
	case 1072:
		return "SideGrossTradeAmt", true
	case 1073:
		return "LegLastForwardPoints", true
	case 1074:
		return "LegCalculatedCcyLastQty", true
	case 1075:
		return "LegGrossTradeAmt", true
	case 1079:
		return "MaturityTime", true
	case 1080:
		return "RefOrderID", true
	case 1081:
		return "RefOrderIDSource", true
	case 1082:
		return "SecondaryDisplayQty", true
	case 1083:
		return "DisplayWhen", true
	case 1084:
		return "DisplayMethod", true
	case 1085:
		return "DisplayLowQty", true
	case 1086:
		return "DisplayHighQty", true
	case 1087:
		return "DisplayMinIncr", true
	case 1088:
		return "RefreshQty", true
	case 1089:
		return "MatchIncrement", true
	case 1090:
		return "MaxPriceLevels", true
	case 1091:
		return "PreTradeAnonymity", true
	case 1092:
		return "PriceProtectionScope", true
	case 1093:
		return "LotType", true
	case 1094:
		return "PegPriceType", true
	case 1095:
		return "PeggedRefPrice", true
	case 1096:
		return "PegSecurityIDSource", true
	case 1097:
		return "PegSecurityID", true
	case 1098:
		return "PegSymbol", true
	case 1099:
		return "PegSecurityDesc", true
	case 1100:
		return "TriggerType", true
	case 1101:
		return "TriggerAction", true
	case 1102:
		return "TriggerPrice", true
	case 1103:
		return "TriggerSymbol", true
	case 1104:
		return "TriggerSecurityID", true
	case 1105:
		return "TriggerSecurityIDSource", true
	case 1106:
		return "TriggerSecurityDesc", true
	case 1107:
		return "TriggerPriceType", true
	case 1108:
		return "TriggerPriceTypeScope", true
	case 1109:
		return "TriggerPriceDirection", true
	case 1110:
		return "TriggerNewPrice", true
	case 1111:
		return "TriggerOrderType", true
	case 1112:
		return "TriggerNewQty", true
	case 1113:
		return "TriggerTradingSessionID", true
	case 1114:
		return "TriggerTradingSessionSubID", true
	case 1115:
		return "OrderCategory", true
	case 1116:
		return "NoRootPartyIDs", true
	case 1117:
		return "RootPartyID", true
	case 1118:
		return "RootPartyIDSource", true
	case 1119:
		return "RootPartyRole", true
	case 1120:
		return "NoRootPartySubIDs", true
	case 1121:
		return "RootPartySubID", true
	case 1122:
		return "RootPartySubIDType", true
	case 1123:
		return "TradeHandlingInstr", true
	case 1124:
		return "OrigTradeHandlingInstr", true
	case 1125:
		return "OrigTradeDate", true
	case 1126:
		return "OrigTradeID", true
	case 1127:
		return "OrigSecondaryTradeID", true
	case 1128:
		return "ApplVerID", true
	case 1129:
		return "CstmApplVerID", true
	case 1130:
		return "RefApplVerID", true
	case 1131:
		return "RefCstmApplVerID", true
	case 1132:
		return "TZTransactTime", true
	case 1133:
		return "ExDestinationIDSource", true
	case 1134:
		return "ReportedPxDiff", true
	case 1135:
		return "RptSys", true
	case 1136:
		return "AllocClearingFeeIndicator", true
	case 1137:
		return "DefaultApplVerID", true
	case 1138:
		return "DisplayQty", true
	case 1139:
		return "ExchangeSpecialInstructions", true
	case 1213:
		return "UnderlyingMaturityTime", true
	case 1212:
		return "LegMaturityTime", true
	case 1140:
		return "MaxTradeVol", true
	case 1141:
		return "NoMDFeedTypes", true
	case 1142:
		return "MatchAlgorithm", true
	case 1143:
		return "MaxPriceVariation", true
	case 1144:
		return "ImpliedMarketIndicator", true
	case 1145:
		return "EventTime", true
	case 1146:
		return "MinPriceIncrementAmount", true
	case 1147:
		return "UnitOfMeasureQty", true
	case 1148:
		return "LowLimitPrice", true
	case 1149:
		return "HighLimitPrice", true
	case 1150:
		return "TradingReferencePrice", true
	case 1151:
		return "SecurityGroup", true
	case 1152:
		return "LegNumber", true
	case 1153:
		return "SettlementCycleNo", true
	case 1154:
		return "SideCurrency", true
	case 1155:
		return "SideSettlCurrency", true
	case 1157:
		return "CcyAmt", true
	case 1158:
		return "NoSettlDetails", true
	case 1159:
		return "SettlObligMode", true
	case 1160:
		return "SettlObligMsgID", true
	case 1161:
		return "SettlObligID", true
	case 1162:
		return "SettlObligTransType", true
	case 1163:
		return "SettlObligRefID", true
	case 1164:
		return "SettlObligSource", true
	case 1165:
		return "NoSettlOblig", true
	case 1166:
		return "QuoteMsgID", true
	case 1167:
		return "QuoteEntryStatus", true
	case 1168:
		return "TotNoCxldQuotes", true
	case 1169:
		return "TotNoAccQuotes", true
	case 1170:
		return "TotNoRejQuotes", true
	case 1171:
		return "PrivateQuote", true
	case 1172:
		return "RespondentType", true
	case 1173:
		return "MDSubBookType", true
	case 1174:
		return "SecurityTradingEvent", true
	case 1175:
		return "NoStatsIndicators", true
	case 1176:
		return "StatsType", true
	case 1177:
		return "NoOfSecSizes", true
	case 1178:
		return "MDSecSizeType", true
	case 1179:
		return "MDSecSize", true
	case 1180:
		return "ApplID", true
	case 1181:
		return "ApplSeqNum", true
	case 1182:
		return "ApplBegSeqNum", true
	case 1183:
		return "ApplEndSeqNum", true
	case 1184:
		return "SecurityXMLLen", true
	case 1185:
		return "SecurityXML", true
	case 1186:
		return "SecurityXMLSchema", true
	case 1187:
		return "RefreshIndicator", true
	case 1188:
		return "Volatility", true
	case 1189:
		return "TimeToExpiration", true
	case 1190:
		return "RiskFreeRate", true
	case 1191:
		return "PriceUnitOfMeasure", true
	case 1192:
		return "PriceUnitOfMeasureQty", true
	case 1193:
		return "SettlMethod", true
	case 1194:
		return "ExerciseStyle", true
	case 1419:
		return "UnderlyingExerciseStyle", true
	case 1420:
		return "LegExerciseStyle", true
	case 1195:
		return "OptPayAmount", true
	case 1196:
		return "PriceQuoteMethod", true
	case 1197:
		return "FuturesValuationMethod", true
	case 1198:
		return "ListMethod", true
	case 1199:
		return "CapPrice", true
	case 1200:
		return "FloorPrice", true
	case 1201:
		return "NoStrikeRules", true
	case 1202:
		return "StartStrikePxRange", true
	case 1203:
		return "EndStrikePxRange", true
	case 1204:
		return "StrikeIncrement", true
	case 1205:
		return "NoTickRules", true
	case 1206:
		return "StartTickPriceRange", true
	case 1207:
		return "EndTickPriceRange", true
	case 1208:
		return "TickIncrement", true
	case 1209:
		return "TickRuleType", true
	case 1210:
		return "NestedInstrAttribType", true
	case 1211:
		return "NestedInstrAttribValue", true
	case 1214:
		return "DerivativeSymbol", true
	case 1215:
		return "DerivativeSymbolSfx", true
	case 1216:
		return "DerivativeSecurityID", true
	case 1217:
		return "DerivativeSecurityIDSource", true
	case 1218:
		return "NoDerivativeSecurityAltID", true
	case 1219:
		return "DerivativeSecurityAltID", true
	case 1220:
		return "DerivativeSecurityAltIDSource", true
	case 1221:
		return "SecondaryLowLimitPrice", true
	case 1230:
		return "SecondaryHighLimitPrice", true
	case 1222:
		return "MaturityRuleID", true
	case 1223:
		return "StrikeRuleID", true
	case 1225:
		return "DerivativeOptPayAmount", true
	case 1226:
		return "EndMaturityMonthYear", true
	case 1227:
		return "ProductComplex", true
	case 1228:
		return "DerivativeProductComplex", true
	case 1229:
		return "MaturityMonthYearIncrement", true
	case 1231:
		return "MinLotSize", true
	case 1232:
		return "NoExecInstRules", true
	case 1234:
		return "NoLotTypeRules", true
	case 1235:
		return "NoMatchRules", true
	case 1236:
		return "NoMaturityRules", true
	case 1237:
		return "NoOrdTypeRules", true
	case 1239:
		return "NoTimeInForceRules", true
	case 1240:
		return "SecondaryTradingReferencePrice", true
	case 1241:
		return "StartMaturityMonthYear", true
	case 1242:
		return "FlexProductEligibilityIndicator", true
	case 1243:
		return "DerivFlexProductEligibilityIndicator", true
	case 1244:
		return "FlexibleIndicator", true
	case 1245:
		return "TradingCurrency", true
	case 1246:
		return "DerivativeProduct", true
	case 1247:
		return "DerivativeSecurityGroup", true
	case 1248:
		return "DerivativeCFICode", true
	case 1249:
		return "DerivativeSecurityType", true
	case 1250:
		return "DerivativeSecuritySubType", true
	case 1251:
		return "DerivativeMaturityMonthYear", true
	case 1252:
		return "DerivativeMaturityDate", true
	case 1253:
		return "DerivativeMaturityTime", true
	case 1254:
		return "DerivativeSettleOnOpenFlag", true
	case 1255:
		return "DerivativeInstrmtAssignmentMethod", true
	case 1256:
		return "DerivativeSecurityStatus", true
	case 1257:
		return "DerivativeInstrRegistry", true
	case 1258:
		return "DerivativeCountryOfIssue", true
	case 1259:
		return "DerivativeStateOrProvinceOfIssue", true
	case 1260:
		return "DerivativeLocaleOfIssue", true
	case 1261:
		return "DerivativeStrikePrice", true
	case 1262:
		return "DerivativeStrikeCurrency", true
	case 1263:
		return "DerivativeStrikeMultiplier", true
	case 1264:
		return "DerivativeStrikeValue", true
	case 1265:
		return "DerivativeOptAttribute", true
	case 1266:
		return "DerivativeContractMultiplier", true
	case 1267:
		return "DerivativeMinPriceIncrement", true
	case 1268:
		return "DerivativeMinPriceIncrementAmount", true
	case 1269:
		return "DerivativeUnitOfMeasure", true
	case 1270:
		return "DerivativeUnitOfMeasureQty", true
	case 1271:
		return "DerivativeTimeUnit", true
	case 1272:
		return "DerivativeSecurityExchange", true
	case 1273:
		return "DerivativePositionLimit", true
	case 1274:
		return "DerivativeNTPositionLimit", true
	case 1275:
		return "DerivativeIssuer", true
	case 1276:
		return "DerivativeIssueDate", true
	case 1277:
		return "DerivativeEncodedIssuerLen", true
	case 1278:
		return "DerivativeEncodedIssuer", true
	case 1279:
		return "DerivativeSecurityDesc", true
	case 1280:
		return "DerivativeEncodedSecurityDescLen", true
	case 1281:
		return "DerivativeEncodedSecurityDesc", true
	case 1282:
		return "DerivativeSecurityXMLLen", true
	case 1283:
		return "DerivativeSecurityXML", true
	case 1284:
		return "DerivativeSecurityXMLSchema", true
	case 1285:
		return "DerivativeContractSettlMonth", true
	case 1286:
		return "NoDerivativeEvents", true
	case 1287:
		return "DerivativeEventType", true
	case 1288:
		return "DerivativeEventDate", true
	case 1289:
		return "DerivativeEventTime", true
	case 1290:
		return "DerivativeEventPx", true
	case 1291:
		return "DerivativeEventText", true
	case 1292:
		return "NoDerivativeInstrumentParties", true
	case 1293:
		return "DerivativeInstrumentPartyID", true
	case 1294:
		return "DerivativeInstrumentPartyIDSource", true
	case 1295:
		return "DerivativeInstrumentPartyRole", true
	case 1296:
		return "NoDerivativeInstrumentPartySubIDs", true
	case 1297:
		return "DerivativeInstrumentPartySubID", true
	case 1298:
		return "DerivativeInstrumentPartySubIDType", true
	case 1299:
		return "DerivativeExerciseStyle", true
	case 1300:
		return "MarketSegmentID", true
	case 1301:
		return "MarketID", true
	case 1302:
		return "MaturityMonthYearIncrementUnits", true
	case 1303:
		return "MaturityMonthYearFormat", true
	case 1304:
		return "StrikeExerciseStyle", true
	case 1305:
		return "SecondaryPriceLimitType", true
	case 1306:
		return "PriceLimitType", true
	case 1308:
		return "ExecInstValue", true
	case 1309:
		return "NoTradingSessionRules", true
	case 1310:
		return "NoMarketSegments", true
	case 1311:
		return "NoDerivativeInstrAttrib", true
	case 1312:
		return "NoNestedInstrAttrib", true
	case 1313:
		return "DerivativeInstrAttribType", true
	case 1314:
		return "DerivativeInstrAttribValue", true
	case 1315:
		return "DerivativePriceUnitOfMeasure", true
	case 1316:
		return "DerivativePriceUnitOfMeasureQty", true
	case 1317:
		return "DerivativeSettlMethod", true
	case 1318:
		return "DerivativePriceQuoteMethod", true
	case 1319:
		return "DerivativeFuturesValuationMethod", true
	case 1320:
		return "DerivativeListMethod", true
	case 1321:
		return "DerivativeCapPrice", true
	case 1322:
		return "DerivativeFloorPrice", true
	case 1323:
		return "DerivativePutOrCall", true
	case 1324:
		return "ListUpdateAction", true
	case 1358:
		return "LegPutOrCall", true
	case 1224:
		return "LegUnitOfMeasureQty", true
	case 1421:
		return "LegPriceUnitOfMeasure", true
	case 1422:
		return "LegPriceUnitOfMeasureQty", true
	case 1423:
		return "UnderlyingUnitOfMeasureQty", true
	case 1424:
		return "UnderlyingPriceUnitOfMeasure", true
	case 1425:
		return "UnderlyingPriceUnitOfMeasureQty", true
	case 1393:
		return "MarketReqID", true
	case 1394:
		return "MarketReportID", true
	case 1395:
		return "MarketUpdateAction", true
	case 1396:
		return "MarketSegmentDesc", true
	case 1397:
		return "EncodedMktSegmDescLen", true
	case 1398:
		return "EncodedMktSegmDesc", true
	case 1325:
		return "ParentMktSegmID", true
	case 1326:
		return "TradingSessionDesc", true
	case 1327:
		return "TradSesUpdateAction", true
	case 1328:
		return "RejectText", true
	case 1329:
		return "FeeMultiplier", true
	case 1330:
		return "UnderlyingLegSymbol", true
	case 1331:
		return "UnderlyingLegSymbolSfx", true
	case 1332:
		return "UnderlyingLegSecurityID", true
	case 1333:
		return "UnderlyingLegSecurityIDSource", true
	case 1334:
		return "NoUnderlyingLegSecurityAltID", true
	case 1335:
		return "UnderlyingLegSecurityAltID", true
	case 1336:
		return "UnderlyingLegSecurityAltIDSource", true
	case 1337:
		return "UnderlyingLegSecurityType", true
	case 1338:
		return "UnderlyingLegSecuritySubType", true
	case 1339:
		return "UnderlyingLegMaturityMonthYear", true
	case 1343:
		return "UnderlyingLegPutOrCall", true
	case 1340:
		return "UnderlyingLegStrikePrice", true
	case 1341:
		return "UnderlyingLegSecurityExchange", true
	case 1342:
		return "NoOfLegUnderlyings", true
	case 1344:
		return "UnderlyingLegCFICode", true
	case 1345:
		return "UnderlyingLegMaturityDate", true
	case 1405:
		return "UnderlyingLegMaturityTime", true
	case 1391:
		return "UnderlyingLegOptAttribute", true
	case 1392:
		return "UnderlyingLegSecurityDesc", true
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
	case 1156:
		return "ApplExtID", true
	case 1406:
		return "RefApplExtID", true
	case 1407:
		return "DefaultApplExtID", true
	case 1408:
		return "DefaultCstmApplVerID", true
	case 1409:
		return "SessionStatus", true
	case 1410:
		return "DefaultVerIndicator", true
	case 809:
		return "NoUsernames", true
	case 1367:
		return "LegAllocSettlCurrency", true
	case 1361:
		return "TotNoFills", true
	case 1362:
		return "NoFills", true
	case 1363:
		return "FillExecID", true
	case 1364:
		return "FillPx", true
	case 1365:
		return "FillQty", true
	case 1366:
		return "LegAllocID", true
	case 1368:
		return "TradSesEvent", true
	case 1369:
		return "MassActionReportID", true
	case 1370:
		return "NoNotAffectedOrders", true
	case 1371:
		return "NotAffectedOrderID", true
	case 1372:
		return "NotAffOrigClOrdID", true
	case 1373:
		return "MassActionType", true
	case 1374:
		return "MassActionScope", true
	case 1375:
		return "MassActionResponse", true
	case 1376:
		return "MassActionRejectReason", true
	case 1377:
		return "MultilegModel", true
	case 1378:
		return "MultilegPriceMethod", true
	case 1379:
		return "LegVolatility", true
	case 1380:
		return "DividendYield", true
	case 1381:
		return "LegDividendYield", true
	case 1382:
		return "CurrencyRatio", true
	case 1383:
		return "LegCurrencyRatio", true
	case 1384:
		return "LegExecInst", true
	case 1385:
		return "ContingencyType", true
	case 1386:
		return "ListRejectReason", true
	case 1387:
		return "NoTrdRepIndicators", true
	case 1388:
		return "TrdRepPartyRole", true
	case 1389:
		return "TrdRepIndicator", true
	case 1390:
		return "TradePublishIndicator", true
	case 1346:
		return "ApplReqID", true
	case 1347:
		return "ApplReqType", true
	case 1348:
		return "ApplResponseType", true
	case 1349:
		return "ApplTotalMessageCount", true
	case 1350:
		return "ApplLastSeqNum", true
	case 1351:
		return "NoApplIDs", true
	case 1352:
		return "ApplResendFlag", true
	case 1353:
		return "ApplResponseID", true
	case 1354:
		return "ApplResponseError", true
	case 1355:
		return "RefApplID", true
	case 1356:
		return "ApplReportID", true
	case 1357:
		return "RefApplLastSeqNum", true
	case 1399:
		return "ApplNewSeqNum", true
	case 1426:
		return "ApplReportType", true
	case 1411:
		return "Nested4PartySubIDType", true
	case 1412:
		return "Nested4PartySubID", true
	case 1413:
		return "NoNested4PartySubIDs", true
	case 1414:
		return "NoNested4PartyIDs", true
	case 1415:
		return "Nested4PartyID", true
	case 1416:
		return "Nested4PartyIDSource", true
	case 1417:
		return "Nested4PartyRole", true
	case 1418:
		return "LegLastQty", true
	default:
		return "", false
	}
}

func (f FIX50SP1) ValueName(tag int, value string) (string, bool) {
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
		case "N":
			return `New`, true
		case "C":
			return `Cancel`, true
		case "R":
			return `Replace`, true
		}

	case 13:
		switch value {
		case "1":
			return `Per Unit (implying shares, par, currency, etc.)`, true
		case "2":
			return `Percent`, true
		case "3":
			return `Absolute (total monetary amount)`, true
		case "4":
			return `Percentage waived - cash discount (for CIV buy orders)`, true
		case "5":
			return `Percentage waived -= enhanced units (for CIV buy orders)`, true
		case "6":
			return `Points per bond or contract (supply ContractMultiplier (231) in the <Instrument> component block if the object security is denominated in a size other than the industry default - 1000 par for bonds)`, true
		}

	case 18:
		switch value {
		case "0":
			return `Stay on offer side`, true
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
			return `Stay on bid side`, true
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
		case "H":
			return `Reinstate on system failure (mutually exclusive with Q and l)`, true
		case "I":
			return `Institutions only`, true
		case "J":
			return `Reinstate on Trading Halt (mutually exclusive with K and m)`, true
		case "K":
			return `Cancel on Trading Halt (mutually exclusive with J and m)`, true
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
		case "Q":
			return `Cancel on system failure (mutually exclusive with H and l)`, true
		case "R":
			return `Primary peg (primary market - buy at bid/sell at offer)`, true
		case "S":
			return `Suspend`, true
		case "T":
			return `Fixed Peg to Local best bid or offer at time of order`, true
		case "U":
			return `Customer Display Instruction (Rule 11Ac1-1/4)`, true
		case "V":
			return `Netting (for Forex)`, true
		case "W":
			return `Peg to VWAP`, true
		case "X":
			return `Trade Along`, true
		case "Y":
			return `Try To Stop`, true
		case "Z":
			return `Cancel if not best`, true
		case "a":
			return `Trailing Stop Peg`, true
		case "b":
			return `Strict Limit (No price improvement)`, true
		case "c":
			return `Ignore Price Validity Checks`, true
		case "d":
			return `Peg to Limit Price`, true
		case "e":
			return `Work to Target Strategy`, true
		case "f":
			return `Intermarket Sweep`, true
		case "g":
			return `External Routing Allowed`, true
		case "h":
			return `External Routing Not Allowed`, true
		case "i":
			return `Imbalance Only`, true
		case "j":
			return `Single execution requested for block trade`, true
		case "k":
			return `Best Execution`, true
		case "l":
			return `Suspend on system failure (mutually exclusive with H and Q)`, true
		case "m":
			return `Suspend on Trading Halt (mutually exclusive with J and K)`, true
		case "n":
			return `Reinstate on connection loss (mutually exclusive with o and p)`, true
		case "o":
			return `Cancel on connection loss (mutually exclusive with n and p)`, true
		case "p":
			return `Suspend on connection loss (mutually exclusive with n and o)`, true
		case "q":
			return `Release from suspension (mutually exclusive with S)`, true
		case "r":
			return `Execute as delta neutral using volatility provided`, true
		case "s":
			return `Execute as duration neutral`, true
		case "t":
			return `Execute as FX neutral`, true
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
		case "8":
			return `Exchange Symbol`, true
		case "9":
			return `Consolidated Tape Association (CTA) Symbol (SIAC CTS/CQS line format)`, true
		case "A":
			return `Bloomberg Symbol`, true
		case "B":
			return `Wertpapier`, true
		case "C":
			return `Dutch`, true
		case "D":
			return `Valoren`, true
		case "E":
			return `Sicovam`, true
		case "F":
			return `Belgian`, true
		case "G":
			return `"Common" (Clearstream and Euroclear)`, true
		case "H":
			return `Clearing House / Clearing Organization`, true
		case "I":
			return `ISDA/FpML Product Specification (XML in EncodedSecurityDesc)`, true
		case "J":
			return `Option Price Reporting Authority`, true
		case "K":
			return `ISDA/FpML Product URL (URL in SecurityID)`, true
		case "L":
			return `Letter of Credit`, true
		case "M":
			return `Marketplace-assigned Identifier`, true
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
		case "S":
			return `Small`, true
		case "M":
			return `Medium`, true
		case "L":
			return `Large`, true
		case "U":
			return `Undisclosed Quantity`, true
		}

	case 28:
		switch value {
		case "N":
			return `New`, true
		case "C":
			return `Cancel`, true
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
			return `Replaced (No longer used)`, true
		case "6":
			return `Pending Cancel (i.e. result of Order Cancel Request)`, true
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
			return `Accepted for Bidding`, true
		case "E":
			return `Pending Replace (i.e. result of Order Cancel/Replace Request)`, true
		}

	case 40:
		switch value {
		case "1":
			return `Market`, true
		case "2":
			return `Limit`, true
		case "3":
			return `Stop / Stop Loss`, true
		case "4":
			return `Stop Limit`, true
		case "5":
			return `Market On Close (No longer used)`, true
		case "6":
			return `With Or Without`, true
		case "7":
			return `Limit Or Better`, true
		case "8":
			return `Limit With Or Without`, true
		case "9":
			return `On Basis`, true
		case "A":
			return `On Close (No longer used)`, true
		case "B":
			return `Limit On Close (No longer used)`, true
		case "C":
			return `Forex Market (No longer used)`, true
		case "D":
			return `Previously Quoted`, true
		case "E":
			return `Previously Indicated`, true
		case "F":
			return `Forex Limit (No longer used)`, true
		case "G":
			return `Forex Swap`, true
		case "H":
			return `Forex Previously Quoted (No longer used)`, true
		case "I":
			return `Funari (Limit day order with unexecuted portion handles as Market On Close. E.g. Japan)`, true
		case "J":
			return `Market If Touched (MIT)`, true
		case "K":
			return `Market With Left Over as Limit (market order with unexecuted quantity becoming limit order at last price)`, true
		case "L":
			return `Previous Fund Valuation Point (Historic pricing; for CIV)`, true
		case "M":
			return `Next Fund Valuation Point (Forward pricing; for CIV)`, true
		case "P":
			return `Pegged`, true
		case "Q":
			return `Counter-order selection`, true
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
			return `Proprietary, Non-Algorithmic Program Trade (non-index arbitrage)`, true
		case "D":
			return `Program order, index arb, for Member firm/org`, true
		case "E":
			return `Short Exempt Transaction for Principal (was incorrectly identified in the FIX spec as "Registered Equity Market Maker trades")`, true
		case "F":
			return `Short exempt transaction (refer to W type)`, true
		case "H":
			return `Short exempt transaction (refer to I type)`, true
		case "I":
			return `Individual Investor, single order`, true
		case "J":
			return `Proprietary, Algorithmic Program Trading (non-index arbitrage)`, true
		case "K":
			return `Agency, Algorithmic Program Trading (non-index arbitrage)`, true
		case "L":
			return `Short exempt transaction for member competing market-maker affliated with the firm clearing the trade (refer to P and O types)`, true
		case "M":
			return `Program Order, index arb, for other member`, true
		case "N":
			return `Agent for other Member, Non-Algorithmic Program Trade (non-index arbitrage)`, true
		case "O":
			return `Proprietary transactions for competing market-maker that is affiliated with the clearing member (was incorrectly identified in the FIX spec as "Competing dealer trades")`, true
		case "P":
			return `Principal`, true
		case "R":
			return `Transactions for the account of a non-member compting market-maker (was incorrectly identified in the FIX spec as "Competing dealer trades")`, true
		case "S":
			return `Specialist trades`, true
		case "T":
			return `Transactions for the account of an unaffiliated member's competing market-maker (was incorrectly identified in the FIX spec as "Competing dealer trades")`, true
		case "U":
			return `Agency, Index Arbitrage (includes Individual, Index Arbitrage trades)`, true
		case "W":
			return `All other orders as agent for other member`, true
		case "X":
			return `Short exempt transaction for member competing market-maker not affiliated with the firm clearing the trade (refer to W and T types)`, true
		case "Y":
			return `Agency, Non-Algorithmic Program Trade (non-index arbitrage)`, true
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
			return `Undisclosed (valid for IOI and List Order messages only)`, true
		case "8":
			return `Cross (orders where counterparty is an exchange, valid for all messages except IOIs)`, true
		case "9":
			return `Cross short`, true
		case "A":
			return `Cross short exxmpt`, true
		case "B":
			return `"As Defined" (for use with multileg instruments)`, true
		case "C":
			return `"Opposite" (for use with multileg instruments)`, true
		case "D":
			return `Subscribe (e.g. CIV)`, true
		case "E":
			return `Redeem (e.g. CIV)`, true
		case "F":
			return `Lend (FINANCING - identifies direction of collateral)`, true
		case "G":
			return `Borrow (FINANCING - identifies direction of collateral)`, true
		}

	case 59:
		switch value {
		case "0":
			return `Day (or session)`, true
		case "1":
			return `Good Till Cancel (GTC)`, true
		case "2":
			return `At the Opening (OPG)`, true
		case "3":
			return `Immediate Or Cancel (IOC)`, true
		case "4":
			return `Fill Or Kill (FOK)`, true
		case "5":
			return `Good Till Crossing (GTX)`, true
		case "6":
			return `Good Till Date (GTD)`, true
		case "7":
			return `At the Close`, true
		case "8":
			return `Good Through Crossing`, true
		case "9":
			return `At Crossing`, true
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
			return `Regular / FX Spot settlement (T+1 or T+2 depending on currency)`, true
		case "1":
			return `Cash (TOD / T+0)`, true
		case "2":
			return `Next Day (TOM / T+1)`, true
		case "3":
			return `T+2`, true
		case "4":
			return `T+3`, true
		case "5":
			return `T+4`, true
		case "6":
			return `Future`, true
		case "7":
			return `When And If Issued`, true
		case "8":
			return `Sellers Option`, true
		case "9":
			return `T+5`, true
		case "B":
			return `Broken date - for FX expressing non-standard tenor, SettlDate (64) must be specified`, true
		case "C":
			return `FX Spot Next settlement (Spot+1, aka next day)`, true
		}

	case 65:
		switch value {
		case "CD":
			return `EUCP with lump-sum interest rather than discount price`, true
		case "WI":
			return `"When Issued" for a security to be reissued under an old CUSIP or ISIN`, true
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
			return `Preliminary (without MiscFees and NetMoney) (Removed/Replaced)`, true
		case "4":
			return `Calculated (includes MiscFees and NetMoney) (Removed/Replaced)`, true
		case "5":
			return `Calculated without Preliminary (sent unsolicited by broker, includes MiscFees and NetMoney) (Removed/Replaced)`, true
		case "6":
			return `Reversal`, true
		}

	case 77:
		switch value {
		case "C":
			return `Close`, true
		case "F":
			return `FIFO`, true
		case "O":
			return `Open`, true
		case "R":
			return `Rolled`, true
		case "N":
			return `Close but notify on open`, true
		case "D":
			return `Default`, true
		}

	case 81:
		switch value {
		case "0":
			return `Regular`, true
		case "1":
			return `Soft Dollar`, true
		case "2":
			return `Step-In`, true
		case "3":
			return `Step-Out`, true
		case "4":
			return `Soft-dollar Step-In`, true
		case "5":
			return `Soft-dollar Step-Out`, true
		case "6":
			return `Plan Sponsor`, true
		}

	case 87:
		switch value {
		case "0":
			return `accepted (successfully processed)`, true
		case "1":
			return `block level reject`, true
		case "2":
			return `account level reject`, true
		case "3":
			return `received (received, not yet processed)`, true
		case "4":
			return `incomplete`, true
		case "5":
			return `rejected by intermediary`, true
		case "6":
			return `allocation pending`, true
		case "7":
			return `reversed`, true
		}

	case 88:
		switch value {
		case "0":
			return `Unknown account(s)`, true
		case "1":
			return `Incorrect quantity`, true
		case "2":
			return `Incorrect averageg price`, true
		case "3":
			return `Unknown executing broker mnemonic`, true
		case "4":
			return `Commission difference`, true
		case "5":
			return `Unknown OrderID (37)`, true
		case "6":
			return `Unknown ListID (66)`, true
		case "7":
			return `Other (further in Text (58))`, true
		case "8":
			return `Incorrect allocated quantity`, true
		case "9":
			return `Calculation difference`, true
		case "10":
			return `Unknown or stale ExecID`, true
		case "11":
			return `Mismatched data`, true
		case "12":
			return `Unknown ClOrdID`, true
		case "13":
			return `Warehouse request rejected`, true
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

	case 102:
		switch value {
		case "0":
			return `Too late to cancel`, true
		case "1":
			return `Unknown order`, true
		case "2":
			return `Broker / Exchange Option`, true
		case "3":
			return `Order already in Pending Cancel or Pending Replace status`, true
		case "4":
			return `Unable to process Order Mass Cancel Request`, true
		case "5":
			return `OrigOrdModTime (586) did not match last TransactTime (60) of order`, true
		case "6":
			return `Duplicate ClOrdID (11) received`, true
		case "7":
			return `Price exceeds current price`, true
		case "8":
			return `Price exceeds current price band`, true
		case "18":
			return `Invalid price increment`, true
		case "99":
			return `Other`, true
		}

	case 103:
		switch value {
		case "0":
			return `Broker / Exchange option`, true
		case "1":
			return `Unknown symbol`, true
		case "2":
			return `Exchange closed`, true
		case "3":
			return `Order exceeds limit`, true
		case "4":
			return `Too late to enter`, true
		case "5":
			return `Unknown order`, true
		case "6":
			return `Duplicate Order (e.g. dupe ClOrdID)`, true
		case "7":
			return `Duplicate of a verbally communicated order`, true
		case "8":
			return `Stale order`, true
		case "9":
			return `Trade along required`, true
		case "10":
			return `Invalid Investor ID`, true
		case "11":
			return `Unsupported order characteristic`, true
		case "12":
			return `Surveillence Option`, true
		case "13":
			return `Incorrect quantity`, true
		case "14":
			return `Incorrect allocated quantity`, true
		case "15":
			return `Unknown account(s)`, true
		case "16":
			return `Price exceeds current price band`, true
		case "18":
			return `Invalid price increment`, true
		case "99":
			return `Other`, true
		}

	case 104:
		switch value {
		case "A":
			return `All or None (AON)`, true
		case "B":
			return `Market On Close (MOC) (held to close)`, true
		case "C":
			return `At the close (around/not held to close)`, true
		case "D":
			return `VWAP (Volume Weighted Average Price)`, true
		case "I":
			return `In touch with`, true
		case "L":
			return `Limit`, true
		case "M":
			return `More Behind`, true
		case "O":
			return `At the Open`, true
		case "P":
			return `Taking a Position`, true
		case "Q":
			return `At the Market (previously called Current Quote)`, true
		case "R":
			return `Ready to Trade`, true
		case "S":
			return `Portfolio Shown`, true
		case "T":
			return `Through the Day`, true
		case "V":
			return `Versus`, true
		case "W":
			return `Indidcation - Working Away`, true
		case "X":
			return `Crossing Opportunity`, true
		case "Y":
			return `At the Midpoint`, true
		case "Z":
			return `Pre-open`, true
		}

	case 113:
		switch value {
		case "N":
			return `Indicates the party sending message will report trade`, true
		case "Y":
			return `Indicates the party receiving message must report trade`, true
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
			return `Do Not Execute Forex After Security Trade`, true
		case "Y":
			return `Execute Forex After Security Trade`, true
		}

	case 123:
		switch value {
		case "N":
			return `Sequence Reset, Ignore Msg Seq Num (N/A For FIXML - Not Used)`, true
		case "Y":
			return `Gap Fill Message, Msg Seq Num Field Valid`, true
		}

	case 127:
		switch value {
		case "A":
			return `Unknown Symbol`, true
		case "B":
			return `Wrong Side`, true
		case "C":
			return `Quantity Exceeds Order`, true
		case "D":
			return `No Matching Order`, true
		case "E":
			return `Price Exceeds Limit`, true
		case "F":
			return `Calculation Difference`, true
		case "Z":
			return `Other`, true
		}

	case 130:
		switch value {
		case "N":
			return `Not Natural`, true
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
		case "9":
			return `Consumption Tax`, true
		case "10":
			return `Per transaction`, true
		case "11":
			return `Conversion`, true
		case "12":
			return `Agent`, true
		case "13":
			return `Transfer Fee`, true
		case "14":
			return `Security Lending`, true
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
		case "3":
			return `Done for day`, true
		case "4":
			return `Canceled`, true
		case "5":
			return `Replaced`, true
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
			return `Restated (Execution Report sent unsolicited by sellside, with ExecRestatementReason (378) set)`, true
		case "E":
			return `Pending Replace (e.g. result of Order Cancel/Replace Request)`, true
		case "F":
			return `Trade (partial fill or fill)`, true
		case "G":
			return `Trade Correct`, true
		case "H":
			return `Trade Cancel`, true
		case "I":
			return `Order Status`, true
		case "J":
			return `Trade in a Clearing Hold`, true
		case "K":
			return `Trade has been released to Clearing`, true
		case "L":
			return `Triggered or Activated by System`, true
		}

	case 156:
		switch value {
		case "M":
			return `Multiply`, true
		case "D":
			return `Divide`, true
		}

	case 160:
		switch value {
		case "0":
			return `Default (Replaced)`, true
		case "1":
			return `Standing Instructions Provided`, true
		case "2":
			return `Specific Allocation Account Overriding (Replaced)`, true
		case "3":
			return `Specific Allocation Account Standing (Replaced)`, true
		case "4":
			return `Specific Order for a single account (for CIV)`, true
		case "5":
			return `Request reject`, true
		}

	case 163:
		switch value {
		case "N":
			return `New`, true
		case "C":
			return `Cancel`, true
		case "R":
			return `Replace`, true
		case "T":
			return `Restate`, true
		}

	case 165:
		switch value {
		case "1":
			return `Broker's Instructions`, true
		case "2":
			return `Institution's Instructions`, true
		case "3":
			return `Investor (e.g. CIV use)`, true
		}

	case 166:
		switch value {
		case "CED":
			return `CEDEL`, true
		case "DTC":
			return `Depository Trust Company`, true
		case "EUR":
			return `Euro clear`, true
		case "FED":
			return `Federal Book Entry`, true
		case "ISO_Country_Code":
			return `Local Market Settle Location`, true
		case "PNY":
			return `Physical`, true
		case "PTC":
			return `Participant Trust Company`, true
		}

	case 167:
		switch value {
		case "UST":
			return `US Treasury Note (Deprecated Value Use TNOTE)`, true
		case "USTB":
			return `US Treasury Bill (Deprecated Value Use TBILL)`, true
		case "EUSUPRA":
			return `Euro Supranational Coupons *`, true
		case "FAC":
			return `Federal Agency Coupon`, true
		case "FADN":
			return `Federal Agency Discount Note`, true
		case "PEF":
			return `Private Export Funding *`, true
		case "SUPRA":
			return `USD Supranational Coupons *`, true
		case "CORP":
			return `Corporate Bond`, true
		case "CPP":
			return `Corporate Private Placement`, true
		case "CB":
			return `Convertible Bond`, true
		case "DUAL":
			return `Dual Currency`, true
		case "EUCORP":
			return `Euro Corporate Bond`, true
		case "EUFRN":
			return `Euro Corporate Floating Rate Notes`, true
		case "FRN":
			return `US Corporate Floating Rate Notes`, true
		case "XLINKD":
			return `Indexed Linked`, true
		case "STRUCT":
			return `Structured Notes`, true
		case "YANK":
			return `Yankee Corporate Bond`, true
		case "FOR":
			return `Foreign Exchange Contract`, true
		case "CDS":
			return `Credit Default Swap`, true
		case "FUT":
			return `Future`, true
		case "OPT":
			return `Option`, true
		case "OOF":
			return `Options on Futures`, true
		case "OOP":
			return `Options on Physical - use not recommended`, true
		case "IRS":
			return `Interest Rate Swap`, true
		case "OOC":
			return `Options on Combo`, true
		case "CS":
			return `Common Stock`, true
		case "PS":
			return `Preferred Stock`, true
		case "REPO":
			return `Repurchase`, true
		case "FORWARD":
			return `Forward`, true
		case "BUYSELL":
			return `Buy Sellback`, true
		case "SECLOAN":
			return `Securities Loan`, true
		case "SECPLEDGE":
			return `Securities Pledge`, true
		case "BRADY":
			return `Brady Bond`, true
		case "CAN":
			return `Canadian Treasury Notes`, true
		case "CTB":
			return `Canadian Treasury Bills`, true
		case "EUSOV":
			return `Euro Sovereigns *`, true
		case "PROV":
			return `Canadian Provincial Bonds`, true
		case "TB":
			return `Treasury Bill - non US`, true
		case "TBOND":
			return `US Treasury Bond`, true
		case "TINT":
			return `Interest Strip From Any Bond Or Note`, true
		case "TBILL":
			return `US Treasury Bill`, true
		case "TIPS":
			return `Treasury Inflation Protected Securities`, true
		case "TCAL":
			return `Principal Strip Of A Callable Bond Or Note`, true
		case "TPRN":
			return `Principal Strip From A Non-Callable Bond Or Note`, true
		case "TNOTE":
			return `US Treasury Note`, true
		case "TERM":
			return `Term Loan`, true
		case "RVLV":
			return `Revolver Loan`, true
		case "RVLVTRM":
			return `Revolver/Term Loan`, true
		case "BRIDGE":
			return `Bridge Loan`, true
		case "LOFC":
			return `Letter Of Credit`, true
		case "SWING":
			return `Swing Line Facility`, true
		case "DINP":
			return `Debtor In Possession`, true
		case "DEFLTED":
			return `Defaulted`, true
		case "WITHDRN":
			return `Withdrawn`, true
		case "REPLACD":
			return `Replaced`, true
		case "MATURED":
			return `Matured`, true
		case "AMENDED":
			return `Amended & Restated`, true
		case "RETIRED":
			return `Retired`, true
		case "BA":
			return `Bankers Acceptance`, true
		case "BDN":
			return `Bank Depository Note`, true
		case "BN":
			return `Bank Notes`, true
		case "BOX":
			return `Bill Of Exchanges`, true
		case "CAMM":
			return `Canadian Money Markets`, true
		case "CD":
			return `Certificate Of Deposit`, true
		case "CL":
			return `Call Loans`, true
		case "CP":
			return `Commercial Paper`, true
		case "DN":
			return `Deposit Notes`, true
		case "EUCD":
			return `Euro Certificate Of Deposit`, true
		case "EUCP":
			return `Euro Commercial Paper`, true
		case "LQN":
			return `Liquidity Note`, true
		case "MTN":
			return `Medium Term Notes`, true
		case "ONITE":
			return `Overnight`, true
		case "PN":
			return `Promissory Note`, true
		case "STN":
			return `Short Term Loan Note`, true
		case "PZFJ":
			return `Plazos Fijos`, true
		case "SLQN":
			return `Secured Liquidity Note`, true
		case "TD":
			return `Time Deposit`, true
		case "TLQN":
			return `Term Liquidity Note`, true
		case "XCN":
			return `Extended Comm Note`, true
		case "YCD":
			return `Yankee Certificate Of Deposit`, true
		case "ABS":
			return `Asset-backed Securities`, true
		case "CMB":
			return `Canadian Mortgage Bonds`, true
		case "CMBS":
			return `Corp. Mortgage-backed Securities`, true
		case "CMO":
			return `Collateralized Mortgage Obligation`, true
		case "IET":
			return `IOETTE Mortgage`, true
		case "MBS":
			return `Mortgage-backed Securities`, true
		case "MIO":
			return `Mortgage Interest Only`, true
		case "MPO":
			return `Mortgage Principal Only`, true
		case "MPP":
			return `Mortgage Private Placement`, true
		case "MPT":
			return `Miscellaneous Pass-through`, true
		case "PFAND":
			return `Pfandbriefe *`, true
		case "TBA":
			return `To Be Announced`, true
		case "AN":
			return `Other Anticipation Notes (BAN, GAN, etc.)`, true
		case "COFO":
			return `Certificate Of Obligation`, true
		case "COFP":
			return `Certificate Of Participation`, true
		case "GO":
			return `General Obligation Bonds`, true
		case "MT":
			return `Mandatory Tender`, true
		case "RAN":
			return `Revenue Anticipation Note`, true
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
		case "TECP":
			return `Tax Exempt Commercial Paper`, true
		case "TMCP":
			return `Taxable Municipal CP`, true
		case "TRAN":
			return `Tax Revenue Anticipation Note`, true
		case "VRDN":
			return `Variable Rate Demand Note`, true
		case "WAR":
			return `Warrant`, true
		case "MF":
			return `Mutual Fund`, true
		case "MLEG":
			return `Multileg Instrument`, true
		case "NONE":
			return `No Security Type`, true
		case "?":
			return `Wildcard entry for use on Security Definition Request`, true
		case "CASH":
			return `Cash`, true
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
			return `A Global Custodian (StandInstDBName (70) must be provided)`, true
		case "4":
			return `AccountNet`, true
		}

	case 172:
		switch value {
		case "0":
			return `"Versus. Payment": Deliver (if Sell) or Receive (if Buy) vs. (Against) Payment`, true
		case "1":
			return `"Free": Deliver (if Sell) or Receive (if Buy) Free`, true
		case "2":
			return `Tri-Party`, true
		case "3":
			return `Hold In Custody`, true
		}

	case 197:
		switch value {
		case "0":
			return `FX Netting`, true
		case "1":
			return `FX Swap`, true
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
		case "1":
			return `CURVE`, true
		case "2":
			return `5YR`, true
		case "3":
			return `OLD5`, true
		case "4":
			return `10YR`, true
		case "5":
			return `OLD10`, true
		case "6":
			return `30YR`, true
		case "7":
			return `OLD30`, true
		case "8":
			return `3MOLIBOR`, true
		case "9":
			return `6MOLIBOR`, true
		}

	case 221:
		switch value {
		case "EONIA":
			return `EONIA`, true
		case "EUREPO":
			return `EUREPO`, true
		case "Euribor":
			return `Euribor`, true
		case "FutureSWAP":
			return `FutureSWAP`, true
		case "LIBID":
			return `LIBID`, true
		case "LIBOR":
			return `LIBOR (London Inter-Bank Offer)`, true
		case "MuniAAA":
			return `MuniAAA`, true
		case "OTHER":
			return `OTHER`, true
		case "Pfandbriefe":
			return `Pfandbriefe`, true
		case "SONIA":
			return `SONIA`, true
		case "SWAP":
			return `SWAP`, true
		case "Treasury":
			return `Treasury`, true
		}

	case 233:
		switch value {
		case "AMT":
			return `Alternative Minimum Tax (Y/N)`, true
		case "AUTOREINV":
			return `Auto Reinvestment at <rate> or better`, true
		case "BANKQUAL":
			return `Bank qualified (Y/N)`, true
		case "BGNCON":
			return `Bargain conditions (see StipulationValue (234) for values)`, true
		case "COUPON":
			return `Coupon range`, true
		case "CURRENCY":
			return `ISO Currency Code`, true
		case "CUSTOMDATE":
			return `Custom start/end date`, true
		case "GEOG":
			return `Geographics and % range (ex. 234=CA 0-80 [minimum of 80% California assets])`, true
		case "HAIRCUT":
			return `Valuation Discount`, true
		case "INSURED":
			return `Insured (Y/N)`, true
		case "ISSUE":
			return `Year Or Year/Month of Issue (ex. 234=2002/09)`, true
		case "ISSUER":
			return `Issuer's ticker`, true
		case "ISSUESIZE":
			return `issue size range`, true
		case "LOOKBACK":
			return `Lookback Days`, true
		case "LOT":
			return `Explicit lot identifier`, true
		case "LOTVAR":
			return `Lot Variance (value in percent maximum over- or under-allocation allowed)`, true
		case "MAT":
			return `Maturity Year And Month`, true
		case "MATURITY":
			return `Maturity range`, true
		case "MAXSUBS":
			return `Maximum substitutions (Repo)`, true
		case "MINDNOM":
			return `Minimum denomination`, true
		case "MININCR":
			return `Minimum increment`, true
		case "MINQTY":
			return `Minimum quantity`, true
		case "PAYFREQ":
			return `Payment frequency, calendar`, true
		case "PIECES":
			return `Number Of Pieces`, true
		case "PMAX":
			return `Pools Maximum`, true
		case "PPL":
			return `Pools per Lot`, true
		case "PPM":
			return `Pools per Million`, true
		case "PPT":
			return `Pools per Trade`, true
		case "PRICE":
			return `Price Range`, true
		case "PRICEFREQ":
			return `Pricing frequency`, true
		case "PROD":
			return `Production Year`, true
		case "PROTECT":
			return `Call protection`, true
		case "PURPOSE":
			return `Purpose`, true
		case "PXSOURCE":
			return `Benchmark price source`, true
		case "RATING":
			return `Rating source and range`, true
		case "REDEMPTION":
			return `Type Of Redemption - values are: NonCallable, Prefunded, EscrowedToMaturity, Putable, Convertible`, true
		case "RESTRICTED":
			return `Restricted (Y/N)`, true
		case "SECTOR":
			return `Market Sector`, true
		case "SECTYPE":
			return `Security Type included or excluded`, true
		case "STRUCT":
			return `Structure`, true
		case "SUBSFREQ":
			return `Substitutions frequency (Repo)`, true
		case "SUBSLEFT":
			return `Substitutions left (Repo)`, true
		case "TEXT":
			return `Freeform Text`, true
		case "TRDVAR":
			return `Trade Variance (value in percent maximum over- or under-allocation allowed)`, true
		case "WAC":
			return `Weighted Average Coupon - value in percent (exact or range) plus "Gross" or "Net" of servicing spread (the default) (ex. 234=6.5-Net [minimum of 6.5% net of servicing fee])`, true
		case "WAL":
			return `Weighted Average Life Coupon - value in percent (exact or range)`, true
		case "WALA":
			return `Weighted Average Loan Age - value in months (exact or range)`, true
		case "WAM":
			return `Weighted Average Maturity - value in months (exact or range)`, true
		case "WHOLE":
			return `Whole Pool (Y/N)`, true
		case "YIELD":
			return `Yield Range`, true
		case "AVFICO":
			return `Average FICO Score`, true
		case "AVSIZE":
			return `Average Loan Size`, true
		case "MAXBAL":
			return `Maximum Loan Balance`, true
		case "POOL":
			return `Pool Identifier`, true
		case "ROLLTYPE":
			return `Type of Roll trade`, true
		case "REFTRADE":
			return `reference to rolling or closing trade`, true
		case "REFPRIN":
			return `principal of rolling or closing trade`, true
		case "REFINT":
			return `interest of rolling or closing trade`, true
		case "AVAILQTY":
			return `Available offer quantity to be shown to the street`, true
		case "BROKERCREDIT":
			return `Broker's sales credit`, true
		case "INTERNALPX":
			return `Offer price to be shown to internal brokers`, true
		case "INTERNALQTY":
			return `Offer quantity to be shown to internal brokers`, true
		case "LEAVEQTY":
			return `The minimum residual offer quantity`, true
		case "MAXORDQTY":
			return `Maximum order size`, true
		case "ORDRINCR":
			return `Order quantity increment`, true
		case "PRIMARY":
			return `Primary or Secondary market indicator`, true
		case "SALESCREDITOVR":
			return `Broker sales credit override`, true
		case "TRADERCREDIT":
			return `Trader's credit`, true
		case "DISCOUNT":
			return `Discount Rate (when price is denominated in percent of par)`, true
		case "YTM":
			return `Yield to Maturity (when YieldType(235) and Yield(236) show a different yield)`, true
		case "ABS":
			return `Absolute Prepayment Speed`, true
		case "CPP":
			return `Constant Prepayment Penalty`, true
		case "CPR":
			return `Constant Prepayment Rate`, true
		case "CPY":
			return `Constant Prepayment Yield`, true
		case "HEP":
			return `final CPR of Home Equity Prepayment Curve`, true
		case "MHP":
			return `Percent of Manufactured Housing Prepayment Curve`, true
		case "MPR":
			return `Monthly Prepayment Rate`, true
		case "PPC":
			return `Percent of Prospectus Prepayment Curve`, true
		case "PSA":
			return `Percent of BMA Prepayment Curve`, true
		case "SMM":
			return `Single Monthly Mortality`, true
		}

	case 235:
		switch value {
		case "AFTERTAX":
			return `After Tax Yield (Municipals)`, true
		case "ANNUAL":
			return `Annual Yield`, true
		case "ATISSUE":
			return `Yield At Issue (Municipals)`, true
		case "AVGMATURITY":
			return `Yield To Avg Maturity`, true
		case "BOOK":
			return `Book Yield`, true
		case "CALL":
			return `Yield to Next Call`, true
		case "CHANGE":
			return `Yield Change Since Close`, true
		case "CLOSE":
			return `Closing Yield`, true
		case "COMPOUND":
			return `Compound Yield`, true
		case "CURRENT":
			return `Current Yield`, true
		case "GOVTEQUIV":
			return `Gvnt Equivalent Yield`, true
		case "GROSS":
			return `True Gross Yield`, true
		case "INFLATION":
			return `Yield with Inflation Assumption`, true
		case "INVERSEFLOATER":
			return `Inverse Floater Bond Yield`, true
		case "LASTCLOSE":
			return `Most Recent Closing Yield`, true
		case "LASTMONTH":
			return `Closing Yield Most Recent Month`, true
		case "LASTQUARTER":
			return `Closing Yield Most Recent Quarter`, true
		case "LASTYEAR":
			return `Closing Yield Most Recent Year`, true
		case "LONGAVGLIFE":
			return `Yield to Longest Average Life`, true
		case "MARK":
			return `Mark to Market Yield`, true
		case "MATURITY":
			return `Yield to Maturity`, true
		case "NEXTREFUND":
			return `Yield to Next Refund (Sinking Fund Bonds)`, true
		case "OPENAVG":
			return `Open Average Yield`, true
		case "PREVCLOSE":
			return `Previous Close Yield`, true
		case "PROCEEDS":
			return `Proceeds Yield`, true
		case "PUT":
			return `Yield to Next Put`, true
		case "SEMIANNUAL":
			return `Semi-annual Yield`, true
		case "SHORTAVGLIFE":
			return `Yield to Shortest Average Life`, true
		case "SIMPLE":
			return `Simple Yield`, true
		case "TAXEQUIV":
			return `Tax Equivalent Yield`, true
		case "TENDER":
			return `Yield to Tender Date`, true
		case "TRUE":
			return `True Yield`, true
		case "VALUE1_32":
			return `Yield Value Of 1/32`, true
		case "WORST":
			return `Yield To Worst`, true
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
		case "0":
			return `Snapshot`, true
		case "1":
			return `Snapshot + Updates (Subscribe)`, true
		case "2":
			return `Disable previous Snapshot + Update Request (Unsubscribe)`, true
		}

	case 265:
		switch value {
		case "0":
			return `Full refersh`, true
		case "1":
			return `Incremental refres`, true
		}

	case 266:
		switch value {
		case "Y":
			return `book entries to be aggregated`, true
		case "N":
			return `book entries should not be aggregated`, true
		}

	case 269:
		switch value {
		case "0":
			return `Bid`, true
		case "1":
			return `Offer`, true
		case "2":
			return `Trade`, true
		case "3":
			return `Index Value`, true
		case "4":
			return `Opening Price`, true
		case "5":
			return `Closing Price`, true
		case "6":
			return `Settlement Price`, true
		case "7":
			return `Trading Session High Price`, true
		case "8":
			return `Trading Session Low Price`, true
		case "9":
			return `Trading Session VWAP Price`, true
		case "A":
			return `Imbalance`, true
		case "B":
			return `Trade Volume`, true
		case "C":
			return `Open Interest`, true
		case "D":
			return `Composite Underlying Price`, true
		case "E":
			return `Simulated Sell Price`, true
		case "F":
			return `Simulated Buy Price`, true
		case "G":
			return `Margin Rate`, true
		case "H":
			return `Mid Price`, true
		case "J":
			return `Empty Book`, true
		case "K":
			return `Settle High Price`, true
		case "L":
			return `Settle Low Price`, true
		case "M":
			return `Prior Settle Price`, true
		case "N":
			return `Session High Bid`, true
		case "O":
			return `Session Low Offer`, true
		case "P":
			return `Early Prices`, true
		case "Q":
			return `Auction Clearing Price`, true
		case "S":
			return `Swap Value Factor (SVP) for swaps cleared through a central counterparty (CCP)`, true
		case "R":
			return `Daily value adjustment for long positions`, true
		case "T":
			return `Cumulative Value Adjustment for long positions`, true
		case "U":
			return `Daily Value Adjustment for Short Positions`, true
		case "V":
			return `Cumulative Value Adjustment for Short Positions`, true
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
		case "A":
			return `Open/Active`, true
		case "B":
			return `Closed/Inactive`, true
		case "C":
			return `Exchange Best`, true
		case "D":
			return `Consolidated Best`, true
		case "E":
			return `Locked`, true
		case "F":
			return `Crossed`, true
		case "G":
			return `Depth`, true
		case "H":
			return `Fast Trading`, true
		case "I":
			return `Non-Firm`, true
		case "L":
			return `Manual/Slow Quote`, true
		case "J":
			return `Outright Price`, true
		case "K":
			return `Implied Price`, true
		case "M":
			return `Depth on Offer`, true
		case "N":
			return `Depth on Bid`, true
		case "O":
			return `Closing`, true
		case "P":
			return `News Dissemination`, true
		case "Q":
			return `Trading Range`, true
		case "R":
			return `Order Influx`, true
		case "S":
			return `Due to Related`, true
		case "T":
			return `News Pending`, true
		case "U":
			return `Additional Info`, true
		case "V":
			return `Additional Info due to related`, true
		case "W":
			return `Resume`, true
		case "X":
			return `View of Common`, true
		case "Y":
			return `Volume Alert`, true
		case "Z":
			return `Order Imbalance`, true
		case "a":
			return `Equipment Changeover`, true
		case "b":
			return `No Open / No Resume`, true
		case "c":
			return `Regular ETH`, true
		case "d":
			return `Automatic Execution`, true
		case "e":
			return `Automatic Execution ETH`, true
		case "f ":
			return `Fast Market ETH`, true
		case "g":
			return `Inactive ETH`, true
		case "h":
			return `Rotation`, true
		case "i":
			return `Rotation ETH`, true
		case "j":
			return `Halt`, true
		case "k":
			return `Halt ETH`, true
		case "l":
			return `Due to News Dissemination`, true
		case "m":
			return `Due to News Pending`, true
		case "n":
			return `Trading Resume`, true
		case "o":
			return `Out of Sequence`, true
		case "p":
			return `Bid Specialist`, true
		case "q":
			return `Offer Specialist`, true
		case "r":
			return `Bid Offer Specialist`, true
		case "s":
			return `End of Day SAM`, true
		case "t":
			return `Forbidden SAM`, true
		case "u":
			return `Frozen SAM`, true
		case "v":
			return `PreOpening SAM`, true
		case "w":
			return `Opening SAM`, true
		case "x":
			return `Open SAM`, true
		case "y":
			return `Surveillance SAM`, true
		case "z":
			return `Suspended SAM`, true
		case "0":
			return `Reserved SAM`, true
		case "1":
			return `No Active SAM`, true
		case "2":
			return `Restricted`, true
		case "3":
			return `Rest of Book VWAP`, true
		case "4":
			return `Better Prices in Conditional Orders`, true
		case "5":
			return `Median Price`, true
		}

	case 277:
		switch value {
		case "A":
			return `Cash (only) Market`, true
		case "B":
			return `Average Price Trade`, true
		case "C":
			return `Cash Trade (same day clearing)`, true
		case "D":
			return `Next Day (only)Market`, true
		case "E":
			return `Opening/Reopening Trade Detail`, true
		case "F":
			return `Intraday Trade Detail`, true
		case "G":
			return `Rule 127 Trade (NYSE)`, true
		case "H":
			return `Rule 155 Trade (AMEX)`, true
		case "I":
			return `Sold Last (late reporting)`, true
		case "J":
			return `Next Day Trade (next day clearing)`, true
		case "K":
			return `Opened (late report of opened trade)`, true
		case "L":
			return `Seller`, true
		case "M":
			return `Sold (out of sequence)`, true
		case "N":
			return `Stopped Stock (guarantee of price but does not execute the order)`, true
		case "P":
			return `Imbalance More Buyers (cannot be used in combination with Q)`, true
		case "Q":
			return `Imbalance More Sellers (cannot be used in combination with P)`, true
		case "R":
			return `Opening Price`, true
		case "S":
			return `Bargain Condition (LSE)`, true
		case "T":
			return `Converted Price Indicator`, true
		case "U":
			return `Exchange Last`, true
		case "V":
			return `Final Price of Session`, true
		case "W":
			return `Ex-pit`, true
		case "X":
			return `Crossed`, true
		case "Y":
			return `Trades resulting from manual/slow quote`, true
		case "Z":
			return `Trades resulting from intermarket sweep`, true
		case "a":
			return `Volume Only`, true
		case "b":
			return `Direct Plus`, true
		case "c":
			return `Acquisition`, true
		case "d":
			return `Bunched`, true
		case "e":
			return `Distribution`, true
		case "f":
			return `Bunched Sale`, true
		case "g":
			return `Split Trade`, true
		case "h":
			return `Cancel Stopped`, true
		case "i":
			return `Cancel ETH`, true
		case "j":
			return `Cancel Stopped ETH`, true
		case "k":
			return `Out of Sequence ETH`, true
		case "l":
			return `Cancel Last ETH`, true
		case "m":
			return `Sold Last Sale ETH`, true
		case "n":
			return `Cancel Last`, true
		case "o":
			return `Sold Last Sale`, true
		case "p":
			return `Cancel Open`, true
		case "q":
			return `Cancel Open ETH`, true
		case "r":
			return `Opened Sale ETH`, true
		case "s":
			return `Cancel Only`, true
		case "t":
			return `Cancel Only ETH`, true
		case "u":
			return `Late Open ETH`, true
		case "v":
			return `Auto Execution ETH`, true
		case "w":
			return `Reopen`, true
		case "x":
			return `Reopen ETH`, true
		case "y":
			return `Adjusted`, true
		case "z":
			return `Adjusted ETH`, true
		case "AA":
			return `Spread`, true
		case "AB":
			return `Spread ETH`, true
		case "AC":
			return `Straddle`, true
		case "AD":
			return `Straddle ETH`, true
		case "AE":
			return `Stopped`, true
		case "AF":
			return `Stopped ETH`, true
		case "AG":
			return `Regular ETH`, true
		case "AH":
			return `Combo`, true
		case "AI":
			return `Combo ETH`, true
		case "AJ":
			return `Official Closing Price`, true
		case "AK":
			return `Prior Reference Price`, true
		case "0":
			return `Cancel`, true
		case "AL":
			return `Stopped Sold Last`, true
		case "AM":
			return `Stopped Out of Sequence`, true
		case "AN":
			return `Offical Closing Price (duplicate enumeration - use 'AJ' instead)`, true
		case "AO":
			return `Crossed (duplicate enumeration - use 'X' instead)`, true
		case "AP":
			return `Fast Market`, true
		case "AQ":
			return `Automatic Execution`, true
		case "AR":
			return `Form T`, true
		case "AS":
			return `Basket Index`, true
		case "AT":
			return `Burst Basket`, true
		case "AV":
			return `Outside Spread`, true
		case "1":
			return `Implied Trade`, true
		case "2":
			return `Marketplace entered trade`, true
		case "3":
			return `Mult Asset Class Multileg Trade`, true
		case "4":
			return `Multileg-to-Multileg Trade`, true
		}

	case 279:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Change`, true
		case "2":
			return `Delete`, true
		case "3":
			return `Delete Thru`, true
		case "4":
			return `Delete From`, true
		case "5":
			return `Overlay`, true
		}

	case 281:
		switch value {
		case "0":
			return `Unknown symbol`, true
		case "1":
			return `Duplicate MDReqID`, true
		case "2":
			return `Insufficient Bandwidth`, true
		case "3":
			return `Insufficient Permissions`, true
		case "4":
			return `Unsupported SubscriptionRequestType`, true
		case "5":
			return `Unsupported MarketDepth`, true
		case "6":
			return `Unsupported MDUpdateType`, true
		case "7":
			return `Unsupported AggregatedBook`, true
		case "8":
			return `Unsupported MDEntryType`, true
		case "9":
			return `Unsupported TradingSessionID`, true
		case "A":
			return `Unsupported Scope`, true
		case "B":
			return `Unsupported OpenCloseSettleFlag`, true
		case "C":
			return `Unsupported MDImplicitDelete`, true
		case "D":
			return `Insufficient credit`, true
		}

	case 285:
		switch value {
		case "0":
			return `Cancellation / Trade Bust`, true
		case "1":
			return `Error`, true
		}

	case 286:
		switch value {
		case "0":
			return `Daily Open / Close / Settlement entry`, true
		case "1":
			return `Session Open / Close / Settlement entry`, true
		case "2":
			return `Delivery Settlement entry`, true
		case "3":
			return `Expected entry`, true
		case "4":
			return `Entry from previous business day`, true
		case "5":
			return `Theoretical Price value`, true
		}

	case 291:
		switch value {
		case "1":
			return `Bankrupt`, true
		case "2":
			return `Pending delisting`, true
		case "3":
			return `Restricted`, true
		}

	case 292:
		switch value {
		case "A":
			return `Ex-Dividend`, true
		case "B":
			return `Ex-Distribution`, true
		case "C":
			return `Ex-Rights`, true
		case "D":
			return `New`, true
		case "E":
			return `Ex-Interest`, true
		case "F":
			return `Cash Dividend`, true
		case "G":
			return `Stock Dividend`, true
		case "H":
			return `Non-Integer Stock Split`, true
		case "I":
			return `Reverse Stock Split`, true
		case "J":
			return `Standard-Integer Stock Split`, true
		case "K":
			return `Position Consolidation`, true
		case "L":
			return `Liquidation Reorganization`, true
		case "M":
			return `Merger Reorganization`, true
		case "N":
			return `Rights Offering`, true
		case "O":
			return `Shareholder Meeting`, true
		case "P":
			return `Spinoff`, true
		case "Q":
			return `Tender Offer`, true
		case "R":
			return `Warrant`, true
		case "S":
			return `Special Action`, true
		case "T":
			return `Symbol Conversion`, true
		case "U":
			return `CUSIP / Name Change`, true
		case "V":
			return `Leap Rollover`, true
		case "W":
			return `Succession Event`, true
		}

	case 297:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Cancel for Symbol(s)`, true
		case "2":
			return `Canceled for Security Type(s)`, true
		case "3":
			return `Canceled for Underlying`, true
		case "4":
			return `Canceled All`, true
		case "5":
			return `Rejected`, true
		case "6":
			return `Removed from Market`, true
		case "7":
			return `Expired`, true
		case "8":
			return `Query`, true
		case "9":
			return `Quote Not Found`, true
		case "10":
			return `Pending`, true
		case "11":
			return `Pass`, true
		case "12":
			return `Locked Market Warning`, true
		case "13":
			return `Cross Market Warning`, true
		case "14":
			return `Canceled Due To Lock Market`, true
		case "15":
			return `Canceled Due To Cross Market`, true
		case "16":
			return `Active`, true
		case "17":
			return `Canceled`, true
		case "18":
			return `Unsolicited Quote Replenishment`, true
		case "19":
			return `Pending End Trade`, true
		case "20":
			return `Too Late to End`, true
		}

	case 298:
		switch value {
		case "1":
			return `Cancel for Symbol(s)`, true
		case "2":
			return `Cancel for Security Type(s)`, true
		case "3":
			return `Cancel for Underlying Symbol`, true
		case "4":
			return `Cancel All Quotes`, true
		case "5":
			return `Cancel quote specified in QuoteID`, true
		}

	case 300:
		switch value {
		case "1":
			return `Unknown Symbol (security)`, true
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
		case "9":
			return `Not authorized to quote security`, true
		case "10":
			return `Price exceeds current price band`, true
		case "11":
			return `Quote Locked - Unable to Update/Cancel`, true
		case "99":
			return `Other`, true
		}

	case 301:
		switch value {
		case "0":
			return `No Acknowledgement`, true
		case "1":
			return `Acknowledge only negative or erroneous quotes`, true
		case "2":
			return `Acknowledge each quote messages`, true
		case "3":
			return `Summary Acknowledgement`, true
		}

	case 303:
		switch value {
		case "1":
			return `Manual`, true
		case "2":
			return `Automatic`, true
		}

	case 321:
		switch value {
		case "0":
			return `Request Security identity and specifications`, true
		case "1":
			return `Request Security identity for the specifications provided (name of the security is not supplied)`, true
		case "2":
			return `Request List Security Types`, true
		case "3":
			return `Request List Securities (can be qualified with Symbol, SecurityType, TradingSessionID, SecurityExchange. If provided then only list Securities for the specific type.)`, true
		case "4":
			return `Symbol`, true
		case "5":
			return `SecurityType and or CFICode`, true
		case "6":
			return `Product`, true
		case "7":
			return `TradingSessionID`, true
		case "8":
			return `All Securities`, true
		case "9":
			return `MarketID or MarketID + MarketSegmentID`, true
		}

	case 323:
		switch value {
		case "1":
			return `Accept security proposal as-is`, true
		case "2":
			return `Accept security proposal with revisions as indicated in the message`, true
		case "3":
			return `List of security types returned per request`, true
		case "4":
			return `List of securities returned per request`, true
		case "5":
			return `Reject security proposal`, true
		case "6":
			return `Cannot match selection criteria`, true
		}

	case 325:
		switch value {
		case "N":
			return `Message is being sent as a result of a prior request`, true
		case "Y":
			return `Message is being secnt unsolicited`, true
		}

	case 326:
		switch value {
		case "1":
			return `Opening delay`, true
		case "2":
			return `Trading halt`, true
		case "3":
			return `Resume`, true
		case "4":
			return `No Open / No Resume`, true
		case "5":
			return `Price indication`, true
		case "6":
			return `Trading Range Indication`, true
		case "7":
			return `Market Imbalance Buy`, true
		case "8":
			return `Market Imbalance Sell`, true
		case "9":
			return `Market on Close Imbalance Buy`, true
		case "10":
			return `Market on Close Imbalance Sell`, true
		case "12":
			return `No Market Imbalance`, true
		case "13":
			return `No Market on Close Imbalance`, true
		case "14":
			return `ITS Pre-opening`, true
		case "15":
			return `New Price Indication`, true
		case "16":
			return `Trade Dissemination Time`, true
		case "17":
			return `Ready to trade (start of session)`, true
		case "18":
			return `Not available for trading (end of session)`, true
		case "19":
			return `Not traded on this market`, true
		case "20":
			return `Unknown or Invalid`, true
		case "21":
			return `Pre-open`, true
		case "22":
			return `Opening Rotation`, true
		case "23":
			return `Fast Market`, true
		case "24":
			return `Pre-Cross - system is in a pre-cross state allowing market to respond to either side of cross`, true
		case "25":
			return `Cross - system has crossed a percentage of the orders and allows market to respond prior to crossing remaining portion`, true
		}

	case 327:
		switch value {
		case "D":
			return `News Dissemination`, true
		case "E":
			return `Order Influx`, true
		case "I":
			return `Order Imbalance`, true
		case "M":
			return `Additional Information`, true
		case "P":
			return `New Pending`, true
		case "X":
			return `Equipment Changeover`, true
		}

	case 328:
		switch value {
		case "N":
			return `Halt was not related to a halt of the common stock`, true
		case "Y":
			return `Half was due to common stock being halted`, true
		}

	case 329:
		switch value {
		case "N":
			return `Halt was not related to a halt of the related security`, true
		case "Y":
			return `Halt was due to related security being halted`, true
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

	case 336:
		switch value {
		case "1":
			return `Day`, true
		case "2":
			return `HalfDay`, true
		case "3":
			return `Morning`, true
		case "4":
			return `Afternoon`, true
		case "5":
			return `Evening`, true
		case "6":
			return `After-hours`, true
		}

	case 338:
		switch value {
		case "1":
			return `Electronic`, true
		case "2":
			return `Open Outcry`, true
		case "3":
			return `Two Party`, true
		}

	case 339:
		switch value {
		case "1":
			return `Testing`, true
		case "2":
			return `Simulated`, true
		case "3":
			return `Production`, true
		}

	case 340:
		switch value {
		case "0":
			return `Unknown`, true
		case "1":
			return `Halted`, true
		case "2":
			return `Open`, true
		case "3":
			return `Closed`, true
		case "4":
			return `Pre-Open`, true
		case "5":
			return `Pre-Close`, true
		case "6":
			return `Request Rejected`, true
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

	case 374:
		switch value {
		case "C":
			return `Cancel`, true
		case "N":
			return `New`, true
		}

	case 377:
		switch value {
		case "N":
			return `Was not solicited`, true
		case "Y":
			return `Was solicited`, true
		}

	case 378:
		switch value {
		case "0":
			return `GT corporate action`, true
		case "1":
			return `GT renewal / restatement (no corporate action)`, true
		case "2":
			return `Verbal change`, true
		case "3":
			return `Repricing of order`, true
		case "4":
			return `Broker option`, true
		case "5":
			return `Partial decline of OrderQty (e.g. exchange initiated partial cancel)`, true
		case "6":
			return `Cancel on Trading Halt`, true
		case "7":
			return `Cancel on System Failure`, true
		case "8":
			return `Market (Exchange) option`, true
		case "9":
			return `Canceled, not best`, true
		case "10":
			return `Warehouse Recap`, true
		case "11":
			return `Peg Refresh`, true
		case "99":
			return `Other`, true
		}

	case 380:
		switch value {
		case "0":
			return `Other`, true
		case "1":
			return `Unknown ID`, true
		case "2":
			return `Unknown Security`, true
		case "3":
			return `Unknown Message Type`, true
		case "4":
			return `Application not available`, true
		case "5":
			return `Conditionally required field missing`, true
		case "6":
			return `Not Authorized`, true
		case "7":
			return `DeliverTo firm not available at this time`, true
		case "18":
			return `Invalid price increment`, true
		}

	case 385:
		switch value {
		case "R":
			return `Receive`, true
		case "S":
			return `Send`, true
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
		case "6":
			return `Related to VWAP`, true
		case "7":
			return `Average Price Guarantee`, true
		}

	case 394:
		switch value {
		case "1":
			return `"Non Disclosed" style (e.g. US/European)`, true
		case "2":
			return `"Disclosed" sytle (e.g. Japanese)`, true
		case "3":
			return `No bidding process`, true
		}

	case 399:
		switch value {
		case "1":
			return `Sector`, true
		case "2":
			return `Country`, true
		case "3":
			return `Index`, true
		}

	case 401:
		switch value {
		case "1":
			return `Side Value 1`, true
		case "2":
			return `Side Value 2`, true
		}

	case 409:
		switch value {
		case "1":
			return `5-day moving average`, true
		case "2":
			return `20-day moving average`, true
		case "3":
			return `Normal market size`, true
		case "4":
			return `Other`, true
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
		case "1":
			return `Buy-side explicitly requests status using Statue Request (default), the sell-side firm can, however, send a DONE status List STatus Response in an unsolicited fashion`, true
		case "2":
			return `Sell-side periodically sends status using List Status. Period optionally specified in ProgressPeriod.`, true
		case "3":
			return `Real-time execution reports (to be discourage)`, true
		}

	case 416:
		switch value {
		case "1":
			return `Net`, true
		case "2":
			return `Gross`, true
		}

	case 418:
		switch value {
		case "A":
			return `Agency`, true
		case "G":
			return `VWAP Guarantee`, true
		case "J":
			return `Guaranteed Close`, true
		case "R":
			return `Risk Trade`, true
		}

	case 419:
		switch value {
		case "2":
			return `Closing price at morning session`, true
		case "3":
			return `Closing price`, true
		case "4":
			return `Current price`, true
		case "5":
			return `SQ`, true
		case "6":
			return `VWAP through a day`, true
		case "7":
			return `VWAP through a morning session`, true
		case "8":
			return `VWAP through an afternoon session`, true
		case "9":
			return `VWAP through a day except "YORI" (an opening auction)`, true
		case "A":
			return `VWAP through a morning session except "YORI" (an opening auction)`, true
		case "B":
			return `VWAP through an afternoon session except "YORI" (an opening auction)`, true
		case "C":
			return `Strike`, true
		case "D":
			return `Open`, true
		case "Z":
			return `Others`, true
		}

	case 423:
		switch value {
		case "1":
			return `Percentage (i.e. percent of par) (often called "dollar price" for fixed income)`, true
		case "2":
			return `Per unit (i.e. per share or contract)`, true
		case "3":
			return `Fixed amount (absolute value)`, true
		case "4":
			return `Discount - percentage points below par`, true
		case "5":
			return `Premium - percentage points over par`, true
		case "6":
			return `Spread (basis points spread)`, true
		case "7":
			return `TED Price`, true
		case "8":
			return `TED Yield`, true
		case "9":
			return `Yield`, true
		case "10":
			return `Fixed cabinet trade price (primarily for listed futures and options)`, true
		case "11":
			return `Variable cabinet trade price (primarily for listed futures and options)`, true
		case "13":
			return `Product ticks in halfs`, true
		case "14":
			return `Product ticks in fourths`, true
		case "15":
			return `Product ticks in eights`, true
		case "16":
			return `Product ticks in sixteenths`, true
		case "17":
			return `Product ticks in thirty-seconds`, true
		case "18":
			return `Product ticks in sixty-forths`, true
		case "19":
			return `Product ticks in one-twenty-eights`, true
		}

	case 427:
		switch value {
		case "0":
			return `Book out all trades on day of execution`, true
		case "1":
			return `Accumulate exectuions until forder is filled or expires`, true
		case "2":
			return `Accumulate until verballly notified otherwise`, true
		}

	case 429:
		switch value {
		case "1":
			return `Ack`, true
		case "2":
			return `Response`, true
		case "3":
			return `Timed`, true
		case "4":
			return `Exec Started`, true
		case "5":
			return `All Done`, true
		case "6":
			return `Alert`, true
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
		case "1":
			return `In bidding process`, true
		case "2":
			return `Received for execution`, true
		case "3":
			return `Executing`, true
		case "4":
			return `Cancelling`, true
		case "5":
			return `Alert`, true
		case "6":
			return `All Done`, true
		case "7":
			return `Reject`, true
		}

	case 433:
		switch value {
		case "1":
			return `Immediate`, true
		case "2":
			return `Wait for Execut Instruction (i.e. a List Execut message or phone call before proceeding with execution of the list)`, true
		case "3":
			return `Exchange/switch CIV order - Sell driven`, true
		case "4":
			return `Exchange/switch CIV order - Buy driven, cash top-up (i.e. additional cash will be provided to fulfill the order)`, true
		case "5":
			return `Exchange/switch CIV order - Buy driven, cash withdraw (i.e. additional cash will not be provided to fulfill the order)`, true
		}

	case 434:
		switch value {
		case "1":
			return `Order cancel request`, true
		case "2":
			return `Order cancel/replace request`, true
		}

	case 442:
		switch value {
		case "1":
			return `Single security (defualt if not specified)`, true
		case "2":
			return `Individual leg of a multi=leg security`, true
		case "3":
			return `Multi-leg security`, true
		}

	case 447:
		switch value {
		case "6":
			return `UK National Insurance or Pension Number`, true
		case "7":
			return `US Social Security Number`, true
		case "8":
			return `US Employer or Tax ID Number`, true
		case "9":
			return `Australian Business Number`, true
		case "A":
			return `Australian Tax File Number`, true
		case "1":
			return `Korean Investor ID`, true
		case "2":
			return `Taiwanese Qualified Foreign Investor ID QFII/FID`, true
		case "3":
			return `Taiwanese Trading Acct`, true
		case "4":
			return `Malaysian Central Depository (MCD) number`, true
		case "5":
			return `Chinese Investor ID`, true
		case "I":
			return `Directed broker three character acronym as defined in ISITC "ETC Best Practice" guidelines document`, true
		case "B":
			return `BIC (Bank Identification Code - SWIFT managed) code (ISO9362 - See "Appendix 6-B")`, true
		case "C":
			return `Generally accepted market participant identifier (e.g. NASD mnemonic)`, true
		case "D":
			return `Proprietary / Custom code`, true
		case "E":
			return `ISO Country Code`, true
		case "F":
			return `Settlement Entity Location (note if Local Market Settlement use "E=ISO Country Code") (see "Appendix 6-G" for valid values)`, true
		case "G":
			return `MIC (ISO 10383 - Market Identificer Code) (See "Appendix 6-C")`, true
		case "H":
			return `CSD participant/member code (e.g.. Euroclear, DTC, CREST or Kassenverein number)`, true
		}

	case 452:
		switch value {
		case "1":
			return `Executing Firm (formerly FIX 4.2 ExecBroker)`, true
		case "2":
			return `Broker of Credit (formerly FIX 4.2 BrokerOfCredit)`, true
		case "3":
			return `Client ID (formerly FIX 4.2 ClientID)`, true
		case "4":
			return `Clearing Firm (formerly FIX 4.2 ClearingFirm)`, true
		case "5":
			return `Investor ID`, true
		case "6":
			return `Introducing Firm`, true
		case "7":
			return `Entering Firm`, true
		case "8":
			return `Locate / Lending Firm (for short-sales)`, true
		case "9":
			return `Fund Manager Client ID (for CIV)`, true
		case "10":
			return `Settlement Location (formerly FIX 4.2 SettlLocation)`, true
		case "11":
			return `Order Origination Trader (associated with Order Origination Firm - i.e. trader who initiates/submits the order)`, true
		case "12":
			return `Executing Trader (associated with Executing Firm - actually executes)`, true
		case "13":
			return `Order Origination Firm (e.g. buy-side firm)`, true
		case "14":
			return `Giveup Clearing Firm (firm to which trade is given up)`, true
		case "15":
			return `Correspondant Clearing Firm`, true
		case "16":
			return `Executing System`, true
		case "17":
			return `Contra Firm`, true
		case "18":
			return `Contra Clearing Firm`, true
		case "19":
			return `Sponsoring Firm`, true
		case "20":
			return `Underlying Contra Firm`, true
		case "21":
			return `Clearing Organization`, true
		case "22":
			return `Exchange`, true
		case "24":
			return `Customer Account`, true
		case "25":
			return `Correspondent Clearing Organization`, true
		case "26":
			return `Correspondent Broker`, true
		case "27":
			return `Buyer/Seller (Receiver/Deliverer)`, true
		case "28":
			return `Custodian`, true
		case "29":
			return `Intermediary`, true
		case "30":
			return `Agent`, true
		case "31":
			return `Sub-custodian`, true
		case "32":
			return `Beneficiary`, true
		case "33":
			return `Interested party`, true
		case "34":
			return `Regulatory body`, true
		case "35":
			return `Liquidity provider`, true
		case "36":
			return `Entering trader`, true
		case "37":
			return `Contra trader`, true
		case "38":
			return `Position account`, true
		case "39":
			return `Contra Investor ID`, true
		case "40":
			return `Transfer to Firm`, true
		case "41":
			return `Contra Position Account`, true
		case "42":
			return `Contra Exchange`, true
		case "43":
			return `Internal Carry Account`, true
		case "44":
			return `Order Entry Operator ID`, true
		case "45":
			return `Secondary Account Number`, true
		case "46":
			return `Foriegn Firm`, true
		case "47":
			return `Third Party Allocation Firm`, true
		case "48":
			return `Claiming Account`, true
		case "49":
			return `Asset Manager`, true
		case "50":
			return `Pledgor Account`, true
		case "51":
			return `Pledgee Account`, true
		case "52":
			return `Large Trader Reportable Account`, true
		case "53":
			return `Trader mnemonic`, true
		case "54":
			return `Sender Location`, true
		case "55":
			return `Session ID`, true
		case "56":
			return `Acceptable Counterparty`, true
		case "57":
			return `Unacceptable Counterparty`, true
		case "58":
			return `Entering Unit`, true
		case "59":
			return `Executing Unit`, true
		case "60":
			return `Introducing Broker`, true
		case "61":
			return `Quote originator`, true
		case "62":
			return `Report originator`, true
		case "63":
			return `Systematic internaliser (SI)`, true
		case "64":
			return `Multilateral Trading Facility (MTF)`, true
		case "65":
			return `Regulated Market (RM)`, true
		case "66":
			return `Market Maker`, true
		case "67":
			return `Investment Firm`, true
		case "68":
			return `Host Competent Authority (Host CA)`, true
		case "69":
			return `Home Competent Authority (Home CA)`, true
		case "70":
			return `Competent Authority of the most relevant market in terms of liquidity (CAL)`, true
		case "71":
			return `Competent Authority of the Transaction (Execution) Venue (CATV)`, true
		case "72":
			return `Reporting intermediary (medium/vendor via which report has been published)`, true
		case "73":
			return `Execution Venue`, true
		case "74":
			return `Market data entry originator`, true
		case "75":
			return `Location ID`, true
		case "76":
			return `Desk ID`, true
		case "77":
			return `Market data market`, true
		case "78":
			return `Allocation Entity`, true
		case "79":
			return `Prime Broker providing General Trade Services`, true
		case "80":
			return `Step-Out Firm (Prime Broker)`, true
		case "81":
			return `BrokerClearingID`, true
		}

	case 460:
		switch value {
		case "1":
			return `AGENCY`, true
		case "2":
			return `COMMODITY`, true
		case "3":
			return `CORPORATE`, true
		case "4":
			return `CURRENCY`, true
		case "5":
			return `EQUITY`, true
		case "6":
			return `GOVERNMENT`, true
		case "7":
			return `INDEX`, true
		case "8":
			return `LOAN`, true
		case "9":
			return `MONEYMARKET`, true
		case "10":
			return `MORTGAGE`, true
		case "11":
			return `MUNICIPAL`, true
		case "12":
			return `OTHER`, true
		case "13":
			return `FINANCING`, true
		}

	case 464:
		switch value {
		case "N":
			return `Fales (Production)`, true
		case "Y":
			return `True (Test)`, true
		}

	case 465:
		switch value {
		case "1":
			return `SHARES`, true
		case "2":
			return `BONDS`, true
		case "3":
			return `CURRENTFACE`, true
		case "4":
			return `ORIGINALFACE`, true
		case "5":
			return `CURRENCY`, true
		case "6":
			return `CONTRACTS`, true
		case "7":
			return `OTHER`, true
		case "8":
			return `PAR`, true
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

	case 477:
		switch value {
		case "1":
			return `CREST`, true
		case "2":
			return `NSCC`, true
		case "3":
			return `Euroclear`, true
		case "4":
			return `Clearstream`, true
		case "5":
			return `Cheque`, true
		case "6":
			return `Telegraphic Transfer`, true
		case "7":
			return `Fed Wire`, true
		case "8":
			return `Direct Credit (BECS, BACS)`, true
		case "9":
			return `ACH Credit`, true
		case "10":
			return `BPAY`, true
		case "11":
			return `High Value Clearing System HVACS`, true
		case "12":
			return `Reinvest In Fund`, true
		}

	case 480:
		switch value {
		case "Y":
			return `Yes`, true
		case "N":
			return `No - Execution Only`, true
		case "M":
			return `No - Waiver agreement`, true
		case "O":
			return `No - Institutional`, true
		}

	case 481:
		switch value {
		case "Y":
			return `Passed`, true
		case "N":
			return `Not Checked`, true
		case "1":
			return `Exempt - Below the Limit`, true
		case "2":
			return `Exempt - Client Money Type exemption`, true
		case "3":
			return `Exempt - Authorised Credit or financial institution`, true
		}

	case 484:
		switch value {
		case "B":
			return `Bid price`, true
		case "C":
			return `Creation price`, true
		case "D":
			return `Creation price plus adjustment percent`, true
		case "E":
			return `Creation price plus adjustment amount`, true
		case "O":
			return `Offer price`, true
		case "P":
			return `Offer price minus adjustment percent`, true
		case "Q":
			return `Offer price minus adjustment amount`, true
		case "S":
			return `Single price`, true
		}

	case 487:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Cancel`, true
		case "2":
			return `Replace`, true
		case "3":
			return `Release`, true
		case "4":
			return `Reverse`, true
		case "5":
			return `Cancel Due To Back Out of Trade`, true
		}

	case 492:
		switch value {
		case "1":
			return `CREST`, true
		case "2":
			return `NSCC`, true
		case "3":
			return `Euroclear`, true
		case "4":
			return `Clearstream`, true
		case "5":
			return `Cheque`, true
		case "6":
			return `Telegraphic Transfer`, true
		case "7":
			return `Fed Wire`, true
		case "8":
			return `Debit Card`, true
		case "9":
			return `Direct Debit (BECS)`, true
		case "10":
			return `Direct Credit (BECS)`, true
		case "11":
			return `Credit Card`, true
		case "12":
			return `ACH Debit`, true
		case "13":
			return `ACH Credit`, true
		case "14":
			return `BPAY`, true
		case "15":
			return `High Value Clearing System (HVACS)`, true
		}

	case 495:
		switch value {
		case "0":
			return `None/Not Applicable (default)`, true
		case "1":
			return `Maxi ISA (UK)`, true
		case "2":
			return `TESSA (UK)`, true
		case "3":
			return `Mini Cash ISA (UK)`, true
		case "4":
			return `Mini Stocks And Shares ISA (UK)`, true
		case "5":
			return `Mini Insurance ISA (UK)`, true
		case "6":
			return `Current Year Payment (US)`, true
		case "7":
			return `Prior Year Payment (US)`, true
		case "8":
			return `Asset Transfer (US)`, true
		case "9":
			return `Employee - prior year (US)`, true
		case "10":
			return `Employee - current year (US)`, true
		case "11":
			return `Employer - prior year (US)`, true
		case "12":
			return `Employer - current year (US)`, true
		case "13":
			return `Non-fund prototype IRA (US)`, true
		case "14":
			return `Non-fund qualified plan (US)`, true
		case "15":
			return `Defined contribution plan (US)`, true
		case "16":
			return `Individual Retirement Account (US)`, true
		case "17":
			return `Individual Retirement Account - Rollover (US)`, true
		case "18":
			return `KEOGH (US)`, true
		case "19":
			return `Profit Sharing Plan (US)`, true
		case "20":
			return `401(k) (US)`, true
		case "21":
			return `Self-directed IRA (US)`, true
		case "22":
			return `403(b) (US)`, true
		case "23":
			return `457 (US)`, true
		case "24":
			return `Roth IRA (Fund Prototype) (US)`, true
		case "25":
			return `Roth IRA (Non-prototype) (US)`, true
		case "26":
			return `Roth Conversion IRA (Fund Prototype) (US)`, true
		case "27":
			return `Roth Conversion IRA (Non-prototype) (US)`, true
		case "28":
			return `Education IRA (Fund Prototype) (US)`, true
		case "29":
			return `Education IRA (Non-prototype) (US)`, true
		case "999":
			return `Other`, true
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
			return `Accepted`, true
		case "R":
			return `Rejected`, true
		case "H":
			return `Held`, true
		case "N":
			return `Reminder - i.e. Registration Instructions are still outstanding`, true
		}

	case 507:
		switch value {
		case "1":
			return `Invalid/unacceptable Account Type`, true
		case "2":
			return `Invalid/unacceptable Tax Exempt Type`, true
		case "3":
			return `Invalid/unacceptable Ownership Type`, true
		case "4":
			return `Invalid/unacceptable No Reg Details`, true
		case "5":
			return `Invalid/unacceptable Reg Seq No`, true
		case "6":
			return `Invalid/unacceptable Reg Details`, true
		case "7":
			return `Invalid/unacceptable Mailing Details`, true
		case "8":
			return `Invalid/unacceptable Mailing Instructions`, true
		case "9":
			return `Invalid/unacceptable Investor ID`, true
		case "10":
			return `Invalid/unaceeptable Investor ID Source`, true
		case "11":
			return `Invalid/unacceptable Date Of Birth`, true
		case "12":
			return `Invalid/unacceptable Investor Country Of Residence`, true
		case "13":
			return `Invalid/unacceptable No Distrib Instns`, true
		case "14":
			return `Invalid/unacceptable Distrib Percentage`, true
		case "15":
			return `Invalid/unacceptable Distrib Payment Method`, true
		case "16":
			return `Invalid/unacceptable Cash Distrib Agent Acct Name`, true
		case "17":
			return `Invalid/unacceptable Cash Distrib Agent Code`, true
		case "18":
			return `Invalid/unacceptable Cash Distrib Agent Acct Num`, true
		case "99":
			return `Other`, true
		}

	case 514:
		switch value {
		case "0":
			return `New`, true
		case "2":
			return `Cancel`, true
		case "1":
			return `Replace`, true
		}

	case 517:
		switch value {
		case "J":
			return `Joint Investors`, true
		case "T":
			return `Tenants in Common`, true
		case "2":
			return `Joint Trustees`, true
		}

	case 519:
		switch value {
		case "1":
			return `Commission amount (actual)`, true
		case "2":
			return `Commission percent (actual)`, true
		case "3":
			return `Initial Charge Amount`, true
		case "4":
			return `Initial Charge Percent`, true
		case "5":
			return `Discount Amount`, true
		case "6":
			return `Discount Percent`, true
		case "7":
			return `Dilution Levy Amount`, true
		case "8":
			return `Dilution Levy Percent`, true
		case "9":
			return `Exit Charge Amount`, true
		case "10":
			return `Exit Charge Percent`, true
		case "11":
			return `Fund-Based Renewal Commission Percent (a.k.a. Trail commission)`, true
		case "12":
			return `Projected Fund Value (i.e. for investments intended to realise or exceed a specific future value)`, true
		case "13":
			return `Fund-Based Renewal Commission Amount (based on Order value)`, true
		case "14":
			return `Fund-Based Renewal Commission Amount (based on Projected Fund value)`, true
		case "15":
			return `Net Settlement Amount`, true
		}

	case 522:
		switch value {
		case "1":
			return `Individual Investor`, true
		case "2":
			return `Public Company`, true
		case "3":
			return `Private Company`, true
		case "4":
			return `Individual Trustee`, true
		case "5":
			return `Company Trustee`, true
		case "6":
			return `Pension Plan`, true
		case "7":
			return `Custodian Under Gifts to Minors Act`, true
		case "8":
			return `Trusts`, true
		case "9":
			return `Fiduciaries`, true
		case "10":
			return `Networking Sub-account`, true
		case "11":
			return `Non-profit organization`, true
		case "12":
			return `Corporate Body`, true
		case "13":
			return `Nominee`, true
		}

	case 528:
		switch value {
		case "A":
			return `Agency`, true
		case "G":
			return `Proprietary`, true
		case "I":
			return `Individual`, true
		case "P":
			return `Principal (Note for CMS purposes, "Principal" includes "Proprietary")`, true
		case "R":
			return `Riskless Principal`, true
		case "W":
			return `Agent for Other Member`, true
		}

	case 529:
		switch value {
		case "1":
			return `Program Trade`, true
		case "2":
			return `Index Arbitrage`, true
		case "3":
			return `Non-Index Arbitrage`, true
		case "4":
			return `Competing Market Maker`, true
		case "5":
			return `Acting as Market Maker or Specialist in the security`, true
		case "6":
			return `Acting as Market Maker of Specialist in the underlying security of a derivative seucirty`, true
		case "7":
			return `Foreign Entity (of foreign government or regulatory jurisdiction)`, true
		case "8":
			return `External Market Participant`, true
		case "9":
			return `Extneral Inter-connected Market Linkage`, true
		case "A":
			return `Riskless Arbitrage`, true
		case "B":
			return `Issuer Holding`, true
		case "C":
			return `Issue Price Stabilization`, true
		case "D":
			return `Non-algorithmic`, true
		case "E":
			return `Algorithmic`, true
		}

	case 530:
		switch value {
		case "1":
			return `Cancel orders for a security`, true
		case "2":
			return `Cancel orders for an underlying security`, true
		case "3":
			return `Cancel orders for a Product`, true
		case "4":
			return `Cancel orders for a CFICode`, true
		case "5":
			return `Cancel orders for a SecurityType`, true
		case "6":
			return `Cancel orders for a trading session`, true
		case "7":
			return `Cancel all orders`, true
		case "8":
			return `Cancel orders for a market`, true
		case "9":
			return `Cancel orders for a market segment`, true
		case "A":
			return `Cancel orders for a security group`, true
		}

	case 531:
		switch value {
		case "0":
			return `Cancel Request Rejected - See MassCancelRejectReason (532)`, true
		case "1":
			return `Cancel orders for a security`, true
		case "2":
			return `Cancel orders for an Underlying Security`, true
		case "3":
			return `Cancel orders for a Product`, true
		case "4":
			return `Cancel orders for a CFICode`, true
		case "5":
			return `Cancel orders for a SecurityType`, true
		case "6":
			return `Cancel orders for a trading session`, true
		case "7":
			return `Cancel All Orders`, true
		case "8":
			return `Cancel orders for a market`, true
		case "9":
			return `Cancel orders for a market segment`, true
		case "A":
			return `Cancel orders for a security group`, true
		}

	case 532:
		switch value {
		case "0":
			return `Mass Cancel Not Supported`, true
		case "1":
			return `Invalid or Unknown Security`, true
		case "2":
			return `Invalid or Unkown Underlying security`, true
		case "3":
			return `Invalid or Unknown Product`, true
		case "4":
			return `Invalid or Unknown CFICode`, true
		case "5":
			return `Invalid or Unknown SecurityType`, true
		case "6":
			return `Invalid or Unknown Trading Session`, true
		case "7":
			return `Invalid or unknown Market`, true
		case "8":
			return `Invalid or unkown Market Segment`, true
		case "9":
			return `Invalid or unknown Security Group`, true
		case "99":
			return `Other`, true
		}

	case 537:
		switch value {
		case "0":
			return `Indicative`, true
		case "1":
			return `Tradeable`, true
		case "2":
			return `Restricted Tradeable`, true
		case "3":
			return `Counter (tradeable)`, true
		}

	case 544:
		switch value {
		case "1":
			return `Cash`, true
		case "2":
			return `Margin Open`, true
		case "3":
			return `Margin Close`, true
		}

	case 546:
		switch value {
		case "1":
			return `Local Market (Exchange, ECN, ATS)`, true
		case "2":
			return `National`, true
		case "3":
			return `Global`, true
		}

	case 547:
		switch value {
		case "N":
			return `Server must send an explicit delete for bids or offers falling outside the requested MarketDepth of the request`, true
		case "Y":
			return `Client has responsibility for implicitly deleting bids or offers falling outside the MarketDepth of the request`, true
		}

	case 549:
		switch value {
		case "1":
			return `Cross AON - cross tade which is executed complete or not. Both sides are treated in the same manner. This is equivalent to an "All or None".`, true
		case "2":
			return `Cross IOC - cross trade which is executed partially and the rest is cancelled. One side is fully executed, the other side is partially executed with the remainder being cancelled. This is equivalent to an IOC on the other side. Note: CrossPrioritization (550) field may be used to indicate which side should fully execute in this scenario.`, true
		case "3":
			return `Cross One Side - cross trade which is partially executed with the unfilled portions remaining active.. One side of the corss is fully executed (as denoted by the CrossPrioritization (550) field), but the unfilled portion remains active.`, true
		case "4":
			return `Cross Same Price - cross trade is executed with existing orders with the same price. In this case other orders exist with the same price, the quantity of the Cross is executed against the existing orders and quotes, the remainder of the corss is executed against the other side of the cross. The two sides potentially have different quantities.`, true
		}

	case 550:
		switch value {
		case "0":
			return `None`, true
		case "1":
			return `Buy side is prioritized`, true
		case "2":
			return `Sell side is prioritized`, true
		}

	case 552:
		switch value {
		case "1":
			return `One Side`, true
		case "2":
			return `Both Sides`, true
		}

	case 559:
		switch value {
		case "0":
			return `Symbol`, true
		case "1":
			return `SecurityType and/or CFICode`, true
		case "2":
			return `Product`, true
		case "3":
			return `TradingSessionID`, true
		case "4":
			return `All Securities`, true
		case "5":
			return `MarketID or MarketID + MarketSegmentID`, true
		}

	case 560:
		switch value {
		case "0":
			return `Valid request`, true
		case "1":
			return `Invalid or unsupported request`, true
		case "2":
			return `No instruments found that match selection criteria`, true
		case "3":
			return `Not authorized to retrieve instrument data`, true
		case "4":
			return `Instrument data temporarily unavailable`, true
		case "5":
			return `Request for instrument data not supported`, true
		}

	case 563:
		switch value {
		case "0":
			return `Report by mulitleg security only (do not report legs)`, true
		case "1":
			return `Report by multileg security and by instrument legs belonging to the multileg security`, true
		case "2":
			return `Report by instrument legs belonging to the multileg security only (do not report status of multileg security)`, true
		}

	case 567:
		switch value {
		case "1":
			return `Unknown or invalid TradingSessionID`, true
		case "99":
			return `Other`, true
		}

	case 569:
		switch value {
		case "0":
			return `All Trades`, true
		case "1":
			return `Matched trades matching criteria provided on request (Parties, ExecID, TradeID, OrderID, Instrument, InputSource, etc.)`, true
		case "2":
			return `Unmatched trades that match criteria`, true
		case "3":
			return `Unreported trades that match criteria`, true
		case "4":
			return `Advisories that match criteria`, true
		}

	case 570:
		switch value {
		case "N":
			return `Not reported to counterparty`, true
		case "Y":
			return `Perviously reported to counterparty`, true
		}

	case 573:
		switch value {
		case "0":
			return `Compared, matched or affirmed`, true
		case "1":
			return `Uncompared, unmatched, or unaffired`, true
		case "2":
			return `Advisory or alert`, true
		}

	case 574:
		switch value {
		case "1":
			return `One-Party Trade Report (privately negotiated trade)`, true
		case "2":
			return `Two-Party Trade Report (privately negotiated trade)`, true
		case "3":
			return `Confirmed Trade Report (reporting from recognized markets)`, true
		case "4":
			return `Auto-match`, true
		case "5":
			return `Cross Auction`, true
		case "6":
			return `Counter-Order Selection`, true
		case "7":
			return `Call Auction`, true
		case "8":
			return `Issuing/Buy Back Auction`, true
		case "M3":
			return `ACT Accepted Trade`, true
		case "M4":
			return `ACT Default Trade`, true
		case "M5":
			return `ACT Default After M2`, true
		case "M6":
			return `ACT M6 Match`, true
		case "A1":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus four badges and execution time (within two-minute window)`, true
		case "A2":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator, plus four badges`, true
		case "A3":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator, plus two badges and execution time (within two-minute window)`, true
		case "A4":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator, plus two badges`, true
		case "A5":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, TradeType, and Special Trade Indicator plus execution time (within two-minute window)`, true
		case "AQ":
			return `Compared records resulting from stamped advisories or specialist accepts/pair-offs`, true
		case "S1":
			return `Summarized match using A1 exact match criteria except quantity is summaried`, true
		case "S2":
			return `Summarized match using A2 exact match criteria except quantity is summarized`, true
		case "S3":
			return `Summarized match using A3 exact match criteria except quantity is summarized`, true
		case "S4":
			return `Summarized match using A4 exact match criteria except quantity is summarized`, true
		case "S5":
			return `Summarized match using A5 exact match criteria except quantity is summarized`, true
		case "M1":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator minus badges And times: ACT M1 match`, true
		case "M2":
			return `Summarized match minus badges and times: ACT M2 Match`, true
		case "MT":
			return `OCS Locked In: Non-ACT`, true
		}

	case 575:
		switch value {
		case "N":
			return `Treat as round lot (default)`, true
		case "Y":
			return `Treat as odd lot`, true
		}

	case 577:
		switch value {
		case "0":
			return `Process normally`, true
		case "1":
			return `Exclude from all netting`, true
		case "2":
			return `Bilateral netting only`, true
		case "3":
			return `Ex clearing`, true
		case "4":
			return `Special trade`, true
		case "5":
			return `Multilateral netting`, true
		case "6":
			return `Clear against central counterparty`, true
		case "7":
			return `Exclude from central counterparty`, true
		case "8":
			return `Manual mode (pre-posting and/or pre-giveup)`, true
		case "9":
			return `Automatic posting mode (trade posting to the position account number specified)`, true
		case "10":
			return `Automatic give-up mode (trade give-up to the give-up destination number specified)`, true
		case "11":
			return `Qualified Service Representative QSR`, true
		case "12":
			return `Customer trade`, true
		case "13":
			return `Self clearing`, true
		}

	case 581:
		switch value {
		case "1":
			return `Account is carried on customer side of the books`, true
		case "2":
			return `Account is carried on non-customer side of books`, true
		case "3":
			return `House Trader`, true
		case "4":
			return `Floor Trader`, true
		case "6":
			return `Account is carried on non-customer side of books and is cross margined`, true
		case "7":
			return `Account is house trader and is cross margined`, true
		case "8":
			return `Joint back office account (JBO)`, true
		}

	case 582:
		switch value {
		case "1":
			return `Member trading for their own account`, true
		case "2":
			return `Clearing Firm trading for its proprietary account`, true
		case "3":
			return `Member trading for another member`, true
		case "4":
			return `All other`, true
		}

	case 585:
		switch value {
		case "1":
			return `Status for orders for a Security`, true
		case "2":
			return `Status for orders for an Underlying Security`, true
		case "3":
			return `Status for orders for a Product`, true
		case "4":
			return `Status for orders for a CFICode`, true
		case "5":
			return `Status for orders for a SecurityType`, true
		case "6":
			return `Status for orders for a trading session`, true
		case "7":
			return `Status for all orders`, true
		case "8":
			return `Status for orders for a PartyID`, true
		}

	case 589:
		switch value {
		case "0":
			return `Can trigger booking without reference to the order initiator ("auto")`, true
		case "1":
			return `Speak with order initiator before booking ("speak first")`, true
		case "2":
			return `Accumulate`, true
		}

	case 590:
		switch value {
		case "0":
			return `Each partial execution is a bookable unit`, true
		case "1":
			return `Aggregate partial executions on this order, and book one trade per order`, true
		case "2":
			return `Aggregate executions for this symbol, side, and settlement date`, true
		}

	case 591:
		switch value {
		case "0":
			return `Pro-rata`, true
		case "1":
			return `Do not pro-rata - discuss first`, true
		}

	case 625:
		switch value {
		case "1":
			return `Pre-Trading`, true
		case "2":
			return `Opening or opening auction`, true
		case "3":
			return `(Continuous) Trading`, true
		case "4":
			return `Closing or closing auction`, true
		case "5":
			return `Post-Trading`, true
		case "6":
			return `Intraday Auction`, true
		case "7":
			return `Quiescent`, true
		}

	case 626:
		switch value {
		case "1":
			return `Calculated (includes MiscFees and NetMoney)`, true
		case "2":
			return `Preliminary (without MiscFees and NetMoney)`, true
		case "3":
			return `Sellside Calculated Using Preliminary (includes MiscFees and NetMoney) (Replaced)`, true
		case "4":
			return `Sellside Calculated Without Preliminary (sent unsolicited by sellside, includes MiscFees and NetMoney) (Replaced)`, true
		case "5":
			return `Ready-To-Book - Single Order`, true
		case "6":
			return `Buyside Ready-To-Book - Combined Set of Orders (Replaced)`, true
		case "7":
			return `Warehouse Instruction`, true
		case "8":
			return `Request to Intermediary`, true
		case "9":
			return `Accept`, true
		case "10":
			return `Reject`, true
		case "11":
			return `Accept Pending`, true
		case "12":
			return `Incomplete Group`, true
		case "13":
			return `Complete Group`, true
		case "14":
			return `Reversal Pending`, true
		}

	case 635:
		switch value {
		case "1":
			return `1st year delegate trading for own account`, true
		case "2":
			return `2nd year delegate trading for own account`, true
		case "3":
			return `3rd year delegate trading for own account`, true
		case "4":
			return `4th year delegate trading for own account`, true
		case "5":
			return `5th year delegate trading for own account`, true
		case "9":
			return `6th year delegate trading for own account`, true
		case "B":
			return `CBOE Member`, true
		case "C":
			return `Non-member and Customer`, true
		case "E":
			return `Equity Member and Clearing Member`, true
		case "F":
			return `Full and Associate Member trading for own account and as floor brokers`, true
		case "H":
			return `106.H and 106.J firms`, true
		case "I":
			return `GIM, IDEM and COM Membership Interest Holders`, true
		case "L":
			return `Lessee 106.F Employees`, true
		case "M":
			return `All other ownership types`, true
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
			return `Priority unchanged`, true
		case "1":
			return `Lost Priority as result of order change`, true
		}

	case 650:
		switch value {
		case "N":
			return `Does not consitute a Legal Confirm`, true
		case "Y":
			return `Legal Confirm`, true
		}

	case 653:
		switch value {
		case "0":
			return `Pending Approval`, true
		case "1":
			return `Approved (Accepted)`, true
		case "2":
			return `Rejected`, true
		case "3":
			return `Unauthorized Request`, true
		case "4":
			return `Invalid Definition Request`, true
		}

	case 658:
		switch value {
		case "1":
			return `Unknown Symbol (Security)`, true
		case "2":
			return `Exchange (Security) Closed`, true
		case "3":
			return `Quote Request Exceeds Limit`, true
		case "4":
			return `Too Late to enter`, true
		case "5":
			return `Invalid Price`, true
		case "6":
			return `Not Authorized To Request Quote`, true
		case "7":
			return `No Match For Inquiry`, true
		case "8":
			return `No Market For Instrument`, true
		case "9":
			return `No Inventory`, true
		case "10":
			return `Pass`, true
		case "11":
			return `Insufficient credit`, true
		case "99":
			return `Other`, true
		}

	case 660:
		switch value {
		case "1":
			return `BIC`, true
		case "2":
			return `SID Code`, true
		case "3":
			return `TFM (GSPTA)`, true
		case "4":
			return `OMGEO (Alert ID)`, true
		case "5":
			return `DTCC Code`, true
		case "99":
			return `Other (custom or proprietary)`, true
		}

	case 665:
		switch value {
		case "1":
			return `Received`, true
		case "2":
			return `Mismatched Account`, true
		case "3":
			return `Missing Settlement Instructions`, true
		case "4":
			return `Confirmed`, true
		case "5":
			return `Request Rejected`, true
		}

	case 666:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Replace`, true
		case "2":
			return `Cancel`, true
		}

	case 668:
		switch value {
		case "1":
			return `Book Entry (default)`, true
		case "2":
			return `Bearer`, true
		}

	case 690:
		switch value {
		case "1":
			return `Par For Par`, true
		case "2":
			return `Modified Duration`, true
		case "4":
			return `Risk`, true
		case "5":
			return `Proceeds`, true
		}

	case 692:
		switch value {
		case "1":
			return `Percent (percent of par)`, true
		case "2":
			return `Per Share (e.g. cents per share)`, true
		case "3":
			return `Fixed Amount (absolute value)`, true
		case "4":
			return `Discount - percentage points below par`, true
		case "5":
			return `Premium - percentage points over par`, true
		case "6":
			return `Spread - basis points relative to benchmark`, true
		case "7":
			return `TED Price`, true
		case "8":
			return `TED Yield`, true
		case "9":
			return `Yield Spread (swaps)`, true
		case "10":
			return `Yield`, true
		}

	case 694:
		switch value {
		case "1":
			return `Hit/Lift`, true
		case "2":
			return `Counter`, true
		case "3":
			return `Expired`, true
		case "4":
			return `Cover`, true
		case "5":
			return `Done Away`, true
		case "6":
			return `Pass`, true
		case "7":
			return `End Trade`, true
		case "8":
			return `Timed Out`, true
		}

	case 703:
		switch value {
		case "ALC":
			return `Allocation Trade Qty`, true
		case "AS":
			return `Option Assignment`, true
		case "ASF":
			return `As-of Trade Qty`, true
		case "DLV":
			return `Delivery Qty`, true
		case "ETR":
			return `Electronic Trade Qty`, true
		case "EX":
			return `Option Exercise Qty`, true
		case "FIN":
			return `End-of-Day Qty`, true
		case "IAS":
			return `Intra-spread Qty`, true
		case "IES":
			return `Inter-spread Qty`, true
		case "PA":
			return `Adjustment Qty`, true
		case "PIT":
			return `Pit Trade Qty`, true
		case "SOD":
			return `Start-of-Day Qty`, true
		case "SPL":
			return `Integral Split`, true
		case "TA":
			return `Transaction from Assignment`, true
		case "TOT":
			return `Total Transaction Qty`, true
		case "TQ":
			return `Transaction Quantity`, true
		case "TRF":
			return `Transfer Trade Qty`, true
		case "TX":
			return `Transaction from Exercise`, true
		case "XM":
			return `Cross Margin Qty`, true
		case "RCV":
			return `Receive Quantity`, true
		case "CAA":
			return `Corporate Action Adjustment`, true
		case "DN":
			return `Delivery Notice Qty`, true
		case "EP":
			return `Exchange for Physical Qty`, true
		case "PNTN":
			return `Privately negotiated Trade Qty (Non-regulated)`, true
		}

	case 706:
		switch value {
		case "0":
			return `Submitted`, true
		case "1":
			return `Accepted`, true
		case "2":
			return `Rejected`, true
		}

	case 707:
		switch value {
		case "CASH":
			return `Cash Amount (Corporate Event)`, true
		case "CRES":
			return `Cash Residual Amount`, true
		case "FMTM":
			return `Final Mark-to-Market Amount`, true
		case "IMTM":
			return `Incremental Mark-to-Market Amount`, true
		case "PREM":
			return `Premium Amount`, true
		case "SMTM":
			return `Start-of-Day Mark-to-Market Amount`, true
		case "TVAR":
			return `Trade Variation Amount`, true
		case "VADJ":
			return `Value Adjusted Amount`, true
		case "SETL":
			return `Settlement Value`, true
		}

	case 709:
		switch value {
		case "1":
			return `Exercise`, true
		case "2":
			return `Do Not Exercise`, true
		case "3":
			return `Position Adjustment`, true
		case "4":
			return `Position Change Submission/Margin Disposition`, true
		case "5":
			return `Pledge`, true
		case "6":
			return `Large Trader Submission`, true
		}

	case 712:
		switch value {
		case "1":
			return `New - used to increment the overall transaction quantity`, true
		case "2":
			return `Replace - used to override the overall transaction quantity or specifi add messages based on the reference ID`, true
		case "3":
			return `Cancel - used to remove the overall transaction or specific add messages based on reference ID`, true
		case "4":
			return `Reverse - used to completelly back-out the transaction such that the transaction never existed`, true
		}

	case 716:
		switch value {
		case "ITD":
			return `Intraday`, true
		case "RTH":
			return `Regular Trading Hours`, true
		case "ETH":
			return `Electronic Trading Hours`, true
		case "EOD":
			return `End Of Day`, true
		}

	case 718:
		switch value {
		case "0":
			return `Process Request As Margin Disposition`, true
		case "1":
			return `Delta Plus`, true
		case "2":
			return `Delta Minus`, true
		case "3":
			return `Final`, true
		}

	case 722:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Accepted With Warnings`, true
		case "2":
			return `Rejected`, true
		case "3":
			return `Completed`, true
		case "4":
			return `Completed With Warnings`, true
		}

	case 723:
		switch value {
		case "0":
			return `Successful Completion - no warnings or errors`, true
		case "1":
			return `Rejected`, true
		case "99":
			return `Other`, true
		}

	case 724:
		switch value {
		case "0":
			return `Positions`, true
		case "1":
			return `Trades`, true
		case "2":
			return `Exercises`, true
		case "3":
			return `Assignments`, true
		case "4":
			return `Settlement Activity`, true
		case "5":
			return `Backout Message`, true
		}

	case 725:
		switch value {
		case "0":
			return `Inband - transport the request was sent over (default)`, true
		case "1":
			return `Out of Band - pre-arranged out-of-band delivery mechanizm (i.e. FTP, HTTP, NDM, etc.) between counterparties. Details specified via ResponseDestination (726).`, true
		}

	case 728:
		switch value {
		case "0":
			return `Valid request`, true
		case "1":
			return `Invalid or unsupported request`, true
		case "2":
			return `No positions found that match criteria`, true
		case "3":
			return `Not authorized to request positions`, true
		case "4":
			return `Request for position not supported`, true
		case "99":
			return `Other (use Text (58) in conjunction with this code for an explaination)`, true
		}

	case 729:
		switch value {
		case "0":
			return `Completed`, true
		case "1":
			return `Completed With Warnings`, true
		case "2":
			return `Rejected`, true
		}

	case 731:
		switch value {
		case "1":
			return `Final`, true
		case "2":
			return `Theoretical`, true
		}

	case 744:
		switch value {
		case "P":
			return `Pro-rata`, true
		case "R":
			return `Random`, true
		}

	case 747:
		switch value {
		case "A":
			return `Automatic`, true
		case "M":
			return `Manual`, true
		}

	case 749:
		switch value {
		case "0":
			return `Successful (default)`, true
		case "1":
			return `Invalid or unknown instrument`, true
		case "2":
			return `Invalid type of trade requested`, true
		case "3":
			return `Invalid parties`, true
		case "4":
			return `Invalid transport type requested`, true
		case "5":
			return `Invalid destination requested`, true
		case "8":
			return `TradeRequestType not supported`, true
		case "9":
			return `Not authorized`, true
		case "99":
			return `Other`, true
		}

	case 750:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Completed`, true
		case "2":
			return `Rejected`, true
		}

	case 751:
		switch value {
		case "0":
			return `Successful (default)`, true
		case "1":
			return `Invalid party onformation`, true
		case "2":
			return `Unknown instrument`, true
		case "3":
			return `Unauthorized to report trades`, true
		case "4":
			return `Invalid trade type`, true
		case "99":
			return `Other`, true
		}

	case 752:
		switch value {
		case "1":
			return `Single Security (default if not specified)`, true
		case "2":
			return `Individual leg of a multileg security`, true
		case "3":
			return `Multileg Security`, true
		}

	case 770:
		switch value {
		case "1":
			return `Execution Time`, true
		case "2":
			return `Time In`, true
		case "3":
			return `Time Out`, true
		case "4":
			return `Broker Receipt`, true
		case "5":
			return `Broker Execution`, true
		case "6":
			return `Desk Receipt`, true
		}

	case 773:
		switch value {
		case "1":
			return `Status`, true
		case "2":
			return `Confirmation`, true
		case "3":
			return `Confirmation Request Rejected (reason can be stated in Text (58) field)`, true
		}

	case 774:
		switch value {
		case "1":
			return `Mismatched account`, true
		case "2":
			return `Missing settlement instructions`, true
		case "99":
			return `Other`, true
		}

	case 775:
		switch value {
		case "0":
			return `Regular booking`, true
		case "1":
			return `CFD (Contract for difference)`, true
		case "2":
			return `Total Return Swap`, true
		}

	case 780:
		switch value {
		case "0":
			return `Use default instructions`, true
		case "1":
			return `Derive from parameters provided`, true
		case "2":
			return `Full details provided`, true
		case "3":
			return `SSI DB IDs provided`, true
		case "4":
			return `Phone for instructions`, true
		}

	case 787:
		switch value {
		case "C":
			return `Cash`, true
		case "S":
			return `Securities`, true
		}

	case 788:
		switch value {
		case "1":
			return `Overnight`, true
		case "2":
			return `Term`, true
		case "3":
			return `Flexible`, true
		case "4":
			return `Open`, true
		}

	case 792:
		switch value {
		case "0":
			return `Unable to process request`, true
		case "1":
			return `Unknown account`, true
		case "2":
			return `No matching settlement instructions found`, true
		case "99":
			return `Other`, true
		}

	case 794:
		switch value {
		case "2":
			return `Preliminary Request to Intermediary`, true
		case "3":
			return `Sellside Calculated Using Preliminary (includes MiscFees and NetMoney)`, true
		case "4":
			return `Sellside Calculated Without Preliminary (sent unsolicited by sellside, includes MiscFees and NetMoney)`, true
		case "5":
			return `Warehouse Recap`, true
		case "8":
			return `Request to Intermediary`, true
		case "9":
			return `Accept`, true
		case "10":
			return `Reject`, true
		case "11":
			return `Accept Pending`, true
		case "12":
			return `Complete`, true
		case "14":
			return `Reverse Pending`, true
		}

	case 796:
		switch value {
		case "1":
			return `Original details incomplete/incorrect`, true
		case "2":
			return `Change in underlying order details`, true
		case "99":
			return `Other`, true
		}

	case 798:
		switch value {
		case "1":
			return `Account is carried pn customer side of books`, true
		case "2":
			return `Account is carried on non-customer side of books`, true
		case "3":
			return `House trader`, true
		case "4":
			return `Floor trader`, true
		case "6":
			return `Account is carried on non-customer side of books and is cross margined`, true
		case "7":
			return `Account is house trader and is cross margined`, true
		case "8":
			return `Joint back office account (JBO)`, true
		}

	case 803:
		switch value {
		case "1":
			return `Firm`, true
		case "2":
			return `Person`, true
		case "3":
			return `System`, true
		case "4":
			return `Application`, true
		case "5":
			return `Full legal name of firm`, true
		case "6":
			return `Postal address`, true
		case "7":
			return `Phone number`, true
		case "8":
			return `Email address`, true
		case "9":
			return `Contact name`, true
		case "10":
			return `Securities account number (for settlement instructions)`, true
		case "11":
			return `Registration number (for settlement instructions and confirmations)`, true
		case "12":
			return `Registered address (for confirmation purposes)`, true
		case "13":
			return `Regulatory status (for confirmation purposes)`, true
		case "14":
			return `Registration name (for settlement instructions)`, true
		case "15":
			return `Cash account number (for settlement instructions)`, true
		case "16":
			return `BIC`, true
		case "17":
			return `CSD participant member code`, true
		case "18":
			return `Registered address`, true
		case "19":
			return `Fund account name`, true
		case "20":
			return `Telex number`, true
		case "21":
			return `Fax number`, true
		case "22":
			return `Securities account name`, true
		case "23":
			return `Cash account name`, true
		case "24":
			return `Department`, true
		case "25":
			return `Location desk`, true
		case "26":
			return `Position account type`, true
		case "27":
			return `Security locate ID`, true
		case "28":
			return `Market maker`, true
		case "29":
			return `Eligible counterparty`, true
		case "30":
			return `Professional client`, true
		case "31":
			return `Location`, true
		case "32":
			return `Execution venue`, true
		case "33":
			return `Currency delivery identifier`, true
		}

	case 808:
		switch value {
		case "1":
			return `Pending Accept`, true
		case "2":
			return `Pending Release`, true
		case "3":
			return `Pending Reversal`, true
		case "4":
			return `Accept`, true
		case "5":
			return `Block Level Reject`, true
		case "6":
			return `Account Level Reject`, true
		}

	case 814:
		switch value {
		case "0":
			return `No Action Taken`, true
		case "1":
			return `Queue Flushed`, true
		case "2":
			return `Overlay Last`, true
		case "3":
			return `End Session`, true
		}

	case 815:
		switch value {
		case "0":
			return `No Action Taken`, true
		case "1":
			return `Queue Flushed`, true
		case "2":
			return `Overlay Last`, true
		case "3":
			return `End Session`, true
		}

	case 819:
		switch value {
		case "0":
			return `No Average Pricing`, true
		case "1":
			return `Trade is part of an average price group identified by the TradeLinkID (820)`, true
		case "2":
			return `Last trade is the average price group identified by the TradeLinkID (820)`, true
		}

	case 826:
		switch value {
		case "0":
			return `Allocation not required`, true
		case "1":
			return `Allocation required (give-up trade) allocation information not provided (incomplete)`, true
		case "2":
			return `Use allocation provided with the trade`, true
		case "3":
			return `Allocation give-up executor`, true
		case "4":
			return `Allocation from executor`, true
		case "5":
			return `Allocation to claim account`, true
		}

	case 827:
		switch value {
		case "0":
			return `Expire on trading session close (default)`, true
		case "1":
			return `Expire on trading session open`, true
		case "2":
			return `Trading eligibility expiration specified in the date and time fields [EventDate(866) and EventTime(1145)] associated with EventType(865)=7(Last Eligible Trade Date)`, true
		}

	case 828:
		switch value {
		case "0":
			return `Regular Trade`, true
		case "1":
			return `Block Trade`, true
		case "2":
			return `EFP (Exchange for physical)`, true
		case "3":
			return `Transfer`, true
		case "4":
			return `Late Trade`, true
		case "5":
			return `T Trade`, true
		case "6":
			return `Weighted Average Price Trade`, true
		case "7":
			return `Bunched Trade`, true
		case "8":
			return `Late Bunched Trade`, true
		case "9":
			return `Prior Reference Price Trade`, true
		case "10":
			return `After Hours Trade`, true
		case "11":
			return `Exchange for Risk (EFR)`, true
		case "12":
			return `Exchange for Swap (EFS )`, true
		case "13":
			return `Exchange of Futures for (in Market) Futures (EFM ) (e,g, full sized for mini)`, true
		case "14":
			return `Exchange of Options for Options (EOO)`, true
		case "15":
			return `Trading at Settlement`, true
		case "16":
			return `All or None`, true
		case "17":
			return `Futures Large Order Execution`, true
		case "18":
			return `Exchange of Futures for Futures (external market) (EFF)`, true
		case "19":
			return `Option Interim Trade`, true
		case "20":
			return `Option Cabinet Trade`, true
		case "22":
			return `Privately Negotiated Trades`, true
		case "23":
			return `Substitution of Futures for Forwards`, true
		case "48":
			return `Non-standard settlement`, true
		case "49":
			return `Derivative Related Transaction`, true
		case "50":
			return `Portfolio Trade`, true
		case "51":
			return `Volume Weighted Average Trade`, true
		case "52":
			return `Exchange Granted Trade`, true
		case "53":
			return `Repurchase Agreement`, true
		case "54":
			return `OTC`, true
		case "55":
			return `Exchange Basis Facility (EBF)`, true
		case "24":
			return `Error trade`, true
		case "25":
			return `Special cum dividend (CD)`, true
		case "26":
			return `Special ex dividend (XD)`, true
		case "27":
			return `Special cum coupon (CC)`, true
		case "28":
			return `Special ex coupon (XC)`, true
		case "29":
			return `Cash settlement (CS)`, true
		case "30":
			return `Special price (usually net- or all-in price) (SP)`, true
		case "31":
			return `Guaranteed delivery (GD)`, true
		case "32":
			return `Special cum rights (CR)`, true
		case "33":
			return `Special ex rights (XR)`, true
		case "34":
			return `Special cum capital repayments (CP)`, true
		case "35":
			return `Special ex capital repayments (XP)`, true
		case "36":
			return `Special cum bonus (CB)`, true
		case "37":
			return `Special ex bonus (XB)`, true
		case "38":
			return `Block trade (same as large trade)`, true
		case "39":
			return `Worked principal trade (UK-specific)`, true
		case "40":
			return `Block Trades - after market`, true
		case "41":
			return `Name change`, true
		case "42":
			return `Portfolio transfer`, true
		case "43":
			return `Prorogation buy - Euronext Paris only. Is used to defer settlement under French SRD (deferred settlement system) . Trades must be reported as crosses at zero price`, true
		case "44":
			return `Prorogation sell - see prorogation buy`, true
		case "45":
			return `Option exercise`, true
		case "46":
			return `Delta neutral transaction`, true
		case "47":
			return `Financing transaction (includes repo and stock lending)`, true
		}

	case 829:
		switch value {
		case "0":
			return `CMTA`, true
		case "1":
			return `Internal transfer or adjustment`, true
		case "2":
			return `External transfer or transfer of account`, true
		case "3":
			return `Reject for submitting side`, true
		case "4":
			return `Advisory for contra side`, true
		case "5":
			return `Offset due to an allocation`, true
		case "6":
			return `Onset dut to an allocation`, true
		case "7":
			return `Differential spread`, true
		case "8":
			return `Implied spread leg executed against an outright`, true
		case "9":
			return `Transaction from exercise`, true
		case "10":
			return `Transaction from assignment`, true
		case "11":
			return `ACATS`, true
		case "33":
			return `Off Hours Trade`, true
		case "34":
			return `On Hours Trade`, true
		case "35":
			return `OTC Quote`, true
		case "36":
			return `Converted SWAP`, true
		case "14":
			return `AI (Automated input facility disabled in response to an exchange request.)`, true
		case "15":
			return `B (Transaction between two member firms where neither member firm is registered as a market maker in the security in question and neither is a designated fund manager. Also used by broker dealers when dealing with another broker which is not a member firm. Non-order book securities only.)`, true
		case "16":
			return `K (Transaction using block trade facility.)`, true
		case "17":
			return `LC (Correction submitted more than three days after publication of the original trade report.)`, true
		case "18":
			return `M (Transaction, other than a transaction resulting from a stock swap or stock switch, between two market makers registered in that security including IDB or a public display system trades. Non-order book securities only.)`, true
		case "19":
			return `N (Non-protected portfolio transaction or a fully disclosed portfolio transaction)`, true
		case "20":
			return `NM ( i) transaction where Exchange has granted permission for non-publication
ii)IDB is reporting as seller
iii) submitting a transaction report to the Exchange, where the transaction report is not also a trade report.)`, true
		case "21":
			return `NR (Non-risk transaction in a SEATS security other than an AIM security)`, true
		case "22":
			return `P (Protected portfolio transaction or a worked principal agreement to effect a portfolio transaction which includes order book securities)`, true
		case "23":
			return `PA (Protected transaction notification)`, true
		case "24":
			return `PC (Contra trade for transaction which took place on a previous day and which was automatically executed on the Exchange trading system)`, true
		case "25":
			return `PN (Worked principal notification for a portfolio transaction which includes order book securities)`, true
		case "26":
			return `R ( (i) riskless principal transaction between non-members where the buying and selling transactions are executed at different prices or on different terms (requires a trade report with trade type indicator R for each transaction)
(ii) market maker is reporting all the legs of a riskless principal transaction where the buying and selling transactions are executed at different prices (requires a trade report with trade type indicator R for each transaction)or
(iii) market maker is reporting the onward leg of a riskless principal transaction where the legs are executed at different prices, and another market maker has submitted a trade report using trade type indicator M for the first leg (this requires a single trade report with trade type indicator R).)`, true
		case "27":
			return `RO (Transaction which resulted from the exercise of a traditional option or a stock-settled covered warrant)`, true
		case "28":
			return `RT (Risk transaction in a SEATS security, (excluding AIM security) reported by a market maker registered in that security)`, true
		case "29":
			return `SW (Transactions resulting from stock swap or a stock switch (one report is required for each line of stock))`, true
		case "30":
			return `T (If reporting a single protected transaction)`, true
		case "31":
			return `WN (Worked principal notification for a single order book security)`, true
		case "32":
			return `WT (Worked principal transaction (other than a portfolio transaction))`, true
		case "37":
			return `Crossed Trade (X)`, true
		case "38":
			return `Interim Protected Trade (I)`, true
		case "39":
			return `Large in Scale (L)`, true
		}

	case 835:
		switch value {
		case "0":
			return `Floating (default)`, true
		case "1":
			return `Fixed`, true
		}

	case 836:
		switch value {
		case "0":
			return `Price (default)`, true
		case "1":
			return `Basis Points`, true
		case "2":
			return `Ticks`, true
		case "3":
			return `Price Tier / Level`, true
		}

	case 837:
		switch value {
		case "0":
			return `Or better (default) - price improvement allowed`, true
		case "1":
			return `Strict - limit is a strict limit`, true
		case "2":
			return `Or worse - for a buy the peg limit is a minimum and for a sell the peg limit is a maximum (for use for orders which have a price range)`, true
		}

	case 838:
		switch value {
		case "1":
			return `More aggressive - on a buy order round the price up to the nearest tick; on a sell order round down to the nearest tick`, true
		case "2":
			return `More passive - on a buy order round down to the nearest tick; on a sell order round up to the nearest tick`, true
		}

	case 840:
		switch value {
		case "1":
			return `Local (Exchange, ECN, ATS)`, true
		case "2":
			return `National`, true
		case "3":
			return `Global`, true
		case "4":
			return `National excluding local`, true
		}

	case 841:
		switch value {
		case "0":
			return `Floating (default)`, true
		case "1":
			return `Fixed`, true
		}

	case 842:
		switch value {
		case "0":
			return `Price (default)`, true
		case "1":
			return `Basis Points`, true
		case "2":
			return `Ticks`, true
		case "3":
			return `Price Tier / Level`, true
		}

	case 843:
		switch value {
		case "0":
			return `Or better (default) - price improvement allowed`, true
		case "1":
			return `Strict - limit is a strict limit`, true
		case "2":
			return `Or worse - for a buy the discretion price is a minimum and for a sell the discretion price is a maximum (for use for orders which have a price range)`, true
		}

	case 844:
		switch value {
		case "1":
			return `More aggressive - on a buy order round the price up to the nearest tick; on a sell round down to the nearest tick`, true
		case "2":
			return `More passive - on a buy order round down to the nearest tick; on a sell order round up to the nearest tick`, true
		}

	case 846:
		switch value {
		case "1":
			return `Local (Exchange, ECN, ATS)`, true
		case "2":
			return `National`, true
		case "3":
			return `Global`, true
		case "4":
			return `National excluding local`, true
		}

	case 847:
		switch value {
		case "1":
			return `VWAP`, true
		case "2":
			return `Participate (i.e. aim to be x percent of the market volume)`, true
		case "3":
			return `Mininize market impact`, true
		}

	case 851:
		switch value {
		case "1":
			return `Added Liquidity`, true
		case "2":
			return `Removed Liquidity`, true
		case "3":
			return `Liquidity Routed Out`, true
		case "4":
			return `Auction`, true
		}

	case 852:
		switch value {
		case "N":
			return `Do Not Report Trade`, true
		case "Y":
			return `Report Trade`, true
		}

	case 853:
		switch value {
		case "0":
			return `Dealer Sold Short`, true
		case "1":
			return `Dealer Sold Short Exempt`, true
		case "2":
			return `Selling Customer Sold Short`, true
		case "3":
			return `Selling Customer Sold Short Exempt`, true
		case "4":
			return `Qualified Service Representative (QSR) or Automatic Give-up (AGU) Contra Side Sold Short`, true
		case "5":
			return `QSR or AGU Contra Side Sold Short Exempt`, true
		}

	case 854:
		switch value {
		case "0":
			return `Units (shares, par, currency)`, true
		case "1":
			return `Contracts (if used - must specify ContractMultiplier (tag 231))`, true
		case "2":
			return `Units of Measure per Time Unit (if used - must specify UnitofMeasure (tag 996) and TimeUnit (tag 997))`, true
		}

	case 856:
		switch value {
		case "0":
			return `Submit`, true
		case "1":
			return `Alleged`, true
		case "2":
			return `Accept`, true
		case "3":
			return `Decline`, true
		case "4":
			return `Addendum`, true
		case "5":
			return `No/Was`, true
		case "6":
			return `Trade Report Cancel`, true
		case "7":
			return `(Locked-In) Trade Break`, true
		case "8":
			return `Defaulted`, true
		case "9":
			return `Invalid CMTA`, true
		case "10":
			return `Pended`, true
		case "11":
			return `Alleged New`, true
		case "12":
			return `Alleged Addendum`, true
		case "13":
			return `Alleged No/Was`, true
		case "14":
			return `Alleged Trade Report Cancel`, true
		case "15":
			return `Alleged (Locked-In) Trade Break`, true
		}

	case 857:
		switch value {
		case "0":
			return `Not Specified`, true
		case "1":
			return `Explicit List Provided`, true
		}

	case 865:
		switch value {
		case "1":
			return `Put`, true
		case "2":
			return `Call`, true
		case "3":
			return `Tender`, true
		case "4":
			return `Sinking Fund Call`, true
		case "5":
			return `Activation`, true
		case "6":
			return `Inactiviation`, true
		case "7":
			return `Last Eligible Trade Date`, true
		case "8":
			return `Swap Start Date`, true
		case "9":
			return `Swap End Date`, true
		case "10":
			return `Swap Roll Date`, true
		case "11":
			return `Swap Next Start Date`, true
		case "12":
			return `Swap Next Roll Date`, true
		case "13":
			return `First Delivery Date`, true
		case "14":
			return `Last Delivery Date`, true
		case "15":
			return `Initial Inventory Due Date`, true
		case "16":
			return `Final Inventory Due Date`, true
		case "17":
			return `First Intent Date`, true
		case "18":
			return `Last Intent Date`, true
		case "19":
			return `Position Removal Date`, true
		case "99":
			return `Other`, true
		}

	case 871:
		switch value {
		case "1":
			return `Flat (securities pay interest on a current basis but are traded without interest)`, true
		case "2":
			return `Zero coupon`, true
		case "3":
			return `Interest bearing (for Euro commercial paper when not issued at discount)`, true
		case "4":
			return `No periodic payments`, true
		case "5":
			return `Variable rate`, true
		case "6":
			return `Less fee for put`, true
		case "7":
			return `Stepped coupon`, true
		case "8":
			return `Coupon period (if not semi-annual). Supply redemption date in the InstrAttribValue (872) field.`, true
		case "9":
			return `When [and if] issued`, true
		case "10":
			return `Original issue discount`, true
		case "11":
			return `Callable, puttable`, true
		case "12":
			return `Escrowed to Maturity`, true
		case "13":
			return `Escrowed to redemption date - callable. Supply redemption date in the InstrAttribValue (872) field`, true
		case "14":
			return `Pre-refunded`, true
		case "15":
			return `In default`, true
		case "16":
			return `Unrated`, true
		case "17":
			return `Taxable`, true
		case "18":
			return `Indexed`, true
		case "19":
			return `Subject To Alternative Minimum Tax`, true
		case "20":
			return `Original issue discount price. Supply price in the InstrAttribValue (872) field`, true
		case "21":
			return `Callable below maturity value`, true
		case "22":
			return `Callable without notice by mail to holder unless registered`, true
		case "23":
			return `Price tick rules for security.`, true
		case "24":
			return `Trade type eligibility details for security.`, true
		case "25":
			return `Instrument Denominator`, true
		case "26":
			return `Instrument Numerator`, true
		case "27":
			return `Instrument Price Precision`, true
		case "28":
			return `Instrument Strike Price`, true
		case "29":
			return `Tradeable Indicator`, true
		case "99":
			return `Text. Supply the text of the attribute or disclaimer in the InstrAttribValue (872) field.`, true
		}

	case 875:
		switch value {
		case "1":
			return `3(a)(3)`, true
		case "2":
			return `4(2)`, true
		case "99":
			return `Other`, true
		}

	case 891:
		switch value {
		case "0":
			return `Absolute`, true
		case "1":
			return `Per Unit`, true
		case "2":
			return `Percentage`, true
		}

	case 893:
		switch value {
		case "N":
			return `Not Last Message`, true
		case "Y":
			return `Last Message`, true
		}

	case 895:
		switch value {
		case "0":
			return `Initial`, true
		case "1":
			return `Scheduled`, true
		case "2":
			return `Time Warning`, true
		case "3":
			return `Margin Deficiency`, true
		case "4":
			return `Margin Excess`, true
		case "5":
			return `Forward Collateral Demand`, true
		case "6":
			return `Event of default`, true
		case "7":
			return `Adverse tax event`, true
		}

	case 896:
		switch value {
		case "0":
			return `Trade Date`, true
		case "1":
			return `GC Instrument`, true
		case "2":
			return `Collateral Instrument`, true
		case "3":
			return `Substitution Eligible`, true
		case "4":
			return `Not Assigned`, true
		case "5":
			return `Partially Assigned`, true
		case "6":
			return `Fully Assigned`, true
		case "7":
			return `Outstanding Trades (Today < end date)`, true
		}

	case 903:
		switch value {
		case "0":
			return `New`, true
		case "1":
			return `Replace`, true
		case "2":
			return `Cancel`, true
		case "3":
			return `Release`, true
		case "4":
			return `Reverse`, true
		}

	case 905:
		switch value {
		case "0":
			return `Received`, true
		case "1":
			return `Accepted`, true
		case "2":
			return `Declined`, true
		case "3":
			return `Rejected`, true
		}

	case 906:
		switch value {
		case "0":
			return `Unknown deal (order / trade)`, true
		case "1":
			return `Unknown or invalid instrument`, true
		case "2":
			return `Unauthorized transaction`, true
		case "3":
			return `Insufficient collateral`, true
		case "4":
			return `Invalid type of collateral`, true
		case "5":
			return `Excessive substitution`, true
		case "99":
			return `Other`, true
		}

	case 910:
		switch value {
		case "0":
			return `Unassigned`, true
		case "1":
			return `Partially Assigned`, true
		case "2":
			return `Assignment Proposed`, true
		case "3":
			return `Assigned (Accepted)`, true
		case "4":
			return `Challenged`, true
		}

	case 912:
		switch value {
		case "N":
			return `Not last message`, true
		case "Y":
			return `Last message`, true
		}

	case 919:
		switch value {
		case "0":
			return `"Versus Payment": Deliver (if sell) or Receive (if buy) vs. (against) Payment`, true
		case "1":
			return `"Free": Deliver (if sell) or Receive (if buy) Free`, true
		case "2":
			return `Tri-Party`, true
		case "3":
			return `Hold In Custody`, true
		}

	case 924:
		switch value {
		case "1":
			return `Log On User`, true
		case "2":
			return `Log Off User`, true
		case "3":
			return `Change Password For User`, true
		case "4":
			return `Request Individual User Status`, true
		}

	case 926:
		switch value {
		case "1":
			return `Logged In`, true
		case "2":
			return `Not Logged In`, true
		case "3":
			return `User Not Recognised`, true
		case "4":
			return `Password Incorrect`, true
		case "5":
			return `Password Changed`, true
		case "6":
			return `Other`, true
		case "7":
			return `Forced user logout by Exchange`, true
		case "8":
			return `Session shutdown warning`, true
		}

	case 928:
		switch value {
		case "1":
			return `Connected`, true
		case "2":
			return `Not Connected - down expected up`, true
		case "3":
			return `Not Connected - down expected down`, true
		case "4":
			return `In Process`, true
		}

	case 935:
		switch value {
		case "1":
			return `Snapshot`, true
		case "2":
			return `Subscribe`, true
		case "4":
			return `Stop Subscribing`, true
		case "8":
			return `Level of Detail, then NoCompID's becomes required`, true
		}

	case 937:
		switch value {
		case "1":
			return `Full`, true
		case "2":
			return `Incremental Update`, true
		}

	case 939:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Rejected`, true
		case "3":
			return `Accepted with errors`, true
		}

	case 940:
		switch value {
		case "1":
			return `Received`, true
		case "2":
			return `Confirm rejected, i.e. not affirmed`, true
		case "3":
			return `Affirmed`, true
		}

	case 944:
		switch value {
		case "0":
			return `Retain`, true
		case "1":
			return `Add`, true
		case "2":
			return `Remove`, true
		}

	case 945:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Accepted With Warnings`, true
		case "2":
			return `Completed`, true
		case "3":
			return `Completed With Warnings`, true
		case "4":
			return `Rejected`, true
		}

	case 946:
		switch value {
		case "0":
			return `Successful (default)`, true
		case "1":
			return `Invalid or unknown instrument`, true
		case "2":
			return `Invalid or unknown collateral type`, true
		case "3":
			return `Invalid Parties`, true
		case "4":
			return `Invalid Transport Type requested`, true
		case "5":
			return `Invalid Destination requested`, true
		case "6":
			return `No collateral found for the trade specified`, true
		case "7":
			return `No collateral found for the order specified`, true
		case "8":
			return `Collateral inquiry type not supported`, true
		case "9":
			return `Unauthorized for collateral inquiry`, true
		case "99":
			return `Other (further information in Text (58) field)`, true
		}

	case 959:
		switch value {
		case "1":
			return `Int`, true
		case "2":
			return `Length`, true
		case "3":
			return `NumInGroup`, true
		case "4":
			return `SeqNum`, true
		case "5":
			return `TagNum`, true
		case "6":
			return `Float`, true
		case "7":
			return `Qty`, true
		case "8":
			return `Price`, true
		case "9":
			return `PriceOffset`, true
		case "10":
			return `Amt`, true
		case "11":
			return `Percentage`, true
		case "12":
			return `Char`, true
		case "13":
			return `Boolean`, true
		case "14":
			return `String`, true
		case "15":
			return `MultipleCharValue`, true
		case "16":
			return `Currency`, true
		case "17":
			return `Exchange`, true
		case "18":
			return `MonthYear`, true
		case "19":
			return `UTCTimeStamp`, true
		case "20":
			return `UTCTimeOnly`, true
		case "21":
			return `LocalMktTime`, true
		case "22":
			return `UTCDate`, true
		case "23":
			return `Data`, true
		case "24":
			return `MultipleStringValue`, true
		}

	case 965:
		switch value {
		case "1":
			return `Active`, true
		case "2":
			return `Inactive`, true
		}

	case 974:
		switch value {
		case "FIXED":
			return `FIXED`, true
		case "DIFF":
			return `DIFF`, true
		}

	case 975:
		switch value {
		case "2":
			return `T+1`, true
		case "4":
			return `T+3`, true
		case "5":
			return `T+4`, true
		}

	case 980:
		switch value {
		case "A":
			return `Add`, true
		case "D":
			return `Delete`, true
		case "M":
			return `Modify`, true
		}

	case 982:
		switch value {
		case "1":
			return `Auto Exercise`, true
		case "2":
			return `Non Auto Exercise`, true
		case "3":
			return `Final Will Be Exercised`, true
		case "4":
			return `Contrary Intention`, true
		case "5":
			return `Difference`, true
		}

	case 992:
		switch value {
		case "1":
			return `Sub Allocate`, true
		case "2":
			return `Third Party Allocation`, true
		}

	case 996:
		switch value {
		case "Bcf":
			return `Billion cubic feet`, true
		case "MMbbl":
			return `Million Barrels`, true
		case "MMBtu":
			return `One Million BTU`, true
		case "MWh":
			return `Megawatt hours`, true
		case "Bbl":
			return `Barrels`, true
		case "Bu":
			return `Bushels`, true
		case "lbs":
			return `pounds`, true
		case "Gal":
			return `Gallons`, true
		case "oz_tr":
			return `Troy Ounces`, true
		case "t":
			return `Metric Tons (aka Tonne)`, true
		case "tn":
			return `Tons (US)`, true
		case "USD":
			return `US Dollars`, true
		}

	case 997:
		switch value {
		case "H":
			return `Hour`, true
		case "Min":
			return `Minute`, true
		case "S":
			return `Second`, true
		case "D":
			return `Day`, true
		case "Wk":
			return `Week`, true
		case "Mo":
			return `Month`, true
		case "Yr":
			return `Year`, true
		}

	case 1002:
		switch value {
		case "1":
			return `Automatic`, true
		case "2":
			return `Guarantor`, true
		case "3":
			return `Manual`, true
		}

	case 1015:
		switch value {
		case "0":
			return `false - trade is not an AsOf trade`, true
		case "1":
			return `true - trade is an AsOf trade`, true
		}

	case 1021:
		switch value {
		case "1":
			return `Top of Book`, true
		case "2":
			return `Price Depth`, true
		case "3":
			return `Order Depth`, true
		}

	case 1024:
		switch value {
		case "0":
			return `Book`, true
		case "1":
			return `Off-Book`, true
		case "2":
			return `Cross`, true
		}

	case 1031:
		switch value {
		case "ADD":
			return `Add-on Order`, true
		case "AON":
			return `All or None`, true
		case "CNH":
			return `Cash Not Held`, true
		case "DIR":
			return `Directed Order`, true
		case "E.W":
			return `Exchange for Physical Transaction`, true
		case "FOK":
			return `Fill or Kill`, true
		case "IO":
			return `Imbalance Only`, true
		case "IOC":
			return `Immediate or Cancel`, true
		case "LOO":
			return `Limit On Open`, true
		case "LOC":
			return `Limit on Close`, true
		case "MAO":
			return `Market at Open`, true
		case "MAC":
			return `Market at Close`, true
		case "MOO":
			return `Market on Open`, true
		case "MOC":
			return `Market On Close`, true
		case "MQT":
			return `Minimum Quantity`, true
		case "NH":
			return `Not Held`, true
		case "OVD":
			return `Over the Day`, true
		case "PEG":
			return `Pegged`, true
		case "RSV":
			return `Reserve Size Order`, true
		case "S.W":
			return `Stop Stock Transaction`, true
		case "SCL":
			return `Scale`, true
		case "TMO":
			return `Time Order`, true
		case "TS":
			return `Trailing Stop`, true
		case "WRK":
			return `Work`, true
		}

	case 1032:
		switch value {
		case "1":
			return `NASD OATS`, true
		}

	case 1033:
		switch value {
		case "A":
			return `Agency`, true
		case "AR":
			return `Arbitrage`, true
		case "D":
			return `Derivatives`, true
		case "IN":
			return `International`, true
		case "IS":
			return `Institutional`, true
		case "O":
			return `Other`, true
		case "PF":
			return `Preferred Trading`, true
		case "PR":
			return `Proprietary`, true
		case "PT":
			return `Program Trading`, true
		case "S":
			return `Sales`, true
		case "T":
			return `Trading`, true
		}

	case 1034:
		switch value {
		case "1":
			return `NASD OATS`, true
		}

	case 1035:
		switch value {
		case "ADD":
			return `Add-on Order`, true
		case "AON":
			return `All or None`, true
		case "CNH":
			return `Cash Not Held`, true
		case "DIR":
			return `Directed Order`, true
		case "E.W":
			return `Exchange for Physical Transaction`, true
		case "FOK":
			return `Fill or Kill`, true
		case "IO":
			return `Imbalance Only`, true
		case "IOC":
			return `Immediate or Cancel`, true
		case "LOO":
			return `Limit On Open`, true
		case "LOC":
			return `Limit on Close`, true
		case "MAO":
			return `Market at Open`, true
		case "MAC":
			return `Market at Close`, true
		case "MOO":
			return `Market on Open`, true
		case "MOC":
			return `Market On Close`, true
		case "MQT":
			return `Minimum Quantity`, true
		case "NH":
			return `Not Held`, true
		case "OVD":
			return `Over the Day`, true
		case "PEG":
			return `Pegged`, true
		case "RSV":
			return `Reserve Size Order`, true
		case "S.W":
			return `Stop Stock Transaction`, true
		case "SCL":
			return `Scale`, true
		case "TMO":
			return `Time Order`, true
		case "TS":
			return `Trailing Stop`, true
		case "WRK":
			return `Work`, true
		}

	case 1036:
		switch value {
		case "0":
			return `Received, not yet processed`, true
		case "1":
			return `Accepted`, true
		case "2":
			return `Don't know / Rejected`, true
		}

	case 1043:
		switch value {
		case "0":
			return `Specific Deposit`, true
		case "1":
			return `General`, true
		}

	case 1046:
		switch value {
		case "D":
			return `Divide`, true
		case "M":
			return `Multiply`, true
		}

	case 1047:
		switch value {
		case "O":
			return `Open`, true
		case "C":
			return `Close`, true
		case "R":
			return `Rolled`, true
		case "F":
			return `FIFO`, true
		}

	case 1057:
		switch value {
		case "Y":
			return `Order initiator is aggressor`, true
		case "N":
			return `Order initiator is passive`, true
		}

	case 1070:
		switch value {
		case "0":
			return `Indicative`, true
		case "1":
			return `Tradeable`, true
		case "2":
			return `Restricted Tradeable`, true
		case "3":
			return `Counter`, true
		case "4":
			return `Indicative and Tradeable`, true
		}

	case 1081:
		switch value {
		case "0":
			return `SecondaryOrderID(198)`, true
		case "1":
			return `OrderID(37)`, true
		case "2":
			return `MDEntryID(278)`, true
		case "3":
			return `QuoteEntryID(299)`, true
		}

	case 1083:
		switch value {
		case "1":
			return `Immediate (after each fill)`, true
		case "2":
			return `Exhaust (when DisplayQty = 0)`, true
		}

	case 1084:
		switch value {
		case "1":
			return `Initial (use original DisplayQty)`, true
		case "2":
			return `New (use RefreshQty)`, true
		case "3":
			return `Random (randomize value)`, true
		}

	case 1092:
		switch value {
		case "0":
			return `None`, true
		case "1":
			return `Local (Exchange, ECN, ATS)`, true
		case "2":
			return `National (Across all national markets)`, true
		case "3":
			return `Global (Across all markets)`, true
		}

	case 1093:
		switch value {
		case "1":
			return `Odd Lot`, true
		case "2":
			return `Round Lot`, true
		case "3":
			return `Block Lot`, true
		}

	case 1094:
		switch value {
		case "1":
			return `Last peg (last sale)`, true
		case "2":
			return `Mid-price peg (midprice of inside quote)`, true
		case "3":
			return `Opening peg`, true
		case "4":
			return `Market peg`, true
		case "5":
			return `Primary peg (primary market - buy at bid or sell at offer)`, true
		case "7":
			return `Peg to VWAP`, true
		case "8":
			return `Trailing Stop Peg`, true
		case "9":
			return `Peg to Limit Price`, true
		}

	case 1100:
		switch value {
		case "1":
			return `Partial Execution`, true
		case "2":
			return `Specified Trading Session`, true
		case "3":
			return `Next Auction`, true
		case "4":
			return `Price Movement`, true
		}

	case 1101:
		switch value {
		case "1":
			return `Activate`, true
		case "2":
			return `Modify`, true
		case "3":
			return `Cancel`, true
		}

	case 1107:
		switch value {
		case "1":
			return `Best Offer`, true
		case "2":
			return `Last Trade`, true
		case "3":
			return `Best Bid`, true
		case "4":
			return `Best Bid or Last Trade`, true
		case "5":
			return `Best Offer or Last Trade`, true
		case "6":
			return `Best Mid`, true
		}

	case 1108:
		switch value {
		case "0":
			return `None`, true
		case "1":
			return `Local (Exchange, ECN, ATS)`, true
		case "2":
			return `National (Across all national markets)`, true
		case "3":
			return `Global (Across all markets)`, true
		}

	case 1109:
		switch value {
		case "U":
			return `Trigger if the price of the specified type goes UP to or through the specified Trigger Price.`, true
		case "D":
			return `Trigger if the price of the specified type goes DOWN to or through the specified Trigger Price.`, true
		}

	case 1111:
		switch value {
		case "1":
			return `Market`, true
		case "2":
			return `Limit`, true
		}

	case 1115:
		switch value {
		case "1":
			return `Order`, true
		case "2":
			return `Quote`, true
		case "3":
			return `Privately Negotiated Trade`, true
		case "4":
			return `Multileg order`, true
		case "5":
			return `Linked order`, true
		case "6":
			return `Quote Request`, true
		case "7":
			return `Implied Order`, true
		case "8":
			return `Cross Order`, true
		case "9":
			return `Streaming price (quote)`, true
		}

	case 1123:
		switch value {
		case "0":
			return `Trade Confirmation`, true
		case "1":
			return `Two-Party Report`, true
		case "2":
			return `One-Party Report for Matching`, true
		case "3":
			return `One-Party Report for Pass Through`, true
		case "4":
			return `Automated Floor Order Routing`, true
		case "5":
			return `Two Party Report for Claim`, true
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

	case 1133:
		switch value {
		case "B":
			return `BIC (Bank Identification Code) (ISO 9362)`, true
		case "C":
			return `Generally accepted market participant identifier (e.g. NASD mnemonic)`, true
		case "D":
			return `Proprietary / Custom code`, true
		case "E":
			return `ISO Country Code`, true
		case "G":
			return `MIC (ISO 10383 - Market Identifier Code)`, true
		}

	case 1144:
		switch value {
		case "0":
			return `Not implied`, true
		case "1":
			return `Implied-in - The existence of a multi-leg instrument is implied by the legs of that instrument`, true
		case "2":
			return `Implied-out - The existence of the underlying legs are implied by the multi-leg instrument`, true
		case "3":
			return `Both Implied-in and Implied-out`, true
		}

	case 1159:
		switch value {
		case "1":
			return `Preliminary`, true
		case "2":
			return `Final`, true
		}

	case 1162:
		switch value {
		case "C":
			return `Cancel`, true
		case "N":
			return `New`, true
		case "R":
			return `Replace`, true
		case "T":
			return `Restate`, true
		}

	case 1164:
		switch value {
		case "1":
			return `Instructions of Broker`, true
		case "2":
			return `Instructions for Institution`, true
		case "3":
			return `Investor`, true
		}

	case 1167:
		switch value {
		case "0":
			return `Accepted`, true
		case "5":
			return `Rejected`, true
		case "6":
			return `Removed from Market`, true
		case "7":
			return `Expired`, true
		case "12":
			return `Locked Market Warning`, true
		case "13":
			return `Cross Market Warning`, true
		case "14":
			return `Canceled due to Lock Market`, true
		case "15":
			return `Canceled due to Cross Market`, true
		case "16":
			return `Active`, true
		}

	case 1172:
		switch value {
		case "1":
			return `All market participants`, true
		case "2":
			return `Specified market participants`, true
		case "3":
			return `All Market Makers`, true
		case "4":
			return `Primary Market Maker(s)`, true
		}

	case 1174:
		switch value {
		case "1":
			return `Order imbalance, auction is extended`, true
		case "2":
			return `Trading resumes (after Halt)`, true
		case "3":
			return `Price Volatility Interruption`, true
		case "4":
			return `Change of Trading Session`, true
		case "5":
			return `Change of Trading Subsession`, true
		case "6":
			return `Change of Security Status`, true
		case "7":
			return `Change of Book Type`, true
		case "8":
			return `Change of Market Depth`, true
		}

	case 1176:
		switch value {
		case "1":
			return `Exchange Last`, true
		case "2":
			return `High / Low Price`, true
		case "3":
			return `Average Price (VWAP, TWAP ... )`, true
		case "4":
			return `Turnover (Price * Qty)`, true
		}

	case 1178:
		switch value {
		case "1":
			return `Customer`, true
		}

	case 1193:
		switch value {
		case "C":
			return `Cash settlement required`, true
		case "P":
			return `Physical settlement required`, true
		}

	case 1194:
		switch value {
		case "0":
			return `European`, true
		case "1":
			return `American`, true
		case "2":
			return `Bermuda`, true
		}

	case 1196:
		switch value {
		case "STD":
			return `Standard, money per unit of a physical`, true
		case "INX":
			return `Index`, true
		case "INT":
			return `Interest rate Index`, true
		}

	case 1197:
		switch value {
		case "EQTY":
			return `premium style`, true
		case "FUT":
			return `futures style mark-to-market`, true
		case "FUTDA":
			return `futures style with an attached cash adjustment`, true
		}

	case 1198:
		switch value {
		case "0":
			return `pre-listed only`, true
		case "1":
			return `user requested`, true
		}

	case 1209:
		switch value {
		case "0":
			return `Regular`, true
		case "1":
			return `Variable`, true
		case "2":
			return `Fixed`, true
		case "3":
			return `Traded as a spread leg`, true
		case "4":
			return `Settled as a spread leg`, true
		}

	case 1302:
		switch value {
		case "0":
			return `Months`, true
		case "1":
			return `Days`, true
		case "2":
			return `Weeks`, true
		case "3":
			return `Years`, true
		}

	case 1303:
		switch value {
		case "0":
			return `YearMonth Only (default)`, true
		case "1":
			return `YearMonthDay`, true
		case "2":
			return `YearMonthWeek`, true
		}

	case 1306:
		switch value {
		case "0":
			return `Price`, true
		case "1":
			return `Ticks`, true
		case "2":
			return `Percentage`, true
		}

	case 1307:
		switch value {
		case "0":
			return `Symbol`, true
		case "1":
			return `SecurityType and or CFICode`, true
		case "2":
			return `Product`, true
		case "3":
			return `TradingSessionID`, true
		case "4":
			return `All Securities`, true
		case "5":
			return `UndelyingSymbol`, true
		case "6":
			return `Underlying SecurityType and or CFICode`, true
		case "7":
			return `Underlying Product`, true
		case "8":
			return `MarketID or MarketID + MarketSegmentID`, true
		}

	case 1395:
		switch value {
		case "A":
			return `Add`, true
		case "D":
			return `Delete`, true
		case "M":
			return `Modify`, true
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

	case 1368:
		switch value {
		case "0":
			return `Trading resumes (after Halt)`, true
		case "1":
			return `Change of Trading Session`, true
		case "2":
			return `Change of Trading Subsession`, true
		case "3":
			return `Change of Trading Status`, true
		}

	case 1373:
		switch value {
		case "1":
			return `Suspend orders`, true
		case "2":
			return `Release orders from suspension`, true
		case "3":
			return `Cancel orders`, true
		}

	case 1374:
		switch value {
		case "1":
			return `All orders for a security`, true
		case "2":
			return `All orders for an underlying security`, true
		case "3":
			return `All orders for a Product`, true
		case "4":
			return `All orders for a CFICode`, true
		case "5":
			return `All orders for a SecurityType`, true
		case "6":
			return `All orders for a trading session`, true
		case "7":
			return `All orders`, true
		case "8":
			return `All orders for a Market`, true
		case "9":
			return `All orders for a Market Segment`, true
		case "10":
			return `All orders for a Security Group`, true
		}

	case 1375:
		switch value {
		case "0":
			return `Rejected - See MassActionRejectReason(1376)`, true
		case "1":
			return `Accepted`, true
		}

	case 1376:
		switch value {
		case "0":
			return `Mass Action Not Supported`, true
		case "1":
			return `Invalid or unknown security`, true
		case "2":
			return `Invalid or unknown underlying security`, true
		case "3":
			return `Invalid or unknown Product`, true
		case "4":
			return `Invalid or unknown CFICode`, true
		case "5":
			return `Invalid or unknown SecurityType`, true
		case "6":
			return `Invalid or unknown trading session`, true
		case "7":
			return `Invalid or unknown Market`, true
		case "8":
			return `Invalid or unknown Market Segment`, true
		case "9":
			return `Invalid or unknown Security Group`, true
		case "99":
			return `Other`, true
		}

	case 1377:
		switch value {
		case "0":
			return `Predefined Multileg Security`, true
		case "1":
			return `User-defined Multleg Security`, true
		case "2":
			return `User-defined, Non-Securitized, Multileg`, true
		}

	case 1378:
		switch value {
		case "0":
			return `Net Price`, true
		case "1":
			return `Reversed Net Price`, true
		case "2":
			return `Yield Difference`, true
		case "3":
			return `Individual`, true
		case "4":
			return `Contract Weighted Average Price`, true
		case "5":
			return `Multiplied Price`, true
		}

	case 1385:
		switch value {
		case "1":
			return `One Cancels the Other (OCO)`, true
		case "2":
			return `One Triggers the Other (OTO)`, true
		case "3":
			return `One Updates the Other (OUO) - Absolute Quantity Reduction`, true
		case "4":
			return `One Updates the Other (OUO) - Proportional Quantity Reduction`, true
		}

	case 1386:
		switch value {
		case "0":
			return `Broker / Exchange option`, true
		case "2":
			return `Exchange closed`, true
		case "4":
			return `Too late to enter`, true
		case "5":
			return `Unknown order`, true
		case "6":
			return `Duplicate Order (e.g. dupe ClOrdID)`, true
		case "11":
			return `Unsupported order characteristic`, true
		case "99":
			return `Other`, true
		}

	case 1390:
		switch value {
		case "0":
			return `Do Not Publish Trade`, true
		case "1":
			return `Publish Trade`, true
		case "2":
			return `Deferred Publication`, true
		}

	case 1347:
		switch value {
		case "0":
			return `Retransmission of application messages for the specified Applications`, true
		case "1":
			return `Subscription to the specified Applications`, true
		case "2":
			return `Request for the last ApplLastSeqNum published for the specified Applications`, true
		case "3":
			return `Request valid set of Applications`, true
		case "4":
			return `Unsubscribe to the specified Applications`, true
		}

	case 1348:
		switch value {
		case "0":
			return `Request successfully processed`, true
		case "1":
			return `Application does not exist`, true
		case "2":
			return `Messages not available`, true
		}

	case 1354:
		switch value {
		case "0":
			return `Application does not exist`, true
		case "1":
			return `Messages requested are not available`, true
		case "2":
			return `User not authorized for application`, true
		}

	case 1426:
		switch value {
		case "0":
			return `Reset ApplSeqNum to new value specified in ApplNewSeqNum(1399)`, true
		case "1":
			return `Reports that the last message has been sent for the ApplIDs Refer to RefApplLastSeqNum(1357) for the application sequence number of the last message.`, true
		case "2":
			return `Heartbeat message indicating that Application identified by RefApplID(1355) is still alive. Refer to RefApplLastSeqNum(1357) for the application sequence number of the previous message.`, true
		}

	case 35:
		switch value {
		case "0":
			return `Heartbeat`, true
		case "1":
			return `TestRequest`, true
		case "2":
			return `ResendRequest`, true
		case "3":
			return `Reject`, true
		case "4":
			return `SequenceReset`, true
		case "5":
			return `Logout`, true
		case "6":
			return `IOI`, true
		case "7":
			return `Advertisement`, true
		case "8":
			return `ExecutionReport`, true
		case "9":
			return `OrderCancelReject`, true
		case "A":
			return `Logon`, true
		case "AA":
			return `DerivativeSecurityList`, true
		case "AB":
			return `NewOrderMultileg`, true
		case "AC":
			return `MultilegOrderCancelReplace`, true
		case "AD":
			return `TradeCaptureReportRequest`, true
		case "AE":
			return `TradeCaptureReport`, true
		case "AF":
			return `OrderMassStatusRequest`, true
		case "AG":
			return `QuoteRequestReject`, true
		case "AH":
			return `RFQRequest`, true
		case "AI":
			return `QuoteStatusReport`, true
		case "AJ":
			return `QuoteResponse`, true
		case "AK":
			return `Confirmation`, true
		case "AL":
			return `PositionMaintenanceRequest`, true
		case "AM":
			return `PositionMaintenanceReport`, true
		case "AN":
			return `RequestForPositions`, true
		case "AO":
			return `RequestForPositionsAck`, true
		case "AP":
			return `PositionReport`, true
		case "AQ":
			return `TradeCaptureReportRequestAck`, true
		case "AR":
			return `TradeCaptureReportAck`, true
		case "AS":
			return `AllocationReport`, true
		case "AT":
			return `AllocationReportAck`, true
		case "AU":
			return `ConfirmationAck`, true
		case "AV":
			return `SettlementInstructionRequest`, true
		case "AW":
			return `AssignmentReport`, true
		case "AX":
			return `CollateralRequest`, true
		case "AY":
			return `CollateralAssignment`, true
		case "AZ":
			return `CollateralResponse`, true
		case "B":
			return `News`, true
		case "BA":
			return `CollateralReport`, true
		case "BB":
			return `CollateralInquiry`, true
		case "BC":
			return `NetworkCounterpartySystemStatusRequest`, true
		case "BD":
			return `NetworkCounterpartySystemStatusResponse`, true
		case "BE":
			return `UserRequest`, true
		case "BF":
			return `UserResponse`, true
		case "BG":
			return `CollateralInquiryAck`, true
		case "BH":
			return `ConfirmationRequest`, true
		case "BI":
			return `TradingSessionListRequest`, true
		case "BJ":
			return `TradingSessionList`, true
		case "BK":
			return `SecurityListUpdateReport`, true
		case "BL":
			return `AdjustedPositionReport`, true
		case "BM":
			return `AllocationInstructionAlert`, true
		case "BN":
			return `ExecutionAcknowledgement`, true
		case "BO":
			return `ContraryIntentionReport`, true
		case "BP":
			return `SecurityDefinitionUpdateReport`, true
		case "BQ":
			return `SettlementObligationReport`, true
		case "BR":
			return `DerivativeSecurityListUpdateReport`, true
		case "BS":
			return `TradingSessionListUpdateReport`, true
		case "BT":
			return `MarketDefinitionRequest`, true
		case "BU":
			return `MarketDefinition`, true
		case "BV":
			return `MarketDefinitionUpdateReport`, true
		case "BW":
			return `ApplicationMessageRequest`, true
		case "BX":
			return `ApplicationMessageRequestAck`, true
		case "BY":
			return `ApplicationMessageReport`, true
		case "BZ":
			return `OrderMassActionReport`, true
		case "C":
			return `Email`, true
		case "CA":
			return `OrderMassActionRequest`, true
		case "CB":
			return `UserNotification`, true
		case "D":
			return `NewOrderSingle`, true
		case "E":
			return `NewOrderList`, true
		case "F":
			return `OrderCancelRequest`, true
		case "G":
			return `OrderCancelReplaceRequest`, true
		case "H":
			return `OrderStatusRequest`, true
		case "J":
			return `AllocationInstruction`, true
		case "K":
			return `ListCancelRequest`, true
		case "L":
			return `ListExecute`, true
		case "M":
			return `ListStatusRequest`, true
		case "N":
			return `ListStatus`, true
		case "P":
			return `AllocationInstructionAck`, true
		case "Q":
			return `DontKnowTrade`, true
		case "R":
			return `QuoteRequest`, true
		case "S":
			return `Quote`, true
		case "T":
			return `SettlementInstructions`, true
		case "V":
			return `MarketDataRequest`, true
		case "W":
			return `MarketDataSnapshotFullRefresh`, true
		case "X":
			return `MarketDataIncrementalRefresh`, true
		case "Y":
			return `MarketDataRequestReject`, true
		case "Z":
			return `QuoteCancel`, true
		case "a":
			return `QuoteStatusRequest`, true
		case "b":
			return `MassQuoteAcknowledgement`, true
		case "c":
			return `SecurityDefinitionRequest`, true
		case "d":
			return `SecurityDefinition`, true
		case "e":
			return `SecurityStatusRequest`, true
		case "f":
			return `SecurityStatus`, true
		case "g":
			return `TradingSessionStatusRequest`, true
		case "h":
			return `TradingSessionStatus`, true
		case "i":
			return `MassQuote`, true
		case "j":
			return `BusinessMessageReject`, true
		case "k":
			return `BidRequest`, true
		case "l":
			return `BidResponse`, true
		case "m":
			return `ListStrikePrice`, true
		case "n":
			return `XMLnonFIX`, true
		case "o":
			return `RegistrationInstructions`, true
		case "p":
			return `RegistrationInstructionsResponse`, true
		case "q":
			return `OrderMassCancelRequest`, true
		case "r":
			return `OrderMassCancelReport`, true
		case "s":
			return `NewOrderCross`, true
		case "t":
			return `CrossOrderCancelReplaceRequest`, true
		case "u":
			return `CrossOrderCancelRequest`, true
		case "v":
			return `SecurityTypeRequest`, true
		case "w":
			return `SecurityTypes`, true
		case "x":
			return `SecurityListRequest`, true
		case "y":
			return `SecurityList`, true
		case "z":
			return `DerivativeSecurityListRequest`, true
		}

	}
	return "", false
}
