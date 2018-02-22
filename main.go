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
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

var allowUsers = map[string]string{
	"U4dba084ef992f2cc6204ccf8e5261ccc": "Esterlism",
	"U48e0dbbce8d4b4fd0690a39ed267484f": "Test",
}

var bot *linebot.Client
var appBaseURL = "https://peaceful-shelf-33227.herokuapp.com"
var downloadDir = "https://peaceful-shelf-33227.herokuapp.com"

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot :", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

//func callbackHandler(w http.ResponseWriter, r *http.Request) {
//	events, err := bot.ParseRequest(r)
//
//	if err != nil {
//		if err == linebot.ErrInvalidSignature {
//			w.WriteHeader(400)
//		} else {
//			w.WriteHeader(500)
//		}
//		return
//	}
//
//	for _, event := range events {
//		if event.Type == linebot.EventTypeMessage {
//			switch message := event.Message.(type) {
//			case *linebot.TextMessage:
//				if strings.EqualFold(message.Text, "test1") {
//
//					leftBtn := linebot.NewMessageTemplateAction("left", "left clicked")
//					rightBtn := linebot.NewMessageTemplateAction("right", "right clicked")
//
//					template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
//					messgage := linebot.NewTemplateMessage("OK YOYO", template)
//
//					if _, err = bot.ReplyMessage(event.ReplyToken, messgage).Do(); err != nil {
//						log.Print(err)
//					}
//
//				} else if strings.EqualFold(message.Text, "test2") {
//
//					messgage := linebot.NewStickerMessage ("2", "175")
//
//					if _, err = bot.ReplyMessage(event.ReplyToken, messgage).Do(); err != nil {
//						log.Print(err)
//					}
//
//				} else if strings.EqualFold(message.Text, "test3") {
//
//					messgage := linebot.NewImageMessage("https://mgronline.com/images/mgr-online-logo.png", "https://mgronline.com/images/mgr-online-logo.png")
//
//					if _, err = bot.ReplyMessage(event.ReplyToken, messgage).Do(); err != nil {
//						log.Print(err)
//					}
//
//				} else if strings.EqualFold(message.Text, "logme") {
//
//					userID := event.Source.UserID
//					groupID := event.Source.GroupID
//					RoomID := event.Source.RoomID
//
//					messgage := "userID = " + userID + ", groupID = " + groupID + ",RoomID = " + RoomID
//
//					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(messgage)).Do(); err != nil {
//						log.Print(err)
//					}
//
//				} else {
//					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.ID+":"+message.Text+" OK! 222")).Do(); err != nil {
//						log.Print(err)
//					}
//				}
//			}
//		}
//	}
//
//}

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
		log.Printf("Got event %v", event)
		switch event.Type {
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if err := handleText(message, event.ReplyToken, event.Source); err != nil {
					log.Print(err)
				}
			case *linebot.ImageMessage:
				if err := handleImage(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.VideoMessage:
				if err := handleVideo(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.AudioMessage:
				if err := handleAudio(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.LocationMessage:
				if err := handleLocation(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				if err := handleSticker(message, event.ReplyToken); err != nil {
					log.Print(err)
				}
			default:
				log.Printf("Unknown message: %v", message)
			}
		case linebot.EventTypeFollow:
			if err := replyText(event.ReplyToken, "Got followed event"); err != nil {
				log.Print(err)
			}
		case linebot.EventTypeUnfollow:
			log.Printf("Unfollowed this bot: %v", event)
		case linebot.EventTypeJoin:
			if err := replyText(event.ReplyToken, "Joined "+string(event.Source.Type)); err != nil {
				log.Print(err)
			}
		case linebot.EventTypeLeave:
			log.Printf("Left: %v", event)
		case linebot.EventTypePostback:
			data := event.Postback.Data
			if data == "DATE" || data == "TIME" || data == "DATETIME" {
				data += fmt.Sprintf("(%v)", *event.Postback)
			}
			if err := replyText(event.ReplyToken, "Got postback: "+data); err != nil {
				log.Print(err)
			}
		case linebot.EventTypeBeacon:
			if err := replyText(event.ReplyToken, "Got beacon: "+event.Beacon.Hwid); err != nil {
				log.Print(err)
			}
		default:
			log.Printf("Unknown event: %v", event)
		}
	}
}

func handleText(message *linebot.TextMessage, replyToken string, source *linebot.EventSource) error {

	//var msgToUserID []linebot.Message
	//msgToUserID[0] = linebot.NewTextMessage("สวัสดี")
	//msgToUserID[1] = linebot.NewTextMessage("จ้า")

	//_, err := bot.PushMessage("U4dba084ef992f2cc6204ccf8e5261ccc", linebot.NewTextMessage("สวัสดี"), linebot.NewTextMessage("จ้า")).Do()
	//if err != nil {
	//	// Do something when some bad happened
	//}

	userID := source.UserID
	//groupID := source.GroupID
	//RoomID := source.RoomID

	if strings.ToLower(message.Text) == "logme" {

		userID := source.UserID
		groupID := source.GroupID
		RoomID := source.RoomID

		//message := "userID = " + userID + ", groupID = " + groupID + ",RoomID = " + RoomID
		message := fmt.Sprintf("userId = %s\ngroupID = %s\nRoomID = %s", userID, groupID, RoomID)

		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	}

	if _, ok := allowUsers[userID]; !ok {
		message := "Sorry, for security reason, only authorized user can run this Bot command\n"
		//for _, v := range allowUsers {
		//	msg += v + ", "
		//}
		//msg += "\n-------------------------\n"
		//msg += ev0.Source.UserId
		//msg += "\n-------------------------"
		//sendLineReply(ev0.ReplyToken, msg)
		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Text), "set ") {
		//SET
		symbol := ReturnStringAfterLastSpace(message.Text)
		symbol = strings.TrimSpace(symbol)
		symbol = strings.ToUpper(symbol)
		message := GetPriceSettrade(symbol)

		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	} else {

		switch strings.ToLower(message.Text) {
		case "profile":
			if source.UserID != "" {
				profile, err := bot.GetProfile(source.UserID).Do()
				if err != nil {
					return replyText(replyToken, err.Error())
				}
				if _, err := bot.ReplyMessage(
					replyToken,
					linebot.NewTextMessage("Display name: "+profile.DisplayName),
					linebot.NewTextMessage("Status message: "+profile.StatusMessage),
				).Do(); err != nil {
					return err
				}
			} else {
				return replyText(replyToken, "Bot can't use profile API without user ID")
			}
		//case "logme":
		//	if source.UserID != "" {
		//
		//		userID := source.UserID
		//		groupID := source.GroupID
		//		RoomID := source.RoomID
		//
		//		message := "userID = " + userID + ", groupID = " + groupID + ",RoomID = " + RoomID
		//
		//		if _, err := bot.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		//			log.Print(err)
		//		}
		//
		//	} else {
		//		return replyText(replyToken, "Bot can't use profile API without user ID")
		//	}
		case "stickerme":

			messgage := linebot.NewStickerMessage("2", "175")

			if _, err := bot.ReplyMessage(replyToken, messgage).Do(); err != nil {
				log.Print(err)
			}

		case "buttons":
			imageURL := appBaseURL + "/static/buttons/1040.jpg"
			template := linebot.NewButtonsTemplate(
				imageURL, "My button sample", "Hello, my button",
				linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
				linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
				linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
				linebot.NewMessageTemplateAction("Say message", "Rice=米"),
			)
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewTemplateMessage("Buttons alt text", template),
			).Do(); err != nil {
				return err
			}
		case "confirm":
			template := linebot.NewConfirmTemplate(
				"Do it?",
				linebot.NewMessageTemplateAction("Yes", "Yes!"),
				linebot.NewMessageTemplateAction("No", "No!"),
			)
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewTemplateMessage("Confirm alt text", template),
			).Do(); err != nil {
				return err
			}
		case "carousel":
			imageURL := appBaseURL + "/static/buttons/1040.jpg"
			template := linebot.NewCarouselTemplate(
				linebot.NewCarouselColumn(
					imageURL, "hoge", "fuga",
					linebot.NewURITemplateAction("Go to line.me", "https://line.me"),
					linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
				),
				linebot.NewCarouselColumn(
					imageURL, "hoge", "fuga",
					linebot.NewPostbackTemplateAction("言 hello2", "hello こんにちは", "hello こんにちは"),
					linebot.NewMessageTemplateAction("Say message", "Rice=米"),
				),
			)
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewTemplateMessage("Carousel alt text", template),
			).Do(); err != nil {
				return err
			}
			//case "image carousel":
			//	imageURL := appBaseURL + "/static/buttons/1040.jpg"
			//	template := linebot.NewImageCarouselTemplate(
			//		linebot.NewImageCarouselColumn(
			//			imageURL,
			//			linebot.NewURITemplateAction("Go to LINE", "https://line.me"),
			//		),
			//		linebot.NewImageCarouselColumn(
			//			imageURL,
			//			linebot.NewPostbackTemplateAction("Say hello1", "hello こんにちは", ""),
			//		),
			//		linebot.NewImageCarouselColumn(
			//			imageURL,
			//			linebot.NewMessageTemplateAction("Say message", "Rice=米"),
			//		),
			//		linebot.NewImageCarouselColumn(
			//			imageURL,
			//			linebot.NewDatetimePickerTemplateAction("datetime", "DATETIME", "datetime", "", "", ""),
			//		),
			//	)
			//	if _, err := bot.ReplyMessage(
			//		replyToken,
			//		linebot.NewTemplateMessage("Image carousel alt text", template),
			//	).Do(); err != nil {
			//		return err
			//	}
			//case "datetime":
			//	template := linebot.NewButtonsTemplate(
			//		"", "", "Select date / time !",
			//		linebot.NewDatetimePickerTemplateAction("date", "DATE", "date", "", "", ""),
			//		linebot.NewDatetimePickerTemplateAction("time", "TIME", "time", "", "", ""),
			//		linebot.NewDatetimePickerTemplateAction("datetime", "DATETIME", "datetime", "", "", ""),
			//	)
			//	if _, err := bot.ReplyMessage(
			//		replyToken,
			//		linebot.NewTemplateMessage("Datetime pickers alt text", template),
			//	).Do(); err != nil {
			//		return err
			//	}
		case "imagemap":
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewImagemapMessage(
					appBaseURL+"/static/rich",
					"Imagemap alt text",
					linebot.ImagemapBaseSize{1040, 1040},
					linebot.NewURIImagemapAction("https://store.line.me/family/manga/en", linebot.ImagemapArea{0, 0, 520, 520}),
					linebot.NewURIImagemapAction("https://store.line.me/family/music/en", linebot.ImagemapArea{520, 0, 520, 520}),
					linebot.NewURIImagemapAction("https://store.line.me/family/play/en", linebot.ImagemapArea{0, 520, 520, 520}),
					linebot.NewMessageImagemapAction("URANAI!", linebot.ImagemapArea{520, 520, 520, 520}),
				),
			).Do(); err != nil {
				return err
			}
		case "bye":
			switch source.Type {
			case linebot.EventSourceTypeUser:
				return replyText(replyToken, "Bot can't leave from 1:1 chat")
			case linebot.EventSourceTypeGroup:
				if err := replyText(replyToken, "Leaving group"); err != nil {
					return err
				}
				if _, err := bot.LeaveGroup(source.GroupID).Do(); err != nil {
					return replyText(replyToken, err.Error())
				}
			case linebot.EventSourceTypeRoom:
				if err := replyText(replyToken, "Leaving room"); err != nil {
					return err
				}
				if _, err := bot.LeaveRoom(source.RoomID).Do(); err != nil {
					return replyText(replyToken, err.Error())
				}
			}
		default:
			log.Printf("Echo message to %s: %s", replyToken, message.Text)
			if _, err := bot.ReplyMessage(
				replyToken,
				linebot.NewTextMessage(message.Text),
			).Do(); err != nil {
				return err
			}
		}
	}

	return nil
}

