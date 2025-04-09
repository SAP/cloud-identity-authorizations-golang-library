POLICY salesOrderItemsCountryCode {
	GRANT * ON salesOrderItems WHERE CountryCode IN ('DE', 'CH', 'AT');
}

POLICY readSalesOrdersLike {
	GRANT read ON salesOrders WHERE Name LIKE '%deal%';
}



/// TESTS

TEST salesOrderItemsCountryCode {
	GRANT read ON salesOrderItems POLICY salesOrderItemsCountryCode INPUT 
	{ CountryCode: 'DE' };
}

TEST NOTsalesOrderItemsCountryCode {
	DENY read ON salesOrderItems POLICY salesOrderItemsCountryCode INPUT 
	{ CountryCode: 'US' };
}

TEST readSalesOrdersLike {
	GRANT read ON salesOrders POLICY readSalesOrdersLike INPUT 
	{ Name: 'Winterdeals' };
}

TEST NOTreadSalesOrdersLike {
	DENY read ON salesOrders POLICY readSalesOrdersLike INPUT 
	{ Name: 'WinterNews' };
}