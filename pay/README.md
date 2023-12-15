## 整体流程简述

> 参考文档：https://help.apple.com/app-store-connect/#/devb57be10e7

1. app从服务端获取待支付订单ID（这个订单ID是自己服务端产生的订单信息）和待付费productId（如果在苹果上配置了多个productId，则需要从服务端拉取自己的商品和苹果商品关联信息，这里的是指在苹果的productId）；
2. app根据productId，发起应用内支付；
3. app得到支付成功结果，并将支付结果的receipt信息以及第一步的待支付订单ID发回自己的服务端，服务端调用SDK的Verify，得到校验结果（表示该支付信息是否是苹果处理的）；
4. 如果校验成功则表示苹果支付完成，根据第三步传回的订单ID来处理自己的后续业务逻辑；
 综上：当使用苹果的应用内支付时，其实是由自己的客户端APP来发起的“支付结果通知”请求来推动业务支付流程的数据流转。

## Apple

- App Store Server API：[官方文档](https://developer.apple.com/documentation/appstoreserverapi)

## Apple ApiClient

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


## Apple支付回调用状态说明

> notificationtype: https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
>
> subtype: https://developer.apple.com/documentation/appstoreservernotifications/subtype

| notification_type         | notification_type说明                                        | subtype              | subtype说明                                                  |
| ------------------------- | ------------------------------------------------------------ | -------------------- | ------------------------------------------------------------ |
| DID_CHANGE_RENEWAL_PREF   | 用户对其订阅计划进行了更改                                   | UPGRADE              | 用户升级其订阅,或交叉分级到具有相同持续时间的订阅。升级立即生效,开始新的计费周期,用户将收到上一周期未使用部分的按比例退款 |
|                           |                                                              | DOWNGRADE            | 用户降级其订阅或交叉分级到具有不同持续时间的订阅。降级将在下一个续订日期生效,并且不会影响当前有效的计划 |
|                           |                                                              | 空                   | 则用户将其续订首选项更改回当前订阅,从而有效地取消降级        |
| DID_CHANGE_RENEWAL_STATUS | 用户对其订阅计划进行了更改                                   | AUTO_RENEW_ENABLED   | 则用户重新启用订阅自动续订                                   |
|                           |                                                              | AUTO_RENEW_DISABLED  | 则用户禁用了订阅自动续费,或者用户申请退款后App Store禁用了订阅自动续费 |
| DID_FAIL_TO_RENEW         | 订阅由于计费问题而未能续订(通知用户他们的账单信息可能存在问题。App Store 将在 60 天内继续重试计费,或者直到用户解决计费问题或取消订阅(以先到者为准)) | GRACE_PERIOD         | 则在宽限期内继续提供服务                                     |
|                           |                                                              | 空                   | 则说明订阅不在宽限期内,您可以停止提供订阅服务                |
| DID_RENEW                 | 订阅已成功续订                                               | BILLING_RECOVERY     | 则之前续订失败的过期订阅已成功续订                           |
|                           |                                                              | 空                   | 则活动订阅已成功自动续订新的交易周期。为客户提供对订阅内容或服务的访问权限 |
| EXPIRED                   | 订阅已过期                                                   | VOLUNTARY            | 则订阅在用户禁用订阅续订后过期                               |
|                           |                                                              | BILLING_RETRY        | 则订阅已过期,因为计费重试期已结束,但没有成功的计费事务       |
|                           |                                                              | PRICE_INCREASE       | 则订阅已过期,因为用户不同意需要用户同意的价格上涨            |
|                           |                                                              | PRODUCT_NOT_FOR_SALE | 则订阅已过期,因为在订阅尝试续订时该产品不可购买              |
| GRACE_PERIOD_EXPIRED      | 通知用户他们的账单信息可能存在问题。App Store 将在 60 天内继续重试计费,或者直到用户解决计费问题或取消订阅(以先到者为准)。 |                      | 指示计费宽限期已结束而无需续订订阅,因此您可以关闭对服务或内容的访问。 |
| OFFER_REDEEMED            | 用户兑换了促销优惠或优惠代码                                 | INITIAL_BUY          | 则用户兑换了首次购买的优惠                                   |
|                           |                                                              | RESUBSCRIBE          | 则用户兑换了重新订阅非活动订阅的优惠                         |
|                           |                                                              | UPGRADE              | 则用户兑换了升级其有效订阅的优惠,该优惠立即生效              |
|                           |                                                              | DOWNGRADE            | 则用户兑换了降级其有效订阅的优惠,该优惠将在下一个续订日期生效 |
|                           |                                                              | 无subtype            | 如果用户兑换了其有效订阅的优惠                               |
| PRICE_INCREASE            | 系统已通知用户自动续订订阅价格上涨                           | PENDING              | 如果涨价需要用户同意,是subtype指PENDING用户没有对涨价做出回应,或者ACCEPTED用户已经同意涨价 |
|                           |                                                              | ACCEPTED             | 如果涨价不需要用户同意,那subtype就是ACCEPTED                 |
| REFUND                    | 包含退款交易的时间戳。并标识原始交易和产品。其中包含原因。revocationDateoriginalTransactionIdproductIdrevocationReason |                      | 指示App Store已成功对消费品应用内购买、非消费品应用内购买、自动续订订阅或非续订订阅的交易进行退款 |
| REFUND_DECLINED           |                                                              |                      | 指示App Store拒绝了应用开发者使用以下任一方法发起的退款请求  |
| REFUND_REVERSED           | 此通知类型可适用于任何应用内购买类型:消耗型、非消耗型、非续订订阅和自动续订订阅。对于自动续订订阅,当App Store撤销退款时,续订日期保持不变 |                      | 表明App Store由于客户提出的争议而撤销了之前授予的退款。如果您的应用因相关退款而撤销了内容或服务,则需要恢复它们 |
| RENEWAL_EXTENDED          |                                                              |                      | 指示App Store延长了特定订阅的订阅续订日期。您可以通过调用App Store Server API中的延长订阅续订日期或为所有活跃订阅者延长订阅续订日期来请求订阅续订日期延期 |
| RENEWAL_EXTENSION         | App Store正在尝试通过调用为所有活跃订阅者延长订阅续订日期来延长您请求的订阅续订日期 | SUMMARY              | App Store已完成为所有符合条件的订阅者延长续订日期            |
|                           |                                                              | FAILURE              | 则特定订阅的续订日期延长未成功                               |
| REVOKE                    |                                                              |                      | 指示用户有权通过“家人共享”进行应用内购买的通知类型不再可通过共享进行。当购买者禁用产品的家庭共享、购买者(或家庭成员)离开家庭群组或购买者要求并收到退款时,App Store会发送此通知。您的应用程序也会收到呼叫。家庭共享适用于非消耗性应用内购买和自动续订订阅。 |
| SUBSCRIBED                | 用户订阅了产品                                               | INITIAL_BUY          | 用户要么首次购买了订阅,要么通过家庭共享首次获得了对该订阅的访问权限 |
|                           |                                                              | RESUBSCRIBE          | 用户要么续订了订阅,要么通过家庭共享获得了对相同订阅或同一订阅组内的另一个订阅的访问权限 |
| TEST                      |                                                              |                      | 当您通过调用请求测试通知端点请求时,App Store服务器发送的通知类型。调用该端点来测试您的服务器是否正在接收通知。仅当您提出请求时,您才会收到此通知。有关故障排除信息,请参阅获取测试通知状态端点。 |
