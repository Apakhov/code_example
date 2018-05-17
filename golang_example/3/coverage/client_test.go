package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type TestCase struct {
	Request  SearchRequest
	TestUser SearchClient
	Response *SearchResponse
	IsError  bool
	Error    string
}

func TestRightWork(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "e",
				OrderField: "Id",
				OrderBy:    0,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "BoydWolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "HildaMayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female"}},
				NextPage: true,
			},
		},
		TestCase{
			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "e",
				OrderField: "Id",
				OrderBy:    1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "BoydWolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "HildaMayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female"}},
				NextPage: true,
			},
		},
		TestCase{

			Request: SearchRequest{
				Limit:      2,
				Offset:     0,
				Query:      "e",
				OrderField: "Id",
				OrderBy:    -1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{
					User{
						Id:     1,
						Name:   "HildaMayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female"},
					User{
						Id:     0,
						Name:   "BoydWolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male",
					}},
				NextPage: true,
			},
		},
		TestCase{
			Request: SearchRequest{
				Limit:      29,
				Offset:     3,
				Query:      "",
				OrderField: "Id",
				OrderBy:    0,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{User{Id: 3, Name: "EverettDillard", Age: 27, About: "Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n", Gender: "male"}, User{Id: 4, Name: "OwenLynn", Age: 30, About: "Elit anim elit eu et deserunt veniam laborum commodo irure nisi ut labore reprehenderit fugiat. Ipsum adipisicing labore ullamco occaecat ut. Ea deserunt ad dolor eiusmod aute non enim adipisicing sit ullamco est ullamco. Elit in proident pariatur elit ullamco quis. Exercitation amet nisi fugiat voluptate esse sit et consequat sit pariatur labore et.\n", Gender: "male"}, User{Id: 5, Name: "BeulahStark", Age: 30, About: "Enim cillum eu cillum velit labore. In sint esse nulla occaecat voluptate pariatur aliqua aliqua non officia nulla aliqua. Fugiat nostrud irure officia minim cupidatat laborum ad incididunt dolore. Fugiat nostrud eiusmod ex ea nulla commodo. Reprehenderit sint qui anim non ad id adipisicing qui officia Lorem.\n", Gender: "female"}, User{Id: 6, Name: "JenningsMays", Age: 39, About: "Veniam consectetur non non aliquip exercitation quis qui. Aliquip duis ut ad commodo consequat ipsum cupidatat id anim voluptate deserunt enim laboris. Sunt nostrud voluptate do est tempor esse anim pariatur. Ea do amet Lorem in mollit ipsum irure Lorem exercitation. Exercitation deserunt adipisicing nulla aute ex amet sint tempor incididunt magna. Quis et consectetur dolor nulla reprehenderit culpa laboris voluptate ut mollit. Qui ipsum nisi ullamco sit exercitation nisi magna fugiat anim consectetur officia.\n", Gender: "male"}, User{Id: 7, Name: "LeannTravis", Age: 34, About: "Lorem magna dolore et velit ut officia. Cupidatat deserunt elit mollit amet nulla voluptate sit. Quis aute aliquip officia deserunt sint sint nisi. Laboris sit et ea dolore consequat laboris non. Consequat do enim excepteur qui mollit consectetur eiusmod laborum ut duis mollit dolor est. Excepteur amet duis enim laborum aliqua nulla ea minim.\n", Gender: "female"}, User{Id: 8, Name: "GlennJordan", Age: 29, About: "Duis reprehenderit sit velit exercitation non aliqua magna quis ad excepteur anim. Eu cillum cupidatat sit magna cillum irure occaecat sunt officia officia deserunt irure. Cupidatat dolor cupidatat ipsum minim consequat Lorem adipisicing. Labore fugiat cupidatat nostrud voluptate ea eu pariatur non. Ipsum quis occaecat irure amet esse eu fugiat deserunt incididunt Lorem esse duis occaecat mollit.\n", Gender: "male"}, User{Id: 9, Name: "RoseCarney", Age: 36, About: "Voluptate ipsum ad consequat elit ipsum tempor irure consectetur amet. Et veniam sunt in sunt ipsum non elit ullamco est est eu. Exercitation ipsum do deserunt do eu adipisicing id deserunt duis nulla ullamco eu. Ad duis voluptate amet quis commodo nostrud occaecat minim occaecat commodo. Irure sint incididunt est cupidatat laborum in duis enim nulla duis ut in ut. Cupidatat ex incididunt do ullamco do laboris eiusmod quis nostrud excepteur quis ea.\n", Gender: "female"}, User{Id: 10, Name: "HendersonMaxwell", Age: 30, About: "Ex et excepteur anim in eiusmod. Cupidatat sunt aliquip exercitation velit minim aliqua ad ipsum cillum dolor do sit dolore cillum. Exercitation eu in ex qui voluptate fugiat amet.\n", Gender: "male"}, User{Id: 11, Name: "GilmoreGuerra", Age: 32, About: "Labore consectetur do sit et mollit non incididunt. Amet aute voluptate enim et sit Lorem elit. Fugiat proident ullamco ullamco sint pariatur deserunt eu nulla consectetur culpa eiusmod. Veniam irure et deserunt consectetur incididunt ad ipsum sint. Consectetur voluptate adipisicing aute fugiat aliquip culpa qui nisi ut ex esse ex. Sint et anim aliqua pariatur.\n", Gender: "male"}, User{Id: 12, Name: "CruzGuerrero", Age: 36, About: "Sunt enim ad fugiat minim id esse proident laborum magna magna. Velit anim aliqua nulla laborum consequat veniam reprehenderit enim fugiat ipsum mollit nisi. Nisi do reprehenderit aute sint sit culpa id Lorem proident id tempor. Irure ut ipsum sit non quis aliqua in voluptate magna. Ipsum non aliquip quis incididunt incididunt aute sint. Minim dolor in mollit aute duis consectetur.\n", Gender: "male"}, User{Id: 13, Name: "WhitleyDavidson", Age: 40, About: "Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n", Gender: "male"}, User{Id: 14, Name: "NicholsonNewman", Age: 23, About: "Tempor minim reprehenderit dolore et ad. Irure id fugiat incididunt do amet veniam ex consequat. Quis ad ipsum excepteur eiusmod mollit nulla amet velit quis duis ut irure.\n", Gender: "male"}, User{Id: 15, Name: "AllisonValdez", Age: 21, About: "Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n", Gender: "male"}, User{Id: 16, Name: "AnnieOsborn", Age: 35, About: "Consequat fugiat veniam commodo nisi nostrud culpa pariatur. Aliquip velit adipisicing dolor et nostrud. Eu nostrud officia velit eiusmod ullamco duis eiusmod ad non do quis.\n", Gender: "female"}, User{Id: 17, Name: "DillardMccoy", Age: 36, About: "Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n", Gender: "male"}, User{Id: 18, Name: "TerrellHall", Age: 27, About: "Ut nostrud est est elit incididunt consequat sunt ut aliqua sunt sunt. Quis consectetur amet occaecat nostrud duis. Fugiat in irure consequat laborum ipsum tempor non deserunt laboris id ullamco cupidatat sit. Officia cupidatat aliqua veniam et ipsum labore eu do aliquip elit cillum. Labore culpa exercitation sint sint.\n", Gender: "male"}, User{Id: 19, Name: "BellBauer", Age: 26, About: "Nulla voluptate nostrud nostrud do ut tempor et quis non aliqua cillum in duis. Sit ipsum sit ut non proident exercitation. Quis consequat laboris deserunt adipisicing eiusmod non cillum magna.\n", Gender: "male"}, User{Id: 20, Name: "LoweryYork", Age: 27, About: "Dolor enim sit id dolore enim sint nostrud deserunt. Occaecat minim enim veniam proident mollit Lorem irure ex. Adipisicing pariatur adipisicing aliqua amet proident velit. Magna commodo culpa sit id.\n", Gender: "male"}, User{Id: 21, Name: "JohnsWhitney", Age: 26, About: "Elit sunt exercitation incididunt est ea quis do ad magna. Commodo laboris nisi aliqua eu incididunt eu irure. Labore ullamco quis deserunt non cupidatat sint aute in incididunt deserunt elit velit. Duis est mollit veniam aliquip. Nulla sunt veniam anim et sint dolore.\n", Gender: "male"}, User{Id: 22, Name: "BethWynn", Age: 31, About: "Proident non nisi dolore id non. Aliquip ex anim cupidatat dolore amet veniam tempor non adipisicing. Aliqua adipisicing eu esse quis reprehenderit est irure cillum duis dolor ex. Laborum do aute commodo amet. Fugiat aute in excepteur ut aliqua sint fugiat do nostrud voluptate duis do deserunt. Elit esse ipsum duis ipsum.\n", Gender: "female"}, User{Id: 23, Name: "GatesSpencer", Age: 21, About: "Dolore magna magna commodo irure. Proident culpa nisi veniam excepteur sunt qui et laborum tempor. Qui proident Lorem commodo dolore ipsum.\n", Gender: "male"}, User{Id: 24, Name: "GonzalezAnderson", Age: 33, About: "Quis consequat incididunt in ex deserunt minim aliqua ea duis. Culpa nisi excepteur sint est fugiat cupidatat nulla magna do id dolore laboris. Aute cillum eiusmod do amet dolore labore commodo do pariatur sit id. Do irure eiusmod reprehenderit non in duis sunt ex. Labore commodo labore pariatur ex minim qui sit elit.\n", Gender: "male"}, User{Id: 25, Name: "KatherynJacobs", Age: 32, About: "Magna excepteur anim amet id consequat tempor dolor sunt id enim ipsum ea est ex. In do ea sint qui in minim mollit anim est et minim dolore velit laborum. Officia commodo duis ut proident laboris fugiat commodo do ex duis consequat exercitation. Ad et excepteur ex ea exercitation id fugiat exercitation amet proident adipisicing laboris id deserunt. Commodo proident laborum elit ex aliqua labore culpa ullamco occaecat voluptate voluptate laboris deserunt magna.\n", Gender: "female"}, User{Id: 26, Name: "SimsCotton", Age: 39, About: "Ex cupidatat est velit consequat ad. Tempor non cillum labore non voluptate. Et proident culpa labore deserunt ut aliquip commodo laborum nostrud. Anim minim occaecat est est minim.\n", Gender: "male"}, User{Id: 27, Name: "RebekahSutton", Age: 26, About: "Aliqua exercitation ad nostrud et exercitation amet quis cupidatat esse nostrud proident. Ullamco voluptate ex minim consectetur ea cupidatat in mollit reprehenderit voluptate labore sint laboris. Minim cillum et incididunt pariatur amet do esse. Amet irure elit deserunt quis culpa ut deserunt minim proident cupidatat nisi consequat ipsum.\n", Gender: "female"}}, NextPage: true}},

		TestCase{
			Request: SearchRequest{
				Limit:      0,
				Offset:     0,
				Query:      "e",
				OrderField: "Id",
				OrderBy:    -1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users:    []User{},
				NextPage: true,
			},
		},
		TestCase{
			Request: SearchRequest{
				Limit:      6,
				Offset:     0,
				Query:      "eshdgfkjhsdgfhgsdkhgfhskdjhgkijshdgjhsdkjhgksdjhgkjhsdgkjhsdkjhgskjh",
				OrderField: "Id",
				OrderBy:    -1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
		},
		TestCase{
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "e",
				OrderField: "Name",
				OrderBy:    -1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{
					User{
						Id:     1,
						Name:   "HildaMayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
					User{Id: 2,
						Name:   "BrooksAguilar",
						Age:    25,
						About:  "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n",
						Gender: "male"},
					User{Id: 0,
						Name:   "BoydWolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male"},
				},
				NextPage: true},
		},
		TestCase{
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "e",
				OrderField: "",
				OrderBy:    -1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{
					User{
						Id:     1,
						Name:   "HildaMayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female",
					},
					User{Id: 2,
						Name:   "BrooksAguilar",
						Age:    25,
						About:  "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n",
						Gender: "male"},
					User{Id: 0,
						Name:   "BoydWolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male"},
				},
				NextPage: true},
		},
		TestCase{
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "e",
				OrderField: "Age",
				OrderBy:    -1,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},

			Response: &SearchResponse{
				Users: []User{

					User{Id: 2,
						Name:   "BrooksAguilar",
						Age:    25,
						About:  "Velit ullamco est aliqua voluptate nisi do. Voluptate magna anim qui cillum aliqua sint veniam reprehenderit consectetur enim. Laborum dolore ut eiusmod ipsum ad anim est do tempor culpa ad do tempor. Nulla id aliqua dolore dolore adipisicing.\n",
						Gender: "male"},
					User{Id: 0,
						Name:   "BoydWolf",
						Age:    22,
						About:  "Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n",
						Gender: "male"},
					User{
						Id:     1,
						Name:   "HildaMayer",
						Age:    21,
						About:  "Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n",
						Gender: "female"},
				},
				NextPage: true},
		},
	}

	s := Server{
		rightAccessToken: "wrongAccessToken",
		datasetPath:      "dataset.xml",
	}

	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	for caseNum, item := range cases {
		item.TestUser.URL = ts.URL
		result, err := item.TestUser.FindUsers(item.Request)

		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if !reflect.DeepEqual(item.Response, result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v", caseNum, item.Response, result)
		}
	}
	ts.Close()
}

