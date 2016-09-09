package main

import(
	"time"
	"strconv"
)

func GetLogValue(log_number string)(string,int64){
	timestamp := time.Now().Unix()
        logvalue = strconv.FormatInt(timestamp, 10) + log_number
	return logvalue,timestamp
}

func AddLogValue(log_number,logvalue string)(string,int64){
	timestamp := time.Now().Unix()
	logvalue = logvalue + "|" + strconv.FormatInt(timestamp, 10) + log_number
	return logvalue,timestamp
}

func GetLogKey(servicename,hostname string,timestamp int64)string{
	logkey = "cmha/service/" + servicename + "/log/" + hostname + "/mha-handlers/" + strconv.FormatInt(timestamp, 10)
	return logkey
}
