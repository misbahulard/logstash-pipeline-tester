build:
	@echo "Start building nuxtjs app"
	@cd web/logstash-pipeline-tester; yarn build
	@echo "Build success"
	@echo ""
	@echo "Start building golang app"
	@go build -o logstash-pipeline-tester
	@chmod +x logstash-pipeline-tester
	@echo "Build app success -> 'logstash-pipeline-tester'"
	@echo ""