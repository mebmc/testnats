#!/bin/bash
http POST dev.local:8080/api/submit	\
	id=1				\
	first_name=alice		\
	last_name=smith			\
	email=alice@example.com		\
	guess=d
