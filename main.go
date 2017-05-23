package main

import (
	"fmt"
	"log"
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
		} else {
			fmt.Printf("%s - %s\n", couleur, etat)
		}
	})
}

func main() {
	gocron.Every(1).Minute().Do(ExampleScrape)
	<-gocron.Start()
}
