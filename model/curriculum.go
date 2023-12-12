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
