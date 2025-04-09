POLICY stackPolicies {
	use readSalesOrdersCountryCode;
	use readSalesOrdersSalesOrderItems;
}

@plan: 'standard'
POLICY useReadSalesOrdersCountryCodeExtraGrant {
	USE readSalesOrdersCountryCode;
	GRANT write ON salesOrders WHERE CountryCode = 'SE';
}


/// TESTS


TEST useReadSalesOrdersCountryCodeExtraGrant {
	GRANT write ON salesOrders POLICY useReadSalesOrdersCountryCodeExtraGrant INPUT 
	{ CountryCode: 'SE' };
}
	
TEST NOTuseReadSalesOrdersCountryCodeExtraGrant {
	DENY write ON salesOrders POLICY useReadSalesOrdersCountryCodeExtraGrant INPUT 
	{ CountryCode: 'FR' };
}