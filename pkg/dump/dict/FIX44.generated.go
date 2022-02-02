package dict

type FIX44 struct {
}

func (f FIX44) Version() string {
	return "FIX.4.4"
}

func (f FIX44) TagName(tag int) (string, bool) {
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
	default:
		return "", false
	}
}

func (f FIX44) ValueName(tag int, value string) (string, bool) {
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
		case "1":
			return `per unit (implying shares, par, currency, etc)`, true
		case "2":
			return `percentage`, true
		case "3":
			return `absolute (total monetary amount)`, true
		case "4":
			return `(for CIV buy orders) percentage waived – cash discount`, true
		case "5":
			return `(for CIV buy orders) percentage waived – enhanced units`, true
		case "6":
			return `points per bond or contract [Supply ContractMultiplier (231) in the <Instrument> component block if the object security is denominated in a size other than the industry default - 1000 par for bonds.]`, true
		}

	case 18:
		switch value {
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
		case "0":
			return `Stay on offerside`, true
		case "A":
			return `No cross (cross is forbidden)`, true
		case "B":
			return `OK to cross`, true
		case "C":
			return `Call first`, true
		case "D":
			return `Percent of volume "(indicates that the sender does not want to be all of the volume on the floor vs. a specific percentage)"`, true
		case "E":
			return `Do not increase - DNI`, true
		case "F":
			return `Do not reduce - DNR`, true
		case "G":
			return `All or none - AON`, true
		case "H":
			return `Reinstate on System Failure (mutually exclusive with Q)`, true
		case "I":
			return `Institutions only`, true
		case "J":
			return `Reinstate on Trading Halt (mutually exclusive with K)`, true
		case "K":
			return `Cancel on Trading Halt (mutually exclusive with L)`, true
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
			return `Cancel on System Failure (mutually exclusive with H)`, true
		case "R":
			return `Primary peg (primary market - buy at bid/sell at offer)`, true
		case "S":
			return `Suspend`, true
		case "U":
			return `Customer Display Instruction (Rule11Ac1-1/4)`, true
		case "V":
			return `Netting (for Forex)`, true
		case "W":
			return `Peg to VWAP`, true
		case "X":
			return `Trade Along`, true
		case "Y":
			return `Try to Stop`, true
		case "Z":
			return `Cancel if Not Best`, true
		case "a":
			return `Trailing Stop Peg`, true
		case "b":
			return `Strict Limit (No Price Improvement)`, true
		case "c":
			return `Ignore Price Validity Checks`, true
		case "d":
			return `Peg to Limit Price`, true
		case "e":
			return `Work to Target Strategy`, true
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
			return `ISDA/FpML Product Specification`, true
		case "J":
			return `Options Price Reporting Authority`, true
		}

	case 25:
		switch value {
		case "L":
			return `Low`, true
		case "M":
			return `Medium`, true
		case "H":
			return `High`, true
		}

	case 27:
		switch value {
		case "S":
			return `Small`, true
		case "M":
			return `Medium`, true
		case "L":
			return `Large`, true
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
			return `Order – Single`, true
		case "E":
			return `Order – List`, true
		case "F":
			return `Order Cancel Request`, true
		case "G":
			return `Order Cancel/Replace Request`, true
		case "H":
			return `Order Status Request`, true
		case "J":
			return `Allocation Instruction`, true
		case "K":
			return `List Cancel Request`, true
		case "L":
			return `List Execute`, true
		case "M":
			return `List Status Request`, true
		case "N":
			return `List Status`, true
		case "P":
			return `Allocation Instruction Ack`, true
		case "Q":
			return `Don’t Know Trade (DK)`, true
		case "R":
			return `Quote Request`, true
		case "S":
			return `Quote`, true
		case "T":
			return `Settlement Instructions`, true
		case "V":
			return `Market Data Request`, true
		case "W":
			return `Market Data-Snapshot/Full Refresh`, true
		case "X":
			return `Market Data-Incremental Refresh`, true
		case "Y":
			return `Market Data Request Reject`, true
		case "Z":
			return `Quote Cancel`, true
		case "a":
			return `Quote Status Request`, true
		case "b":
			return `Mass Quote Acknowledgement`, true
		case "c":
			return `Security Definition Request`, true
		case "d":
			return `Security Definition`, true
		case "e":
			return `Security Status Request`, true
		case "f":
			return `Security Status`, true
		case "g":
			return `Trading Session Status Request`, true
		case "h":
			return `Trading Session Status`, true
		case "i":
			return `Mass Quote`, true
		case "j":
			return `Business Message Reject`, true
		case "k":
			return `Bid Request`, true
		case "l":
			return `Bid Response (lowercase L)`, true
		case "m":
			return `List Strike Price`, true
		case "n":
			return `XML message (e.g. non-FIX MsgType)`, true
		case "o":
			return `Registration Instructions`, true
		case "p":
			return `Registration Instructions Response`, true
		case "q":
			return `Order Mass Cancel Request`, true
		case "r":
			return `Order Mass Cancel Report`, true
		case "s":
			return `New Order - Cross`, true
		case "t":
			return `Cross Order Cancel/Replace Request (a.k.a. Cross Order Modification Request)`, true
		case "u":
			return `Cross Order Cancel Request`, true
		case "v":
			return `Security Type Request`, true
		case "w":
			return `Security Types`, true
		case "x":
			return `Security List Request`, true
		case "y":
			return `Security List`, true
		case "z":
			return `Derivative Security List Request`, true
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
		case "AJ":
			return `Quote Response`, true
		case "AK":
			return `Confirmation`, true
		case "AL":
			return `Position Maintenance Request`, true
		case "AM":
			return `Position Maintenance Report`, true
		case "AN":
			return `Request For Positions`, true
		case "AO":
			return `Request For Positions Ack`, true
		case "AP":
			return `Position Report`, true
		case "AQ":
			return `Trade Capture Report Request Ack`, true
		case "AR":
			return `Trade Capture Report Ack`, true
		case "AS":
			return `Allocation Report (aka Allocation Claim)`, true
		case "AT":
			return `Allocation Report Ack (aka Allocation Claim Ack)`, true
		case "AU":
			return `Confirmation Ack (aka Affirmation)`, true
		case "AV":
			return `Settlement Instruction Request`, true
		case "AW":
			return `Assignment Report`, true
		case "AX":
			return `Collateral Request`, true
		case "AY":
			return `Collateral Assignment`, true
		case "AZ":
			return `Collateral Response`, true
		case "BA":
			return `Collateral Report`, true
		case "BB":
			return `Collateral Inquiry`, true
		case "BC":
			return `Network (Counterparty System) Status Request`, true
		case "BD":
			return `Network (Counterparty System) Status Response`, true
		case "BE":
			return `User Request`, true
		case "BF":
			return `User Response`, true
		case "BG":
			return `Collateral Inquiry Ack`, true
		case "BH":
			return `Confirmation Request`, true
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
		case "6":
			return `With or without`, true
		case "7":
			return `Limit or better (Deprecated)`, true
		case "8":
			return `Limit with or without`, true
		case "9":
			return `On basis`, true
		case "D":
			return `Previously quoted`, true
		case "E":
			return `Previously indicated`, true
		case "G":
			return `Forex - Swap`, true
		case "I":
			return `Funari (Limit Day Order with unexecuted portion handled as Market On Close. E.g. Japan)`, true
		case "J":
			return `Market If Touched (MIT)`, true
		case "K":
			return `Market with Leftover as Limit (market order then unexecuted quantity becomes limit order at last price)`, true
		case "L":
			return `Previous Fund Valuation Point (Historic pricing) (for CIV)`, true
		case "M":
			return `Next Fund Valuation Point –(Forward pricing) (for CIV)`, true
		case "P":
			return `Pegged`, true
		}

	case 43:
		switch value {
		case "Y":
			return `Possible duplicate`, true
		case "N":
			return `Original transmission`, true
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
			return `Short Exempt Transaction for Principal (was incorrectly identified in the FIX spec as "Registered Equity Market Maker trades")`, true
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
			return `Proprietary transactions for competing market-maker that is affiliated with the clearing member (was incorrectly identified in the FIX spec as "Competing dealer trades")`, true
		case "P":
			return `Principal`, true
		case "R":
			return `Transactions for the account of a non-member competing market maker (was incorrectly identified in the FIX spec as "Competing dealer trades")`, true
		case "S":
			return `Specialist trades`, true
		case "T":
			return `Transactions for the account of an unaffiliated member’s competing market maker (was incorrectly identified in the FIX spec as "Competing dealer trades")`, true
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
			return `Undisclosed (valid for IOI and List Order messages only)`, true
		case "8":
			return `Cross (orders where counterparty is an exchange, valid for all messages except IOIs)`, true
		case "9":
			return `Cross short`, true
		case "A":
			return `Cross short exempt`, true
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
			return `Immediate or Cancel (IOC)`, true
		case "4":
			return `Fill or Kill (FOK)`, true
		case "5":
			return `Good Till Crossing (GTX)`, true
		case "6":
			return `Good Till Date`, true
		case "7":
			return `At the Close`, true
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
			return `Next Day (T+1)`, true
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
		}

	case 77:
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
			return `block level reject`, true
		case "2":
			return `account level reject`, true
		case "3":
			return `received (received, not yet processed)`, true
		case "4":
			return `incomplete`, true
		case "5":
			return `rejected by intermediary`, true
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
			return `unknown OrderID (37)`, true
		case "6":
			return `unknown ListID (66)`, true
		case "7":
			return `other (further in Note 58=)`, true
		case "8":
			return `incorrect allocated quantity`, true
		case "9":
			return `calculation difference`, true
		case "10":
			return `unknown or stale ExecID (17)`, true
		case "11":
			return `mismatched data value (further in Note 58=)`, true
		case "12":
			return `unknown ClOrdID (11)`, true
		case "13":
			return `warehouse request rejected`, true
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
		case "Y":
			return `Possible resend`, true
		case "N":
			return `Original transmission`, true
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
			return `Unknown Order`, true
		case "6":
			return `Duplicate Order (e.g. dupe ClOrdID (11))`, true
		case "7":
			return `Duplicate of a verbally communicated order`, true
		case "8":
			return `Stale Order`, true
		case "9":
			return `Trade Along required`, true
		case "10":
			return `Invalid Investor ID`, true
		case "11":
			return `Unsupported order characteristic12 = Surveillence Option`, true
		case "13":
			return `Incorrect quantity`, true
		case "14":
			return `Incorrect allocated quantity`, true
		case "15":
			return `Unknown account(s)`, true
		case "99":
			return `Other`, true
		}

	case 104:
		switch value {
		case "A":
			return `All or none`, true
		case "B":
			return `Market On Close (MOC) (held to close)`, true
		case "C":
			return `At the close (around/not held to close)`, true
		case "D":
			return `VWAP (Volume Weighted Avg Price)`, true
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
		case "R":
			return `Ready to trade`, true
		case "S":
			return `Portfolio shown`, true
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
		case "F":
			return `Calculation difference`, true
		case "Z":
			return `Other`, true
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
		case "0":
			return `New`, true
		case "3":
			return `Done for day`, true
		case "4":
			return `Canceled`, true
		case "5":
			return `Replace`, true
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
			return `Restated (ExecutionRpt sent unsolicited by sellside, with ExecRestatementReason (378) set)`, true
		case "E":
			return `Pending Replace (e.g. result of Order Cancel/Replace Request)`, true
		case "F":
			return `Trade (partial fill or fill)`, true
		case "G":
			return `Trade Correct (formerly an ExecTransType (20))`, true
		case "H":
			return `Trade Cancel (formerly an ExecTransType)`, true
		case "I":
			return `Order Status (formerly an ExecTransType)`, true
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
		case "1":
			return `Standing Instructions Provided`, true
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
			return `Restate (used where the Settlement Instruction is being used to communicate standing instructions which have not been changed or added to)`, true
		}

	case 165:
		switch value {
		case "1":
			return `Broker’s Instructions`, true
		case "2":
			return `Institution’s Instructions`, true
		case "3":
			return `Investor (e.g. CIV use)`, true
		}

	case 167:
		switch value {
		case "FUT":
			return `Future`, true
		case "OPT":
			return `Option`, true
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
		case "XLINKD":
			return `Indexed Linked`, true
		case "STRUCT":
			return `Structured Notes`, true
		case "YANK":
			return `Yankee Corporate Bond`, true
		case "FOR":
			return `Foreign Exchange Contract`, true
		case "CS":
			return `Common Stock`, true
		case "PS":
			return `Preferred Stock`, true
		case "BRADY":
			return `Brady Bond`, true
		case "EUSOV":
			return `Euro Sovereigns *`, true
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
			return `US Treasury Note (deprecated value, use "TNOTE")`, true
		case "USTB":
			return `US Treasury Bill (deprecated value, use "TBILL")`, true
		case "TNOTE":
			return `US Treasury Note`, true
		case "TBILL":
			return `US Treasury Bill`, true
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
		case "TERM":
			return `Term Loan`, true
		case "RVLV":
			return `Revolver Loan`, true
		case "RVLVTRM":
			return `Revolver/Term Loan`, true
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
		case "BN":
			return `Bank Notes`, true
		case "BOX":
			return `Bill of Exchanges`, true
		case "CD":
			return `Certificate of Deposit`, true
		case "CL":
			return `Call Loans`, true
		case "CP":
			return `Commercial Paper`, true
		case "DN":
			return `Deposit Notes`, true
		case "EUCD":
			return `Euro Certificate of Deposit`, true
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
		case "PZFJ":
			return `Plazos Fijos`, true
		case "STN":
			return `Short Term Loan Note`, true
		case "TD":
			return `Time Deposit`, true
		case "XCN":
			return `Extended Comm Note`, true
		case "YCD":
			return `Yankee Certificate of Deposit`, true
		case "ABS":
			return `Asset-backed Securities`, true
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
			return `To be Announced`, true
		case "AN":
			return `Other Anticipation Notes BAN, GAN, etc.`, true
		case "COFO":
			return `Certificate of Obligation`, true
		case "COFP":
			return `Certificate of Participation`, true
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
		case "TRAN":
			return `Tax & Revenue Anticipation Note`, true
		case "VRDN":
			return `Variable Rate Demand Note`, true
		case "WAR":
			return `Warrant`, true
		case "MF":
			return `Mutual Fund (i.e. any kind of open-ended "Collective Investment Vehicle")`, true
		case "MLEG":
			return `Multi-leg instrument (e.g. options strategy or futures spread. CFICode (461) can be used to identify if options-based, futures-based, etc.)`, true
		case "NONE":
			return `No Security Type`, true
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
			return `A Global Custodian (StandInstDbName (170) must be provided)`, true
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

	case 208:
		switch value {
		case "Y":
			return `Details should be communicated`, true
		case "N":
			return `Details should not be communicated`, true
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
			return `5-YR`, true
		case "3":
			return `OLD-5`, true
		case "4":
			return `10-YR`, true
		case "5":
			return `OLD-10`, true
		case "6":
			return `30-YR`, true
		case "7":
			return `OLD-30`, true
		case "8":
			return `3-MO-LIBOR`, true
		case "9":
			return `6-MO-LIBOR`, true
		}

	case 233:
		switch value {
		case "AMT":
			return `AMT (y/n)`, true
		case "AUTOREINV":
			return `Auto Reinvestment at <rate> or better`, true
		case "BANKQUAL":
			return `Bank qualified (y/n)`, true
		case "BGNCON":
			return `Bargain Conditions– see (234) for values`, true
		case "COUPON":
			return `Coupon range`, true
		case "CURRENCY":
			return `ISO Currency code`, true
		case "CUSTOMDATE":
			return `Custom start/end date`, true
		case "GEOG":
			return `Geographics and % Range (ex. 234=CA 0-80 [minimum of 80% California assets])`, true
		case "HAIRCUT":
			return `Valuation discount`, true
		case "INSURED":
			return `Insured (y/n)`, true
		case "ISSUE":
			return `Year or Year/Month of Issue (ex. 234=2002/09)`, true
		case "ISSUER":
			return `Issuer’s ticker`, true
		case "ISSUESIZE":
			return `issue size range`, true
		case "LOOKBACK":
			return `Lookback days`, true
		case "LOT":
			return `Explicit lot identifier`, true
		case "LOTVAR":
			return `Lot Variance (value in percent maximum over- or under-allocation allowed)`, true
		case "MAT":
			return `Maturity Year and Month`, true
		case "MATURITY":
			return `Maturity range`, true
		case "MAXSUBS":
			return `Maximum substitutions (Repo)`, true
		case "MINQTY":
			return `Minimum quantity`, true
		case "MININCR":
			return `Minimum increment`, true
		case "MINDNOM":
			return `Minimum denomination`, true
		case "PAYFREQ":
			return `Payment frequency, calendar`, true
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
		case "PRICE":
			return `Price range`, true
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
			return `Type of redemption – values are: NonCallable, Callable, Prefunded, EscrowedToMaturity, Putable, Convertible`, true
		case "RESTRICTED":
			return `Restricted (y/n)`, true
		case "SECTOR":
			return `Market sector`, true
		case "SECTYPE":
			return `SecurityType included or excluded`, true
		case "STRUCT":
			return `Structure`, true
		case "SUBSFREQ":
			return `Substitutions frequency (Repo)`, true
		case "SUBSLEFT":
			return `Substitutions left (Repo)`, true
		case "TEXT":
			return `Freeform text`, true
		case "TRDVAR":
			return `Trade Variance (value in percent maximum over- or under-allocation allowed)`, true
		case "WAC":
			return `Weighted Average Coupon:value in percent (exact or range) plus ‘Gross’ or ‘Net’ of servicing spread (the default) (ex. 234=6.5- Net [minimum of 6.5% net of servicing fee])`, true
		case "WAL":
			return `Weighted Average Life Coupon: value in percent (exact or range)`, true
		case "WALA":
			return `Weighted Average Loan Age: value in months (exact or range)`, true
		case "WAM":
			return `Weighted Average Maturity : value in months (exact or range)`, true
		case "WHOLE":
			return `Whole Pool (y/n)`, true
		case "YIELD":
			return `Yield range`, true
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
			return `Yield To Average Maturity`, true
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
		case "GROSS":
			return `True Gross Yield`, true
		case "GOVTEQUIV":
			return `Government Equivalent Yield`, true
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
			return `Mark To Market Yield`, true
		case "MATURITY":
			return `Yield to Maturity`, true
		case "NEXTREFUND":
			return `Yield To Next Refund (Sinking Fund Bonds)`, true
		case "OPENAVG":
			return `Open Average Yield`, true
		case "PUT":
			return `Yield to Next Put`, true
		case "PREVCLOSE":
			return `Previous Close Yield`, true
		case "PROCEEDS":
			return `Proceeds Yield`, true
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
		case "VALUE1/32":
			return `Yield Value Of 1/32`, true
		case "WORST":
			return `Yield To Worst`, true
		}

	case 258:
		switch value {
		case "Y":
			return `Traded Flat`, true
		case "N":
			return `Not Traded Flat`, true
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
			return `Open / Active`, true
		case "B":
			return `Closed / Inactive`, true
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
			return `Next Day (only) Market`, true
		case "E":
			return `Opening / Reopening Trade Detail`, true
		case "F":
			return `Intraday Trade Detail`, true
		case "G":
			return `Rule 127 Trade (NYSE)`, true
		case "H":
			return `Rule 155 Trade (Amex)`, true
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
			return `Imbalance More Buyers (Cannot be used in combination with Q)`, true
		case "Q":
			return `Imbalance More Sellers (Cannot be used in combination with P)`, true
		case "R":
			return `Opening Price`, true
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
		}

	case 297:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Canceled for Symbol(s)`, true
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
			return `Canceled due to lock market`, true
		case "15":
			return `Canceled due to cross market`, true
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
		}

	case 300:
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
			return `Unknown Quote`, true
		case "6":
			return `Duplicate Quote`, true
		case "7":
			return `Invalid bid/ask spread`, true
		case "8":
			return `Invalid price`, true
		case "9":
			return `Not authorized to quote security`, true
		case "99":
			return `Other`, true
		}

	case 301:
		switch value {
		case "0":
			return `No Acknowledgement (Default)`, true
		case "1":
			return `Acknowledge only negative or erroneous quotes`, true
		case "2":
			return `Acknowledge each quote messages`, true
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
			return `Request Security identity for the specifications provided (Name of the security is not supplied)`, true
		case "2":
			return `Request List Security Types`, true
		case "3":
			return `Request List Securities (Can be qualified with Symbol, SecurityType, TradingSessionID, SecurityExchange. If provided then only list Securities for the specific type)`, true
		}

	case 323:
		switch value {
		case "1":
			return `Accept security proposal as is`, true
		case "2":
			return `Accept security proposal with revisions as indicated in the message`, true
		case "5":
			return `Reject security proposal`, true
		case "6":
			return `Can not match selection criteria`, true
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
		case "1":
			return `Opening Delay`, true
		case "2":
			return `Trading Halt`, true
		case "3":
			return `Resume`, true
		case "4":
			return `No Open/No Resume`, true
		case "5":
			return `Price Indication`, true
		case "6":
			return `Trading Range Indication`, true
		case "7":
			return `Market Imbalance Buy`, true
		case "8":
			return `Market Imbalance Sell`, true
		case "9":
			return `Market On Close Imbalance Buy`, true
		case "10":
			return `Market On Close Imbalance Sell`, true
		case "12":
			return `No Market Imbalance`, true
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
		case "18":
			return `Not Available for trading (end of session)`, true
		case "19":
			return `Not Traded on this Market`, true
		case "20":
			return `Unknown or Invalid`, true
		case "21":
			return `Pre-Open`, true
		case "22":
			return `Opening Rotation`, true
		case "23":
			return `Fast Market`, true
		}

	case 327:
		switch value {
		case "I":
			return `Order Imbalance`, true
		case "X":
			return `Equipment Changeover`, true
		case "P":
			return `News Pending`, true
		case "D":
			return `News Dissemination`, true
		case "E":
			return `Order Influx`, true
		case "M":
			return `Additional Information`, true
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

	case 347:
		switch value {
		case "ISO-2022-JP":
			return `JIS`, true
		case "EUC-JP":
			return `EUC`, true
		case "Shift_JIS":
			return `for using SJIS`, true
		case "UTF-8":
			return `Unicode`, true
		}

	case 373:
		switch value {
		case "0":
			return `Invalid tag number`, true
		case "1":
			return `Required tag missing`, true
		case "2":
			return `Tag not defined for this message type`, true
		case "3":
			return `Undefined Tag`, true
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
			return `SendingTime accuracy problem`, true
		case "11":
			return `Invalid MsgType`, true
		case "12":
			return `XML Validation error`, true
		case "13":
			return `Tag appears more than once`, true
		case "14":
			return `Tag specified out of required order`, true
		case "15":
			return `Repeating group fields out of order`, true
		case "16":
			return `Incorrect NumInGroup count for repeating group`, true
		case "17":
			return `Non "data" value includes field delimiter (SOH character)`, true
		case "99":
			return `Other`, true
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
		case "Y":
			return `Was solcitied`, true
		case "N":
			return `Was not solicited`, true
		}

	case 378:
		switch value {
		case "0":
			return `GT Corporate action`, true
		case "1":
			return `GT renewal / restatement (no corporate action)`, true
		case "2":
			return `Verbal change`, true
		case "3":
			return `Repricing of order`, true
		case "4":
			return `Broker option`, true
		case "5":
			return `Partial decline of OrderQty (e.g. exchange-initiated partial cancel)`, true
		case "6":
			return `Cancel on Trading Halt`, true
		case "7":
			return `Cancel on System Failure`, true
		case "8":
			return `Market (Exchange) Option`, true
		case "9":
			return `Canceled, Not Best`, true
		case "10":
			return `Warehouse recap`, true
		case "99":
			return `Other`, true
		}

	case 380:
		switch value {
		case "0":
			return `Other`, true
		case "1":
			return `Unkown ID`, true
		case "2":
			return `Unknown Security`, true
		case "3":
			return `Unsupported Message Type`, true
		case "4":
			return `Application not available`, true
		case "5":
			return `Conditionally Required Field Missing`, true
		case "6":
			return `Not authorized`, true
		case "7":
			return `DeliverTo firm not available at this time`, true
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
		case "6":
			return `Related to VWAP`, true
		}

	case 394:
		switch value {
		case "1":
			return `"Non Disclosed" Style (e.g. US/European)`, true
		case "2":
			return `"Disclosed" Style (e.g. Japanese)`, true
		case "3":
			return `No Bidding Process`, true
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
			return `SideValue1`, true
		case "2":
			return `SideValue 2`, true
		}

	case 409:
		switch value {
		case "1":
			return `5day moving average`, true
		case "2":
			return `20 day moving average`, true
		case "3":
			return `Normal Market Size`, true
		case "4":
			return `Other`, true
		}

	case 411:
		switch value {
		case "Y":
			return `True`, true
		case "N":
			return `False`, true
		}

	case 414:
		switch value {
		case "1":
			return `BuySide explicitly requests status using StatusRequest (Default) The sell-side firm can however, send a DONE status List Status Response in an unsolicited fashion`, true
		case "2":
			return `SellSide periodically sends status using ListStatus. Period optionally specified in ProgressPeriod`, true
		case "3":
			return `Real-time execution reports (to be discouraged)`, true
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
		case "R":
			return `Risk Trade`, true
		case "G":
			return `VWAP Guarantee`, true
		case "A":
			return `Agency`, true
		case "J":
			return `Guaranteed Close`, true
		}

	case 419:
		switch value {
		case "2":
			return `Closing Price at morning session`, true
		case "3":
			return `Closing Price`, true
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
			return `Percentage (e.g. percent of par) (often called "dollar price" for fixed income)`, true
		case "2":
			return `Per unit (i.e. per share or contract)`, true
		case "3":
			return `Fixed Amount (absolute value)`, true
		case "4":
			return `Discount – percentage points below par`, true
		case "5":
			return `Premium – percentage points over par`, true
		case "6":
			return `Spread`, true
		case "7":
			return `TED price`, true
		case "8":
			return `TED yield`, true
		case "9":
			return `Yield`, true
		case "10":
			return `Fixed cabinet trade price (primarily for listed futures and options)`, true
		case "11":
			return `Variable cabinet trade price (primarily for listed futures and options)`, true
		}

	case 427:
		switch value {
		case "0":
			return `book out all trades on day of execution`, true
		case "1":
			return `accumulate executions until order is filled or expires`, true
		case "2":
			return `accumulate until verbally notified otherwise`, true
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
			return `ExecStarted`, true
		case "5":
			return `AllDone`, true
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
			return `InBiddingProcess`, true
		case "2":
			return `ReceivedForExecution`, true
		case "3":
			return `Executing`, true
		case "4":
			return `Canceling`, true
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
			return `Wait for Execute Instruction (e.g. a List Execute message or phone call before proceeding with execution of the list)`, true
		case "3":
			return `Exchange/switch CIV order – Sell driven`, true
		case "4":
			return `Exchange/switch CIV order – Buy driven, cash top-up (i.e. additional cash will be provided to fulfil the order)`, true
		case "5":
			return `Exchange/switch CIV order – Buy driven, cash withdraw (i.e. additional cash will not be provided to fulfil the order)`, true
		}

	case 434:
		switch value {
		case "1":
			return `Order Cancel Request`, true
		case "2":
			return `Order Cancel/Replace Request`, true
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
		case "B":
			return `BIC (Bank Identification Code—Swift managed) code (ISO 9362 - See "Appendix 6-B")`, true
		case "C":
			return `Generally accepted market participant identifier (e.g. NASD mnemonic)`, true
		case "D":
			return `Proprietary/Custom code`, true
		case "E":
			return `ISO Country Code`, true
		case "F":
			return `Settlement Entity Location (note if Local Market Settlement use "E = ISO Country Code") (see "Appendix 6-G" for valid values)`, true
		case "G":
			return `MIC (ISO 10383 - Market Identifier Code) (See "Appendix 6-C")`, true
		case "H":
			return `CSD participant/member code (e.g. Euroclear, DTC, CREST or Kassenverein number)`, true
		case "1":
			return `Korean Investor ID`, true
		case "2":
			return `Taiwanese Qualified Foreign Investor ID QFII / FID`, true
		case "3":
			return `Taiwanese Trading Account`, true
		case "4":
			return `Malaysian Central Depository (MCD) number`, true
		case "5":
			return `Chinese B Share (Shezhen and Shanghai)`, true
		case "6":
			return `UK National Insurance or Pension Number`, true
		case "7":
			return `US Social Security Number`, true
		case "8":
			return `US Employer Identification Number`, true
		case "9":
			return `Australian Business Number`, true
		case "A":
			return `Australian Tax File Number`, true
		case "I":
			return `Directed broker three character acronym as defined in ISITC ‘ETC Best Practice’ guidelines document`, true
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
			return `Locate/Lending Firm (for short-sales)`, true
		case "9":
			return `Fund manager Client ID (for CIV)`, true
		case "10":
			return `Settlement Location (formerly FIX 4.2 SettlLocation)`, true
		case "11":
			return `Order Origination Trader (associated with Order Origination Firm – e.g. trader who initiates/submits the order)`, true
		case "12":
			return `Executing Trader (associated with Executing Firm - actually executes)`, true
		case "13":
			return `Order Origination Firm (e.g. buyside firm)`, true
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
			return `Sub custodian`, true
		case "32":
			return `Beneficiary`, true
		case "33":
			return `Interested party`, true
		case "34":
			return `Regulatory body`, true
		case "35":
			return `Liquidity provider`, true
		case "36":
			return `Entering Trader`, true
		case "37":
			return `Contra Trader`, true
		case "38":
			return `Position Account`, true
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
		case "Y":
			return `True (Test)`, true
		case "N":
			return `False (Production)`, true
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
			return `FedWire`, true
		case "8":
			return `Direct Credit (BECS, BACS)`, true
		case "9":
			return `ACH Credit`, true
		case "10":
			return `BPAY`, true
		case "11":
			return `High Value Clearing System (HVACS)`, true
		case "12":
			return `Reinvest in fund`, true
		}

	case 480:
		switch value {
		case "Y":
			return `Yes`, true
		case "N":
			return `No – execution only`, true
		case "M":
			return `No – waiver agreement`, true
		case "O":
			return `No – institutional.`, true
		}

	case 481:
		switch value {
		case "Y":
			return `Passed`, true
		case "N":
			return `Not checked`, true
		case "1":
			return `Exempt – Below The Limit`, true
		case "2":
			return `Exempt – Client Money Type Exemption`, true
		case "3":
			return `Exempt – Authorised Credit or Financial Institution.`, true
		}

	case 484:
		switch value {
		case "B":
			return `Bid price`, true
		case "C":
			return `Creation price`, true
		case "D":
			return `Creation price plus adjustment %`, true
		case "E":
			return `Creation price plus adjustment amount`, true
		case "O":
			return `Offer price`, true
		case "P":
			return `Offer price minus adjustment %`, true
		case "Q":
			return `Offer price minus adjustment amount`, true
		case "S":
			return `Single price`, true
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
			return `FedWire`, true
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
			return `Mini Stocks and Shares ISA (UK)`, true
		case "5":
			return `Mini Insurance ISA (UK)`, true
		case "6":
			return `Current year payment (US)`, true
		case "7":
			return `Prior year payment (US)`, true
		case "8":
			return `Asset transfer (US)`, true
		case "9":
			return `Employee - prior year (US)`, true
		case "10":
			return `Employee – current year (US)`, true
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
		case "16":
			return `Individual Retirement Account (US)`, true
		case "17":
			return `Individual Retirement Account – Rollover (US)`, true
		case "18":
			return `KEOGH (US)`, true
		case "19":
			return `Profit Sharing Plan (US)`, true
		case "20":
			return `401K (US)`, true
		case "21":
			return `Self-Directed IRA (US)`, true
		case "22":
			return `403(b) (US)`, true
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
		}

	case 497:
		switch value {
		case "Y":
			return `Yes`, true
		case "N":
			return `No`, true
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
			return `Reminder – i.e. Registration Instructions are still outstanding`, true
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
			return `Invalid/unacceptable No Reg Detls`, true
		case "5":
			return `Invalid/unacceptable Reg Seq No`, true
		case "6":
			return `Invalid/unacceptable Reg Dtls`, true
		case "7":
			return `Invalid/unacceptable Mailing Dtls`, true
		case "8":
			return `Invalid/unacceptable Mailing Inst`, true
		case "9":
			return `Invalid/unacceptable Investor ID`, true
		case "10":
			return `Invalid/unacceptable Investor ID Source`, true
		case "11":
			return `Invalid/unacceptable Date of Birth`, true
		case "12":
			return `Invalid/unacceptable Investor Country Of Residence`, true
		case "13":
			return `Invalid/unacceptable NoDistribInstns`, true
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
		case "1":
			return `Replace`, true
		case "2":
			return `Cancel`, true
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
		case "13":
			return `Fund-based Renewal Commission Amount (based on Order value)`, true
		case "14":
			return `Fund-based Renewal Commission Amount (based on Projected Fund value)`, true
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
			return `Networking Sub-Account`, true
		case "11":
			return `Non-Profit Organization`, true
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
			return `Principal (Note for CMS purposes, Principal includes Proprietary)`, true
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
			return `Acting as Market Maker or Specialist in the underlying security of a derivative security`, true
		case "7":
			return `Foreign Entity (of foreign governmnet or regulatory jurisdiction)`, true
		case "8":
			return `External Market Participant`, true
		case "9":
			return `External Inter-connected Market Linkage`, true
		case "A":
			return `Riskless Arbitrage`, true
		}

	case 530:
		switch value {
		case "1":
			return `Cancel orders for a security`, true
		case "2":
			return `Cancel orders for an Underlying security`, true
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
		}

	case 531:
		switch value {
		case "0":
			return `Cancel Request Rejected -- See MassCancelRejectReason (532)`, true
		case "1":
			return `Cancel orders for a security`, true
		case "2":
			return `Cancel orders for an Underlying security`, true
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
		}

	case 532:
		switch value {
		case "0":
			return `Mass Cancel Not Supported`, true
		case "1":
			return `Invalid or unknown Security`, true
		case "2":
			return `Invalid or unknown underlying`, true
		case "3":
			return `Invalid or unknown Product`, true
		case "4":
			return `Invalid or unknown CFICode`, true
		case "5":
			return `Invalid or unknown Security Type`, true
		case "6":
			return `Invalid or unknown trading session`, true
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
			return `Counter (tradable)`, true
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
			return `Cross Trade which is executed partially and the rest is cancelled. One side is fully executed, the other side is partially executed with the remainder being cancelled. This is equivalent to an Immediate or Cancel on the other side. Note: The CrossPrioritzation (550) field may be used to indicate which side should fully execute in this scenario.`, true
		case "3":
			return `Cross trade which is partially executed with the unfilled portions remaining active. One side of the cross is fully executed (as denoted with the CrossPrioritization field), but the unfilled portion remains active.`, true
		case "4":
			return `Cross trade is executed with existing orders with the same price. In the case other orders exist with the same price, the quantity of the Cross is executed against the existing orders and quotes, the remainder of the cross is executed against the other side of the cross. The two sides potentially have different quantities.`, true
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
			return `one side`, true
		case "2":
			return `both sides`, true
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
			return `Report by mulitleg security only (Do not report legs)`, true
		case "1":
			return `Report by multileg security and by instrument legs belonging to the multileg security.`, true
		case "2":
			return `Report by instrument legs belonging to the multileg security only (Do not report status of multileg security)`, true
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
			return `All trades`, true
		case "1":
			return `Matched trades matching Criteria provided on request (parties, exec id, trade id, order id, instrument, input source, etc.)`, true
		case "2":
			return `Unmatched trades that match criteria`, true
		case "3":
			return `Unreported trades that match criteria`, true
		case "4":
			return `Advisories that match criteria`, true
		}

	case 570:
		switch value {
		case "Y":
			return `previously reported to counterparty`, true
		case "N":
			return `not reported to counterparty`, true
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
		case "A1":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus four badges and execution time (within two-minute window)`, true
		case "A2":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus four badges`, true
		case "A3":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus two badges and execution time (within two-minute window)`, true
		case "A4":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus two badges`, true
		case "A5":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator plus execution time (within two-minute window)`, true
		case "AQ":
			return `Compared records resulting from stamped advisories or specialist accepts/pair-offs`, true
		case "S1":
			return `Summarized Match using A1 exact match criteria except quantity is summarized`, true
		case "S2":
			return `Summarized Match using A2 exact match criteria except quantity is summarized`, true
		case "S3":
			return `Summarized Match using A3 exact match criteria except quantity is summarized`, true
		case "S4":
			return `Summarized Match using A4 exact match criteria except quantity is summarized`, true
		case "S5":
			return `Summarized Match using A5 exact match criteria except quantity is summarized`, true
		case "M1":
			return `Exact match on Trade Date, Stock Symbol, Quantity, Price, Trade Type, and Special Trade Indicator minus badges And times: ACT M1 match`, true
		case "M2":
			return `Summarized match minus badges and times: ACT M2 Match`, true
		case "MT":
			return `OCS Locked In: Non-ACT`, true
		case "M3":
			return `ACT Accepted Trade`, true
		case "M4":
			return `ACT Default Trade`, true
		case "M5":
			return `ACT Default After M2`, true
		case "M6":
			return `ACT M6 Match`, true
		}

	case 575:
		switch value {
		case "Y":
			return `treat as odd lot`, true
		case "N":
			return `treat as round lot`, true
		}

	case 577:
		switch value {
		case "0":
			return `process normally`, true
		case "1":
			return `exclude from all netting`, true
		case "2":
			return `bilateral netting only`, true
		case "3":
			return `ex clearing`, true
		case "4":
			return `special trade`, true
		case "5":
			return `multilateral netting`, true
		case "6":
			return `clear against central counterparty`, true
		case "7":
			return `exclude from central counterparty`, true
		case "8":
			return `Manual mode (pre-posting and/or pre-giveup)`, true
		case "9":
			return `Automatic posting mode (trade posting to the position account number specified)`, true
		case "10":
			return `Automatic give-up mode (trade give-up to the give-up destination number specified)`, true
		case "11":
			return `Qualified Service Representative (QSR) -`, true
		case "12":
			return `Customer Trade`, true
		case "13":
			return `Self clearing`, true
		}

	case 581:
		switch value {
		case "1":
			return `Account is carried on customer Side of Books`, true
		case "2":
			return `Account is carried on non-Customer Side of books`, true
		case "3":
			return `House Trader`, true
		case "4":
			return `Floor Trader`, true
		case "6":
			return `Account is carried on non-customer side of books and is cross margined`, true
		case "7":
			return `Account is house trader and is cross margined`, true
		case "8":
			return `Joint Backoffice Account (JBO)`, true
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
			return `Do not pro-rata = discuss first`, true
		}

	case 626:
		switch value {
		case "1":
			return `Calculated (includes MiscFees and NetMoney)`, true
		case "2":
			return `Preliminary (without MiscFees and NetMoney)`, true
		case "5":
			return `Ready-To-Book`, true
		case "7":
			return `Warehouse instruction`, true
		case "8":
			return `Request to Intermediary`, true
		}

	case 635:
		switch value {
		case "B":
			return `CBOE Member`, true
		case "C":
			return `Non-member and Customer`, true
		case "E":
			return `Equity Member and Clearing Member`, true
		case "F":
			return `Full and Associate Member trading for own account and as floor Brokers`, true
		case "H":
			return `106.H and 106.J Firms`, true
		case "I":
			return `GIM, IDEM and COM Membership Interest Holders`, true
		case "L":
			return `Lessee and 106.F Employees`, true
		case "M":
			return `All other ownership types`, true
		case "1":
			return `1st year delegate trading for his own account`, true
		case "2":
			return `2nd year delegate trading for his own account`, true
		case "3":
			return `3rd year delegate trading for his own account`, true
		case "4":
			return `4th year delegate trading for his own account`, true
		case "5":
			return `5th year delegate trading for his own account`, true
		case "9":
			return `6th year and beyond delegate trading for his own account`, true
		}

	case 636:
		switch value {
		case "Y":
			return `Order is currently being worked`, true
		case "N":
			return `Order has been accepted but not yet in a working state`, true
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
		case "7":
			return `No match for inquiry`, true
		case "8":
			return `No market for instrument`, true
		case "9":
			return `No inventory`, true
		case "10":
			return `Pass`, true
		case "99":
			return `Other`, true
		}

	case 660:
		switch value {
		case "1":
			return `BIC`, true
		case "2":
			return `SID code`, true
		case "3":
			return `TFM (GSPTA)`, true
		case "4":
			return `OMGEO (AlertID)`, true
		case "5":
			return `DTCC code`, true
		case "99":
			return `Other (custom or proprietary)`, true
		}

	case 665:
		switch value {
		case "1":
			return `Received`, true
		case "2":
			return `Mismatched account`, true
		case "3":
			return `Missing settlement instructions`, true
		case "4":
			return `Confirmed`, true
		case "5":
			return `Request rejected`, true
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
			return `BookEntry`, true
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
			return `percent (percent of par)`, true
		case "2":
			return `per share (e.g. cents per share)`, true
		case "3":
			return `fixed amount (absolute value)`, true
		case "4":
			return `discount – percentage points below par`, true
		case "5":
			return `premium – percentage points over par`, true
		case "6":
			return `basis points relative to benchmark`, true
		case "7":
			return `TED price`, true
		case "8":
			return `TED yield`, true
		case "9":
			return `Yield spread (swaps)`, true
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
		}

	case 703:
		switch value {
		case "TQ":
			return `Transaction Quantity`, true
		case "IAS":
			return `Intra-Spread Qty`, true
		case "IES":
			return `Inter-Spread Qty`, true
		case "FIN":
			return `End-of-Day Qty`, true
		case "SOD":
			return `Start-of-Day Qty`, true
		case "EX":
			return `Option Exercise Qty`, true
		case "AS":
			return `Option Assignment`, true
		case "TX":
			return `Transaction from Exercise`, true
		case "TA":
			return `Transaction from Assignment`, true
		case "PIT":
			return `Pit Trade Qty`, true
		case "TRF":
			return `Transfer Trade Qty`, true
		case "ETR":
			return `Electronic Trade Qty`, true
		case "ALC":
			return `Allocation Trade Qty`, true
		case "PA":
			return `Adjustment Qty`, true
		case "ASF":
			return `As-of Trade Qty`, true
		case "DLV":
			return `Delivery Qty`, true
		case "TOT":
			return `Total Transaction Qty`, true
		case "XM":
			return `Cross Margin Qty`, true
		case "SPL":
			return `Integral Split`, true
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
		case "FMTM":
			return `Final Mark-to-Market Amount`, true
		case "IMTM":
			return `Incremental Mark-to-Market Amount`, true
		case "TVAR":
			return `Trade Variation Amount`, true
		case "SMTM":
			return `Start-of-Day Mark-to-Market Amount`, true
		case "PREM":
			return `Premium Amount`, true
		case "CRES":
			return `Cash Residual Amount`, true
		case "CASH":
			return `Cash Amount (Corporate Event)`, true
		case "VADJ":
			return `Value Adjusted Amount`, true
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
		}

	case 712:
		switch value {
		case "1":
			return `New: used to increment the overall transaction quantity`, true
		case "2":
			return `Replace: used to override the overall transaction quantity or specific add messages based on the reference id`, true
		case "3":
			return `Cancel: used to remove the overall transaction or specific add messages based on reference id`, true
		}

	case 716:
		switch value {
		case "ITD":
			return `Intraday`, true
		case "RTH":
			return `Regular Trading Hours`, true
		case "ETH":
			return `Electronic Trading Hours`, true
		}

	case 718:
		switch value {
		case "0":
			return `Process request as Margin Disposition`, true
		case "1":
			return `Delta_plus`, true
		case "2":
			return `Delta_minus`, true
		case "3":
			return `Final`, true
		}

	case 722:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Accepted with Warnings`, true
		case "2":
			return `Rejected`, true
		case "3":
			return `Completed`, true
		case "4":
			return `Completed with Warnings`, true
		}

	case 723:
		switch value {
		case "0":
			return `Successful completion - no warnings or errors`, true
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
		}

	case 725:
		switch value {
		case "0":
			return `Inband: transport the request was sent over (Default)`, true
		case "1":
			return `Out-of-Band: pre-arranged out of band delivery mechanism (i.e. FTP, HTTP, NDM, etc) between counterparties. Details specified via ResponseDestination (726).`, true
		}

	case 728:
		switch value {
		case "0":
			return `Valid Request`, true
		case "1":
			return `Invalid or unsupported Request`, true
		case "2":
			return `No positions found that match criteria`, true
		case "3":
			return `Not authorized to request positions`, true
		case "4":
			return `Request for Position not supported`, true
		case "99":
			return `Other (use Text(58) in conjunction with this code for an explanation)`, true
		}

	case 729:
		switch value {
		case "0":
			return `Completed`, true
		case "1":
			return `Completed with Warnings`, true
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
		case "R":
			return `Random`, true
		case "P":
			return `ProRata`, true
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
			return `Successful (Default)`, true
		case "1":
			return `Invalid or unknown instrument`, true
		case "2":
			return `Invalid type of trade requested`, true
		case "3":
			return `Invalid parties`, true
		case "4":
			return `Invalid Transport Type requested`, true
		case "5":
			return `Invalid Destination requested`, true
		case "8":
			return `TradeRequestType not supported`, true
		case "9":
			return `Unauthorized for Trade Capture Report Request`, true
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
			return `Successful (Default)`, true
		case "1":
			return `Invalid party information`, true
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
			return `Individual leg of a multi-leg security`, true
		case "3":
			return `Multi-leg security`, true
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
		}

	case 773:
		switch value {
		case "1":
			return `Status`, true
		case "2":
			return `Confirmation`, true
		case "3":
			return `Confirmation Request Rejected (reason can be stated in Text field)`, true
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
			return `CFD (Contract For Difference)`, true
		case "2":
			return `Total return swap`, true
		}

	case 780:
		switch value {
		case "0":
			return `use default instructions`, true
		case "1":
			return `derive from parameters provided`, true
		case "2":
			return `full details provided`, true
		case "3":
			return `SSI db ids provided`, true
		case "4":
			return `phone for instructions`, true
		}

	case 787:
		switch value {
		case "S":
			return `securities`, true
		case "C":
			return `cash`, true
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
			return `unable to process request (e.g. database unavailable)`, true
		case "1":
			return `unknown account`, true
		case "2":
			return `no matching settlement instructions found`, true
		case "99":
			return `other`, true
		}

	case 794:
		switch value {
		case "3":
			return `Sellside Calculated Using Preliminary (includes MiscFees and NetMoney)`, true
		case "4":
			return `Sellside Calculated Without Preliminary (sent unsolicited by sellside, includes MiscFees and NetMoney)`, true
		case "5":
			return `Warehouse recap`, true
		case "8":
			return `Request to Intermediary`, true
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
			return `Account is carried on customer Side of Books`, true
		case "2":
			return `Account is carried on non-Customer Side of books`, true
		case "3":
			return `House Trader`, true
		case "4":
			return `Floor Trader`, true
		case "6":
			return `Account is carried on non-customer side of books and is cross margined`, true
		case "7":
			return `Account is house trader and is cross margined`, true
		case "8":
			return `Joint Backoffice Account (JBO)`, true
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
			return `Postal address (inclusive of street address, location, and postal code)`, true
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
			return `CSD participant/member code (e.g. Euroclear, DTC, CREST or Kassenverein number)`, true
		case "18":
			return `Registered address`, true
		case "19":
			return `Fund/account name`, true
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
			return `Location / Desk`, true
		case "26":
			return `Position Account Type`, true
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
			return `No action taken`, true
		case "1":
			return `Queue flushed`, true
		case "2":
			return `Overlay last`, true
		case "3":
			return `End session`, true
		}

	case 815:
		switch value {
		case "0":
			return `No action taken`, true
		case "1":
			return `Queue flushed`, true
		case "2":
			return `Overlay last`, true
		case "3":
			return `End session`, true
		}

	case 819:
		switch value {
		case "0":
			return `No Average Pricing`, true
		case "1":
			return `Trade is part of an average price group identified by the TradeLinkID`, true
		case "2":
			return `Last Trade in the average price group identified by the TradeLinkID`, true
		}

	case 826:
		switch value {
		case "0":
			return `Allocation not required`, true
		case "1":
			return `Allocation required (give up trade) allocation information not provided (incomplete)`, true
		case "2":
			return `Use allocation provided with the trade`, true
		}

	case 827:
		switch value {
		case "0":
			return `Expire on trading session close (default)`, true
		case "1":
			return `Expire on trading session open`, true
		}

	case 828:
		switch value {
		case "0":
			return `Regular Trade`, true
		case "1":
			return `Block Trade`, true
		case "2":
			return `EFP (Exchange for Physical)`, true
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
			return `Strict – limit is a strict limit`, true
		case "2":
			return `Or worse – for a buy the peg limit is a minimum and for a sell the peg limit is a maximum (for use for orders which have a price range)`, true
		}

	case 838:
		switch value {
		case "1":
			return `More aggressive – on a buy order round the price up round up to the nearest tick, on a sell round down to the nearest tick`, true
		case "2":
			return `More passive – on a buy order round down to nearest tick on a sell order round up to nearest tick`, true
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
			return `Strict – limit is a strict limit`, true
		case "2":
			return `Or worse – for a buy the discretion price is a minimum and for a sell the discretion price is a maximum (for use for orders which have a price range)`, true
		}

	case 844:
		switch value {
		case "1":
			return `More aggressive – on a buy order round the price up round up to the nearest tick, on a sell round down to the nearest tick`, true
		case "2":
			return `More passive – on a buy order round down to nearest tick on a sell order round up to nearest tick`, true
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
		}

	case 852:
		switch value {
		case "Y":
			return `Report trade`, true
		case "N":
			return `Do not report trade`, true
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
			return `Qualifed Service Representative (QSR) or Automatic Giveup (AGU) Contra Side Sold Short`, true
		case "5":
			return `QSR or AGU Contra Side Sold Short Exempt`, true
		}

	case 854:
		switch value {
		case "0":
			return `Units (shares, par, currency)`, true
		case "1":
			return `Contracts (if used - should specify ContractMultiplier (tag 231))`, true
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
			return `Locked In Trade Break`, true
		}

	case 857:
		switch value {
		case "0":
			return `Not specified`, true
		case "1":
			return `Explicit list provided`, true
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
			return `Coupon period (if not semi-annual). Supply redemption date in the InstrAttribValue (872) field`, true
		case "9":
			return `When [and if] issued`, true
		case "10":
			return `Original issue discount`, true
		case "11":
			return `Callable, puttable`, true
		case "12":
			return `Escrowed to Maturity`, true
		case "13":
			return `Escrowed to redemption date – callable. Supply redemption date in the InstrAttribValue (872) field`, true
		case "14":
			return `Prerefunded`, true
		case "15":
			return `In default`, true
		case "16":
			return `Unrated`, true
		case "17":
			return `Taxable`, true
		case "18":
			return `Indexed`, true
		case "19":
			return `Subject to Alternative Minimum Tax`, true
		case "20":
			return `Original issue discount price. Supply price in the InstrAttribValue (872) field`, true
		case "21":
			return `Callable below maturity value`, true
		case "22":
			return `Callable without notice by mail to holder unless registered`, true
		case "99":
			return `Text. Supply the text of the attribute or disclaimer in the InstrAttribValue (872) field`, true
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
			return `Per unit`, true
		case "2":
			return `Percentage`, true
		}

	case 893:
		switch value {
		case "Y":
			return `Last message`, true
		case "N":
			return `Not last message`, true
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
			return `TradeDate`, true
		case "1":
			return `GC Instrument`, true
		case "2":
			return `CollateralInstrument`, true
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

	case 919:
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

	case 924:
		switch value {
		case "1":
			return `LogOnUser`, true
		case "2":
			return `LogOffUser`, true
		case "3":
			return `ChangePasswordForUser`, true
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
		}

	case 928:
		switch value {
		case "1":
			return `Connected`, true
		case "2":
			return `Not connected – down expected up`, true
		case "3":
			return `Not connected – down expected down`, true
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
			return `Stop subscribing`, true
		case "8":
			return `Level of detail, then NoCompID’s becomes required`, true
		}

	case 937:
		switch value {
		case "1":
			return `Full`, true
		case "2":
			return `Incremental update`, true
		}

	case 939:
		switch value {
		case "0":
			return `Accepted`, true
		case "1":
			return `Rejected`, true
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
			return `Accepted with Warnings`, true
		case "2":
			return `Completed`, true
		case "3":
			return `Completed with Warnings`, true
		case "4":
			return `Rejected`, true
		}

	case 946:
		switch value {
		case "0":
			return `Successful (Default)`, true
		case "1":
			return `Invalid or unknown instrument`, true
		case "2":
			return `Invalid or unknown collateral type`, true
		case "3":
			return `Invalid parties`, true
		case "4":
			return `Invalid Transport Type requested`, true
		case "5":
			return `Invalid Destination requested`, true
		case "6":
			return `No collateral found for the trade specified`, true
		case "7":
			return `No collateral found for the order specified`, true
		case "8":
			return `Collateral Inquiry type not supported`, true
		case "9":
			return `Unauthorized for collateral inquiry`, true
		case "99":
			return `Other (further information in Text (58) field)`, true
		}

	}
	return "", false
}