func TestErrorsClientUse(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Request: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "e",
				OrderField: "Id",
				OrderBy:    0,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
			IsError: true,
			Error:   "limit must be > 0",
		},
		TestCase{
			Request: SearchRequest{
				Limit:      2,
				Offset:     -1,
				Query:      "e",
				OrderField: "Id",
				OrderBy:    0,
			},

			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
			IsError: true,
			Error:   "offset must be > 0",
		},
	}

	s := Server{
		rightAccessToken: "wrongAccessToken",
		datasetPath:      "dataset.xml",
	}

	ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))

	for caseNum, item := range cases {
		item.TestUser.URL = ts.URL
		_, err := item.TestUser.FindUsers(item.Request)

		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		if item.Error != err.Error() {
			t.Errorf("[%d] expected error: %s,\ngot: %s", caseNum, item.Error, err.Error())
		}
	}
	ts.Close()
}

type ErrorTest struct {
	Request  SearchRequest
	TestUser SearchClient
	Handler  func(w http.ResponseWriter, r *http.Request)
	Error    string
}

func TestErrorsClientResp(t *testing.T) {
	cases := []ErrorTest{
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				_, err := ioutil.ReadFile("sdfsdfsdfsdfds")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			},
			Error: "SearchServer fatal error",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				res, _ := json.Marshal(SearchErrorResponse{"bad Accesstoken"})
				http.Error(w, string(res), http.StatusUnauthorized)
			},
			Error: "Bad AccessToken",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				res, _ := json.Marshal(SearchErrorResponse{"bad method"})
				http.Error(w, string(res), http.StatusMethodNotAllowed)
			},
			Error: "cant unpack result json: json: cannot unmarshal object into Go value of type []main.User",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				_, err := strconv.Atoi("anton")
				res, _ := json.Marshal(err)
				http.Error(w, string(res), http.StatusBadRequest)
			},
			Error: "unknown bad request error: ",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				_, err := strconv.Atoi("anton")
				http.Error(w, err.Error(), http.StatusBadRequest)
			},
			Error: "cant unpack error json: invalid character 's' looking for beginning of value",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, `{ "Error": "ErrorBadOrderField"}`, http.StatusBadRequest)
			},
			Error: "OrderFeld  invalid",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				http.Error(w, `{ "Error": "ErrorBadOrderField"}`, http.StatusBadRequest)
			},
			Error: "timeout for limit=1&offset=0&order_by=0&order_field=&query=",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8080/",
			},
		},
		ErrorTest{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				log.Println(":-(")
			},
			Error: "unknown error Get http://127.0.0.1:8081/?limit=1&offset=0&order_by=0&order_field=&query=: dial tcp 127.0.0.1:8081: connect: connection refused",
			TestUser: SearchClient{
				AccessToken: "wrongAccessToken",
				URL:         "http://127.0.0.1:8081/",
			},
		},
	}
	for caseNum, item := range cases {
		ts := httptest.NewServer(http.HandlerFunc(item.Handler))
		if item.TestUser.URL != "http://127.0.0.1:8081/" {
			item.TestUser.URL = ts.URL
		}
		_, err := item.TestUser.FindUsers(item.Request)

		if item.Error != err.Error() {
			t.Errorf("[%d] expected error: %s,\ngot: %s", caseNum, item.Error, err.Error())
		}
		ts.Close()
	}

}

