package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"NuaikitTracking_BE.com/model"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
)

var PASS_GRADE = []string{"A", "B", "C", "D", "B+", "C+", "D+", "S", "CX"}
var COOPcourse = "261495"

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func readMockData(mockFile string) string {
	jsonFile, err := os.Open(mockFile + ".json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalln(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	c, error := ioutil.ReadAll(jsonFile)
	if error != nil {
		log.Fatalln("Error is : ", err)
	}

	return string(c)
}

func getCirriculum(year string, curriculumProgram string, isCOOP string) (string, error) {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/curriculum"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return "", err
	}

	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("year", year)
	q.Add("curriculumProgram", curriculumProgram)
	q.Add("isCOOPPlan", isCOOP)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error is : ", err)

		return "", err
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
		return "", err
	}

	return string(c), nil
}

func getCirriculumJSON(year string, curriculumProgram string, isCOOP string) (model.CurriculumModel, error) {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/curriculum"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.CurriculumModel{}, err
	}

	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("year", year)
	q.Add("curriculumProgram", curriculumProgram)
	q.Add("isCOOPPlan", isCOOP)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.CurriculumModel{}, err
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
		return model.CurriculumModel{}, err
	}

	curriculum := model.CurriculumModel{}
	err = json.Unmarshal(c, &curriculum)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.CurriculumModel{}, err
	}

	return curriculum, nil

}

func getTermDetail(year string, curriculumProgram string, isCOOP string, studyYear string, studySemester string) (string, model.CurriculumModel, error) {

	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/curriculum"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return "", model.CurriculumModel{}, err
	}

	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("year", year)
	q.Add("curriculumProgram", curriculumProgram)
	q.Add("isCOOPPlan", isCOOP)
	q.Add("studyYear", studyYear)
	q.Add("studySemester", studySemester)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return "", model.CurriculumModel{}, err
	}

	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return "", model.CurriculumModel{}, err
	}

	term := model.CurriculumModel{}
	err = json.Unmarshal(c, &term)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return "", model.CurriculumModel{}, err
	}

	return string(c), term, nil

}

func getCourseDetail(courseNo string) (model.GetCourseDetail, error) {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/course/detail"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.GetCourseDetail{}, err
	}

	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("courseNo", courseNo)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.GetCourseDetail{}, err
	}

	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.GetCourseDetail{}, err
	}

	detail := model.GetCourseDetail{}

	err = json.Unmarshal(c, &detail)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.GetCourseDetail{}, err
	}

	return detail, nil

}

func getRawTranscript(studentId string) (model.CourseGrade, error) {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/private/student/" + studentId + "/courseGrade"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.CourseGrade{}, err
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
		return model.CourseGrade{}, error
	}

	courseGrade := model.CourseGrade{}

	err = json.Unmarshal(c, &courseGrade)
	if err != nil {
		log.Fatalln("Error is : ", err)
		return model.CourseGrade{}, err
	}

	return courseGrade, nil
}

func getTranscript(studentId string) model.StudentTranscript {

	split := strings.SplitAfter(studentId, "")
	idText := "25" + split[0] + split[1]
	idNum, _ := strconv.Atoi(idText)

	rawTranscript, _ := getRawTranscript(studentId)
	if !rawTranscript.Ok {
		return model.StudentTranscript{
			Status:     false,
			StudentId:  studentId,
			Transcript: []model.TranscriptYear{},
		}
	}

	courseGrade := map[int]map[int][]model.TranscriptCourse{}

	for _, c := range rawTranscript.CourseGrades {
		semesterList, b := courseGrade[c.Year]
		if b {
			courseList, b := semesterList[c.Semester]
			if b {
				courseList = append(courseList, model.TranscriptCourse{
					Code:  c.CourseNo,
					Grade: c.Grade,
				})

				courseGrade[c.Year][c.Semester] = courseList

			} else {
				courses := []model.TranscriptCourse{}
				courses = append(courses, model.TranscriptCourse{
					Code:  c.CourseNo,
					Grade: c.Grade,
				})

				semesterList[c.Semester] = courses
				courseGrade[c.Year] = semesterList

			}

		} else {
			courses := []model.TranscriptCourse{}
			courses = append(courses, model.TranscriptCourse{
				Code:  c.CourseNo,
				Grade: c.Grade,
			})
			semesterList = map[int][]model.TranscriptCourse{}
			semesterList[c.Semester] = courses
			courseGrade[c.Year] = semesterList
		}
	}

	transcriptYear := []model.TranscriptYear{}

	for i := 0; i < len(courseGrade); i++ {
		transcriptSemester := []model.TranscriptSemester{}

		_, b := courseGrade[idNum]
		if b {

			for j := 1; j <= len(courseGrade[idNum]); j++ {
				log.Println("id Num : ", idNum)
				log.Println("j: ", j)
				detail, b := courseGrade[idNum][j]
				if b {
					semesterDetail := model.TranscriptSemester{
						Semester: j,
						Details:  detail,
					}
					transcriptSemester = append(transcriptSemester, semesterDetail)
				} else {
					semesterDetail := model.TranscriptSemester{
						Semester: j,
						Details:  []model.TranscriptCourse{},
					}
					transcriptSemester = append(transcriptSemester, semesterDetail)
				}
			}

			yearDetail := model.TranscriptYear{
				Year:        idNum,
				YearDetails: transcriptSemester,
			}
			idNum++

			transcriptYear = append(transcriptYear, yearDetail)
		} else {
			yearDetail := model.TranscriptYear{
				Year:        idNum,
				YearDetails: transcriptSemester,
			}
			transcriptYear = append(transcriptYear, yearDetail)
		}

	}

	transcriptFinal := model.StudentTranscript{
		Status:     true,
		StudentId:  studentId,
		Transcript: transcriptYear,
	}

	log.Println(transcriptFinal)

	return transcriptFinal
}

func getTranscriptWithCredit(studentId string) model.StudentTranscript {

	split := strings.SplitAfter(studentId, "")
	idText := "25" + split[0] + split[1]
	idNum, _ := strconv.Atoi(idText)

	rawTranscript, _ := getRawTranscript(studentId)
	if !rawTranscript.Ok {
		return model.StudentTranscript{
			Status:     false,
			StudentId:  studentId,
			Transcript: []model.TranscriptYear{},
		}
	}

	courseGrade := map[int]map[int][]model.TranscriptCourse{}

	for _, c := range rawTranscript.CourseGrades {

		detail, err := getCourseDetail(c.CourseNo)
		if err != nil {
			log.Fatalln("Error is : ", err)
		}

		credit := 3
		if len(detail.CourseDetail) != 0 {
			credit = detail.CourseDetail[0].Credits.Credits
		}

		semesterList, b := courseGrade[c.Year]
		if b {
			courseList, b := semesterList[c.Semester]
			if b {
				courseList = append(courseList, model.TranscriptCourse{
					Code:   c.CourseNo,
					Credit: credit,
					Grade:  c.Grade,
				})

				courseGrade[c.Year][c.Semester] = courseList

			} else {
				courses := []model.TranscriptCourse{}
				courses = append(courses, model.TranscriptCourse{
					Code:   c.CourseNo,
					Credit: credit,
					Grade:  c.Grade,
				})

				semesterList[c.Semester] = courses
				courseGrade[c.Year] = semesterList

			}

		} else {
			courses := []model.TranscriptCourse{}
			courses = append(courses, model.TranscriptCourse{
				Code:   c.CourseNo,
				Credit: credit,
				Grade:  c.Grade,
			})
			semesterList = map[int][]model.TranscriptCourse{}
			semesterList[c.Semester] = courses
			courseGrade[c.Year] = semesterList
		}
	}

	transcriptYear := []model.TranscriptYear{}

	for i := 0; i < len(courseGrade); i++ {
		transcriptSemester := []model.TranscriptSemester{}

		_, b := courseGrade[idNum]
		if b {

			for j := 1; j <= len(courseGrade[idNum]); j++ {
				log.Println("id Num : ", idNum)
				log.Println("j: ", j)
				detail, b := courseGrade[idNum][j]
				if b {
					semesterDetail := model.TranscriptSemester{
						Semester: j,
						Details:  detail,
					}
					transcriptSemester = append(transcriptSemester, semesterDetail)
				} else {
					semesterDetail := model.TranscriptSemester{
						Semester: j,
						Details:  []model.TranscriptCourse{},
					}
					transcriptSemester = append(transcriptSemester, semesterDetail)
				}
			}

			yearDetail := model.TranscriptYear{
				Year:        idNum,
				YearDetails: transcriptSemester,
			}
			idNum++

			transcriptYear = append(transcriptYear, yearDetail)
		} else {
			yearDetail := model.TranscriptYear{
				Year:        idNum,
				YearDetails: transcriptSemester,
			}
			transcriptYear = append(transcriptYear, yearDetail)
		}

	}

	transcriptFinal := model.StudentTranscript{
		Status:     true,
		StudentId:  studentId,
		Transcript: transcriptYear,
	}

	return transcriptFinal
}

