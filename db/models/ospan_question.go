package models

/*
CREATE TABLE OSPANQuestion(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                  Question TEXT, OPTION TEXT, isTrue INTEGER);
*/
type OSPANQuestion struct {
	ID       int
	Question string
	Option   string
	IsTrue   int
}
