POLICY MultipleActions {
    GRANT *      ON * WHERE restr = 0;
    GRANT M1     ON * WHERE restr = 1;
    GRANT M1, M2 ON * WHERE restr = 2;
}

TEST MultiActionsTest {
    GRANT *     POLICY MultipleActions INPUT { restr: 0      };
    GRANT M1,M2 POLICY MultipleActions INPUT { restr: 0      };
    DENY  *     POLICY MultipleActions INPUT {               };
    GRANT *     POLICY MultipleActions INPUT { restr: IGNORE };
    EXPECT restr = 0 FOR * ON * POLICY MultipleActions INPUT { restr: UNKNOWN };

    DENY  *     POLICY MultipleActions INPUT { restr: 1 };
    GRANT M1    POLICY MultipleActions INPUT { restr: 1 };
    DENY  M2    POLICY MultipleActions INPUT { restr: 1 };
    DENY  M1    POLICY MultipleActions INPUT {               };
    GRANT M1    POLICY MultipleActions INPUT { restr: IGNORE };
    EXPECT restr = 0 OR restr = 1 OR restr = 2 FOR M1 ON * POLICY MultipleActions INPUT { restr: UNKNOWN };


    DENY  *     POLICY MultipleActions INPUT { restr: 2 };
    GRANT M1,M2 POLICY MultipleActions INPUT { restr: 2 };
    DENY  M1,M2 POLICY MultipleActions INPUT {               };
    GRANT M1,M2 POLICY MultipleActions INPUT { restr: IGNORE };
    EXPECT restr = 0 OR restr = 2              FOR M2 ON * POLICY MultipleActions INPUT { restr: UNKNOWN };
}