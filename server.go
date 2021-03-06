package main

import (
	"ProjectLive/binarytree"
	"ProjectLive/database/logger"
	"ProjectLive/database/quotation"
	"ProjectLive/database/secretkey"
	"ProjectLive/database/submissions"
	"ProjectLive/database/transactions"
	"ProjectLive/database/users"
	hashtable "ProjectLive/hashTable"
	"ProjectLive/secure"
	"ProjectLive/sorting"
	"ProjectLive/url"
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

var (
	tpl          *template.Template
	sessionMap   = hashtable.Init()  //uuid as the key, value as the username
	userTrackMap = hashtable.Init()  //key is the username, value is the jwt token
	sellerMap    = binarytree.Init() //stores the information on the seller information
	sqluser      string
	sqlpassword  string
	userMap      map[string]users.User //key is username, value is the User struct
	key          string                //encryption key
	apiKey       string
	clientID     string //oauth
	clientSecret string //oauth
)

type condition struct {
	Storage             []string
	Screen              []string
	Housing             []string
	AnyOtherIssues      []string
	OriginalAccessories []string
}

//data contain the fields for the data parsing into the template for selection
type data struct {
	NameOfPhone string
	ID          string
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	sqluser = os.Getenv("SQL_USER")
	sqlpassword = os.Getenv("SQL_PASSWORD")
	clientID = os.Getenv("GITHUB_ID")
	clientSecret = os.Getenv("GITHUB_SECRET")
	tpl = template.Must(template.ParseGlob("templates/*.html"))
	//connect to database and fill the datastructures with info from database
	db := connectDB()
	defer db.Close()
	userMap, err = users.GetRecord(db)
	if err != nil {
		logger.Logging(db, "Failed to retrieve record from database: init")
	}
	var sCompany []users.User
	sCompany, err = users.RetrieveSeller(db)
	if err != nil {
		log.Println(err)
	}
	for _, v := range sCompany {
		//populate the binary tree with company
		err = sellerMap.Insert(v.Company, v.Username)
		if err != nil {
			log.Println(err)
			return
		}
	}
	key, err = secretkey.GetKey(db, "encryptionKey")
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
	go http.HandleFunc(urlPattern.Search, searchSeller)
	go http.HandleFunc(urlPattern.List, listSellers)
	go http.HandleFunc(urlPattern.Signup, signup)
	go http.HandleFunc(urlPattern.Login, login)
	go http.HandleFunc(urlPattern.CustomerSell, customerSell)
	go http.HandleFunc(urlPattern.OrderList, orderList)
	go http.HandleFunc(urlPattern.Logout, logout)
	go http.HandleFunc(urlPattern.InsertQuotation, insertQuotation)
	go http.HandleFunc(urlPattern.ViewResponse, viewResponse)
	go http.HandleFunc(urlPattern.SellerTransaction, sellerViewTransaction)
	go http.HandleFunc(urlPattern.CustomerTransaction, customerViewTransaction)
	go http.HandleFunc(urlPattern.ForgetPassword, forgetPassword)
	go http.HandleFunc(urlPattern.ResetPassword, resetPassword)
	// log.Fatalln(http.ListenAndServe(":5000", nil))
	log.Fatalln(http.ListenAndServeTLS(":5000", "cert.pem", "key.pem", nil))
}

func searchSeller(w http.ResponseWriter, r *http.Request) {
	//take the binary tree and call the lookup function
	if r.Method == http.MethodPost {
		sellerName := r.FormValue("seller")
		var s binarytree.ReturnSellerInfo
		s, err := sellerMap.Lookup(sellerName)
		if err != nil {
			tpl.ExecuteTemplate(w, "redirect.html", "Company not found")
			return
		}
		tpl.ExecuteTemplate(w, "displaySeller.html", s)
		return
	}
	tpl.ExecuteTemplate(w, "searchSeller.html", nil)
}

func listSellers(w http.ResponseWriter, r *http.Request) {
	//take the binary tree and call the listallnodes function
	binarytree.ResetSlice()
	s := sellerMap.ListAllNodes(sellerMap.Root)
	sorted := sorting.Split(s)
	tpl.ExecuteTemplate(w, "list.html", sorted)
}

//getUsername returns the username or the email of the user
func getUsername(r *http.Request) string {
	myCookie, _ := r.Cookie("myCookie")
	username, err := sessionMap.Search(myCookie.Value)
	if err != nil {
		logger.Logging(connectDB(), "failed to get username from sessionMap: getUsername")
		return ""
	}
	return username
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
	myUser = userMap[username]
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
		fmt.Println(userTrackMap.Delete(username))
		fmt.Println(sessionMap.Delete(myCookie.Value))
		return false
	}
	//check if the userMap contain the user information
	username, err := sessionMap.Search(myCookie.Value)
	if err != nil {
		return false
	}
	_, ok := userMap[username]
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
		logger.Logging(connectDB(), "Failed to execute template: index")
		return
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
				logger.Logging(connectDB(), "Failure in Bcrypting password: SignUp")
				return
			}
			if isCompany == "true" {
				userMap[username] = users.User{Username: username, Password: string(bPassword), Company: company, IsCompany: isCompany}
				sellerMap.Insert(company, username)
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
			logger.Logging(connectDB(), "Recovered from login: login")
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
			//there is another account used, delete the other account

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
			logger.Logging(connectDB(), "Error in generating token from login: login")
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
			logger.Logging(connectDB(), "Failed to insert signedToken into sessionMap: login")
			return
		}
		//userTrackMap contains the JWT token of user
		err = userTrackMap.Insert(username, signedToken)
		if err != nil {
			logger.Logging(connectDB(), "Failed to insert into userTrackMap: login")
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
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		logger.Logging(connectDB(), "Failed to parse login.html template: login")
	}
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
			logger.Logging(db, "Failed to insert into submissions table: customerSell")
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
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	//get information from database and parse into templates for sellers to view
	//only show information that the seller did not reply to
	db := connectDB()
	defer db.Close()
	orders, err := submissions.GetDetails(db)
	if err != nil {
		logger.Logging(connectDB(), "Failed to retrieve info from submissions table: orderlist")
	}
	selleruser := getUsername(r)
	//returns the id that seller responds to already
	checkedID, err := quotation.SearchSeller(db, selleruser)
	dataSlice := []data{}
	for _, v := range orders {
		skip := false
		for _, v1 := range checkedID {
			if v.ID == v1 {
				skip = true
				break
			}
		}
		if skip == true {
			continue
		}
		var values data
		values.ID = v.ID
		values.NameOfPhone = v.Name
		dataSlice = append(dataSlice, values)
	}
	if r.Method == http.MethodPost {
		//get the sessionID from the cookie and search inside sessionmap, delete the current value
		//create new value containing the transaction id
		//id of the order
		myCookie, _ := r.Cookie("myCookie")
		sessionMap.Delete(myCookie.Value)
		id := r.FormValue("order")
		err := sessionMap.InsertTransaction(myCookie.Value, selleruser, id)
		if err != nil {
			logger.Logging(db, "Failed to enter info into sessionMap: orderList")
			return
		}
		http.Redirect(w, r, "/insertQuotation", http.StatusSeeOther)
		return
	}
	err = tpl.ExecuteTemplate(w, "orderList.html", dataSlice)
	if err != nil {
		logger.Logging(connectDB(), "Failed to parse orderList.html template: orderList")
	}
}

