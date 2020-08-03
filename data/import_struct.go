type КоммерческаяИнформация struct {
	XMLName          xml.Name `xml:"КоммерческаяИнформация"`
	Text             string   `xml:",chardata"`
	ВерсияСхемы      string   `xml:"ВерсияСхемы,attr"`
	ДатаФормирования string   `xml:"ДатаФормирования,attr"`
	Классификатор    struct {
		Text         string `xml:",chardata"`
		Ид           string `xml:"Ид"`
		Наименование string `xml:"Наименование"`
		Владелец     struct {
			Text                    string `xml:",chardata"`
			Ид                      string `xml:"Ид"`
			Наименование            string `xml:"Наименование"`
			ОфициальноеНаименование string `xml:"ОфициальноеНаименование"`
			ЮридическийАдрес        struct {
				Text          string `xml:",chardata"`
				Представление string `xml:"Представление"`
				АдресноеПоле  []struct {
					Text     string `xml:",chardata"`
					Тип      string `xml:"Тип"`
					Значение string `xml:"Значение"`
				} `xml:"АдресноеПоле"`
			} `xml:"ЮридическийАдрес"`
			ИНН            string `xml:"ИНН"`
			КПП            string `xml:"КПП"`
			ОКПО           string `xml:"ОКПО"`
			РасчетныеСчета struct {
				Text          string `xml:",chardata"`
				РасчетныйСчет struct {
					Text       string `xml:",chardata"`
					НомерСчета string `xml:"НомерСчета"`
					Банк       struct {
						Text                  string `xml:",chardata"`
						СчетКорреспондентский string `xml:"СчетКорреспондентский"`
						Наименование          string `xml:"Наименование"`
						Адрес                 struct {
							Text          string `xml:",chardata"`
							Представление string `xml:"Представление"`
						} `xml:"Адрес"`
						БИК string `xml:"БИК"`
					} `xml:"Банк"`
				} `xml:"РасчетныйСчет"`
			} `xml:"РасчетныеСчета"`
		} `xml:"Владелец"`
		Производители struct {
			Text          string `xml:",chardata"`
			Производитель []struct {
				Text            string `xml:",chardata"`
				Ид              string `xml:"Ид"`
				ПометкаУдаления string `xml:"ПометкаУдаления"`
				НашБренд        string `xml:"НашБренд"`
				Наименование    string `xml:"Наименование"`
			} `xml:"Производитель"`
		} `xml:"Производители"`
		КатегорииYML struct {
			Text      string `xml:",chardata"`
			Категория []struct {
				Text            string `xml:",chardata"`
				Ид              string `xml:"Ид"`
				ПометкаУдаления string `xml:"ПометкаУдаления"`
				Наименование    string `xml:"Наименование"`
				ИдГруппы        string `xml:"ИдГруппы"`
			} `xml:"Категория"`
		} `xml:"КатегорииYML"`
		Группы struct {
			Text   string `xml:",chardata"`
			Группа struct {
				Text            string `xml:",chardata"`
				Ид              string `xml:"Ид"`
				Наименование    string `xml:"Наименование"`
				Слаг            string `xml:"Слаг"`
				ПометкаУдаления string `xml:"ПометкаУдаления"`
				Группы          struct {
					Text   string `xml:",chardata"`
					Группа struct {
						Text            string `xml:",chardata"`
						Ид              string `xml:"Ид"`
						Наименование    string `xml:"Наименование"`
						Слаг            string `xml:"Слаг"`
						ПометкаУдаления string `xml:"ПометкаУдаления"`
						Группы          struct {
							Text   string `xml:",chardata"`
							Группа struct {
								Text            string `xml:",chardata"`
								Ид              string `xml:"Ид"`
								Наименование    string `xml:"Наименование"`
								Слаг            string `xml:"Слаг"`
								ПометкаУдаления string `xml:"ПометкаУдаления"`
								Группы          struct {
									Text   string `xml:",chardata"`
									Группа struct {
										Text            string `xml:",chardata"`
										Ид              string `xml:"Ид"`
										Наименование    string `xml:"Наименование"`
										Слаг            string `xml:"Слаг"`
										ПометкаУдаления string `xml:"ПометкаУдаления"`
									} `xml:"Группа"`
								} `xml:"Группы"`
							} `xml:"Группа"`
						} `xml:"Группы"`
					} `xml:"Группа"`
				} `xml:"Группы"`
			} `xml:"Группа"`
		} `xml:"Группы"`
		Свойства struct {
			Text     string `xml:",chardata"`
			Свойство []struct {
				Text             string `xml:",chardata"`
				Ид               string `xml:"Ид"`
				ПометкаУдаления  string `xml:"ПометкаУдаления"`
				Наименование     string `xml:"Наименование"`
				ТипЗначений      string `xml:"ТипЗначений"`
				ВариантыЗначений struct {
					Text       string `xml:",chardata"`
					Справочник []struct {
						Text            string `xml:",chardata"`
						ИдЗначения      string `xml:"ИдЗначения"`
						ПометкаУдаления string `xml:"ПометкаУдаления"`
						Значение        string `xml:"Значение"`
					} `xml:"Справочник"`
				} `xml:"ВариантыЗначений"`
				ДляТоваров string `xml:"ДляТоваров"`
				Основное   string `xml:"Основное"`
			} `xml:"Свойство"`
		} `xml:"Свойства"`
	} `xml:"Классификатор"`
	Каталог struct {
		Text                    string `xml:",chardata"`
		СодержитТолькоИзменения string `xml:"СодержитТолькоИзменения,attr"`
		Ид                      string `xml:"Ид"`
		ИдКлассификатора        string `xml:"ИдКлассификатора"`
		Наименование            string `xml:"Наименование"`
		Владелец                struct {
			Text                    string `xml:",chardata"`
			Ид                      string `xml:"Ид"`
			Наименование            string `xml:"Наименование"`
			ОфициальноеНаименование string `xml:"ОфициальноеНаименование"`
			ЮридическийАдрес        struct {
				Text          string `xml:",chardata"`
				Представление string `xml:"Представление"`
				АдресноеПоле  []struct {
					Text     string `xml:",chardata"`
					Тип      string `xml:"Тип"`
					Значение string `xml:"Значение"`
				} `xml:"АдресноеПоле"`
			} `xml:"ЮридическийАдрес"`
			ИНН            string `xml:"ИНН"`
			КПП            string `xml:"КПП"`
			ОКПО           string `xml:"ОКПО"`
			РасчетныеСчета struct {
				Text          string `xml:",chardata"`
				РасчетныйСчет struct {
					Text       string `xml:",chardata"`
					НомерСчета string `xml:"НомерСчета"`
					Банк       struct {
						Text                  string `xml:",chardata"`
						СчетКорреспондентский string `xml:"СчетКорреспондентский"`
						Наименование          string `xml:"Наименование"`
						Адрес                 struct {
							Text          string `xml:",chardata"`
							Представление string `xml:"Представление"`
						} `xml:"Адрес"`
						БИК string `xml:"БИК"`
					} `xml:"Банк"`
				} `xml:"РасчетныйСчет"`
			} `xml:"РасчетныеСчета"`
		} `xml:"Владелец"`
		Товары struct {
			Text  string `xml:",chardata"`
			Товар struct {
				Text                                string `xml:",chardata"`
				Ид                                  string `xml:"Ид"`
				Артикул                             string `xml:"Артикул"`
				Код                                 string `xml:"Код"`
				Наименование                        string `xml:"Наименование"`
				ПродажаУпаковками                   string `xml:"ПродажаУпаковками"`
				КоличествоШтукВУпаковке             string `xml:"КоличествоШтукВУпаковке"`
				КоличествоШтукВКоробке              string `xml:"КоличествоШтукВКоробке"`
				ПометкаУдаления                     string `xml:"ПометкаУдаления"`
				Слаг                                string `xml:"Слаг"`
				Тип                                 string `xml:"Тип"`
				Вид                                 string `xml:"Вид"`
				Бренд                               string `xml:"Бренд"`
				ИдПроизводителя                     string `xml:"ИдПроизводителя"`
				КатегорияYML                        string `xml:"КатегорияYML"`
				Состав                              string `xml:"Состав"`
				Новинка                             string `xml:"Новинка"`
				СкороВПродаже                       string `xml:"СкороВПродаже"`
				РекомендоватьПокупателю             string `xml:"РекомендоватьПокупателю"`
				Распродажа                          string `xml:"Распродажа"`
				КоэффициентСкидки                   string `xml:"КоэффициентСкидки"`
				РазмернаяСетка                      string `xml:"РазмернаяСетка"`
				ОтклонениеРазмера                   string `xml:"ОтклонениеРазмера"`
				ПоказыватьЦветаНаСайтеТолькоНаличия string `xml:"ПоказыватьЦветаНаСайтеТолькоНаличия"`
				ВидеоВыводитьПервым                 string `xml:"ВидеоВыводитьПервым"`
				ЗащитаВатермарк                     string `xml:"ЗащитаВатермарк"`
				БазоваяЕдиница                      struct {
					Text                    string `xml:",chardata"`
					Код                     string `xml:"Код,attr"`
					НаименованиеПолное      string `xml:"НаименованиеПолное,attr"`
					МеждународноеСокращение string `xml:"МеждународноеСокращение,attr"`
				} `xml:"БазоваяЕдиница"`
				Группы struct {
					Text string `xml:",chardata"`
					Ид   string `xml:"Ид"`
				} `xml:"Группы"`
				Описание string `xml:"Описание"`
				Картинки struct {
					Text                string `xml:",chardata"`
					ОсновноеИзображение struct {
						Text  string `xml:",chardata"`
						Номер string `xml:"Номер"`
						Файл  string `xml:"Файл"`
					} `xml:"ОсновноеИзображение"`
					Цвета struct {
						Text string `xml:",chardata"`
						Цвет struct {
							Text         string `xml:",chardata"`
							Наименование string `xml:"Наименование"`
							Файл         string `xml:"Файл"`
						} `xml:"Цвет"`
					} `xml:"Цвета"`
				} `xml:"Картинки"`
				ЗначенияСвойств struct {
					Text             string `xml:",chardata"`
					ЗначенияСвойства []struct {
						Text     string `xml:",chardata"`
						Ид       string `xml:"Ид"`
						Значение string `xml:"Значение"`
					} `xml:"ЗначенияСвойства"`
				} `xml:"ЗначенияСвойств"`
				СтавкиНалогов struct {
					Text         string `xml:",chardata"`
					СтавкаНалога struct {
						Text         string `xml:",chardata"`
						Наименование string `xml:"Наименование"`
						Ставка       string `xml:"Ставка"`
					} `xml:"СтавкаНалога"`
				} `xml:"СтавкиНалогов"`
				ЗначенияРеквизитов struct {
					Text              string `xml:",chardata"`
					ЗначениеРеквизита []struct {
						Text         string `xml:",chardata"`
						Наименование string `xml:"Наименование"`
						Значение     string `xml:"Значение"`
					} `xml:"ЗначениеРеквизита"`
				} `xml:"ЗначенияРеквизитов"`
			} `xml:"Товар"`
		} `xml:"Товары"`
	} `xml:"Каталог"`
} 
