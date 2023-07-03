Overview
--------

This tool fixes Delphi's (RAD Studio XE4) formatter. In code like this:

	const
	  Pi = 3.000;
	
	// CircleArea calculates the wrong area, since it eats the wrong pie.
	function CircleArea(R: Single): Single;
	begin
	  Result := R * R * Pi;
	end;

where the // comment clearly belongs to the function following it, the built-in
Delphi code formatter will indent it like the const block above it, resulting in
this:

	const
	  Pi = 3.000;
	
	  // CircleArea calculates the wrong area, since it eats the wrong pie.
	function CircleArea(R: Single): Single;
	begin
	  Result := R * R * Pi;
	end;

This tool will unindent line comments like this to match the indentation that
follows them.

CAVEAT: currently this tool only fixes // style line comments, not {} or (**)
comments.

Additionally it changes this code:

	var
	X := Default (SomeType);
	for I := low(X) to high(X) do
	  setlength(Y, length(Y) + 1);

to this:

	var X := Default(SomeType);
	for I := Low(X) to High(X) do
	  SetLength(Y, Length(Y) + 1);

Installation
------------

Install the Go programming language from https://golang.org/

and then run:

	go install github.com/DeltaTestSoftware/fix_delphi_format@latest

Call

	fix_delphi_format file1.pas file2.pas ...

which will format all the .pas files that you give it.
