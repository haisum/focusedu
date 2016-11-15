/*
User
	- ID
	- Name
	- Age
	- RollNo
	- Gender
	- RegisteredAt
	- MidtermScore
	- CurrentStep

	OSPANScore      int
	TotalCorrect    int
	SpeedErrors     int
	AccuracyErrors  int
	MathErrors      int
	FOREIGN KEY(trackartist) REFERENCES artist(artistid)
*/
CREATE TABLE User(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                         Name TEXT DEFAULT "", Age INTEGER DEFAULT 0, RollNo TEXT UNIQUE, 
                                         Gender INTEGER DEFAULT 0, RegisteredAt INTEGER DEFAULT 0,
                                         MidtermScore INTEGER DEFAULT 0, CurrentScore INTEGER DEFAULT 0,
                                         CurrentStep INTEGER DEFAULT 0, QuestionTimeout INTEGER DEFAULT 0,
                                         OSPANScore INTEGER DEFAULT 0, TotalCorrect INTEGER DEFAULT 0,
                                         SpeedErrors INTEGER DEFAULT 0, AccuracyErrors INTEGER DEFAULT 0,
                                         MathErrors INTEGER DEFAULT 0, Type INTEGER DEFAULT 1,
                                         ModuleOneDistractionCount INTEGER DEFAULT 0,
                                         ModuleOneExampleCount INTEGER DEFAULT 0,
                                         ModuleOneGraspingCount INTEGER DEFAULT 0,
                                         ModuleOneCorrect INTEGER DEFAULT 0,
                                         ModuleTwoDistractionCount INTEGER DEFAULT 0,
                                         ModuleTwoExampleCount INTEGER DEFAULT 0,
                                         ModuleTwoGraspingCount INTEGER DEFAULT 0,
                                         ModuleTwoCorrect INTEGER DEFAULT 0,
                                         UsedQuestions TEXT DEFAULT "");
CREATE INDEX User_RollNo ON USER (RollNo);

/*
OSPANQuestion - Stores questions and there answers for ospan
	- ID
	- Question
	- Option
	- isTrue
*/
CREATE TABLE OSPANQuestion(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                  Question TEXT DEFAULT "", Option TEXT DEFAULT "", IsTrue INTEGER DEFAULT 0);

INSERT INTO OSPANQuestion (Question, Option, IsTrue) VALUES
	("(2*9) – 9", "9", 1),
	("(4*2) + 6", "14", 1),
	("(1/1) - 1", "2", 0),
	("(8*2) - 8", "2", 0),
	("(7*3) – 3", "18", 1),
	("(2*9) – 9", "9", 1),
	("(8/2) + 9", "7", 0),
	("(4*5) – 5 ", "11", 0),
	("(1*2) + 1", "3", 1),
	("(3/3) + 2 ", "1", 1),
	("(2*6)- 4 ", "6", 0),
	("(4*3) - 4", "16", 0),
	("(4/4) + 7", "12", 0),
	("(3*8) - 1 ", "23", 1),
	("(8*9) - 8 ", "64", 1),
	("(6/3) +1 ", "3", 1),
	("(1*6) – 3", "9", 0),
	("(8*2) + 5 ", "26", 0),
	("(9/9) – 1 ", "0", 1),
	("(3*4) + 6 ", "20", 0),
	("(7*2) + 1 ", "15", 1),
	("(3/3) +1 ", "2", 1),
	("(8/4) – 2 ", "0", 1),
	("(1*9) – 6 ", "3", 1),
	("(9/3) – 1 ", "1", 0),
	("(2/2) +1 ", "0", 0),
	("(5*3) – 8", "0", 0),
	("(2*6) + 4 ", "16", 1),
	("(1/1) +7 ", "8", 1),
	("(3*2) +3 ", "12", 0),
	("(2*3) + 2", "8", 1),
	("(2*2) – 1", "3", 0),
	("(6/2) + 9", "12", 1),
	("(3/1) +9", "3", 0),
	("(1*3) -1  ", "0", 0),
	("(1*7) - 6", "2", 0),
	("(5*2) – 9", "8", 0),
	("(7/1) - 4 ", "3", 1),
	("(2*5) – 7", "2", 0),
	("(8/2) -1", "11", 0),
	("(2*7) +2", "16", 1),
	("(1*2) + 5", "7", 1),
	("(5/5) + 6", "1", 0),
	("(1*8) + 1", "4", 1),
	("(1*2) + 5", "7", 1),
	("(4/4) + 3", "4", 1),
	("(3*3) + 2", "20", 0);
	

/*
Feedback
	- ID
	- UserID
	- Content
*/
CREATE TABLE Feedback(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                             UserID INTEGER, Question INTEGER, Answer INTEGER
                                             Comments TEXT, FOREIGN KEY(UserID) REFERENCES User(ID));
CREATE INDEX Feedback_UserID ON Feedback (UserID);