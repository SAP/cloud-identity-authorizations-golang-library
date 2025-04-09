POLICY PolInLiteral {
    GRANT IsIn         ON * WHERE restr     IN (1,2);
    GRANT IsNotIn      ON * WHERE restr NOT IN (1,2);
    GRANT IsEmptyIn    ON * WHERE restr     IN ();
    GRANT IsEmptyNotIn ON * WHERE restr NOT IN ();
}

TEST CheckInLiteral {
    GRANT IsIn         POLICY PolInLiteral INPUT { restr : 1 }, { restr : 2 };
    DENY  IsIn         POLICY PolInLiteral INPUT { restr : 0 }, { restr : 3 };
    DENY  IsIn         POLICY PolInLiteral INPUT {                };
    //GRANT IsIn         POLICY PolInLiteral INPUT { restr : IGNORE };
    EXPECT restr = 1 OR restr = 2 
      FOR IsIn         POLICY PolInLiteral INPUT { restr : UNKNOWN };

    GRANT IsNotIn      POLICY PolInLiteral INPUT { restr : 0 }, { restr : 3 };
    DENY  IsNotIn      POLICY PolInLiteral INPUT { restr : 1 }, { restr : 2 };
    DENY  IsNotIn      POLICY PolInLiteral INPUT {                };
    //GRANT IsNotIn      POLICY PolInLiteral INPUT { restr : IGNORE };
    EXPECT restr <> 1 AND restr <> 2 
      FOR IsNotIn      POLICY PolInLiteral INPUT { restr : UNKNOWN };

    DENY  IsEmptyIn    POLICY PolInLiteral INPUT { restr : 1      };
    DENY  IsEmptyIn    POLICY PolInLiteral INPUT {                };
    //GRANT IsEmptyIn    POLICY PolInLiteral INPUT { restr : IGNORE };
    DENY  IsEmptyIn    POLICY PolInLiteral INPUT { restr : UNKNOWN };

    //GRANT IsEmptyNotIn POLICY PolInLiteral INPUT { restr : 1      };
    //GRANT IsEmptyNotIn POLICY PolInLiteral INPUT {                };
    //GRANT IsEmptyNotIn POLICY PolInLiteral INPUT { restr : IGNORE };
    //GRANT  IsEmptyNotIn POLICY PolInLiteral INPUT { restr : UNKNOWN };
}

POLICY PolInDynamic {
    GRANT IsIn         ON num WHERE restr     IN numArray;
    GRANT IsNotIn      ON num WHERE restr NOT IN numArray;
    
    GRANT IsIn         ON bool WHERE bool1     IN boolArray;
    GRANT IsNotIn      ON bool WHERE bool1 NOT IN boolArray;
    
    GRANT IsIn         ON str WHERE str1     IN stringArray;
    GRANT IsNotIn      ON str WHERE str1 NOT IN stringArray;
}

TEST CheckInDynamicNum {
    DENY    IsIn ON num POLICY PolInDynamic INPUT {                 }; //numArray: UNSET
    DENY    IsIn ON num POLICY PolInDynamic INPUT { restr: 1        }; //numArray: UNSET
    DENY    IsIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  }; //numArray: UNSET
    DENY    IsIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN }; //numArray: UNSET
    
    DENY IsNotIn ON num POLICY PolInDynamic INPUT {                 }; //numArray: UNSET
    DENY IsNotIn ON num POLICY PolInDynamic INPUT { restr: 1        }; //numArray: UNSET
    DENY IsNotIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  }; //numArray: UNSET
    DENY IsNotIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN }; //numArray: UNSET
    
    //GRANT   IsIn ON num POLICY PolInDynamic INPUT {                   numArray: IGNORE };
    //DENY    IsIn ON num POLICY PolInDynamic INPUT { restr: 1        , numArray: IGNORE };
    //DENY    IsIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  , numArray: IGNORE };
    //DENY    IsIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN , numArray: IGNORE };
    
    //GRANT   IsNotIn ON num POLICY PolInDynamic INPUT {                   numArray: IGNORE };
    //DENY    IsNotIn ON num POLICY PolInDynamic INPUT { restr: 1        , numArray: IGNORE };
    //DENY    IsNotIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  , numArray: IGNORE };
    //DENY    IsNotIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN , numArray: IGNORE };
    
    DENY    IsIn ON num POLICY PolInDynamic INPUT {                   numArray: [] };
    DENY    IsIn ON num POLICY PolInDynamic INPUT { restr: 1        , numArray: [] };
    //GRANT   IsIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  , numArray: [] };
    DENY    IsIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN , numArray: [] };
    
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT {                   numArray: [] };
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT { restr: 1        , numArray: [] };
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  , numArray: [] };
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN , numArray: [] };
    
    DENY    IsIn ON num POLICY PolInDynamic INPUT {                   numArray: [1,2] };
    DENY    IsIn ON num POLICY PolInDynamic INPUT { restr: 42        , numArray: [1,2] };
    //GRANT   IsIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  , numArray: [1,2] };
    EXPECT restr = 1 OR restr = 2 
       FOR IsIn  ON num POLICY PolInDynamic INPUT { restr : UNKNOWN , numArray: [1,2] };
    
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT {                   numArray: [1,2] };
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT { restr: 42       , numArray: [1,2] };
    //GRANT IsNotIn ON num POLICY PolInDynamic INPUT { restr : IGNORE  , numArray: [1,2] };
    //EXPECT restr <> 1 AND restr <> 2
    //  FOR IsNotIn ON num POLICY PolInDynamic INPUT { restr : UNKNOWN , numArray: [1,2] };
}

