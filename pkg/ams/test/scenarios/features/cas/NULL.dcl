POLICY PolNull {
    GRANT IsNull    ON * WHERE nullval IS NULL;
    GRANT IsNotNull ON * WHERE nullval IS NOT NULL;
}

TEST NullCheck {
    DENY  IsNull    POLICY PolNull INPUT { nullval : true   };
    GRANT IsNull    POLICY PolNull INPUT {                  };
    GRANT IsNull    POLICY PolNull INPUT { nullval : IGNORE };
    EXPECT nullval IS NULL FOR IsNull ON * POLICY PolNull INPUT { nullval : UNKNOWN };
    
    GRANT IsNotNull POLICY PolNull INPUT { nullval : true   };
    DENY  IsNotNull POLICY PolNull INPUT {                  };
    GRANT IsNotNull POLICY PolNull INPUT { nullval : IGNORE };
    EXPECT nullval IS NOT NULL FOR IsNotNull ON * POLICY PolNull INPUT { nullval : UNKNOWN };
}
