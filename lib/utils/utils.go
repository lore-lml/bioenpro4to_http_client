package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const(
	Authenticate    = "id-manager/authenticate"
	ValidCredential = "id-manager/is-credential-valid"
	DailyChannel    = "channel-manager/daily-channel"
)

type Method string
const (
	Get  Method = "GET"
	Post        = "POST"
)

type Category string
const (
	Trucks Category = "trucks"
	Scales          = "weighing_scales"
	BioCells 		= "biocells"
)

type Credential []byte

const dateLayout = "02/01/2006 03:04:05 PM"

func CheckDateFormat(date string) bool{
	regex, err := regexp.Compile("^([0-2][0-9]|(3)[0-1])(/)(((0)[0-9])|((1)[0-2]))(/)\\d{4}$")
	if err != nil{
		return false
	}
	return regex.MatchString(date)
}

func DateToTimestamp(date string) (int64, error){
	newDate := strings.ReplaceAll(date, "-", "/")
	if !CheckDateFormat(newDate){
		return -1, errors.New("date should be in the format dd-mm-yyyy or dd/mm/yyyy")
	}
	snapshot := fmt.Sprintf("%s 00:00:00 AM", newDate)
	t, err := time.Parse(dateLayout, snapshot)
	if err != nil {
		return -1, errors.New("date should be in the format dd-mm-yyyy or dd/mm/yyyy")
	}
	return t.Unix(), nil
}

