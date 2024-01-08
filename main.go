package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"NuaikitTracking_BE.com/model"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
)

var PASS_GRADE = []string{"A", "B", "C", "D", "S"}
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

func getCirriculum(year string, curriculumProgram string, isCOOP string) string {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/curriculum"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
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
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
	}

	return string(c)
}

func getCirriculumJSON(year string, curriculumProgram string, isCOOP string) model.CurriculumModel {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/curriculum"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
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
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
	}

	curriculum := model.CurriculumModel{}
	err = json.Unmarshal(c, &curriculum)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	return curriculum

}

func getTermDetail(year string, curriculumProgram string, isCOOP string, studyYear string, studySemester string) (string, model.CurriculumModel) {

	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/curriculum"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
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
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
	}

	term := model.CurriculumModel{}
	err = json.Unmarshal(c, &term)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	return string(c), term

}

func getCourseDetail(courseNo string) model.GetCourseDetail {
	client := &http.Client{}

	cpeAPI := goDotEnvVariable("CPE_API_URL")
	cpeToken := goDotEnvVariable("CPE_API_TOKEN")

	url := cpeAPI + "/course/detail"
	bearer := "Bearer " + cpeToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	req.Header.Add("Authorization", bearer)
	q := req.URL.Query()
	q.Add("courseNo", courseNo)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	c, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		log.Fatalln("Error is : ", err)
	}

	detail := model.GetCourseDetail{}

	err = json.Unmarshal(c, &detail)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	return detail

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

func getSummaryCredits(c model.CurriculumModel, curriculumString string, mockData string, isCOOP string) (model.CategoryResponseV2, error) {

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

	transcript := readMockData(mockData)
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

	err = json.Unmarshal([]byte(template), &t)
	if err != nil {
		return t, err
	}

	t.SummaryCredits = summaryCredits

	return t, nil
}

func getCategoryTemplate(c model.CurriculumModel, curriculumString string, mockData string, isCOOP string) (string, int, error) {

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

	transcript := readMockData(mockData)
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
								return "", 0, err
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
									return "", 0, err
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
									return "", 0, err
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
										return "", 0, err
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
										return "", 0, err
									}

									courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
										CourseNo:  code.String(),
										Credits:   int(credit),
										GroupName: group,
										IsPass:    true,
									})
								}
							}

							template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").electiveCourseList`, courseList)
						}

					}
					summaryCredits += int(credit)
				}

			}

			i++
		}
	}

	return template, summaryCredits, nil

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

func RemoveIndex(s *[]string, value string) {
	index := slices.Index[[]string](*s, value)
	*s = append((*s)[:index], (*s)[index+1:]...)
}

