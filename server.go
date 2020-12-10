package main

import (
	"ProjectLive/database/submissions"
	"ProjectLive/database/users"
	hashtable "ProjectLive/hashTable"
	"ProjectLive/logger"
	"ProjectLive/secure"
	"ProjectLive/url"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type condition struct {
	Storage             []string
	Screen              []string
	Housing             []string
	AnyOtherIssues      []string
	OriginalAccessories []string
}

//data contain the fields for the data parsing into the template
type data struct {
	NameOfPhone  string
	ID           string
	CustomerName string
	SellerName   string
}

var (
	tpl          *template.Template
	sessionMap   = hashtable.Init() //uuid as the key, value as the username
	userTrackMap = hashtable.Init() //key is the username, value is the
	sqluser      = "root"
	sqlpassword  = "password"
	userMap      map[string]users.User
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	//connect to database and fill the datastructures with info from database
	db := connectDB()
	defer db.Close()
	var err error
	userMap, err = users.GetRecord(db)
	if err != nil {
		log.Println(err)
	}

}
func connectDB() *sql.DB {
	connectionString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:8888)/store", sqluser, sqlpassword)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Println("Good to go")
	}
	return db
}
func main() {
	urlPattern := url.ReturnURL() //gets url string from url package
	//handles static css files
	http.Handle(urlPattern.Static, http.StripPrefix(urlPattern.Static, http.FileServer(http.Dir("."+urlPattern.Static))))
	go http.HandleFunc(urlPattern.Home, index)
	go http.HandleFunc(urlPattern.Signup, signup)
	go http.HandleFunc(urlPattern.Login, login)
	go http.HandleFunc(urlPattern.CustomerSell, customerSell)
	go http.HandleFunc(urlPattern.OrderList, orderList)
	go http.HandleFunc(urlPattern.Logout, logout)
	// go http.HandleFunc(urlPattern.AutoLogout, autoLogout)

	log.Fatalln(http.ListenAndServe(":5000", nil))
}

