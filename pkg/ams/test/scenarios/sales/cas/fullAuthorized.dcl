@description: 'fully authorized'
POLICY fullAuthorized {
	GRANT * ON *;
}

/// TESTS

TEST fullAuthorized_read {
	GRANT read ON salesOrders POLICY fullAuthorized;
}

TEST fullAuthorized_write {
	GRANT write ON salesOrders POLICY fullAuthorized;
}

TEST fullAuthorized_activate {
	GRANT activate ON salesOrders POLICY fullAuthorized;
}

TEST fullAuthorized_activate_01 {
	GRANT activate ON salesOrdersItems POLICY fullAuthorized;
}