POLICY salesOrderItems2Rules {
	GRANT read ON salesOrderItems WHERE CountryCode = 'FR' AND salesOrderId = 200;
	GRANT write ON salesOrderItems WHERE CountryCode = 'DE' AND salesOrderId = 300;
}

// TESTS

TEST salesOrderItems2Rules {
	GRANT read ON salesOrderItems POLICY salesOrderItems2Rules INPUT 
	{
		CountryCode: 'FR',
		salesOrderId: 200
	};
	GRANT write ON salesOrderItems POLICY salesOrderItems2Rules INPUT 
	{
		CountryCode: 'DE',
		salesOrderId: 300
	};
}

TEST NOTsalesOrderItems2Rules {
	DENY read ON salesOrders POLICY salesOrderItems2Rules INPUT 
	{
		CountryCode: 'IT',
		salesOrderId: 300
	},{
		CountryCode: 'US'
	};
}
