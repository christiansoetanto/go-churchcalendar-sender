how to run:
replace the omitted env with actual value.
run cmd/main.go
curl http://localhost:8080
enjoy

file explanation:
i was using pubsub and its zamn hard (read: i dont understand how) to test locally, hence the main.go
but i migrated to HTTPFunction because its easier (read: i was able) to execute it in local, hence the http_functions
from now on, main.go is used for new feature while the actual cloud functions is in http_functions package