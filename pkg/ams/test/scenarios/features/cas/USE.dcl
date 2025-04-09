POLICY SimpleRestrictPolicy {
    GRANT A ON * WHERE num1 IS RESTRICTED;
}

POLICY SimpleRestrictPolicyString {
    GRANT A ON * WHERE stringval IS RESTRICTED;
}

POLICY SimpleRestrictPolicyBoolean {
    GRANT A ON * WHERE bool1 IS RESTRICTED;
}

POLICY SimpleRestrictPolicyAbsenceRestrict {
    USE SimpleRestrictPolicy;
}

TEST SimpleRestrictPolicyAbsenceRestrictTest {
    DENY A ON * POLICY SimpleRestrictPolicyAbsenceRestrict;
}

POLICY SimpleRestrictPolicyWithIn {
    USE SimpleRestrictPolicy RESTRICT num1 IN(1,5);
}

TEST SimpleRestrictPolicyWithInTest {
    GRANT A ON * POLICY SimpleRestrictPolicyWithIn INPUT { num1 : 1      };
    DENY  A ON * POLICY SimpleRestrictPolicyWithIn INPUT { num1 : 2      };
    DENY  A ON * POLICY SimpleRestrictPolicyWithIn INPUT {               };
    GRANT A ON * POLICY SimpleRestrictPolicyWithIn INPUT { num1 : IGNORE };
}

POLICY SimpleRestrictPolicyWithBetween {
    USE SimpleRestrictPolicy RESTRICT num1 BETWEEN 1 AND 5;
}

TEST SimpleRestrictPolicyWithBetweenTest {
    GRANT A ON * POLICY SimpleRestrictPolicyWithBetween INPUT { num1 : 2      };
    DENY  A ON * POLICY SimpleRestrictPolicyWithBetween INPUT { num1 : 6      };
    DENY  A ON * POLICY SimpleRestrictPolicyWithBetween INPUT {               };
    GRANT A ON * POLICY SimpleRestrictPolicyWithBetween INPUT { num1 : IGNORE };
}

POLICY SimpleRestrictPolicyWithLike {
    USE SimpleRestrictPolicyString RESTRICT stringval LIKE 'a%';
}

TEST SimpleRestrictPolicyWithLikeTest {
    GRANT A ON * POLICY SimpleRestrictPolicyWithLike INPUT { stringval : 'ab'   };
    DENY  A ON * POLICY SimpleRestrictPolicyWithLike INPUT { stringval : 'bb'   };
    DENY  A ON * POLICY SimpleRestrictPolicyWithLike INPUT {                    };
    GRANT A ON * POLICY SimpleRestrictPolicyWithLike INPUT { stringval : IGNORE };
}

POLICY SimpleRestrictPolicyWithOperator {
    USE SimpleRestrictPolicy        RESTRICT num1 > 1;
    USE SimpleRestrictPolicyString  RESTRICT stringval < 'B';
    USE SimpleRestrictPolicyBoolean RESTRICT bool1 = TRUE;
}

TEST SimpleRestrictPolicyWithOperatorTest {
    GRANT A ON * POLICY SimpleRestrictPolicyWithOperator               INPUT { num1 : 2       };
    DENY  A ON * POLICY SimpleRestrictPolicyWithOperator               INPUT { num1 : 0       };
    DENY  A ON * POLICY SimpleRestrictPolicyWithOperator               INPUT {                };
    GRANT A ON * POLICY SimpleRestrictPolicyWithOperator               INPUT { num1 : IGNORE  };
    EXPECT num1 > 1 FOR A ON * POLICY SimpleRestrictPolicyWithOperator INPUT { num1 : UNKNOWN };

    GRANT A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                      { stringval : 'A'     };
    DENY  A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                      { stringval : 'C'     };
    DENY  A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                      {                     };
    GRANT A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                      { stringval : IGNORE  };
    EXPECT stringval < 'B' FOR A ON * POLICY SimpleRestrictPolicyWithOperator INPUT { stringval : UNKNOWN };

    GRANT A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                   { bool1 : TRUE    };
    DENY  A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                   { bool1 : FALSE   };
    DENY  A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                   {                 };
    GRANT A ON * POLICY SimpleRestrictPolicyWithOperator INPUT                   { bool1 : IGNORE  };
    EXPECT bool1 = TRUE FOR A ON * POLICY SimpleRestrictPolicyWithOperator INPUT { bool1 : UNKNOWN };
}

POLICY SimpleNotRestrictPolicy {
    GRANT A ON * WHERE num1 IS NOT RESTRICTED;
}

POLICY SimpleNotRestrictPolicyString {
    GRANT A ON * WHERE stringval IS NOT RESTRICTED;
}

POLICY SimpleNotRestrictPolicyBoolean {
    GRANT A ON * WHERE bool1 IS NOT RESTRICTED;
}

POLICY SimpleNotRestrictPolicyWithIn {
    USE SimpleNotRestrictPolicy RESTRICT num1 IN(1,5);
}

TEST SimpleNotRestrictPolicyWithInTest {
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithIn INPUT { num1 : 1      };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithIn INPUT { num1 : 2      };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithIn INPUT {               };
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithIn INPUT { num1 : IGNORE };
}

POLICY SimpleNotRestrictPolicyAbsenceRestrict {
    USE SimpleNotRestrictPolicy;
}

