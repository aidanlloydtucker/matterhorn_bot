package commands

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"net/url"

	"strconv"

	"fmt"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

func newTestBot() (*tgbotapi.BotAPI, chan TestMsgCfg) {
	output := make(chan TestMsgCfg, 1)
	return &tgbotapi.BotAPI{
		Client: &http.Client{
			Transport: RoundTTest{
				Response: output,
			},
		},
		Token: "foobar",
		Self: tgbotapi.User{
			ID:        1234,
			FirstName: "foo",
			LastName:  "bar",
			UserName:  "foobar",
		},
	}, output
}

func newMessageTmpl() *tgbotapi.Message {
	return &tgbotapi.Message{
		MessageID: 1111,
		From: &tgbotapi.User{
			ID:        2222,
			FirstName: "test",
			LastName:  "user",
			UserName:  "testuser",
		},
		Date: int(time.Now().Unix()),
		Chat: &tgbotapi.Chat{
			ID:    5678,
			Type:  "group",
			Title: "test",
		},
	}
}

func newCommand(command string) *tgbotapi.Message {
	msg := newMessageTmpl()
	msg.Text = "/" + command
	return msg
}

type TestMsgCfg struct {
	Message *tgbotapi.MessageConfig
	Photo   *tgbotapi.PhotoConfig
	Voice   *tgbotapi.VoiceConfig
}

type RoundTTest struct {
	Response chan TestMsgCfg
}

func (rt RoundTTest) RoundTrip(r *http.Request) (*http.Response, error) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		rt.Response <- TestMsgCfg{
			Message: valsToMessageConfig(r.Form),
		}
	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse form (10k max allowed in ram)
		err := r.ParseMultipartForm((1 << 10) * 10)
		if err != nil {
			panic(err)
		}

		msgCfg := TestMsgCfg{}

		photo, ok := r.MultipartForm.File["photo"]
		if ok {
			photoCfg := tgbotapi.PhotoConfig{
				Caption: r.Form.Get("caption"),
			}
			photoCfg.File = photo[0]
			msgCfg.Photo = &photoCfg

			rt.Response <- msgCfg
		} else {
			voice, ok := r.MultipartForm.File["voice"]
			if ok {
				voiceCfg := tgbotapi.VoiceConfig{}

				voiceCfg.File = voice[0]
				msgCfg.Voice = &voiceCfg

				rt.Response <- msgCfg
			}
		}
	} else {
		panic("Unknown content type: " + contentType)
	}

	// Returns an error so the request stops and the bot doesn't proceed
	return nil, errors.New("STOP")
}

/*func parseFormURLEncoded(body io.ReadCloser) (*tgbotapi.MessageConfig, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	vals, err := url.ParseQuery(buf.String())
	if err != nil {
		return nil, err
	}

	return valsToMessageConfig(vals), nil
}*/

func valsToMessageConfig(vals url.Values) *tgbotapi.MessageConfig {
	msgConf := tgbotapi.MessageConfig{
		Text:      vals.Get("text"),
		ParseMode: vals.Get("parse_mode"),
	}
	dWP, err := strconv.ParseBool(vals.Get("disable_web_page_preview"))
	if err == nil {
		msgConf.DisableWebPagePreview = dWP
	}
	return &msgConf
}

func checkMessageError(msg TestMsgCfg, resType string) error {
	if resType == "message" {
		if msg.Message == nil {
			return errors.New("No message found")
		}
		if msg.Message.Text == "" {
			return errors.New("Message has no text")
		}
		if strings.HasPrefix(msg.Message.Text, "Error! ") {
			return errors.New("Message has error: " + strings.TrimPrefix(msg.Message.Text, "Error! "))
		}
	} else if resType == "photo" {
		if msg.Photo == nil {
			return errors.New("No photo found")
		}
		if msg.Photo.File == nil {
			return errors.New("No file found")
		}
	} else if resType == "voice" {
		if msg.Voice == nil {
			return errors.New("No voice found")
		}
		if msg.Voice.File == nil {
			return errors.New("No file found")
		}
	}

	return nil
}

/*
 THE START OF TEST FUNCTIONS
   - All commands should have at least one test function
   - Functions are in alphabetical order by filename
*/

