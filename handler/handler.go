package handler

import (
	"encoding/json"
	"net/http"

	"github.com/misbahulard/logstash-pipeline-tester/logstash"
	"go.uber.org/zap"
)

type HealthCheckResponse struct {
	Alive  bool `json:"alive"`
	Status struct {
		Logstash bool `json:"logstash"`
	} `json:"status"`
	Message struct {
		Logstash string `json:"logstash"`
	} `json:"message"`
}

type PipelineResponse struct {
	Data struct {
		Output string `json:"output"`
	} `json:"data"`
	Message string `json:"message"`
}

type DefaultResponse struct {
	Message string `json:"message"`
}

func errorResponse(w http.ResponseWriter, status int, msg string) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	var ds DefaultResponse
	ds.Message = msg

	resp, err := json.Marshal(ds)
	if err != nil {
		return err
	}

	w.Write(resp)
	return nil
}

func PipelineHandler(logger *zap.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.Error("an error occured",
				zap.String("stacktrace", err.Error()))
			errorResponse(w, http.StatusInternalServerError, "an error occured, stacktrace: "+err.Error())
			return
		}

		var pr logstash.PipelineRequest

		pr.UUID = r.FormValue("uuid")
		pr.PipelineInput = r.FormValue("pipeline_input")
		pr.LogInput = r.FormValue("log_input")

		output, err := logstash.ExecutePipeline(pr, logger)
		if err != nil {
			logger.Error("an error occured",
				zap.String("stacktrace", err.Error()))

			var pipelineResponse PipelineResponse
			pipelineResponse.Data.Output = err.Error()
			pipelineResponse.Message = "an error occured"

			resp, err := json.Marshal(&pipelineResponse)
			if err != nil {
				logger.Error("an error occured",
					zap.String("stacktrace", err.Error()))
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
			return
		}

		var pipelineResponse PipelineResponse
		pipelineResponse.Data.Output = output
		pipelineResponse.Message = "logstash pipeline has been executed successfully"

		resp, err := json.Marshal(&pipelineResponse)
		if err != nil {
			logger.Error("an error occured",
				zap.String("stacktrace", err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})
}

func HealthCheckHandler(logger *zap.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var health HealthCheckResponse

		stdout, err := logstash.Version(logger)
		if err != nil {
			logger.Error("an error occured",
				zap.String("stacktrace", err.Error()))
			health.Message.Logstash = err.Error()
		} else {
			logger.Info("logstash version",
				zap.String("stdout", stdout))
			health.Message.Logstash = stdout
		}

		if err == nil {
			health.Status.Logstash = true
		}

		if health.Status.Logstash {
			health.Alive = true
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		resp, err := json.Marshal(&health)
		if err != nil {
			logger.Error("an error occured",
				zap.String("stacktrace", err.Error()))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	})
}

func RootHandler(logger *zap.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("it's works!"))
	})
}