func getUser(w http.ResponseWriter, r *http.Request) users.User {
	//get current session cookie
	myCookie, err := r.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		//create session cookie
		myCookie = &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
	}
	http.SetCookie(w, myCookie)
	var myUser users.User
	//if user already exists, use the cookie value to extract username as key to user struct
	username, err := sessionMap.Search(myCookie.Value)
	if err != nil {
		return myUser
	}
	myUser = userMap[username.(string)]
	return myUser
}
func alreadyLoggedIn(r *http.Request) bool {
	//check if session cookie is still present
	myCookie, err := r.Cookie("myCookie")
	if err != nil {
		return false
	}
	_, err = secure.ParseToken(myCookie.Value)
	if err != nil {
		//sign out for user, delete from sessionMap and userTrackMap
		username, err := sessionMap.Search(myCookie.Value)
		if err != nil {
			return false
		}
		fmt.Println(userTrackMap.Delete(username.(string)))
		fmt.Println(sessionMap.Delete(myCookie.Value))
		return false
	}
	//check if the userMap contain the user information
	username, err := sessionMap.Search(myCookie.Value)
	if err != nil {
		return false
	}
	_, ok := userMap[username.(string)]
	if !ok {
		return false
	}
	return true
}
func index(w http.ResponseWriter, r *http.Request) {
	myUser := getUser(w, r)
	// fmt.Printf("Type of username: %T", myUser.Username)
	err := tpl.ExecuteTemplate(w, "index.html", myUser)
	if err != nil {
		log.Println(err)
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//form submission of user details
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		company := r.FormValue("company")
		isCompany := r.FormValue("isCompany")
		//input of username is not empty, check if it is taken
		if username != "" {
			_, ok := userMap[username]
			if ok {
				tpl.ExecuteTemplate(w, "redirect.html", "Username was taken please select another")
				return
			}

			if !secure.InputValidate(username) {
				tpl.ExecuteTemplate(w, "redirect.html", "Your username should not contain ', \",  <, >, tabs or empty spaces")
				return
			}
			if !secure.InputValidate(password) {
				tpl.ExecuteTemplate(w, "redirect.html", "Your password should not contain ', \",  <, >, tabs or empty spaces")
				return
			}
			if !secure.InputValidate(company) {
				tpl.ExecuteTemplate(w, "redirect.html", "Your user details should not contain ', \",  <, >, tabs or empty spaces")
				return
			}

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				logger.Logging("Bcrypting password from SignUp")
				return
			}
			if isCompany == "true" {
				userMap[username] = users.User{Username: username, Password: string(bPassword), Company: company, IsCompany: isCompany}
			} else {
				company = ""
				isCompany = "false"
				userMap[username] = users.User{Username: username, Password: string(bPassword), Company: company, IsCompany: isCompany}
			}
			//Put the user info in database
			db := connectDB()
			defer db.Close()
			users.InsertRecord(db, username, string(bPassword), company, isCompany)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//if method is not post, then user has not signup yet
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			logger.Logging("Recovered from login")
		}
	}()
	if alreadyLoggedIn(r) {
		myUser := getUser(w, r)
		if myUser.Username == "admin@gmail.com" {
			http.Redirect(w, r, "/adminOnly", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		myUser, ok := userMap[username]
		if !ok {
			tpl.ExecuteTemplate(w, "redirect.html", "No Username or password found")
			return
		}
		_, err := userTrackMap.Search(username)
		if err == nil {
			tpl.ExecuteTemplate(w, "redirect.html", "There is a similar account currently logged in, please logout first")
			return
		}
		//if username matches, compare password with userMap
		err = bcrypt.CompareHashAndPassword([]byte(myUser.Password), []byte(password))
		if err != nil {
			tpl.ExecuteTemplate(w, "redirect.html", "No Username or password found")
			return
		}
		//after user has logged in, create cookie
		id := uuid.NewV4()
		// encrypted, err := enDecrypt(id.String())
		claiming := &secure.MyClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
			},
			SessionID: id.String(),
		}
		signedToken, err := secure.GenerateJWT(claiming)
		if err != nil {
			logger.Logging("Error in generating token from login")
			panic("Error in generating token")
		}
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: signedToken,
		}
		http.SetCookie(w, myCookie)
		//create the session
		err = sessionMap.Insert(signedToken, username)
		if err != nil {
			log.Println(err)
		}
		//userTrackMap contains the unencrypted token of user
		err = userTrackMap.Insert(username, signedToken)
		if err != nil {
			log.Println(err)
		}
		if username == "admin" {
			http.Redirect(w, r, "/adminOnly", http.StatusSeeOther)
			return
		}
		//check if user is company or customer
		if myUser.Company == "" {
			//it is a customer
			http.Redirect(w, r, "/customer", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/sellers", http.StatusSeeOther)
		}
		return
	}
	//if method is not post, means user not logged in yet
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	myCookie, err := r.Cookie("myCookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	myUser := getUser(w, r)
	userTrackMap.Delete(myUser.Username)
	//delete the session from sessionMap
	sessionMap.Delete(myCookie.Value)
	//delete cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, myCookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//customerSell display the information for capturing phone and create order for sellers
func customerSell(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	if r.Method == http.MethodPost {
		//obtain information about the user phone selling details
		id := uuid.NewV4()
		deviceName := r.FormValue("devicename")
		storage := r.FormValue("storage")
		screen := r.FormValue("screen")
		housing := r.FormValue("housing")
		otherIssues := r.Form["otherissues"]
		issues := strings.Join(otherIssues, ",")
		accessories := r.Form["accessories"]
		acc := strings.Join(accessories, ",")
		//insert the data in the database under the submissions table
		db := connectDB()
		defer db.Close()
		myUser := getUser(w, r)
		err := submissions.InsertDetails(db, myUser.Username, deviceName, storage, housing, screen, acc, issues, id.String())
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//create a struct type with slices [device storage], [screen], [housing], [other issues], [original accessories]
	phoneCondition := condition{
		Storage:             []string{"32GB", "64GB", "128GB", "512GB"},
		Screen:              []string{"Cracked or chipped", "Moderate scratches", "Minor scratches", "Flawless"},
		Housing:             []string{"Cracked or chipped", "Moderate Scratches", "Minor scratches", "Flawless"},
		AnyOtherIssues:      []string{"Unable to power on", "LED display defective", "Camera faulty", "Touchscreen faulty", "Fingerprint/Face sensor faulty", "Flawless"},
		OriginalAccessories: []string{"Box", "Charging cable", "Power adaptor", "Earphones"},
	}
	tpl.ExecuteTemplate(w, "sellPhones.html", phoneCondition)
}

//orderList show the list of orders from the database that customer submitted
//generates a submission id, input into database
func orderList(w http.ResponseWriter, r *http.Request) {
	//get information from database and parse into templates for sellers to view
	//only show information that the seller did not reply to
	db := connectDB()
	defer db.Close()
	orders, err := submissions.GetDetails(db)
	if err != nil {
		log.Println(err)
	}
	dataSlice := []data{}
	for _, v := range orders {
		var values data
		values.ID = v.ID
		values.NameOfPhone = v.Name
		dataSlice = append(dataSlice, values)
	}
	// if r.Method == http.MethodPost {
	// 	//id of the order
	// 	myUser := getUser(w, r)
	// 	seller := myUser.Username
	// 	id := r.FormValue("id")
	// 	details := orders[id]
	// 	customer := details.Customer
	// 	err := quotation.InsertQuotation(db, customer, seller, id, quotation)

	// }
	tpl.ExecuteTemplate(w, "orderList.html", dataSlice)
}

//submittedOrder shows the list of submitted orders
func submittedOrder(w http.ResponseWriter, r *http.Request) {

}
