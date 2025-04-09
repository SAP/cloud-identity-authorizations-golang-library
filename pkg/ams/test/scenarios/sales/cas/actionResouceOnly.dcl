POLICY readSalesOrders {
	GRANT read ON salesOrders;
}

POLICY "writeDeleteSalesOrderItems" {
	GRANT write, "delete" ON "salesOrderItems";
}

POLICY readSalesOrdersSalesOrderItems {
	GRANT read ON salesOrders, salesOrderItems, salesOrderItemsLists;
}


/// Tests

TEST readSalesOrders {
	GRANT read ON salesOrders POLICY readSalesOrders;
}

TEST NOTreadSalesOrders {
	DENY write ON salesOrders POLICY readSalesOrders;
}

TEST "writeDeleteSalesOrderItems" {
	GRANT write, delete ON salesOrderItems POLICY writeDeleteSalesOrderItems;
}

TEST "NOTwriteDeleteSalesOrderItems" {
	DENY read ON salesOrderItems POLICY writeDeleteSalesOrderItems;
}