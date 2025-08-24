# brew install node-sass or npm install sass -g
export rootDir="public/assets/frontend/"

sass $rootDir/css/custom.css:$rootDir/css/custom.min.css\
 $rootDir/blogdefault/css/style.css:$rootDir/blogdefault/css/style.min.css\
 $rootDir/blogdefault/css/dark.css:$rootDir/blogdefault/css/dark.min.css\
 --style=compressed --no-source-map --watch

# $rootDir/css/frontend.scss:$rootDir/css/frontend.min.css\ #bug
# ~/minifyCSS.sh frontend.css frontend.min.css