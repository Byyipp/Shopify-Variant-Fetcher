package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type prod struct {
	Product urlprod `json:"product"`
}

type urlprod struct {
	Title    string     `json:"title"`
	Variants []sitevars `json:"variants"`
	Images   []img      `json:"images"`
}

type sitevars struct {
	ID   int    `json:"id"`
	Size string `json:"option1"`
}

type img struct {
	Src string `json:"src"`
}

func handMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, "-var") {
		if m.Author.ID == BotID {
			return
		}

		str := m.Content
		newString := strings.TrimPrefix(str, "-var ")

		jsonString := newString + ".json"
		resp, err := http.Get(jsonString)
		if err != nil {
			errormsg := make([]*discordgo.MessageEmbed, 0, 1)
			errormsg = append(errormsg, &discordgo.MessageEmbed{
				URL:         newString,
				Title:       newString,
				Description: "Uh Oh Link does not exist\n", //the variants
				Color:       14177041,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Sponsored by Nuggie variants",
					IconURL: "https://media.discordapp.net/attachments/631341893472747520/797945348781768744/tyboWVm.jpg",
				},
			})

			webhook := discordgo.WebhookParams{
				Username:  "Variants",
				Embeds:    errormsg,
				AvatarURL: "https://media.discordapp.net/attachments/631341893472747520/797945348781768744/tyboWVm.jpg",
			}

			_, _ = s.WebhookExecute("<Insert Webhook ID>", "<Insert Webhook Token>", true, &webhook)
			log.Printf("json error: \n", err)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var check prod
		err = json.Unmarshal(body, &check)
		if err != nil {
			log.Fatal(err)
		}

		var varstring string = "**Variants**```\n"
		for i := range check.Product.Variants {
			varstring += check.Product.Variants[i].Size + " - " + strconv.Itoa(check.Product.Variants[i].ID) + "\n"
		}
		varstring += "```"
		msg := make([]*discordgo.MessageEmbed, 0, 1)
		msg = append(msg, &discordgo.MessageEmbed{
			URL:         newString,
			Title:       check.Product.Title,
			Description: varstring, //the variants
			Color:       14177041,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Sponsored by Nuggie variants",
				IconURL: "https://media.discordapp.net/attachments/631341893472747520/797945348781768744/tyboWVm.jpg",
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: check.Product.Images[0].Src,
			},
		})

		webhook := discordgo.WebhookParams{
			Username:  "Variants",
			Embeds:    msg,
			AvatarURL: "https://media.discordapp.net/attachments/631341893472747520/797945348781768744/tyboWVm.jpg",
		}

		_, _ = s.WebhookExecute("<Insert Webhook ID>", "<Insert Webhook Token>", true, &webhook)
		// Get From Webhook URL
	}

}

const token string = "<Insert Bot Token>"
// Insert Bot token from Bot developer site

var BotID string

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	dg.AddHandler(handMessage)

	err = dg.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("bot running..")

	<-make(chan struct{})
	return
}
