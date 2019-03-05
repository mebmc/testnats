#!/bin/bash
siege -c5000 -t1s --content-type "application/json"	\
	'http://dev.local:8080/api/submit POST {"id": "1", "first_name": "alice", "last_name": "Smith", "email": "alice@example.com", "guess": "a"}'
