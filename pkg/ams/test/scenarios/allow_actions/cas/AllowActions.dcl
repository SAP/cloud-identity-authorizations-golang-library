POLICY Empty {
}

POLICY Base {
	GRANT A ON *;
	GRANT B ON Y;
	GRANT C ON * WHERE dummy = '';
	GRANT D ON * WHERE dummy IS NOT RESTRICTED;
}

POLICY Use1 {
	USE Empty;
}

POLICY Use2 {
	USE Empty;
	GRANT E ON *;
}

POLICY Use3 {
	USE Base;
}

POLICY Use4 {
	USE Base RESTRICT dummy = 'x';
	GRANT F ON Z;
}

POLICY Use5 {
	USE Base RESTRICT dummy = 'x';
	GRANT G ON *;
}

POLICY Use6 {
	USE Use3;
}
