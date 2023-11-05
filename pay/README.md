## Apple

- App Store Server API：[官方文档](https://developer.apple.com/documentation/appstoreserverapi)

### Apple ApiClient

初始化Apple客户端
~~~
privateKeyByte, err := os.ReadFile(privateKeyFile)
if err != nil {
  return
}
// 初始化通联客户端
// iss：issuer ID
// bid：bundle ID
// kid：private key ID
// privateKeyFile：私钥文件读取后的字节内容
// isProduction：是否是正式环境
api, err := pay.NewApiClient(iss, bid, keyID, privateKeyByte,false)
~~~

### App Store Server API Client Function

- `client.GetTransactionInfo()` => [Get Transaction Info](https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info)
- `client.GetTransactionHistory()` => [Get Transaction History](https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history)
- `client.GetAllSubscriptionStatuses()` =>[get_all_subscription_statuses](https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses)
- `client.SendConsumptionInformation()` => [Send Consumption Information](https://developer.apple.com/documentation/appstoreserverapi/send_consumption_information)
- `client.GetNotificationHistory()` => [Get Notification History](https://developer.apple.com/documentation/appstoreserverapi/get_notification_history)
- `client.LookUpOrderId()` => [Look Up Order ID](https://developer.apple.com/documentation/appstoreserverapi/look_up_order_id)
- `client.GetRefundHistory()` => [Get Refund History](https://developer.apple.com/documentation/appstoreserverapi/get_refund_history)

### Apple Function

* `apple.VerifyReceipt()` => [验证支付凭证](https://developer.apple.com/documentation/appstorereceipts/verifyreceipt)
* `apple.ExtractClaims()` => 解析signedPayload
* `apple.DecodeSignedPayload()` => 解析notification signedPayload
