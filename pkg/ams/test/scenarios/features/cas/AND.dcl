///
Policy PolANDTrue {
    GRANT IsAND ON * WHERE bool1 = true AND bool2 = true;
}

Test ANDCheckTrue {
    GRANT                   IsAND POLICY PolANDTrue INPUT { bool1: true,    bool2: true    };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: true,    bool2: false   };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: true                    };
    GRANT                   IsAND POLICY PolANDTrue INPUT { bool1: true,    bool2: IGNORE  };
    EXPECT bool2 = true FOR IsAND POLICY PolANDTrue INPUT { bool1: true,    bool2: UNKNOWN };

    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: false,   bool2: true    };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: false,   bool2: false   };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: false                   };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: false,   bool2: IGNORE  };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: false,   bool2: UNKNOWN };

    DENY                    IsAND POLICY PolANDTrue INPUT {                 bool2: true    };
    DENY                    IsAND POLICY PolANDTrue INPUT {                 bool2: false   };
    DENY                    IsAND POLICY PolANDTrue INPUT {                                };
    DENY                    IsAND POLICY PolANDTrue INPUT {                 bool2: IGNORE  };
    DENY                    IsAND POLICY PolANDTrue INPUT {                 bool2: UNKNOWN };


    GRANT                   IsAND POLICY PolANDTrue INPUT { bool1: IGNORE,  bool2: true     };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: IGNORE,  bool2: false    };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: IGNORE                   };
    GRANT                   IsAND POLICY PolANDTrue INPUT { bool1: IGNORE,  bool2: IGNORE   };
    EXPECT bool2 = true FOR IsAND POLICY PolANDTrue INPUT { bool1: IGNORE,  bool2: UNKNOWN  };

    EXPECT bool1 = true FOR IsAND POLICY PolANDTrue INPUT { bool1: UNKNOWN, bool2: true     };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: UNKNOWN, bool2: false    };
    DENY                    IsAND POLICY PolANDTrue INPUT { bool1: UNKNOWN                  };
    EXPECT bool1 = true FOR IsAND POLICY PolANDTrue INPUT { bool1: UNKNOWN, bool2: IGNORE   };
    EXPECT bool1 = true AND
           bool2 = true FOR IsAND POLICY PolANDTrue INPUT { bool1: UNKNOWN, bool2: UNKNOWN  };
}

///
Policy PolANDFalse {
    GRANT IsAND ON * WHERE bool1 = false AND bool2 = false;
}

Test ANDCheckFalse {
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: true,    bool2: true    };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: true,    bool2: false   };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: true                    };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: true,    bool2: IGNORE  };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: true,    bool2: UNKNOWN };

    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: false,   bool2: true    };
    GRANT                    IsAND POLICY PolANDFalse INPUT { bool1: false,   bool2: false   };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: false                   };
    GRANT                    IsAND POLICY PolANDFalse INPUT { bool1: false,   bool2: IGNORE  };
    EXPECT bool2 = false FOR IsAND POLICY PolANDFalse INPUT { bool1: false,   bool2: UNKNOWN };


    DENY                     IsAND POLICY PolANDFalse INPUT {                 bool2: true    };
    DENY                     IsAND POLICY PolANDFalse INPUT {                 bool2: false   };
    DENY                     IsAND POLICY PolANDFalse INPUT {                                };
    DENY                     IsAND POLICY PolANDFalse INPUT {                 bool2: IGNORE  };
    DENY                     IsAND POLICY PolANDFalse INPUT {                 bool2: UNKNOWN };

    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: IGNORE,  bool2: true    };
    GRANT                    IsAND POLICY PolANDFalse INPUT { bool1: IGNORE,  bool2: false   };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: IGNORE                  };
    GRANT                    IsAND POLICY PolANDFalse INPUT { bool1: IGNORE,  bool2: IGNORE  };
    EXPECT bool2 = false FOR IsAND POLICY PolANDFalse INPUT { bool1: IGNORE,  bool2: UNKNOWN };

    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: UNKNOWN, bool2: true    };
    EXPECT bool1 = false FOR IsAND POLICY PolANDFalse INPUT { bool1: UNKNOWN, bool2: false   };
    DENY                     IsAND POLICY PolANDFalse INPUT { bool1: UNKNOWN                 };
    EXPECT bool1 = false FOR IsAND POLICY PolANDFalse INPUT { bool1: UNKNOWN, bool2: IGNORE  };
    EXPECT bool1 = false AND
           bool2 = false FOR IsAND POLICY PolANDFalse INPUT { bool1: UNKNOWN, bool2: UNKNOWN };
}