func handleImage(message *linebot.ImageMessage, replyToken string) error {
	return handleHeavyContent(message.ID, func(originalContent *os.File) error {
		// You need to install ImageMagick.
		// And you should consider about security and scalability.
		previewImagePath := originalContent.Name() + "-preview"
		_, err := exec.Command("convert", "-resize", "240x", "jpeg:"+originalContent.Name(), "jpeg:"+previewImagePath).Output()
		if err != nil {
			return err
		}

		originalContentURL := appBaseURL + "/downloaded/" + filepath.Base(originalContent.Name())
		previewImageURL := appBaseURL + "/downloaded/" + filepath.Base(previewImagePath)
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewImageMessage(originalContentURL, previewImageURL),
		).Do(); err != nil {
			return err
		}
		return nil
	})
}

func handleVideo(message *linebot.VideoMessage, replyToken string) error {
	return handleHeavyContent(message.ID, func(originalContent *os.File) error {
		// You need to install FFmpeg and ImageMagick.
		// And you should consider about security and scalability.
		previewImagePath := originalContent.Name() + "-preview"
		_, err := exec.Command("convert", "mp4:"+originalContent.Name()+"[0]", "jpeg:"+previewImagePath).Output()
		if err != nil {
			return err
		}

		originalContentURL := appBaseURL + "/downloaded/" + filepath.Base(originalContent.Name())
		previewImageURL := appBaseURL + "/downloaded/" + filepath.Base(previewImagePath)
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewVideoMessage(originalContentURL, previewImageURL),
		).Do(); err != nil {
			return err
		}
		return nil
	})
}

