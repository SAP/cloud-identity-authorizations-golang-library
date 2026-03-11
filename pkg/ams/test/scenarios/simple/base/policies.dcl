
DEFAULT POLICY PublicStuff {
    GRANT read on r1 where entity1.public;
    GRANT read on r2 where $user.email in entity2.owners;
    GRANT read on r1 where entity1.group in $user.groups;
}

POLICY R1NameBased {
    GRANT read on r1 where entity1.name is RESTRICTED;
}

POLICY R2SubentityBased {
    GRANT read on r2 where entity2.subentity.subNumberField is RESTRICTED;
}

POLICY R1SizeBased {
    GRANT read on r1 where entity1.size is RESTRICTED;
}

POLICY R1SmallSized {
    USE R1SizeBased RESTRICT entity1.size < 100;
}
