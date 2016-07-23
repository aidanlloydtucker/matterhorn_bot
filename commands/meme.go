package commands

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"

	"io"

	"os"

	"io/ioutil"

	"github.com/fogleman/gg"
	"gopkg.in/telegram-bot-api.v4"
)

const FontMin = 40
const LineSpacing = 1.5

func makeMeme(imgFileBytes []byte, topText string, bottomText string) (error, *io.Reader) {
	imgFile := bytes.NewReader(imgFileBytes)

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err, nil
	}

	imgRect := img.Bounds()
	imgW := imgRect.Dx()
	imgH := imgRect.Dy()

	ctx := gg.NewContext(imgW, imgH)

	ctx.DrawImage(img, 0, 0)
	ctx.LoadFontFace("../resources/Impact.ttf", FontMin)

	// Top Text

	lineSpacing := float64((FontMin * len(ctx.WordWrap(topText, float64(imgW)))) + 10)

	// Apply black stroke
	ctx.SetHexColor("#000")
	strokeSize := 6
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			// give it rounded corners
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := float64(imgW/2 + dx)
			y := lineSpacing - float64(dy)
			ctx.DrawStringWrapped(topText, x, y, 0.5, 1, float64(imgW), LineSpacing, gg.AlignCenter)
		}
	}

	// Apply white fill for actual text
	ctx.SetHexColor("#FFF")
	ctx.DrawStringWrapped(topText, float64(imgW)/2, lineSpacing, 0.5, 1, float64(imgW), LineSpacing, gg.AlignCenter)

	// Bottom Text

	// Apply black stroke
	ctx.SetHexColor("#000")
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			// give it rounded corners
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				continue
			}
			x := float64(imgW/2 + dx)
			y := float64(imgH - FontMin + dy)
			ctx.DrawStringWrapped(bottomText, x, y, 0.5, 1, float64(imgW), LineSpacing, gg.AlignCenter)
		}
	}

	// Apply white fill for actual text
	ctx.SetHexColor("#FFF")
	ctx.DrawStringWrapped(bottomText, float64(imgW)/2, float64(imgH)-FontMin, 0.5, 1, float64(imgW), LineSpacing, gg.AlignCenter)

	var b bytes.Buffer
	outWriter := bufio.NewWriter(&b)

	jpeg.Encode(outWriter, ctx.Image(), &jpeg.Options{Quality: jpeg.DefaultQuality})

	outReader := bufio.NewReader(&b)

	return nil, outReader
}

type MemeHandler struct {
}

var memeHandlerInfo = CommandInfo{
	Command:     "meme",
	Args:        `(.+) ["'](.*?)["'] ["'](.*?)["']`,
	Permission:  3,
	Description: "makes a meme",
	LongDesc:    "",
	Usage:       `/meme [meme] "(top text)" "(bottom text)"`,
	Examples: []string{
		`/meme success "uses /meme for the first time" "it works"`,
	},
	ResType: "photo",
}

func (responder MemeHandler) HandleCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, args []string) {
	var errMsg tgbotapi.MessageConfig
	var msg tgbotapi.PhotoConfig

	defer func(errMsg *tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) {
		if *errMsg != nil {
			bot.Send(*errMsg)
		}
	}(&errMsg, bot)

	memeFn := "../resources/meme-tmpl/" + args[0]

	if _, err := os.Stat(memeFn); err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return
	}

	imgFileBytes, err := ioutil.ReadFile(memeFn)
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return

	}

	err, memeImg := makeMeme(imgFileBytes, args[1], args[2])
	if err != nil {
		errMsg = NewErrorMessage(message.Chat.ID, err)
		return

	}
	msg = tgbotapi.NewPhotoUpload(message.Chat.ID, *memeImg)

	bot.Send(msg)
}

func (responder MemeHandler) Info() *CommandInfo {
	return &memeHandlerInfo
}