func handleAudio(message *linebot.AudioMessage, replyToken string) error {
	return handleHeavyContent(message.ID, func(originalContent *os.File) error {
		originalContentURL := appBaseURL + "/downloaded/" + filepath.Base(originalContent.Name())
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewAudioMessage(originalContentURL, 100),
		).Do(); err != nil {
			return err
		}
		return nil
	})
}

func handleLocation(message *linebot.LocationMessage, replyToken string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewLocationMessage(message.Title, message.Address, message.Latitude, message.Longitude),
	).Do(); err != nil {
		return err
	}
	return nil
}

func handleSticker(message *linebot.StickerMessage, replyToken string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewStickerMessage(message.PackageID, message.StickerID),
	).Do(); err != nil {
		return err
	}
	return nil
}

func replyText(replyToken, text string) error {
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do(); err != nil {
		return err
	}
	return nil
}

func handleHeavyContent(messageID string, callback func(*os.File) error) error {
	content, err := bot.GetMessageContent(messageID).Do()
	if err != nil {
		return err
	}
	defer content.Content.Close()
	log.Printf("Got file: %s", content.ContentType)
	originalConent, err := saveContent(content.Content)
	if err != nil {
		return err
	}
	return callback(originalConent)
}

func saveContent(content io.ReadCloser) (*os.File, error) {
	file, err := ioutil.TempFile(downloadDir, "")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return nil, err
	}
	log.Printf("Saved %s", file.Name())
	return file, nil
}