TEST CheckInDynamicBool {
    DENY    IsIn ON bool POLICY PolInDynamic INPUT {                 }; //boolArray: UNSET
    DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : false   }; //boolArray: UNSET
    DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  }; //boolArray: UNSET
    DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN }; //boolArray: UNSET
    
    DENY IsNotIn ON bool POLICY PolInDynamic INPUT {                 }; //boolArray: UNSET
    DENY IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : false   }; //boolArray: UNSET
    DENY IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  }; //boolArray: UNSET
    DENY IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN }; //boolArray: UNSET
    
    //GRANT   IsIn ON bool POLICY PolInDynamic INPUT {                   boolArray: IGNORE };
    //DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : false   , boolArray: IGNORE };
    //DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  , boolArray: IGNORE };
    //DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN , boolArray: IGNORE };
    
    //GRANT   IsNotIn ON bool POLICY PolInDynamic INPUT {                   boolArray: IGNORE };
    //DENY    IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : false   , boolArray: IGNORE };
    //DENY    IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  , boolArray: IGNORE };
    //DENY    IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN , boolArray: IGNORE };
    
    //DENY    IsIn ON bool POLICY PolInDynamic INPUT {                   boolArray: [] };
    //DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1: false    , boolArray: [] };
    //GRANT   IsIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  , boolArray: [] };
    DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN , boolArray: [] };
    
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT {                   boolArray: [] };
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : false   , boolArray: [] };
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  , boolArray: [] };
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN , boolArray: [] };
    
    DENY    IsIn ON bool POLICY PolInDynamic INPUT {                   boolArray: [true] };
    DENY    IsIn ON bool POLICY PolInDynamic INPUT { bool1 : false   , boolArray: [true] };
    //GRANT   IsIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  , boolArray: [true] };
    EXPECT bool1 = true 
       FOR IsIn  ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN , boolArray: [true] };
    
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT {                   boolArray: [true] };
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : false   , boolArray: [true] };
    //GRANT IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : IGNORE  , boolArray: [true] };
    //EXPECT bool1 <> true
    //  FOR IsNotIn ON bool POLICY PolInDynamic INPUT { bool1 : UNKNOWN , boolArray: [true] };
}

