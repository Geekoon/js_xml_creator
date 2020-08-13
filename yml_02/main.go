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

const (
	filialURL   = "https://moscow.js-company.ru/"
	numberParam = 6
)

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
	id           int
	groupID      int
	parentID     int
	nameProduct  string
	nameOffer    string
	url          string
	description  string
	barcode      string
	uuid         string
	amount       int
	price        float32
	whosalePrice float32
	code         string
	brand        string
	kind         string
	structure    string
	sex          int
	age          int
	size         string
	color        string
	rgb          string
}

// GroupProduct is a folder (from 1C) or group of goods collected together by one feature or BRAND
type GroupProduct struct {
	id       int
	parentID int
	name     string
}

type Offer struct {
	XMLName       xml.Name            `xml:"offer"`
	ID            int                 `xml:"id,attr"`
	Available     string              `xml:"available,attr"`
	Type          string              `xml:"type,attr"`
	GroupID       int                 `xml:"group_id,attr"`
	Name          string              `xml:"model"`
	Brand         string              `xml:"vendor"`
	URL           string              `xml:"url"`
	CategoryID    int                 `xml:"categoryId"`
	CountItems    int                 `xml:"countItems"`
	Price         float32             `xml:"price"`
	WhosalePrice  float32             `xml:"whosaleprice"`
	RealBarcode   string              `xml:"realBarCode"`
	ProductCode1C string              `xml:"productCode1C"`
	UUID          string              `xml:"uuid"`
	Kind          string              `xml:"typePrefix"`
	Description   *Description        `xml:"description"`
	Parametres    *[numberParam]Param `xml:"param"`
}

type Description struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}

type Param struct {
	Name string `xml:"name,attr"`
	Text string `xml:",chardata"`
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

func changeAmount(m int) int {
	switch {
	case m < 20:
		return m
	case m < 100:
		return (m / 10) * 10
	default:
		return 100
	}
}

func getListProperty(q string) map[int]string {
	rows, err := database.Query(q)
	if err != nil {
		ErrorLogger.Println("MySQL in getListProperty:", err)
	}
	defer rows.Close()

	newList := make(map[int]string)

	for rows.Next() {
		id := 0
		value := ""
		err := rows.Scan(&id, &value)
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}
		newList[id] = value
	}

	return newList
}

