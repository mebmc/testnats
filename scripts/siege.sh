#!/bin/bash
siege -c50 -t15m --content-type "application/json"	\
	'http://dev.local:8080/api/submit POST {"id": "1", "first_name": "alice", "last_name": "Smith", "email": "alice@example.com", "guess": "a"}'
