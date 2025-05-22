-- MySQL dump 10.13  Distrib 8.0.42, for Linux (x86_64)
--
-- Host: localhost    Database: gongchang
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `designer_profiles`
--

DROP TABLE IF EXISTS `designer_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `designer_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` varchar(191) DEFAULT NULL,
  `company_name` longtext,
  `address` longtext,
  `website` longtext,
  `bio` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_designer_profiles_user_id` (`user_id`),
  KEY `idx_designer_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `designer_profiles`
--

LOCK TABLES `designer_profiles` WRITE;
/*!40000 ALTER TABLE `designer_profiles` DISABLE KEYS */;
INSERT INTO `designer_profiles` VALUES (1,'2025-05-14 15:49:18.578','2025-05-14 15:49:18.578',NULL,'1','设计工作室1','北京市朝阳区','http://designer1.com','专业服装设计工作室，专注于高端定制');
/*!40000 ALTER TABLE `designer_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `factory_profiles`
--

DROP TABLE IF EXISTS `factory_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `factory_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` varchar(191) DEFAULT NULL,
  `company_name` longtext,
  `address` longtext,
  `capacity` bigint DEFAULT NULL,
  `equipment` longtext,
  `certificates` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_factory_profiles_user_id` (`user_id`),
  KEY `idx_factory_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `factory_profiles`
--

LOCK TABLES `factory_profiles` WRITE;
/*!40000 ALTER TABLE `factory_profiles` DISABLE KEYS */;
INSERT INTO `factory_profiles` VALUES (1,'2025-05-14 15:49:18.581','2025-05-14 15:49:18.581',NULL,'2','服装厂1','广东省深圳市',1000,'全自动裁剪机,工业缝纫机','ISO9001,质量管理体系认证');
/*!40000 ALTER TABLE `factory_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `files` (
  `id` varchar(191) NOT NULL,
  `name` longtext,
  `path` longtext,
  `order_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_files_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `files`
--

LOCK TABLES `files` WRITE;
/*!40000 ALTER TABLE `files` DISABLE KEYS */;
/*!40000 ALTER TABLE `files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_attachments`
--

DROP TABLE IF EXISTS `order_attachments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_attachments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` bigint unsigned NOT NULL,
  `file_name` longtext NOT NULL,
  `file_path` longtext NOT NULL,
  `file_type` longtext NOT NULL,
  `uploaded_by` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_attachments_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_attachments`
--

LOCK TABLES `order_attachments` WRITE;
/*!40000 ALTER TABLE `order_attachments` DISABLE KEYS */;
/*!40000 ALTER TABLE `order_attachments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `order_progresses`
--

DROP TABLE IF EXISTS `order_progresses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_progresses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NOT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_id` bigint unsigned NOT NULL,
  `status` longtext NOT NULL,
  `description` text,
  `created_by` varchar(191) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_progresses_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `order_progresses`
--

LOCK TABLES `order_progresses` WRITE;
/*!40000 ALTER TABLE `order_progresses` DISABLE KEYS */;
/*!40000 ALTER TABLE `order_progresses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` longtext NOT NULL,
  `description` longtext,
  `fabric` longtext,
  `quantity` bigint DEFAULT NULL,
  `factory_id` varchar(191) DEFAULT NULL,
  `status` varchar(191) DEFAULT 'draft',
  `designer_id` longtext,
  `customer_id` longtext,
  `unit_price` double DEFAULT NULL,
  `total_price` double DEFAULT NULL,
  `payment_status` longtext,
  `shipping_address` longtext,
  `order_type` longtext,
  `fabrics` longtext,
  `delivery_date` datetime(3) DEFAULT NULL,
  `order_date` datetime(3) DEFAULT NULL,
  `special_requirements` longtext,
  `attachments` json DEFAULT NULL,
  `models` json DEFAULT NULL,
  `images` json DEFAULT NULL,
  `videos` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `orders`
--

LOCK TABLES `orders` WRITE;
/*!40000 ALTER TABLE `orders` DISABLE KEYS */;
/*!40000 ALTER TABLE `orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS `products`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext NOT NULL,
  `description` longtext,
  `category` longtext NOT NULL,
  `price` double NOT NULL,
  `stock` bigint NOT NULL,
  `status` varchar(191) DEFAULT 'active',
  `created_by` varchar(191) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES `products` WRITE;
/*!40000 ALTER TABLE `products` DISABLE KEYS */;
/*!40000 ALTER TABLE `products` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `supplier_profiles`
--

DROP TABLE IF EXISTS `supplier_profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `supplier_profiles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` varchar(191) DEFAULT NULL,
  `company_name` longtext,
  `address` longtext,
  `main_products` longtext,
  `certificates` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_supplier_profiles_user_id` (`user_id`),
  KEY `idx_supplier_profiles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `supplier_profiles`
--

LOCK TABLES `supplier_profiles` WRITE;
/*!40000 ALTER TABLE `supplier_profiles` DISABLE KEYS */;
INSERT INTO `supplier_profiles` VALUES (1,'2025-05-14 15:49:18.583','2025-05-14 15:49:18.583',NULL,'3','面料供应商1','浙江省绍兴市','棉料,丝绸,化纤','环保认证,质量认证');
/*!40000 ALTER TABLE `supplier_profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` varchar(191) NOT NULL,
  `username` varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `email` longtext NOT NULL,
  `role` varchar(191) DEFAULT 'user',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('451383e4-4636-4dd9-9a9f-f3e2bb312140','supplier1','$2a$10$kQPA8LYqGkWOUzGPAoPOMO/fcdv7KtR4RyZCFWwKdnd.ipSyzxPOi','supplier1@test.com','supplier','2025-05-14 15:49:18.575','2025-05-14 15:49:18.575',NULL),('b968438e-a699-4d96-920c-a9adc3d65906','factory1','$2a$10$kQPA8LYqGkWOUzGPAoPOMO/fcdv7KtR4RyZCFWwKdnd.ipSyzxPOi','factory1@test.com','factory','2025-05-14 15:49:18.573','2025-05-14 15:49:18.573',NULL),('ee505807-15b7-49fe-ad99-d14e33c2b580','designer1','$2a$10$kQPA8LYqGkWOUzGPAoPOMO/fcdv7KtR4RyZCFWwKdnd.ipSyzxPOi','designer1@test.com','designer','2025-05-14 15:49:18.571','2025-05-14 15:49:18.571',NULL);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-05-14  7:55:41
