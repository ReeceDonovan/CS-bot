package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// Req constructs a request and returns status code and response body
func Req(method, slug, token string, body []byte) (int, []byte) {

	r, err := http.NewRequest(method,
		fmt.Sprintf("%s/%s", viper.GetString("canvas.domain"), slug),
		bytes.NewReader(body),
	)
	if err != nil {
		log.Println("Error making http request <", method, ">", slug, "\n", err)
		return 0, []byte{}
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	if method == "POST" {
		r.Header.Add("Content-Type", "application/json")
	}
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("Error sending http request <", method, ">", slug, "\n", err)
		return 0, []byte{}
	}
	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Couldn't read response \n", err)
		return 0, []byte{}
	}
	return resp.StatusCode, bd
}

// TODO: Might move everything from here down to different file for tidyness if we need other query commands, but ah sur. Also gotta change default struct variable names from the autogen

func QueryAssignments() (parsedData CourseAssignment) {

	_, res := Req("POST", "/api/graphql",
		fmt.Sprintf("%s", viper.GetString("canvas.token")),
		[]byte(`
			{"query": "query CourseAssignments {
				allCourses {
				  _id
				  name
				  state
				  courseCode
				  term {
					_id
					name
				  }
				  assignmentsConnection {
					  nodes {
						_id
						name
						dueAt
						htmlUrl
					  }
				  }
				  enrollmentsConnection {
					edges {
					  node {
						id
					  }
					}
				  }
				}
			  }
			  "}`))

	jsonErr := json.Unmarshal(res, &parsedData)
	if jsonErr != nil {
		log.Println("Error parsing response\n", jsonErr)
	}

	for _, course := range parsedData.Data.AllCourses {
		if course.Term.ID != "44" || course.EnrollmentsConnection == nil {
			continue
		}
		_, res := Req("GET", "/api/v1/courses/"+course.ID+"/assignments?include[]=submission&include[]=score_statistics", fmt.Sprintf("%s", viper.GetString("canvas.token")), nil)

		assStat := ScoreRaw{}

		jsonErr := json.Unmarshal(res, &assStat)
		if jsonErr != nil {
			log.Println("Error parsing response\n", jsonErr)
		}

		for x, _ := range course.AssignmentsConnection.Nodes {
			course.AssignmentsConnection.Nodes[x].ScoreStatistics.Min = assStat[x].ScoreStatistics.Min
			course.AssignmentsConnection.Nodes[x].ScoreStatistics.Mean = assStat[x].ScoreStatistics.Mean
			course.AssignmentsConnection.Nodes[x].ScoreStatistics.Max = assStat[x].ScoreStatistics.Max
		}
	}
	return parsedData
}