POLICY readAllCountryCode {
	GRANT read ON * WHERE CountryCode = 'DE' OR CountryCode = 'FR';
}



@description: 'not equal'
POLICY readAllCountryCodeNe {
	GRANT read ON * WHERE CountryCode <> 'BE';
	GRANT read ON * WHERE CountryCode =  'BE';
}

TEST NOTreadAllCountryCodeNe {
	DENY read ON * POLICY readAllCountryCodeNe;
}




POLICY readSalesOrdersCountryCodeSalesIdBetween {
	GRANT read ON * WHERE CountryCode = 'DE' AND salesOrderId BETWEEN 100 AND 200;
}

///


TEST readAllCountryCode {
	GRANT read ON salesOrders POLICY readAllCountryCode INPUT
	{ CountryCode: 'DE' };
}

TEST readAllCountryCode_01 {
	GRANT read ON salesOrders POLICY readAllCountryCode INPUT
	{ CountryCode: 'FR' };
}

TEST NOTreadAllCountryCode {
	DENY read ON salesOrders POLICY readAllCountryCode INPUT
	{ CountryCode: 'US' };
}

//TEST readAllCountryCodeNe {
//	GRANT read ON * POLICY readAllCountryCodeNe INPUT 
//	{ CountryCode: 'DE' };
//}


TEST readSalesOrdersCountryCodeSalesIdBetween {
	GRANT read ON salesOrders POLICY readSalesOrdersCountryCodeSalesIdBetween INPUT 
	{
		CountryCode: 'DE', 
	  	salesOrderId: 133
	};
}

TEST NOTreadSalesOrdersCountryCodeSalesIdBetween {
	DENY read ON salesOrders POLICY readSalesOrdersCountryCodeSalesIdBetween INPUT 
	{ 
		CountryCode: 'DE',
	  	salesOrderId: 233
	};
}
