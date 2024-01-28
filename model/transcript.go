package model

type CourseGrade struct {
	Ok           bool     `json:"ok"`
	CourseGrades []Course `json:"courseGrades"`
}

type Course struct {
	CourseNo     string `json:"courseNo"`
	Grade        string `json:"grade"`
	Semester     int    `json:"semester"`
	Year         int    `json:"year"`
	CourseNameEN string `json:"courseNameEN"`
	CourseNameTH string `json:"courseNameTH"`
}

type StudentTranscript struct {
	StudentId  string           `json:"studenteId"`
	Curriculum string           `json:"curriculum"`
	Transcript []TranscriptYear `json:"transcript"`
}

type TranscriptYear struct {
	Year        int                  `json:"year"`
	YearDetails []TranscriptSemester `json:"yearDetails"`
}

type TranscriptSemester struct {
	Semester int                `json:"semester"`
	Details  []TranscriptCourse `json:"details"`
}

type TranscriptCourse struct {
	Code   string `json:"code"`
	Credit int    `json:"credit"`
	Grade  string `json:"grade"`
}
