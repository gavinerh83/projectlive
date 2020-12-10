package url

type url struct {
	Static           string
	Login            string
	Signup           string
	Home             string
	CustomerSell     string
	CustomerResponse string
	OrderList        string
	Logout           string
	AutoLogout       string
}

//ReturnURL returns the url for the controller
func ReturnURL() url {
	var urlPattern url
	urlPattern.Home = "/"                             //show the before and after login page, segregated for customer and sellers
	urlPattern.Login = "/login"                       //login form to fill
	urlPattern.Signup = "/signup"                     //signup form to fill
	urlPattern.Static = "/static/"                    //serves the css files
	urlPattern.CustomerSell = "/customersell"         //shows the page to fill up form to sell phone
	urlPattern.CustomerResponse = "/customerResponse" //show the page to view responses from sellers
	urlPattern.OrderList = "/orderList"               //show the past transactions and new orders
	urlPattern.Logout = "/logout"
	urlPattern.AutoLogout = "/autologout"
	return urlPattern
}
