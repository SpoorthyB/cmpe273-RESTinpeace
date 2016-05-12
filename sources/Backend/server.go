package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/drone/routes"
	"github.com/mkilling/goejdb"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var myEJDB *goejdb.Ejdb
var myEJDBColl *goejdb.EjColl

func main() {

	go createService("3005")
	go checkStatus()

	var input string
	fmt.Scanln(&input)
}

func checkStatus() {
	for {
		if nil != myEJDBColl {
			res, _ := myEJDBColl.FindOne("students")

			if nil != res {
				log.Println("checkStatus_found student information")
				var profileMap interface{}
				bson.Unmarshal(*res, &profileMap)
				bJSONString, _ := json.Marshal(profileMap)
				requestParser, _ := simplejson.NewJson(bJSONString)

				fmt.Println("array size ", len(requestParser.Get("students").MustArray()))
				now_time := time.Now().Unix()

				for i := 0; i < len(requestParser.Get("students").MustArray()); i++ {
					stime := requestParser.Get("students").GetIndex(i).Get("time").MustString()
					t, _ := strconv.Atoi(stime)
					id := requestParser.Get("students").GetIndex(i).Get("id").MustString()
					attended := requestParser.Get("students").GetIndex(i).Get("attended").MustString()
					if int(now_time)-t > 10 && len(attended) != 0 {
						myEJDBColl.Update(getUpdateAttendedQueryString(id, "no"))
					}
				}

				//fmt.Println(parser.String())
			} else {
				fmt.Println("checkStatus_not found")
			}
		}
		time.Sleep(time.Second * 10)

	}
}
func createService(port string) {

	dbFileName := "server_" + port + ".db"

	log.Println(dbFileName)
	thisEJDB, err := goejdb.Open(dbFileName, goejdb.JBOWRITER|goejdb.JBOCREAT|goejdb.JBOTRUNC)
	//thisEJDB, err := goejdb.Open(dbFileName, goejdb.JBOWRITER|goejdb.JBOCREAT)
	if err != nil {
		log.Println("goejdb.Open fail: ", err)
		return
	}

	myEJDB = thisEJDB
	thisEJDBColl, _ := myEJDB.CreateColl("coll", nil)
	defer myEJDB.Close()

	myEJDBColl = thisEJDBColl

	log.Println("server running")
	mux := routes.New()

	mux.Put("/:id/:name", putRegisterAPI)
	mux.Put("/update/:id/attended", putAttendedAPI)
	mux.Put("/update/:id/left", putLeftAPI)
	mux.Get("/", getAllAPI)

	// use localhost:port/kill to terminate service
	mux.Put("/kill", kill)

	http.Handle("/", mux)

	http.ListenAndServe(":"+port, nil)
}

func getAllAPI(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllAPI")

	res, _ := myEJDBColl.FindOne("students")

	if nil != res {
		log.Println("found student information")
		var profileMap interface{}
		bson.Unmarshal(*res, &profileMap)
		bJSONString, _ := json.Marshal(profileMap)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bJSONString)
		return
	} else {
		log.Println("not found")
		w.WriteHeader(http.StatusNotFound)
	}
}

func putRegisterAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	name := params.Get(":name")

	log.Println("putRegister: id:" + id + " name: " + name)

	w.WriteHeader(putRegister(id, name))
}

func putRegister(id, name string) int {
	res, _ := myEJDBColl.Find("students")

	if 0 != len(res) {
		resID, _ := myEJDBColl.Find(getElemMatchString(id))

		if 0 != len(resID) {
			log.Println("User exist")
			return http.StatusForbidden
		}

		log.Println("DB has data, update it")
		log.Println("String: " + getUpdateQueryString(id, name))

		myEJDBColl.Update(getUpdateQueryString(id, name))
	} else { // Add new profile
		log.Println("DB empty, create it")

		bodyParser, _ := simplejson.NewJson([]byte(`{"students":[]}`))
		profileMap, _ := bodyParser.Map()
		bProfileMap, _ := bson.Marshal(profileMap)
		myEJDBColl.SaveBson(bProfileMap)
		myEJDBColl.Update(getUpdateQueryString(id, name))
	}

	printAllProfiles()

	return http.StatusNoContent
}

func putAttendedAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")

	log.Println("putAttendedAPI: id:" + id)

	w.WriteHeader(putAttended(id))
}

func putAttended(id string) int {
	res, _ := myEJDBColl.Find(getElemMatchString(id))

	if 0 != len(res) {
		log.Println("DB has data, update it")
		t := getTime()
		log.Println("String: " + getUpdateAttendedQueryString(id, "yes"))

		myEJDBColl.Update(getUpdateAttendedQueryString(id, "yes"))
		myEJDBColl.Update(getUpdateTimeStampQueryString(id, t))
		myEJDBColl.Update(getUpdateTimeQueryString(id, strconv.Itoa(int(time.Now().Unix()))))

		return http.StatusNoContent
	} else {
		return http.StatusNotFound
	}
}

func putLeftAPI(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id := params.Get(":id")
	log.Println("putLeftAPI: id:" + id)

	w.WriteHeader(putLeft(id))
}

func putLeft(id string) int {
	res, _ := myEJDBColl.Find(getElemMatchString(id))

	if 0 != len(res) {
		log.Println("DB has data, update it")
		t := getTime()
		log.Println("String: " + getUpdateAttendedQueryString(id, "no"))

		myEJDBColl.Update(getUpdateAttendedQueryString(id, "no"))
		myEJDBColl.Update(getUpdateTimeStampQueryString(id, t))
		myEJDBColl.Update(getUpdateTimeQueryString(id, strconv.Itoa(int(time.Now().Unix()))))

		return http.StatusNoContent
	}

	return http.StatusNotFound
}

func getTime() string {
	t := time.Now().Add(-7 * time.Hour)
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second())
}

func printAllProfiles() {
	res, _ := myEJDBColl.Find("")

	log.Println("\n******************\nRecords found: %d", len(res))

	// Now print the result set records
	for _, bs := range res {
		var m map[string]interface{}
		bson.Unmarshal(bs, &m)
		log.Println(m)
		log.Println("----------------------")
	}

	log.Println("******************")
}

func kill(w http.ResponseWriter, r *http.Request) {
	log.Println("server killed")
	os.Exit(0)
}

func getUpdateQueryString(id, name string) string {
	return `{"$addToSetAll":{"students":[{"id":"` + id + `", "name":"` + name + `", "attended":"", "time_stamp":"", "time":""}]}}`
}

func getUpdateAttendedQueryString(id, attended string) string {
	return `{"students":{"$elemMatch":{"id":"` + id + `"}},"$upsert":{"students.$.attended":"` + attended + `"}}`
}

func getUpdateTimeStampQueryString(id, time string) string {
	return `{"students":{"$elemMatch":{"id":"` + id + `"}},"$upsert":{"students.$.time_stamp":"` + time + `"}}`
}

func getUpdateTimeQueryString(id, time string) string {
	return `{"students":{"$elemMatch":{"id":"` + id + `"}},"$upsert":{"students.$.time":"` + time + `"}}`
}

func getElemMatchString(id string) string {
	return `{"students":{"$elemMatch":{"id":"` + id + `"}}}`
}
