package logs

import (
	"ChatterinoLogExplorer/models"
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var ROOTPATH string = "C:/Users/USERPROFILENAME/AppData/Roaming/Chatterino2/Logs/Twitch/Channels"
var LOGFILE string = "./logs.txt"
var CDICT map[string][]string = make(map[string][]string)

func Search(formData models.FormData) string {

	//fmt.Printf("Performing Search on %s", formData)
	if err := sanitize(); err != nil {
		panic(err)
	}

	buildDict(formData.Terms)
	fileList := buildFileList(formData.Channel, formData.Duration)
	loadAllLogs(fileList)

	results := searchForTerms(formData)

	return results
}

func searchForTerms(formData models.FormData) string {

	cleanedLogs := ""

	f, _ := os.Open(LOGFILE)
	defer f.Close()

	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
	date := ""

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}
		// If its a Start/End Logging message discard
		if line[0] == '#' {
			continue
		}

		//If its a date 2025-09-17 set as current date
		if re.MatchString(line) {
			date = line
			continue
		}

		for _, kw := range formData.Terms {

			record, err := parseLine(line, date)
			if err != nil {
				continue
			}

			if validateLine(record, formData, kw) {
				CDICT[kw] = append(CDICT[kw], assembleLineData(record))
			}
		}
	}

	var headerParts []string
	var bodyParts []string
	for k, slice := range CDICT {

		numFoundLogs := len(slice)

		headerParts = append(headerParts, fmt.Sprintf("%s - %d", k, numFoundLogs))

		bpart := fmt.Sprintf("%s - %d\n", k, numFoundLogs) + strings.Join(CDICT[k], "")
		bodyParts = append(bodyParts, bpart)
	}

	header := "Found Instances: " + strings.Join(headerParts, ", ") + "\n"
	body := strings.Join(bodyParts, "\n"+strings.Repeat("_", 200)+"\n\n")
	cleanedLogs = header + "\n" + body

	return cleanedLogs
}

func parseLine(message string, date string) (models.Record, error) {
	var record models.Record

	timestamp := message[1:9]
	rest := message[11:]

	firstColon := strings.IndexByte(rest, ':')
	if firstColon == -1 {
		return record, fmt.Errorf("invalid message")
	}

	username := rest[:firstColon]
	msg := rest[firstColon+2:]

	record.Date = date + "-" + timestamp
	record.User = username
	record.Message = msg

	return record, nil
}

func validateLine(record models.Record, formData models.FormData, kw string) bool {

	if len(formData.Terms) == 1 && formData.Terms[0] == "" {
		return strings.EqualFold(formData.User, record.User)
	}

	if formData.User != "" {
		if strings.EqualFold(formData.User, record.User) {
			return false
		}
	}

	if formData.ExactMatch {
		for _, word := range strings.Fields(record.Message) {
			if word == kw {
				return true
			}
		}
		return false
	}

	if strings.Contains(record.Message, kw) && kw != "" {
		return true
	}

	return false
}

func assembleLineData(record models.Record) string {

	var b strings.Builder

	msgLength := len(record.Date) + 1 + len(record.User) + 2 + len(record.Message) + 1

	b.Grow(msgLength)
	b.WriteString(record.Date)
	b.WriteByte(' ')
	b.WriteString(record.User)
	b.WriteString(": ")
	b.WriteString(record.Message)
	b.WriteByte('\n')

	return b.String()

}

func loadAllLogs(fileList []string) {

	for _, fp := range fileList {

		data, _ := os.ReadFile(fp)
		f, _ := os.OpenFile(LOGFILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
		date := re.FindString(fp)
		dateString := date + "\n"
		f.Write([]byte(dateString))
		f.Write(data)
		f.Close()
	}
}

func buildFileList(channel string, timeSpan int) []string {
	var fileList []string

	for i := 0; i < timeSpan+1; i++ {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		fileName := channel + "-" + date + ".log"
		fp := filepath.Join(ROOTPATH, channel, fileName)
		fileList = append(fileList, fp)
	}

	return fileList
}

func sanitize() error {
	if err := os.Remove(LOGFILE); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	CDICT = make(map[string][]string)

	return nil
}

func buildDict(terms []string) {
	for _, term := range terms {
		CDICT[term] = []string{}
	}
}
