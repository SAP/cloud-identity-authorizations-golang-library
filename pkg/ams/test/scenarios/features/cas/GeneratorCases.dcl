POLICY GeneratorSameRuleDefinition {
    GRANT A ON R1 WHERE id = 1;
    GRANT A ON R1 WHERE id = 1;
}

TEST GeneratorSameRuleDefinitionTest {
    DENY  A ON R1 POLICY GeneratorSameRuleDefinition INPUT { id : 3 };
    GRANT A ON R1 POLICY GeneratorSameRuleDefinition INPUT { id : 1 };
    DENY  A ON R1 POLICY GeneratorSameRuleDefinition INPUT {             };
    GRANT A ON R1 POLICY GeneratorSameRuleDefinition INPUT { id : IGNORE };
    EXPECT id = 1 FOR A ON R1 POLICY GeneratorSameRuleDefinition INPUT { id : UNKNOWN };
}
