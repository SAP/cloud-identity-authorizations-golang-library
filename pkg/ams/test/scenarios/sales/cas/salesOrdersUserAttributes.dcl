POLICY readSalesOrdersUserCountryCode {
	GRANT read ON salesOrders WHERE CountryCode = $user.country;
}

POLICY writeSalesOrdersUserName {
	GRANT write ON salesOrders WHERE changedBy = $user."name";
}


/// TESTS


TEST readSalesOrdersUserCountryCode {
	GRANT read ON salesOrders POLICY readSalesOrdersUserCountryCode INPUT 
	{
		$user: {
			country: 'DE'
		},
		CountryCode: 'DE'
	};
}

TEST NOTreadSalesOrdersUserCountryCode {
	DENY read ON salesOrders POLICY readSalesOrdersUserCountryCode INPUT 
	{
		$user: {
			country: 'DE'
		},
		CountryCode: 'FR'
	};
}
TEST writeSalesOrdersUserName {
	GRANT write ON salesOrders POLICY writeSalesOrdersUserName INPUT 
	{
		$user: {
			"name": 'Alice'
		},
		changedBy: 'Alice'
	};
}