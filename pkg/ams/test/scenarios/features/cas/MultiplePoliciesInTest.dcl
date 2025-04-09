POLICY PolA {
    GRANT A ON *;
}

POLICY PolB {
    GRANT B ON *;
}

TEST MultiplePoliciesInTest {
    GRANT A POLICY PolA;
    GRANT B POLICY PolB;
    DENY  B ON * POLICY PolA;
    DENY  A ON * POLICY PolB;
    GRANT A,B POLICY PolA, PolB;
}
