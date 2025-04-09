POLICY SameInput {
    GRANT A ON * WHERE $same.x = ''     and $same.x = '';
    GRANT B ON * WHERE "stringval" = '' and stringval = '';
}

TEST SameInputTest {
    GRANT A POLICY SameInput INPUT { $same:   {  x : ''  } },
                                   { "$same": {  x : ''  } },
                                   { $same:   { "x": ''  } },
                                   { "$same": { "x": ''  } };
    DENY  A POLICY SameInput INPUT { $same:   {  x : 'a' } },
                                   { "$same": {  x : 'a' } },
                                   { $same:   { "x": 'a' } },
                                   { "$same": { "x": 'a' } };
    DENY  A POLICY SameInput INPUT {                       };
    GRANT A POLICY SameInput INPUT { $same:   {  x : IGNORE } },
                                   { "$same": {  x : IGNORE } },
                                   { $same:   { "x": IGNORE } },
                                   { "$same": { "x": IGNORE } };
    EXPECT $same.x = '' FOR A ON * POLICY SameInput INPUT { $same:   {  x : UNKNOWN } },
                                                          { "$same": {  x : UNKNOWN } },
                                                          { $same:   { "x": UNKNOWN } },
                                                          { "$same": { "x": UNKNOWN } };

    GRANT B POLICY SameInput INPUT { "stringval" : ''      },
                                   { stringval   : ''      };
    DENY  B POLICY SameInput INPUT { "stringval" : 'a'     },
                                   { stringval   : 'a'     };
    DENY  B POLICY SameInput INPUT {                       };
    GRANT B POLICY SameInput INPUT { "stringval" : IGNORE  },
                                   { stringval   : IGNORE  };
    EXPECT "stringval" = '' FOR B ON * POLICY SameInput INPUT { "stringval" : UNKNOWN };
    EXPECT stringval = ''   FOR B ON * POLICY SameInput INPUT { stringval : UNKNOWN   };
}

POLICY Access_Quoted_SubName {
    GRANT * ON * WHERE $Struct."Quoted-sub-name" IS NOT NULL;
}

TEST Access_Quoted_SubName_Test{
    GRANT * POLICY Access_Quoted_SubName INPUT { $Struct: { "Quoted-sub-name": 'a'                     } };
    DENY  * POLICY Access_Quoted_SubName INPUT { $Struct: {                         anyOtherName: 'x'  } };
    DENY  * POLICY Access_Quoted_SubName INPUT {                                                         };
    GRANT * POLICY Access_Quoted_SubName INPUT { $Struct: { "Quoted-sub-name": IGNORE                  } };
    EXPECT $Struct."Quoted-sub-name" IS NOT NULL FOR * ON * POLICY Access_Quoted_SubName INPUT { $Struct: { "Quoted-sub-name": UNKNOWN} };
}

// Added with version 0.9.x
// 
POLICY Access$user {
    GRANT A on * WHERE $user.user_uuid = 'a';
}

TEST Access$userTest{
    GRANT A POLICY Access$user INPUT { $user: { "user_uuid": 'a'    } };
    DENY  A POLICY Access$user INPUT { $user: { "user_uuid": 'b'    } };
    DENY  A POLICY Access$user INPUT {                                };
    GRANT A POLICY Access$user INPUT { $user: { "user_uuid": IGNORE } };
    EXPECT $user."user_uuid" = 'a' FOR A ON * POLICY Access$user INPUT { $user: { "user_uuid": UNKNOWN} };
}

///////////////////////////////////////////////////////////


POLICY Jürgen_s_unicode_test_with_dotted_and_dotless_I__Iiİı_ {
    GRANT A on *;
}

POLICY QuotedReference {
    USE "Jürgen_s_unicode_test_with_dotted_and_dotless_I__Iiİı_";
}

