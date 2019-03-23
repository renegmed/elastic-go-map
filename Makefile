
.PHONY: cmd
cmd:
	go install ./match

.PHONY: update
update:	
	match producer delete index -i producer
	match producer index file -f producer.json
	match consumer delete index -i consumer
	match consumer index file -f consumer.json
	curl -XGET 'localhost:9200/_cat/indices?v&pretty'


build: 
	if [ -a ./web.exe ]; then  rm ./web.exe; fi;   # remove main if it exists 
	go build -o ./web.exe
	./web.exe