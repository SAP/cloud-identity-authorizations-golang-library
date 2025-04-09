default POLICY DefaultPolicy {
	GRANT default, test, DefaultAction ON DefaultResource;
}

TEST DefaultPolicyTest {
   GRANT default, test, DefaultAction ON DefaultResource;
   DENY  write, read                  ON DefaultResource;
   DENY  DefaultAction                ON SalesOrder;
}