## FilePlexer
This tool watches over a target folder where users with no access to NFS upload folders to.
Then it handles the responsability of transferring newly added content to the NFS.

### How it works:
This tool relies on fsnotify events to track what's happening to our target folder.
If a behavior pattern occurs that implies a file have been added, the main function launches a go coroutine to handle the transfer.

