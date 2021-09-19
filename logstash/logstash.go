package logstash

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type PipelineRequest struct {
	UUID          string `json:"uuid"`
	PipelineInput string `json:"pipelineInput"`
	LogInput      string `json:"logInput"`
}

func ExecutePipeline(pr PipelineRequest, logger *zap.Logger) (string, error) {
	err := createDir(pr.UUID)
	if err != nil {
		return "", err
	}

	err = initializeLogstashConfig(pr.UUID)
	if err != nil {
		return "", err
	}

	// append input output for logstash pipeline
	pipelineString, err := appendInputOutput(pr.PipelineInput, pr.UUID)
	if err != nil {
		return "", err
	}

	// write to pipeline file
	pipelinePath := viper.GetString("config.logstash.pipeline_dir") + "/" + pr.UUID + "/pipeline.conf"
	err = writeToFile(pipelineString, pipelinePath)
	if err != nil {
		return "", err
	}

	// write to log file
	logPath := viper.GetString("config.logstash.pipeline_dir") + "/" + pr.UUID + "/sample.log"
	err = writeToFile(pr.LogInput, logPath)
	if err != nil {
		return "", err
	}

	settingsPath := viper.GetString("config.logstash.pipeline_dir") + "/" + pr.UUID

	output, err := execute(pipelinePath, settingsPath, logger)
	if err != nil {
		return "", err
	}

	return output, nil
}

func Version(logger *zap.Logger) (string, error) {
	cmd := exec.Command(viper.GetString("config.logstash.bin_path"), "--version")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(fmt.Sprint(err.Error() + ": " + string(stdout)))
	}

	return string(stdout), nil
}

func createDir(uuid string) error {
	p := viper.GetString("config.logstash.pipeline_dir") + "/" + uuid
	err := os.MkdirAll(p, 0777)
	if err != nil {
		return err
	}

	p = viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/data"
	err = os.MkdirAll(p, 0777)
	if err != nil {
		return err
	}

	p = viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/log"
	err = os.MkdirAll(p, 0777)
	if err != nil {
		return err
	}

	return nil
}

func initializeLogstashConfig(uuid string) error {
	type logstashConfig struct {
		LogstashData string
		LogstashLog  string
	}

	var lc logstashConfig
	lc.LogstashData = viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/data"
	lc.LogstashLog = viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/log"

	const confTemplate = `path.data: {{ .LogstashData }}
pipeline.ordered: auto
path.logs: {{ .LogstashLog }}`

	t := template.Must(template.New("config").Parse(confTemplate))

	var logstashConfigString bytes.Buffer

	err := t.Execute(&logstashConfigString, lc)
	if err != nil {
		return err
	}

	logstashConfigPath := viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/logstash.yml"
	writeToFile(logstashConfigString.String(), logstashConfigPath)

	return nil
}

func execute(pipelinePath string, settingsPath string, logger *zap.Logger) (string, error) {
	logger.Info("execute command",
		zap.String("cmd", fmt.Sprint(viper.GetString("config.logstash.bin_path"), "-f", pipelinePath, "--path.settings", settingsPath)))

	cmd := exec.Command(viper.GetString("config.logstash.bin_path"), "-f", pipelinePath, "--path.settings", settingsPath)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(fmt.Sprint(err.Error() + ": " + string(stdout)))
	}

	return string(stdout), nil
}

func appendInputOutput(pipeline string, uuid string) (string, error) {
	logPath := viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/sample.log"

	const inTemplate = `input {
  file {
    path => "{{ . }}"
    start_position => "beginning"
    sincedb_path => "/dev/null"
    mode => "read"
    exit_after_read => true
  }
}

`
	t := template.Must(template.New("input").Parse(inTemplate))

	var inString bytes.Buffer

	err := t.Execute(&inString, logPath)
	if err != nil {
		return "", err
	}

	outString := `

output {
  stdout { codec => rubydebug }
}`

	return inString.String() + pipeline + outString, nil
}

func writeToFile(input string, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	dw := bufio.NewWriter(file)
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		dw.WriteString(line + "\n")
	}

	dw.Flush()
	file.Close()

	return nil
}

func readPipeline(pipeline string, uuid string) error {
	fmt.Println("read pipeline")

	fileName := viper.GetString("config.logstash.pipeline_dir") + "/" + uuid + "/" + "pipeline.conf"

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	dw := bufio.NewWriter(file)

	lines := strings.Split(pipeline, "\n")
	startPrint := false

	for i, line := range lines {
		if equal := strings.Index(line, "input {"); equal >= 0 {
			startPrint = false
		}

		if equal := strings.Index(line, "filter {"); equal >= 0 {
			startPrint = true
		}

		if equal := strings.Index(line, "output {"); equal >= 0 {
			startPrint = false
		}

		if startPrint {
			if viper.GetBool("config.log.debug") {
				fmt.Println(i, line)
			}

			dw.WriteString(line + "\n")
		}
	}

	dw.Flush()
	file.Close()

	return nil
}
