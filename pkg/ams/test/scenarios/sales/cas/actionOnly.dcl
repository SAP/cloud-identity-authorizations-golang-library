POLICY readWriteAll {
	GRANT read, write ON *;
}

POLICY activateAll {
	GRANT activate ON *;
}


/// Tests

TEST readwriteAll_read {
	GRANT read ON * POLICY readWriteAll;
}

TEST readwriteAll_write {
	GRANT "write" ON * POLICY readWriteAll;
}

TEST readwriteAll_read_write{
	GRANT read, write ON * POLICY readWriteAll;
}

TEST NOTreadWriteAll_activate { 
	DENY "activate" ON * POLICY readWriteAll;
}

TEST activateAll {
	GRANT activate ON * POLICY activateAll;
}

TEST NOTactivateAll {
	DENY read ON * POLICY activateAll;
}

TEST readSalesOrdersSalesOrderItems {
	GRANT read ON salesOrderItemsLists POLICY readSalesOrdersSalesOrderItems;
}

TEST NOTreadSalesOrdersSalesOrderItems {
	DENY delete ON salesOrderItemsLists POLICY readSalesOrdersSalesOrderItems;
}