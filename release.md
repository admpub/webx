# Release Notes

## [v0.5.19] — 2026-08-25

### Changed
- 构建工具更新：goreleaser 改用 `GITHUB_TOKEN`
- 构建配置：新增 `config/i18n` 打包到发行版

---

## [v0.5.18] — 2026-08-25

### Added
- 新增后台线下支付管理功能（`/official/customer/offline_pay`）
- 新增前台用户线下充值功能（`recharge_offline.go`）
- 新增地区国家管理功能（`tool/area_country.go`）
- 新增翻译工具（`tool/translation.go`）
- 新增多语言设置（`settings_multilingual.go`）
- 新增 OAuth2 登录支持（COSCMS、GitHub）
- 新增 `ip2region_v6.xdb` IPv6 地理数据库
- 新增英文 README 文档（README.en-US.md）

### Changed
- 评论系统重构（`article/comment.go`）
- 前台路由页面重构（`manager/frontend_route_page.go`）
- 客户分组包管理重构（`group_package/group_package.go`）
- 客户角色管理重构（`customer/role/role.go`）
- 客户等级管理重构（`customer/level/level.go`）
- 广告管理重构（`advert/ad_item.go`, `advert/ad_position.go`）
- 页面模板管理重构（`page/template.go`, `page/navigate.go`）
- 标签管理重构（`tags/tags.go`）
- 地区管理重构（`tool/area.go`, `tool/area_group.go`）
- 前台文章多语言支持（`i18nm.GetModelsTranslations`）
- 图片配置优化（`image/config.go`）
- 用户登录逻辑重构（`index/customer.go`）
- 攻击防护模板优化（`under_attack.html`）
- 大量前端模板重构（167 文件）
- 导航系统全面更新
- 依赖升级（go 1.25.3→1.26.2、webcore v0.13.x、webfront 更新）

---

## [v0.5.17] — 2025-11-28

### Changed
- 文章/分类/标签/导航 handler 重构（迁移至 `formbuilder`）
- 前台文章用户管理重构（`article/user/manage.go`）
- 前端模板优化（52 文件）
- 依赖更新

---

## [v0.5.16] — 2025-08-21

### Changed
- 大量前端模板重构（188 文件）
- 导航系统全面更新
- 前端 CSS/JS 资产更新
- 依赖升级

---

## [v0.5.15] — 2025-08-19

### Added
- 新增 `ip2region_v6.xdb` IPv6 地理数据库
- 前台文章详情新增多语言支持

### Changed
- 图片配置优化（支持移除水印、缓存清理）
- 用户登录逻辑重构（移除 `top.ParseDuration` 依赖）
- 配置文件更新（`config.yaml.sample`）
- 依赖升级（go 1.24.5 → 1.26.2）

---

## [v0.5.13] — 2025-08-19

### Changed
- 前端模板样式调整（导航、用户资料、钱包等）
- 字体文件更新（Merriweather、Ionicons 等）
- 前端 CSS 优化
- 依赖更新

---

## [v0.5.12] — 2025-08-19

### Changed
- 大量前端模板重构（40+ 文件）
- 用户注册/登录表单优化
- 用户个人资料页面重构
- 用户钱包充值页面改进
- 前端 CSS/JS 资产更新
- 依赖升级

---

## [v0.5.11] — 2025-08-19

### Changed
- 文章分类管理功能重构
- 导航管理优化
- 前端忘记密码页面改进
- 用户资料和钱包模板更新
- 依赖更新

---

## [v0.5.8] — 2025-08-13

### Changed
- 依赖包更新（go.mod/go.sum）

---

## [v0.5.6] — 2025-08-06

### Changed
- 前端默认模板和博客模板布局重构
- 404/500 错误页面优化（新增 `404in.html`）
- 依赖更新

---

## [v0.5.5] — 2025-02-11

### Changed
- 客户钱包流水处理优化（`AddRepeatableFlow`）
- 页面模板相关修复

---

## [v0.5.4] — 2025-01-19

### Changed
- 页面模板导航修复
- 前端首页初始化修复

---

## [v0.5.3] — 2025-01-08

### Added
- 用户等级设置新增"金额类型"字段（支持余额/累积总收入）
- 后台仪表盘文本国际化（`echo.T`）
- 账户绑定国际化支持
- 新增"提交"/"重置"翻译键

### Changed
- 客户钱包流水查询优化
- 导航配置更新