func putInTemplate(templateArr [][]string, x int, corequisiteList []string, noPreList []string, havePreList []string, haveRequisite map[string][]string, corequisite string, courseNo string, prerequisites []gjson.Result) ([][]string, []string, []string, []string, map[string][]string, []string) {

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

			// have 2 prerequisites
			havePreList = append(havePreList, courseNo)
			prerequisitesList = append(prerequisitesList, prerequisites[0].String())
			arr, h := haveRequisite[prerequisites[0].String()]
			if !h {
				haveRequisite[prerequisites[0].String()] = []string{courseNo}
			} else {
				arr = append(arr, courseNo)
				haveRequisite[prerequisites[0].String()] = arr
			}
			prerequisitesList = append(prerequisitesList, prerequisites[1].String())
			arr, h = haveRequisite[prerequisites[1].String()]
			if !h {
				haveRequisite[prerequisites[1].String()] = []string{courseNo}
			} else {
				arr = append(arr, courseNo)
				haveRequisite[prerequisites[1].String()] = arr
			}

			// find position of prerequisites
			preRow1 := 0
			preCol1 := 0

			preRow2 := 0
			preCol2 := 0
			preCourseNO2 := ""
			for _, p := range prerequisites {
				if slices.Contains[[]string](havePreList, p.String()) {
					for col := range templateArr {
						row := slices.Index[[]string](templateArr[col], p.String())
						if row != -1 {
							preRow1 = row
							preCol1 = col
							break
						}
					}
				} else {
					for col := range templateArr {
						row := slices.Index[[]string](templateArr[col], p.String())
						if row != -1 {
							if preCourseNO2 == "" {
								preRow2 = row
								preCol2 = col
								preCourseNO2 = p.String()

							} else {
								preRow1 = row
								preCol1 = col

							}
							break
						}
					}
				}
			}

			if preRow1+1 == preRow2 && preCol1 == preCol2 && templateArr[x][preRow1] == "000000" {
				templateArr[x][preRow1] = courseNo
			} else {

				// เช็คว่าตต pre1 ไปจน column ล่าสุด มี course ไหนมี corequisite ไหม
				haveCoreq := true
				for o := preCol1; o < x; o++ {
					if slices.Contains[[]string](corequisiteList, templateArr[o][preRow1+1]) {
						haveCoreq = false
						break
					}
				}

				if !haveCoreq {

					// มี course ที่มี corequisite
					// เพิ่ม 2 แถวลงใต้ pre1 => pre1Row + 1

					// เช็คว่าช่องต่อจาก pre1 ว่างไหม
					available := true
					for o := preCol1 + 1; o < x+1; o++ {
						if templateArr[o][preRow1] != "000000" {
							available = false
							break
						}
					}

					if available {
						// ช่องต่อจาก pre1 ว่าง
						insertRow(&templateArr, preRow1+1, corequisiteList)
						insertRow(&templateArr, preRow1+1, corequisiteList)
						templateArr[x][preRow1] = courseNo
						templateArr[preCol2][preRow1+2] = preCourseNO2

						for v := preCol2 + 1; v < x; v++ {

							templateArr[v][preRow1+2] = "111111"

						}

					} else {

						// ช่องต่อจาก pre1 ไม่ว่าง
						// ไปใส่ช่องของ pre2 แทน
						insertRow(&templateArr, preRow1+1, corequisiteList)
						insertRow(&templateArr, preRow1+1, corequisiteList)
						templateArr[x][preRow1+2] = courseNo
						templateArr[preCol2][preRow1+2] = preCourseNO2

						for v := preCol2 + 1; v < x; v++ {

							templateArr[v][preRow1+2] = "111111"

						}
					}

					// update ช่องเดิมของ pre2 เป็น 000000
					if preRow2 < preRow1 {
						templateArr[preCol2][preRow2] = "000000"
					} else {
						templateArr[preCol2][preRow2+2] = "000000"
					}

				} else {

					// ไม่มี course ไหนมี corequisite
					// เพิ่มแถวลงใต้ pre1 => pre1Row + 1
					available := true
					for o := preCol1 + 1; o < x+1; o++ {
						if templateArr[o][preRow1] != "000000" {
							available = false
							break
						}
					}

					if available {

						// ช่องต่อจาก pre1 ว่าง
						insertRow(&templateArr, preRow1+1, corequisiteList)
						templateArr[x][preRow1] = courseNo
						templateArr[preCol2][preRow1+1] = preCourseNO2

						for v := preCol2 + 1; v < x; v++ {

							templateArr[v][preRow1+1] = "111111"

						}

					} else {

						// ช่องต่อจาก pre1 ไม่ว่าง
						// ไปใส่ช่องของ pre2 แทน
						insertRow(&templateArr, preRow1+1, corequisiteList)
						templateArr[x][preRow1+1] = courseNo
						templateArr[preCol2][preRow1+1] = preCourseNO2

						for v := preCol2 + 1; v < x; v++ {

							templateArr[v][preRow1+1] = "111111"

						}
					}

					// update ช่องเดิมของ pre2 เป็น 000000
					if preRow2 < preRow1 {
						templateArr[preCol2][preRow2] = "000000"
					} else {
						templateArr[preCol2][preRow2+1] = "000000"
					}

				}

			}

		}

	}

	return templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList
}

func checkTermAndIndex(templateArr [][]string, course string) (int, int) {

	log.Println("course : ", course)
	for t, term := range templateArr {
		index := slices.Index[[]string](term, course)
		if index != -1 {
			return t, index
		}
	}

	return -1, -1

}

func updateTemplate(templateArr [][]string, x int, lastRow int, updateIndex int, haveRequisite map[string][]string) [][]string {

	l := lastRow - 1
	reqRow := -1
	reqCol := -1
	for l >= x {

		reqList, b := haveRequisite[templateArr[l][updateIndex]]

		if b {
			for _, req := range reqList {
				col, index := checkTermAndIndex(templateArr, req)

				if index == -1 || col == -1 {
					break
				}
				if index != updateIndex {
					log.Println("not equal")

					if reqRow == -1 || reqRow > index {
						reqRow = index
						reqCol = col
					}
				} else {
					log.Println("equal : ")
				}
			}
		}

		templateArr[l+2][updateIndex] = templateArr[l][updateIndex]

		templateArr[l][updateIndex] = "000000"

		l = l - 1

	}

	if reqCol != -1 && reqRow != -1 {
		for q := reqRow; q < updateIndex; q++ {

			templateArr = updateTemplate(templateArr, reqCol, lastRow, q, haveRequisite)

		}
	}

	return templateArr
}