TEST QuotedReferenceTest {
    GRANT A POLICY "Jürgen_s_unicode_test_with_dotted_and_dotless_I__Iiİı_";
    DENY  B POLICY "Jürgen_s_unicode_test_with_dotted_and_dotless_I__Iiİı_";
    GRANT A POLICY QuotedReference;
    DENY  B POLICY QuotedReference;
}

POLICY "$Policy" {
    GRANT A on *;
}

POLICY DollarPolicy {
    USE "$Policy";
}

Test DollarPolicyTest {
    GRANT A POLICY "$Policy";
    DENY  B POLICY "$Policy";
    GRANT A POLICY DollarPolicy;
    DENY  B POLICY DollarPolicy;
}

POLICY Quote {
    GRANT action on * where "\"quoted\"" = 'a' and "\"quoted2\"".findme ='';
}

TEST QuoteTest {
    GRANT action ON * POLICY Quote INPUT { "\"quoted\"": 'a',    "\"quoted2\"": { findme: ''      }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": 'b',    "\"quoted2\"": { findme: ''      }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": 'a',    "\"quoted2\"": { findme: 'a'     }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": 'b',    "\"quoted2\"": { findme: 'b'     }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": 'a'                                       };
    DENY  action ON * POLICY Quote INPUT {                       "\"quoted2\"": { findme: ''      }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": IGNORE                                    };
    DENY  action ON * POLICY Quote INPUT {                       "\"quoted2\"": { findme: IGNORE  }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": UNKNOWN                                   };
    DENY  action ON * POLICY Quote INPUT {                       "\"quoted2\"": { findme: UNKNOWN }};
    DENY  action ON * POLICY Quote INPUT {                                                         };
    GRANT action ON * POLICY Quote INPUT { "\"quoted\"": 'a',    "\"quoted2\"": { findme: IGNORE  }};
    GRANT action ON * POLICY Quote INPUT { "\"quoted\"": IGNORE, "\"quoted2\"": { findme: ''      }};
    DENY  action ON * POLICY Quote INPUT { "\"quoted\"": 'b',    "\"quoted2\"": { findme: IGNORE  }};
    GRANT action ON * POLICY Quote INPUT { "\"quoted\"": IGNORE, "\"quoted2\"": { findme: IGNORE  }};
    EXPECT "\"quoted2\"".findme ='' FOR action ON * POLICY Quote INPUT { "\"quoted\"": 'a', "\"quoted2\"": { findme: UNKNOWN  }};
    EXPECT "\"quoted\"" = 'a'       FOR action ON * POLICY Quote INPUT { "\"quoted\"": UNKNOWN, "\"quoted2\"": { findme: ''   }};
    EXPECT "\"quoted\"" = 'a' and "\"quoted2\"".findme ='' FOR action ON * POLICY Quote INPUT { "\"quoted\"": UNKNOWN, "\"quoted2\"": { findme: UNKNOWN }};
}


POLICY "polücywithoutquÖtes" { USE "$Policy" ; }

POLICY "REF__polücywithoutquÖtes" { USE "polücywithoutquÖtes"; }

Test "REF__polücywithoutquÖtes_TEST" {
    GRANT A ON * POLICY "REF__polücywithoutquÖtes";
    DENY  B ON * POLICY "REF__polücywithoutquÖtes";
    GRANT A ON * POLICY polücywithoutquÖtes;
    DENY  B ON * POLICY polücywithoutquÖtes;
}

// All Special with Special names
//
//FUNCTION Boolean "special: \"\n*"() {
//	return 1=1;
//}
//
//POLICY "special: \"\n*" { USE "$Policy" where "special: \"\n*"(); }
//INPUT "special: \"\n*" {
//	"\"quoted2\"": {
//		findme: ''
//	}
//}
//
//TEST "special: \"\n*" {
//	DENY "special: \"\n*" ON "special: \"\n*" POLICY "special: \"\n*" INPUT "special: \"\n*";
//}
// All Special



//