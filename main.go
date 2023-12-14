package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"NuaikitTracking_BE.com/model"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
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

func checkGroup(cirriculum string, courseNo string) string {

	groupList := gjson.Get(cirriculum, "curriculum.geGroups.#.groupName")
	for _, groupName := range groupList.Array() {

		queryReqCourse := `curriculum.geGroups.#(groupName=="` + groupName.String() + `").requiredCourses.#(courseNo=="` + courseNo + `")`
		valueReqCourse := gjson.Get(cirriculum, queryReqCourse)

		if valueReqCourse.Exists() {
			return groupName.String()
		}

		queryElecCourse := `curriculum.geGroups.#(groupName=="` + groupName.String() + `").electiveCourses.#(courseNo=="` + courseNo + `")`
		valueElecCourse := gjson.Get(cirriculum, queryElecCourse)

		if valueElecCourse.Exists() {
			return groupName.String()
		}
	}

	groupList = gjson.Get(cirriculum, "curriculum.coreAndMajorGroups.#.groupName")
	for _, groupName := range groupList.Array() {

		queryReqCourse := `curriculum.coreAndMajorGroups.#(groupName=="` + groupName.String() + `").requiredCourses.#(courseNo=="` + courseNo + `")`
		valueReqCourse := gjson.Get(cirriculum, queryReqCourse)

		if valueReqCourse.Exists() {
			return groupName.String()
		}

		queryElecCourse := `curriculum.coreAndMajorGroups.#(groupName=="` + groupName.String() + `").electiveCourses.#(courseNo=="` + courseNo + `")`
		valueElecCourse := gjson.Get(cirriculum, queryElecCourse)

		if valueElecCourse.Exists() {
			return groupName.String()
		}
	}

	return "Free"
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

		CoreList := []model.CourseDetailResponse{}
		CoreCredits := 0
		MajorList := []model.CourseDetailResponse{}
		MajorCredits := 0
		GEList := []model.CourseDetailResponse{}
		GECredits := 0
		FreeList := []model.CourseDetailResponse{}
		FreeCredits := 0

		// studentID := c.QueryParam("studentID")
		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		cirriculum := getCirriculum(year, curriculumProgram, isCOOP)

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
					credit := gjson.Get(c.String(), "credit")

					if slices.Contains(PASS_GRADE, grade.String()) {
						group := checkGroup(cirriculum, code.String())

						if group == "Free" {
							FreeCredits += int(credit.Int())
							FreeList = append(FreeList, model.CourseDetailResponse{
								CourseNo:   code.String(),
								CourseName: course.String(),
								GroupName:  group,
								IsPass:     true,
							})
						} else if group == "Core" {
							CoreCredits += int(credit.Int())
							CoreList = append(CoreList, model.CourseDetailResponse{
								CourseNo:   code.String(),
								CourseName: course.String(),
								GroupName:  group,
								IsPass:     true,
							})
						} else if group == "Major Required" || group == "Major Elective" {
							MajorCredits += int(credit.Int())
							MajorList = append(MajorList, model.CourseDetailResponse{
								CourseNo:   code.String(),
								CourseName: course.String(),
								GroupName:  group,
								IsPass:     true,
							})
						} else {
							GECredits += int(credit.Int())
							GEList = append(GEList, model.CourseDetailResponse{
								CourseNo:   code.String(),
								CourseName: course.String(),
								GroupName:  group,
								IsPass:     true,
							})
						}
					}
				}

				i++
			}
		}

		requiredCredits := gjson.Get(cirriculum, "curriculum.requiredCredits")
		credits := getRequiredCredits(cirriculum)
		sumCredits := GECredits + CoreCredits + MajorCredits + FreeCredits

		response := model.CategoryResponse{
			SummaryCredits:  sumCredits,
			RequiredCredits: int(requiredCredits.Int()),
			CoreCategory: model.CategoryDetail{
				SummaryCredits:  CoreCredits,
				RequiredCredits: credits.CoreCredits,
				CourseList:      CoreList,
			},
			MajorCategory: model.CategoryDetail{
				SummaryCredits:  MajorCredits,
				RequiredCredits: credits.MajorCredits,
				CourseList:      MajorList,
			},
			GECategory: model.CategoryDetail{
				SummaryCredits:  GECredits,
				RequiredCredits: credits.GeCredits,
				CourseList:      GEList,
			},
			FreeCategory: model.CategoryDetail{
				SummaryCredits:  FreeCredits,
				RequiredCredits: credits.FreeCredits,
				CourseList:      FreeList,
			},
		}

		return c.JSON(http.StatusOK, response)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
