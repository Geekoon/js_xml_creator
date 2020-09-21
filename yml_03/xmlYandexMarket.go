package main

import (
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"os"
	"time"

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
	id          int
	groupID     int
	parentID    int
	nameProduct string
	nameOffer   string
	url         string
	description string
	amount      int
	price       float32
	brand       string
	kind        string
	structure   string
	sizeType    string
	sex         int
	age         int
	size        string
	color       string
	sizeID      string
	colorID     string
	parentRGB   string
}

// GroupProduct is a folder (from 1C) or group of goods collected together by one feature or BRAND
type GroupProduct struct {
	id       int
	parentID string
	name     string
}

type Offer struct {
	XMLName    xml.Name `xml:"offer"`
	ID         int      `xml:"id,attr"`
	Available  string   `xml:"available,attr"`
	Type       string   `xml:"type,attr"`
	GroupID    int      `xml:"group_id,attr"`
	Name       string   `xml:"model"`
	Brand      string   `xml:"vendor"`
	URL        string   `xml:"url"`
	CategoryID int      `xml:"categoryId"`
	//CountItems  int                 `xml:"countItems"`
	Price       float32             `xml:"price"`
	CurrencyID  string              `xml:"currencyId"`
	Kind        string              `xml:"typePrefix"`
	ImageURL    string              `xml:"picture"`
	Description *Description        `xml:"description"`
	Parametres  *[numberParam]Param `xml:"param"`
}

type Description struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}

type Param struct {
	Name    string `xml:"name,attr"`
	AddAttr string `xml:"unit,attr,omitempty"`
	Text    string `xml:",chardata"`
}

type Category struct {
	XMLName  xml.Name `xml:"category"`
	ID       int      `xml:"id,attr"`
	ParentID string   `xml:"parentId,attr,omitempty"`
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

type Currencies struct {
	XMLName  xml.Name `xml:"currencies"`
	Currency struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
		Rate string `xml:"rate,attr"`
	} `xml:"currency"`
}

type TagShop struct {
	XMLName    xml.Name `xml:"shop"`
	Name       string   `xml:"name"`
	Company    string   `xml:"company"`
	URL        string   `xml:"url"`
	Author     string   `xml:"agency"`
	Email      string   `xml:"email"`
	Currencies Currencies
	Categories CategoryArray
	Delivery   string `xml:"local_delivery_cost"`
	AllOffers  OfferArray
}

type YmlCatalog struct {
	XMLName  xml.Name `xml:"yml_catalog"`
	DataTime string   `xml:"date,attr"`
	Shop     TagShop
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

	var (
		id    int
		value string
	)

	for rows.Next() {
		err := rows.Scan(&id, &value)
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}
		newList[id] = value
	}

	return newList
}

func getListSS(q string) map[string]string {
	rows, err := database.Query(q)
	if err != nil {
		ErrorLogger.Println("MySQL in getListProperty:", err)
	}
	defer rows.Close()

	newList := make(map[string]string)

	var name, value string

	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}
		newList[value] = name
	}

	return newList
}

type Images struct {
	productID int
	color     string
	URL       string
}

func getImagesURL() []Images {
	q := "SELECT pim.id_product_item, pim.color, IFNULL(pim.image_raw, '') " +
		"FROM tbl_product_images AS pim WHERE pim.act = 1 AND pim.color IS NOT NULL "
		//"GROUP BY pim.color"

	rows, err := database.Query(q)
	if err != nil {
		ErrorLogger.Println("MySQL in getImagesURL:", err)
	}
	defer rows.Close()

	newList := make([]Images, 0)
	for rows.Next() {
		var el Images
		err := rows.Scan(&el.productID, &el.color, &el.URL)
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}
		newList = append(newList, el)
	}

	return newList
}

func getProduct() []Product {
	q := "SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, IFNULL(c.content, ''), " +
		"CAST(ob.value AS UNSIGNED), pr.value, pid.size_type, " +
		"pid.id_property_sex, pid.id_property_age, IFNULL(pid.structure, ''), pa.name, pik.name, " +
		"MAX(CASE WHEN fv.id_feature=1 THEN fv.value END) AS size, " +
		"MAX(CASE WHEN fv.id_feature=2 THEN fv.value END) AS color, " +
		"IFNULL(MAX(CASE WHEN fv.id_feature=2 THEN fv.parent_color END), 'multi') AS parentColor, " +
		"MAX(CASE WHEN fv.id_feature=1 THEN fv.id END) AS sizeID, " +
		"MAX(CASE WHEN fv.id_feature=2 THEN fv.id END) AS colorID " +
		"FROM tbl_offers AS o LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id " +
		"LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer " +
		"LEFT OUTER JOIN tbl_offer_prices AS pr ON o.id = pr.id_offer " +
		"LEFT OUTER JOIN tbl_product_item_detail AS pid ON c.id = pid.id_product_item " +
		"LEFT OUTER JOIN tbl_product_articles AS pa ON pid.id_article = pa.id " +
		"LEFT OUTER JOIN tbl_product_item_kind AS pik ON pik.id = pid.kind_id " +
		"LEFT OUTER JOIN tbl_offer_features AS of ON o.id = of.id_offer " +
		"LEFT OUTER JOIN tbl_feature_values AS fv ON of.id_feature_value = fv.id " +
		"WHERE c.act=1 AND o.act=1 AND o.id_1c_offer != 0 AND ob.id_storage=2 AND ob.value != 0 AND pr.id_price = 3 GROUP BY o.id" // LIMIT 3000"
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
			&p.amount,
			&p.price,
			&p.sizeType,
			&p.sex,
			&p.age,
			&p.structure,
			&p.brand,
			&p.kind,
			&p.size,
			&p.color,
			&p.parentRGB,
			&p.sizeID,
			&p.colorID,
		)
		if err != nil {
			WarningLogger.Println("offer name:", p.nameOffer, "Err:", err)
			continue
		}
		p.amount = changeAmount(p.amount)
		p.url = filialURL + p.url + "?color=" + p.colorID + "&size=" + p.sizeID
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

		if p.parentID == "2" {
			p.parentID = ""
		}

		gp = append(gp, p)
	}
	return gp
}

