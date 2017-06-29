# squirrelchopper

Yet another tile server...why? 

There are 3 main goals for SC that differentiate it from other (maybe not all) tile servers floating around:

### Focused
SC's goal is to provide a lighweight, portable, offline, basemap capability. To put it simply: it serves Vector tiles quickly.
No rendering, interpretation, security or other fluff. Just read-only access to open data. 

### Fast
SC embraces new standards in order to serve out tiles as fast as possible.  This includes:
- using HTTP2 with server push to predictively send 
tiles to the client before they are needed
- agressive in-memory caching memory storage of commonly used tiles (Zoom 1-8) + additional configurable LRU

### Small 
Thanks to Golang we end up with a single 11mb binary and a 186mb Alpine docker image which includes 86mb of data and sample apps. It provides
everything needed for a self-contained offline basemap. There is plenty to fat left to trim so expect the size to keep shrinking. At 
the moment this still uses glibc instead of Musl (native to Alpine) and provides multiple versions of demos which could all be optimized. 
SC has a small memory footprint as well -- using about 100mb of ram even with a full cache. 

### What's up with the name?
I've got a garden-wind-powered-metal-thing in my back yard and the squirrels love to climb on it. It rotates and looks like its going to 
chop them up -- hence the name squirrelchopper.  Anyways, watching it while drinking coffee one day I was struck at how efficiently it moved
and thought it would be a good name for what this is trying to accomplish

### References, projects and data this leverages

- https://github.com/lukasmartinelli/osm-liberty
- https://openmaptiles.org
- https://github.com/mapbox/mapbox-gl-js 
- https://github.com/elazarl/go-bindata-assetfs
