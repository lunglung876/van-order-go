CREATE TABLE `order` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `distance` int(11) DEFAULT NULL,
  `status` varchar(30) DEFAULT NULL,
  `origin_latitude` varchar(20) DEFAULT NULL,
  `origin_longitude` varchar(20) DEFAULT NULL,
  `destination_latitude` varchar(20) DEFAULT NULL,
  `destination_longitude` varchar(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;