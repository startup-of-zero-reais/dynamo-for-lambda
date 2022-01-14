package mocks

type Mocktable struct {
	PK           string `diinamo:"type:string;hash"`
	SK           string `diinamo:"type:string;range"`
	Owner        string `diinamo:"type:string;gsi:CourseOwnerIndex;keyPairs:PK=Owner"`
	Title        string `diinamo:"type:string;gsi:CourseTitleIndex;keyPairs:Title=SK"`
	ParentCourse string `diinamo:"type:string;gsi:CourseLessonsIndex;keyPairs:ParentCourse=SK"`
	ParentModule string `diinamo:"type:string;lsi:ModuleLessonsIndex;keyPairs:ParentModule=SK"`
}
