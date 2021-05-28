package src_v2

// This package is strip down version of old code base with removed support for Go, Java,
//	CockroachDB, etc. Only Rust and MySQL is supported in order to keep abstraction closer
//  to our needs and have a smaller source codebase in order to better reason about and
// 	improve.

// Notes:
//	+ Mysql extractor and output types has been seperated as they address different
//		domains and each has it's own edge cases.

/*
	Contemplation:
	+ Should we assume null to default in rust or support null?

*/
