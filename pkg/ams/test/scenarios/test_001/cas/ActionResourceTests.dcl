POLICY StarOnStar {
    GRANT * ON *;
}

POLICY StarOnSingleResource {
    GRANT * ON Resource1;
}

POLICY SingleActionOnStar {
    GRANT Action1 ON *;
}

POLICY SingleActionOnSingleResource {
    GRANT Action1 ON Resource1;
}

POLICY MultiActionOnSingleResource {
    GRANT Action1, Action2 ON Resource1;
}

POLICY SingleActionOnMultiResource {
    GRANT Action1 ON Resource1, Resource2;
}

POLICY MultiActionOnMultiResource {
    GRANT Action1, Action2 ON Resource1, Resource2;
}

POLICY UseSingleActionOnSingleResource {
    USE SingleActionOnSingleResource;
}

POLICY A1onR1 {
    GRANT A1 ON R1;
}

POLICY A1onR2 {
    GRANT A1 ON R2;
}

POLICY A2onR1 {
    GRANT A2 ON R1;
}

POLICY A2onR2 {
    GRANT A2 ON R2;
}
TEST ActionResource {
    GRANT * POLICY StarOnStar INPUT {};
    
    DENY * POLICY StarOnSingleResource INPUT {};
    
    GRANT Action1 POLICY SingleActionOnStar INPUT {};
    DENY Action2 POLICY SingleActionOnStar INPUT {};
    DENY * POLICY SingleActionOnStar INPUT {};

    GRANT Action1 ON Resource1 POLICY SingleActionOnSingleResource;
    DENY Action2 ON Resource1 POLICY SingleActionOnSingleResource;
    DENY Action1 ON Resource2 POLICY SingleActionOnSingleResource;
    DENY * ON Resource1 POLICY SingleActionOnSingleResource; 
    DENY Action1 ON * POLICY SingleActionOnSingleResource; 

    GRANT Action1 ON Resource1 POLICY MultiActionOnSingleResource;
    GRANT Action2 ON Resource1 POLICY MultiActionOnSingleResource;
    GRANT Action1, Action2 ON Resource1 POLICY MultiActionOnSingleResource;
    DENY * ON Resource1 POLICY MultiActionOnSingleResource;
    //GRANT Action1, Action3 ON Resource1 POLICY MultiActionOnSingleResource;
    GRANT Action1 ON Resource1 POLICY MultiActionOnSingleResource;
    DENY Action3 ON Resource1 POLICY MultiActionOnSingleResource;

    DENY Action1, Action2 ON * POLICY MultiActionOnSingleResource; 
    
    GRANT Action1 ON Resource1 POLICY SingleActionOnMultiResource;
    GRANT Action1 ON Resource2 POLICY SingleActionOnMultiResource;
    GRANT Action1 ON Resource1, Resource2 POLICY SingleActionOnMultiResource;
    DENY Action1 ON * POLICY SingleActionOnMultiResource;
    //GRANT Action1 ON Resource1, Resource3 POLICY SingleActionOnMultiResource;
    //expands into two assertion with different results -> always fails regardless which outcome we expect

    GRANT Action1 ON Resource1 POLICY SingleActionOnMultiResource;
    DENY Action1 ON Resource3 POLICY SingleActionOnMultiResource; 

    DENY * ON Resource1, Resource2 POLICY SingleActionOnMultiResource; 

    GRANT Action1 ON Resource1 POLICY MultiActionOnMultiResource;
    GRANT Action2 ON Resource2 POLICY MultiActionOnMultiResource;
    GRANT Action1, Action2 ON Resource2 POLICY MultiActionOnMultiResource;
    GRANT Action1, Action2 ON Resource1 POLICY MultiActionOnMultiResource;
    GRANT Action1 ON Resource1, Resource2 POLICY MultiActionOnMultiResource;
    DENY * ON Resource1, Resource2 POLICY MultiActionOnMultiResource;
    DENY Action3 ON Resource1, Resource2 POLICY MultiActionOnMultiResource;
    DENY Action1, Action2 ON Resource3 POLICY MultiActionOnMultiResource;
    DENY Action1, Action2 ON * POLICY MultiActionOnMultiResource;
}

//TEST CompositeA {
    //GRANT A1,A2 ON R1,R2 POLICY A1onR1;
    // short hand for
    // GRANT A1 ON R1 POLICY A1onR1; -> true
    // GRANT A1 ON R2 POLICY A1onR1; -> false
    // GRANT A2 ON R1 POLICY A1onR1; -> false
    // GRANT A2 ON R2 POLICY A1onR1; -> false
    // -> Rego test fails on second assertion
//}

//TEST NotCompositeA {
    // DENY A1,A2 ON R1,R2 POLICY A1onR1;
    // short hand for
    // DENY A1 ON R1 POLICY A1onR1; -> false
    // DENY A1 ON R2 POLICY A1onR1; -> true
    // DENY A2 ON R1 POLICY A1onR1; -> true
    // DENY A2 ON R2 POLICY A1onR1; -> true
    // -> Rego test fails on first assertion
//}

//TEST CompositeB {
//    GRANT A1,A2 ON R1,R2 POLICY A1onR1,A1onR2;
//}

//TEST CompositeC {
//    GRANT A1,A2 ON R1,R2 POLICY A1onR1,A1onR2,A2onR1;
//}

TEST CompositeA {
    GRANT A1,A2 ON R1,R2 POLICY A1onR1,A1onR2,A2onR1,A2onR2;
}

TEST CompositeB {
    DENY A1,A2 ON R1,R2;
}

TEST UsePolicy {
    GRANT Action1 ON Resource1 POLICY UseSingleActionOnSingleResource;
    DENY Action2 ON Resource1 POLICY UseSingleActionOnSingleResource;
    DENY Action1 ON Resource2 POLICY UseSingleActionOnSingleResource;
    DENY * ON Resource1 POLICY UseSingleActionOnSingleResource; 
    DENY Action1 ON * POLICY UseSingleActionOnSingleResource; 
}