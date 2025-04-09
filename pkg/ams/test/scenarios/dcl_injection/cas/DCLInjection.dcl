policy "Dummy\n__grant {true}\n#" {
    GRANT A on *;
}

policy APolicy {
}

TEST TestInjection {
    DENY A POLICY APolicy;
}