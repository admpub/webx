-- MySQL dump 10.13  Distrib 9.0.1, for macos12.7 (x86_64)
--
-- Host: 127.0.0.1    Database: nging
-- ------------------------------------------------------
-- Server version	9.0.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `official_ad_item`
--

DROP TABLE IF EXISTS `official_ad_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_ad_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '广告名称',
  `publisher_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '广告商ID',
  `position_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '广告位ID',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '广告内容',
  `contype` enum('text','image','video','audio') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'image' COMMENT '内容类型',
  `mode` enum('CPA','CPM','CPC','CPS','CPT') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'CPS' COMMENT '广告模式',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '广告链接',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `start` int unsigned NOT NULL DEFAULT '0' COMMENT '生效起始时间',
  `end` int unsigned NOT NULL DEFAULT '0' COMMENT '生效结束时间',
  `sort` int NOT NULL DEFAULT '500' COMMENT '序号',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ad_item_disabled` (`disabled`),
  KEY `ad_item_position_id` (`position_id`),
  KEY `ad_item_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='广告';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_ad_position`
--

DROP TABLE IF EXISTS `official_ad_position`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_ad_position` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `ident` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '唯一标识',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '位置名称',
  `width` int unsigned NOT NULL DEFAULT '0' COMMENT '宽度',
  `height` int unsigned NOT NULL DEFAULT '0' COMMENT '高度',
  `content` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '默认内容',
  `contype` enum('text','image','video','audio') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'image' COMMENT '内容类型',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '广告链接',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ident` (`ident`),
  KEY `disabled` (`disabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='广告位置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_ad_publisher`
--

DROP TABLE IF EXISTS `official_ad_publisher`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_ad_publisher` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所有者ID',
  `owner_type` enum('user','customer') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '所有者类型(customer-前台客户;user-后台用户)',
  `deposit` decimal(12,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '押金',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='广告主';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_ad_settings`
--

DROP TABLE IF EXISTS `official_ad_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_ad_settings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `advert_id` bigint unsigned NOT NULL COMMENT '广告ID',
  `type` enum('area','age','time','client','gendar') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'area' COMMENT '设置类型(area-地区;age-年龄;time-时段;client-客户端类型;gendar-性别)',
  `value` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '值',
  `v_start` int unsigned NOT NULL DEFAULT '0' COMMENT '起始值',
  `v_end` int unsigned NOT NULL DEFAULT '0' COMMENT '结束值',
  `t_start` int unsigned NOT NULL DEFAULT '0' COMMENT '起始时间',
  `t_end` int unsigned NOT NULL DEFAULT '0' COMMENT '结束时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='广告设置';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_api_account`
--

DROP TABLE IF EXISTS `official_common_api_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_api_account` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_type` enum('user','customer') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '所有者类型(user-后台用户;customer-前台客户)',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所有者ID',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '分组',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口名称',
  `url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口网址(生产环境)',
  `url_dev` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '接口网址(测试环境)',
  `app_id` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'APP ID',
  `app_secret` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'APP 密钥',
  `public_key` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '公钥',
  `encryption` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '加密方式',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '扩展数据(JSON格式)',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `common_api_account_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_area`
--

DROP TABLE IF EXISTS `official_common_area`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_area` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `pid` int unsigned NOT NULL DEFAULT '0' COMMENT '父id',
  `short` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '简称',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `merged` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '全称',
  `level` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '层级(1,2,3-省,市,区县)',
  `pinyin` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '拼音',
  `code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '长途区号',
  `zip` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮编',
  `first` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '首字母',
  `lng` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '经度',
  `lat` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '纬度',
  `country_abbr` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'CN' COMMENT '国家缩写',
  PRIMARY KEY (`id`),
  KEY `common_area_pid` (`pid`),
  KEY `common_area_pinyin` (`pinyin`),
  KEY `common_area_first` (`first`),
  KEY `common_area_country_abbr` (`country_abbr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='地区表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_area_group`
--

DROP TABLE IF EXISTS `official_common_area_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_area_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `country_abbr` char(2) COLLATE utf8mb4_general_ci NOT NULL COMMENT '国家缩写',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名称',
  `abbr` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '组缩写',
  `area_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '根地区ID',
  `sort` int NOT NULL DEFAULT '0' COMMENT '排序编号',
  PRIMARY KEY (`id`),
  UNIQUE KEY `common_area_group_uniq` (`country_abbr`,`abbr`),
  KEY `common_area_group_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='地区分组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_article`
--

DROP TABLE IF EXISTS `official_common_article`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_article` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `category1` int unsigned NOT NULL DEFAULT '0' COMMENT '顶级分类ID',
  `category2` int unsigned NOT NULL DEFAULT '0' COMMENT '二级分类ID',
  `category3` int unsigned NOT NULL DEFAULT '0' COMMENT '三级分类ID',
  `category_id` int unsigned NOT NULL DEFAULT '0' COMMENT '最底层分类ID',
  `source_id` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源ID(空代表不限)',
  `source_table` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源表(不含official_前缀)',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '新闻发布者',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '所有者类型(customer-前台客户;user-后台用户)',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '新闻标题',
  `keywords` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '关键词',
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '缩略图',
  `image_original` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '原始图',
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '摘要',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '内容',
  `contype` enum('text','html','markdown') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'markdown' COMMENT '内容类型',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `display` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否显示',
  `template` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模版',
  `comments` bigint unsigned NOT NULL DEFAULT '0' COMMENT '评论数量',
  `close_comment` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '关闭评论',
  `comment_auto_display` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '自动显示评论',
  `comment_allow_user` enum('all','buyer','author','admin','allAgent','curAgent','none','designated') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'all' COMMENT '允许评论的用户(all-所有人;buyer-当前商品买家;author-当前文章作者;admin-管理员;allAgent-所有代理;curAgent-当前产品代理;none-无人;designated-指定人员)',
  `likes` bigint unsigned NOT NULL DEFAULT '0' COMMENT '好评数量',
  `hates` bigint unsigned NOT NULL DEFAULT '0' COMMENT '差评数量',
  `views` bigint unsigned NOT NULL DEFAULT '0' COMMENT '浏览次数',
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标签',
  `price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '价格',
  `slugify` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'SEO-friendly URLs with Slugify',
  PRIMARY KEY (`id`),
  KEY `common_article_category1` (`category1`),
  KEY `common_article_category2` (`category2`),
  KEY `common_article_category3` (`category3`),
  KEY `common_article_category_id` (`category_id`),
  KEY `common_article_source` (`source_table`,`source_id`),
  KEY `common_article_display` (`display`),
  KEY `common_article_slugify` (`slugify`),
  KEY `common_article_owner` (`owner_id`,`owner_type`,`created`),
  KEY `common_article_likes` (`likes` DESC),
  KEY `common_article_comments` (`comments` DESC),
  KEY `common_article_updated` (`updated` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='官方新闻';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_category`
--

DROP TABLE IF EXISTS `official_common_category`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '上级分类ID',
  `has_child` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否有子分类',
  `level` int unsigned NOT NULL DEFAULT '0' COMMENT '层级',
  `name` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '分类名称',
  `keywords` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类页面关键词',
  `description` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类说明',
  `cover` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类封面图',
  `type` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'article' COMMENT '类型',
  `sort` int NOT NULL DEFAULT '5000' COMMENT '排序编号(从小到大)',
  `template` varchar(120) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分类列表页模版',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `show_on_menu` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否(Y/N)显示在导航菜单上',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `slugify` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'SEO-friendly URLs with Slugify',
  PRIMARY KEY (`id`),
  KEY `common_category_parent_id` (`parent_id`),
  KEY `common_category_disabled` (`disabled`),
  KEY `common_category_sort` (`sort`,`id`),
  KEY `common_category_show_on_menu` (`show_on_menu`),
  KEY `common_category_slugify` (`slugify`),
  KEY `common_category_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='分类';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_click_flow`
--

DROP TABLE IF EXISTS `official_common_click_flow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_click_flow` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `target_type` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'article' COMMENT '目标类型',
  `target_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '目标ID',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '用户类型(customer-前台客户;user-后台用户)',
  `type` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型(例如:like,hate)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `common_click_flow_uniqid` (`target_type`,`target_id`,`owner_id`,`owner_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='点击流水记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_comment`
--

DROP TABLE IF EXISTS `official_common_comment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_comment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `reply_comment_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '回复评论ID',
  `reply_owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '回复用户ID',
  `reply_owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '回复用户类型(customer-前台客户;user-后台用户)',
  `root_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '根评论ID',
  `target_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'article' COMMENT '评论目标类型(article,product...)',
  `target_subtype` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '评论目标子类型',
  `target_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '评论目标ID',
  `target_owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '目标作者ID',
  `target_owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '目标作者类型(customer-前台客户;user-后台用户)',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '评论者ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '评论者类型(customer-前台客户;user-后台用户)',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论内容',
  `contype` enum('text','html','markdown') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'text' COMMENT '内容类型',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '编辑时间',
  `display` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '显示',
  `level` int unsigned NOT NULL DEFAULT '0' COMMENT '层数',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '路径',
  `replies` bigint unsigned NOT NULL DEFAULT '0' COMMENT '回复数',
  `likes` bigint unsigned NOT NULL DEFAULT '0' COMMENT '喜欢数量',
  `hates` bigint unsigned NOT NULL DEFAULT '0' COMMENT '不喜欢数量',
  PRIMARY KEY (`id`),
  KEY `common_comment_owner` (`owner_type`,`owner_id`,`created`),
  KEY `common_comment_target` (`target_type`,`target_subtype`,`target_id`,`target_owner_id`,`target_owner_type`),
  KEY `common_comment_display` (`display`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='评论表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_complaint`
--

DROP TABLE IF EXISTS `official_common_complaint`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_complaint` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户ID',
  `target_name` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '对象名称',
  `target_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '投诉对象ID',
  `target_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '投诉对象类型',
  `target_ident` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '投诉对象标识',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '投诉类型',
  `content` mediumtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '投诉内容',
  `process` enum('idle','reject','done','queue') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'idle' COMMENT '处理状态(idle-空闲;reject-驳回;done-已处理;queue-等待处理中)',
  `result` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '处理结果说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='投诉信息';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_friendlink`
--

DROP TABLE IF EXISTS `official_common_friendlink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_friendlink` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `category_id` int unsigned NOT NULL DEFAULT '0' COMMENT '分类',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户ID',
  `logo` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'LOGO',
  `logo_original` varchar(150) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'LOGO原图',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网站名称',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网站说明',
  `url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网址',
  `host` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网址主机名(域名)',
  `verify_time` int unsigned NOT NULL DEFAULT '0' COMMENT '验证时间',
  `verify_fail_count` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '验证失败次数',
  `verify_result` enum('ok','invalid','none') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'none' COMMENT '验证结果(ok-成功;invalid-无效;none-未验证)',
  `process` enum('idle','success','reject') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'idle' COMMENT '处理结果(idle-待处理;success-成功;reject-拒绝)',
  `process_remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '处理备注',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `return_time` int unsigned NOT NULL DEFAULT '0' COMMENT '回访时间',
  `return_count` int unsigned NOT NULL DEFAULT '0' COMMENT '回访次数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `common_friendlink _host` (`host`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='友情链接';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_group`
--

DROP TABLE IF EXISTS `official_common_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '上级ID',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `name` varchar(120) COLLATE utf8mb4_general_ci NOT NULL COMMENT '组名',
  `type` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '类型(customer-客户组;cert-证书组;order-订单组;product-产品组;attr-产品属性组;openapp-开放平台应用;api-外部接口组)',
  `description` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `common_group_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='分组';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_message`
--

DROP TABLE IF EXISTS `official_common_message`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `type` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '消息类型',
  `customer_a` bigint unsigned NOT NULL DEFAULT '0' COMMENT '发信人ID(0为系统消息)',
  `customer_b` bigint unsigned NOT NULL DEFAULT '0' COMMENT '收信人ID',
  `customer_group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '客户组消息',
  `user_a` int unsigned NOT NULL DEFAULT '0' COMMENT '发信人ID(后台用户ID，用于系统消息)',
  `user_b` int unsigned NOT NULL DEFAULT '0' COMMENT '收信人ID(后台用户ID，用于后台消息)',
  `user_role_id` int unsigned NOT NULL DEFAULT '0' COMMENT '后台角色消息',
  `title` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '消息标题',
  `content` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息内容',
  `contype` enum('text','html','markdown') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'text' COMMENT '内容类型',
  `encrypted` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否为加密消息',
  `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码',
  `created` int unsigned NOT NULL COMMENT '发送时间',
  `url` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网址',
  `root_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '根ID',
  `reply_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '回复ID',
  `has_new_reply` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '是否(1/0)有新回复',
  `view_progress` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '查看总进度(100为100%)',
  PRIMARY KEY (`id`),
  KEY `common_message_from` (`customer_a`,`user_a`),
  KEY `common_message_to_customer` (`customer_b`,`customer_group_id`),
  KEY `common_message_to_user` (`user_b`,`user_role_id`),
  KEY `common_message_encrypted` (`encrypted`),
  KEY `common_message_view_progress` (`view_progress`),
  KEY `common_message_has_new_reply` (`has_new_reply` DESC),
  FULLTEXT KEY `common_message_title_content` (`title`,`content`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='站内信';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_message_viewed`
--

DROP TABLE IF EXISTS `official_common_message_viewed`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_message_viewed` (
  `message_id` bigint unsigned NOT NULL COMMENT '消息ID',
  `viewer_id` bigint unsigned NOT NULL COMMENT '浏览者ID',
  `viewer_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '浏览者类型',
  `created` int unsigned NOT NULL COMMENT '查看时间',
  PRIMARY KEY (`message_id`,`viewer_id`,`viewer_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='消息浏览记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_navigate`
--

DROP TABLE IF EXISTS `official_common_navigate`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_navigate` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `type` enum('default','userCenter','backend','other') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'default' COMMENT '导航类型(default-前后默认;userCenter-用户中心;backend-后台)',
  `link_type` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'custom' COMMENT '菜单类型(category-分类;custom-自定义链接)',
  `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '上级ID',
  `has_child` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否有子菜单',
  `level` int unsigned NOT NULL DEFAULT '0' COMMENT '层级',
  `title` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '菜单标题',
  `cover` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片封面',
  `url` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '网址',
  `ident` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标识',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `sort` int NOT NULL DEFAULT '5000' COMMENT '排序',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `target` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '打开目标(_self/_blank/_parent/_top)',
  `direction` enum('X','Y') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '非自定义链接的排列方向(X-横向;Y-纵向)',
  `badge` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '徽标文本',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='导航连接';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_remark`
--

DROP TABLE IF EXISTS `official_common_remark`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_remark` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所有者ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '所有者类型(customer-前台客户;user-后台用户)',
  `source_type` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源类型(组)',
  `source_table` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源表',
  `source_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '来源ID',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '简短描述',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='备注';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_route_page`
--

DROP TABLE IF EXISTS `official_common_route_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_route_page` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(120) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '页面名称',
  `route` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '路由网址',
  `method` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'GET' COMMENT '路由方法(GET/POST/PUT...)',
  `page_content` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '页面内容',
  `page_vars` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '页面变量(JSON)',
  `page_type` enum('html','json','text','xml','redirect') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'html' COMMENT '页面类型',
  `page_id` int unsigned NOT NULL DEFAULT '0' COMMENT '页面ID(可选,0为不关联)',
  `template_enabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否使用模板',
  `template_file` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模板文件(位于route_page文件夹)',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='自定义路由页面';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_route_rewrite`
--

DROP TABLE IF EXISTS `official_common_route_rewrite`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_route_rewrite` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `route` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '路由原网址',
  `rewrite_to` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '重写为网址',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '规则名称',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='网址重写规则';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_sensitive`
--

DROP TABLE IF EXISTS `official_common_sensitive`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_sensitive` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `words` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '敏感词',
  `type` enum('bad','noise') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'bad' COMMENT '类型(bad-敏感词;noise-噪音词)',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='敏感词';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_common_tags`
--

DROP TABLE IF EXISTS `official_common_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_common_tags` (
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '标签名',
  `num` bigint unsigned NOT NULL DEFAULT '0' COMMENT '数量',
  `group` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '分组标识',
  `display` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否显示',
  PRIMARY KEY (`name`,`group`),
  KEY `common_tags_group` (`group`,`display`,`num` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='标签库';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer`
--

DROP TABLE IF EXISTS `official_customer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '系统用户ID',
  `group_id` int unsigned NOT NULL DEFAULT '0' COMMENT '分组ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `password` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`密码',
  `salt` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`盐值',
  `safe_pwd` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`安全密码',
  `session_id` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '`omit:encode`session id',
  `real_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '真实姓名',
  `mobile` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `mobile_bind` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '手机是否已绑定',
  `email` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '邮箱',
  `email_bind` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '邮箱是否已绑定',
  `online` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否在线',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `gender` enum('male','female','secret') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'secret' COMMENT '性别(male-男;female-女;secret-保密)',
  `id_card_no` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '身份证号',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `avatar` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `licenses` bigint unsigned NOT NULL DEFAULT '0' COMMENT '有效证书数量',
  `login_fails` int unsigned NOT NULL DEFAULT '0' COMMENT '连续登录失败次数',
  `level_id` int unsigned NOT NULL DEFAULT '0' COMMENT '客户等级',
  `agent_level` int unsigned NOT NULL DEFAULT '0' COMMENT '代理等级',
  `inviter_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '推荐人(代理)ID',
  `following` int unsigned NOT NULL DEFAULT '0' COMMENT '我关注的人数',
  `followers` int unsigned NOT NULL DEFAULT '0' COMMENT '关注我的人数',
  `role_ids` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色ID(多个用“,”分隔开)',
  `file_size` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上传文件总大小',
  `file_num` bigint unsigned NOT NULL DEFAULT '0' COMMENT '上传文件数量',
  `registered_by` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '注册来源',
  PRIMARY KEY (`id`),
  KEY `customer_name` (`name`),
  KEY `customer_mobile` (`mobile`,`mobile_bind`),
  KEY `customer_email` (`email`,`email_bind`),
  KEY `customer_disabled` (`disabled`),
  KEY `customer_updated` (`updated` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户资料';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_agent_level`
--

DROP TABLE IF EXISTS `official_customer_agent_level`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_agent_level` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '等级名称',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '等级说明',
  `agency_fee` decimal(18,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '代理费',
  `agency_fee_rebate_ratio` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '下级代理费回扣比例',
  `sales_commission_ratio_1` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(1级下线)',
  `sales_commission_ratio_2` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(2级下线)',
  `sales_commission_ratio_3` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(3级下线)',
  `sales_commission_ratio_4` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(4级下线)',
  `sales_commission_ratio_5` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(5级下线)',
  `sales_commission_ratio_6` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(6级下线)',
  `sales_commission_ratio_7` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(7级下线)',
  `sales_commission_ratio_8` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(8级下线)',
  `sales_commission_ratio_9` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(9级下线)',
  `sales_commission_ratio_10` decimal(4,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '销售提成比例(10级下线)',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `role_ids` varchar(225) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色ID(多个用“,”分隔开)',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='代理等级';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_agent_product`
--

DROP TABLE IF EXISTS `official_customer_agent_product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_agent_product` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `agent_id` bigint unsigned NOT NULL COMMENT '代理商UID(customer表id)',
  `product_id` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商品ID',
  `product_table` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '商品表名称(不含official_前缀)',
  `sold` bigint unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `performance` decimal(23,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '业绩',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `expired` int unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  PRIMARY KEY (`id`),
  KEY `customer_agent_product_agent_id_product_id_table` (`agent_id`,`product_id`,`product_table`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='代理产品列表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_agent_profile`
--

DROP TABLE IF EXISTS `official_customer_agent_profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_agent_profile` (
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `earning_balance` decimal(20,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '可提现收入金额',
  `freeze_amount` decimal(20,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '冻结金额(提现中)',
  `margin_amount` bigint unsigned NOT NULL DEFAULT '0' COMMENT '保证金金额',
  `sold` bigint unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `members` bigint unsigned NOT NULL DEFAULT '0' COMMENT '成员统计',
  `status` enum('idle','pending','paid','unconfirm','success','reject','cheat','signedContract','unsignedContract') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'pending' COMMENT '状态(idle:空闲/草稿;pending:待付款;paid:已付款;unconfirm:未确认;success:申请成功;reject:拒绝;cheat:作弊封号;unsignedContract-未签合同;signedContract-已签合同)',
  `apply_level` int unsigned NOT NULL DEFAULT '0' COMMENT '申请代理等级ID',
  `remark` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '备注',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '申请时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '资料更新时间',
  PRIMARY KEY (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_agent_recv`
--

DROP TABLE IF EXISTS `official_customer_agent_recv`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_agent_recv` (
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `recv_money_method` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'alipay' COMMENT '收款方式(提现)',
  `recv_money_branch` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款银行分行(收款方式为银行类时有效)',
  `recv_money_account` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款账号(提现)',
  `recv_money_owner` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '收款人户名(提现)',
  `recv_money_times` int unsigned NOT NULL DEFAULT '0' COMMENT '收款次数(提现)',
  `recv_money_total` decimal(22,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '提现总额',
  `recv_contract_address` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合同接收地址',
  `recv_contract_addressee` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合同接收人姓名',
  `recv_contract_tel` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '合同接收人电话',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '申请时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '资料更新时间',
  PRIMARY KEY (`customer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_counter`
--

DROP TABLE IF EXISTS `official_customer_counter`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_counter` (
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `target` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '目标',
  `total` bigint unsigned NOT NULL DEFAULT '0' COMMENT '统计',
  UNIQUE KEY `customer_counter_key` (`customer_id`,`target`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户的其它数据计数';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_device`
--

DROP TABLE IF EXISTS `official_customer_device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_device` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户ID',
  `session_id` char(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'session id',
  `scense` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '场景标识',
  `platform` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '系统平台',
  `device_no` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '设备编号',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '登录时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `expired` int unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_device_customer_id` (`customer_id`,`scense`,`platform`,`device_no`),
  KEY `customer_device_updated` (`updated`),
  KEY `customer_device_expired` (`expired`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户登录设备';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_following`
--

DROP TABLE IF EXISTS `official_customer_following`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_following` (
  `customer_a` bigint unsigned NOT NULL COMMENT '关注人ID',
  `customer_b` bigint unsigned NOT NULL COMMENT '被关注人ID',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `mutual` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否相互关注',
  UNIQUE KEY `customer_following_uniqid` (`customer_a`,`customer_b`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='关注';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_group_package`
--

DROP TABLE IF EXISTS `official_customer_group_package`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_group_package` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `group` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '等级组',
  `title` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '标题',
  `description` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '说明',
  `price` decimal(12,2) unsigned NOT NULL COMMENT '价格',
  `time_duration` int unsigned NOT NULL DEFAULT '0' COMMENT '时间长度',
  `time_unit` enum('day','week','month','year','forever') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'forever' COMMENT '时间单位',
  `sort` int NOT NULL DEFAULT '5000' COMMENT '排序',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `recommend` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)推荐',
  `sold` int unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `icon_image` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图标图片',
  `icon_class` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图标class',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='等级组套餐价格';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_invitation`
--

DROP TABLE IF EXISTS `official_customer_invitation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_invitation` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'user' COMMENT '创建者类型',
  `code` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '邀请码',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `start` int unsigned NOT NULL DEFAULT '0' COMMENT '有效时间',
  `end` int unsigned NOT NULL DEFAULT '0' COMMENT '失效时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `level_id` int unsigned NOT NULL DEFAULT '0' COMMENT '客户等级ID',
  `agent_level_id` int unsigned NOT NULL DEFAULT '0' COMMENT '代理等级ID',
  `role_ids` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '注册为角色(多个用“,”分隔开)',
  `used_num` int unsigned NOT NULL DEFAULT '0' COMMENT '已使用次数',
  `allow_num` int unsigned NOT NULL DEFAULT '1' COMMENT '剩余允许使用次数',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_invitation_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='邀请码';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_invitation_used`
--

DROP TABLE IF EXISTS `official_customer_invitation_used`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_invitation_used` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户ID',
  `invitation_id` int unsigned NOT NULL DEFAULT '0' COMMENT '邀请码ID',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `level_id` int unsigned NOT NULL DEFAULT '0' COMMENT '客户等级ID',
  `agent_level_id` int unsigned NOT NULL DEFAULT '0' COMMENT '代理等级ID',
  `role_ids` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '注册为角色(多个用“,”分隔开)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_invitation_used _id` (`invitation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='邀请客户';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_level`
--

DROP TABLE IF EXISTS `official_customer_level`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_level` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(30) COLLATE utf8mb4_general_ci NOT NULL COMMENT '等级名称',
  `short` varchar(10) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '等级简称',
  `description` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '等级简介',
  `icon_image` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图标图片',
  `icon_class` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '图片class名',
  `color` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '颜色',
  `bgcolor` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '背景色',
  `price` decimal(10,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '升级价格(0为免费)',
  `integral_asset` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'integral' COMMENT '当作升级积分的资产',
  `integral_min` decimal(30,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '最小积分',
  `integral_max` decimal(30,3) unsigned NOT NULL DEFAULT '0.000' COMMENT '最大积分',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `score` int NOT NULL DEFAULT '50000' COMMENT '分值(分值越大等级越高)',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '扩展配置(JSON)',
  `group` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'base' COMMENT '扩展组(base-基础组,其它名称为扩展组。客户只能有一个基础组等级,可以有多个扩展组等级)',
  `role_ids` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '角色ID(多个用“,”分隔开)',
  PRIMARY KEY (`id`),
  KEY `customer_level_group` (`group`),
  KEY `customer_level_score` (`score` DESC),
  KEY `customer_level_disabled` (`disabled`,`group`,`price`,`integral_asset`,`integral_min`,`integral_max`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户等级';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_level_relation`
--

DROP TABLE IF EXISTS `official_customer_level_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_level_relation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `level_id` int unsigned NOT NULL COMMENT '等级ID',
  `status` enum('actived','expired') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'actived' COMMENT '状态(actived-有效;expired-已过期)',
  `expired` int unsigned NOT NULL DEFAULT '0' COMMENT '过期时间(0为永不过期)',
  `accumulated_days` int unsigned NOT NULL DEFAULT '0' COMMENT '累计天数',
  `last_renewal_at` int unsigned NOT NULL DEFAULT '0' COMMENT '最近续费时间',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int unsigned NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_level_relation_customer_level` (`customer_id`,`level_id`),
  KEY `customer_level_relation_status` (`status`,`expired`),
  KEY `customer_level_relation_updated` (`updated` DESC),
  KEY `customer_level_relation_last` (`last_renewal_at` DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户等级关联';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_oauth`
--

DROP TABLE IF EXISTS `official_customer_oauth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_oauth` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户名',
  `nick_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `union_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'UNION ID',
  `open_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'OPEN ID',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'wechat' COMMENT '类型(例如:wechat/qq/alipay)',
  `avatar` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '头像',
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'E-mail',
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机号',
  `access_token` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Access Token',
  `refresh_token` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Refresh Token',
  `expired` int unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_oauth_uniqid` (`customer_id`,`union_id`,`open_id`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='第三方登录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_prepaid_card`
--

DROP TABLE IF EXISTS `official_customer_prepaid_card`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_prepaid_card` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int unsigned NOT NULL DEFAULT '0' COMMENT '创建者ID',
  `customer_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '客户ID(使用者)',
  `amount` int unsigned NOT NULL COMMENT '面值',
  `sale_price` decimal(12,2) unsigned NOT NULL COMMENT '售价',
  `number` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '卡号',
  `password` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '卡密',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `start` int unsigned NOT NULL DEFAULT '0' COMMENT '有效时间',
  `end` int unsigned NOT NULL DEFAULT '0' COMMENT '失效时间',
  `used` int unsigned NOT NULL DEFAULT '0' COMMENT '使用时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `bg_image` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '背景图片',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_prepaid_card_number` (`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='充值卡';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_role`
--

DROP TABLE IF EXISTS `official_customer_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `description` tinytext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '说明',
  `created` int unsigned NOT NULL COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `is_default` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否作为新用户注册时的默认角色',
  `parent_id` int unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  PRIMARY KEY (`id`),
  KEY `customer_role_disabled` (`disabled`),
  KEY `customer_role_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户角色';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_role_permission`
--

DROP TABLE IF EXISTS `official_customer_role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_role_permission` (
  `role_id` int unsigned NOT NULL COMMENT '角色ID',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '权限类型',
  `permission` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '权限值',
  UNIQUE KEY `customer_role_permission_uniqid` (`role_id`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_u2f`
--

DROP TABLE IF EXISTS `official_customer_u2f`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_u2f` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '签名',
  `type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '类型',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '扩展设置',
  `step` tinyint unsigned NOT NULL DEFAULT '2' COMMENT '第几步',
  `precondition` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'password' COMMENT '除了密码登录外的其它前置条件(仅step=2时有效),用半角逗号分隔',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '绑定时间',
  PRIMARY KEY (`id`),
  KEY `customer_u2f_uid_typ_stepe` (`customer_id`,`type`,`step`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='两步验证';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_wallet`
--

DROP TABLE IF EXISTS `official_customer_wallet`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_wallet` (
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `asset_type` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'money' COMMENT '资产类型(money-钱;point-点数;credit-信用分;integral-积分;gold-金币;silver-银币;copper-铜币;experience-经验)',
  `balance` decimal(30,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '余额',
  `freeze` decimal(30,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '冻结金额',
  `accumulated` decimal(30,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '累计总金额',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`customer_id`,`asset_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='钱包';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_customer_wallet_flow`
--

DROP TABLE IF EXISTS `official_customer_wallet_flow`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_customer_wallet_flow` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `customer_id` bigint unsigned NOT NULL COMMENT '客户ID',
  `asset_type` char(10) COLLATE utf8mb4_general_ci NOT NULL COMMENT '资产类型',
  `amount_type` enum('balance','freeze') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'balance' COMMENT '金额类型(balance-余额;freeze-冻结额)',
  `amount` decimal(30,4) NOT NULL COMMENT '金额(正数为收入;负数为支出)',
  `wallet_amount` decimal(30,4) NOT NULL DEFAULT '0.0000' COMMENT '变动后钱包总金额',
  `source_customer` bigint unsigned NOT NULL DEFAULT '0' COMMENT '来自谁',
  `source_type` varchar(60) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源类型(组)',
  `source_table` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源表(来自物品表)',
  `source_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '来源ID(来自物品ID)',
  `number` bigint unsigned NOT NULL DEFAULT '0' COMMENT '备用编号',
  `trade_no` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '交易号(来自哪个交易)',
  `status` enum('pending','confirmed','refunded','failed','succeed','canceled') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'confirmed' COMMENT '状态(pending-待确认;confirmed-已确认;canceled-已取消)',
  `description` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '简短描述',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `customer_wallet_flow_uniqid` (`customer_id`,`asset_type`,`amount_type`,`source_type`,`source_table`,`source_id`,`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='钱包流水记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_page`
--

DROP TABLE IF EXISTS `official_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_page` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '名称',
  `ident` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '英文唯一标识',
  `template` varchar(200) COLLATE utf8mb4_general_ci NOT NULL COMMENT '模版文件',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `page_ident` (`ident`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='页面布局';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_page_block`
--

DROP TABLE IF EXISTS `official_page_block`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_page_block` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '区块名称',
  `style` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '区块自定义样式',
  `with_items` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '包含项目',
  `item_configs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '项目配置',
  `template` varchar(200) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '模版文件',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `page_block_disabled` (`disabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='页面区块';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_page_layout`
--

DROP TABLE IF EXISTS `official_page_layout`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_page_layout` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `block_id` int unsigned NOT NULL COMMENT '区块ID',
  `page_id` int unsigned NOT NULL COMMENT '页面ID',
  `configs` text COLLATE utf8mb4_general_ci NOT NULL COMMENT '区块在布局中的配置',
  `sort` int NOT NULL DEFAULT '5000' COMMENT '排序',
  `disabled` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否禁用',
  `created` int unsigned NOT NULL COMMENT '添加时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `page_layout_page_id` (`page_id`),
  KEY `page_layout_block_id` (`block_id`),
  KEY `page_layout_disabed` (`disabled`),
  KEY `page_layout_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='页面布局所含区块';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_short_url`
--

DROP TABLE IF EXISTS `official_short_url`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_short_url` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '添加者ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '所有者类型(customer-前台客户;user-后台用户)',
  `long_url` varchar(10240) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '长网址',
  `long_hash` char(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '长网址MD5值',
  `short_url` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '短网址',
  `domain_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '域名ID',
  `visited` int unsigned NOT NULL DEFAULT '0' COMMENT '最近访问时间',
  `visits` bigint unsigned NOT NULL DEFAULT '0' COMMENT '访问次数',
  `available` enum('Y','N') CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'Y' COMMENT '是否有效',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  `expired` int unsigned NOT NULL DEFAULT '0' COMMENT '过期时间(0为不限制)',
  `password` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '访问密码md5(空代表无需密码)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `short_url_short_url` (`short_url`),
  KEY `short_url_available` (`short_url`,`available`),
  KEY `short_url_long_hash` (`long_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='创建时间';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_short_url_domain`
--

DROP TABLE IF EXISTS `official_short_url_domain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_short_url_domain` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所有者客户ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '所有者类型(customer-前台客户;user-后台用户)',
  `domain` varchar(120) COLLATE utf8mb4_general_ci NOT NULL COMMENT '域名',
  `url_count` bigint unsigned NOT NULL DEFAULT '0' COMMENT '网址统计',
  `disabled` enum('Y','N') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'N' COMMENT '是否(Y/N)禁用',
  `created` int unsigned NOT NULL COMMENT '创建时间',
  `updated` int unsigned NOT NULL DEFAULT '0' COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `short_url_domain` (`domain`),
  KEY `short_url_domain_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='短网址域名';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `official_short_url_visit`
--

DROP TABLE IF EXISTS `official_short_url_visit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `official_short_url_visit` (
  `owner_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '所有者ID',
  `owner_type` enum('user','customer') COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'customer' COMMENT '所有者类型(customer-前台客户;user-后台用户)',
  `url_id` bigint unsigned NOT NULL COMMENT '网址ID',
  `domain_id` bigint unsigned NOT NULL COMMENT '域名ID',
  `year` mediumint unsigned NOT NULL COMMENT '年',
  `month` tinyint unsigned NOT NULL COMMENT '月',
  `day` tinyint unsigned NOT NULL COMMENT '日',
  `hour` tinyint unsigned NOT NULL COMMENT '时',
  `ip` varchar(20) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'IP',
  `referer` varchar(120) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '来源网址',
  `language` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '语言',
  `country` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '国家',
  `region` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '行政区',
  `province` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '省份',
  `city` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '城市',
  `isp` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'ISP网络',
  `os` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作系统',
  `os_version` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '操作系统版本',
  `browser` varchar(30) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器',
  `browser_version` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '浏览器版本',
  `created` int NOT NULL COMMENT '创建时间',
  KEY `short_url_visit_created` (`created`),
  KEY `short_url_visit_owner_id` (`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='网址访问日志';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-30 14:46:36
