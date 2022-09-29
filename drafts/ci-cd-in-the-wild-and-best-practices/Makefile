.PHONY: images

images: design.pu
	plantuml -tpng design.pu
	mv *.png images/

build_watcher:
	bash -c "while true; do inotifywait -e modify design.pu; make; pkill -SIGHUP mupdf; done"
