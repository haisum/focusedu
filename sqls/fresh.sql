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
	FOREIGN KEY(trackartist) REFERENCES artist(artistid)
*/
CREATE TABLE User(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                         Name TEXT DEFAULT "", Age INTEGER DEFAULT 0, RollNo TEXT UNIQUE , Gender INTEGER DEFAULT 0, RegisteredAt INTEGER DEFAULT 0, MidtermScore INTEGER DEFAULT 0, CurrentScore INTEGER DEFAULT 0, CurrentStep INTEGER DEFAULT 0, QuestionTimeout INTEGER DEFAULT 0);
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

/*
DemoResult
	- ID
	- UserID
	- AverageTime
	- Questions
*/
CREATE TABLE DemoResult(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                               UserID INTEGER, AverageTime INTEGER, Questions TEXT,
                                               FOREIGN KEY(UserID) REFERENCES User(ID));
CREATE INDEX DemoResult_UserID ON DemoResult (UserID);

/*
OSPANResult
	- ID
	- UserID
	- Score
	- Accuracy
	- TimedOut
*/
CREATE TABLE OSPANResult(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                UserID INTEGER, Score INTEGER, Accuracy INTEGER, Timeout INTEGER,
                                                FOREIGN KEY(UserID) REFERENCES User(ID));
CREATE INDEX OSPANResult_UserID ON OSPANResult (UserID);

/*
Module
	- ID
	- Content
	- Example
	- Timeout
*/
CREATE TABLE Module(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                           Content TEXT, Example TEXT, Timeout INTEGER);

/*
ModuleReading
	- ID
	- UserID
	- TotalTimeTaken
	- ModuleID
*/
CREATE TABLE ModuleReading(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                  UserID INTEGER, TotalTimeTaken INTEGER, ModuleID INTEGER,
                                                  FOREIGN KEY(UserID) REFERENCES User(ID),
                                                  FOREIGN KEY(ModuleID) REFERENCES ModuleID(ID));
CREATE INDEX ModuleReading_UserID ON ModuleReading (UserID);

/*
ModuleQuestion
	- ID
	- Question
	- Options
	- Answer
	- ModuleID
*/
CREATE TABLE ModuleQuestion(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                   Question TEXT, OPTIONS TEXT, Answer TEXT, ModuleID INTEGER,
                                                   FOREIGN KEY(ModuleID) REFERENCES ModuleID(ID));
CREATE INDEX ModuleQuestion_ModuleID ON ModuleQuestion (ModuleID);

/*
ModuleResult
	- ID
	- ModuleQuestionID
	- Score
	- GivenAnswer
	- UserID
	- ModuleID
*/
CREATE TABLE ModuleResult(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                                 ModuleQuestionID INTEGER, Score INTEGER, GivenAnswer TEXT, UserID INTEGER, ModuleID INTEGER, FOREIGN KEY(UserID) REFERENCES User(ID),
                                                 FOREIGN KEY(ModuleID) REFERENCES ModuleID(ID),
                                                 FOREIGN KEY(ModuleQuestionID) REFERENCES ModuleQuestion(ID));
CREATE INDEX ModuleResult_UserID ON ModuleResult (UserID);
CREATE INDEX ModuleResult_ModuleID ON ModuleResult (ModuleID);
CREATE INDEX ModuleResult_ModuleQuestionID ON ModuleResult (UserID);

/*
Feedback
	- ID
	- UserID
	- Content
*/
CREATE TABLE Feedback(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                             UserID INTEGER, Content TEXT, FOREIGN KEY(UserID) REFERENCES User(ID));
CREATE INDEX Feedback_UserID ON Feedback (UserID);

/*
Settings (ospan.total.sets.given=5, ospan.demo.sets.given=3)
	- ID
	- Name
	- Value
*/
CREATE TABLE Settings(ID INTEGER PRIMARY KEY AUTOINCREMENT,
                                             Name TEXT, Value TEXT);
INSERT INTO Settings(Name, Value) VALUES ("ospan.total.sets.given", "5");
INSERT INTO Settings(Name, Value) VALUES ("ospan.demo.sets.given", "3");
CREATE INDEX Settings_Name ON Settings (Name);

