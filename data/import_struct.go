type ���������������������� struct {
	XMLName          xml.Name `xml:"����������������������"`
	Text             string   `xml:",chardata"`
	�����������      string   `xml:"�����������,attr"`
	���������������� string   `xml:"����������������,attr"`
	�������������    struct {
		Text         string `xml:",chardata"`
		��           string `xml:"��"`
		������������ string `xml:"������������"`
		��������     struct {
			Text                    string `xml:",chardata"`
			��                      string `xml:"��"`
			������������            string `xml:"������������"`
			����������������������� string `xml:"�����������������������"`
			����������������        struct {
				Text          string `xml:",chardata"`
				������������� string `xml:"�������������"`
				������������  []struct {
					Text     string `xml:",chardata"`
					���      string `xml:"���"`
					�������� string `xml:"��������"`
				} `xml:"������������"`
			} `xml:"����������������"`
			���            string `xml:"���"`
			���            string `xml:"���"`
			����           string `xml:"����"`
			�������������� struct {
				Text          string `xml:",chardata"`
				������������� struct {
					Text       string `xml:",chardata"`
					���������� string `xml:"����������"`
					����       struct {
						Text                  string `xml:",chardata"`
						��������������������� string `xml:"���������������������"`
						������������          string `xml:"������������"`
						�����                 struct {
							Text          string `xml:",chardata"`
							������������� string `xml:"�������������"`
						} `xml:"�����"`
						��� string `xml:"���"`
					} `xml:"����"`
				} `xml:"�������������"`
			} `xml:"��������������"`
		} `xml:"��������"`
		������������� struct {
			Text          string `xml:",chardata"`
			������������� []struct {
				Text            string `xml:",chardata"`
				��              string `xml:"��"`
				��������������� string `xml:"���������������"`
				��������        string `xml:"��������"`
				������������    string `xml:"������������"`
			} `xml:"�������������"`
		} `xml:"�������������"`
		���������YML struct {
			Text      string `xml:",chardata"`
			��������� []struct {
				Text            string `xml:",chardata"`
				��              string `xml:"��"`
				��������������� string `xml:"���������������"`
				������������    string `xml:"������������"`
				��������        string `xml:"��������"`
			} `xml:"���������"`
		} `xml:"���������YML"`
		������ struct {
			Text   string `xml:",chardata"`
			������ struct {
				Text            string `xml:",chardata"`
				��              string `xml:"��"`
				������������    string `xml:"������������"`
				����            string `xml:"����"`
				��������������� string `xml:"���������������"`
				������          struct {
					Text   string `xml:",chardata"`
					������ struct {
						Text            string `xml:",chardata"`
						��              string `xml:"��"`
						������������    string `xml:"������������"`
						����            string `xml:"����"`
						��������������� string `xml:"���������������"`
						������          struct {
							Text   string `xml:",chardata"`
							������ struct {
								Text            string `xml:",chardata"`
								��              string `xml:"��"`
								������������    string `xml:"������������"`
								����            string `xml:"����"`
								��������������� string `xml:"���������������"`
								������          struct {
									Text   string `xml:",chardata"`
									������ struct {
										Text            string `xml:",chardata"`
										��              string `xml:"��"`
										������������    string `xml:"������������"`
										����            string `xml:"����"`
										��������������� string `xml:"���������������"`
									} `xml:"������"`
								} `xml:"������"`
							} `xml:"������"`
						} `xml:"������"`
					} `xml:"������"`
				} `xml:"������"`
			} `xml:"������"`
		} `xml:"������"`
		�������� struct {
			Text     string `xml:",chardata"`
			�������� []struct {
				Text             string `xml:",chardata"`
				��               string `xml:"��"`
				���������������  string `xml:"���������������"`
				������������     string `xml:"������������"`
				�����������      string `xml:"�����������"`
				���������������� struct {
					Text       string `xml:",chardata"`
					���������� []struct {
						Text            string `xml:",chardata"`
						����������      string `xml:"����������"`
						��������������� string `xml:"���������������"`
						��������        string `xml:"��������"`
					} `xml:"����������"`
				} `xml:"����������������"`
				���������� string `xml:"����������"`
				��������   string `xml:"��������"`
			} `xml:"��������"`
		} `xml:"��������"`
	} `xml:"�������������"`
	������� struct {
		Text                    string `xml:",chardata"`
		����������������������� string `xml:"�����������������������,attr"`
		��                      string `xml:"��"`
		����������������        string `xml:"����������������"`
		������������            string `xml:"������������"`
		��������                struct {
			Text                    string `xml:",chardata"`
			��                      string `xml:"��"`
			������������            string `xml:"������������"`
			����������������������� string `xml:"�����������������������"`
			����������������        struct {
				Text          string `xml:",chardata"`
				������������� string `xml:"�������������"`
				������������  []struct {
					Text     string `xml:",chardata"`
					���      string `xml:"���"`
					�������� string `xml:"��������"`
				} `xml:"������������"`
			} `xml:"����������������"`
			���            string `xml:"���"`
			���            string `xml:"���"`
			����           string `xml:"����"`
			�������������� struct {
				Text          string `xml:",chardata"`
				������������� struct {
					Text       string `xml:",chardata"`
					���������� string `xml:"����������"`
					����       struct {
						Text                  string `xml:",chardata"`
						��������������������� string `xml:"���������������������"`
						������������          string `xml:"������������"`
						�����                 struct {
							Text          string `xml:",chardata"`
							������������� string `xml:"�������������"`
						} `xml:"�����"`
						��� string `xml:"���"`
					} `xml:"����"`
				} `xml:"�������������"`
			} `xml:"��������������"`
		} `xml:"��������"`
		������ struct {
			Text  string `xml:",chardata"`
			����� struct {
				Text                                string `xml:",chardata"`
				��                                  string `xml:"��"`
				�������                             string `xml:"�������"`
				���                                 string `xml:"���"`
				������������                        string `xml:"������������"`
				�����������������                   string `xml:"�����������������"`
				�����������������������             string `xml:"�����������������������"`
				����������������������              string `xml:"����������������������"`
				���������������                     string `xml:"���������������"`
				����                                string `xml:"����"`
				���                                 string `xml:"���"`
				���                                 string `xml:"���"`
				�����                               string `xml:"�����"`
				���������������                     string `xml:"���������������"`
				���������YML                        string `xml:"���������YML"`
				������                              string `xml:"������"`
				�������                             string `xml:"�������"`
				�������������                       string `xml:"�������������"`
				�����������������������             string `xml:"�����������������������"`
				����������                          string `xml:"����������"`
				�����������������                   string `xml:"�����������������"`
				��������������                      string `xml:"��������������"`
				�����������������                   string `xml:"�����������������"`
				����������������������������������� string `xml:"�����������������������������������"`
				�������������������                 string `xml:"�������������������"`
				���������������                     string `xml:"���������������"`
				��������������                      struct {
					Text                    string `xml:",chardata"`
					���                     string `xml:"���,attr"`
					������������������      string `xml:"������������������,attr"`
					����������������������� string `xml:"�����������������������,attr"`
				} `xml:"��������������"`
				������ struct {
					Text string `xml:",chardata"`
					��   string `xml:"��"`
				} `xml:"������"`
				�������� string `xml:"��������"`
				�������� struct {
					Text                string `xml:",chardata"`
					������������������� struct {
						Text  string `xml:",chardata"`
						����� string `xml:"�����"`
						����  string `xml:"����"`
					} `xml:"�������������������"`
					����� struct {
						Text string `xml:",chardata"`
						���� struct {
							Text         string `xml:",chardata"`
							������������ string `xml:"������������"`
							����         string `xml:"����"`
						} `xml:"����"`
					} `xml:"�����"`
				} `xml:"��������"`
				��������������� struct {
					Text             string `xml:",chardata"`
					���������������� []struct {
						Text     string `xml:",chardata"`
						��       string `xml:"��"`
						�������� string `xml:"��������"`
					} `xml:"����������������"`
				} `xml:"���������������"`
				������������� struct {
					Text         string `xml:",chardata"`
					������������ struct {
						Text         string `xml:",chardata"`
						������������ string `xml:"������������"`
						������       string `xml:"������"`
					} `xml:"������������"`
				} `xml:"�������������"`
				������������������ struct {
					Text              string `xml:",chardata"`
					����������������� []struct {
						Text         string `xml:",chardata"`
						������������ string `xml:"������������"`
						��������     string `xml:"��������"`
					} `xml:"�����������������"`
				} `xml:"������������������"`
			} `xml:"�����"`
		} `xml:"������"`
	} `xml:"�������"`
} 
