1. 扫码登录功能中生成二维码的数据新增两种模式：缓存模式和默认模式，可以通过在配置文件中的 `extend` 节点中添加 `QRSignIn:{case:"cache"}` 来指定 
2. 会员中心的 `扫二维码` 功能除了用来扫码登录外，还支持其它二维码的识别，也可以通过调用 `user.RegisterQRCodeDecoder(name, decoder)` 来扩展自己的二维码数据处理逻辑
3. 更新依赖到最新版本