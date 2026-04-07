build:
	go build -o seedbank
cp:
	cp seedbank ~/.local/bin/
	
install: build cp