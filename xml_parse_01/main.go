package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var database *sql.DB

type Product struct {
	id   int
	name string
	url  string
}

type GroupProduct struct {
	id   int
	name string
}

type Offer struct {
	XMLName xml.Name `xml:"offer"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:"name"`
	Url     string   `xml:"url"`
}

type Category struct {
	XMLName xml.Name `xml:"category"`
	Id      int      `xml:"id,attr"`
	Name    string   `xml:",chardata"`
}

type OfferArray struct {
	XMLName xml.Name `xml:"offers"`
	Offers  []Offer
}

type CategoryArray struct {
	XMLName    xml.Name `xml:"categories"`
	Categories []Category
}

type TagShop struct {
	XMLName    xml.Name `xml:"shop"`
	Categories CategoryArray
	AllOffers  OfferArray
}

type YmlCatalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Shop    TagShop
}

func getProduct() []Product {
	//	rows, err := database.Query("select id, name, url from js78base.tbl_core AS c WHERE c.model='ProductItem'")
	rows, err := database.Query("SELECT barcode, name, id_product_item FROM js78base.tbl_offers AS o WHERE o.barcode IS NOT NULL AND o.id_1c_offer > 0 AND act=1 LIMIT 500")

	if err != nil {
		log.Println("MySQL Error:", err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.id, &p.name, &p.url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	return products
}

func main() {
	xmlFileName := "import.xml"
	xmlFile, err := os.Open("data/" + xmlFileName)
	if err != nil {
		fmt.Println("Unable to open XML file:", err)
		os.Exit(1)
	}
	defer xmlFile.Close()

	/*	enc := xml.NewEncoder(xmlWriter)
		enc.Indent("", "    ")
		if err := enc.Encode(v); err != nil {
			fmt.Printf("error: %v\n", err)
		}


		db, err := sql.Open("mysql", "antor:Yehat13@/js78base")
		if err != nil {
			log.Println(err)
		}

		database = db
		defer db.Close()

		pDB := getProduct()

		/*	filename := "ppp.txt"
			file, err := os.Create(filename)
			if err != nil {
				fmt.Println("Unable to save file:", err)
				os.Exit(1)
			}
			defer file.Close()
			for i := 0; i < len(catDB); i++ {
				fmt.Fprintln(file, catDB[i].id, catDB[i].name)
			}
	*/
	//	var v = YmlCatalog{}
	/*	xmlString, err := xml.MarshalIndent(v, "", "    ")
		if err != nil {
			fmt.Println("Error in MarshalIndent: ", err)
		}
		fmt.Printf("%s \n", string(xmlString))
	*/

}