func (s *CategoryArray) AddCategory(sID int, sParentID string, sName string) {
	staffRecord := Category{ID: sID, ParentID: sParentID, Name: sName}
	s.Categories = append(s.Categories, staffRecord)
}

func (s *OfferArray) AddOffer(
	sID int,
	sGroupID int,
	sAmount int,
	sName string,
	sURL string,
	sPrice float32,
	sParentID int,
	sBrand string,
	sKind string,
	sImage string,
	sDescription string,
	sParametres [numberParam]Param,
) {
	staffRecord := Offer{
		ID:        sID,
		Available: "true",
		Type:      "vendor.model",
		GroupID:   sGroupID,
		//CountItems:  sAmount,
		Name:        sName,
		Brand:       sBrand,
		URL:         sURL,
		Price:       sPrice,
		CurrencyID:  "RUR",
		CategoryID:  sParentID,
		Kind:        sKind,
		ImageURL:    sImage,
		Description: &Description{Text: sDescription},
		Parametres:  &sParametres,
	}
	s.Offers = append(s.Offers, staffRecord)
}

func main() {

	//db, err := sql.Open("mysql", "root:pass123@/js78base")
	db, err := sql.Open("mysql", "admitex:8E5s3T7y2Y0w2W5y@/js2base")
	if err != nil {
		ErrorLogger.Println(err)
	}

	database = db
	defer db.Close()

	pDB := getProduct()
	catDB := getGroups()
	propertyDB := getListProperty("SELECT id, name FROM tbl_property_values WHERE act=1")
	realColorDB := getListSS("SELECT name, value FROM tbl_reference WHERE model='ParentColorItem'")
	imagesDB := getImagesURL()

	numberOffers := len(pDB)
	InfoLogger.Println("Found offers:", numberOffers)

	dtnow := time.Now().Format("2006-01-02 15:04:05")
	v := &YmlCatalog{DataTime: dtnow}

	v.Shop.Name = "JS-Company"
	v.Shop.Company = "ДжиЭс Групп - Москва"
	v.Shop.URL = filialURL
	v.Shop.Author = "A. Orlovskikh"
	v.Shop.Email = "js-admin@mail.ru"

	v.Shop.Currencies.Currency.ID = "RUR"
	v.Shop.Currencies.Currency.Rate = "1"

	for i := 0; i < len(catDB); i++ {
		v.Shop.Categories.AddCategory(catDB[i].id, catDB[i].parentID, catDB[i].name)
	}

	v.Shop.Delivery = "0"

	for i := 0; i < numberOffers; i++ {
		var props [numberParam]Param
		props[0].Name = "Пол"
		props[0].Text = propertyDB[pDB[i].sex]
		props[1].Name = "Возраст"
		props[1].Text = propertyDB[pDB[i].age]
		props[2].Name = "Состав"
		props[2].Text = pDB[i].structure
		props[3].Name = "Размер"
		props[3].Text = pDB[i].size
		props[3].AddAttr = pDB[i].sizeType
		props[4].Name = "Цвет"
		prgb := pDB[i].parentRGB
		if prgb == "multi" || prgb == "" {
			props[4].Text = "Разноцветный"
		} else {
			props[4].Text = realColorDB[prgb]
		}
		props[5].Name = "Color name"
		props[5].Text = pDB[i].color

		imgURL := ""
		for k := 0; k < len(imagesDB); k++ {
			if imagesDB[k].productID == pDB[i].groupID && (imagesDB[k].color == pDB[i].color) {
				imgURL = filialURL + imagesDB[k].URL
			}
		}

		v.Shop.AllOffers.AddOffer(
			pDB[i].id,
			pDB[i].groupID,
			int(pDB[i].amount),
			pDB[i].nameOffer,
			pDB[i].url,
			pDB[i].price,
			pDB[i].parentID,
			pDB[i].brand,
			pDB[i].kind,
			imgURL,
			pDB[i].description,
			props,
		)
	}

	xmlFileName := "shop_moscow.yml"
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
