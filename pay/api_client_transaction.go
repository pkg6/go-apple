package pay

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/pkg6/go-apple"
)

// GetTransactionHistory Get Transaction History
// Doc: https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (a *ApiClient) GetTransactionHistory(transactionId string, body apple.BodyMap) (resp *ResponseTransactionHistory, err error) {
	resp = new(ResponseTransactionHistory)
	path := fmt.Sprintf("/inApps/v1/history/%s", transactionId) + "?" + body.EncodeURLParams()
	err = a.WithTokenGet(path, nil, &resp)
	return
}

type SignedTransaction string

func (s *SignedTransaction) DecodeSignedTransaction() (ti *TransactionsItem, err error) {
	if *s == "" {
		return nil, fmt.Errorf("signedTransactions is empty")
	}
	ti = new(TransactionsItem)
	if err = ExtractClaims(string(*s), ti); err != nil {
		return nil, err
	}
	return ti, nil
}

// ResponseTransactionHistory
//Doc: HistoryResponse https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type ResponseTransactionHistory struct {
	ResponseErrorMessage
	AppAppleId         int                 `json:"appAppleId"`
	BundleId           string              `json:"bundleId"`
	Environment        string              `json:"environment"`
	HasMore            bool                `json:"hasMore"`
	Revision           string              `json:"revision"`
	SignedTransactions []SignedTransaction `json:"signedTransactions"`
}

// TransactionsItem
// Doc: https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type TransactionsItem struct {
	jwt.StandardClaims
	TransactionId               string `json:"transactionId"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	WebOrderLineItemId          string `json:"webOrderLineItemId"`
	BundleId                    string `json:"bundleId"`
	ProductId                   string `json:"productId"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	PurchaseDate                int64  `json:"purchaseDate"`
	OriginalPurchaseDate        int64  `json:"originalPurchaseDate"`
	ExpiresDate                 int64  `json:"expiresDate"`
	Quantity                    int    `json:"quantity"`
	Type                        string `json:"type"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	SignedDate                  int64  `json:"signedDate"`
	OfferType                   int    `json:"offerType"`
	Environment                 string `json:"environment"`
	AppAccountToken             string `json:"appAccountToken"`
}

// GetTransactionInfo Get Transaction Info
// Doc: https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (a *ApiClient) GetTransactionInfo(transactionId string) (resp *ResponseTransactionInfo, err error) {
	resp = new(ResponseTransactionInfo)
	path := fmt.Sprintf("/inApps/v1/transactions/%s", transactionId)
	err = a.WithTokenGet(path, nil, &resp)
	return
}

// ResponseTransactionInfo
//Doc: https://developer.apple.com/documentation/appstoreserverapi/transactioninforesponse
type ResponseTransactionInfo struct {
	ResponseErrorMessage
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

func (t *ResponseTransactionInfo) DecodeSignedTransaction() (ti *TransactionsItem, err error) {
	if t.SignedTransactionInfo == "" {
		return nil, fmt.Errorf("signedTransactionInfo is empty")
	}
	ti = &TransactionsItem{}
	if err = ExtractClaims(t.SignedTransactionInfo, ti); err != nil {
		return nil, err
	}
	return ti, nil
}
