package main

import (
	"database/sql"
	"encoding/xml"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/gcfg.v1"

	_ "github.com/go-sql-driver/mysql"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger
)

var (
	database           *sql.DB
	cfg                Config
	pathResult, DBpath string
)

const (
	numberParam = 6
)

func init() {
	filelog, err := os.OpenFile("xml4customers.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Couldn't open log file 'xml4customers.log'!", err)
	}

	InfoLogger = log.New(filelog, "INFO: ", log.Ldate|log.Lmicroseconds)
	WarningLogger = log.New(filelog, "WARNING: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	ErrorLogger = log.New(filelog, "ERROR: ", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	FatalLogger = log.New(filelog, "FATAL: ", log.Ldate|log.Lmicroseconds)

	InfoLogger.Println("############### Starting new session ###############")

	err = gcfg.ReadFileInto(&cfg, "config.gcfg")
	if err != nil {
		FatalLogger.Println("Couldn't open config file! Exit...", err)
		os.Exit(1)
	}

	if cfg.MainSection.DBname == "" || cfg.MainSection.Username == "" || cfg.MainSection.Passuser == "" {
		FatalLogger.Println("Wrong parameter in config file! Exit...")
		os.Exit(1)
	}

	if cfg.MainSection.Key != 123 {
		FatalLogger.Println("Key has expired! Purchase a new one.")
		os.Exit(1)
	}

	if cfg.MainSection.XmlPath == "" {
		pathResult = "xml_result/"
	} else {
		pathResult = cfg.MainSection.XmlPath
	}

	DBpath = cfg.MainSection.Username + ":" + cfg.MainSection.Passuser + "@/" + cfg.MainSection.DBname

}

type Config struct {
	MainSection struct {
		DBname   string
		Username string
		Passuser string
		Key      int
		XmlPath  string
	}
	Client struct {
		Documentation string
	}
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
	amount      int
	code        string
	brand       string
	kind        string
	structure   string
	sex         int
	age         int
}

// GroupProduct is a folder (from 1C) or group of goods collected together by one feature or BRAND
type GroupProduct struct {
	id       int
	parentID string
	name     string
}

type Offer struct {
	XMLName        xml.Name            `xml:"offer"`
	ID             int                 `xml:"id,attr"`
	Available      string              `xml:"available,attr"`
	GroupID        int                 `xml:"group_id,attr"`
	Name           string              `xml:"model"`
	Brand          string              `xml:"vendor"`
	URL            string              `xml:"url"`
	CategoryID     int                 `xml:"categoryId"`
	CountItems     int                 `xml:"countItems"`
	Price          float32             `xml:"price"`
	WhosalePrice   float32             `xml:"whosalePrice"`
	WholesalePrice float32             `xml:"wholesalePrice"`
	CurrencyID     string              `xml:"currencyId"`
	RealBarcode    string              `xml:"realBarCode"`
	ProductCode1C  string              `xml:"productCode1C"`
	UUID           string              `xml:"uuid"`
	Kind           string              `xml:"typePrefix"`
	ImageURL       *[]string           `xml:"picture"`
	BoxImageURL    *[]string           `xml:"boxImage"`
	Description    *Description        `xml:"description"`
	Parametres     *[numberParam]Param `xml:"param"`
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
	ParentID string   `xml:"parent_id,attr,omitempty"`
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
	DocsXML    string   `xml:"documentation"`
	Currencies Currencies
	Categories CategoryArray
	AllOffers  OfferArray
}

