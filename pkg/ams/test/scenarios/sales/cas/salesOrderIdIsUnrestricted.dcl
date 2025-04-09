POLICY readSalesOrderIdIsUnrestricted {
	GRANT read ON salesOrders WHERE salesOrderId IS NOT RESTRICTED;
}

POLICY UseReadSalesOrderIsUnrestrictedBetween {
	USE readSalesOrderIdIsUnrestricted RESTRICT salesOrderId BETWEEN 100 AND 300;
}

POLICY UseReadSalesOrderIsUnrestrictedEqual {
	USE readSalesOrderIdIsUnrestricted RESTRICT salesOrderId = 3456
									 RESTRICT salesOrderId = 5678;
}

POLICY UseReadSalesOrderIsUnrestrictedLt { 
	USE readSalesOrderIdIsUnrestricted RESTRICT salesOrderId < 5;
}

POLICY UseReadSalesOrderIsUnrestrictedGt {
	use readSalesOrderIdIsUnrestricted RESTRICT salesOrderId > 7;
} 

POLICY writeSalesOrderCountryCodeIsUnrestricted {
	GRANT write ON salesOrders WHERE salesOrderId IS NOT RESTRICTED AND CountryCode IS NOT RESTRICTED;
}

@description: 'use a policy with 2 attributes'
POLICY UseWriteSalesOrderCountryCodeIsUnrestricted {
	USE writeSalesOrderCountryCodeIsUnrestricted 
		RESTRICT salesOrderId = 5498 , CountryCode = 'DE'
		RESTRICT salesOrderId BETWEEN 1000 AND 3000, CountryCode = $user.country;
}


/// TESTS

TEST readSalesOrderIdIsUnrestricted {
	GRANT read ON salesOrders POLICY readSalesOrderIdIsUnrestricted INPUT
	{ salesOrderId: 354689 },
	{ salesOrderId: 32435 };
}

TEST NOTreadSalesOrderIdIsUnrestricted {
	DENY delete ON salesOrders POLICY readSalesOrderIdIsUnrestricted INPUT 
	{ salesOrderId: 354689 },
	{ salesOrderId: 32435 };
}

TEST UseReadSalesOrderIsUnrestrictedBetween {
	GRANT read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedBetween INPUT 
	{ salesOrderId: 112 },
	{ salesOrderId: 299 };
}

TEST NOTUseReadSalesOrderIsUnrestrictedBetween {
	DENY read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedBetween INPUT 
	{ salesOrderId: 301 },
	{ salesOrderId: 99 };
}

TEST NOTUseReadSalesOrderIsUnrestrictedBetween_01 {
	DENY write ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedBetween INPUT 
	{ salesOrderId: 112 },
	{ salesOrderId: 299 };
}

TEST UseWriteSalesOrderCountryCodeIsUnrestricted {
	GRANT write ON salesOrders POLICY UseWriteSalesOrderCountryCodeIsUnrestricted INPUT 
	{
		salesOrderId: 5498,
		CountryCode: 'DE'
	};
}

TEST NOTUseWriteSalesOrderCountryCodeIsUnrestricted {
	DENY write ON salesOrders POLICY UseWriteSalesOrderCountryCodeIsUnrestricted INPUT 
	{
		salesOrderId: 6000,
		CountryCode: 'DE'
	},{
		salesOrderId: 4598,
		CountryCode: 'FR'
	};
}

TEST UseReadSalesOrderIsUnrestrictedEqual {
	GRANT read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedEqual INPUT 
	{ salesOrderId: 3456 },
	{ salesOrderId: 5678 };
}
	
TEST NOTUseReadSalesOrderIsUnrestrictedEqual {
	DENY read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedEqual INPUT 
	{ salesOrderId: 345678 };
}

TEST UseReadSalesOrderIsUnrestrictedLt {
	GRANT read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedLt INPUT 
	{ salesOrderId: 3 };
}

TEST NOTUseReadSalesOrderIsUnrestrictedLt {
	DENY read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedLt INPUT
    { salesOrderId: 88 };
}

// -------------------------------------------	
// UseReadSalesOrderIsUnrestrictedGt
// -------------------------------------------	
TEST UseReadSalesOrderIsUnrestrictedGt {
	GRANT read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedGt INPUT 
	{ salesOrderId: 9 };	
}

TEST NOTUseReadSalesOrderIsUnrestrictedGt {
	DENY read ON salesOrders POLICY UseReadSalesOrderIsUnrestrictedGt INPUT 
	{ salesOrderId: -1 };	
}

// -------------------------------------------
// writeSalesOrderCountryCodeIsUnrestricted
// -------------------------------------------
TEST writeSalesOrderCountryCodeIsUnrestricted {
	GRANT write ON salesOrders POLICY writeSalesOrderCountryCodeIsUnrestricted INPUT 
	{
		salesOrderId: 435,
		CountryCode: 'FR'
	};
}

TEST writeSalesOrderCountryCodeIsUnrestricted_01 {
	GRANT write ON salesOrders POLICY writeSalesOrderCountryCodeIsUnrestricted;
}

TEST NOTwriteSalesOrderCountryCodeIsUnrestricted {
	DENY read ON salesOrders POLICY writeSalesOrderCountryCodeIsUnrestricted INPUT 
	{
		salesOrderId: 435,
		CountryCode: 'FR'
	};
	DENY delete ON salesOrders POLICY writeSalesOrderCountryCodeIsUnrestricted INPUT 
	{
		salesOrderId: 435,
		CountryCode: 'FR'
	};
}