func checkGroup(cirriculum string, courseNo string) (string, string) {

	groupList := gjson.Get(cirriculum, "curriculum.geGroups.#.groupName")
	for _, groupName := range groupList.Array() {

		queryReqCourse := `curriculum.geGroups.#(groupName=="` + groupName.String() + `").requiredCourses.#(courseNo=="` + courseNo + `")`
		valueReqCourse := gjson.Get(cirriculum, queryReqCourse)

		if valueReqCourse.Exists() {
			return groupName.String(), "requiredCourses"
		}

		queryElecCourse := `curriculum.geGroups.#(groupName=="` + groupName.String() + `").electiveCourses.#(courseNo=="` + courseNo + `")`
		valueElecCourse := gjson.Get(cirriculum, queryElecCourse)

		if valueElecCourse.Exists() {
			return groupName.String(), "electiveCourses"
		}
	}

	groupList = gjson.Get(cirriculum, "curriculum.coreAndMajorGroups.#.groupName")
	for _, groupName := range groupList.Array() {

		queryReqCourse := `curriculum.coreAndMajorGroups.#(groupName=="` + groupName.String() + `").requiredCourses.#(courseNo=="` + courseNo + `")`
		valueReqCourse := gjson.Get(cirriculum, queryReqCourse)

		if valueReqCourse.Exists() {
			return groupName.String(), "requiredCourses"
		}

		queryElecCourse := `curriculum.coreAndMajorGroups.#(groupName=="` + groupName.String() + `").electiveCourses.#(courseNo=="` + courseNo + `")`
		valueElecCourse := gjson.Get(cirriculum, queryElecCourse)

		if valueElecCourse.Exists() {
			return groupName.String(), "electiveCourses"
		}
	}

	return "Free", "electiveCourses"
}

