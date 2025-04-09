Policy PolORTrue {
    GRANT IsOR ON * WHERE bool1 = true OR  bool2 = true;
}

Test ORCheckTrue {
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: true,    bool2: true    };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: true,    bool2: false   };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: true                    };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: true,    bool2: IGNORE  };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: true,    bool2: UNKNOWN };

    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: false,   bool2: true    };
    DENY                    IsOR POLICY PolORTrue INPUT { bool1: false,   bool2: false   };
    DENY                    IsOR POLICY PolORTrue INPUT { bool1: false                   };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: false,   bool2: IGNORE  };
    EXPECT bool2 = true FOR IsOR POLICY PolORTrue INPUT { bool1: false,   bool2: UNKNOWN };

    GRANT                   IsOR POLICY PolORTrue INPUT {                 bool2: true    };
    DENY                    IsOR POLICY PolORTrue INPUT {                 bool2: false   };
    DENY                    IsOR POLICY PolORTrue INPUT {                                };
    GRANT                   IsOR POLICY PolORTrue INPUT {                 bool2: IGNORE  };
    EXPECT bool2 = true FOR IsOR POLICY PolORTrue INPUT {                 bool2: UNKNOWN };

    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: IGNORE,  bool2: true    };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: IGNORE,  bool2: false   };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: IGNORE                  };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: IGNORE,  bool2: IGNORE  };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: IGNORE,  bool2: UNKNOWN };

    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: UNKNOWN, bool2: true    };
    EXPECT bool1 = true FOR IsOR POLICY PolORTrue INPUT { bool1: UNKNOWN, bool2: false   };
    EXPECT bool1 = true FOR IsOR POLICY PolORTrue INPUT { bool1: UNKNOWN                 };
    GRANT                   IsOR POLICY PolORTrue INPUT { bool1: UNKNOWN, bool2: IGNORE  };
    // By default (and with the @comparison: 'default' annotation), checks don't consider the ordering of the following operators: AND, OR, = and <>.
    // Also, attributes are expected to be sorted lexicographically.
    //
    // For order preserving checks, set @comparison: 'strict'.
    EXPECT bool2 = true OR
           bool1 = true FOR IsOR POLICY PolORTrue INPUT { bool1: UNKNOWN, bool2: UNKNOWN };
}

/* 
 * Revised definition for IGNORE. See ../documentation/concept/Ignore.md
 */
POLICY ORCornerCases {
    // No constantly neutral element
    GRANT CC1 ON * WHERE bool1 = true;
    
    // One constantly neutral element, one "OR true"
    GRANT CC2a ON * WHERE 0 = 1 OR bool1 = true;
    GRANT CC2b ON * WHERE (bool2 IS NULL AND bool2 IS NOT NULL) OR bool1 = true;
    GRANT CC2c ON * WHERE (bool2 IS NULL AND bool2 = true) OR bool1 = true;
    
    // One constantly neutral element, two "OR true"
    GRANT CC3 ON * WHERE 0 = 1 OR bool1 = true OR bool2 = true;
    
    // Four "OR true"
    GRANT CC4 ON * WHERE bool1 = true OR bool2 = true OR bool3 = true OR bool4 = true;
    
}

Test ORCornerCases {
    GRANT                   CC1 POLICY ORCornerCases INPUT { bool1: IGNORE }; 
    EXPECT bool1 = true FOR CC1 POLICY ORCornerCases INPUT { bool1: UNKNOWN }; 

    GRANT                   CC2a,CC2b,CC2c POLICY ORCornerCases INPUT { bool1: IGNORE }; 
    EXPECT bool1 = true FOR CC2a,CC2b,CC2c POLICY ORCornerCases INPUT { bool1: UNKNOWN }; 
    
    GRANT                   CC3 POLICY ORCornerCases INPUT { bool1: IGNORE, bool2: IGNORE }; 
    GRANT                   CC3 POLICY ORCornerCases INPUT { bool1: IGNORE, bool2: false }; 
    DENY                    CC3 POLICY ORCornerCases INPUT { bool1: false,  bool2: false }; 
    GRANT                   CC3 POLICY ORCornerCases INPUT { bool1: IGNORE, bool2: UNKNOWN }; 
    EXPECT bool2 = true FOR CC3 POLICY ORCornerCases INPUT { bool1: false,  bool2: UNKNOWN }; 

    GRANT                   CC4 POLICY ORCornerCases INPUT { bool1: IGNORE, bool2: IGNORE, bool3: IGNORE,  bool4: IGNORE }; 
    DENY                    CC4 POLICY ORCornerCases INPUT { bool1: false,  bool2: false,  bool3: false,   bool4: false }; 
    GRANT                   CC4 POLICY ORCornerCases INPUT { bool1: IGNORE, bool2: IGNORE, bool3: UNKNOWN, bool4: UNKNOWN }; 
}

///
Policy PolORFalse {
    GRANT IsOR ON * WHERE bool1 = false OR  bool2 = false;
}

Test ORCheckFalse {
    GRANT IsOR POLICY PolORFalse INPUT { bool1: false, bool2: false  };
    GRANT IsOR POLICY PolORFalse INPUT { bool1: false, bool2: true   };
    GRANT IsOR POLICY PolORFalse INPUT { bool1: false                };
    GRANT IsOR POLICY PolORFalse INPUT { bool1: true,  bool2: false  };
    DENY  IsOR POLICY PolORFalse INPUT { bool1: true,  bool2: true   };
    DENY  IsOR POLICY PolORFalse INPUT { bool1: true                 };
    GRANT IsOR POLICY PolORFalse INPUT {               bool2: false  };
    DENY  IsOR POLICY PolORFalse INPUT {               bool2: true   };
    DENY  IsOR POLICY PolORFalse INPUT {                             };
}
