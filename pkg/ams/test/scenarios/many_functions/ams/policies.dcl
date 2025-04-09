//
// c = carry bit
// v = value bit
// cxT = carry bit x (value)
// cxN = carry bit x (negated value)
//
FUNCTION Boolean c0T() { return a0=true  AND b0= true; }
FUNCTION Boolean c0N() { return a0<>true OR  b0<>true; }
FUNCTION Boolean v0T() { return a0<>b0;                }
FUNCTION Boolean v0N() { return a0= b0;                }


FUNCTION Boolean c1T() { return (c0N() AND a1=true  AND b1=true )  OR (c0T() AND (a1=true  OR b1=true )); }
FUNCTION Boolean c1N() { return (c0N() AND (a1=false OR b1=false)) OR (c0T() AND a1=false AND b1=false);  }
FUNCTION Boolean v1T() { return (c0N() AND a1<>b1)                 OR (c0T() AND a1= b1);                 }
FUNCTION Boolean v1N() { return (c0N() AND a1= b1)                 OR (c0T() AND a1<>b1);                 }


FUNCTION Boolean c2T() { return (c1N() AND a2=true  AND b2=true )  OR (c1T() AND (a2=true  OR b2=true )); }
FUNCTION Boolean c2N() { return (c1N() AND (a2=false OR b2=false)) OR (c1T() AND a2=false AND b2=false);  }
FUNCTION Boolean v2T() { return (c1N() AND a2<>b2)                 OR (c1T() AND a2= b2);                 }
FUNCTION Boolean v2N() { return (c1N() AND a2= b2)                 OR (c1T() AND a2<>b2);                 }



FUNCTION Boolean r0()  { return (v0T() AND r0=true) OR (v0N() AND r0=false); }
FUNCTION Boolean r1()  { return (v1T() AND r1=true) OR (v1N() AND r1=false); }
FUNCTION Boolean r2()  { return (v2T() AND r2=true) OR (v2N() AND r2=false); }
FUNCTION Boolean r3()  { return (c2T() AND r3=true) OR (c2N() AND r3=false); }


DEFAULT POLICY X       { GRANT * ON * WHERE r0() AND r1() AND r2() AND r3(); }

TEST X {
    GRANT * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:false, r2:false, r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:false, r2:true , r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:true , r2:false, r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:true , r2:true , r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:false, r2:false, r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:false, r2:true , r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:true , r2:false, r3: false };
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:true , r2:true , r3: false };
    EXPECT r0=false AND r1=false AND r2=false AND r3=false    FOR * ON * INPUT { a0: false, a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:UNKNOWN , r1:UNKNOWN , r2:UNKNOWN, r3:UNKNOWN  };

    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:false, r2:false, r3: false };
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:false, r2:true , r3: false };
    GRANT * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:true , r2:false, r3: false };
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:true , r2:true , r3: false };
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:false, r2:false, r3: false };
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:false, r2:true , r3: false };
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:true , r2:false, r3: false };
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:true , r2:true , r3: false };
    EXPECT r0=false AND r1=true AND r2=false AND r3=false     FOR * ON * INPUT { a0: true , a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:UNKNOWN , r1:UNKNOWN , r2:UNKNOWN, r3:UNKNOWN  };

    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:false, r1:false, r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:false, r1:false, r2:true  , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:false, r1:true , r2:false , r3: false};
    GRANT * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:false, r1:true , r2:true  , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:true , r1:false, r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:true , r1:false, r2:true  , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:true , r1:true , r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:true , r1:true , r2:true  , r3: false};
    EXPECT r0=false AND r1=true AND r2=true AND r3=false      FOR * ON * INPUT { a0: true , a1: true , a2: false,   b0:true , b1:true , b2: false,   r0:UNKNOWN , r1:UNKNOWN , r2:UNKNOWN, r3:UNKNOWN  };

    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:false, r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:false, r2:true  , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:true , r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:false, r1:true , r2:true  , r3: false};
    GRANT * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:false, r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:false, r2:true  , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:true , r2:false , r3: false};
    DENY  * ON * INPUT { a0: true,  a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:true , r1:true , r2:true  , r3: false};
    EXPECT r0=true AND r1=false AND r2=false AND r3=false     FOR * ON * INPUT { a0: true , a1: false, a2: false,   b0:false, b1:false, b2: false,   r0:UNKNOWN , r1:UNKNOWN , r2:UNKNOWN, r3:UNKNOWN  };

    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:false, r2:false , r3: false};
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:false, r2:true  , r3: false};
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:true , r2:false , r3: false};
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:false, r1:true , r2:true  , r3: false};
    GRANT * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:false, r2:false , r3: false};
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:false, r2:true  , r3: false};
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:true , r2:false , r3: false};
    DENY  * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:true , r1:true , r2:true  , r3: false};
    EXPECT r0=true AND r1=false AND r2=false AND r3=false     FOR * ON * INPUT { a0: false, a1: false, a2: false,   b0:true , b1:false, b2: false,   r0:UNKNOWN , r1:UNKNOWN , r2:UNKNOWN, r3:UNKNOWN  };


//
// The output represents the minimal possible solution.
// Simple engines may have a hard time to reach this simplification  
//
//  EXPECT a0=true AND a1=false AND a2=false     FOR * ON * INPUT { a0: UNKNOWN, a1: UNKNOWN, a2: UNKNOWN,   b0:true , b1:false, b2: false,   r0:false , r1:true , r2:false, r3:false  };

 }