type YmlCatalog struct {
	XMLName      xml.Name `xml:"yml_catalog"`
	DataTime     string   `xml:"datetime,attr"`
	NumberOffers int      `xml:"number_offers,attr"`
	Author       string   `xml:"author"`
	Email        string   `xml:"email"`
	Shop         TagShop
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

func getListPrice(q string) map[int]float32 {
	rows, err := database.Query(q)
	if err != nil {
		ErrorLogger.Println("MySQL in getListProperty:", err)
	}
	defer rows.Close()

	newList := make(map[int]float32)

	var (
		id    int
		value float32
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

type Images struct {
	productID int
	color     string
	URL       string
}

type BoxImages struct {
	productID int
	URL       string
}

func getImagesURL() []Images {
	q := "SELECT pim.id_product_item, pim.color, IFNULL(pim.image_raw, '') " +
		"FROM tbl_product_images AS pim WHERE pim.act = 1 AND pim.color IS NOT NULL "

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

func getBoxImagesURL() []BoxImages {
	q := "SELECT pim.id_product_item, pim.image_raw " +
		"FROM tbl_product_images AS pim WHERE pim.act = 1 AND pim.color IS NULL AND pim.image_raw IS NOT NULL "

	rows, err := database.Query(q)
	if err != nil {
		ErrorLogger.Println("MySQL in getImagesURL:", err)
	}
	defer rows.Close()

	newList := make([]BoxImages, 0)
	for rows.Next() {
		var el BoxImages
		err := rows.Scan(&el.productID, &el.URL)
		if err != nil {
			ErrorLogger.Println(err)
			continue
		}
		newList = append(newList, el)
	}

	return newList
}

func getProduct(num string, fn string) []Product {
	q := "SELECT o.id, c.id, c.parent_id, c.name, o.name, c.url, IFNULL(c.content, ''), IFNULL(o.barcode, ''), " +
		"o.id_1c_offer, CAST(ob.value AS SIGNED), " +
		"pid.code, pid.id_property_sex, pid.id_property_age, IFNULL(pid.structure, ''), pa.name, pik.name " +
		"FROM tbl_offers AS o LEFT OUTER JOIN tbl_core AS c ON o.id_product_item = c.id " +
		"LEFT OUTER JOIN tbl_offer_balance AS ob ON o.id = ob.id_offer " +
		"LEFT OUTER JOIN tbl_product_item_detail AS pid ON c.id = pid.id_product_item " +
		"LEFT OUTER JOIN tbl_product_articles AS pa ON pid.id_article = pa.id " +
		"LEFT OUTER JOIN tbl_product_item_kind AS pik ON pik.id = pid.kind_id " +
		"WHERE c.act=1 AND o.act=1 AND o.id_1c_offer != '00000000-0000-0000-0000-000000000000' AND ob.id_storage=" + num + " AND ob.value != 0 " //LIMIT 1500"
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
			&p.code,
			&p.sex,
			&p.age,
			&p.structure,
			&p.brand,
			&p.kind,
		)
		if err != nil {
			WarningLogger.Println("offer name:", p.nameOffer, "Err:", err)
			continue
		}
		p.amount = changeAmount(p.amount)
		p.url = "https://" + fn + ".js-company.ru/" + p.url
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
			WarningLogger.Println(err)
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
	sImage *[]string,
	sBoxImage *[]string,
	sDescription string,
	sParametres [numberParam]Param,
) {
	staffRecord := Offer{
		ID:             sID,
		Available:      "true",
		GroupID:        sGroupID,
		UUID:           sUUID,
		CountItems:     sAmount,
		Name:           sName,
		Brand:          sBrand,
		URL:            sURL,
		Price:          sPrice,
		WhosalePrice:   sWhosalePrice,
		WholesalePrice: sWhosalePrice,
		CurrencyID:     "RUR",
		CategoryID:     sParentID,
		RealBarcode:    sBarcode,
		ProductCode1C:  sCode1C,
		Kind:           sKind,
		ImageURL:       sImage,
		BoxImageURL:    sBoxImage,
		Description:    &Description{Text: sDescription},
		Parametres:     &sParametres,
	}
	s.Offers = append(s.Offers, staffRecord)
}

func main() {

	db, err := sql.Open("mysql", DBpath)
	if err != nil {
		ErrorLogger.Println(err)
	}

	database = db
	defer db.Close()

	catDB := getGroups()
	propertyDB := getListProperty("SELECT id, name FROM tbl_property_values WHERE act=1")
	filialsDB := getListProperty("SELECT id, slug FROM tbl_storage WHERE slug <> '' AND in_yml_api = 1")
	imagesDB := getImagesURL()
	boxImagesDB := getBoxImagesURL()
	sizeDB := getListProperty("SELECT of.id_offer, fv.value FROM tbl_offer_features AS of LEFT OUTER JOIN tbl_feature_values AS fv ON of.id_feature_value = fv.id WHERE fv.id_feature=1")
	colorDB := getListProperty("SELECT of.id_offer, fv.value FROM tbl_offer_features AS of LEFT OUTER JOIN tbl_feature_values AS fv ON of.id_feature_value = fv.id WHERE fv.id_feature=2")
	rgbDB := getListProperty("SELECT of.id_offer, IFNULL(fv.rgb, '') FROM tbl_offer_features AS of LEFT OUTER JOIN tbl_feature_values AS fv ON of.id_feature_value = fv.id WHERE fv.id_feature=2")
	regularPriceDB := getListPrice("SELECT op.id_offer, op.value FROM tbl_offer_prices AS op WHERE op.id_price=3")
	salesPriceDB := getListPrice("SELECT op.id_offer, op.value FROM tbl_offer_prices AS op WHERE op.id_price=1")

	for j, fn := range filialsDB {

		pDB := getProduct(strconv.Itoa(j), fn)
		numberOffers := len(pDB)
		InfoLogger.Println("Found offers in "+fn+":", numberOffers)

		dtnow := time.Now().Format("2006-01-02 15:04:05")
		v := &YmlCatalog{DataTime: dtnow, NumberOffers: numberOffers, Author: "A. Orlovskikh", Email: "js-admin@mail.ru"}

		v.Shop.Name = "JS-Company"
		v.Shop.Company = "ООО 'ДжиЭс Групп'"
		v.Shop.URL = "https://" + fn + ".js-company.ru"
		v.Shop.DocsXML = "https://" + cfg.Client.Documentation

		v.Shop.Currencies.Currency.ID = "RUR"
		v.Shop.Currencies.Currency.Rate = "1"

		for i := 0; i < len(catDB); i++ {
			v.Shop.Categories.AddCategory(catDB[i].id, catDB[i].parentID, catDB[i].name)
		}

		for i := 0; i < numberOffers; i++ {
			var props [numberParam]Param
			props[0].Name = "Пол"
			props[0].Text = propertyDB[pDB[i].sex]
			props[1].Name = "Возраст"
			props[1].Text = propertyDB[pDB[i].age]
			props[2].Name = "Состав"
			props[2].Text = pDB[i].structure
			props[3].Name = "Размер"
			props[3].Text = sizeDB[pDB[i].id]
			props[4].Name = "Цвет"
			props[4].Text = colorDB[pDB[i].id]
			props[5].Name = "RGB"
			props[5].Text = rgbDB[pDB[i].id]

			imgURL := []string{}
			for k := 0; k < len(imagesDB); k++ {
				if imagesDB[k].productID == pDB[i].groupID && (imagesDB[k].color == colorDB[pDB[i].id]) {
					imgURL = append(imgURL, "https://www.js-company.ru/"+imagesDB[k].URL)
				}
			}

			boxImgURL := []string{}
			for k := 0; k < len(boxImagesDB); k++ {
				if boxImagesDB[k].productID == pDB[i].groupID {
					boxImgURL = append(boxImgURL, "https://www.js-company.ru/"+boxImagesDB[k].URL)
				}
			}

			v.Shop.AllOffers.AddOffer(
				pDB[i].id,
				pDB[i].groupID,
				pDB[i].uuid,
				int(pDB[i].amount),
				pDB[i].nameOffer,
				pDB[i].url,
				regularPriceDB[pDB[i].id],
				salesPriceDB[pDB[i].id],
				pDB[i].parentID,
				pDB[i].barcode,
				pDB[i].code,
				pDB[i].brand,
				pDB[i].kind,
				&imgURL,
				&boxImgURL,
				pDB[i].description,
				props,
			)
		}

		xmlFileName := pathResult + fn + ".xml"
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
}
