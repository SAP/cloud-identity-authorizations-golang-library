POLICY salesOrders {
	GRANT * ON salesOrders;
}

POLICY salesOrdersSalesOrderItems {
	GRANT * ON salesOrders, salesOrderItems, salesOrderLists;
}

/// TESTS

TEST "salesOrders" {
	GRANT read ON salesOrders POLICY salesOrders;
}

TEST "NOTsalesOrders" {
	DENY read ON salesOrderItems POLICY salesOrders;
}

TEST salesOrdersSalesOrderItems {
	GRANT read ON salesOrderLists POLICY salesOrdersSalesOrderItems;
}

TEST salesOrdersSalesOrderItems_01 {
	DENY read ON salesOrderItemsLists POLICY salesOrdersSalesOrderItems;	
}