package custom

func LoadCustom() {
	addToCustomCommands(QuotesHandler{})
	addToCustomCommands(QuoteHandler{})
}
