POLICY readSalesOrdersCountryCode {
	GRANT read ON salesOrders WHERE CountryCode IN ('IT','AT', 'BE');
}

POLICY readwriteSalesOrdersCountryCode {
	GRANT read, write ON salesOrders WHERE CountryCode IN ('IT','AT');
}

/// TESTS

TEST readSalesOrdersCountryCode {
	GRANT read ON salesOrders POLICY readSalesOrdersCountryCode INPUT 
	{ CountryCode: 'IT' };
}

TEST NOTreadSalesOrdersCountryCode {
	DENY write ON salesOrders POLICY readSalesOrdersCountryCode INPUT 
	{ CountryCode: 'IT' };
}

TEST readwriteSalesOrdersCountryCode {
	GRANT read, write ON salesOrders POLICY readwriteSalesOrdersCountryCode INPUT 
	{ CountryCode: 'IT' };
}

TEST readwriteSalesOrdersCountryCode_01 {
	GRANT read ON salesOrders POLICY readwriteSalesOrdersCountryCode INPUT 
	{ CountryCode: 'IT' };
}

TEST readwriteSalesOrdersCountryCode_02{
	GRANT write ON salesOrders POLICY readwriteSalesOrdersCountryCode INPUT 
	{ CountryCode: 'AT' };
}

TEST NOTreadwriteSalesOrdersCountryCode {
	DENY read, write ON salesOrders POLICY readwriteSalesOrdersCountryCode INPUT 
	{ CountryCode: 'US' };
}