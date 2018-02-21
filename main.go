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
	"errors"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
	"strings"
	//"io/ioutil"
	//"encoding/json"
	"encoding/json"
	"bytes"
	"time"
	"io/ioutil"
)

const (
	lineApiReplyUrl = "https://api.line.me/v2/bot/message/reply"
	lineApiPushUrl  = "https://api.line.me/v2/bot/message/push"
)

var lineToken = os.Getenv("ChannelAccessToken")
var lineSecret = os.Getenv("ChannelSecret")
var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot YOYO :", bot, " err:", err)
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
				//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK!")).Do(); err != nil {
				//	log.Print(err)
				//}

				if strings.EqualFold(message.Text, "ter") || message.Text == "เต๋อ" {
					//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("0989782592")).Do(); err != nil {
					//	log.Print(err)
					//}
					if err := sendLinePush(event.ReplyToken, "0989782592"); err != nil {
						log.Print(err)
					}
				} else {
					//if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! 222")).Do(); err != nil {
					//	log.Print(err)
					//}
					if err := sendLinePush(event.ReplyToken, message.ID+":"+message.Text+" OK! 222"); err != nil {
						log.Print(err)
					}
				}


			}
		}
	}
}

func sendLinePush(userId string, messages ...string) error {
	type Msgs struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type postJSON struct {
		To      string `json:"to"`
		Message []Msgs `json:"messages"`
	}
	var res postJSON
	res.To = userId

	for _, m := range messages {
		res.Message = append(res.Message, Msgs{"text", m})
	}

	return postLINEAPI(res, lineApiPushUrl)
}

func postLINEAPI(res interface{}, lineApiUrl string) error {
	body, err := json.Marshal(res)
	if err != nil {
		log.Println("Cannot marshal JSON:", err)
		return err
	}
	log.Println("Prepare result:", string(body))
	post, err := http.NewRequest("POST", lineApiUrl, bytes.NewBuffer(body))
	post.Header.Set("Content-Type", "application/json")
	post.Header.Set("Authorization", "Bearer "+lineToken)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	apiRes, err := client.Do(post)
	if err != nil {
		log.Println("Cannot post API:", err)
		return err
	}
	defer apiRes.Body.Close()

	if apiRes.StatusCode != 200 {
		log.Println("API return:", apiRes.StatusCode)
		bodyBytes, _ := ioutil.ReadAll(apiRes.Body)
		log.Println("API body:", string(bodyBytes))
		return errors.New("API return error:" + strconv.Itoa(apiRes.StatusCode) + "," + string(bodyBytes))
	}

	return nil
}

type reqLINEAPIJSON struct {
	Events []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source struct {
			Type    string `json:"type"`
			UserId  string `json:"userId"`
			RoomId  string `json:"roomId"`
			GroupId string `json:"groupId"`
		} `json:"source"`
		Message struct {
			Id        string `json:"id"`
			Type      string `json:"type"`
			Text      string `json:"text"`
			PackageId string `json:"packageId"`
			StickerId string `json:"stickerId"`
		} `json:"message"`
	} `json:"events"`
}

func sendLineReply(replyToken string, messages ...string) error {
	type Msgs struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	type postJSON struct {
		ReplyToken string `json:"replyToken"`
		Message    []Msgs `json:"messages"`
	}
	var res postJSON
	res.ReplyToken = replyToken

	for _, m := range messages {
		res.Message = append(res.Message, Msgs{"text", m})
	}

	return postLINEAPI(res, lineApiReplyUrl)
}

func sendLineSticker(replyToken, packageId, stickerId string) error {
	type Msgs struct {
		Type      string `json:"type"`
		PackageId string `json:"packageId"`
		StickerId string `json:"stickerId"`
	}
	type postJSON struct {
		ReplyToken string `json:"replyToken"`
		Message    []Msgs `json:"messages"`
	}
	var res postJSON
	res.ReplyToken = replyToken
	res.Message = append(res.Message, Msgs{"sticker", packageId, stickerId})

	return postLINEAPI(res, lineApiReplyUrl)
}

func checkLeave(req reqLINEAPIJSON) bool {
	ev0 := req.Events[0]
	if ev0.Message.Type != "text" {
		return false
	}

	if (strings.Contains(ev0.Message.Text, "ขอลา")) ||
		(strings.Contains(ev0.Message.Text, "ลาป่วย")) ||
		(strings.Contains(ev0.Message.Text, "อนุญาตลา")) ||
		(strings.Contains(ev0.Message.Text, "อนุญาติลา")) {
		sendLineSticker(ev0.ReplyToken, "2", "175")
		//sendLineSticker(ev0.ReplyToken, "1004278", "235263")
		return true
	}
	return false
}