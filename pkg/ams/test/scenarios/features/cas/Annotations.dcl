@description: 'Annotated Function with special characters ßÖÜÄ'
FUNCTION Boolean AnnotatedFunction() {
    return stringval = 'a';
}

POLICY UsedPolicyWithAnnotations {
}


@label: 'Annotated Policy with special characters ßÖÜÄ'
@description: 'POLICY description'
POLICY AnnotatedPolicy {
    @Fast: TRUE
    GRANT * ON *;
    
    @ForNow: True
    USE UsedPolicyWithAnnotations;
}


@description: 'TEST description'
TEST testWithAnnotation {
    @IgnoreOnWindows8: TRUE
    DENY * ON *;
}
