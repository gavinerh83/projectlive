package url

type url struct {
	Static              string
	Login               string
	Signup              string
	Home                string
	CustomerSell        string
	ViewResponse        string
	OrderList           string
	Logout              string
	SellerTransaction   string
	InsertQuotation     string
	CustomerTransaction string
	ForgetPassword      string
	ResetPassword       string
	Search              string
	List                string
}

//ReturnURL returns the url for the controller
func ReturnURL() url {
	var urlPattern url
	urlPattern.Home = "/"                     //show the before and after login page, segregated for customer and sellers
	urlPattern.Login = "/login"               //login form to fill
	urlPattern.Signup = "/signup"             //signup form to fill
	urlPattern.Static = "/static/"            //serves the css files
	urlPattern.CustomerSell = "/customersell" //shows the page to fill up form to sell phone
	urlPattern.ViewResponse = "/viewResponse" //show the page to view responses from sellers
	urlPattern.OrderList = "/orderList"       //show the past transactions and new orders
	urlPattern.Logout = "/logout"
	urlPattern.InsertQuotation = "/insertQuotation"
	urlPattern.SellerTransaction = "/successTransactions"
	urlPattern.CustomerTransaction = "/viewSubmitted"
	urlPattern.ForgetPassword = "/forgetpassword"
	urlPattern.ResetPassword = "/resetpassword"
	urlPattern.Search = "/search"
	urlPattern.List = "/list"
	return urlPattern
}