func getTermTemplateV2(transcript string, year string, curriculumProgram string, isCOOP string) ([][]string, map[string]*model.CurriculumCourseDetail2, []int) {

	templateArr := [][]string{}

	corequisiteList := []string{}
	havePreList := []string{}
	var listOfCourse = map[string]*model.CurriculumCourseDetail2{}
	var haveRequisite = map[string][]string{}

	fullCurriculum := getCirriculum(year, curriculumProgram, isCOOP)

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

			termString, _ := getTermDetail(year, curriculumProgram, isCOOP, strconv.Itoa(i+1), strconv.Itoa(j+1))

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

					templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites)

					listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
						CourseNo:      courseNo,
						Prerequisites: prerequisitesList,
						Corequisite:   corequisite,
						Credits:       int(core.Get("credits").Int()),
					}
				}

				majorList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Major Required").requiredCourses`).Array()
				for _, major := range majorList {

					courseNo := major.Get("courseNo").String()
					prerequisites := major.Get("prerequisites").Array()
					corequisite := major.Get("corequisite").String()

					templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites)

					listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
						CourseNo:      courseNo,
						Prerequisites: prerequisitesList,
						Corequisite:   corequisite,
						Credits:       int(major.Get("credits").Int()),
					}
				}

				if isCOOP == "true" {
					majorList := gjson.Get(termString, `curriculum.coreAndMajorGroups.#(groupName=="Major Elective").electiveCourses`).Array()
					for _, major := range majorList {

						courseNo := major.Get("courseNo").String()
						prerequisites := major.Get("prerequisites").Array()
						corequisite := major.Get("corequisite").String()

						templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites)

						listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
							CourseNo:      courseNo,
							Prerequisites: prerequisitesList,
							Corequisite:   corequisite,
							Credits:       int(major.Get("credits").Int()),
						}
					}
				}

				numberGE := gjson.Get(termString, `curriculum.geGroups.#`).Int()
				n := 0
				for n < int(numberGE) {
					geList := gjson.Get(termString, `curriculum.geGroups.`+strconv.Itoa(n)+`.requiredCourses`).Array()
					for _, ge := range geList {

						courseNo := ge.Get("courseNo").String()
						prerequisites := ge.Get("prerequisites").Array()
						corequisite := ge.Get("corequisite").String()

						templateArr, corequisiteList, noPreList, havePreList, haveRequisite, prerequisitesList = putInTemplate(templateArr, x, corequisiteList, noPreList, havePreList, haveRequisite, corequisite, courseNo, prerequisites)

						listOfCourse[courseNo] = &model.CurriculumCourseDetail2{
							CourseNo:      courseNo,
							Prerequisites: prerequisitesList,
							Corequisite:   corequisite,
							Credits:       int(ge.Get("credits").Int()),
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
					}
				}

				numberGE := gjson.Get(termString, `curriculum.geGroups.#`).Int()
				n := 0
				for n < int(numberGE) {
					geList := gjson.Get(termString, `curriculum.geGroups.`+strconv.Itoa(n)+`.requiredCourses`).Array()
					for _, ge := range geList {
						no := ge.Get("courseNo").String()
						term = append(term, no)

						listOfCourse[no] = &model.CurriculumCourseDetail2{
							CourseNo:      no,
							Prerequisites: []string{},
							Corequisite:   ge.Get("corequisite").String(),
							Credits:       int(ge.Get("credits").Int()),
						}
					}
					n++
				}

				templateArr = append(templateArr, term)
			}

			if len(noPreList) != 0 {
				// lastRow := len(templateArr[x]) - 1
				// for l := lastRow; l >= 0; l-- {
				// 	if templateArr[x][l] != "000000" {
				// 		lastRow++
				// 		break
				// 	} else {
				// 		lastRow--
				// 	}
				// }

				// if lastRow < 0 {
				// 	lastRow = 0
				// }

				// for _, c := range noPreList {
				// 	if lastRow >= len(templateArr[x]) {
				// 		insertRow(&templateArr, lastRow, corequisiteList)
				// 		templateArr[x][lastRow] = c
				// 	} else {
				// 		templateArr[x][lastRow] = c
				// 	}
				// 	lastRow++
				// }

				for _, c := range noPreList {
					for n := 0; n < len(templateArr[x])+1; n++ {
						if n >= len(templateArr[x]) {
							insertRow(&templateArr, n, corequisiteList)
							templateArr[x][n] = c
							break
						}

						if templateArr[x][n] == "000000" {
							templateArr[x][n] = c
							break
						}
					}
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

	// check with student enroll
	yearListNum := gjson.Get(transcript, "transcript.#").Int()
	numOfTerm := []int{}
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
							credit := courseDetail.Get("credit").Int()
							courseDetail := getCourseDetail(code)

							// add to list of study course
							listOfCourse[code] = &model.CurriculumCourseDetail2{
								CourseNo:          courseDetail.CourseDetail[0].CourseNo,
								RecommendSemester: 0,
								RecommendYear:     0,
								Prerequisites:     []string{},
								Corequisite:       "",
								Credits:           int(credit),
								GroupName:         group,
							}

							// edit credit
							if numberFree[group] <= 0 {
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

				log.Println("study in term ", x+1, " : ", pass, x+1)

				// map study course into template
				// check if summer term
				if t == 2 {
					// summer term add 1 row
					lenX := len(templateArr[x])
					term3 := []string{}
					for k := 0; k < lenX; k++ {

						term3 = append(term3, "000000")
					}
					templateArr = slices.Insert[[][]string](templateArr, x, term3)

					for _, c := range pass {
						term, index := checkTermAndIndex(templateArr, c)
						templateArr[x][index] = c
						templateArr[term][index] = "000000"
					}

				} else {

					first := true
					for index, temp := range templateArr[x] {

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
							if t != 1 && len(termList) != 3 && !slices.Contains[[]string](summerList, temp) {
								templateArr = updateTemplate(templateArr, x, last, index, haveRequisite)
							}

						}

						if index == requiredRow-1 {
							break
						}

					}

				}

				// map free elective that pass in template
				for _, f := range freePass {
					lenX := len(templateArr[x])
					if lenX == requiredRow {
						insertRow(&templateArr, lenX, corequisiteList)
						templateArr[x][lenX] = f
					} else {
						addRow := true
						for u := requiredRow; u < lenX; u++ {
							if templateArr[x][u] == "000000" {
								templateArr[x][u] = f
								addRow = false
								break
							}
						}

						if addRow {
							insertRow(&templateArr, lenX, corequisiteList)
							templateArr[x][lenX] = f
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

	// map free elective ที่เหลือเข้าไปใน template
	// GE
	for l := 0; l < int(geNum); l++ {

		groupName := gjson.Get(elective, `curriculum.geGroups.`+strconv.Itoa(l)+`.groupName`).String()
		geCourse := gjson.Get(elective, `curriculum.geGroups.`+strconv.Itoa(l)+`.electiveCourses`).Array()

		// check need more credit
		needMore := numberFree[groupName] / 3
		if numberFree[groupName]%3 != 0 {
			needMore = needMore + 1
		}

		needMore = len(geCourse) - needMore

		for _, ge := range geCourse {

			if needMore <= 0 {

				// คำนวณเทอมใหม่ อิงจากเทอมที่ควรจะอยู่
				term := ge.Get("recommendSemester").Int()
				year := ge.Get("recommendYear").Int()
				x := ((int(year) - 1) * 2) + int(term) - 1 + addNew

				lenX := len(templateArr[x])
				if lenX == requiredRow {
					insertRow(&templateArr, lenX, corequisiteList)
					templateArr[x][lenX] = groupName
				} else {
					addRow := true
					for u := requiredRow; u < lenX; u++ {
						if templateArr[x][u] == "000000" {
							templateArr[x][u] = groupName
							addRow = false
							break
						}
					}

					if addRow {
						insertRow(&templateArr, lenX, corequisiteList)
						templateArr[x][lenX] = groupName
					}
				}
			}

			needMore--

		}

	}

	// map free elective ที่เหลือเข้าไปใน template
	// Major
	majorCourse := gjson.Get(elective, `curriculum.coreAndMajorGroups.0.electiveCourses`).Array()

	needMore := numberFree["Major Elective"] / 3
	if numberFree["Major Elective"]%3 != 0 {
		needMore = needMore + 1
	}

	needMore = len(majorCourse) - needMore

	for _, major := range majorCourse {

		if needMore <= 0 {

			term := major.Get("recommendSemester").Int()
			year := major.Get("recommendYear").Int()
			x := ((int(year) - 1) * 2) + int(term) - 1 + addNew

			lenX := len(templateArr[x])
			if lenX == requiredRow {
				insertRow(&templateArr, lenX, corequisiteList)
				templateArr[x][lenX] = "Major Elective"
			} else {
				addRow := true
				for u := requiredRow; u < lenX; u++ {
					if templateArr[x][u] == "000000" {
						templateArr[x][u] = "Major Elective"
						addRow = false
						break
					}
				}

				if addRow {
					insertRow(&templateArr, lenX, corequisiteList)
					templateArr[x][lenX] = "Major Elective"
				}
			}
		}
		needMore--

	}

	// map free elective ที่เหลือเข้าไปใน template
	// Free
	freeCourse := gjson.Get(elective, `curriculum.freeGroups.0.electiveCourses`).Array()

	needMore = numberFree["Free"] / 3
	if numberFree["Free"]%3 != 0 {
		needMore = needMore + 1
	}

	needMore = len(freeCourse) - needMore
	for _, free := range freeCourse {
		if needMore <= 0 {
			term := free.Get("recommendSemester").Int()
			year := free.Get("recommendYear").Int()
			x := ((int(year) - 1) * 2) + int(term) - 1 + addNew

			lenX := len(templateArr[x])
			if lenX == requiredRow {
				insertRow(&templateArr, lenX, corequisiteList)
				templateArr[x][lenX] = "Free"
			} else {
				addRow := true
				for u := requiredRow; u < lenX; u++ {
					if templateArr[x][u] == "000000" {
						templateArr[x][u] = "Free"
						addRow = false
						break
					}
				}

				if addRow {
					insertRow(&templateArr, lenX, corequisiteList)
					templateArr[x][lenX] = "Free"
				}
			}
		}
		needMore--

	}

	log.Println("Final templateArr : ", templateArr)

	return templateArr, listOfCourse, numOfTerm
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

		cirriculum := getCirriculum(year, curriculumProgram, isCOOP)
		query := `curriculum.coreAndMajorGroups.#(groupName=="Major Elective").electiveCourses`
		value := gjson.Get(cirriculum, query)

		return c.JSON(http.StatusOK, echo.Map{"courseLists": value.Value()})
	})

	e.GET("/ge/elective", func(c echo.Context) error {

		groupName := c.QueryParam("groupName")

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		cirriculum := getCirriculum(year, curriculumProgram, isCOOP)
		query := `curriculum.geGroups.#(groupName=="` + groupName + `").electiveCourses`
		value := gjson.Get(cirriculum, query)

		return c.JSON(http.StatusOK, echo.Map{"courseLists": value.Value()})
	})

	e.GET("/categoryView", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")
		mockData := c.QueryParam("mockData")

		cirriculumJSON := getCirriculumJSON(year, curriculumProgram, isCOOP)
		curriculumString := getCirriculum(year, curriculumProgram, isCOOP)

		template, summaryCredits, err := getCategoryTemplate(cirriculumJSON, curriculumString, mockData, isCOOP)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		reponse := model.CategoryResponse{}
		err = json.Unmarshal([]byte(template), &reponse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		reponse.SummaryCredits = summaryCredits

		return c.JSON(http.StatusOK, reponse)

	})

	e.GET("/termView", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		mockData := c.QueryParam("mockData")

		// log.Println(getTermTemplate(year, curriculumProgram, isCOOP))

		transcript := readMockData(mockData)

		templateArr, listOfCourse, numOfTerm := getTermTemplateV2(transcript, year, curriculumProgram, isCOOP)

		// term := len(templateArr)
		// num := len(templateArr[0])

		// for i := 0; i < num; i++ {
		// 	line := ""
		// 	for j := 0; j < term; j++ {
		// 		line = line + templateArr[j][i] + " "
		// 	}
		// 	log.Println(line)
		// }

		return c.JSON(http.StatusOK, echo.Map{"study term": numOfTerm, "template": templateArr, "list of course": listOfCourse})
	})

	e.GET("/summaryCredits", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		mockData := c.QueryParam("mockData")

		cirriculumJSON := getCirriculumJSON(year, curriculumProgram, isCOOP)
		curriculumString := getCirriculum(year, curriculumProgram, isCOOP)

		summaryCredits, err := getSummaryCredits(cirriculumJSON, curriculumString, mockData, isCOOP)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, summaryCredits)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