POLICY ANDCornerCases {
    GRANT CC1a ON * WHERE bool1 = true AND 5 > 7;
    GRANT CC1b ON * WHERE bool1 = true AND 7 > 5;

    GRANT CC2  ON * WHERE bool1 = true AND bool1 = false;

    GRANT CC3a ON * WHERE bool1 = true AND bool1 IS NULL;
    GRANT CC3b ON * WHERE bool1 = true AND bool2 IS NULL;
    
}

TEST ANDCornerCases {
    DENY                     CC1a POLICY ANDCornerCases INPUT { bool1: true       };
    DENY                     CC1a POLICY ANDCornerCases INPUT { bool1: false      };
    DENY                     CC1a POLICY ANDCornerCases INPUT {                   };
    DENY                     CC1a POLICY ANDCornerCases INPUT { bool1: IGNORE     };
    DENY                     CC1a POLICY ANDCornerCases INPUT { bool1: UNKNOWN    };
//
    GRANT                    CC1b POLICY ANDCornerCases INPUT { bool1: true       };
    DENY                     CC1b POLICY ANDCornerCases INPUT { bool1: false      };
    DENY                     CC1b POLICY ANDCornerCases INPUT {                   };
    GRANT                    CC1b POLICY ANDCornerCases INPUT { bool1: IGNORE     };
    EXPECT bool1 = true FOR  CC1b POLICY ANDCornerCases INPUT { bool1: UNKNOWN    };
//
    DENY                     CC2  POLICY ANDCornerCases INPUT { bool1: true       };
    DENY                     CC2  POLICY ANDCornerCases INPUT { bool1: false      };
    DENY                     CC2  POLICY ANDCornerCases INPUT {                   };
    GRANT                    CC2  POLICY ANDCornerCases INPUT { bool1: IGNORE     };
//    DENY                     CC2  POLICY ANDCornerCases INPUT { bool1: UNKNOWN    };
//
    DENY                     CC3a POLICY ANDCornerCases INPUT { bool1: true       };
    DENY                     CC3a POLICY ANDCornerCases INPUT { bool1: false      };
    DENY                     CC3a POLICY ANDCornerCases INPUT {                   };
    GRANT                    CC3a POLICY ANDCornerCases INPUT { bool1: IGNORE     };
    EXPECT bool1 = true AND
           bool1 IS NULL FOR CC3a POLICY ANDCornerCases INPUT { bool1: UNKNOWN    };
//
    GRANT                    CC3b POLICY ANDCornerCases INPUT { bool1: true       };
    DENY                     CC3b POLICY ANDCornerCases INPUT { bool1: false      };
    DENY                     CC3b POLICY ANDCornerCases INPUT {                   };
    GRANT                    CC3b POLICY ANDCornerCases INPUT { bool1: IGNORE     };
    EXPECT bool1 = true FOR  CC3b POLICY ANDCornerCases INPUT { bool1: UNKNOWN    };
    DENY                     CC3b POLICY ANDCornerCases INPUT { bool1: IGNORE,  bool2: true    };
    EXPECT bool2 IS NULL FOR CC3b POLICY ANDCornerCases INPUT { bool1: IGNORE,  bool2: UNKNOWN };
    DENY                     CC3b POLICY ANDCornerCases INPUT { bool1: UNKNOWN, bool2: true    };
    EXPECT bool1 = true FOR  CC3b POLICY ANDCornerCases INPUT { bool1: UNKNOWN, bool2: IGNORE  };
//
}