TEST CheckInDynamicStr {
    DENY    IsIn ON str POLICY PolInDynamic INPUT {                }; //stringArray: UNSET
    DENY    IsIn ON str POLICY PolInDynamic INPUT { str1 : ''      }; //stringArray: UNSET
    DENY    IsIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  }; //stringArray: UNSET
    DENY    IsIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN }; //stringArray: UNSET
    
    DENY IsNotIn ON str POLICY PolInDynamic INPUT {                }; //stringArray: UNSET
    DENY IsNotIn ON str POLICY PolInDynamic INPUT { str1 : ''      }; //stringArray: UNSET
    DENY IsNotIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  }; //stringArray: UNSET
    DENY IsNotIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN }; //stringArray: UNSET
    
    DENY   IsIn ON str POLICY PolInDynamic INPUT {                   stringArray: IGNORE };
    //GRANT    IsIn ON str POLICY PolInDynamic INPUT { str1 : ''       , stringArray: IGNORE };
    //GRANT    IsIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  , stringArray: IGNORE };
    //GRANT    IsIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN , stringArray: IGNORE };
    
    //DENY    IsNotIn ON str POLICY PolInDynamic INPUT {                   stringArray: IGNORE };
    //GRANT   IsNotIn ON str POLICY PolInDynamic INPUT { str1 : ''       , stringArray: IGNORE };
    //GRANT   IsNotIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  , stringArray: IGNORE };
    //GRANT   IsNotIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN , stringArray: IGNORE };
    
    DENY    IsIn ON str POLICY PolInDynamic INPUT {                   stringArray: [] };
    DENY    IsIn ON str POLICY PolInDynamic INPUT { str1: ''        , stringArray: [] };
    //GRANT   IsIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  , stringArray: [] };
    DENY    IsIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN , stringArray: [] };
    
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT {                   stringArray: [] };
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT { str1 : ''       , stringArray: [] };
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  , stringArray: [] };
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN , stringArray: [] };
    
    DENY    IsIn ON str POLICY PolInDynamic INPUT {                   stringArray: ['a', 'b'] };
    DENY    IsIn ON str POLICY PolInDynamic INPUT { str1 : ''       , stringArray: ['a', 'b'] };
    GRANT   IsIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  , stringArray: ['a', 'b'] };
    EXPECT str1 = 'a' OR str1 = 'b' 
       FOR IsIn  ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN , stringArray: ['a', 'b'] };
    
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT {                   stringArray: ['a', 'b'] };
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT { str1 : ''      , stringArray: ['a', 'b'] };
    //GRANT IsNotIn ON str POLICY PolInDynamic INPUT { str1 : IGNORE  , stringArray: ['a', 'b'] };
    //EXPECT str1 <> 'a' AND str1 <> 'b'
    //  FOR IsNotIn ON str POLICY PolInDynamic INPUT { str1 : UNKNOWN , stringArray: ['a', 'b'] };
}

POLICY ArrayIN {
    GRANT A ON * WHERE stringval IN ['A', 'B'];
    GRANT B ON * WHERE stringval IN stringArray;
    GRANT C ON * WHERE 'A'       IN stringArray;
    GRANT D ON * WHERE 'A'       IN ['A', 'B'];
    GRANT E ON * WHERE 'A'       IN ['B'];
    GRANT F ON * WHERE stringval IN ['A', 'B'] AND stringval IN ['B', 'C'];
}
TEST ArrayIN {
    // identifier IN constant array
    GRANT A POLICY ArrayIN INPUT { stringval: 'A'    },{ stringval: 'B' };
    DENY  A POLICY ArrayIN INPUT { stringval: 'C'    };
    DENY  A POLICY ArrayIN INPUT {                   }; 
    GRANT A POLICY ArrayIN INPUT { stringval: IGNORE };

    // identifier IN identifier
    GRANT B POLICY ArrayIN INPUT { stringval: 'A',     stringArray: ['A', 'B'] };
    GRANT B POLICY ArrayIN INPUT { stringval: 'B',     stringArray: ['A', 'B'] };
    DENY  B POLICY ArrayIN INPUT { stringval: 'C',     stringArray: ['A', 'B'] };
    DENY  B POLICY ArrayIN INPUT {                     stringArray: ['A', 'B'] }; 
    DENY  B POLICY ArrayIN INPUT { stringval: 'A'                              };
    GRANT B POLICY ArrayIN INPUT { stringval: IGNORE , stringArray: ['A', 'B'] };
    GRANT B POLICY ArrayIN INPUT { stringval: 'A',     stringArray: IGNORE     };

    // const IN identifier
    GRANT C POLICY ArrayIN INPUT { stringArray: ['A']  },{ stringArray: ['A', 'B'] };
    DENY  C POLICY ArrayIN INPUT { stringArray: []     },{ stringArray: ['B'] };
    DENY  C POLICY ArrayIN INPUT {                     }; 
    GRANT C POLICY ArrayIN INPUT { stringArray: IGNORE };

    // const IN const
    GRANT D POLICY ArrayIN INPUT {}; 
    DENY  E POLICY ArrayIN INPUT {}; 

    GRANT F POLICY ArrayIN INPUT { stringval: 'B'    };
    DENY  F POLICY ArrayIN INPUT { stringval: 'A'    },{ stringval: 'C' };
    DENY  F POLICY ArrayIN INPUT {                   };
    GRANT F POLICY ArrayIN INPUT { stringval: IGNORE };
}