func insertQuotation(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	//extract transaction id from sessionMap
	myCookie, _ := r.Cookie("myCookie")
	id, err := sessionMap.SearchTransaction(myCookie.Value)
	if err != nil {
		logger.Logging(connectDB(), "Error in searching for id in sessionMap: insertQuotation")
		return
	}
	seller, err := sessionMap.Search(myCookie.Value)
	if err != nil {
		logger.Logging(connectDB(), "Error in searching for seller user in sessionMap: insertQuotation")
		return
	}
	db := connectDB()
	defer db.Close()
	//search the submissions database
	var phoneinfo submissions.Condition
	phoneinfo, err = submissions.GetID(db, id)
	if err != nil {
		logger.Logging(db, "Error in getting id from submission table: insertQuotation")
		return
	}
	if r.Method == http.MethodPost {
		price := r.FormValue("quotation")
		err = quotation.InsertQuotation(db, phoneinfo.Customer, seller, phoneinfo.ID, price, phoneinfo.Name)
		if err != nil {
			logger.Logging(db, "Error in inserting quotation from seller: insertQuotation")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = tpl.ExecuteTemplate(w, "showDetails.html", phoneinfo)
	if err != nil {
		logger.Logging(db, "Error in executing showDetails.html template: insertQuotation")
	}
}

//submittedOrder shows the list of submitted quotations to the customer
func viewResponse(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	//take the current values from the quotation table and display
	db := connectDB()
	customer := getUsername(r)
	var tableData []quotation.QuoteTable
	tableData, err := quotation.GetCustomerQuote(db, customer)
	if err != nil {
		logger.Logging(db, "Error in getting customer quotation: viewResponse")
		return
	}
	if r.Method == http.MethodPost {
		//insert into pastsubmissions table
		sellerWithID := r.FormValue("choice")
		ss := strings.Split(sellerWithID, "\\")
		tNow := time.Now().String()
		tNow = tNow[:28]
		c, err := submissions.GetID(db, ss[1])
		if err != nil {
			logger.Logging(db, "Error in getting id from submissions table: viewResponse")
		}
		err = transactions.InsertTransaction(db, ss[1], c.Customer, ss[0], c.Name, c.Storage, c.Housing, c.Screen, c.OriginalAccessories, c.OtherIssues, ss[2], tNow)
		if err != nil {
			logger.Logging(db, "Error in inserting data into postSubmissions table: viewResponse")
			return
		}
		//delete from submissions
		err = submissions.Delete(db, ss[1])
		if err != nil {
			logger.Logging(db, "Error in deleting from submissions table: viewResponse")
		}
		//delete from quotations
		err = quotation.Delete(db, ss[1])
		if err != nil {
			logger.Logging(db, "Error in deleting from quotations table: viewResponse")
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	err = tpl.ExecuteTemplate(w, "displayQuotes.html", tableData)
	if err != nil {
		logger.Logging(db, "Error in executing displayQuotes.html template: viewResponse")
	}
}

//sellerViewTransaction obtain data from pastsubmissions table and display
func sellerViewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	username := getUsername(r)
	var err error
	var t []transactions.PSubmissions
	t, err = transactions.GetSeller(db, username)
	if err != nil {
		logger.Logging(db, "Error in getting transaction info from postSubmissions table: sellerViewTransaction")
	}
	err = tpl.ExecuteTemplate(w, "displayPastSubmission.html", t)
	if err != nil {
		logger.Logging(db, "Error in executing displayPastSubmission.html template: sellerViewTransaction")
	}
}

func customerViewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connectDB()
	username := getUsername(r)
	var err error
	var t []transactions.PSubmissions
	t, err = transactions.GetCustomer(db, username)
	if err != nil {
		logger.Logging(db, "Error in getting info from postSubmissions table: customerViewTransaction")
	}
	err = tpl.ExecuteTemplate(w, "displayPastSubmission.html", t)
	if err != nil {
		logger.Logging(db, "Error in executing displayPastSubmission.html template: customerViewTransaction")
	}
}

func forgetPassword(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewV4()
	claiming := &secure.MyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		SessionID: id.String(),
	}
	signedToken, err := secure.GenerateJWT(claiming)
	if err != nil {
		logger.Logging(connectDB(), "Error in generating token from login: login")
		return
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		if _, ok := userMap[email]; !ok {
			tpl.ExecuteTemplate(w, "redirect.html", "Email not found")
			return
		}
		apiKey, err := secretkey.GetKey(connectDB(), "apiKey")
		if err != nil {
			log.Println(err)
			return
		}
		from := mail.NewEmail("Upseller", "gavinerh@gmail.com")
		subject := "Password reset"
		to := mail.NewEmail(email, email)
		plainTextContent := "Click on this link to reset your password: https://localhost:5000/resetpassword?token=" + signedToken + "&user=" + email
		htmlContent := "Please reset within 5mins of receiving this email " + "https://localhost:5000/resetpassword?token=" + signedToken + "&user=" + email
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		client := sendgrid.NewSendClient(apiKey)
		_, err = client.Send(message)
		if err != nil {
			logger.Logging(connectDB(), "Error in sending email to reset password: forgetPassword")
		}
		err = tpl.ExecuteTemplate(w, "redirect.html", "Password reset sent to your email")
		if err != nil {
			logger.Logging(connectDB(), "Error in parsing template redirect.html: forgetPassword")
		}
		return
	}
	//serve the password reset page to get the email for them to enter
	err = tpl.ExecuteTemplate(w, "forgetPassword.html", nil)
	if err != nil {
		logger.Logging(connectDB(), "Error in executing forget password template: forgetPassword")
	}
}

//display html for creating new password for user
func resetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	username := r.FormValue("user") //get the username
	_, err := secure.ParseToken(token)
	if err != nil {
		io.WriteString(w, "Link has expired please reset again")
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}
	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		if !secure.InputValidate(password) {
			tpl.ExecuteTemplate(w, "redirect.html", "Your password should not contain ', \",  <, >, tabs or empty spaces")
			return
		}
		bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			logger.Logging(connectDB(), "Failure in Bcrypting password: SignUp")
			return
		}
		//update the users table and usermap
		err = users.UpdateRecord(connectDB(), string(bPassword), username)
		if err != nil {
			log.Println(err)
			return
		}
		myUser := userMap[username]
		myUser.Password = string(bPassword)
		userMap[username] = myUser
		tpl.ExecuteTemplate(w, "redirect.html", "Please login again")
		return

	}
	err = tpl.ExecuteTemplate(w, "resetPassword.html", nil)
	if err != nil {
		log.Println(err)
	}
}

