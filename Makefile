build-js:
	cd assets/src && npm run build

build-go:
	go build -o stonk-watcher.tmp stonk-watcher

build: build-js build-go

deploy: build
	ssh root@143.198.220.97 "(killall stonk-watcher.tmp || echo 'no process found')"
	scp stonk-watcher.tmp root@143.198.220.97:/root/go
	scp morningstarKey.tmp.json root@143.198.220.97:/root/go
	ssh root@143.198.220.97 "cd ~/go && (./stonk-watcher.tmp > ./stock-watcher.log 2>&1 &)"
