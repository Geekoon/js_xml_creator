package main

import (
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

var database *sql.DB

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

// Product is an offer comes from base with some common good's elements
type Product struct {
	id          int
	groupID     int
	parentID    int
	nameProduct string
	nameOffer   string
	url         string
	description string
	barcode     string
	uuid        string
	amount      float32
}

// GroupProduct is a folder (from 1C) or group of goods collected together by one feature or BRAND
type GroupProduct struct {
	id       int
	parentID int
	name     string
}

type Offer struct {
	XMLName     xml.Name     `xml:"offer"`
	ID          int          `xml:"id,attr"`
	GroupID     int          `xml:"group_id,attr"`
	UUID        string       `xml:"uuid"`
	URL         string       `xml:"url"`
	Name        string       `xml:"name"`
	CountItems  int          `xml:"countItems"`
	Description *Description `xml:"description"`
}

type Description struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}

type Category struct {
	XMLName  xml.Name `xml:"category"`
	ID       int      `xml:"id,attr"`
	ParentID int      `xml:"parent_id,attr"`
	Name     string   `xml:",chardata"`
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

	rows, err := database.Query("SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, IFNULL(c.content, '') AS content, IFNULL(o.barcode, '') AS barcode, o.id_1c_offer, ob.value FROM tbl_offers AS o LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer WHERE o.act=1 AND o.id_1c_offer != 0 AND ob.id_storage=2 AND ob.value != 0 LIMIT 3000")
	if err != nil {
		ErrorLogger.Println("MySQL Error:", err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.id, &p.groupID, &p.parentID, &p.nameOffer, &p.nameProduct, &p.url, &p.description, &p.barcode, &p.uuid, &p.amount)
		if err != nil {
			ErrorLogger.Println(err)
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
		ErrorLogger.Println("MySQL Error:", err)
	}
	defer rows.Close()

	var gp []GroupProduct

	for rows.Next() {
		p := GroupProduct{}
		err := rows.Scan(&p.id, &p.parentID, &p.name)
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}
		gp = append(gp, p)
	}
	return gp
}

func (s *CategoryArray) AddCategory(sID int, sParentID int, sName string) {
	staffRecord := Category{ID: sID, ParentID: sParentID, Name: sName}
	s.Categories = append(s.Categories, staffRecord)
}

func (s *OfferArray) AddOffer(sID int, sGroupID int, sUUID string, sAmount int, sName string, sURL string, sDescription string) {
	staffRecord := Offer{ID: sID, GroupID: sGroupID, UUID: sUUID, CountItems: sAmount, Name: sName, URL: sURL, Description: &Description{Text: sDescription}}
	s.Offers = append(s.Offers, staffRecord)
}

func main() {

	db, err := sql.Open("mysql", "antor:Yehat13@/js78base")
	if err != nil {
		ErrorLogger.Println(err)
	}

	database = db
	defer db.Close()

	pDB := getProduct()
	catDB := getGroups()

	var v = YmlCatalog{}
	for i := 0; i < len(catDB); i++ {
		v.Shop.Categories.AddCategory(catDB[i].id, catDB[i].parentID, catDB[i].name)
	}
	for i := 0; i < len(pDB); i++ {
		v.Shop.AllOffers.AddOffer(pDB[i].id, pDB[i].groupID, pDB[i].uuid, int(pDB[i].amount), pDB[i].nameOffer, pDB[i].url, pDB[i].description)
	}

	xmlFileName := "offers.xml"
	xmlFile, err := os.Create(xmlFileName)
	if err != nil {
		ErrorLogger.Println("Unable to save XML file:", err)
		os.Exit(1)
	}
	defer xmlFile.Close()
	xmlWriter := io.Writer(xmlFile)

	xmlWriter.Write([]byte(xml.Header))

	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("", "    ")
	if err := enc.Encode(v); err != nil {
		ErrorLogger.Printf("error: %v\n", err)
	}

}
