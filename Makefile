login:
	docker login registry.lab.madeof.glass

build:
	docker build -t registry.lab.madeof.glass/windowpane:v1.0.0 .
	docker push registry.lab.madeof.glass/windowpane:v1.0.0

run:
	docker run --rm -p 8080:8080 registry.lab.madeof.glass/windowpane:v1.0.0
