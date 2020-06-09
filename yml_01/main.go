package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"os"
)

var database *sql.DB

type Product struct {
	Id   int
	Name string
	Url  string
}
type Offer struct {
	XMLName xml.Name `xml:"offer"`
	Id      int      `xml:"id"`
	Name    string   `xml:"name"`
	Url     string   `xml:"url"`
}

type OfferArray struct {
	Offers []Offer
}

type AllOffers struct {
	XMLName xml.Name `xml:"offers"`
	Offers  OfferArray
}

func getProduct() []Product {
	rows, err := database.Query("select id, name, url from js78base.tbl_core AS c WHERE c.model='ProductItem'")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Name, &p.Url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	return products
}

func (s *OfferArray) AddOffer(sId int, sName string, sUrl string) {
	staffRecord := Offer{Id: sId, Name: sName, Url: sUrl}
	s.Offers = append(s.Offers, staffRecord)
}

func main() {

	db, err := sql.Open("mysql", "antor:Yehat13@/js78base")
	if err != nil {
		log.Println(err)
	}

	database = db
	defer db.Close()

	//fmt.Println(getProduct()[100])

	filename := "ppp.txt"
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Unable to save file:", err)
		os.Exit(1)
	}
	defer file.Close()

	p := getProduct()
	//fmt.Printf("%T", p)
	for i := 0; i < len(p); i++ {
		fmt.Fprintln(file, p[i].Id, p[i].Name)
		//fmt.Println(p)
	}

	var v = AllOffers{}
	for i := 0; i < len(p); i++ {
		v.Offers.AddOffer(p[i].Id, p[i].Name, p[i].Url)
	}

	xmlString, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Println("Error in Marshal: ", err)
	}

	fmt.Printf("%s \n", string(xmlString))

	// everything ok now, write to file.
	xmlFileName := "offers.xml"
	xmlFile, err := os.Create(xmlFileName)
	if err != nil {
		fmt.Println("Unable to save XML file:", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	xmlWriter := io.Writer(xmlFile)

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
