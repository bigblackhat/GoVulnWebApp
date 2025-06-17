/*
 Navicat Premium Data Transfer

 Source Server         : locahost 本地安装
 Source Server Type    : MySQL
 Source Server Version : 80042
 Source Host           : localhost:3306
 Source Schema         : gwva

 Target Server Type    : MySQL
 Target Server Version : 80042
 File Encoding         : 65001

 Date: 12/06/2025 11:31:30
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for manager
-- ----------------------------
DROP TABLE IF EXISTS `manager`;
CREATE TABLE `manager` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `passwd` varchar(255) NOT NULL,
  `level` int DEFAULT '1',
  `mail` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mail` (`mail`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of manager
-- ----------------------------
BEGIN;
INSERT INTO `manager` (`id`, `name`, `passwd`, `level`, `mail`) VALUES (1, 'admin', '123456', 1, 'admin@qq.com');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
