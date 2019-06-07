package spaceshipsuper

import (
	"time"
	"encoding/json"
	
	"github.com/michaelmcallister/spaceshipsuper/http"
)

const apiEndpointBase = "https://api.spaceship.com.au"

// Client holds the username and password needed to fetch the access token.
// Additionally it also holds some internal state to ensure the accessToken
// remains valid.
type Client struct {
	Username string
	Password string
	accessToken string
	expiry time.Time
}

// Account represents the return structure of the "GetAccount" method.
type Account struct {
	Member struct {
		Title                       string 
		FirstName                   string    `json:"first_name"`
		LastName                    string    `json:"last_name"`
		MemberGroup                 string    `json:"member_group"`
		AccountNumber               string    `json:"account_number"`
		SupertickStatus             string    `json:"supertick_status"`
		Age                         int
		TfnStatus                   string    `json:"tfn_status"`
		Phone                       string 
		AccountBalance              float64   `json:"account_balance"`
		AccountBalanceContributions float64   `json:"account_balance_contributions"`
		AccountBalanceFees          float64   `json:"account_balance_fees"`
		AccountBalanceIncome        float64   `json:"account_balance_income"`
		AccountBalanceInsurance     float64   `json:"account_balance_insurance"`
		AccountBalanceTaxes         float64   `json:"account_balance_taxes"`
		Gender                      string
		Email                       string
		Status                      string
		JoinedDate                  time.Time `json:"joined_date"`
		PostalAddress               struct {
			Line1    string `json:"line_1"`
			Line2    string `json:"line_2"`
			State    string
			Suburb   string
			Postcode string
		} `json:"postal_address"`
		ResidentialAddress struct {
			Line1    string `json:"line_1"`
			Line2    string `json:"line_2"`
			State    string
			Suburb   string
			Postcode string
		} `json:"residential_address"`
		CurrentInvestmentSelection []struct {
			InvestmentOption string `json:"investment_option"`
			Distribution     int
		} `json:"current_investment_selection"`
		Investments []struct {
			InvestmentOption string `json:"investment_option"`
			SellPrice        string `json:"sell_price"`
			Units            string
		}
	}
	Rollovers []string
}

// GetAccount returns a populated pointer to Account, refreshing the
// authentication token if necessary.
func (c *Client) GetAccount() (*Account, error) {
	const acctEndpoint = apiEndpointBase + "/v1/super/account"
	c.refreshAuth(false)

	resp, err := http.DoGet(acctEndpoint, c.accessToken)
	if err != nil {
		return nil, err
	}

	acct := &Account{}
	if err := json.Unmarshal(resp, acct); err != nil {
		return nil, err
	}

	return acct, nil
}