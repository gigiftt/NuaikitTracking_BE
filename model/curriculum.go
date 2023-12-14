package model

type CourseDetail struct {
	CourseNo          string   `json:"courseNo"`
	RecommendSemester int      `json:"recommendSemester"`
	RecommendYear     int      `json:"recommendYear"`
	Prerequisites     []string `json:"prerequisites"`
	Corequisite       string   `json:"corequisite"`
	Credits           int      `json:"credits"`
}

type GroupDetails struct {
	RequiredCredits int            `json:"requiredCredits"`
	GroupName       string         `json:"groupName"`
	RequiredCourses []CourseDetail `json:"requiredCourses"`
	ElectiveCourses []CourseDetail `json:"electiveCourses"`
}

type Curriculum struct {
	Program             string         `json:"curriculumProgram"`
	Year                int            `json:"year"`
	IsCOOPPlan          bool           `json:"isCOOPPlan"`
	RequiredCredits     int            `json:"requiredCredits"`
	FreeElectiveCredits int            `json:"freeElectiveCredits"`
	CoreAndMajorGroups  []GroupDetails `json:"coreAndMajorGroups"`
	GeGroups            []GroupDetails `json:"geGroups"`
}

type CurriculumModel struct {
	Ok         bool       `json:"ok"`
	Curriculum Curriculum `json:"curriculum"`
}

type CourseDetailResponse struct {
	CourseNo   string `json:"courseNo"`
	CourseName string `json:"courseName"`
	GroupName  string `json:"groupName"`
	IsPass     bool   `json:"isPass"`
	X          int    `json:"x,omitempty"`
	Y          int    `json:"y,omitempty"`
}

type CategoryDetail struct {
	SummaryCredits  int                    `json:"summaryCredit"`
	RequiredCredits int                    `json:"requiredCredits"`
	CourseList      []CourseDetailResponse `json:"courseList"`
}

type CategoryResponse struct {
	SummaryCredits  int            `json:"summaryCredit"`
	RequiredCredits int            `json:"requiredCredits"`
	CoreCategory    CategoryDetail `json:"coreCategory"`
	MajorCategory   CategoryDetail `json:"majorCategory"`
	GECategory      CategoryDetail `json:"geCategory"`
	FreeCategory    CategoryDetail `json:"freeCategory"`
}

type Credits struct {
	CoreCredits  int
	MajorCredits int
	GeCredits    int
	FreeCredits  int
}
