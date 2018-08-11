## mac-unison-fsmonitor

The Unison has file system change monitoring (in the Unison config, `repeat = watch`) through a process
that runs and monitors file system activity. The protocal is defined in https://github.com/bcpierce00/unison/blob/master/src/fswatch.ml#L19. The Unison distribution includes file system monitors for Window and Linux.

There is a Python implementation of the file system monitor, https://github.com/hnsl/unox which I have used for quite some time now, but it stopped working on my MacBook Pro running High Sierra on an APFS file system.

I thought it would be fun to implement the watcher in Go. The only dependency is on https://github.com/fsnotify/fsevents but the version in the vendor directory has some bug fixes applied that have not yet made it into master yet (https://github.com/fsnotify/fsevents/pull/38 and https://github.com/fsnotify/fsevents/pull/39).