type ServerErrorTest struct {
	Keys        ParamsKeys
	Values      ParamsValues
	Response    SearchErrorResponse
	AccessToken string
	Method      string
	DatasetPath string
}

type ParamsKeys []string
type ParamsValues []string

func TestErrorsServerResp(t *testing.T) {
	cases := []ServerErrorTest{
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{"ErrorBadOrderField"},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"lmit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{`strconv.Atoi: parsing "": invalid syntax`},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"ffset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{`strconv.Atoi: parsing "": invalid syntax`},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"uery",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{"ErrorBadOrderField"},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"lmit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{`bad method`},
			AccessToken: "wrongAccessToken",
			Method:      "POST",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"lmit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{`strconv.Atoi: parsing "": invalid syntax`},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"6",
			},
			Response:    SearchErrorResponse{`bad Accesstoken`},
			AccessToken: "rightAccesstoken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"0",
			},
			Response:    SearchErrorResponse{"open broken_path.xml: no such file or directory"},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
			DatasetPath: "broken_path.xml",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"e",
				"Gender",
				"1",
			},
			Response:    SearchErrorResponse{"ErrorBadOrderField"},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"e",
				"Name",
				"retjerktjkertjjetkr",
			},
			Response:    SearchErrorResponse{`strconv.Atoi: parsing "retjerktjkertjjetkr": invalid syntax`},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
		},
		ServerErrorTest{
			Keys: ParamsKeys{
				"limit",
				"offset",
				"query",
				"order_field",
				"order_by",
			},
			Values: ParamsValues{
				"2",
				"0",
				"id",
				"Age",
				"0",
			},
			Response:    SearchErrorResponse{"EOF"},
			AccessToken: "wrongAccessToken",
			Method:      "GET",
			DatasetPath: "broken.xml",
		},
	}
	s := Server{
		rightAccessToken: "wrongAccessToken",
		datasetPath:      "dataset.xml",
	}
	log.Println("starting last")

	for caseNum, item := range cases {
		if item.DatasetPath != "" {
			s.datasetPath = item.DatasetPath
		}
		ts := httptest.NewServer(http.HandlerFunc(s.SearchServer))
		searcherParams := url.Values{}
		for i, x := range item.Keys {
			searcherParams.Add(x, item.Values[i])
		}

		searcherReq, _ := http.NewRequest(item.Method, ts.URL+"?"+searcherParams.Encode(), nil)
		searcherReq.Header.Add("AccessToken", item.AccessToken)
		resp, err := client.Do(searcherReq)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		errResp := SearchErrorResponse{}
		err = json.Unmarshal(body, &errResp)

		if item.Response.Error != errResp.Error {
			t.Errorf("[%d] expected error: %s,\ngot: %s", caseNum, item.Response.Error, errResp.Error)
		}
		s.datasetPath = "dataset.xml"
		ts.Close()
	}

}