// var githubConfig = &oauth2.Config{}

//additional login with oauth
// func oauthLogin(w http.ResponseWriter, r *http.Request) {
// 	id := uuid.NewV4()
// 	//save the id in the database
// 	oauth.InsertTempID(connectDB(), id.String())
// 	githubConfig.ClientID = clientID
// 	githubConfig.ClientSecret = clientSecret
// 	githubConfig.Endpoint = github.Endpoint
// 	redirectURL := githubConfig.AuthCodeURL(id.String())
// 	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
// }

//oauth redirect page
// func oauthRedirect(w http.ResponseWriter, r *http.Request) {
// 	// check if the code you sent over to github is the same when redirected back
// 	code := r.FormValue("code")
// 	state := r.FormValue("state")
// 	m, err := oauth.GetTempID(connectDB())
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	if _, ok := m[state]; !ok {
// 		http.Error(w, "State not found", http.StatusBadRequest)
// 		return
// 	}
// 	token, err := githubConfig.Exchange(r.Context(), code)
// 	if err != nil {
// 		http.Error(w, "Could not login", http.StatusInternalServerError)
// 		return
// 	}
// 	ts := githubConfig.TokenSource(r.Context(), token)
// 	client := oauth2.NewClient(r.Context(), ts)
// 	requestBody := strings.NewReader(`{"query": "query {viewer {id}}"}`)
// 	res, err := client.Post("https://api.github.com/graphql", "application/json", requestBody)
// 	if err != nil {
// 		http.Error(w, "Could not get user", http.StatusInternalServerError)
// 		return
// 	}
// 	defer res.Body.Close()
// 	bs, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		http.Error(w, "Could not read github information", http.StatusInternalServerError)
// 		return
// 	}
// 	var gr githubResponse
// 	err = json.NewDecoder(res.Body).Decode(&gr)
// 	if err != nil {
// 		http.Error(w, "Github invalid response", http.StatusInternalServerError)
// 		return
// 	}
// 	githubLoginID := gr.Data.Viewer.ID
// }

// //if oauth is not created, we can continue to do the testing script for main and the packages
// //
