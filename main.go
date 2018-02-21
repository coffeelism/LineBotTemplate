 //Licensed under the Apache License, Version 2.0 (the "License");
 //you may not use this file except in compliance with the License.
 //You may obtain a copy of the License at
 //
 //    http://www.apache.org/licenses/LICENSE-2.0
 //
 //Unless required by applicable law or agreed to in writing, software
 //distributed under the License is distributed on an "AS IS" BASIS,
 //WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 //See the License for the specific language governing permissions and
 //limitations under the License.

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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("0989782592")).Do(); err != nil {
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

//package main
//
//import (
//	"bytes"
//	//"crypto/hmac"
//	//"crypto/sha256"
//	//"crypto/tls"
//	//"encoding/base64"
//	"encoding/json"
//	"errors"
//	//"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	//"os/exec"
//	"strconv"
//	//"strings"
//	"time"
//
//	//"encoding/hex"
//	//"github.com/Roasbeef/perm-crypt"
//	"github.com/gin-gonic/gin"
//	"math/rand"
//	//"regexp"
//	//"sync/atomic"
//
//	//"github.com/olebedev/when"
//	//"github.com/olebedev/when/rules/common"
//	//"github.com/olebedev/when/rules/en"
//	//"gopkg.in/mgo.v2"
//	//"gopkg.in/mgo.v2/bson"
//	"strings"
//	//"fmt"
//	//"crypto/tls"
//)
//
//const (
//	lineApiReplyUrl = "https://api.line.me/v2/bot/message/reply"
//	lineApiPushUrl  = "https://api.line.me/v2/bot/message/push"
//)
//
//var lineToken = os.Getenv("ChannelAccessToken")
//var lineSecret = os.Getenv("ChannelSecret")
//
//func main() {
//	//loc, _ := time.LoadLocation("Asia/Bangkok")
//
//	//if len(os.Args) > 1 && os.Args[1] == "batch" {
//	//	startTime := time.Date(2017, time.June, 17, 22, 35, 0, 0, loc)
//	//	endTime := time.Date(2017, time.June, 17, 22, 55, 0, 0, loc)
//	//	currentTime := time.Now()
//	//	log.Println("Current time:", currentTime)
//	//	if currentTime.After(startTime) && currentTime.Before(endTime) {
//	//		log.Println("Matched criteria:", currentTime)
//	//		sendLinePush("U8256b3232f36110d799607f5b19b7684", "ทดสอบเรื่อยๆ", time.Now().String())
//	//	}
//	//	return
//	//}
//
//	//isHeroku := true
//	//isApp01 := false
//	port := os.Getenv("PORT")
//	if port == "" {
//		port = "443"
//	}
//	//if port == "" {
//	//	port = os.Args[1]
//	//	if port == "" {
//	//		log.Fatal("$PORT must be set")
//	//	} else {
//	//		isApp01 = true
//	//	}
//	//} else {
//	//	isHeroku = true
//	//}
//	//log.Println("isHeroku:", isHeroku, "isApp01:", isApp01)
//
//	rand.Seed(time.Now().UnixNano())
//
//	//if isHeroku {
//	//	// connect mongo
//	//	var err error
//	//	session, err = mgo.Dial(mongoHerokuURL)
//	//	if err != nil {
//	//		panic(err)
//	//	}
//	//	// defer session.Close()
//	//	session.SetMode(mgo.Monotonic, true)
//	//
//	//	// run loop for send push
//	//	go loopPushNotify()
//	//}
//
//	router := gin.New()
//	router.Use(gin.Logger())
//	router.LoadHTMLGlob("templates/*.tmpl.html")
//	router.Static("/static", "static")
//
//	router.GET("/", func(c *gin.Context) {
//		c.HTML(http.StatusOK, "index.tmpl.html", nil)
//	})
//
//	// for use in Heroku server
//	router.POST("/callback", callbackHandler)
//
//	// for use in App01 server
//	//router.POST("/run-bot01", runBot01)
//
//	router.Run(":" + port)
//}
//
//func sendLinePush(userId string, messages ...string) error {
//	type Msgs struct {
//		Type string `json:"type"`
//		Text string `json:"text"`
//	}
//	type postJSON struct {
//		To      string `json:"to"`
//		Message []Msgs `json:"messages"`
//	}
//	var res postJSON
//	res.To = userId
//
//	for _, m := range messages {
//		res.Message = append(res.Message, Msgs{"text", m})
//	}
//
//	return postLINEAPI(res, lineApiPushUrl)
//}
//
//func postLINEAPI(res interface{}, lineApiUrl string) error {
//	body, err := json.Marshal(res)
//	if err != nil {
//		log.Println("Cannot marshal JSON:", err)
//		return err
//	}
//	log.Println("Prepare result:", string(body))
//	post, err := http.NewRequest("POST", lineApiUrl, bytes.NewBuffer(body))
//	post.Header.Set("Content-Type", "application/json")
//	post.Header.Set("Authorization", "Bearer "+lineToken)
//	client := &http.Client{
//		Timeout: 10 * time.Second,
//	}
//	apiRes, err := client.Do(post)
//	if err != nil {
//		log.Println("Cannot post API:", err)
//		return err
//	}
//	defer apiRes.Body.Close()
//
//	if apiRes.StatusCode != 200 {
//		log.Println("API return:", apiRes.StatusCode)
//		bodyBytes, _ := ioutil.ReadAll(apiRes.Body)
//		log.Println("API body:", string(bodyBytes))
//		return errors.New("API return error:" + strconv.Itoa(apiRes.StatusCode) + "," + string(bodyBytes))
//	}
//
//	return nil
//}
//
//type reqLINEAPIJSON struct {
//	Events []struct {
//		ReplyToken string `json:"replyToken"`
//		Type       string `json:"type"`
//		Timestamp  int64  `json:"timestamp"`
//		Source struct {
//			Type    string `json:"type"`
//			UserId  string `json:"userId"`
//			RoomId  string `json:"roomId"`
//			GroupId string `json:"groupId"`
//		} `json:"source"`
//		Message struct {
//			Id        string `json:"id"`
//			Type      string `json:"type"`
//			Text      string `json:"text"`
//			PackageId string `json:"packageId"`
//			StickerId string `json:"stickerId"`
//		} `json:"message"`
//	} `json:"events"`
//}
//
//func callbackHandler(c *gin.Context) {
//	log.Println("Start J-BOT01")
//
//	// =========== read request ==============
//	defer c.Request.Body.Close()
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		log.Println("Cannot read body:", err)
//		c.String(http.StatusBadRequest, "Cannot read body:%s", err)
//		return
//	}
//	//if !validateSignature(lineSecret, c.Request.Header.Get("X-Line-Signature"), body) {
//	//	log.Println("Invalid Signature")
//	//	c.String(http.StatusOK, "Invalid signature")
//	//	return
//	//}
//
//	log.Println("Request Body:", string(body))
//
//	var req reqLINEAPIJSON
//	//err = c.BindJSON(&req)
//	err = json.Unmarshal(body, &req)
//	if err != nil {
//		log.Println("Error when bindJSON:", err)
//		c.String(http.StatusBadRequest, "Cannot parse request input:%s", c.Request.RequestURI, "..", err)
//		return
//	}
//	log.Printf("LINE Unmarshal result:%+v\n", req)
//	//ev0 := req.Events[0]
//	//msgTime := time.Unix(0, ev0.Timestamp*int64(time.Millisecond))
//	// ========= end read request ==========
//
//	// ========= write post to API ==========
//	// validate j-bot request
//	// sticker
//	//if ev0.Message.Type == "sticker" && ev0.Source.UserId == "U8256b3232f36110d799607f5b19b7684" {
//	//	msg1 := ev0.Message.Type + "," + ev0.Message.PackageId + "," + ev0.Message.StickerId
//	//	sendLineReply(ev0.ReplyToken, msg1)
//	//	return
//	//}
//	//if socialShare(req) {
//	//	return
//	//}
//	//if mukSeaw(req) {
//	//	return
//	//}
//	if checkLeave(req) {
//		return
//	}
//	//if eatRaiD(req) {
//	//	return
//	//}
//	//if encrypt(req) {
//	//	return
//	//}
//	//if pleaseAdd(req) {
//	//	return
//	//}
//	//if remind(req) {
//	//	return
//	//}
//	//
//	//if ev0.Message.Type != "text" || !strings.HasPrefix(strings.ToLower(ev0.Message.Text), "jbot ") {
//	//	c.String(http.StatusOK, "")
//	//	return
//	//}
//	//jbotCmd := ev0.Message.Text[5:]
//	//if len(jbotCmd) <= 1 {
//	//	sendLineReply(ev0.ReplyToken, "Sorry, invalid J-BOT01 command")
//	//	return
//	//}
//	//if _, ok := allowUsers[ev0.Source.UserId]; !ok {
//	//	msg := "Sorry, for security reason, only authorized user can run J-Bot01 command\n"
//	//	for _, v := range allowUsers {
//	//		msg += v + ", "
//	//	}
//	//	msg += "\n-------------------------\n"
//	//	msg += ev0.Source.UserId
//	//	msg += "\n-------------------------"
//	//	sendLineReply(ev0.ReplyToken, msg)
//	//	return
//	//}
//	//if jbotCmd == "leave" || jbotCmd == "ออกไป" {
//	//	sendLineReply(ev0.ReplyToken, "บ๊ายบาย")
//	//
//	//	leaveUrl := fmt.Sprintf("https://api.line.me/v2/bot/group/%s/leave", ev0.Source.GroupId)
//	//	if ev0.Source.Type == "room" {
//	//		leaveUrl = fmt.Sprintf("https://api.line.me/v2/bot/room/%s/leave", ev0.Source.RoomId)
//	//	}
//	//	post, err := http.NewRequest("POST", leaveUrl, nil)
//	//	post.Header.Set("Authorization", "Bearer "+lineToken)
//	//	client := &http.Client{
//	//		Timeout: 10 * time.Second,
//	//	}
//	//	apiRes, err := client.Do(post)
//	//	if err != nil {
//	//		log.Println("Cannot post API leave group:", err)
//	//	}
//	//	defer apiRes.Body.Close()
//	//
//	//	return
//	//}
//	//
//	//// send to Web01 for run command in SIT
//	//var reqSIT reqSITJSON
//	//reqSIT.Key = lineToken
//	//reqSIT.Command = command(jbotCmd)
//	//postSITBody, err := json.Marshal(reqSIT)
//	//if err != nil {
//	//	log.Println("Cannot marshal JSON:", err)
//	//	sendLineReply(ev0.ReplyToken, "Cannot marshal command in SIT:"+err.Error())
//	//	return
//	//}
//	//log.Println("Request Body:", string(postSITBody))
//	//postSIT, err := http.NewRequest("POST", "https://203.146.225.161/j-bot01/", bytes.NewBuffer(postSITBody))
//	//postSIT.Header.Set("Content-Type", "application/json")
//	//client := &http.Client{
//	//	Timeout: 3 * time.Minute,
//	//	Transport: &http.Transport{
//	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//	//	},
//	//}
//	//sitRes, err := client.Do(postSIT)
//	//if err != nil {
//	//	log.Println("Cannot post to SIT:", err)
//	//	sendLineReply(ev0.ReplyToken, "Cannot run command in SIT:"+err.Error())
//	//	return
//	//}
//	//defer sitRes.Body.Close()
//	//
//	//bodySIT, err := ioutil.ReadAll(sitRes.Body)
//	//if err != nil {
//	//	log.Println("Cannot read SIT body:", err)
//	//	sendLineReply(ev0.ReplyToken, "Cannot get response from SIT:"+err.Error())
//	//	return
//	//}
//	//log.Printf("SIT Body:%s\n", string(bodySIT))
//	//
//	//if sitRes.StatusCode == 404 {
//	//	log.Println("HTTP SIT return 404")
//	//	sendLineReply(ev0.ReplyToken, "SIT service not start")
//	//	return
//	//}
//	//
//	//var resSIT res
//	//err = json.Unmarshal(bodySIT, &resSIT)
//	//if err != nil {
//	//	log.Println("Error when unmarshal SIT:", err)
//	//	sendLineReply(ev0.ReplyToken, "Cannot unmarshal SIT:"+err.Error())
//	//	return
//	//}
//	//log.Printf("SIT Unmarshal result:%+v\n", resSIT)
//	//
//	//if resSIT.Code != "00" {
//	//	msg1 := "Run failed"
//	//	msg2 := resSIT.Message
//	//	sendLineReply(ev0.ReplyToken, msg1, msg2)
//	//	return
//	//}
//	//
//	////msg1 := "Run success, here is the result"
//	//msg2 := resSIT.Message
//	//
//	//err = sendLineReply(ev0.ReplyToken, msg2)
//	log.Println("Success post:", err)
//}
//
//func sendLineReply(replyToken string, messages ...string) error {
//	type Msgs struct {
//		Type string `json:"type"`
//		Text string `json:"text"`
//	}
//	type postJSON struct {
//		ReplyToken string `json:"replyToken"`
//		Message    []Msgs `json:"messages"`
//	}
//	var res postJSON
//	res.ReplyToken = replyToken
//
//	for _, m := range messages {
//		res.Message = append(res.Message, Msgs{"text", m})
//	}
//
//	return postLINEAPI(res, lineApiReplyUrl)
//}
//
//func sendLineSticker(replyToken, packageId, stickerId string) error {
//	type Msgs struct {
//		Type      string `json:"type"`
//		PackageId string `json:"packageId"`
//		StickerId string `json:"stickerId"`
//	}
//	type postJSON struct {
//		ReplyToken string `json:"replyToken"`
//		Message    []Msgs `json:"messages"`
//	}
//	var res postJSON
//	res.ReplyToken = replyToken
//	res.Message = append(res.Message, Msgs{"sticker", packageId, stickerId})
//
//	return postLINEAPI(res, lineApiReplyUrl)
//}
//
//func checkLeave(req reqLINEAPIJSON) bool {
//	ev0 := req.Events[0]
//	if ev0.Message.Type != "text" {
//		return false
//	}
//
//	if (strings.Contains(ev0.Message.Text, "ขอลา")) ||
//		(strings.Contains(ev0.Message.Text, "ลาป่วย")) ||
//		(strings.Contains(ev0.Message.Text, "อนุญาตลา")) ||
//		(strings.Contains(ev0.Message.Text, "อนุญาติลา")) {
//		sendLineSticker(ev0.ReplyToken, "2", "175")
//		//sendLineSticker(ev0.ReplyToken, "1004278", "235263")
//		return true
//	}
//	return false
//}
