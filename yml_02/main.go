package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

type Product struct {
	id          int
	group_id    int
	parent_id   int
	nameProduct string
	nameOffer   string
	url         string
	description string
	barcode     int
	uuid        string
	amount      float32
}

type GroupProduct struct {
	id        int
	parent_id int
	name      string
}

type Offer struct {
	XMLName    xml.Name `xml:"offer"`
	Id         int      `xml:"id,attr"`
	Group_id   int      `xml:"group_id,attr"`
	Uuid       string   `xml:"uuid"`
	Url        string   `xml:"url"`
	Name       string   `xml:"name"`
	CountItems int      `xml:"countItems"`
}

type Category struct {
	XMLName   xml.Name `xml:"category"`
	Id        int      `xml:"id,attr"`
	Parent_id int      `xml:"parent_id,attr"`
	Name      string   `xml:",chardata"`
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
	// rows, err := database.Query("select id, name, url from js78base.tbl_core AS c WHERE c.model='ProductItem'")
	// rows, err := database.Query("SELECT barcode, name, id_product_item FROM tbl_offers AS o WHERE o.barcode IS NOT NULL AND o.id_1c_offer > 0 AND act=1 LIMIT 200")
	rows, err := database.Query("SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, c.content, o.barcode, o.id_1c_offer, ob.value FROM tbl_offers AS o LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer WHERE o.act=1 AND o.id_1c_offer != 0 AND ob.id_storage=2 AND ob.value != 0 LIMIT 3000")
	if err != nil {
		log.Println("MySQL Error:", err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.id, &p.group_id, &p.parent_id, &p.nameOffer, &p.nameProduct, &p.url, &p.description, &p.barcode, &p.uuid, &p.amount)
		if err != nil {
			//fmt.Println(err)
			continue
		}
		products = append(products, p)
	}
	return products
}

func getGroups() []GroupProduct {
	//rows, err := database.Query("select id, name from js78base.tbl_product_item_kind AS pk")
	rows, err := database.Query("SELECT c.id, c.parent_id, c.name FROM tbl_core AS c WHERE model='ProductGroup' AND act=1")
	if err != nil {
		log.Println("MySQL Error:", err)
	}
	defer rows.Close()

	var gp []GroupProduct

	for rows.Next() {
		p := GroupProduct{}
		err := rows.Scan(&p.id, &p.parent_id, &p.name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		gp = append(gp, p)
	}
	return gp
}

func (s *CategoryArray) AddCategory(sId int, sParent_id int, sName string) {
	staffRecord := Category{Id: sId, Parent_id: sParent_id, Name: sName}
	s.Categories = append(s.Categories, staffRecord)
}

func (s *OfferArray) AddOffer(sId int, sGroup_id int, sUuid string, sAmount int, sName string, sUrl string) {
	staffRecord := Offer{Id: sId, Group_id: sGroup_id, Uuid: sUuid, CountItems: sAmount, Name: sName, Url: sUrl}
	s.Offers = append(s.Offers, staffRecord)
}

func main() {

	db, err := sql.Open("mysql", "antor:Yehat13@/js78base")
	if err != nil {
		log.Println(err)
	}

	database = db
	defer db.Close()

	pDB := getProduct()
	catDB := getGroups()

	var v = YmlCatalog{}
	for i := 0; i < len(catDB); i++ {
		v.Shop.Categories.AddCategory(catDB[i].id, catDB[i].parent_id, catDB[i].name)
	}
	for i := 0; i < len(pDB); i++ {
		v.Shop.AllOffers.AddOffer(pDB[i].id, pDB[i].group_id, pDB[i].uuid, int(pDB[i].amount), pDB[i].nameOffer, pDB[i].url)
	}

	xmlFileName := "offers.xml"
	xmlFile, err := os.Create(xmlFileName)
	if err != nil {
		fmt.Println("Unable to save XML file:", err)
		os.Exit(1)
	}
	defer xmlFile.Close()
	xmlWriter := io.Writer(xmlFile)

	xmlWriter.Write([]byte(xml.Header))

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "    ")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("error: %v\n", err)
	}

}
