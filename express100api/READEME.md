#sign

1-所有参数进行按照首字母先后顺序排列
2-把排序后的结果按照参数名+参数值的方式拼接
3-拼装好的字符串首尾拼接 client_secret 进行 md5 加密后转大写，client_secret 的值是申请 API 服务时获取
的 App Secret
4-关于 client_id 和 client_secret 的获取，请前往快递管家官网[申请 API 服务]获取。


#Order

SendOrderData

UpdateOrderData
1.尚未打印出库得订单可修改
2.打印且产生单号得订单不可修改
3.如2有不可修改订单，可以删除订单重新导入