func getSummaryCredits(c model.CurriculumModel, curriculumString string, isCOOP string, studentId string, mockData string) (model.CategoryResponseV2, error) {

	transcript := ""
	transcriptModel := model.StudentTranscript{}
	if studentId == "" {
		transcript = readMockData(mockData)

	} else {
		transcriptModel = getTranscriptWithCredit(studentId)

		tm, err := json.Marshal(transcriptModel)
		if err != nil {
			log.Fatalln("Error is : ", err)
		}

		transcript = string(tm)

		if !transcriptModel.Status {
			transcript = ""
		}
	}

	if strings.Contains(transcript, COOPcourse) {
		isCOOP = "true"
	}

	t := model.CategoryResponseV2{}
	curriculumRequiredCredits := c.Curriculum.RequiredCredits
	freeRequiredCredits := c.Curriculum.FreeElectiveCredits
	freeCategory := []model.CategoryDetailV2{}
	freeCategory = append(freeCategory, model.CategoryDetailV2{
		GroupName:           "Free Elective",
		RequiredCreditsNeed: 0,
		RequiredCreditsGet:  0,
		ElectiveCreditsNeed: freeRequiredCredits,
		ElectiveCreditsGet:  0,
	})
	coreCategory := []model.CategoryDetailV2{}
	majorCategory := []model.CategoryDetailV2{}
	geCategory := []model.CategoryDetailV2{}
	coopCourse := model.CourseDetailResponse{}

	//core and major template
	for _, g := range c.Curriculum.CoreAndMajorGroups {

		groupName := g.GroupName
		reqCredit := 0

		for _, c := range g.RequiredCourses {
			reqCredit += c.Credits
		}

		if isCOOP == "true" && groupName == "Major Required" {
			c := gjson.Get(curriculumString, `curriculum.coreAndMajorGroups.#(groupName="Major Elective").electiveCourses.#(courseNo="`+COOPcourse+`")`)

			reqCredit += int(c.Get("credits").Int())
			coopCourse = model.CourseDetailResponse{
				CourseNo:  c.Get("courseNo").String(),
				Credits:   int(c.Get("credits").Int()),
				GroupName: "Major Required",
				IsPass:    false,
			}

		}

		if groupName == "Core" {
			coreCategory = append(coreCategory, model.CategoryDetailV2{
				GroupName:           groupName,
				RequiredCreditsNeed: reqCredit,
				RequiredCreditsGet:  0,
				ElectiveCreditsNeed: g.RequiredCredits - reqCredit,
				ElectiveCreditsGet:  0,
			})
		} else if groupName == "Major Elective" {

			if isCOOP == "true" {
				g.RequiredCredits -= coopCourse.Credits
			}
			majorCategory = append(majorCategory, model.CategoryDetailV2{
				GroupName:           groupName,
				RequiredCreditsNeed: reqCredit,
				RequiredCreditsGet:  0,
				ElectiveCreditsNeed: g.RequiredCredits,
				ElectiveCreditsGet:  0,
			})
		} else {

			majorCategory = append(majorCategory, model.CategoryDetailV2{
				GroupName:           groupName,
				RequiredCreditsNeed: reqCredit,
				RequiredCreditsGet:  0,
				ElectiveCreditsNeed: 0,
				ElectiveCreditsGet:  0,
			})
		}
	}

	//ge template
	for _, g := range c.Curriculum.GeGroups {

		groupName := g.GroupName
		reqCredit := 0

		for _, c := range g.RequiredCourses {

			reqCredit += c.Credits
		}

		geCategory = append(geCategory, model.CategoryDetailV2{
			GroupName:           groupName,
			RequiredCreditsNeed: reqCredit,
			RequiredCreditsGet:  0,
			ElectiveCreditsNeed: g.RequiredCredits - reqCredit,
			ElectiveCreditsGet:  0,
		})

	}

	t = model.CategoryResponseV2{
		SummaryCredits:  0,
		RequiredCredits: curriculumRequiredCredits,
		CoreCategory:    coreCategory,
		MajorCategory:   majorCategory,
		GECategory:      geCategory,
		FreeCategory:    freeCategory,
	}

	tt, err := json.Marshal(t)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	template := string(tt)

	summaryCredits := 0

	if transcript != "" {
		yearList := gjson.Get(transcript, "transcript.#.year")
		for _, y := range yearList.Array() {
			semester := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#`)
			i := 1
			for i < (int(semester.Int()) + 1) {
				courseList := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#(semester==`+strconv.Itoa(i)+`).details`)
				for _, c := range courseList.Array() {

					code := gjson.Get(c.String(), "code")
					grade := gjson.Get(c.String(), "grade")
					credit := gjson.Get(c.String(), "credit").Int()

					if slices.Contains(PASS_GRADE, grade.String()) {

						group, courseType := checkGroup(curriculumString, code.String())

						if group == "Free" {
							oldCredit := gjson.Get(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`).Int()
							newCredit := oldCredit + credit
							template, _ = sjson.Set(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`, newCredit)

						} else if group == "Core" {
							if courseType == "requiredCourses" {

								oldCredit := gjson.Get(template, `coreCategory.#(groupName="Core").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `coreCategory.#(groupName="Core").requiredCreditsGet`, newCredit)

							} else {

								oldCredit := gjson.Get(template, `coreCategory.#(groupName="Core").electiveCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `coreCategory.#(groupName="Core").electiveCreditsGet`, newCredit)

							}

						} else if group == "Major Required" || group == "Major Elective" {

							if courseType == "requiredCourses" {

								oldCredit := gjson.Get(template, `majorCategory.#(groupName="`+group+`").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `majorCategory.#(groupName="`+group+`").requiredCreditsGet`, newCredit)

							} else if code.String() == COOPcourse {
								oldCredit := gjson.Get(template, `majorCategory.#(groupName="Major Required").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `majorCategory.#(groupName="Major Required").requiredCreditsGet`, newCredit)
							} else {

								oldCredit := gjson.Get(template, `majorCategory.#(groupName="`+group+`").electiveCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `majorCategory.#(groupName="`+group+`").electiveCreditsGet`, newCredit)
							}

						} else {

							if courseType == "requiredCourses" {

								oldCredit := gjson.Get(template, `geCategory.#(groupName="`+group+`").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").requiredCreditsGet`, newCredit)

							} else {

								oldCredit := gjson.Get(template, `geCategory.#(groupName="`+group+`").electiveCreditsGet`).Int()
								requiredCredit := gjson.Get(template, `geCategory.#(groupName="`+group+`").electiveCreditsNeed`).Int()

								if oldCredit >= requiredCredit {

									oldCredit := gjson.Get(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`).Int()
									newCredit := oldCredit + credit
									template, _ = sjson.Set(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`, newCredit)

								} else {

									newCredit := oldCredit + credit
									template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").electiveCreditsGet`, newCredit)

								}

							}
						}
						summaryCredits += int(credit)
					}

				}

				i++
			}
		}
	}

	err = json.Unmarshal([]byte(template), &t)
	if err != nil {
		return t, err
	}

	t.SummaryCredits = summaryCredits
	t.IsCoop = isCOOP

	return t, nil
}

func getCategoryTemplate(c model.CurriculumModel, curriculumString string, isCOOP string, studentId string, mockData string) (string, int, string, error) {

	transcript := ""
	if studentId == "" {
		transcript = readMockData(mockData)

	} else {
		transcriptModel := getTranscriptWithCredit(studentId)

		tm, err := json.Marshal(transcriptModel)
		if err != nil {
			log.Fatalln("Error is : ", err)
		}

		transcript = string(tm)

		if !transcriptModel.Status {
			transcript = ""
		}
	}

	if strings.Contains(transcript, COOPcourse) {
		isCOOP = "true"
	}

	curriculumRequiredCredits := c.Curriculum.RequiredCredits
	freeRequiredCredits := c.Curriculum.FreeElectiveCredits
	freeCategory := []model.CategoryDetail{}
	freeCategory = append(freeCategory, model.CategoryDetail{
		GroupName:           "Free Elective",
		RequiredCreditsNeed: 0,
		RequiredCreditsGet:  0,
		ElectiveCreditsNeed: freeRequiredCredits,
		ElectiveCreditsGet:  0,
		RequiredCourseList:  []model.CourseDetailResponse{},
		ElectiveCourseList:  []model.CourseDetailResponse{},
	})
	coreCategory := []model.CategoryDetail{}
	majorCategory := []model.CategoryDetail{}
	geCategory := []model.CategoryDetail{}

	coopCourse := model.CourseDetailResponse{}

	//core and major template
	for _, g := range c.Curriculum.CoreAndMajorGroups {

		groupName := g.GroupName
		reqCourseList := []model.CourseDetailResponse{}
		reqCredit := 0
		elecCourseList := []model.CourseDetailResponse{}

		for _, c := range g.RequiredCourses {
			reqCourseList = append(reqCourseList, model.CourseDetailResponse{
				CourseNo:  c.CourseNo,
				Credits:   c.Credits,
				GroupName: groupName,
				IsPass:    false,
			})
			reqCredit += c.Credits
		}

		if isCOOP == "true" && groupName == "Major Required" {
			c := gjson.Get(curriculumString, `curriculum.coreAndMajorGroups.#(groupName="Major Elective").electiveCourses.#(courseNo="`+COOPcourse+`")`)

			reqCredit += int(c.Get("credits").Int())
			coopCourse = model.CourseDetailResponse{
				CourseNo:  c.Get("courseNo").String(),
				Credits:   int(c.Get("credits").Int()),
				GroupName: "Major Required",
				IsPass:    false,
			}

			reqCourseList = append(reqCourseList, coopCourse)
		}

		if groupName == "Core" {
			coreCategory = append(coreCategory, model.CategoryDetail{
				GroupName:           groupName,
				RequiredCreditsNeed: reqCredit,
				RequiredCreditsGet:  0,
				ElectiveCreditsNeed: g.RequiredCredits - reqCredit,
				ElectiveCreditsGet:  0,
				RequiredCourseList:  reqCourseList,
				ElectiveCourseList:  elecCourseList,
			})
		} else if groupName == "Major Elective" {
			if isCOOP == "true" {
				g.RequiredCredits -= coopCourse.Credits
			}

			majorCategory = append(majorCategory, model.CategoryDetail{
				GroupName:           groupName,
				RequiredCreditsNeed: reqCredit,
				RequiredCreditsGet:  0,
				ElectiveCreditsNeed: g.RequiredCredits,
				ElectiveCreditsGet:  0,
				RequiredCourseList:  reqCourseList,
				ElectiveCourseList:  elecCourseList,
			})
		} else {

			majorCategory = append(majorCategory, model.CategoryDetail{
				GroupName:           groupName,
				RequiredCreditsNeed: reqCredit,
				RequiredCreditsGet:  0,
				ElectiveCreditsNeed: 0,
				ElectiveCreditsGet:  0,
				RequiredCourseList:  reqCourseList,
				ElectiveCourseList:  elecCourseList,
			})

		}
	}

	//ge template
	for _, g := range c.Curriculum.GeGroups {

		groupName := g.GroupName
		reqCourseList := []model.CourseDetailResponse{}
		reqCredit := 0
		elecCourseList := []model.CourseDetailResponse{}

		for _, c := range g.RequiredCourses {

			// detail := getCourseDetail(c.CourseNo)
			reqCourseList = append(reqCourseList, model.CourseDetailResponse{
				CourseNo:  c.CourseNo,
				Credits:   c.Credits,
				GroupName: groupName,
				IsPass:    false,
			})
			reqCredit += c.Credits
		}

		geCategory = append(geCategory, model.CategoryDetail{
			GroupName:           groupName,
			RequiredCreditsNeed: reqCredit,
			RequiredCreditsGet:  0,
			ElectiveCreditsNeed: g.RequiredCredits - reqCredit,
			ElectiveCreditsGet:  0,
			RequiredCourseList:  reqCourseList,
			ElectiveCourseList:  elecCourseList,
		})

	}

	t := model.CategoryResponse{
		SummaryCredits:  0,
		RequiredCredits: curriculumRequiredCredits,
		CoreCategory:    coreCategory,
		MajorCategory:   majorCategory,
		GECategory:      geCategory,
		FreeCategory:    freeCategory,
	}

	tt, err := json.Marshal(t)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	template := string(tt)

	summaryCredits := 0

	if transcript != "" {

		yearList := gjson.Get(transcript, "transcript.#.year")
		for _, y := range yearList.Array() {
			semester := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#`)
			i := 1
			for i < (int(semester.Int()) + 1) {
				courseList := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#(semester==`+strconv.Itoa(i)+`).details`)
				for _, c := range courseList.Array() {

					code := gjson.Get(c.String(), "code")
					grade := gjson.Get(c.String(), "grade")
					credit := gjson.Get(c.String(), "credit").Int()

					if slices.Contains(PASS_GRADE, grade.String()) {

						group, courseType := checkGroup(curriculumString, code.String())

						courseList := []model.CourseDetailResponse{}

						if group == "Free" {

							oldCredit := gjson.Get(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`).Int()
							newCredit := oldCredit + credit
							template, _ = sjson.Set(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`, newCredit)

							if gjson.Get(template, `freeCategory.#(groupName="Free Elective").electiveCourseList.#`).Int() == 0 {
								courseList = append(courseList, model.CourseDetailResponse{
									CourseNo:  code.String(),
									Credits:   int(credit),
									GroupName: group,
									IsPass:    true,
								})
							} else {

								oldValue := gjson.Get(template, `freeCategory.#(groupName="Free Elective")`).String()
								categoryDetail := model.CategoryDetail{}
								err := json.Unmarshal([]byte(oldValue), &categoryDetail)
								if err != nil {
									return "", 0, "", err
								}

								courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
									CourseNo:  code.String(),
									Credits:   int(credit),
									GroupName: group,
									IsPass:    true,
								})
							}

							template, _ = sjson.Set(template, `freeCategory.#(groupName="Free Elective").electiveCourseList`, courseList)

						} else if group == "Core" {

							if courseType == "requiredCourses" {

								oldCredit := gjson.Get(template, `coreCategory.#(groupName="Core").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `coreCategory.#(groupName="Core").requiredCreditsGet`, newCredit)

								template, _ = sjson.Set(template, `coreCategory.#(groupName="Core").requiredCourseList.#(courseNo="`+code.String()+`").isPass`, true)

							} else {

								oldCredit := gjson.Get(template, `coreCategory.#(groupName="Core").electiveCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `coreCategory.#(groupName="Core").electiveCreditsGet`, newCredit)

								if gjson.Get(template, `coreCategory.#(groupName="Core").electiveCourseList.#`).Int() == 0 {
									courseList = append(courseList, model.CourseDetailResponse{
										CourseNo:  code.String(),
										Credits:   int(credit),
										GroupName: group,
										IsPass:    true,
									})
								} else {

									oldValue := gjson.Get(template, `coreCategory.#(groupName="Core")`).String()
									categoryDetail := model.CategoryDetail{}
									err := json.Unmarshal([]byte(oldValue), &categoryDetail)
									if err != nil {
										return "", 0, "", err
									}

									courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
										CourseNo:  code.String(),
										Credits:   int(credit),
										GroupName: group,
										IsPass:    true,
									})
								}

								template, _ = sjson.Set(template, `coreCategory.#(groupName="Core").electiveCourseList`, courseList)
							}

						} else if group == "Major Required" || group == "Major Elective" {

							if courseType == "requiredCourses" {

								oldCredit := gjson.Get(template, `majorCategory.#(groupName="`+group+`").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `majorCategory.#(groupName="`+group+`").requiredCreditsGet`, newCredit)

								template, _ = sjson.Set(template, `majorCategory.#(groupName="`+group+`").requiredCourseList.#(courseNo="`+code.String()+`").isPass`, true)

							} else if code.String() == COOPcourse {

								oldCredit := gjson.Get(template, `majorCategory.#(groupName="Major Required").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `majorCategory.#(groupName="Major Required").requiredCreditsGet`, newCredit)

								template, _ = sjson.Set(template, `majorCategory.#(groupName="Major Required").requiredCourseList.#(courseNo="`+code.String()+`").isPass`, true)
							} else {

								oldCredit := gjson.Get(template, `majorCategory.#(groupName="`+group+`").electiveCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `majorCategory.#(groupName="`+group+`").electiveCreditsGet`, newCredit)

								if gjson.Get(template, `majorCategory.#(groupName="`+group+`").electiveCourseList.#`).Int() == 0 {
									courseList = append(courseList, model.CourseDetailResponse{
										CourseNo:  code.String(),
										Credits:   int(credit),
										GroupName: group,
										IsPass:    true,
									})
								} else {

									oldValue := gjson.Get(template, `majorCategory.#(groupName="`+group+`")`).String()
									categoryDetail := model.CategoryDetail{}
									err := json.Unmarshal([]byte(oldValue), &categoryDetail)
									if err != nil {
										return "", 0, "", err
									}

									courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
										CourseNo:  code.String(),
										Credits:   int(credit),
										GroupName: group,
										IsPass:    true,
									})
								}

								template, _ = sjson.Set(template, `majorCategory.#(groupName="`+group+`").electiveCourseList`, courseList)
							}

						} else {

							if courseType == "requiredCourses" {

								oldCredit := gjson.Get(template, `geCategory.#(groupName="`+group+`").requiredCreditsGet`).Int()
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").requiredCreditsGet`, newCredit)

								template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").requiredCourseList.#(courseNo="`+code.String()+`").isPass`, true)

							} else {

								oldCredit := gjson.Get(template, `geCategory.#(groupName="`+group+`").electiveCreditsGet`).Int()
								requiredCredit := gjson.Get(template, `geCategory.#(groupName="`+group+`").electiveCreditsNeed`).Int()

								if oldCredit >= requiredCredit {

									oldCredit := gjson.Get(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`).Int()
									newCredit := oldCredit + credit
									template, _ = sjson.Set(template, `freeCategory.#(groupName="Free Elective").electiveCreditsGet`, newCredit)

									if gjson.Get(template, `freeCategory.#(groupName="Free Elective").electiveCourseList.#`).Int() == 0 {
										courseList = append(courseList, model.CourseDetailResponse{
											CourseNo:  code.String(),
											Credits:   int(credit),
											GroupName: group,
											IsPass:    true,
										})
									} else {

										oldValue := gjson.Get(template, `freeCategory.#(groupName="Free Elective")`).String()
										categoryDetail := model.CategoryDetail{}
										err := json.Unmarshal([]byte(oldValue), &categoryDetail)
										if err != nil {
											return "", 0, "", err
										}

										courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
											CourseNo:  code.String(),
											Credits:   int(credit),
											GroupName: group,
											IsPass:    true,
										})
									}

									template, _ = sjson.Set(template, `freeCategory.#(groupName="Free Elective").electiveCourseList`, courseList)
								} else {

									newCredit := oldCredit + credit
									template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").electiveCreditsGet`, newCredit)

									log.Println("elective code.String() : ", code.String())
									if gjson.Get(template, `geCategory.#(groupName="`+group+`").electiveCourseList.#`).Int() == 0 {
										courseList = append(courseList, model.CourseDetailResponse{
											CourseNo:  code.String(),
											Credits:   int(credit),
											GroupName: group,
											IsPass:    true,
										})
									} else {

										oldValue := gjson.Get(template, `geCategory.#(groupName="`+group+`")`).String()
										categoryDetail := model.CategoryDetail{}
										err := json.Unmarshal([]byte(oldValue), &categoryDetail)
										if err != nil {
											return "", 0, "", err
										}

										courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
											CourseNo:  code.String(),
											Credits:   int(credit),
											GroupName: group,
											IsPass:    true,
										})
									}

									template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").electiveCourseList`, courseList)

								}

							}

						}
						summaryCredits += int(credit)
					}

				}

				i++
			}
		}
	}

	return template, summaryCredits, isCOOP, nil

}

func insertRow(template *[][]string, index int, corepList []string) {

	for col := range *template {
		row := index

		if len((*template)[col]) <= index {

			(*template)[col] = append((*template)[col], "000000")
		} else {

			if slices.Contains[[]string](corepList, (*template)[col][index]) {
				row = index + 1
			}
			(*template)[col] = slices.Insert[[]string]((*template)[col], row, "000000")
		}
	}

}

func removeIndex(s *[]string, index int) {
	*s = append((*s)[:index], (*s)[index+1:]...)
}

func putInTemplate(templateArr [][]string, x int, corequisiteList []string, noPreList []string, havePreList []string, haveRequisite map[string][]string, corequisite string, courseNo string, prerequisites []gjson.Result, listOfCourse map[string]*model.CurriculumCourseDetail2) ([][]string, []string, []string, []string, map[string][]string, []string) {

	prerequisitesList := []string{}

	if slices.Contains[[]string](templateArr[x], corequisite) {

		corequisiteList = append(corequisiteList, []string{courseNo}...)
		row := slices.Index[[]string](templateArr[x], corequisite)

		if row+1 < len(templateArr[x]) {
			if templateArr[x][row+1] != "000000" {
				insertRow(&templateArr, row+1, corequisiteList)

			}
		} else {
			insertRow(&templateArr, row+1, corequisiteList)
		}

		templateArr[x][row+1] = courseNo

	} else {

		log.Println(courseNo)

		if len(prerequisites) == 0 {
			// have 0 prerequisites
			noPreList = append(noPreList, courseNo)

		} else if len(prerequisites) == 1 {
			// have 1 prerequisites
			havePreList = append(havePreList, courseNo)
			prerequisitesList = append(prerequisitesList, prerequisites[0].String())
			arr, h := haveRequisite[prerequisites[0].String()]
			if !h {
				haveRequisite[prerequisites[0].String()] = []string{courseNo}
			} else {
				arr = append(arr, courseNo)
				haveRequisite[prerequisites[0].String()] = arr
			}

			preRow := 0
			preCol := 0

			// find position of prerequisites
			for col := range templateArr {
				row := slices.Index[[]string](templateArr[col], prerequisites[0].String())
				if row != -1 {
					preRow = row
					preCol = col
					break
				}
			}

			nowRow := preRow

			if x == preCol+1 {
				// pre อยู่คอลัมก่อนหน้า
				if templateArr[x][preRow] != "000000" {
					if preRow == 0 {
						insertRow(&templateArr, preRow, corequisiteList)

					} else {
						if templateArr[x][preRow-1] == "000000" {
							nowRow = preRow - 1
						} else {
							insertRow(&templateArr, preRow, corequisiteList)

						}
					}

				}
				templateArr[x][nowRow] = courseNo

			} else {

				// pre ไม่ได้อยู่คอลัมก่อนหน้า
				available := true
				for o := preCol + 1; o < x+1; o++ {
					if templateArr[o][preRow] != "000000" {
						available = false
						break
					}
				}

				if !available {

					insertRow(&templateArr, preRow, corequisiteList)

				}
				templateArr[x][nowRow] = courseNo

				for v := preCol + 1; v < x; v++ {
					templateArr[v][nowRow] = "111111"
				}

			}

		} else {
			// เอาวิชานี้ใส่ใน havePreList
			havePreList = append(havePreList, courseNo)

			thisPreList := map[string]model.PreReqListInfo{}

			// course ที่จะเอาตัวนี้มาต่อ
			headCourse := ""

			// prerequisites => prereq ของ course นี้
			for _, c := range prerequisites {

				course := c.String()

				// เพิ่ม pre
				prerequisitesList = append(prerequisitesList, course)

				havePreReq := slices.Contains[[]string](havePreList, course)

				// ถ้าตัวพรีตัวนี้มีตัวพรีตัวก่อนหน้าอีก
				if havePreReq {
					log.Println("course : ", course)
					term, row := checkTermAndIndex(templateArr, course)
					log.Println("term : ", term)
					log.Println("row : ", row)
					b, PreReqCourseList := getAllListCourse(templateArr, course, haveRequisite, listOfCourse, row)
					if !b && headCourse == "" && templateArr[x][row] == "" {
						headCourse = course
					} else if b && headCourse == "" && len(PreReqCourseList) > 0 && templateArr[x][row] == "" {
						headCourse = course
					}

					thisPreList[course] = model.PreReqListInfo{
						Col:              term,
						Row:              row,
						HavePreReq:       havePreReq,
						Move:             b,
						PreReqCourseList: PreReqCourseList,
					}
				} else {

					term, row := checkTermAndIndex(templateArr, course)
					thisPreList[course] = model.PreReqListInfo{
						Col:        term,
						Row:        row,
						Move:       true,
						HavePreReq: havePreReq,
					}
				}
			}

			if headCourse == "" {
				for key, value := range thisPreList {
					if !value.Move {
						headCourse = key
						break
					} else {
						if len(value.PreReqCourseList) > 0 {
							headCourse = key
							break
						} else {
							headCourse = key
						}
					}
				}
			}

			log.Println("headCourse : ", headCourse)
			detail, b := thisPreList[headCourse]
			if b {

				if templateArr[x][detail.Row] == "000000" {
					templateArr[x][detail.Row] = courseNo

					// check ว่าในแถวมี coreq ไหม
					haveCoreq := true
					for o := detail.Col; o < x; o++ {
						if slices.Contains[[]string](corequisiteList, templateArr[o][detail.Row+1]) {
							haveCoreq = false
							break
						}
					}

					log.Println("haveCoreq : ", haveCoreq)

					for key, value := range thisPreList {

						log.Println("key : ", key)

						if !haveCoreq {
							// ในแถวมี coreq จะแทรก 2 แถว
							log.Println("have co course")
							log.Println("key : ", key)
							log.Println("value.Move : ", value.Move)

							if key != headCourse && value.Move {

								insertRow(&templateArr, detail.Row+1, corequisiteList)
								insertRow(&templateArr, detail.Row+1, corequisiteList)

								log.Println("templateArr after insert : ", templateArr)
								farestCol := value.Col
								oldRow := value.Row
								if value.Row > detail.Row {
									oldRow = value.Row + 2
								} else {
									oldRow = value.Row
								}

								// ย้านที่ตัวแรก
								templateArr[value.Col][detail.Row+2] = key

								log.Println("detail : ", detail)
								log.Println("PreReqCourseList : ", value.PreReqCourseList)

								// ถ้าเป็นสายก็ย้ายที่
								for _, d := range value.PreReqCourseList {

									templateArr[d.Term][detail.Row+2] = d.CourseCode

									if farestCol < d.Term {
										farestCol = d.Term
									}

								}

								// update เส้นเชื่อม
								for v := farestCol + 1; v < x; v++ {
									templateArr[v][detail.Row+2] = "111111"
								}

								// ลบก่อนย้ายทิ้ง
								for i := 0; i <= x; i++ {
									removeIndex(&templateArr[i], oldRow)
								}

							} else {

								available := true
								for o := value.Col + 1; o < x+1; o++ {
									if templateArr[o][value.Row] != "000000" {
										available = false
										break
									}
								}

								log.Println("available : ", available)

								if available {
									for v := value.Col + 1; v < x; v++ {
										templateArr[v][detail.Row+2] = "111111"
									}
								}

							}

						} else {

							// ในแถววไม่มี coreq จะแทรก 1 แถว
							if key != headCourse && value.Move {
								insertRow(&templateArr, detail.Row+1, corequisiteList)

								farestCol := value.Col
								oldRow := value.Row

								for _, d := range value.PreReqCourseList {

									templateArr[d.Term][detail.Row+1] = d.CourseCode

									if farestCol < d.Term {
										farestCol = d.Term
									}
									if d.Row > detail.Row {
										oldRow = d.Row + 1
									} else {
										oldRow = d.Row
									}
								}

								// update เส้นเชื่อม
								for v := farestCol + 1; v < x; v++ {
									templateArr[v][detail.Row+1] = "111111"
								}

								// ลบก่อนย้ายทิ้ง
								for i := 0; i <= x; i++ {

									removeIndex(&templateArr[i], oldRow)
								}

							} else {

								available := true
								for o := value.Col + 1; o < x+1; o++ {
									if templateArr[o][value.Row] != "000000" {
										available = false
										break
									}
								}

								if available {
									for v := value.Col + 1; v < x; v++ {
										templateArr[v][detail.Row+1] = "111111"
									}
								}

							}
						}

					}
				} else {
					insertRow(&templateArr, detail.Row-1, corequisiteList)
					templateArr[x][detail.Row] = courseNo

					for key, value := range thisPreList {
						// ในแถววไม่มี coreq จะแทรก 1 แถว
						if key != headCourse && value.Move {
							insertRow(&templateArr, detail.Row, corequisiteList)

							farestCol := value.Col
							oldRow := value.Row

							for _, d := range value.PreReqCourseList {

								templateArr[d.Term][detail.Row] = d.CourseCode

								if farestCol < d.Term {
									farestCol = d.Term
								}
								if d.Row > detail.Row {
									oldRow = d.Row + 1
								} else {
									oldRow = d.Row
								}
							}

							// update เส้นเชื่อม
							for v := farestCol + 1; v < x; v++ {
								templateArr[v][detail.Row] = "111111"
							}

							// ลบก่อนย้ายทิ้ง
							for i := 0; i <= x; i++ {

								removeIndex(&templateArr[i], oldRow)
							}

						} else {

							available := true
							for o := value.Col + 1; o < x+1; o++ {
								if templateArr[o][value.Row] != "000000" {
									available = false
									break
								}
							}

							if available {
								for v := value.Col + 1; v < x; v++ {
									templateArr[v][detail.Row] = "111111"
								}
							}

						}
					}

				}
			}

		}

	}

	return templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList
}

func getAllListCourse(templateArr [][]string, startCourse string, haveRequisite map[string][]string, listOfCourse map[string]*model.CurriculumCourseDetail2, column int) (bool, []model.CourseList) {

	// true เมื่อสามารถย้ายทั้งแถวตัวต่อได้
	// false เมื่อไม่สามารถย้ายได้ => ตัวที่ตัวต่อทั้งหมดไม่ได้อยู่บนแถวเดียวกัน

	// ถ้ามีตัว pre 2 ตัวก็จะไม่ย้าย
	// ถ้าตัวก่อน มีตัวต่อ 2 ตัวก็ไม่ย้าย

	// check หา detail ของ course
	detail, b := listOfCourse[startCourse]
	list := []model.CourseList{}
	if b {

		// หาตำแหน่งของ course
		term, row := checkTermAndIndex(templateArr, startCourse)
		list = append(list, model.CourseList{
			CourseCode: startCourse,
			Term:       term,
			Row:        row,
		})

		req, bb := haveRequisite[startCourse]
		if bb {
			if len(req) > 1 {
				return false, list
			} else if len(detail.Prerequisites) > 1 {
				return false, list
			} else if len(detail.Corequisite) != 0 {
				return false, list
			} else {
				for _, c := range detail.Prerequisites {
					bb, f := getAllListCourse(templateArr, c, haveRequisite, listOfCourse, row)
					list = append(list, f...)
					if !bb {
						return bb, list
					}
				}
			}
		} else {
			if len(detail.Prerequisites) > 1 {
				return false, list
			} else if len(detail.Corequisite) != 0 {
				return false, list
			} else {
				for _, c := range detail.Prerequisites {
					bb, f := getAllListCourse(templateArr, c, haveRequisite, listOfCourse, row)
					list = append(list, f...)
					if !bb {
						return bb, list
					}
				}
			}
		}
	}

	return true, list
}

func checkTermAndIndex(templateArr [][]string, course string) (int, int) {

	for t, term := range templateArr {
		index := slices.Index[[]string](term, course)
		if index != -1 {
			return t, index
		}
	}

	return -1, -1

}

func updateTemplate(templateArr [][]string, x int, numberTerm int, updateIndex int, haveRequisite map[string][]string, listOfCourse map[string]*model.CurriculumCourseDetail2, exceptCoreq bool) [][]string {

	// x = term ปจจ (นับตาม arr)/ เทอมที่จะเลื่อน
	// updateIndex = index ของตัวที่ต้องเลื่อน
	// numberTerm = จำนวนเทอมทั้งหมด

	// เลื่อนแถวนั้นจากท้าย
	l := numberTerm - 1
	// reqRow := -1
	// reqCol := -1
	// เริ่มลูปจากตัวสุดท้ายของแถวนั้น

	for l >= x {

		//เช็คว่าตัวนี้เลื่อนมีตัวต่อไหม
		log.Println("l : ", l)
		log.Println("updateIndex : ", updateIndex)
		_, bb := haveRequisite[templateArr[l][updateIndex]]

		// ถ้ามีตัวต่อด้วยให้มีเส้นเชื่อม
		templateArr[l+2][updateIndex] = templateArr[l][updateIndex]

		detail, b := listOfCourse[templateArr[l+2][updateIndex]]
		if bb && len(detail.Prerequisites) != 0 {
			templateArr[l][updateIndex] = "111111"
			if l == x {
				templateArr[l+1][updateIndex] = "111111"
			}
		} else {
			templateArr[l][updateIndex] = "000000"
		}

		// ถ้ามีตัว coreq ให้เลื่อนตัว coreq ด้วย
		// detail, b := listOfCourse[templateArr[l+2][updateIndex]]
		if b && !exceptCoreq {
			if detail.Corequisite != "" {
				templateArr[l+2][updateIndex+1] = templateArr[l][updateIndex+1]
				templateArr[l][updateIndex+1] = "000000"
			}
		}

		l = l - 1

	}

	// ไล่เช็คจากตัวแรกว่ามีตัวไหนมีตัวต่อไหม
	// เลื่อนตัวต่อของมันด้วย
	l = numberTerm - 1
	start := x

	// เริ่มลูปจากตัวสุดท้ายของแถวนั้น
	updateRow := updateIndex

	for start <= l {

		reqList, b := haveRequisite[templateArr[start][updateIndex]]
		if b {

			// เช็คว่าตัวต่ออยู่ในแถวเดียวกันไหม
			// ถ้าไม่อยู่ก็เลื่อนตรงก้อนนั้นทั้งหมด
			if len(reqList) > 0 {
				for _, req := range reqList {
					col, index := checkTermAndIndex(templateArr, req)
					log.Println(templateArr)
					log.Println("req : ", req)
					log.Println("index 1621 : ", index)
					log.Println("updateIndex 1622 : ", updateIndex)

					if index != updateIndex {
						if index < updateRow {
							for i := index; i < updateIndex; i++ {
								templateArr = updateTemplate(templateArr, col, numberTerm, i, haveRequisite, listOfCourse, true)
							}
							updateRow = index
						}

						pre := listOfCourse[req]

						// อัปเดตเส้นเชื่อม
						templateArr[col][index] = "111111"
						templateArr[col+1][index] = "111111"

						if len(pre.Prerequisites) < 2 {
							for h := col - 1; h >= 0; h-- {
								if templateArr[h][index] == "111111" {
									templateArr[h][index] = "000000"
								} else {
									break
								}
							}
						}

					}

				}

			}

		}

		preReq, b := listOfCourse[templateArr[start][updateIndex]]
		// ถ้ามี prerequisites 2 ตัว
		// update เส้นเชื่อสำหรับตัวที่ผ่านแล้ว แต่อีกตัวไม่ผ่าน
		if b {
		}
		if b && len(preReq.Prerequisites) == 2 {
			for _, preReq := range preReq.Prerequisites {
				col, index := checkTermAndIndex(templateArr, preReq)

				for col = col + 1; col < start; col++ {
					templateArr[col][index] = "111111"
				}
			}
		}

		start++
	}

	return templateArr
}

func getTermTemplateV2(year string, curriculumProgram string, isCOOP string, studentId string, mockData string) ([][]string, map[string]*model.CurriculumCourseDetail2, []int, string, map[string][]string) {

	transcript := ""

	if studentId == "" {
		transcript = readMockData(mockData)
	} else {
		transcriptModel := getTranscriptWithCredit(studentId)
		tm, err := json.Marshal(transcriptModel)
		if err != nil {
			log.Fatalln("Error is : ", err)
		}

		transcript = string(tm)

		if !transcriptModel.Status {
			transcript = ""
		}
	}
	log.Println(transcript)
	if strings.Contains(transcript, COOPcourse) {
		isCOOP = "true"
	}

	templateArr := [][]string{}

	corequisiteList := []string{}
	havePreList := []string{}
	var listOfCourse = map[string]*model.CurriculumCourseDetail2{}
	var haveRequisite = map[string][]string{}

	fullCurriculum, _ := getCirriculum(year, curriculumProgram, isCOOP)

	i := 0
	x := 0
	// for คนยังไม่ได้เรียน

	// term template according to plan
	// loop year
	for i < 4 {

		j := 0
		// loop term
		for j < 2 {

			log.Println("x : ", x)

			term := []string{}
			noPreList := []string{}

			termString, _, _ := getTermDetail(year, curriculumProgram, isCOOP, strconv.Itoa(i+1), strconv.Itoa(j+1))

			if x != 0 {

				k := 0

				for k < len(templateArr[x-1]) {
					term = append(term, "000000")
					k++
				}
				templateArr = append(templateArr, term)
				prerequisitesList := []string{}

				coreList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Core").requiredCourses`)
				for _, core := range coreList.Array() {

					courseNo := core.Get("courseNo").String()
					prerequisites := core.Get("prerequisites").Array()
					corequisite := core.Get("corequisite").String()

					templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites, listOfCourse)

					listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
						CourseNo:      courseNo,
						Prerequisites: prerequisitesList,
						Corequisite:   corequisite,
						Credits:       int(core.Get("credits").Int()),
						GroupName:     "Core",
					}
				}

				majorList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Major Required").requiredCourses`).Array()
				for _, major := range majorList {

					courseNo := major.Get("courseNo").String()
					prerequisites := major.Get("prerequisites").Array()
					corequisite := major.Get("corequisite").String()

					templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites, listOfCourse)

					listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
						CourseNo:      courseNo,
						Prerequisites: prerequisitesList,
						Corequisite:   corequisite,
						Credits:       int(major.Get("credits").Int()),
						GroupName:     "Major Required",
					}
				}

				if isCOOP == "true" {
					majorList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Major Elective").electiveCourses`).Array()
					for _, major := range majorList {

						courseNo := major.Get("courseNo").String()
						prerequisites := major.Get("prerequisites").Array()
						corequisite := major.Get("corequisite").String()

						templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites, listOfCourse)

						listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
							CourseNo:      courseNo,
							Prerequisites: prerequisitesList,
							Corequisite:   corequisite,
							Credits:       int(major.Get("credits").Int()),
							GroupName:     "Major Required",
						}
					}
				}

				numberGE := gjson.Get(termString, `curriculum.geGroups.#`).Int()
				n := 0
				for n < int(numberGE) {
					groupname := gjson.Get(termString, `curriculum.geGroups.`+strconv.Itoa(n)+`.groupName`).String()
					geList := gjson.Get(termString, `curriculum.geGroups.`+strconv.Itoa(n)+`.requiredCourses`).Array()
					for _, ge := range geList {

						courseNo := ge.Get("courseNo").String()
						prerequisites := ge.Get("prerequisites").Array()
						corequisite := ge.Get("corequisite").String()

						templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites, listOfCourse)

						listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
							CourseNo:      courseNo,
							Prerequisites: prerequisitesList,
							Corequisite:   corequisite,
							Credits:       int(ge.Get("credits").Int()),
							GroupName:     groupname,
						}
					}
					n++
				}

			} else {
				coreList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Core").requiredCourses`)
				for _, core := range coreList.Array() {

					no := core.Get("courseNo").String()
					term = append(term, no)

					listOfCourse[no] = &model.CurriculumCourseDetail2{
						CourseNo:      no,
						Prerequisites: []string{},
						Corequisite:   core.Get("corequisite").String(),
						Credits:       int(core.Get("credits").Int()),
						GroupName:     "Core",
					}

				}

				majorList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Major Required").requiredCourses`).Array()
				for _, major := range majorList {
					no := major.Get("courseNo").String()
					term = append(term, no)

					listOfCourse[no] = &model.CurriculumCourseDetail2{
						CourseNo:      no,
						Prerequisites: []string{},
						Corequisite:   major.Get("corequisite").String(),
						Credits:       int(major.Get("credits").Int()),
						GroupName:     "Major Required",
					}
				}

				numberGE := gjson.Get(termString, `curriculum.geGroups.#`).Int()
				n := 0
				for n < int(numberGE) {
					geList := gjson.Get(termString, `curriculum.geGroups.`+strconv.Itoa(n)+`.requiredCourses`).Array()
					groupname := gjson.Get(termString, `curriculum.geGroups.`+strconv.Itoa(n)+`.groupName`).String()

					for _, ge := range geList {
						no := ge.Get("courseNo").String()
						term = append(term, no)

						listOfCourse[no] = &model.CurriculumCourseDetail2{
							CourseNo:      no,
							Prerequisites: []string{},
							Corequisite:   ge.Get("corequisite").String(),
							Credits:       int(ge.Get("credits").Int()),
							GroupName:     groupname,
						}
					}
					n++
				}

				templateArr = append(templateArr, term)
			}

			if len(noPreList) != 0 {

				for _, c := range noPreList {

					n := len(templateArr[x])
					insertRow(&templateArr, n, corequisiteList)
					templateArr[x][n] = c

				}
			}
			log.Println("templateArr : ", templateArr)
			j++
			x++

		}
		i++
	}

	// check program that user choose
	// get elective course for this program
	elective := readMockData("freeNormalPlan")
	if isCOOP == "true" {
		elective = readMockData("freeCoopPlan")
	}

	requiredRow := len(templateArr[0])

	// get requirded credit of elective course
	geNum := gjson.Get(elective, "curriculum.geGroups.#").Int()
	var numberFree = map[string]int{}
	for l := 0; l < int(geNum); l++ {
		groupName := gjson.Get(elective, `curriculum.geGroups.`+strconv.Itoa(l)+`.groupName`).String()
		numberFree[groupName] = int(gjson.Get(elective, `curriculum.geGroups.`+strconv.Itoa(l)+`.requiredCredits`).Int())
	}
	numberFree["Major Elective"] = int(gjson.Get(elective, `curriculum.coreAndMajorGroups.0.requiredCredits`).Int())
	numberFree["Free"] = int(gjson.Get(elective, `curriculum.freeGroups.0.requiredCredits`).Int())

	numOfTerm := []int{}
	if transcript != "" {

		// check with student enroll
		yearListNum := gjson.Get(transcript, "transcript.#").Int()

		x = 0
		if yearListNum != 0 {

			yearList := gjson.Get(transcript, "transcript").Array()

			// loop in year
			for y, yearDetail := range yearList {

				t := 0
				summerList := []string{}
				termList := yearDetail.Get("yearDetails").Array()
				if len(termList) == 3 {
					summerTerm := gjson.Get(transcript, `transcript.`+strconv.Itoa(y)+`.yearDetails.2.details`).Array()

					for _, courseDetail := range summerTerm {

						if slices.Contains[[]string](PASS_GRADE, courseDetail.Get("grade").String()) {
							summerList = append(summerList, courseDetail.Get("code").String())
						}
					}
				}

				// loop in semester
				for _, termDetail := range termList {

					pass := []string{}
					freePass := []string{}

					detail := termDetail.Get("details").Array()

					// add success course in pass[]
					for _, courseDetail := range detail {

						grade := courseDetail.Get("grade").String()

						//check if success
						if slices.Contains[[]string](PASS_GRADE, grade) {

							// check if it elective course
							code := courseDetail.Get("code").String()
							_, isReq := listOfCourse[code]
							if !isReq {

								// check elective group
								group, _ := checkGroup(fullCurriculum, code)
								log.Println("group : ", group)
								credit := courseDetail.Get("credit").Int()
								// courseDetail := getCourseDetail(code)

								// add to list of study course
								if numberFree[group] > 0 {
									listOfCourse[code] = &model.CurriculumCourseDetail2{
										CourseNo:          code,
										RecommendSemester: 0,
										RecommendYear:     0,
										Prerequisites:     []string{},
										Corequisite:       "",
										Credits:           int(credit),
										GroupName:         group,
									}
								} else {
									listOfCourse[code] = &model.CurriculumCourseDetail2{
										CourseNo:          code,
										RecommendSemester: 0,
										RecommendYear:     0,
										Prerequisites:     []string{},
										Corequisite:       "",
										Credits:           int(credit),
										GroupName:         "Free Elective",
									}
								}

								// edit credit
								if numberFree[group] > 0 {
									numberFree[group] = numberFree[group] - int(credit)
								} else {
									numberFree["Free"] = numberFree["Free"] - int(credit)
								}

								freePass = append(freePass, code)
								pass = append(pass, group)

							} else {
								pass = append(pass, code)
							}

						}
					}

					// map study course into template
					// check if summer term
					if t == 2 {
						// summer term add 1 row
						lenX := len(templateArr[x])
						term3 := []string{}
						// ใส่ summer term โดยตรวจว่ามีเส้นเชื่อมไหม
						for k := 0; k < lenX; k++ {

							if templateArr[x-1][k] != "000000" && templateArr[x][k] != "000000" {
								term3 = append(term3, "111111")
							} else {
								term3 = append(term3, "000000")
							}

						}
						templateArr = slices.Insert[[][]string](templateArr, x, term3)

						// ใส่ตัวที่มีใน template ก่อน
						for _, c := range pass {
							term, index := checkTermAndIndex(templateArr, c)
							if term != -1 && index != -1 {
								templateArr[x][index] = c
								templateArr[term][index] = "111111"

							}
						}

					} else {

						first := true
						for index, temp := range templateArr[x] {

							log.Println(temp)

							contain := slices.Contains[[]string](pass, temp)
							if !contain && temp != "000000" && temp != "111111" {

								last := len(templateArr)
								lenX := len(templateArr[x])

								// check ว่าช่องสุดท้ายที่จะเลื่อนไปว่างไหม ถ้าว่างก็ไม่ต้องเพิ่มแถว
								// && (templateArr[last-1][lenX-1] != "000000" || templateArr[last-2][lenX-1] != "000000")
								if !first {
									last = last - 2
								}

								if first {

									term := []string{}
									term2 := []string{}
									for k := 0; k < lenX; k++ {
										term = append(term, "000000")
										term2 = append(term2, "000000")
									}

									templateArr = append(templateArr, term)
									templateArr = append(templateArr, term2)
									first = false
								}

								// loop เลื่อน course ที่ยังไม่ได้เรียน
								// check ใน แถวที่เลื่อนว่าตัวไหนมี pre
								// if t != 1 && len(termList) != 3 && !slices.Contains[[]string](summerList, temp) {

								if len(termList) == 3 && t == 1 && slices.Contains[[]string](summerList, temp) {
									// สำหรับการณีที่มี summer และเป็น term 2 และเรียนใน summer
									// do notting
								} else {
									// สำหรับการณีที่มี summer และเป็น term 1
									log.Println("index : ", index)
									templateArr = updateTemplate(templateArr, x, last, index, haveRequisite, listOfCourse, false)
								}

							}

							if index == requiredRow-1 {
								break
							}

						}

					}

					for _, f := range freePass {
						for i, temp := range templateArr[x] {
							if temp == "000000" {
								templateArr[x][i] = f
								break
							}
						}
					}

					log.Println("after map to template term ", x+1, " : ", templateArr)

					t++
					x++

				}

				// เก็บถึงเทอมที่เรียนเสร็จ
				numOfTerm = append(numOfTerm, t)
			}
		}

		// ตรวจว่าเลื่อนไปกี่เทอม
		addNew := len(templateArr) - 8
		nowTerm := len(templateArr) - 1
		log.Println("addNew : ", addNew)

		// map free elective ที่เหลือเข้าไปใน template
		// GE
		for l := 0; l < int(geNum); l++ {

			groupName := gjson.Get(elective, `curriculum.geGroups.`+strconv.Itoa(l)+`.groupName`).String()
			geCourse := gjson.Get(elective, `curriculum.geGroups.`+strconv.Itoa(l)+`.electiveCourses`).Array()

			log.Println(`numberFree[`+groupName+`] : `, numberFree[groupName])
			// check need more credit
			have := -1
			if numberFree[groupName] > 0 {
				needMore := numberFree[groupName] / 3
				if numberFree[groupName]%3 != 0 {
					needMore = needMore + 1
				}

				have = len(geCourse) - needMore
			}

			for _, ge := range geCourse {

				if have == 0 {

					// คำนวณเทอมใหม่ อิงจากเทอมที่ควรจะอยู่
					term := ge.Get("recommendSemester").Int()
					year := ge.Get("recommendYear").Int()
					x := ((int(year) - 1) * 2) + int(term) - 1 + addNew
					if x < nowTerm {
						x = nowTerm
					}

					success := false
					for i, temp := range templateArr[x] {
						if temp == "000000" {
							templateArr[x][i] = groupName
							success = true
							break
						}
					}

					// ถ้าไม่ม่ช่องให้เติมก็เพิ่มช่องเข้าไป
					if !success {
						lenX := len(templateArr[x])
						insertRow(&templateArr, lenX, corequisiteList)
						templateArr[x][lenX] = "Free"
					}

				} else {
					have--
				}

			}

		}

		// map free elective ที่เหลือเข้าไปใน template
		// Major
		majorCourse := gjson.Get(elective, `curriculum.coreAndMajorGroups.0.electiveCourses`).Array()

		log.Println("numberFree[Major Elective] : ", numberFree["Major Elective"])

		have := -1
		if numberFree["Major Elective"] > 0 {
			needMore := numberFree["Major Elective"] / 3
			if numberFree["Major Elective"]%3 != 0 {
				needMore = needMore + 1
			}

			have = len(majorCourse) - needMore
			log.Println("need more : ", needMore)

		}

		for _, major := range majorCourse {

			if have == 0 {

				term := major.Get("recommendSemester").Int()
				year := major.Get("recommendYear").Int()
				x := ((int(year) - 1) * 2) + int(term) - 1 + addNew
				if x < nowTerm {
					x = nowTerm
				}

				success := false
				for i, temp := range templateArr[x] {
					if temp == "000000" {
						templateArr[x][i] = "Major Elective"
						success = true
						break
					}
				}

				// ถ้าไม่ม่ช่องให้เติมก็เพิ่มช่องเข้าไป
				if !success {
					lenX := len(templateArr[x])
					insertRow(&templateArr, lenX, corequisiteList)
					templateArr[x][lenX] = "Major Elective"
				}

			} else {
				have--
			}

		}

		// map free elective ที่เหลือเข้าไปใน template
		// Free
		freeCourse := gjson.Get(elective, `curriculum.freeGroups.0.electiveCourses`).Array()

		log.Println("numberFree[Free] : ", numberFree["Free"])

		have = -1
		if numberFree["Free"] > 0 {
			needMore := numberFree["Free"] / 3
			if numberFree["Free"]%3 != 0 {
				needMore = needMore + 1
			}

			log.Println("need more : ", needMore)

			have = len(freeCourse) - needMore
		}

		for _, free := range freeCourse {

			if have == 0 {
				term := free.Get("recommendSemester").Int()
				year := free.Get("recommendYear").Int()
				x := ((int(year) - 1) * 2) + int(term) - 1 + addNew
				if x < nowTerm {
					x = nowTerm
				}

				success := false
				for i, temp := range templateArr[x] {
					if temp == "000000" {
						templateArr[x][i] = "Free"
						success = true
						break
					}
				}

				// ถ้าไม่ม่ช่องให้เติมก็เพิ่มช่องเข้าไป
				if !success {
					lenX := len(templateArr[x])
					insertRow(&templateArr, lenX, corequisiteList)
					templateArr[x][lenX] = "Free"
				}

			} else {
				have--
			}

		}
	}
	log.Println("Final templateArr : ", templateArr)

	return templateArr, listOfCourse, numOfTerm, isCOOP, haveRequisite
}

