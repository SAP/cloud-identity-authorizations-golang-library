FUNCTION Boolean T()    {
	return n = 1;
}

POLICY x {
	GRANT A on X WHERE a and b;
}

POLICY y {
	GRANT A on X WHERE T() and b;
}


TEST x{
	GRANT A on X POLICY x INPUT {a:true,b:true};
	EXPECT a FOR A ON X POLICY x INPUT{a:UNKNOWN,b:true};
	EXPECT a and b FOR A ON X POLICY x INPUT{a:UNKNOWN,b:UNKNOWN};
}

TEST y{
	GRANT A on X POLICY y INPUT {n:1,b:true};
	EXPECT n=1 FOR A ON X POLICY y INPUT{n:UNKNOWN,b:true};
	EXPECT n=1 and b FOR A ON X POLICY y INPUT{n:UNKNOWN,b:UNKNOWN};
}