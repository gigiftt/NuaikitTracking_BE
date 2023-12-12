package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/tidwall/gjson"
)

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

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/core-major/elective", func(c echo.Context) error {
		groupName := c.QueryParam("groupName")

		year := c.QueryParam("year")
		curriculumProgram := c.QueryParam("curriculumProgram")
		isCOOP := c.QueryParam("isCOOP")

		cirriculum := getCirriculum(year, curriculumProgram, isCOOP)
		query := `curriculum.coreAndMajorGroups.#(groupName=="` + groupName + `").electiveCourses`
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

	e.GET("/viewTerm", func(c echo.Context) error {
		// res:= getCirriculum(2563, "CPE", false)
		// log.Println((res))
		return c.JSON(http.StatusOK, nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
