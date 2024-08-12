package rule

var MAP_OPERATOR = map[string]string{
	"pmFromFile":        "@pmFromFile",        //": Checks if the value matches a pattern defined in a file.",
	"rx":                "@rx",                //": Performs a regular expression match against the value.",
	"in":                "@pm",                //": Performs a pattern match against the value.",
	"not in":            "@pm",                //": Performs a pattern match against the value.",
	"streq":             "@streq",             //": Checks if the value is exactly equal to the provided string.",
	"gt":                "@gt",                //": Checks if the value is greater than the provided value.",
	"ge":                "@ge",                //": Checks if the value is greater than or equal to the provided value.",
	"lt":                "@lt",                //": Checks if the value is less than the provided value.",
	"le":                "@le",                //": Checks if the value is less than or equal to the provided value.",
	"within":            "@within",            //": Checks if the value is within the specified range.",
	"contains":          "@contains",          //": Checks if the value contains a specific string.",
	"beginsWith":        "@beginsWith",        //": Checks if the value starts with a specific string.",
	"endsWith":          "@endsWith",          //": Checks if the value ends with a specific string.",
	"validateByteRange": "@validateByteRange", //": Checks if the value contains only characters within the specified byte range.",
}
