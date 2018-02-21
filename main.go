// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	"strings"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot :", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.EqualFold(message.Text, "test1") {

					leftBtn := linebot.NewMessageTemplateAction("left", "left clicked")
					rightBtn := linebot.NewMessageTemplateAction("right", "right clicked")

					template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
					messgage := linebot.NewTemplateMessage("OK YOYO", template)

					if _, err = bot.ReplyMessage(event.ReplyToken, messgage).Do(); err != nil {
						log.Print(err)
					}

				} else if strings.EqualFold(message.Text, "test2") {

					messgage := linebot.NewStickerMessage ("2", "175")

					if _, err = bot.ReplyMessage(event.ReplyToken, messgage).Do(); err != nil {
						log.Print(err)
					}

				} else if strings.EqualFold(message.Text, "test3") {

					messgage := linebot.NewImageMessage("https://mgronline.com/images/mgr-online-logo.png", "https://mgronline.com/images/mgr-online-logo.png")

					if _, err = bot.ReplyMessage(event.ReplyToken, messgage).Do(); err != nil {
						log.Print(err)
					}

				} else if strings.EqualFold(message.Text, "logme") {

					userID := event.Source.UserID
					groupID := event.Source.GroupID
					RoomID := event.Source.RoomID

					messgage := "userID = " + userID + ", groupID = " + groupID + ",RoomID = " + RoomID

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messgage)).Do(); err != nil {
						log.Print(err)
					}

				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! 222")).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}

}
