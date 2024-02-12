package model

type CurriculumCourseDetail struct {
	CourseNo          string   `json:"courseNo"`
	RecommendSemester int      `json:"recommendSemester"`
	RecommendYear     int      `json:"recommendYear"`
	Prerequisites     []string `json:"prerequisites"`
	Corequisite       string   `json:"corequisite"`
	Credits           int      `json:"credits"`
}

type CurriculumCourseDetail2 struct {
	CourseNo          string   `json:"courseNo"`
	RecommendSemester int      `json:"recommendSemester"`
	RecommendYear     int      `json:"recommendYear"`
	Prerequisites     []string `json:"prerequisites"`
	Corequisite       string   `json:"corequisite"`
	Credits           int      `json:"credits"`
	GroupName         string   `json:"groupName"`
}

type CurriculumGroupDetails struct {
	RequiredCredits int                      `json:"requiredCredits"`
	GroupName       string                   `json:"groupName"`
	RequiredCourses []CurriculumCourseDetail `json:"requiredCourses"`
	ElectiveCourses []CurriculumCourseDetail `json:"electiveCourses"`
}

type Curriculum struct {
	Program             string                   `json:"curriculumProgram"`
	Year                int                      `json:"year"`
	IsCOOPPlan          bool                     `json:"isCOOPPlan"`
	RequiredCredits     int                      `json:"requiredCredits"`
	FreeElectiveCredits int                      `json:"freeElectiveCredits"`
	CoreAndMajorGroups  []CurriculumGroupDetails `json:"coreAndMajorGroups"`
	GeGroups            []CurriculumGroupDetails `json:"geGroups"`
}

type CurriculumModel struct {
	Ok         bool       `json:"ok"`
	Curriculum Curriculum `json:"curriculum"`
}

type PreReqListInfo struct {
	Col              int
	Row              int
	HavePreReq       bool
	Move             bool
	PreReqCourseList []CourseList
}
