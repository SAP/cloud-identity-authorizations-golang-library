Policy PolLike {
    GRANT IsLike     ON * WHERE stringval     LIKE '%TEST%';
    GRANT IsNotLike  ON * WHERE stringval NOT LIKE '%TEST%';

    GRANT IsLike2    ON * WHERE stringval     LIKE '_TEST_';
    GRANT IsLike3    ON * WHERE stringval     LIKE '_T*E{S}T{1}X++%_';
    GRANT IsLike4    ON * WHERE stringval     LIKE '^_X';
    GRANT IsLike5    ON * WHERE stringval     LIKE 'A^_X';
    GRANT IsLike6    ON * WHERE stringval     LIKE '$';
    GRANT IsLikeEscpape ON * WHERE stringval LIKE 'xy' ESCAPE 'b';
}


TEST LikeCheck2 {
    GRANT IsLike2    POLICY PolLike INPUT { stringval: '1TEST1'    };
    DENY  IsLike2    POLICY PolLike INPUT { stringval: '1TEST'     };
    DENY  IsLike2    POLICY PolLike INPUT { stringval: 'TEST1'     };
    DENY  IsLike2    POLICY PolLike INPUT {                        };
    GRANT IsLike2    POLICY PolLike INPUT { stringval:  IGNORE	   };

    GRANT IsLike3    POLICY PolLike INPUT { stringval: '_T*E{S}T{1}X++!_' };
    GRANT IsLike3    POLICY PolLike INPUT { stringval: '_T*E{S}T{1}X++_'  };
    DENY  IsLike3    POLICY PolLike INPUT {                               };
    GRANT IsLike3    POLICY PolLike INPUT { stringval:  IGNORE	          };

    GRANT IsLike4    POLICY PolLike INPUT { stringval: '^AX'       };
    DENY  IsLike4    POLICY PolLike INPUT { stringval: '^AY'       };
    DENY  IsLike4    POLICY PolLike INPUT {                        };
    GRANT IsLike4    POLICY PolLike INPUT { stringval:  IGNORE	   };
    
    GRANT IsLike5    POLICY PolLike INPUT { stringval: 'A^AX'      };
    DENY  IsLike5    POLICY PolLike INPUT { stringval: 'B^AX'      };
    DENY  IsLike5    POLICY PolLike INPUT {                        };
    GRANT IsLike5    POLICY PolLike INPUT { stringval:  IGNORE	   };
    
    GRANT IsLike6    POLICY PolLike INPUT { stringval: '$'         };
    DENY  IsLike6    POLICY PolLike INPUT { stringval: 'A'         };
    DENY  IsLike6    POLICY PolLike INPUT {                        };
    GRANT IsLike6    POLICY PolLike INPUT { stringval:  IGNORE	   };
}


TEST LikeCheck {
    GRANT IsLike    POLICY PolLike INPUT { stringval: 'TEST'      };
    GRANT IsLike    POLICY PolLike INPUT { stringval: 'xTESTx'    };
    DENY  IsLike	POLICY PolLike INPUT {                        };
    GRANT IsLike 	POLICY PolLike INPUT { stringval:  IGNORE	  };

    DENY  IsNotLike POLICY PolLike INPUT { stringval: 'TEST'      };
    DENY  IsNotLike POLICY PolLike INPUT { stringval: 'xTESTx'    };
    DENY  IsNotLike POLICY PolLike INPUT {                        };
    GRANT IsNotLike POLICY PolLike INPUT { stringval:  IGNORE	  };
    
    DENY  IsLike    POLICY PolLike INPUT { stringval: 'MISSING'   };
    DENY  IsLike    POLICY PolLike INPUT { stringval: 'xMISSINGx' };
    GRANT IsNotLike POLICY PolLike INPUT { stringval: 'MISSING'   };
    GRANT IsNotLike POLICY PolLike INPUT { stringval: 'xMISSINGx' };
}
