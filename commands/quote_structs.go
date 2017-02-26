package commands

import (
	"cloud.google.com/go/datastore"
	"errors"
	chatpkg "github.com/billybobjoeaglt/matterhorn_bot/chat"
	"math/rand"
	"time"
)

type Quote struct {
	Text     string `datastore:",noindex"`
	Date     time.Time
	Author   string
	Document int
	Manual   bool `datastore:",noindex"`
}

const QuoteKeyKind = "Quote"

func addQuote(ds *chatpkg.Datastore, chatID int64, quote Quote) error {
	theChat, _, err := ds.GetChat(chatID)
	if err != nil {
		return err
	}

	quotesDoc := theChat.Settings.QuotesDoc

	if theChat.Settings.QuotesDoc == 0 {
		quotesDoc, err = addChatQuotesDoc(ds, chatID)
		if err != nil {
			return err
		}
	}

	quote.Document = quotesDoc

	quoteKey := datastore.IncompleteKey(QuoteKeyKind, nil)
	quoteKey.Namespace = chatpkg.KeyNamespace

	_, err = ds.Client.Put(ds.Ctx, quoteKey, &quote)
	return err
}

func getQuote(ds *chatpkg.Datastore, chatID int64) (Quote, error) {
	theChat, _, err := ds.GetChat(chatID)
	if err != nil {
		return Quote{}, err
	}

	quotesDoc := theChat.Settings.QuotesDoc

	if theChat.Settings.QuotesDoc == 0 {
		quotesDoc, err = addChatQuotesDoc(ds, chatID)
		if err != nil {
			return Quote{}, err
		}
		return Quote{}, errors.New("No quotes found")
	}

	quotes := []Quote{}

	// TODO: Fix the latency problems with this eventually
	_, err = ds.Client.GetAll(ds.Ctx, datastore.NewQuery(QuoteKeyKind).Namespace(chatpkg.KeyNamespace).Filter("Document =", quotesDoc), &quotes)
	if err != nil {
		return Quote{}, err
	}

	if len(quotes) == 0 {
		return Quote{}, errors.New("No quotes found")
	}

	return quotes[rand.Intn(len(quotes))], nil
}

func addChatQuotesDoc(ds *chatpkg.Datastore, chatID int64) (int, error) {
	updChat, err := ds.UpdateChat(func(oldChat chatpkg.Chat) chatpkg.Chat {
		quotesDoc := oldChat.Settings.QuotesDoc
		if quotesDoc == 0 {
			quotesDoc = int(rand.Int31())
		}
		newChat := oldChat
		newChat.Settings.QuotesDoc = quotesDoc
		return newChat
	}, chatID)
	return updChat.Settings.QuotesDoc, err
}
