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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/exp/slices"
)

var PASS_GRADE = []string{"A", "B", "C", "D", "S"}

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
		log.Println(err)
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

func getRequiredCredits(cirriculum string) model.Credits {

	freeRequiredCredits := gjson.Get(cirriculum, "curriculum.freeElectiveCredits")
	coreRequiredCredits := gjson.Get(cirriculum, `curriculum.coreAndMajorGroups.#(groupName=="Core").requiredCredits`)
	majorRequiredCredits := gjson.Get(cirriculum, `curriculum.coreAndMajorGroups.#(groupName=="Major Required").requiredCredits`)
	majorElectiveCredits := gjson.Get(cirriculum, `curriculum.coreAndMajorGroups.#(groupName=="Major Elective").requiredCredits`)

	geCredits := 0
	groupList := gjson.Get(cirriculum, "curriculum.geGroups.#.groupName")
	for _, groupName := range groupList.Array() {

		queryReqCourse := `curriculum.geGroups.#(groupName=="` + groupName.String() + `").requiredCredits`
		credits := gjson.Get(cirriculum, queryReqCourse)
		geCredits += int(credits.Int())
	}

	return model.Credits{
		CoreCredits:  int(coreRequiredCredits.Int()),
		MajorCredits: int(majorRequiredCredits.Int() + majorElectiveCredits.Int()),
		GeCredits:    geCredits,
		FreeCredits:  int(freeRequiredCredits.Int()),
	}
}

func getCategoryTemplate(c model.CurriculumModel) string {

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

	//core and major template
	for _, g := range c.Curriculum.CoreAndMajorGroups {

		groupName := g.GroupName
		reqCourseList := []model.CourseDetailResponse{}
		reqCredit := 0
		elecCourseList := []model.CourseDetailResponse{}

		for _, c := range g.RequiredCourses {

			detail := getCourseDetail(c.CourseNo)
			reqCourseList = append(reqCourseList, model.CourseDetailResponse{
				CourseNo:   c.CourseNo,
				CourseName: detail.CourseDetail[0].CourseNameEN,
				GroupName:  groupName,
				IsPass:     false,
			})
			reqCredit += c.Credits
		}

		// for _, c := range g.ElectiveCourses {

		// 	detail := getCourseDetail(c.CourseNo)
		// 	elecCourseList = append(elecCourseList, model.CourseDetailResponse{
		// 		CourseNo:   c.CourseNo,
		// 		CourseName: detail.CourseDetail[0].CourseNameEN,
		// 		GroupName:  groupName,
		// 		IsPass:     false,
		// 	})

		// }

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

			detail := getCourseDetail(c.CourseNo)
			reqCourseList = append(reqCourseList, model.CourseDetailResponse{
				CourseNo:   c.CourseNo,
				CourseName: detail.CourseDetail[0].CourseNameEN,
				GroupName:  groupName,
				IsPass:     false,
			})
			reqCredit += c.Credits
		}

		// for _, c := range g.ElectiveCourses {

		// 	detail := getCourseDetail(c.CourseNo)
		// 	elecCourseList = append(elecCourseList, model.CourseDetailResponse{
		// 		CourseNo:   c.CourseNo,
		// 		CourseName: detail.CourseDetail[0].CourseNameEN,
		// 		GroupName:  groupName,
		// 		IsPass:     false,
		// 	})
		// }

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

	template := model.CategoryResponse{
		SummaryCredits:  0,
		RequiredCredits: curriculumRequiredCredits,
		CoreCategory:    coreCategory,
		MajorCategory:   majorCategory,
		GECategory:      geCategory,
		FreeCategory:    freeCategory,
	}

	t, err := json.Marshal(template)
	if err != nil {
		log.Fatalln("Error is : ", err)
	}

	return string(t)

}

