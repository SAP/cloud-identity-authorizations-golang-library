POLICY MultiResources {
    GRANT A ON * WHERE restr = 0;
    GRANT A ON R1 WHERE restr = 1;
    GRANT A ON R1,R2 WHERE restr = 2;
}

TEST MultiResources {
    GRANT A ON *     POLICY MultiResources INPUT { restr: 0 };
    GRANT A ON R1,R2 POLICY MultiResources INPUT { restr: 0 };
    DENY  B ON *     POLICY MultiResources INPUT { restr: 0 };
    DENY  A ON *     POLICY MultiResources INPUT {          };
    GRANT A ON *     POLICY MultiResources INPUT { restr: IGNORE };
    EXPECT restr = 0 FOR A ON * POLICY MultiResources INPUT { restr:UNKNOWN };

    DENY  A ON *     POLICY MultiResources INPUT { restr: 1 };
    GRANT A ON R1    POLICY MultiResources INPUT { restr: 1 };
    GRANT A ON R1    POLICY MultiResources INPUT { restr: 0 };
    DENY  A ON R2    POLICY MultiResources INPUT { restr: 1 };
    DENY  A ON R1    POLICY MultiResources INPUT {          };
    GRANT A ON R1    POLICY MultiResources INPUT { restr: IGNORE };
    EXPECT restr = 1 OR restr = 0 OR restr = 2 FOR A ON R1 POLICY MultiResources INPUT { restr:UNKNOWN };

    DENY  A ON *     POLICY MultiResources INPUT { restr: 2 };
    GRANT A ON R1,R2 POLICY MultiResources INPUT { restr: 2 };
    DENY  A ON R1,R2 POLICY MultiResources INPUT {          };
    GRANT A ON R1,R2 POLICY MultiResources INPUT { restr: IGNORE };
    EXPECT restr = 0 OR restr = 2 FOR A ON R2 POLICY MultiResources INPUT { restr:UNKNOWN };
}