# CowFS
A pure go copy-on-write FUSE file system. The file mount can be initialised with files from the host
file system. These files will appear as if copied already inside the file mount. Any reads will be 
directed to the host file system, however writes will cause a sideways copy to a temporary directory. 

This can improve performance massively for access patterns that are mostly read-only, however writes
are still permitted on occasion. 

# Origin
This was designed for use in [please](https://github.com/thought-machine/please) for managing our 
build directories. The build directories are hermetic directories where we move source files into
to restrict the files that build actions can access. For performance reasons, we symlink these files
to avoid excessive disk writes. This breaks the hermetic builds as build actions can modify the 
source tree. CowFS aims to rectify this by copying source files sideways when build actions attempt 
to modify them. 