func getProduct() []Product {
	q := "SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, IFNULL(c.content, ''), IFNULL(o.barcode, ''), " +
		"o.id_1c_offer, CAST(ob.value AS UNSIGNED), SUM(IF(pr.id_price=1, pr.value, NULL)),	SUM(IF(pr.id_price=3, pr.value, NULL)), " +
		"pid.code, pid.id_property_sex, pid.id_property_age, IFNULL(pid.structure, ''), pib.name, pik.name, " +
		"MAX(CASE WHEN fv.id_feature=1 THEN fv.value END) AS size, " +
		"MAX(CASE WHEN fv.id_feature=2 THEN fv.value END) AS color, " +
		"IFNULL(MAX(CASE WHEN fv.id_feature=2 THEN fv.rgb END), '') AS rgb " +
		"FROM tbl_offers AS o LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id " +
		"LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer " +
		"LEFT OUTER JOIN tbl_offer_prices AS pr ON o.id = pr.id_offer " +
		"LEFT OUTER JOIN tbl_product_item_detail AS pid ON c.id = pid.id_product_item " +
		"LEFT OUTER JOIN tbl_product_item_brand AS pib ON pid.brand_id = pib.id " +
		"LEFT OUTER JOIN tbl_product_item_kind AS pik ON pik.id = pid.kind_id " +
		"LEFT OUTER JOIN tbl_offer_features AS of ON o.id = of.id_offer " +
		"LEFT OUTER JOIN tbl_feature_values AS fv ON of.id_feature_value = fv.id " +
		"WHERE c.act=1 AND o.act=1 AND o.id_1c_offer != 0 AND ob.id_storage=2 AND ob.value != 0 AND pr.id_price != 2 GROUP BY o.id LIMIT 1000"
	rows, err := database.Query(q)
	if err != nil {
		ErrorLogger.Println("MySQL in getProduct:", err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		p := Product{}
		err := rows.Scan(
			&p.id,
			&p.groupID,
			&p.parentID,
			&p.nameOffer,
			&p.nameProduct,
			&p.url,
			&p.description,
			&p.barcode,
			&p.uuid,
			&p.amount,
			&p.whosalePrice,
			&p.price,
			&p.code,
			&p.sex,
			&p.age,
			&p.structure,
			&p.brand,
			&p.kind,
			&p.size,
			&p.color,
			&p.rgb,
		)
		if err != nil {
			WarningLogger.Println("offer name:", p.nameOffer, "Err:", err)
			continue
		}
		p.amount = changeAmount(p.amount)
		p.url = filialURL + p.url
		products = append(products, p)
	}
	return products
}

func getGroups() []GroupProduct {
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

func (s *OfferArray) AddOffer(
	sID int,
	sGroupID int,
	sUUID string,
	sAmount int,
	sName string,
	sURL string,
	sPrice float32,
	sWhosalePrice float32,
	sParentID int,
	sBarcode string,
	sCode1C string,
	sBrand string,
	sKind string,
	sDescription string,
	sParametres [numberParam]Param,
) {
	staffRecord := Offer{
		ID:            sID,
		Available:     "true",
		Type:          "vendor.model",
		GroupID:       sGroupID,
		UUID:          sUUID,
		CountItems:    sAmount,
		Name:          sName,
		Brand:         sBrand,
		URL:           sURL,
		Price:         sPrice,
		WhosalePrice:  sWhosalePrice,
		CategoryID:    sParentID,
		RealBarcode:   sBarcode,
		ProductCode1C: sCode1C,
		Kind:          sKind,
		Description:   &Description{Text: sDescription},
		Parametres:    &sParametres,
	}
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
	propertyDB := getListProperty("SELECT id, name FROM tbl_property_values WHERE act=1")
	//	sizeDB := getListProperty("SELECT id, value FROM tbl_feature_values WHERE id_feature=1")
	//	colorDB := getListProperty("SELECT id, value FROM tbl_feature_values WHERE id_feature=2")
	//	rgbDB := getListProperty("SELECT id, rgb FROM tbl_feature_values WHERE id_feature=2")

	var v = YmlCatalog{}
	for i := 0; i < len(catDB); i++ {
		v.Shop.Categories.AddCategory(catDB[i].id, catDB[i].parentID, catDB[i].name)
	}

	InfoLogger.Println("Found offers:", len(pDB))

	for i := 0; i < len(pDB); i++ {
		var props [numberParam]Param
		props[0].Name = "Пол"
		props[0].Text = propertyDB[pDB[i].sex]
		props[1].Name = "Возраст"
		props[1].Text = propertyDB[pDB[i].age]
		props[2].Name = "Состав"
		props[2].Text = pDB[i].structure
		props[3].Name = "Размер"
		props[3].Text = pDB[i].size
		props[4].Name = "Цвет"
		props[4].Text = pDB[i].color
		props[5].Name = "RGB"
		props[5].Text = pDB[i].rgb

		v.Shop.AllOffers.AddOffer(
			pDB[i].id,
			pDB[i].groupID,
			pDB[i].uuid,
			int(pDB[i].amount),
			pDB[i].nameOffer,
			pDB[i].url,
			pDB[i].price,
			pDB[i].whosalePrice,
			pDB[i].parentID,
			pDB[i].barcode,
			pDB[i].code,
			pDB[i].brand,
			pDB[i].kind,
			pDB[i].description,
			props,
		)
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
		ErrorLogger.Printf("Encoding err: %v\n", err)
	}

}
