package main

import (
	"fmt"
    "net/http"
    "log"
    "time"
	"regexp"
	"math/rand"
	"strconv"
	"strings"
	"io/ioutil"
)


type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pat *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pat, handler})
}

func (h *RegexpHandler) HandleFunc (instr string, handler func(http.ResponseWriter, *http.Request)) {
    pat := regexp.MustCompile(instr)
	h.routes = append(h.routes, &route{pat, http.HandlerFunc(handler)})
}

//Lookup for the pattern and use the methods
func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path){
			route.handler.ServeHTTP(w,r)
			return
		} 
	}
	http.NotFound(w, r)
	
}

//Serve the main index.html
func getindexfile(w http.ResponseWriter, r *http.Request) {
    	http.ServeFile(w, r, "./index.html")
}

//Serve CSS/.js files
func getassets(w http.ResponseWriter, r *http.Request) {
   	http.ServeFile(w, r, r.URL.Path[1:])
}

//Read "stat" file to get the CPU usage
func getcpu() (idle uint64, total_cpu uint64) {
	
	stat_data,err := ioutil.ReadFile("/proc/stat")
    
	if err != nil {
        return 
	}
	
	lines := strings.Split(string(stat_data), "\n")
	for _,line := range(lines) {
	    fields := strings.Fields(line)
        if fields[0] == "cpu" {
			numFields := len(fields)
			for cnt := 1; cnt< numFields; cnt++ {
				data,err := strconv.ParseUint(fields[cnt], 10, 64)
				if err != nil {
					fmt.Println("Err Parsing")
				}
				total_cpu += data
				if cnt == 4 {
					idle = data
				}
			}
			return 
		}
	}
	return	
}
//Calculate CPU usage
func getCPUusage(w http.ResponseWriter, r *http.Request) {
	
	idle1, total1 := getcpu()
	time.Sleep(3 * time.Second)
	idle2, total2 := getcpu()
		
	idleTicks := float64(idle2 - idle1)
    totalTicks := float64(total2 - total1)
	
	if totalTicks != 0 {
		cpuUsage := (totalTicks - idleTicks) / totalTicks
		w.Write([]byte(strconv.Itoa(int(cpuUsage))))
	} else { //Test method if data is not able to read from "stat" file
		w.Write([]byte(strconv.Itoa(rand.Intn(100))))
	}			
}


func main() {
		
    reHandler := new(RegexpHandler)
	reHandler.HandleFunc("/$", getindexfile) // set router
	reHandler.HandleFunc("/getCPUusage$", getCPUusage);
	
	reHandler.HandleFunc(".*.[js|css|png|eof|svg|ttf|woff]", getassets)
		
    err := http.ListenAndServe(":9092", reHandler) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}