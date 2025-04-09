POLICY PolBetween {
    GRANT IsBetween    ON * WHERE restr BETWEEN     1 AND 2;
    GRANT IsNotBetween ON * WHERE restr NOT BETWEEN 1 AND 2;
}

POLICY PolBetweenAnd {
    GRANT IsBetweenAndCondition ON * WHERE restr BETWEEN 1 AND 5 AND restr > 3;
}

POLICY PolBetweenOR {
    GRANT IsBetweenOrCondition  ON * WHERE restr BETWEEN 1 AND 5 OR restr > 6;
}
TEST BetweenCheck {
    GRANT IsBetween    POLICY PolBetween INPUT { restr : 1 }, { restr : 2 };
    DENY  IsBetween    POLICY PolBetween INPUT { restr : 0 }, { restr : 3 };
    DENY  IsBetween    POLICY PolBetween INPUT {                };
    GRANT IsBetween    POLICY PolBetween INPUT { restr : IGNORE };
    
    GRANT IsNotBetween POLICY PolBetween INPUT { restr : 0 }, { restr : 3 };
    DENY  IsNotBetween POLICY PolBetween INPUT { restr : 1 }, { restr : 2 };
    DENY  IsNotBetween POLICY PolBetween INPUT {                };
    GRANT IsNotBetween POLICY PolBetween INPUT { restr : IGNORE };
}

TEST PolBetweenAndTest{
    DENY  IsBetweenAndCondition POLICY PolBetweenAnd INPUT { restr : 0      };
    DENY  IsBetweenAndCondition POLICY PolBetweenAnd INPUT { restr : 2      };
    GRANT IsBetweenAndCondition POLICY PolBetweenAnd INPUT { restr : 4      };
    DENY  IsBetweenAndCondition POLICY PolBetweenAnd INPUT { restr : 6      };
    DENY  IsBetweenAndCondition POLICY PolBetweenAnd INPUT {                };
    GRANT IsBetweenAndCondition POLICY PolBetweenAnd INPUT { restr : IGNORE };
}

TEST PolBetweenOrTest{
    DENY  IsBetweenOrCondition POLICY PolBetweenOR INPUT { restr : 0      };
    GRANT IsBetweenOrCondition POLICY PolBetweenOR INPUT { restr : 2      };
    GRANT IsBetweenOrCondition POLICY PolBetweenOR INPUT { restr : 4      };
    GRANT IsBetweenOrCondition POLICY PolBetweenOR INPUT { restr : 7      }; 
    DENY  IsBetweenOrCondition POLICY PolBetweenOR INPUT {                }; 
    GRANT IsBetweenOrCondition POLICY PolBetweenOR INPUT { restr : IGNORE };
}