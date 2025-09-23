package sql

//go:generate stringer -type=Keyword -linecomment=true

type Keyword int

const (
	Select            Keyword = iota // SELECT
	Insert                           // INSERT
	Update                           // UPDATE
	Delete                           // DELETE
	From                             // FROM
	Where                            // WHERE
	Set                              // SET
	Values                           // VALUES
	OrderBy                          // ORDER BY
	Asc                              // ASC
	Desc                             // DESC
	Limit                            // LIMIT
	Offset                           // OFFSET
	GroupBy                          // GROUP BY
	Having                           // HAVING
	Join                             // JOIN
	LeftJoin                         // LEFT JOIN
	RightJoin                        // RIGHT JOIN
	FullJoin                         // FULL JOIN
	CrossJoin                        // CROSS JOIN
	On                               // ON
	As                               // AS
	InnerJoin                        // INNER JOIN
	OuterJoin                        // OUTER JOIN
	Union                            // UNION
	Intersect                        // INTERSECT
	Except                           // EXCEPT
	UnionAll                         // UNION ALL
	IntersectAll                     // INTERSECT ALL
	ExceptAll                        // EXCEPT ALL
	Not                              // NOT
	In                               // IN
	Between                          // BETWEEN
	And                              // AND
	Or                               // OR
	IsNull                           // IS NULL
	IsNotNull                        // IS NOT NULL
	IsTrue                           // IS TRUE
	IsFalse                          // IS FALSE
	IsUnknown                        // IS UNKNOWN
	IsDistinctFrom                   // IS DISTINCT FROM
	IsNotDistinctFrom                // IS NOT DISTINCT FROM
	Like                             // LIKE
	NotLike                          // NOT LIKE
	Ilike                            // ILIKE
	NotIlike                         // NOT ILIKE
	Any                              // ANY
	All                              // ALL
	Exists                           // EXISTS
	Some                             // SOME
	Unique                           // UNIQUE
	PrimaryKey                       // PRIMARY KEY
	ForeignKey                       // FOREIGN KEY
	Check                            // CHECK
	Default                          // DEFAULT
	Null                             // NULL
	True                             // TRUE
	False                            // FALSE
	Unknown                          // UNKNOWN

	Call    // CALL
	Exec    // EXEC
	Execute // EXECUTE
	Pragma  // PRAGMA
)
