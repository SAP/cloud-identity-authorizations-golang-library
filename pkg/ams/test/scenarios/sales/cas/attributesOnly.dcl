POLICY fullAuthorizedCountryCode {
	GRANT * ON * WHERE CountryCode IN ('FR','DE','GB');
}

/// Tests

TEST fullAuthorizedCountryCode {
	GRANT read ON salesOrders POLICY fullAuthorizedCountryCode INPUT 
	{ CountryCode: 'FR' };
}

TEST fullAuthorizedCountryCode_01 {
	GRANT read ON salesOrders POLICY fullAuthorizedCountryCode INPUT 
	{ CountryCode: 'DE' };
}

TEST fullAuthorizedCountryCode_02 {
	GRANT write ON salesOrders POLICY fullAuthorizedCountryCode INPUT 
	{
		CountryCode: 'GB',
		salesOrderId: 567 //should be ignored
	};
}

TEST NOTfullAuthorizedCountryCode {
	DENY read ON salesOrders POLICY fullAuthorizedCountryCode INPUT 
	{
		CountryCode: 'CH',
		salesOrderId: 345 //should be ignored
	};
}

TEST NOTfullAuthorizedCountryCode_01 {
	DENY read ON salesOrders POLICY fullAuthorizedCountryCode INPUT 
	{ CountryCode: 'CH' };
}