TEST SimpleNotRestrictPolicyAbsenceRestrictTest {
    GRANT A ON * POLICY SimpleNotRestrictPolicyAbsenceRestrict;
}

POLICY SimpleNotRestrictPolicyWithBetween {
    USE SimpleNotRestrictPolicy RESTRICT num1 BETWEEN 1 AND 5;
}

TEST SimpleNotRestrictPolicyWithBetweenTest {
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithBetween INPUT { num1 : 2      };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithBetween INPUT { num1 : 6      };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithBetween INPUT {               };
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithBetween INPUT { num1 : IGNORE };
}

POLICY SimpleNotRestrictPolicyWithLike {
    USE SimpleNotRestrictPolicyString RESTRICT stringval LIKE 'a%';
}

TEST SimpleNotRestrictPolicyWithLikeTest {
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithLike INPUT { stringval : 'ab'   };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithLike INPUT { stringval : 'bb'   };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithLike INPUT {                    };
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithLike INPUT { stringval : IGNORE };
}

POLICY SimpleNotRestrictPolicyWithOperator {
    USE SimpleNotRestrictPolicy        RESTRICT num1 > 1;
    USE SimpleNotRestrictPolicyString  RESTRICT stringval < 'B';
    USE SimpleNotRestrictPolicyBoolean RESTRICT bool1 = TRUE;
}

TEST SimpleNotRestrictPolicyWithOperatorTest {
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT               { num1 : 2       };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT               { num1 : 0       };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT               {                };
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT               { num1 : IGNORE  };
    EXPECT num1 > 1 FOR A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT { num1 : UNKNOWN };

    GRANT A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT                      { stringval : 'A'     };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT                      { stringval : 'C'     };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT                      {                     };
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT                      { stringval : IGNORE  };
    EXPECT stringval < 'B' FOR A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT { stringval : UNKNOWN };

    GRANT A ON * POLICY SimpleNotRestrictPolicyWithOperator                   INPUT { bool1 : TRUE    };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithOperator                   INPUT { bool1 : FALSE   };
    DENY  A ON * POLICY SimpleNotRestrictPolicyWithOperator                   INPUT {                 };
    GRANT A ON * POLICY SimpleNotRestrictPolicyWithOperator                   INPUT { bool1 : IGNORE  };
    EXPECT bool1 = TRUE FOR A ON * POLICY SimpleNotRestrictPolicyWithOperator INPUT { bool1 : UNKNOWN };
}

POLICY Empty {
}

POLICY RestrBase {
    USE Empty;
    USE Empty;

    GRANT X on *;
    GRANT Restr    ON * WHERE restr is RESTRICTED;
    GRANT RestrNot ON * WHERE restr is NOT RESTRICTED;
}

POLICY RestrUse_Restricted {
    Use RestrBase RESTRICT restr IS RESTRICTED;
}

POLICY RestrUse_NOT_Restricted {
    Use RestrBase RESTRICT restr IS NOT RESTRICTED;
}

POLICY RestrBaseIndirect {
    USE RestrBase;
}

POLICY RestrUse_NOT_Restricted_And_Unrestricted {
    use RestrBaseIndirect;
    Use RestrBase RESTRICT restr IS NOT RESTRICTED;
}

POLICY RestrUse_Restricted_And_Unrestricted {
    use RestrBaseIndirect;
    Use RestrBase RESTRICT restr IS RESTRICTED;
}

POLICY BaseAllUnrestricted {
    GRANT AllUR ON * WHERE ur1 IS RESTRICTED 
                       AND ur2 IS RESTRICTED 
                       AND ur3 IS RESTRICTED 
                       AND ur4 IS RESTRICTED 
                       AND ur5 IS RESTRICTED;
}

POLICY RestrictBaseAllUnrestricted {
    USE BaseAllUnrestricted RESTRICT ur1 LIKE 'a', 
                                     ur2 in ('a', 'b'),
                                     ur3 BETWEEN 'a' AND 'c',
                                     ur4 = 'a' ,
                                     ur5 is NOT NULL;
}

TEST RestrictionCheck {
    GRANT Restr, RestrNot POLICY RestrUse_NOT_Restricted;
    DENY  Restr, RestrNot POLICY RestrUse_Restricted;
    GRANT RestrNot POLICY RestrUse_Restricted_And_Unrestricted;
    DENY  Restr    POLICY RestrUse_Restricted_And_Unrestricted;
    GRANT Restr, RestrNot POLICY RestrUse_NOT_Restricted_And_Unrestricted;
}

TEST RestrictBaseAllUnrestricted{
    GRANT AllUR POLICY RestrictBaseAllUnrestricted INPUT {
        ur1: 'a',
        ur2: 'a',
        ur3: 'a',
        ur4: 'a',
        ur5: 'a'
    };
    DENY AllUR POLICY RestrictBaseAllUnrestricted INPUT {
        ur1: 'a',
        ur2: 'a',
        ur3: 'a',
        ur4: 'a',
        ur5: null
    };
    GRANT AllUR POLICY RestrictBaseAllUnrestricted INPUT {
        ur1: 'a',
        ur2: 'b',
        ur3: 'b',
        ur4: 'a',
        ur5: 'b'
    };
    DENY AllUR POLICY RestrictBaseAllUnrestricted INPUT {
        ur1: 'b',
        ur2: 'a',
        ur3: 'a',
        ur4: 'a',
        ur5: 'b'
    };
}
