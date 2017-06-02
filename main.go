package main

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/jasonlvhit/gocron"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://www.stm.info/fr/infos/etat-du-service/metro")

	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.content-block section.item").Each(func(i int, s *goquery.Selection) {
		couleur := s.Find("h2").Text()
		couleur = strings.Trim(couleur, " ")
		couleur = strings.Replace(couleur, "\n", "", -1)
		etat := s.Find("p").Text()
		etat = strings.Trim(etat, " ")
		etat = strings.Replace(etat, "\n", "", -1)
		if strings.Contains(etat, "Interruption de service") {
			color.Red("%s - %s\n", couleur, etat)
			sendmailInteruption(couleur+etat, couleur)
		} else {
			fmt.Printf("%s - %s\n", couleur, etat)
		}
	})
}

func sendmailInteruption(body string, couleur string) {
	from := "stminfogo@gmail.com"
	pass := "allostminfo"
	to := "fprieur@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Interruption sur la ligne " + couleur + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("mail sent")
}

func main() {

	ExampleScrape()

	gocron.Every(1).Minute().Do(ExampleScrape)
	<-gocron.Start()
}