func main() {
	e := echo.New()
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

		// CoreList := []model.CourseDetailResponse{}
		// CoreCredits := 0
		// MajorList := []model.CourseDetailResponse{}
		// MajorCredits := 0
		// GEList := []model.CourseDetailResponse{}
		// GECredits := 0
		// FreeList := []model.CourseDetailResponse{}
		// FreeCredits := 0

		// // studentID := c.QueryParam("studentID")
		// year := c.QueryParam("year")
		// curriculumProgram := c.QueryParam("curriculumProgram")
		// isCOOP := c.QueryParam("isCOOP")

		// cirriculum := getCirriculum(year, curriculumProgram, isCOOP)

		// transcript := readMockData("mockData1")
		// yearList := gjson.Get(transcript, "transcript.#.year")
		// for _, y := range yearList.Array() {
		// 	semester := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#`)
		// 	i := 1
		// 	for i < (int(semester.Int()) + 1) {
		// 		courseList := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#(semester==`+strconv.Itoa(i)+`).details`)
		// 		for _, c := range courseList.Array() {

		// 			code := gjson.Get(c.String(), "code")
		// 			course := gjson.Get(c.String(), "course")
		// 			grade := gjson.Get(c.String(), "grade")
		// 			credit := gjson.Get(c.String(), "credit")

		// 			if slices.Contains(PASS_GRADE, grade.String()) {
		// 				group := checkGroup(cirriculum, code.String())

		// 				if group == "Free" {
		// 					FreeCredits += int(credit.Int())
		// 					FreeList = append(FreeList, model.CourseDetailResponse{
		// 						CourseNo:   code.String(),
		// 						CourseName: course.String(),
		// 						GroupName:  group,
		// 						IsPass:     true,
		// 					})
		// 				} else if group == "Core" {
		// 					CoreCredits += int(credit.Int())
		// 					CoreList = append(CoreList, model.CourseDetailResponse{
		// 						CourseNo:   code.String(),
		// 						CourseName: course.String(),
		// 						GroupName:  group,
		// 						IsPass:     true,
		// 					})
		// 				} else if group == "Major Required" || group == "Major Elective" {
		// 					MajorCredits += int(credit.Int())
		// 					MajorList = append(MajorList, model.CourseDetailResponse{
		// 						CourseNo:   code.String(),
		// 						CourseName: course.String(),
		// 						GroupName:  group,
		// 						IsPass:     true,
		// 					})
		// 				} else {
		// 					GECredits += int(credit.Int())
		// 					GEList = append(GEList, model.CourseDetailResponse{
		// 						CourseNo:   code.String(),
		// 						CourseName: course.String(),
		// 						GroupName:  group,
		// 						IsPass:     true,
		// 					})
		// 				}
		// 			}
		// 		}

		// 		i++
		// 	}
		// }

		// // requiredCredits := gjson.Get(cirriculum, "curriculum.requiredCredits")
		// // credits := getRequiredCredits(cirriculum)
		// // sumCredits := GECredits + CoreCredits + MajorCredits + FreeCredits

		// response := model.CategoryResponse{
		// 	// SummaryCredits:  sumCredits,
		// 	// RequiredCredits: int(requiredCredits.Int()),
		// 	// CoreCategory: model.CategoryDetail{
		// 	// 	SummaryCredits:  CoreCredits,
		// 	// 	RequiredCredits: credits.CoreCredits,
		// 	// },
		// 	// MajorCategory: model.CategoryDetail{
		// 	// 	SummaryCredits:  MajorCredits,
		// 	// 	RequiredCredits: credits.MajorCredits,
		// 	// },
		// 	// GECategory: model.CategoryDetail{
		// 	// 	SummaryCredits:  GECredits,
		// 	// 	RequiredCredits: credits.GeCredits,
		// 	// },
		// 	// FreeCategory: model.CategoryDetail{
		// 	// 	SummaryCredits:  FreeCredits,
		// 	// 	RequiredCredits: credits.FreeCredits,
		// 	// },
		// }

		return c.JSON(http.StatusOK, nil)
	})

	e.GET("/v2/categoryView", func(c echo.Context) error {

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		cirriculumJSON := getCirriculumJSON(year, curriculumProgram, isCOOP)
		template := getCategoryTemplate(cirriculumJSON)

		summaryCredits := 0

		curriculumString := getCirriculum(year, curriculumProgram, isCOOP)
		transcript := readMockData("mockData1")
		yearList := gjson.Get(transcript, "transcript.#.year")
		for _, y := range yearList.Array() {
			semester := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#`)
			i := 1
			for i < (int(semester.Int()) + 1) {
				courseList := gjson.Get(transcript, `transcript.#(year=="`+y.String()+`").yearDetails.#(semester==`+strconv.Itoa(i)+`).details`)
				for _, c := range courseList.Array() {

					code := gjson.Get(c.String(), "code")
					course := gjson.Get(c.String(), "course")
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
									CourseNo:   code.String(),
									CourseName: course.String(),
									GroupName:  group,
									IsPass:     true,
								})
							} else {

								oldValue := gjson.Get(template, `freeCategory.#(groupName="Free Elective")`).String()
								categoryDetail := model.CategoryDetail{}
								err := json.Unmarshal([]byte(oldValue), &categoryDetail)
								if err != nil {
									return nil
								}

								courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
									CourseNo:   code.String(),
									CourseName: course.String(),
									GroupName:  group,
									IsPass:     true,
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
										CourseNo:   code.String(),
										CourseName: course.String(),
										GroupName:  group,
										IsPass:     true,
									})
								} else {

									oldValue := gjson.Get(template, `coreCategory.#(groupName="Core")`).String()
									categoryDetail := model.CategoryDetail{}
									err := json.Unmarshal([]byte(oldValue), &categoryDetail)
									if err != nil {
										return nil
									}

									courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
										CourseNo:   code.String(),
										CourseName: course.String(),
										GroupName:  group,
										IsPass:     true,
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
										CourseNo:   code.String(),
										CourseName: course.String(),
										GroupName:  group,
										IsPass:     true,
									})
								} else {

									oldValue := gjson.Get(template, `majorCategory.#(groupName="`+group+`")`).String()
									categoryDetail := model.CategoryDetail{}
									err := json.Unmarshal([]byte(oldValue), &categoryDetail)
									if err != nil {
										return nil
									}

									courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
										CourseNo:   code.String(),
										CourseName: course.String(),
										GroupName:  group,
										IsPass:     true,
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
								newCredit := oldCredit + credit
								template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").electiveCreditsGet`, newCredit)

								if gjson.Get(template, `geCategory.#(groupName="`+group+`").electiveCourseList.#`).Int() == 0 {
									courseList = append(courseList, model.CourseDetailResponse{
										CourseNo:   code.String(),
										CourseName: course.String(),
										GroupName:  group,
										IsPass:     true,
									})
								} else {

									oldValue := gjson.Get(template, `geCategory.#(groupName="`+group+`")`).String()
									categoryDetail := model.CategoryDetail{}
									err := json.Unmarshal([]byte(oldValue), &categoryDetail)
									if err != nil {
										return nil
									}

									courseList = append(categoryDetail.ElectiveCourseList, model.CourseDetailResponse{
										CourseNo:   code.String(),
										CourseName: course.String(),
										GroupName:  group,
										IsPass:     true,
									})
								}

								template, _ = sjson.Set(template, `geCategory.#(groupName="`+group+`").electiveCourseList`, courseList)
							}

						}
					}

					summaryCredits += int(credit)
				}

				i++
			}
		}
		log.Println(template)

		reponse := model.CategoryResponse{}
		err := json.Unmarshal([]byte(template), &reponse)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		reponse.SummaryCredits = summaryCredits

		return c.JSON(http.StatusOK, reponse)

	})

	e.Logger.Fatal(e.Start(":8080"))
}