// File: 8ball.go
// Checks to make sure that the message text exists
func TestMagicBallHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(MagicBallHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: bash.go
// Checks to make sure that the message text exists
func TestBashHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(BashHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: batman.go
// Checks if the resulting text is valid
func TestBatmanHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(BatmanHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if out.Message.Text != "test user is Batman" {
		t.Fatal("Message doesn't match expected output:", out.Message.Text)
	}
}

// File: bench.go
// Checks if the timestamp recived in the message is at most 0.1s away from the current timestamp
func TestBenchHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(BenchHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	parsedBench, err := strconv.ParseInt(out.Message.Text, 10, 64)
	if err != nil {
		t.Fatal(err)
	}

	benchTime := time.Unix(0, parsedBench)
	if time.Since(benchTime).Seconds() > 0.1 {
		t.Fatal("The seconds elapsed from the message benchmark time to now is greater than 0.1")
	}
}

// File: bitcoin.go
// Checks to make sure that the bitcoin price is parsable as a float64
func TestBitcoinHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(BitcoinHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
	_, err := strconv.ParseFloat(strings.TrimPrefix(out.Message.Text, "Bitcoin's most recent price is $"), 64)
	if err != nil {
		t.Fatal(err)
	}
}

// File: botfather.go
// Checks to make sure that help and echo are mentioned in the result but botfather (the actual command) isn't because
// botfather (the actual command) is a hidden command
func TestBotFatherHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(BotFatherHandler)

	cmdMap := map[string]*CommandInfo{
		"help":      &helpHandlerInfo,
		"echo":      &echoHandlerInfo,
		"botfather": &botFatherHandlerInfo,
	}
	ch.Setup(map[string]interface{}{
		"commandMap": cmdMap,
	})

	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	resFmt := "%s - %s\n%s - %s\n"
	expRes := fmt.Sprintf(resFmt,
		echoHandlerInfo.Command,
		echoHandlerInfo.Description,
		helpHandlerInfo.Command,
		helpHandlerInfo.Description)
	expRes2 := fmt.Sprintf(resFmt,
		helpHandlerInfo.Command,
		helpHandlerInfo.Description,
		echoHandlerInfo.Command,
		echoHandlerInfo.Description)

	if out.Message.Text != expRes && out.Message.Text != expRes2 {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: cat.go
// Checks to make sure that the photo exists
func TestCatHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(CatHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: clear.go
// Checks to make sure that the message text exists
func TestClear_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(ClearHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: echo.go
// Checks to make sure that the resulting text is valid
func TestEchoHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	echoStr := "test 123"

	ch := new(EchoHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{echoStr})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if out.Message.Text != echoStr {
		t.Fatal("Echo is incorrect")
	}
}

// File: fortune.go
// Checks to make sure that the message text exists
func TestFortuneHandler_HandleCommand(t *testing.T) {
	// TODO: Remove this skip once fortune is on all machines
	t.SkipNow()

	bot, output := newTestBot()

	ch := new(FortuneHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: help.go
// Checks to make sure that help (the actual command) and echo are mentioned in the result but botfather isn't because
// botfather is a hidden command
func TestHelpHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(HelpHandler)

	cmdMap := map[string]*CommandInfo{
		"help":      &helpHandlerInfo,
		"echo":      &echoHandlerInfo,
		"botfather": &botFatherHandlerInfo,
	}
	ch.Setup(map[string]interface{}{
		"commandMap": cmdMap,
	})

	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	fmtTxt := "<b>Commands:</b>\n• %s - %s\n• %s - %s\n"

	expRes := fmt.Sprintf(fmtTxt,
		helpHandlerInfo.Command,
		helpHandlerInfo.Description,
		echoHandlerInfo.Command,
		echoHandlerInfo.Description)

	expRes2 := fmt.Sprintf(fmtTxt,
		echoHandlerInfo.Command,
		echoHandlerInfo.Description,
		helpHandlerInfo.Command,
		helpHandlerInfo.Description)

	if out.Message.Text != expRes && out.Message.Text != expRes2 {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: help.go
// Checks to make sure that the detailed command info of echo is valid
func TestHelpHandler_HandleCommand_DetailedInfo(t *testing.T) {
	bot, output := newTestBot()

	ch := new(HelpHandler)

	cmdMap := map[string]*CommandInfo{
		"echo": &echoHandlerInfo,
	}
	ch.Setup(map[string]interface{}{
		"commandMap": cmdMap,
	})
	ch.HandleCommand(bot, newCommand(ch.Info().Command+" echo"), []string{"echo"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if !strings.HasPrefix(out.Message.Text, "<b>echo</b>") {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: info.go
// Checks to make sure that the message is valid
func TestInfoHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(InfoHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if !strings.HasPrefix(out.Message.Text, "<b>foo bar</b>") {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: lenny.go
// Checks to make sure that the message is valid
func TestLennyHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(LennyHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if out.Message.Text != "( ͡° ͜ʖ ͡°)" {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: lines.go
// Checks to make sure that the message is valid
func TestLinesHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(LinesHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"foo"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	expRes := `F O O
O
O`

	if out.Message.Text != expRes {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: lmgtfy.go
// Checks to make sure that the message is valid
func TestLmgtfyHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(LmgtfyHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"foo"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if out.Message.Text != "http://lmgtfy.com/?q=foo" {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: meme.go
// Checks to make sure that the photo exists
// TODO: Fix This: It tries to open ./resources/meme_tmpl/batman.jpg, which works while running but while testing the path should be ../resources/meme_tmpl/batman.jpg
func TestMemeHandler_HandleCommand(t *testing.T) {
	// TODO: Remove this when bug is fixed
	t.SkipNow()

	bot, output := newTestBot()

	ch := new(MemeHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"batman", "foo", "bar"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: memelist.go
// Checks to make sure that the message exists
// TODO: Fix the path bug. Look at TestMemeHandler_HandleCommand for more info
func TestMemeListHandler_HandleCommand(t *testing.T) {
	// TODO: Remove this when bug is fixed
	t.SkipNow()

	bot, output := newTestBot()

	ch := new(MemeListHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: ping.go
// Checks to make sure that the message is valid
func TestPingHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(PingHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if out.Message.Text != "PONG!" {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: random.go
// Checks to make sure that the random int is between [0, 100]. No arguments given to the command
func TestRandomHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(RandomHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	randInt, err := strconv.Atoi(out.Message.Text)
	if err != nil {
		t.Fatal(err)
	}

	if randInt > 100 || randInt < 0 {
		t.Fatal("Random number is not between [0, 100]:", out.Message.Text)
	}
}

// File: random.go
// Checks to make sure that the random int is between [0, 10]. One arguments given to the command
func TestRandomHandler_HandleCommand2(t *testing.T) {
	bot, output := newTestBot()

	ch := new(RandomHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"10"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	randInt, err := strconv.Atoi(out.Message.Text)
	if err != nil {
		t.Fatal(err)
	}

	if randInt > 10 || randInt < 0 {
		t.Fatal("Random number is not between [0, 10]:", out.Message.Text)
	}
}

// File: random.go
// Checks to make sure that the random int is between [3, 7]. Two arguments given to the command
func TestRandomHandler_HandleCommand3(t *testing.T) {
	bot, output := newTestBot()

	ch := new(RandomHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"3", "7"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	randInt, err := strconv.Atoi(out.Message.Text)
	if err != nil {
		t.Fatal(err)
	}

	if randInt > 7 || randInt < 3 {
		t.Fatal("Random number is not between [3, 7]:", out.Message.Text)
	}
}

// File: reddit.go
// Checks to make sure that the message exists
func TestRedditHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(RedditHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"funny"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: rekt.go
// Checks to make sure that the message is valid
func TestRektHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(RektHandler)
	ch.Setup(map[string]interface{}{
		"path": "",
		"reks": []string{"$USER was rekt"},
	})
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"foo"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	if out.Message.Text != "foo was rekt" {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: settings.go
// Checks to make sure that the message exists
func TestSettingsHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(SettingsHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: shame.go
// Checks to make sure that the voice exists
// TODO: Fix This: It tries to open ./resources/shame_sfx.mp3, which works while running but while testing the path should be ./resources/shame_sfx.mp3
func TestShameHandler_HandleCommand(t *testing.T) {
	// TODO: Remove this when bug is fixed
	t.SkipNow()

	bot, output := newTestBot()

	ch := new(ShameHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	fmt.Println("1")
	out := <-output
	fmt.Println("2")
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: square.go
// Checks to make sure that the message is valid
func TestSquareHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(SquareHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"foo"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}

	expRes := `foo
oof
ofo`

	if out.Message.Text != expRes {
		t.Fatal("Message is not the expected result:", out.Message.Text)
	}
}

// File: start.go
// Checks to make sure that the message exists
func TestStartHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(StartHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: urban.go
// Checks to make sure that the message exists
func TestUrbanHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(UrbanHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{"foobar"})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}

// File: xkcd.go
// Checks to make sure that the message exists
func TestXkcdHandler_HandleCommand(t *testing.T) {
	bot, output := newTestBot()

	ch := new(XkcdHandler)
	ch.HandleCommand(bot, newCommand(ch.Info().Command), []string{})

	out := <-output
	if msgErr := checkMessageError(out, ch.Info().ResType); msgErr != nil {
		t.Fatal(msgErr)
	}
}
