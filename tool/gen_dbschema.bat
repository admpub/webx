go get github.com/webx-top/db
go install github.com/webx-top/db/cmd/dbgenerator
dbgenerator -d nging -p root -o ../application/dbschema -match "^official_(ad|common|customer|page|short_url)($|_)" -backup "../application/library/setup/install.sql"

pause