func main() {
	e := echo.New()

	e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/major/elective", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		cirriculum, _ := getCirriculum(year, curriculumProgram, isCOOP)
		query := `curriculum.coreAndMajorGroups.#(groupName=="Major Elective").electiveCourses`
		value := gjson.Get(cirriculum, query)

		return c.JSON(http.StatusOK, echo.Map{"courseLists": value.Value()})
	})

	e.GET("/ge/elective", func(c echo.Context) error {

		groupName := c.QueryParam("groupName")

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		cirriculum, _ := getCirriculum(year, curriculumProgram, isCOOP)
		query := `curriculum.geGroups.#(groupName=="` + groupName + `").electiveCourses`
		value := gjson.Get(cirriculum, query)

		return c.JSON(http.StatusOK, echo.Map{"courseLists": value.Value()})
	})

	e.GET("/categoryView", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")
		mockData := c.QueryParam("mockData")
		studentId := c.QueryParam("studentId")

		cirriculumJSON, _ := getCirriculumJSON(year, curriculumProgram, isCOOP)
		curriculumString, _ := getCirriculum(year, curriculumProgram, isCOOP)

		template, summaryCredits, isCoop, err := getCategoryTemplate(cirriculumJSON, curriculumString, isCOOP, studentId, mockData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		reponse := model.CategoryResponse{}
		err = json.Unmarshal([]byte(template), &reponse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		reponse.IsCoop = isCoop
		reponse.SummaryCredits = summaryCredits

		return c.JSON(http.StatusOK, reponse)

	})

	e.GET("/termView", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")
		mockData := c.QueryParam("mockData")
		studentId := c.QueryParam("studentId")

		templateArr, listOfCourse, numOfTerm, isCoop, haveRequisite := getTermTemplateV2(year, curriculumProgram, isCOOP, studentId, mockData)

		return c.JSON(http.StatusOK, echo.Map{"isCoop": isCoop, "study term": numOfTerm, "template": templateArr, "list of course": listOfCourse, "haveRequisite": haveRequisite})
	})

	e.GET("/summaryCredits", func(c echo.Context) error {

		studentId := c.QueryParam("studentId")
		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")
		mockData := c.QueryParam("mockData")

		cirriculumJSON, _ := getCirriculumJSON(year, curriculumProgram, isCOOP)
		curriculumString, _ := getCirriculum(year, curriculumProgram, isCOOP)

		summaryCredits, err := getSummaryCredits(cirriculumJSON, curriculumString, isCOOP, studentId, mockData)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, summaryCredits)
	})

	e.GET("/checkGroup", func(c echo.Context) error {

		courseNo := c.QueryParam("courseNo")

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		curriculumString, _ := getCirriculum(year, curriculumProgram, isCOOP)
		group, courseType := checkGroup(curriculumString, courseNo)

		return c.JSON(http.StatusOK, echo.Map{"group": group, "courseType": courseType})
	})

	e.GET("/test", func(c echo.Context) error {

		studentId := c.QueryParam("studentId")

		mo := getTranscript(studentId)
		return c.JSON(http.StatusOK, mo)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
