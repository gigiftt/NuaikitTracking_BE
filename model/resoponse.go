package model

type GetCourseDetail struct {
	Ok           bool           `json:"ok"`
	CourseDetail []CourseDetail `json:"courseDetails"`
}
type CourseDetail struct {
	CourseNo              string                       `json:"courseNo"`
	UpdatedYear           int                          `json:"updatedYear"`
	UpdatedSemester       int                          `json:"updatedSemester"`
	CourseNameEN          string                       `json:"courseNameEN"`
	CourseNameTH          string                       `json:"courseNameTH"`
	CurCodeEN             string                       `json:"curCodeEN"`
	CurCodeTH             string                       `json:"curCodeTH"`
	DetailEN              string                       `json:"detailEN"`
	DetailTH              string                       `json:"detailTH"`
	Credits               CreditDetail                 `json:"credits"`
	SelectedTopicSubjects []SelectedTopicSubjectDetail `json:"selectedTopicSubjects"`
}

type CreditDetail struct {
	Credits   int `json:"credits"`
	Lecture   int `json:"lecture"`
	Practice  int `json:"practice"`
	SelfStudy int `json:"selfStudy"`
}

type SelectedTopicSubjectDetail struct {
	SubjectId    string `json:"subjectId"`
	SubjectTitle string `json:"subjectTitle"`
	IsActive     bool   `json:"isActive"`
}

type CourseDetailResponse struct {
	CourseNo      string   `json:"courseNo"`
	Credits       int      `json:"credits"`
	GroupName     string   `json:"groupName"`
	Prerequisites []string `json:"prerequisites"`
	Corequisite   string   `json:"corequisite"`
	IsPass        bool     `json:"isPass"`
	X             int      `json:"x,omitempty"`
	Y             int      `json:"y,omitempty"`
}

type CategoryResponse struct {
	SummaryCredits  int              `json:"summaryCredit"`
	RequiredCredits int              `json:"requiredCredits"`
	CoreCategory    []CategoryDetail `json:"coreCategory"`
	MajorCategory   []CategoryDetail `json:"majorCategory"`
	GECategory      []CategoryDetail `json:"geCategory"`
	FreeCategory    []CategoryDetail `json:"freeCategory"`
}

type CategoryResponseV2 struct {
	SummaryCredits  int                `json:"summaryCredit"`
	RequiredCredits int                `json:"requiredCredits"`
	CoreCategory    []CategoryDetailV2 `json:"coreCategory"`
	MajorCategory   []CategoryDetailV2 `json:"majorCategory"`
	GECategory      []CategoryDetailV2 `json:"geCategory"`
	FreeCategory    []CategoryDetailV2 `json:"freeCategory"`
}

type CategoryDetail struct {
	GroupName           string                 `json:"groupName"`
	RequiredCreditsNeed int                    `json:"requiredCreditsNeed"`
	RequiredCreditsGet  int                    `json:"requiredCreditsGet"`
	ElectiveCreditsNeed int                    `json:"electiveCreditsNeed"`
	ElectiveCreditsGet  int                    `json:"electiveCreditsGet"`
	RequiredCourseList  []CourseDetailResponse `json:"requiredCourseList"`
	ElectiveCourseList  []CourseDetailResponse `json:"electiveCourseList"`
}

type CategoryDetailV2 struct {
	GroupName           string                 `json:"groupName"`
	RequiredCreditsNeed int                    `json:"requiredCreditsNeed"`
	RequiredCreditsGet  int                    `json:"requiredCreditsGet"`
	ElectiveCreditsNeed int                    `json:"electiveCreditsNeed"`
	ElectiveCreditsGet  int                    `json:"electiveCreditsGet"`
	RequiredCourseList  []CourseDetailResponse `json:"requiredCourseList,omitempty"`
	ElectiveCourseList  []CourseDetailResponse `json:"electiveCourseList,omitempty"`
}

type TermResponse struct {
	Year    string    `json:"year"`
	Plan    string    `json:"plan"`
	Details []Details `json:"details"`
}

type Details struct {
	StudyYear        int                `json:"studyYear"`
	StudyYearDetails []StudyYearDetails `json:"studyYearDetails"`
}

type StudyYearDetails struct {
	StudySemester int                    `json:"studySemester"`
	SummaryCredit int                    `json:"summaryCredit"`
	CourseList    []CourseDetailResponse `json:"courseList"`
}

type Credits struct {
	CoreCredits  int
	MajorCredits int
	GeCredits    int
	FreeCredits  int
}
