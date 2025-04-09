POLICY salesOrderItemsComlexRule {
	GRANT read, write ON salesOrderItems WHERE CountryCode = 'GB' AND 
											   (salesOrderId = 233 OR salesOrderId = 677);
}

POLICY salesOrderItemsComlexRuleLike {
	GRANT read, write ON salesOrders, saledOrderItems WHERE (CountryCode ='DE' AND 
        ( salesOrderId = 200 OR Name LIKE '%Winter%' )) OR
        ( CountryCode = 'EN' AND (salesOrderId = 300 OR Name LIKE '%SOMMER!%..%' ESCAPE '!'));
}

POLICY readSalesOrdersNotNull {
	GRANT read ON salesOrderLists WHERE amount is not null;
}

/// TESTS

TEST salesOrderItemsComlexRule {
	GRANT read ON salesOrderItems POLICY salesOrderItemsComlexRule INPUT {
		CountryCode: 'GB',
		salesOrderId: 233
	},{		
		CountryCode: 'GB',
		salesOrderId: 677
	};
}

TEST NOTsalesOrderItemsComlexRule {
	DENY read ON salesOrderItems POLICY salesOrderItemsComlexRule INPUT {
		CountryCode: 'IT',
		salesOrderId: 677
	},{
		CountryCode: 'DE',
		salesOrderId: 233
		
	};
}

TEST salesOrderItemsComlexRuleLike {
	GRANT read, write ON salesOrders POLICY salesOrderItemsComlexRuleLike INPUT 
	{
		Name: 'SommerSales',
		salesOrderId: 300,
		CountryCode: 'EN'
	};
}
	
TEST NOTsalesOrderItemsComlexRuleLike {
	DENY read ON salesOrders POLICY salesOrderItemsComlexRuleLike INPUT 
	{
		Name: 'SommerSales',
		salesOrderId: 200,
		CountryCode: 'EN'
	};
}
	
TEST NOTsalesOrderItemsComlexRuleLike_01 {
	DENY read ON salesOrders POLICY salesOrderItemsComlexRuleLike INPUT 
	{
		Name: 'WinterSales',
		salesOrderId: 200,
		CountryCode: 'EN'
	};	
}
	
TEST salesOrderItemsComlexRuleLike_01 {
	GRANT read ON salesOrders POLICY salesOrderItemsComlexRuleLike INPUT 
	{
		Name: 'SoMMerSales',
		salesOrderId: 300,
		CountryCode: 'EN'
	};	
}