func ReturnStringAfterLastSpace(s string) string {
	if len(s) > 0 && strings.Count(s, " ") > 0 {
		sRune := string([]rune(s))
		lastSpace := strings.LastIndex(sRune, " ")
		name := sRune[lastSpace:]
		return name
	} else {
		return s
	}

}

func GetPriceSettrade(symbol string) string {
	var answer string

	symbol = strings.TrimSpace(symbol)
	symbol = strings.ToUpper(symbol)

	doc, err := goquery.NewDocument("http://www.settrade.com/C04_02_stock_historical_p1.jsp?txtSymbol=" + symbol + "&ssoPageId=10&selectPage=2")
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	//doc.Find(".col-xs-12 .round-border .content-block").Each(func(i int, s *goquery.Selection) {
	//	// For each item found, get the band and title
	//	price := s.Find("<h1>").Text()
	//	fmt.Printf("%s\n", price)
	//})

	doc.Find(".col-xs-12 .round-border .row .col-xs-6").Each(func(index int, item *goquery.Selection) {
		linkTag := item
		linkText := strings.TrimSpace(linkTag.Text())
		linkText = strings.TrimSpace(ReturnStringAfterLastSpace(linkText))
		//linkText = strings.Replace(linkText, "'", "", 0)
		//fmt.Printf("Link #%d: %s\n", index, linkText)

		switch index {
		case 0:
			answer = answer + fmt.Sprintf("หุ้น %s\n", linkText)
		case 2:
			answer = answer + fmt.Sprintf("ราคาล่าสุด %s\n", linkText)
		case 3:
			answer = answer + fmt.Sprintf("เปลี่ยนแปลง %s\n", linkText)
		case 4:
			answer = answer + fmt.Sprintf("%%เปลี่ยนแปลง %s", linkText)
		}
	})

	//fmt.Printf("%s\n", answer)
	return answer
}
